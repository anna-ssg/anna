package main_test

import (
	"testing"

	"github.com/anna-ssg/anna/v2/cmd/anna"
)

func BenchmarkMain(b *testing.B) {
	annaCmd := anna.Cmd{
		RenderDrafts: true,
	}
	for i := 0; i < b.N; i++ {
		annaCmd.VanillaRenderManager()
	}
}
