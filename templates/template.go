package templates

import (
	"errors"
	"io"
	"slingshot/config"

	"github.com/flosch/pongo2/v6"
	"github.com/labstack/echo/v4"
)

// TODO: pongo2 template engine
// TODO: pongo2 template embed

// We need to wrap the renderer because we need a different signature for echo.
type RenderWrapper struct {
	typ       string
	debugMode bool
	baseDir   string
	template  *pongo2.Template
}

func (r *RenderWrapper) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	templatesFS := pongo2.MustNewLocalFileSystemLoader(r.baseDir)
	templateSet := pongo2.NewSet(r.typ, templatesFS)

	templateSet.Debug = r.debugMode

	var err error
	// // If in debug mode, always use FromFile, else use FromCache
	// if r.debugMode {
	// 	r.template, err = templateSet.FromFile(name)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }

	r.template, err = templateSet.FromCache(name)
	if err != nil {
		return err
	}

	d, ok := data.(map[string]interface{})
	if !ok {
		return errors.New("incorrect data format. Should be map[string]interface{}")
	}

	err = r.template.ExecuteWriter(pongo2.Context{"data": d}, w)
	if err != nil {
		return err
	}

	return nil
}

var HTMLRender = &RenderWrapper{
	typ:       "html",
	debugMode: config.Cfg.Server.Debug,
	baseDir:   "templates/templates",
}
