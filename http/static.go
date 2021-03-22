package http

import (
	"fmt"
	"github.com/aaronland/go-http-rewrite"
	"github.com/whosonfirst/go-whosonfirst-spatial-www/static"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"io/fs"
	"log"
	gohttp "net/http"
	"path/filepath"
	"strings"
)

func StaticAssetsHandler() (gohttp.Handler, error) {
	http_fs := gohttp.FS(static.FS)
	return gohttp.FileServer(http_fs), nil
}

func StaticAssetsHandlerWithPrefix(prefix string) (gohttp.Handler, error) {

	fs_handler, err := StaticAssetsHandler()

	if err != nil {
		return nil, err
	}

	fs_handler = gohttp.StripPrefix(prefix, fs_handler)
	return fs_handler, nil
}

func AppendStaticResourcesHandler(next gohttp.Handler) gohttp.Handler {
	return AppendStaticResourcesHandlerWithPrefix(next, "")
}

func AppendStaticResourcesHandlerWithPrefix(next gohttp.Handler, prefix string) gohttp.Handler {

	// Only do this once...

	stylesheets := make([]string, 0)
	javascript := make([]string, 0)

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

		switch filepath.Ext(path) {
		case ".css":
			stylesheets = append(stylesheets, path)
		case ".js":
			javascript = append(javascript, path)
		default:
			// pass
		}

		return nil
	}

	err := fs.WalkDir(static.FS, ".", walk_func)

	if err != nil {
		log.Println(err)
	}

	var cb rewrite.RewriteHTMLFunc

	cb = func(n *html.Node, w io.Writer) {

		if n.Type == html.ElementNode && n.Data == "head" {

			for _, js := range javascript {

				script_type := html.Attribute{"", "type", "text/javascript"}
				script_src := html.Attribute{"", "src", js}

				script := html.Node{
					Type:      html.ElementNode,
					DataAtom:  atom.Script,
					Data:      "script",
					Namespace: "",
					Attr:      []html.Attribute{script_type, script_src},
				}

				n.AppendChild(&script)
			}

			for _, css := range stylesheets {

				link_type := html.Attribute{"", "type", "text/css"}
				link_rel := html.Attribute{"", "rel", "stylesheet"}
				link_href := html.Attribute{"", "href", css}

				link := html.Node{
					Type:      html.ElementNode,
					DataAtom:  atom.Link,
					Data:      "link",
					Namespace: "",
					Attr:      []html.Attribute{link_type, link_rel, link_href},
				}

				n.AppendChild(&link)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			cb(c, w)
		}
	}

	return rewrite.RewriteHTMLHandler(next, cb)
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

		log.Println("APPEND", path)

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
