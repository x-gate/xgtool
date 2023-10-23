package mediaserver

import "xgtool/pkg"

var versionResourcePath = map[string]struct {
	name string
	gif  string
	gf   string
	aif  string
	af   string
	pf   string
}{
	"v1": {
		name: "神獸傳奇 + 魔弓傳奇",
		gif:  "testdata/graphic_info/GraphicInfo_66.bin",
		gf:   "testdata/graphic/Graphic_66.bin",
		aif:  "testdata/anime_info/AnimeInfo_4.bin",
		af:   "testdata/anime/Anime_4.bin",
		pf:   "testdata/palette/palet_00.cgp",
	},
	"v2": {
		name: "龍之沙漏",
		gif:  "testdata/graphic_info/GraphicInfoEx_5.bin",
		gf:   "testdata/graphic/GraphicEx_5.bin",
		aif:  "testdata/anime_info/AnimeInfoEx_1.bin",
		af:   "testdata/anime/AnimeEx_1.bin",
		pf:   "testdata/palette/palet_00.cgp",
	},
	"v3-0": {
		name: "樂園之卵（精靈）",
		gif:  "testdata/graphic_info/GraphicInfoV3_19.bin",
		gf:   "testdata/graphic/GraphicV3_19.bin",
		aif:  "testdata/anime_info/AnimeInfoV3_8.bin",
		af:   "testdata/anime/AnimeV3_8.bin",
		pf:   "testdata/palette/palet_00.cgp",
	},
	"v3-1": {
		name: "樂園之卵",
		gif:  "testdata/graphic_info/GraphicInfo_PUK2_2.bin",
		gf:   "testdata/graphic/Graphic_PUK2_2.bin",
		aif:  "testdata/anime_info/AnimeInfo_PUK2_4.bin",
		af:   "testdata/anime/Anime_PUK2_4.bin",
		pf:   "testdata/palette/palet_00.cgp",
	},
	"v4": {
		name: "天界之騎士與星詠之歌姬",
		gif:  "testdata/graphic_info/GraphicInfo_PUK3_1.bin",
		gf:   "testdata/graphic/Graphic_PUK3_1.bin",
		aif:  "testdata/anime_info/AnimeInfo_PUK3_2.bin",
		af:   "testdata/anime/Anime_PUK3_2.bin",
		pf:   "testdata/palette/palet_00.cgp",
	},
	"v5": {
		name: "砂之記憶與覺醒之光",
		gif:  "testdata/graphic_info/GraphicInfo_Joy_125.bin",
		gf:   "testdata/graphic/Graphic_Joy_125.bin",
		aif:  "testdata/anime_info/AnimeInfo_Joy_91.bin",
		af:   "testdata/anime/Anime_Joy_91.bin",
		pf:   "testdata/palette/palet_00.cgp",
	},
	"v6": {
		name: "輪迴之守",
		gif:  "testdata/graphic_info/GraphicInfo_Joy_CH1.bin",
		gf:   "testdata/graphic/Graphic_Joy_CH1.bin",
		aif:  "testdata/anime_info/AnimeInfo_Joy_CH1.Bin",
		af:   "testdata/anime/Anime_Joy_CH1.Bin",
		pf:   "testdata/palette/palet_00.cgp",
	},
	"v7": {
		name: "天使之降臨",
		gif:  "testdata/graphic_info/GraphicInfo_Joy_EX_152.bin",
		gf:   "testdata/graphic/Graphic_Joy_EX_152.bin",
		aif:  "testdata/anime_info/AnimeInfo_Joy_EX_146.bin",
		af:   "testdata/anime/Anime_Joy_EX_146.bin",
		pf:   "testdata/palette/palet_00.cgp",
	},
}

var versionResources = map[string]pkg.Resources{}

func openVersionResources() {
	for ver, res := range versionResourcePath {
		tmp := pkg.Resources{}

		if err := tmp.OpenGraphicInfo(res.gif); err != nil {
			panic(err)
		}
		if err := tmp.OpenGraphic(res.gf); err != nil {
			panic(err)
		}
		if err := tmp.OpenAnimeInfo(res.aif); err != nil {
			panic(err)
		}
		if err := tmp.OpenAnime(res.af); err != nil {
			panic(err)
		}
		if err := tmp.OpenPalette(res.pf); err != nil {
			panic(err)
		}

		versionResources[ver] = tmp
	}
}

func closeVersionResources() {
	for _, res := range versionResources {
		res.Close()
	}
}
