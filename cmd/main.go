package main

import (
	"fmt"
	"github.com/cristalhq/acmd"
	"xgtool/cmd/convertmap"
	"xgtool/cmd/dumpanime"
	"xgtool/cmd/dumpgraphic"
)

var appVersion = ""
var buildTime = ""

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
		Version:        fmt.Sprintf("XGTool %s (built at: %s)", appVersion, buildTime),
	})

	if err := r.Run(); err != nil {
		r.Exit(err)
	}
}
