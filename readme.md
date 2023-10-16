# XGTool

The toolchain of x-gate, maybe works with CrossGate(?).

## Available Tools

### Dump Graphic

```shell
$ export GRAPHIC_INFO_FILE="/Game/Crossgate/bin/GraphicInfo_66.bin" && \
  export GRAPHIC_FILE="/Game/Crossgate/bin/Graphic_66.bin" && \
  export PALETTE_PATH="/Game/Crossgate/bin/pal/Palette" && \

$ go run ./cmd/root.go dump-graphic \
    -gif $GRAPHIC_INFO_FILE \
    -gf  $GRAPHIC_FILE \
    -pf  $PALETTE_PATH \
    -dry-run
```