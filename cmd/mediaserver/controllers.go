package mediaserver

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"io"
	"xgtool/pkg"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var errInvalidVersion = errors.New("invalid version")
var errMissingID = errors.New("missing id or map_id")

// VerURI is a version string binding.
type VerURI struct {
	Ver string `uri:"ver" binding:"required"`
}

// QueryImg is a query string binding for querying specific graphic id.
type QueryImg struct {
	ID    *int32 `form:"id"`
	MapID *int32 `form:"map_id"`
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

func dumpGraphic(c *gin.Context) {
	uri, query, err := bindParams(c)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	res := versionResources[uri.Ver]

	var id int32
	if query.ID != nil {
		id = *query.ID
	} else {
		id = *query.MapID
	}

	if _, ok := res.GraphicIDIndex[id]; !ok {
		c.JSON(404, gin.H{"error": "id not found"})
		return
	}

	var graphic *pkg.Graphic
	if graphic, err = res.GraphicIDIndex[id].LoadGraphic(res.GraphicFile); err != nil {
		log.Err(err).Send()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Stream(func(w io.Writer) bool {
		var img image.Image
		if img, err = graphic.ImgRGBA(nil); err != nil {
			log.Err(err).Send()
			c.JSON(500, gin.H{"error": err.Error()})
			return true
		}
		if err = png.Encode(w, img); err != nil {
			log.Err(err).Send()
			c.JSON(500, gin.H{"error": err.Error()})
			return true
		}

		return false
	})

}

func dumpAnime(c *gin.Context) {
	uri, query, err := bindParams(c)
	if err != nil {
		c.JSON(404, gin.H{"error": "unknown version or id"})
		return
	}

	res := versionResources[uri.Ver]

	var id int32
	if query.ID != nil {
		id = *query.ID
	} else {
		id = *query.MapID
	}

	var animes []*pkg.Anime
	animes, err = res.AnimeInfoIndex[id].LoadAllAnimes(res.AnimeFile, res.GraphicIDIndex, res.GraphicFile)

	c.Stream(func(w io.Writer) bool {
		var img *gif.GIF
		if img, err = animes[0].GIF(res.Palette); err != nil {
			log.Err(err).Send()
			c.JSON(500, gin.H{"error": err.Error()})
			return true
		}
		if err = gif.EncodeAll(w, img); err != nil {
			log.Err(err).Send()
			c.JSON(500, gin.H{"error": err.Error()})
			return true
		}

		return false
	})

}

func bindParams(c *gin.Context) (uri VerURI, img QueryImg, err error) {
	if err = c.ShouldBind(&uri); err != nil {
		return
	}
	if _, ok := versionResources[uri.Ver]; !ok {
		err = fmt.Errorf("%w: %s", errInvalidVersion, uri.Ver)
		return
	}

	if err = c.ShouldBind(&img); err != nil {
		return
	}
	if img.MapID == nil && img.ID == nil {
		err = errMissingID
		return
	}

	return
}
