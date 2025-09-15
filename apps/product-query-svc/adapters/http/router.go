package httpadapter

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(si ServerInterface, validator func(http.Handler) http.Handler) http.Handler {
	r := chi.NewRouter()
	// 你可以在这儿挂日志、中间件、recover等
	opts := ChiServerOptions{}
	if validator != nil {
		opts.Middlewares = []MiddlewareFunc{validator}
	}
	// 使用生成器的 HandlerWithOptions 挂载路由
	r.Mount("/", HandlerWithOptions(si, opts))
	return r
}
