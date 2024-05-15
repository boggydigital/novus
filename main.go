package main

import (
	"bytes"
	"embed"
	_ "embed"
	"github.com/boggydigital/clo"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/novus/cli"
	"github.com/boggydigital/novus/data"
	"github.com/boggydigital/novus/rest"
	"github.com/boggydigital/pathways"
	"html/template"
	"os"
	"sync"
)

var (
	once = sync.Once{}
	//go:embed "templates/*.gohtml"
	templates embed.FS
	//go:embed "stencil_app/styles/css.gohtml"
	stencilAppStyles embed.FS
	//go:embed "cli-commands.txt"
	cliCommands []byte
	//go:embed "cli-help.txt"
	cliHelp []byte
)

var tmplFuncs = template.FuncMap{
	"empty": empty,
}

const (
	directoriesFilename = "directories.txt"
)

func main() {

	nod.EnableStdOutPresenter()

	ns := nod.Begin("novus is checking for any news")
	defer ns.End()

	if err := pathways.Setup(
		directoriesFilename,
		data.DefaultNovusRootDir,
		nil,
		data.AllAbsDirs()...); err != nil {
		_ = ns.EndWithError(err)
		os.Exit(1)
	}

	once.Do(func() {
		if err := rest.Init(templates, tmplFuncs, stencilAppStyles); err != nil {
			_ = ns.EndWithError(err)
			os.Exit(1)
		}
		cli.Init(templates, tmplFuncs)
	})

	defs, err := clo.Load(
		bytes.NewBuffer(cliCommands),
		bytes.NewBuffer(cliHelp),
		nil)
	if err != nil {
		_ = ns.EndWithError(err)
		os.Exit(1)
	}

	clo.HandleFuncs(map[string]clo.Handler{
		"decode":        cli.DecodeHandler,
		"diff":          cli.DiffHandler,
		"get-content":   cli.GetContentHandler,
		"match-content": cli.MatchContentHandler,
		"publish-atom":  cli.PublishAtomHandler,
		"reduce":        cli.ReduceHandler,
		"reset-changes": cli.ResetChangesHandler,
		"reset-errors":  cli.ResetErrorsHandler,
		"serve":         cli.ServeHandler,
		"sync":          cli.SyncHandler,
		"test-sources":  cli.TestSourcesHandler,
		"version":       cli.VersionHandler,
	})

	if err := defs.AssertCommandsHaveHandlers(); err != nil {
		_ = ns.EndWithError(err)
		os.Exit(1)
	}

	if err := defs.Serve(os.Args[1:]); err != nil {
		_ = ns.EndWithError(err)
		os.Exit(1)
	}
}

func empty(maps ...map[string][]string) bool {
	for _, m := range maps {
		if len(m) > 0 {
			return false
		}
	}
	return true
}
