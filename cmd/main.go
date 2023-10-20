package main

import (
	"github.com/cristalhq/acmd"
	"xgtool/cmd/convertmap"
	"xgtool/cmd/dumpanime"
	"xgtool/cmd/dumpgraphic"
)

// Version it can be set by ldflags="main.Version=x.x.x"
var Version = "0.0.0"

func main() {
	cmds := []acmd.Command{
		{
			Name:        "dump-graphic",
			Description: "Dump graphic from graphic & graphic info file",
			ExecFunc:    dumpgraphic.DumpGraphic,
		},
		{
			Name:        "dump-anime",
			Description: "Dump anime from anime & anime info file",
			ExecFunc:    dumpanime.DumpAnime,
		},
		{
			Name:        "convert-map",
			Description: "Convert map into TMX format",
			ExecFunc:    convertmap.ConvertMap,
		},
	}

	r := acmd.RunnerOf(cmds, acmd.Config{
		AppName:        "xgtool",
		AppDescription: "The toolchain of x-gate",
		Version:        Version,
	})

	if err := r.Run(); err != nil {
		r.Exit(err)
	}
}
