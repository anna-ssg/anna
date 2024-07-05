package main

import (
	"log"
	"os"
	"time"

	"github.com/anna-ssg/anna/v3/cmd/anna"
	"github.com/spf13/cobra"
)

func main() {
	var addr string
	var prof bool
	var renderDrafts bool
	var serve string
	var webconsole bool
	var version bool
	var validateHTMLLayouts bool
	var renderSpecificSite string

	Version := "v2.0.0" // to be set at build time $(git describe --tags)

	rootCmd := &cobra.Command{
		Use:   "anna",
		Short: "Static Site Generator",
		Run: func(cmd *cobra.Command, args []string) {
			annaCmd := anna.Cmd{
				RenderDrafts:       renderDrafts,
				Addr:               addr,
				RenderSpecificSite: renderSpecificSite,
				ServeSpecificSite:  serve,
				ErrorLogger:        log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
				InfoLogger:         log.New(os.Stderr, "LOG\t", log.Ldate|log.Ltime),
			}

			if serve != "" {
				annaCmd.LiveReload = true
				annaCmd.LiveReloadManager()
			}

			if prof {
				startTime := time.Now()
				annaCmd.VanillaRenderManager()
				elapsedTime := time.Since(startTime)
				annaCmd.PrintStats(elapsedTime)
			}

			if version {
				annaCmd.InfoLogger.Println("Current version:", Version)
			}

			if validateHTMLLayouts {
				annaCmd.ValidateHTMLManager()
			}

			if webconsole {
				server := anna.NewWizardServer(":8080")
				go server.Start()
				<-anna.FormSubmittedCh // wait for response
				if err := server.Stop(); err != nil {
					annaCmd.InfoLogger.Println(err)
				}
				annaCmd.LiveReloadManager()
			}

			annaCmd.VanillaRenderManager()
		},
	}

	rootCmd.Flags().StringVarP(&addr, "addr", "a", "8000", "specify port to serve rendered content to")
	rootCmd.Flags().BoolVarP(&renderDrafts, "draft", "d", false, "renders draft posts")
	rootCmd.Flags().BoolVarP(&validateHTMLLayouts, "layout", "l", false, "validates html layouts")
	// Do not set default values for string flags
	rootCmd.Flags().StringVarP(&renderSpecificSite, "render-site", "r", "", "specify the specific site directory to render")
	rootCmd.Flags().BoolVarP(&prof, "prof", "p", false, "enable profiling")
	rootCmd.Flags().StringVarP(&serve, "serve", "s", "", "specify the specific site directory to serve")
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "prints current version number")
	rootCmd.Flags().BoolVarP(&webconsole, "webconsole", "w", false, "wizard to setup anna")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
