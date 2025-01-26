package main

import (
	"flag"
	"os"
	"path"

	"github.com/evanw/esbuild/pkg/api"
)

var entryPoints = [...]string{"index.js"}

// Global options
var srcDir string
var outDir string
var sourceMap api.SourceMap
var joinedEntryPoints []string

func main() {
	parseFlags()

	result := api.Build(api.BuildOptions{
		EntryPoints: joinedEntryPoints,
		Outdir:      outDir,
		Sourcemap:   sourceMap,
		Bundle:      true,
		Write:       true,

		LogLevel: api.LogLevelInfo,
	})

	if len(result.Errors) > 0 {
		os.Exit(1)
	}
}

func parseFlags() {
	flag.StringVar(&srcDir, "src", "web/src", "JS source directory")
	flag.StringVar(&outDir, "out", "web/dist", "Output directory")
	release := flag.Bool("release", false, "Whether to build for release")

	flag.Parse()

	if *release {
		sourceMap = api.SourceMapNone
	} else {
		sourceMap = api.SourceMapLinked
	}

	for _, entry := range entryPoints {
		joinedEntry := path.Join(srcDir, entry)
		joinedEntryPoints = append(joinedEntryPoints, joinedEntry)
	}
}
