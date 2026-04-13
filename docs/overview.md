# xgtool 專案概覽

`xgtool` 是一個用於解析 x-gate 遊戲資源的 Go 語言工具集。它能夠處理調色盤、圖形、動畫和地圖等核心資源格式，並提供 CLI 工具將這些資源轉換為現代通用的格式（如 PNG, GIF, TMX）。

## 核心功能

- **圖形提取** (`dump-graphic`): 從遊戲資料檔中提取圖形並轉換為 JPG/PNG 格式。
- **動畫轉換** (`dump-anime`): 將遊戲動畫資源轉換為 GIF 格式。
- **地圖轉換** (`convert-map`): 將遊戲地圖轉換為 Tiled Map Editor 支援的 TMX 格式。

## 專案結構

- `cmd/`: 包含 CLI 程式的實作，分為 `dump-graphic`, `dump-anime`, `convert-map` 三個主要子命令。
- `pkg/`: 核心邏輯層，定義了各類資源的資料結構與解析邏輯。
  - `palette.go`: 調色盤處理。
  - `graphic.go`: 圖形資源處理。
  - `anime.go`: 動畫資源處理。
  - `map.go`: 地圖資源處理。
- `internal/`: 內部工具與輔助編碼邏輯。
  - `codec.go`: 實作了遊戲特有的 RLE 壓縮/解壓縮演算法。
  - `tmx/`: TMX 格式的相關資料結構。
  - `mat/`: 矩陣運算，用於地圖旋轉與座標轉換。

## 資源處理流程

1. **載入索引 (Info File)**: 首先讀取索引檔（如 `GraphicInfo`），獲取資源在資料檔中的偏移量、長度與屬性。
2. **定位並讀取資料 (Data File)**: 根據索引資訊，在資料檔中定位原始位元組。
3. **解碼 (Decoding)**: 若資料經過壓縮，則調用 `internal/codec` 進行解壓縮。
4. **渲染/轉換 (Rendering/Conversion)**: 結合調色盤資訊，將原始資料轉換為標準影像格式或結構化資料。
