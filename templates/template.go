package templates

import (
	"embed"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/unrolled/render"
)

//go:embed templates
var templatesFS embed.FS

//go:embed assets
var assetsFS embed.FS

var layoutHTMLRender = render.New(render.Options{
	Directory: "templates",
	FileSystem: &render.EmbedFileSystem{
		FS: templatesFS,
	},
	DisableHTTPErrorRendering: true,
	Layout:                    "layout/layout",
	Extensions:                []string{".html"},
	IndentJSON:                true,
})

var rawRender = render.New(render.Options{
	Directory: "templates",
	FileSystem: &render.EmbedFileSystem{
		FS: templatesFS,
	},
	DisableHTTPErrorRendering: true,
	Extensions:                []string{".html"},
	IndentJSON:                true,
})

type RenderWrapper struct { // We need to wrap the renderer because we need a different signature for echo.
	rnd *render.Render
}

func (r *RenderWrapper) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.rnd.HTML(w, 0, name, data) // The zero status code is overwritten by echo.
}

var HTMLRender = &RenderWrapper{rnd: layoutHTMLRender}
var RawHTMLRender = &RenderWrapper{rnd: rawRender}
