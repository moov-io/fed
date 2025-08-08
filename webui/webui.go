package webui

import (
	"embed"
	"net/http"
	"os"

	"github.com/moov-io/base/log"

	"github.com/gorilla/mux"
)

//go:embed *.html
var WebRoot embed.FS

type Controller interface {
	AppendRoutes(router *mux.Router) *mux.Router
}

func NewController(logger log.Logger) Controller {
	return &controller{
		logger:   logger,
		basePath: os.Getenv("BASE_PATH"),
	}
}

type controller struct {
	logger   log.Logger
	basePath string
}

func (c *controller) AppendRoutes(router *mux.Router) *mux.Router {
	staticFS := http.FileServer(http.FS(WebRoot))
	router.PathPrefix(c.basePath).Handler(http.StripPrefix(c.basePath, staticFS))

	return router
}
