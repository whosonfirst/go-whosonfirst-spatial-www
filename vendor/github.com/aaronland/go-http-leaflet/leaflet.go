package leaflet

import (
	"fmt"
	"github.com/aaronland/go-http-leaflet/static"
	"github.com/aaronland/go-http-rewrite"
	"io/fs"
	_ "log"
	"net/http"
	"path/filepath"
	"strings"
)

type LeafletOptions struct {
	JS  []string
	CSS []string
}

func DefaultLeafletOptions() *LeafletOptions {

	opts := &LeafletOptions{
		CSS: []string{
			"/css/leaflet.css",
			"/css/leaflet.fullscreen.css",
		},
		JS: []string{
			"/javascript/leaflet.js",
			"/javascript/leaflet.fullscreen.min.js",
			"/javascript/leaflet-hash.js",
		},
	}

	return opts
}

func AppendResourcesHandler(next http.Handler, opts *LeafletOptions) http.Handler {
	return AppendResourcesHandlerWithPrefix(next, opts, "")
}

func AppendResourcesHandlerWithPrefix(next http.Handler, opts *LeafletOptions, prefix string) http.Handler {

	js := opts.JS
	css := opts.CSS

	if prefix != "" {

		for i, path := range js {
			js[i] = appendPrefix(prefix, path)
		}

		for i, path := range css {
			css[i] = appendPrefix(prefix, path)
		}
	}

	ext_opts := &rewrite.AppendResourcesOptions{
		JavaScript:  js,
		Stylesheets: css,
	}

	return rewrite.AppendResourcesHandler(next, ext_opts)
}

func AssetsHandler() (http.Handler, error) {

	http_fs := http.FS(static.FS)
	return http.FileServer(http_fs), nil
}

func AssetsHandlerWithPrefix(prefix string) (http.Handler, error) {

	fs_handler, err := AssetsHandler()

	if err != nil {
		return nil, err
	}

	prefix = strings.TrimRight(prefix, "/")

	if prefix == "" {
		return fs_handler, nil
	}

	rewrite_func := func(req *http.Request) (*http.Request, error) {
		req.URL.Path = strings.Replace(req.URL.Path, prefix, "", 1)
		return req, nil
	}

	rewrite_handler := rewrite.RewriteRequestHandler(fs_handler, rewrite_func)
	return rewrite_handler, nil
}

func AppendAssetHandlers(mux *http.ServeMux) error {
	return AppendAssetHandlersWithPrefix(mux, "")
}

func AppendAssetHandlersWithPrefix(mux *http.ServeMux, prefix string) error {

	asset_handler, err := AssetsHandlerWithPrefix(prefix)

	if err != nil {
		return nil
	}

	walk_func := func(path string, info fs.DirEntry, err error) error {

		if path == "." {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if prefix != "" {
			path = appendPrefix(prefix, path)
		}

		if !strings.HasPrefix(path, "/") {
			path = fmt.Sprintf("/%s", path)
		}

		// log.Println("APPEND", path)

		mux.Handle(path, asset_handler)
		return nil
	}

	return fs.WalkDir(static.FS, ".", walk_func)
}

func appendPrefix(prefix string, path string) string {

	prefix = strings.TrimRight(prefix, "/")

	if prefix != "" {
		path = strings.TrimLeft(path, "/")
		path = filepath.Join(prefix, path)
	}

	return path
}
