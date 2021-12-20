package gin

import (
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	Debug bool
}

func NewHTTPServer(debug bool) *HTTPServer {
	return &HTTPServer{Debug: debug}
}

func (h *HTTPServer) HTTP(addr string, route func(e *gin.Engine)) error {
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
	return g.Run(addr)
}
