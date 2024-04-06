package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/acmpesuecc/anna/cmd/anna"
	"github.com/spf13/cobra"
)

func main() {
	var serve bool
	var addr string
	var webconsole bool
	var renderDrafts bool
	var version bool
	var prof bool

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

			if version, err := exec.Command("git", "describe", "--tags").Output(); err == nil {
				fmt.Println("Current version:", strings.TrimSpace(string(version)))
			} else {
				fmt.Println("Error:", err)
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

	rootCmd.Flags().BoolVarP(&serve, "serve", "s", false, "serve the rendered content")
	rootCmd.Flags().StringVarP(&addr, "addr", "a", "8000", "ip address to serve rendered content to")
	rootCmd.Flags().BoolVarP(&renderDrafts, "draft", "d", false, "renders draft posts")
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "prints current version number")
	rootCmd.Flags().BoolVarP(&prof, "prof", "p", false, "enable profiling")
	rootCmd.Flags().BoolVarP(&webconsole, "webconsole", "w", false, "wizard to setup anna")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
