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

// Populated at build time by GoReleaser.
var (
	Version        = "dev"
	FullCommitHash = ""
)

func main() {
	var addr string
	var prof bool
	var renderDrafts bool
	var serve bool
	var version bool
	var siteDirPath string

	rootCmd := &cobra.Command{
		Use:   "anna",
		Short: "Static Site Generator",
		Run: func(cmd *cobra.Command, args []string) {
			if version {
				ver := Version
				if FullCommitHash != "" {
					ver = fmt.Sprintf("%s (commit: %s)", Version, FullCommitHash)
				}
				fmt.Println(ver)
				return
			}

			siteDirPath = path.Clean(siteDirPath) + "/"

			// Automatically bootstrap a new site if one doesn't exist.
			if _, err := os.Stat(siteDirPath); os.IsNotExist(err) {
				helper := helpers.Helper{
					ErrorLogger: logger.New(os.Stderr),
				}
				if err := helper.BootstrapEmbedded(false); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
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
				start := time.Now()
				count := annaCmd.VanillaRender(siteDirPath)
				elapsed := time.Since(start)
				annaCmd.PrintStats(elapsed)
				annaCmd.InfoLogger.Printf("Built %d pages in %s\n", count, elapsed)
				built = true
			}

			if !built {
				start := time.Now()
				count := annaCmd.VanillaRender(siteDirPath)
				elapsed := time.Since(start)
				annaCmd.InfoLogger.Printf("Built %d pages in %s\n", count, elapsed)
			}
		},
	}

	rootCmd.Flags().StringVarP(
		&addr,
		"addr",
		"a",
		"localhost:8000",
		"specify address over which rendered content is served",
	)

	rootCmd.Flags().BoolVarP(
		&renderDrafts,
		"draft",
		"d",
		false,
		"renders draft posts",
	)

	rootCmd.Flags().StringVarP(
		&siteDirPath,
		"path",
		"p",
		"site",
		"specify the specific site directory to render",
	)

	rootCmd.Flags().BoolVar(
		&prof,
		"prof",
		false,
		"enable profiling",
	)

	rootCmd.Flags().BoolVarP(
		&serve,
		"serve",
		"s",
		false,
		"serve the rendered site and watch for file updates",
	)

	rootCmd.Flags().BoolVarP(
		&version,
		"version",
		"v",
		false,
		"prints current version number",
	)

	var bsYes bool

	bootstrapCmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Extract the embedded starter site",
		Run: func(cmd *cobra.Command, args []string) {
			info := logger.New(os.Stderr)

			if _, err := os.Stat("site"); err == nil && !bsYes {
				info.Println("site/ already exists. Use --yes to overwrite.")
				return
			}

			if !bsYes {
				fmt.Print("Extract the embedded starter site into ./site? (y/N): ")

				reader := bufio.NewReader(os.Stdin)
				line, _ := reader.ReadString('\n')
				line = strings.TrimSpace(strings.ToLower(line))

				if line != "y" && line != "yes" {
					info.Println("Aborted.")
					return
				}
			}

			helper := helpers.Helper{
				ErrorLogger: logger.New(os.Stderr),
			}

			if err := helper.BootstrapEmbedded(bsYes); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			info.Println("Bootstrapped embedded site successfully.")
		},
	}

	bootstrapCmd.Flags().BoolVarP(
		&bsYes,
		"yes",
		"y",
		false,
		"overwrite existing site",
	)

	rootCmd.AddCommand(bootstrapCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}