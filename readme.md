# XGTool

The toolchain of x-gate, maybe works with CrossGate(?).

## Available Tools

### Dump Graphic

Dump graphics from `GraphicInfo.bin` and `Graphic.bin` to `./output` directory.

```shell
$ export GRAPHIC_INFO_FILE="/Game/Crossgate/bin/GraphicInfo_66.bin" && \
  export GRAPHIC_FILE="/Game/Crossgate/bin/Graphic_66.bin" && \
  export PALETTE_PATH="/Game/Crossgate/bin/pal/Palette"

$ go run ./cmd/dump_graphic.go \
    -gif $GRAPHIC_INFO_FILE \
    -gf  $GRAPHIC_FILE \
    -pf  $PALETTE_PATH \
    -dry-run
```

### Convert Map

Convert map into tmx (json) format.

```shell
$ export GRAPHIC_INFO_FILE="/Game/Crossgate/bin/GraphicInfo_66.bin" && \
  export GRAPHIC_FILE="/Game/Crossgate/bin/Graphic_66.bin" && \
  export PALETTE_PATH="/Game/Crossgate/bin/pal/Palette" && \
  export MAP_FILE="/Game/Crossgate/map/1000.bin" # 法蘭城

$ go run ./cmd/convert_map.go \
    -gif $GRAPHIC_INFO_FILE \
    -gf  $GRAPHIC_FILE \
    -pf  $PALETTE_PATH \
    -mf  $MAP_FILE \
    -o output -n map.json
```