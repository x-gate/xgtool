package main

import (
	"github.com/cristalhq/acmd"
	"xgtool/cmd/convert_map"
	"xgtool/cmd/dump_graphic"
)

var Version = "0.0.0"

func main() {
	cmds := []acmd.Command{
		{
			Name:        "dump-graphic",
			Description: "Dump graphic from graphic & graphic info file",
			ExecFunc:    dump_graphic.DumpGraphic,
		},
		{
			Name:        "convert-map",
			Description: "Convert map into TMX format",
			ExecFunc:    convert_map.ConvertMap,
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
