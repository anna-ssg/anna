package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/anna-ssg/anna/v2/cmd/anna"
	"github.com/spf13/cobra"
)

func main() {
	var addr string
	var prof bool
	var renderDrafts bool
	var serve bool
	var webconsole bool
	var version bool
	var validateHTMLLayouts bool
	var renderAllSites bool
	var renderSite string

	Version := "v2.0.0" // to be set at build time $(git describe --tags)

	rootCmd := &cobra.Command{
		Use:   "anna",
		Short: "Static Site Generator",
		Run: func(cmd *cobra.Command, args []string) {
			annaCmd := anna.Cmd{
				RenderDrafts:      renderDrafts,
				Addr:              addr,
				RenderAll:         renderAllSites,
				ServeSpecificSite: renderSite,
				ErrorLogger:       log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
			}

			if renderAllSites {
				// todo
			} else {
				if serve {
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
					fmt.Println("Current version:", Version)
				}

				if validateHTMLLayouts {
					annaCmd.ValidateHTMLManager()
				}

				if webconsole {
					server := anna.NewWizardServer(":8080")
					go server.Start()
					<-anna.FormSubmittedCh // wait for response
					if err := server.Stop(); err != nil {
						log.Println(err)
					}
					annaCmd.VanillaRenderManager()
					annaCmd.LiveReloadManager()
				}
				annaCmd.VanillaRenderManager()
			}
		},
	}

	rootCmd.Flags().StringVarP(&addr, "addr", "a", "8000", "specify port to serve rendered content to")
	rootCmd.Flags().BoolVarP(&renderDrafts, "draft", "d", false, "renders draft posts")
	rootCmd.Flags().BoolVarP(&renderAllSites, "entire", "e", false, "render all the sites mentioned in anna.yml")
	rootCmd.Flags().BoolVarP(&validateHTMLLayouts, "layout", "l", false, "validates html layouts")
	rootCmd.Flags().StringVarP(&renderSite, "site", "o", "site", "specify the site directory to render")
	rootCmd.Flags().BoolVarP(&prof, "prof", "p", false, "enable profiling")
	rootCmd.Flags().BoolVarP(&serve, "serve", "s", false, "serve the rendered content")
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "prints current version number")
	rootCmd.Flags().BoolVarP(&webconsole, "webconsole", "w", false, "wizard to setup anna")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
