package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/anna-ssg/anna/v3/cmd/anna"
	"github.com/anna-ssg/anna/v3/pkg/helpers"
	"github.com/anna-ssg/anna/v3/pkg/logger"
	"github.com/spf13/cobra"
)

// CommitHash can be set at build time with -ldflags "-X main.CommitHash=<hash-or-tag>".
// Example: go build -ldflags "-X main.CommitHash=1a2b3c4"
var CommitHash string

func main() {
	var addr string
	var prof bool
	var renderDrafts bool
	var serve bool
	var webconsole bool
	var version bool
	var siteDirPath string

	Version := "v3.0.0" // to be set at build time $(git describe --tags)

	rootCmd := &cobra.Command{
		Use:   "anna",
		Short: "Static Site Generator",
		Run: func(cmd *cobra.Command, args []string) {
			siteDirPath = path.Clean(siteDirPath) + "/"

			// If the binary was built with an embedded commit hash, expose it to
			// the rest of the process via ANNA_COMMIT so non-CLI flows can use it.
			if CommitHash != "" {
				os.Setenv("ANNA_COMMIT", CommitHash)
			}
			annaCmd := anna.Cmd{
				RenderDrafts: renderDrafts,
				Addr:         addr,
				LiveReload:   serve,
				SiteDirPath:  siteDirPath,
				ErrorLogger:  logger.New(os.Stderr),
				InfoLogger:   logger.New(os.Stderr),
			}

			if serve {
				annaCmd.StartLiveReload(siteDirPath)
			}

			built := false

			if prof {
				startTime := time.Now()
				count := annaCmd.VanillaRender(siteDirPath)
				elapsedTime := time.Since(startTime)
				annaCmd.PrintStats(elapsedTime)
				annaCmd.InfoLogger.Printf("Built %d pages in %s\n", count, elapsedTime)
				built = true
			}

			if version {
				ver := Version
				if CommitHash != "" {
					ver = fmt.Sprintf("%s (commit: %s)", Version, CommitHash)
				}
				annaCmd.InfoLogger.Println("Current version:", ver)
			}

			if webconsole {
				server := anna.NewWizardServer(":8080")
				go server.Start()
				<-anna.FormSubmittedCh // wait for response
				if err := server.Stop(); err != nil {
					annaCmd.InfoLogger.Println(err)
				}
				annaCmd.StartLiveReload(siteDirPath)
			}

			if !built {
				startTime := time.Now()
				count := annaCmd.VanillaRender(siteDirPath)
				elapsedTime := time.Since(startTime)
				annaCmd.InfoLogger.Printf("Built %d pages in %s\n", count, elapsedTime)
			}
		},
	}

	rootCmd.Flags().StringVarP(&addr, "addr", "a", "localhost:8000", "specify address over which rendered content is served")
	rootCmd.Flags().BoolVarP(&renderDrafts, "draft", "d", false, "renders draft posts")
	// Do not set default values for string flags
	rootCmd.Flags().StringVarP(&siteDirPath, "path", "p", "site", "specify the specific site directory to render")
	rootCmd.Flags().BoolVar(&prof, "prof", false, "enable profiling")
	rootCmd.Flags().BoolVarP(&serve, "serve", "s", false, "serve the rendered site and watch for file updates")
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "prints current version number")
	// rootCmd.Flags().BoolVarP(&webconsole, "webconsole", "w", false, "wizard to setup anna")

	// bootstrap subcommand: download and extract `site/` from the upstream archive
	var bsYes bool
	var bsURL string
	bootstrapCmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Download and extract the default site/ layout from upstream",
		Run: func(cmd *cobra.Command, args []string) {
			info := logger.New(os.Stderr)

			// If a build-time CommitHash was embedded prefer that over --url
			if CommitHash != "" {
				trim := strings.TrimSpace(CommitHash)
				// If the embedded value already looks like a URL, use it directly.
				if strings.HasPrefix(trim, "http://") || strings.HasPrefix(trim, "https://") {
					bsURL = trim
					info.Printf("Bootstrapping from embedded URL -> %s\n", bsURL)
				} else {
					bsURL = fmt.Sprintf("https://github.com/anna-ssg/anna/archive/%s.zip", trim)
					info.Printf("Bootstrapping from embedded commit %s -> %s\n", CommitHash, bsURL)
				}
			}

			if bsURL == "" {
				bsURL = "https://github.com/anna-ssg/anna/archive/refs/heads/main.zip"
			}

			// Safety check: refuse to overwrite an existing `site/` dir unless
			// --yes (`bsYes`) was supplied by the user.
			dest := "site"
			if _, err := os.Stat(dest); err == nil {
				if !bsYes {
					info.Printf("Refusing to bootstrap: %q already exists. Use --yes to overwrite.\n", dest)
					return
				}
				info.Printf("%q already exists; proceeding to overwrite because --yes was passed.\n", dest)
			}

			if !bsYes {
				fmt.Printf("This will download %s and extract the `site/` directory into ./site/. Continue? (y/N): ", bsURL)
				reader := bufio.NewReader(os.Stdin)
				line, _ := reader.ReadString('\n')
				line = strings.TrimSpace(line)
				if strings.ToLower(line) != "y" && strings.ToLower(line) != "yes" {
					info.Println("Aborted.")
					return
				}
			}

			helper := helpers.Helper{ErrorLogger: logger.New(os.Stderr)}
			if err := helper.BootstrapFromURL(bsURL); err != nil {
				// If the initial download fails with a 404 and we attempted a
				// commit-specific archive, try falling back to main.zip.
				if strings.Contains(err.Error(), "404 Not Found") && !strings.Contains(bsURL, "refs/heads/main.zip") {
					info.Printf("Download %s returned 404; falling back to main branch archive and retrying...\n", bsURL)
					fallback := "https://github.com/anna-ssg/anna/archive/refs/heads/main.zip"
					if err2 := helper.BootstrapFromURL(fallback); err2 != nil {
						fmt.Fprintln(os.Stderr, "fallback failed:", err2)
						os.Exit(1)
					}
				} else {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}

			info.Println("Bootstrapped site/ successfully")
		},
	}
	bootstrapCmd.Flags().BoolVarP(&bsYes, "yes", "y", false, "do not prompt; proceed non-interactively")
	bootstrapCmd.Flags().StringVar(&bsURL, "url", "https://github.com/anna-ssg/anna/archive/refs/heads/main.zip", "zip archive url to bootstrap from")
	rootCmd.AddCommand(bootstrapCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
