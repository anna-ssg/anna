package main

import (
	"flag"
	"log"
	"os"

	"github.com/acmpesuecc/ssg/cmd/ssg"
)

func main() {
	serve := flag.Bool("serve", false, "serve the rendered content")
	addr := flag.String("addr", ":8000", "ip address to serve rendered content to")
	flag.Parse()

	generator := ssg.Generator{
		ErrorLogger: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	generator.RenderSite()

	if *serve {
		generator.ServeSite(*addr)
	}
}
