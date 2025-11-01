package main

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/anna-ssg/anna/v3/cmd/anna"
	"github.com/spf13/cobra"
)

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

			annaCmd := anna.Cmd{
				RenderDrafts: renderDrafts,
				Addr:         addr,
				LiveReload:   serve,
				SiteDirPath:  siteDirPath,
				ErrorLogger:  log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime),
				InfoLogger:   log.New(os.Stderr, "LOG\t", log.Ldate|log.Ltime),
			}

			if serve {
				annaCmd.StartLiveReload(siteDirPath)
			}

			if prof {
				startTime := time.Now()
				annaCmd.VanillaRender(siteDirPath)
				elapsedTime := time.Since(startTime)
				annaCmd.PrintStats(elapsedTime)
			}

			if version {
				annaCmd.InfoLogger.Println("Current version:", Version)
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

			annaCmd.VanillaRender(siteDirPath)
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

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
