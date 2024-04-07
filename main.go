package main

import (
	"fmt"
	"log"
	"time"

	"github.com/acmpesuecc/anna/cmd/anna"
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
	var Version string = "v1.0.0-alpha-15-g1eb8e48" // to be set at build time $(git describe --tags)

	rootCmd := &cobra.Command{
		Use:   "anna",
		Short: "Static Site Generator",
		Run: func(cmd *cobra.Command, args []string) {
			annaCmd := anna.Cmd{
				RenderDrafts: renderDrafts,
				Addr:         addr,
			}

			if serve {
				annaCmd.StartLiveReload()
			}

			if prof {
				startTime := time.Now()
				anna.StartProfiling(&annaCmd)

				elapsedTime := time.Since(startTime)
				go anna.PrintStats(elapsedTime)
				anna.RunProfilingServer()
			}

			if version {
				fmt.Println("Current version:", Version)
			}

			if validateHTMLLayouts {
				annaCmd.ValidateHTMLContent()
			}

			if webconsole {
				server := anna.NewWizardServer(":8080")
				go server.Start()
				<-anna.FormSubmittedCh // wait for response
				server.Stop()          // stop the server
				annaCmd.VanillaRender()
				annaCmd.StartLiveReload()
			}
			annaCmd.VanillaRender()
		},
	}

	rootCmd.Flags().StringVarP(&addr, "addr", "a", "8000", "ip address to serve rendered content to")
	rootCmd.Flags().BoolVarP(&renderDrafts, "draft", "d", false, "renders draft posts")
	rootCmd.Flags().BoolVarP(&validateHTMLLayouts, "layout", "l", false, "validates html layouts")
	rootCmd.Flags().BoolVarP(&prof, "prof", "p", false, "enable profiling")
	rootCmd.Flags().BoolVarP(&serve, "serve", "s", false, "serve the rendered content")
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "prints current version number")
	rootCmd.Flags().BoolVarP(&webconsole, "webconsole", "w", false, "wizard to setup anna")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
