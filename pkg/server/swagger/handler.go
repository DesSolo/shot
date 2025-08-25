package swagger

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/swaggest/swgui"
	"github.com/swaggest/swgui/v5emb"
)

const (
	swaggerPath     = "/docs/"
	swaggerJSONPath = swaggerPath + "apidocs.swagger.json"

	interceptorJSTemplate = `(req) => {
		if (req.loadSpec) return req;
		const url = new URL(req.url, window.location.origin);
		url.hostname = window.location.hostname;
		url.port = %s;
		url.protocol = window.location.protocol;
		req.url = url.toString();
		return req;
	}`
)

// Init ...
func Init(mux *chi.Mux, port string) {
	swaggerUI := v5emb.NewHandlerWithConfig(swgui.Config{
		Title:       "API Documentation",
		SwaggerJSON: swaggerJSONPath,
		BasePath:    swaggerPath,
		SettingsUI: map[string]string{
			"requestInterceptor": fmt.Sprintf(interceptorJSTemplate, port),
		},
	})

	mux.Handle(swaggerPath+"*", swaggerUI)

	mux.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, swaggerPath, http.StatusFound)
	})

	mux.Handle(swaggerJSONPath,
		http.StripPrefix(swaggerPath,
			http.FileServer(http.Dir("."+swaggerPath)),
		),
	)
}
