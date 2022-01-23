package gin

import (
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	G     *gin.Engine
	Debug bool
}

func NewHTTPServer(debug bool) *HTTPServer {
	return &HTTPServer{Debug: debug}
}

func (h *HTTPServer) LoadRoute(route func(e *gin.Engine)) error {
	if h.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	if err := LoadValidatorLocal("zh"); err != nil {
		return err
	}
	g := gin.New()
	route(g)
	h.G = g
	return nil
}

func (h *HTTPServer) Run(addr string) error {
	return h.G.Run(addr)
}
