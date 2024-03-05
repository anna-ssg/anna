package main

import (
	"html/template"
	"log"
	"os"

	"github.com/acmpesuecc/ssg/cmd/ssg"
	"github.com/spf13/cobra"
)

func main() {
	var serve bool
	var addr string
	var draft bool
	var validateHTML bool

	rootCmd := &cobra.Command{
		Use:   "ssg",
		Short: "Static Site Generator",
		Run: func(cmd *cobra.Command, args []string) {
			generator := ssg.Generator{
				ErrorLogger: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
				Templates:   make(map[template.URL]ssg.TemplateData),
			}
			if draft {
				generator.RenderDrafts = true
			}

			if serve {
				generator.StartLiveReload(addr)
			}
			generator.RenderSite(addr)

			if validateHTML {
				ssg.ValidateHTMLContent()
			}
		},
	}

	rootCmd.Flags().BoolVarP(&serve, "serve", "s", false, "serve the rendered content")
	rootCmd.Flags().StringVarP(&addr, "addr", "a", "8000", "ip address to serve rendered content to")
	rootCmd.Flags().BoolVarP(&draft, "draft", "d", false, "renders draft posts")
	rootCmd.Flags().BoolVarP(&validateHTML, "validate-html", "v", false, "validate semantic HTML")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
