# 地圖格式 (MAP)

x-gate 地圖檔案 (.MAP) 用於儲存遊戲場景的地形、物件與元數據資訊。

## 地圖結構

每個地圖檔案包含一個 20 位元組的標頭，後接三個大型資料塊：

1. **標頭 (MapHeader)** (20 bytes):
   - `Magic` (12 bytes): 開頭為 `"MAP"`，後接 9 位元組的保留空間。
   - `Width/Height` (int32): 地圖的圖磚數量。
2. **地面層 (Ground Layer)**:
   - `uint16` 陣列，長度為 `Width * Height`。
   - 儲存每個座標對應的地面圖磚 ID。
3. **物件層 (Object Layer)**:
   - `uint16` 陣列，長度為 `Width * Height`。
   - 儲存地圖上的物件（如樹木、房屋、柵欄）的圖磚 ID。
4. **元數據層 (Meta Layer)**:
   - `uint16` 陣列，長度為 `Width * Height`。
   - 儲存遊戲邏輯（如碰撞檢測、傳送點、NPC 觸發區域）。

## TMX 轉換邏輯

`xgtool` 能將地圖轉換為 [Tiled Map Editor](https://www.mapeditor.org/) 的 `.tmx` 格式。

### 1. 座標轉換
由於 x-gate 使用 45 度角的等距投影（Isometric Projection），且與 Tiled 的座標系不完全一致，轉換過程中需要進行以下處理：
- **地圖旋轉**: 原始地圖資料需要逆時針旋轉 90 度，以適應 Tiled 的顯示效果。
- **物件定位**: 根據 `GraphicInfo` 中的寬度、高度與偏移量 (`OffX/OffY`)，計算物件在 TMX 中的精確 X/Y 座標。

### 2. 圖磚資源與 GID
- 為地面與物件層分別建立圖磚集（Tileset）。
- 在轉換時，將每個原始 ID 映射至全域圖磚 ID (Global Tile ID, GID)。
- 將關聯的圖形資源匯出為 PNG 圖檔，供 TMX 引用。

## 技術參數
- **圖磚尺寸**: 基準為 `64x47` 像素。
- **地圖方向**: `Isometric` (左上為原點)。
