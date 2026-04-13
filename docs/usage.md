# CLI 使用指南

`xgtool` 提供了三個主要的子命令，用於資源提取與轉換。

## 提取圖形 (dump-graphic)

將遊戲圖形資源匯出為 JPEG 影像。

```bash
xgtool dump-graphic \
  -gif GraphicInfo.info \
  -gf Graphic.bin \
  -pf Pal.cgp \
  -o ./output
```

- `-gif`: 圖形索引檔路徑。
- `-gf`: 圖形資料檔路徑。
- `-pf`: 調色盤檔案路徑。
- `-o`: 輸出目錄（預設為 `output`）。
- `-dry-run`: 僅執行流程而不實際寫入檔案。

## 提取動畫 (dump-anime)

將遊戲動畫資源匯出為 GIF 影像。

```bash
xgtool dump-anime \
  -aif AnimeInfo.info \
  -af Anime.bin \
  -gif GraphicInfo.info \
  -gf Graphic.bin \
  -pf Pal.cgp \
  -o ./output
```

- `-aif`: 動畫索引檔路徑。
- `-af`: 動畫資料檔路徑。
- 其餘參數與 `dump-graphic` 相同，用於定位動畫引用的圖形資源。

## 轉換地圖 (convert-map)

將遊戲地圖轉換為 TMX 格式與對應的 PNG 圖磚。

```bash
xgtool convert-map \
  -mf 100.map \
  -gif GraphicInfo.info \
  -gf Graphic.bin \
  -pf Pal.cgp \
  -o ./output
```

- `-mf`: 地圖檔案路徑 (.MAP)。
- 輸出目錄中將包含一個 `.tmx` 檔案以及所有引用的圖磚影像檔案。
