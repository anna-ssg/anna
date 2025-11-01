package main_test

import (
	"testing"

	"github.com/anna-ssg/anna/v3/cmd/anna"
)

const siteDirPath = "site/"

func BenchmarkMain(b *testing.B) {
	annaCmd := anna.Cmd{
		RenderDrafts: true,
	}
	for i := 0; i < b.N; i++ {
		annaCmd.VanillaRender(siteDirPath)
	}
}
