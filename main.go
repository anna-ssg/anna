package main

import (
	"log"
	"os"

	"github.com/anirudhsudhir/ssg/cmd/ssg"
)

func main() {
	generator := ssg.Generator{
		ErrorLogger: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	generator.ParseMarkdown()
	generator.RenderSite()
}
