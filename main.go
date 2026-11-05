package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/anna-ssg/anna/cmd/anna"
	"github.com/anna-ssg/anna/pkg/helpers"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// Populated at build time by GoReleaser.
var (
	Version        = "dev"
	FullCommitHash = ""
)

func main() {
	var addr string
	var prof bool
	var renderDrafts bool
	var serve bool
	var version bool
	var siteDirPath string

	console := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.Kitchen,
		NoColor:    false,
	}
	globalLogger := zerolog.New(console).With().Timestamp().Stack().Logger()

	rootCmd := &cobra.Command{
		Use:   "anna",
		Short: "Static Site Generator",
		Run: func(cmd *cobra.Command, args []string) {
			if version {
				ver := Version
				if FullCommitHash != "" {
					ver = fmt.Sprintf("%s (commit: %s)", Version, FullCommitHash)
				}
				globalLogger.Info().Msg(ver)
				return
			}

			siteDirPath = path.Clean(siteDirPath) + "/"

			helper := helpers.Helper{
				Logger: &globalLogger,
			}

			// Bootstrap a new site if one doesn't exist.
			if err := helper.EnsureSiteExists(); err != nil {
				globalLogger.Fatal().Err(err).Send()
			}

			annaCmd := anna.Cmd{
				RenderDrafts: renderDrafts,
				Addr:         addr,
				LiveReload:   serve,
				SiteDirPath:  siteDirPath,
				Logger:       &globalLogger,
			}

			if serve {
				annaCmd.StartLiveReload(siteDirPath)
			}

			built := false

			if prof {
				start := time.Now()
				count := annaCmd.VanillaRender(siteDirPath)
				elapsed := time.Since(start)
				annaCmd.PrintStats(elapsed)
				globalLogger.Info().Msgf("Built %d pages in %s\n", count, elapsed)
				built = true
			}

			if !built {
				start := time.Now()
				count := annaCmd.VanillaRender(siteDirPath)
				elapsed := time.Since(start)
				globalLogger.Info().Msgf("Built %d pages in %s\n", count, elapsed)
			}
		},
	}

	rootCmd.Flags().StringVarP(
		&addr,
		"addr",
		"a",
		"localhost:8000",
		"specify address over which rendered content is served",
	)

	rootCmd.Flags().BoolVarP(
		&renderDrafts,
		"draft",
		"d",
		false,
		"renders draft posts",
	)

	rootCmd.Flags().StringVarP(
		&siteDirPath,
		"path",
		"p",
		"site",
		"specify the specific site directory to render",
	)

	rootCmd.Flags().BoolVar(
		&prof,
		"prof",
		false,
		"enable profiling",
	)

	rootCmd.Flags().BoolVarP(
		&serve,
		"serve",
		"s",
		false,
		"serve the rendered site and watch for file updates",
	)

	rootCmd.Flags().BoolVarP(
		&version,
		"version",
		"v",
		false,
		"prints current version number",
	)

	var bsYes bool

	bootstrapCmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Extract the embedded starter site",
		Run: func(cmd *cobra.Command, args []string) {
			if _, err := os.Stat("site"); err == nil && !bsYes {
				globalLogger.Info().Msg("site/ already exists. Use --yes to overwrite.")
				return
			}

			if !bsYes {
				globalLogger.Info().Msg("Extract the embedded starter site into ./site? (y/N): ")

				reader := bufio.NewReader(os.Stdin)
				line, _ := reader.ReadString('\n')
				line = strings.TrimSpace(strings.ToLower(line))

				if line != "y" && line != "yes" {
					globalLogger.Info().Msg("Aborted.")
					return
				}
			}

			helper := helpers.Helper{
				Logger: &globalLogger,
			}

			if err := helper.BootstrapEmbedded(bsYes); err != nil {
				globalLogger.Fatal().Err(err).Send()
			}

			globalLogger.Info().Msg("Bootstrapped embedded site successfully.")
		},
	}

	bootstrapCmd.Flags().BoolVarP(
		&bsYes,
		"yes",
		"y",
		false,
		"overwrite existing site",
	)

	rootCmd.AddCommand(bootstrapCmd)

	if err := rootCmd.Execute(); err != nil {
		globalLogger.Fatal().Err(err).Send()
	}
}
