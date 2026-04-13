# xgtool 技術文件

歡迎使用 `xgtool` 技術文件。本文件詳細介紹了 x-gate 資源格式的解析實作與工具使用。

## 目錄

- [專案概覽](overview.md) - 專案架構與核心功能介紹。
- [CLI 使用指南](usage.md) - 如何使用 `xgtool` 工具進行資源提取。

### 資源格式規範

- [調色盤 (CGP)](formats/palette.md) - 8 位元索引顏色的構成與透明度處理。
- [圖形 (RD)](formats/graphic.md) - 圖像索引與資料儲存格式。
- [RLE 編解碼 (Codec)](formats/codec.md) - 遊戲特有的壓縮演算法詳解。
- [動畫 (Anime)](formats/anime.md) - 多層級動畫結構與 GIF 合成邏輯。
- [地圖 (MAP)](formats/map.md) - 地圖層級結構與 TMX 轉換坐標映射。
