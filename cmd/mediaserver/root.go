package mediaserver

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
)

func MediaServer(ctx context.Context, args []string) error {
	openVersionResources()
	defer closeVersionResources()

	r := gin.Default()

	r.GET("/ping", ping)
	r.GET("/graphic/:version", dumpGraphic)
	r.GET("/anime/:version", dumpAnime)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Send()
		}
	}()

	// handle Ctrl+C signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Send()
	}

	return nil
}
