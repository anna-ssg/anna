package main

import (
	"html/template"
	"log"
	"os"

	"github.com/acmpesuecc/anna/pkg/parser"
	"github.com/spf13/cobra"
)

func main() {
	var serve bool
	var addr string
	var draft bool
	var validateHTML bool
	var prof bool
	rootCmd := &cobra.Command{
		Use:   "anna",
		Short: "Static Site Generator",
		Run: func(cmd *cobra.Command, args []string) {
			parser := parser.Parser{
				ErrorLogger: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
				Templates:   make(map[template.URL]parser.TemplateData),
				TagsMap:     make(map[string][]parser.TemplateData),
			}

			if draft {
				parser.RenderDrafts = true
			}

			if serve {
				// parser.StartLiveReload(addr)
			}

			if prof {
				// TODO: To be filled later
			}

			if validateHTML {
				// anna.ValidateHTMLContent()
			}
			// parser.RenderSite("")
		},
	}

	rootCmd.Flags().BoolVarP(&serve, "serve", "s", false, "serve the rendered content")
	rootCmd.Flags().StringVarP(&addr, "addr", "a", "8000", "ip address to serve rendered content to")
	rootCmd.Flags().BoolVarP(&draft, "draft", "d", false, "renders draft posts")
	rootCmd.Flags().BoolVarP(&validateHTML, "validate-html", "v", false, "validate semantic HTML")
	rootCmd.Flags().BoolVar(&prof, "prof", false, "enable profiling")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
