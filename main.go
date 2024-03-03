package main

import (
	"log"
	"os"

	"github.com/acmpesuecc/ssg/cmd/ssg"
	"github.com/spf13/cobra"
)

func main() {
	var serve bool
	var addr string
    var draft bool

	rootCmd := &cobra.Command{
		Use:   "ssg",
		Short: "Static Site Generator",
		Run: func(cmd *cobra.Command, args []string) {
			generator := ssg.Generator{
				ErrorLogger: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
			}
            if draft { 
                generator.Draft = true
            }
			generator.RenderSite(addr)

			if serve {
				generator.ServeSite(addr)
			}
		},
	}

	rootCmd.Flags().BoolVarP(&serve, "serve", "s", false, "serve the rendered content")
	rootCmd.Flags().StringVarP(&addr, "addr", "a", "8000", "ip address to serve rendered content to")
	rootCmd.Flags().BoolVarP(&draft, "draft", "d", false, "renders draft posts")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
