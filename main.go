package main

import (
	"log"
	"time"

	"github.com/acmpesuecc/anna/cmd/anna"
	"github.com/spf13/cobra"
)

func main() {
	var serve bool
	var addr string
	var webconsole bool
	var renderDrafts bool
	var validateHTML bool
	var prof bool
	anna.StartProfiling()
	startTime := time.Now()
	rootCmd := &cobra.Command{
		Use:   "anna",
		Short: "Static Site Generator",
		Run: func(cmd *cobra.Command, args []string) {
			annaCmd := anna.Cmd{
				RenderDrafts: renderDrafts,
				Addr:         addr,
			}

			if serve {
				if prof {
					go func() {
						for {
							time.Sleep(1 * time.Second) //change as per needed
							anna.PrintStats(time.Since(startTime))
						}
					}()
				}
				annaCmd.StartLiveReload()
			}
			if !prof {
				annaCmd.VanillaRender()

			}
			if prof {
				annaCmd.VanillaRender()

				// go anna.StartProfiling()
				// startTime := time.Now()
				// annaCmd.VanillaRender()
				elapsedTime := time.Now().Sub(startTime)
				// aPrintStats(elapsedTime)
				go anna.PrintStats(elapsedTime)
				defer anna.StopProfiling()
			}

			if validateHTML {
				// anna.ValidateHTMLContent()
				cmd.Println("TODO: To be filled later")
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
	// go anna.CreateIndex()

	rootCmd.Flags().BoolVarP(&serve, "serve", "s", false, "serve the rendered content")
	rootCmd.Flags().StringVarP(&addr, "addr", "a", "8000", "ip address to serve rendered content to")
	rootCmd.Flags().BoolVarP(&renderDrafts, "draft", "d", false, "renders draft posts")
	rootCmd.Flags().BoolVarP(&validateHTML, "validate-html", "v", false, "validate semantic HTML")
	rootCmd.Flags().BoolVarP(&prof, "prof", "p", false, "enable profiling")
	rootCmd.Flags().BoolVarP(&webconsole, "webconsole", "w", false, "wizard to setup anna")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
