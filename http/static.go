package http

import (
	"github.com/aaronland/go-http-rewrite"
	"github.com/whosonfirst/go-whosonfirst-spatial-http/static"
	gohttp "net/http"
	"path/filepath"
	"strings"
	"io/fs"
)

func StaticFileSystem() (gohttp.FileSystem, error) {
	fs := assetFS()
	return fs, nil
}

func StaticAssetsHandler() (gohttp.Handler, error) {

	http_fs := http.FS(static.FS)
	return http.FileServer(http_fs), nil
}

func StaticAssetsHandlerWithPrefix(prefix string) (gohttp.Handler, error) {

	fs_handler, err := StaticAssetsHandler()

	if err != nil {
		return nil, err
	}

	prefix = strings.TrimRight(prefix, "/")

	if prefix == "" {
		return fs_handler, nil
	}

	rewrite_func := func(req *gohttp.Request) (*gohttp.Request, error) {
		req.URL.Path = strings.Replace(req.URL.Path, prefix, "", 1)
		return req, nil
	}

	rewrite_handler := rewrite.RewriteRequestHandler(fs_handler, rewrite_func)
	return rewrite_handler, nil
}

func AppendStaticAssetHandlers(mux *gohttp.ServeMux) error {
	return AppendStaticAssetHandlersWithPrefix(mux, "")
}

func AppendStaticAssetHandlersWithPrefix(mux *gohttp.ServeMux, prefix string) error {

	asset_handler, err := StaticAssetsHandlerWithPrefix(prefix)

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
