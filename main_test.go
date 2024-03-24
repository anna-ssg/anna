package main_test

// Benchmark currently does not work as project is being refactored

// import (
// 	"html/template"
// 	"log"
// 	"os"
// 	"testing"
//
// 	"github.com/acmpesuecc/anna/cmd/anna"
// )
//
// func BenchmarkMain(b *testing.B) {
// 	gen := anna.Generator{
// 		ErrorLogger: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
// 		Templates:   make(map[template.URL]anna.TemplateData),
// 		TagsMap:     make(map[string][]anna.TemplateData),
// 	}
// 	for i := 0; i < b.N; i++ {
// 		gen.RenderSite("")
// 	}
// }
