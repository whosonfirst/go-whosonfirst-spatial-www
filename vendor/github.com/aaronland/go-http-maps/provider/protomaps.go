package provider

import (
	_ "gocloud.dev/blob/fileblob"
)

import (
	"context"
	"fmt"
	"github.com/aaronland/go-http-leaflet"
	"github.com/protomaps/go-pmtiles/pmtiles"
	"github.com/sfomuseum/go-http-protomaps"
	pmhttp "github.com/sfomuseum/go-sfomuseum-pmtiles/http"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const PROTOMAPS_SCHEME string = "protomaps"

type ProtomapsProvider struct {
	Provider
	leafletOptions   *leaflet.LeafletOptions
	protomapsOptions *protomaps.ProtomapsOptions
	logger           *log.Logger
	serve_tiles      bool
	cache_size       int
	bucket_uri       string
	path_tiles       string
	database         string
}

func init() {
	protomaps.APPEND_LEAFLET_RESOURCES = false
	protomaps.APPEND_LEAFLET_ASSETS = false

	ctx := context.Background()
	RegisterProvider(ctx, PROTOMAPS_SCHEME, NewProtomapsProvider)
}

func ProtomapsOptionsFromURL(u *url.URL) (*protomaps.ProtomapsOptions, error) {
	opts := protomaps.DefaultProtomapsOptions()
	return opts, nil
}

func NewProtomapsProvider(ctx context.Context, uri string) (Provider, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	leaflet_opts, err := LeafletOptionsFromURL(u)

	if err != nil {
		return nil, fmt.Errorf("Failed to create leaflet options, %w", err)
	}

	protomaps_opts, err := ProtomapsOptionsFromURL(u)

	if err != nil {
		return nil, fmt.Errorf("Failed to create protomaps options, %w", err)
	}

	q := u.Query()

	q_tile_url := q.Get(ProtomapsTileURLFlag)
	protomaps_opts.TileURL = q_tile_url

	logger := log.New(io.Discard, "", 0)

	p := &ProtomapsProvider{
		leafletOptions:   leaflet_opts,
		protomapsOptions: protomaps_opts,
		logger:           logger,
	}

	serve_tiles := false

	q_serve_tiles := q.Get(ProtomapsServeTilesFlag)

	if q_serve_tiles != "" {

		v, err := strconv.ParseBool(q_serve_tiles)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?protomaps-serve-tiles= parameter, %w")
		}

		serve_tiles = v
	}

	if serve_tiles {

		q_cache_size := q.Get(ProtomapsCacheSizeFlag)
		q_bucket_uri := q.Get(ProtomapsBucketURIFlag)
		q_database := q.Get(ProtomapsDatabaseFlag)

		sz, err := strconv.Atoi(q_cache_size)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?%s= parameter, %w", ProtomapsCacheSizeFlag, err)
		}

		p.cache_size = sz
		p.bucket_uri = q_bucket_uri
		p.database = q_database
		p.path_tiles = q_tile_url
		p.serve_tiles = true
	}

	return p, nil
}

func (p *ProtomapsProvider) Scheme() string {
	return PROTOMAPS_SCHEME
}

func (p *ProtomapsProvider) AppendResourcesHandler(handler http.Handler) http.Handler {
	handler = leaflet.AppendResourcesHandler(handler, p.leafletOptions)
	handler = protomaps.AppendResourcesHandler(handler, p.protomapsOptions)
	return handler
}

func (p *ProtomapsProvider) AppendAssetHandlers(mux *http.ServeMux) error {

	err := leaflet.AppendAssetHandlers(mux)

	if err != nil {
		return fmt.Errorf("Failed to append leaflet asset handler, %w", err)
	}

	err = protomaps.AppendAssetHandlers(mux)

	if err != nil {
		return fmt.Errorf("Failed to append protomaps asset handler, %w", err)
	}

	if p.serve_tiles {

		loop, err := pmtiles.NewLoop(p.bucket_uri, p.logger, p.cache_size, "")

		if err != nil {
			return fmt.Errorf("Failed to create pmtiles.Loop, %w", err)
		}

		loop.Start()

		pmtiles_handler := pmhttp.TileHandler(loop, p.logger)

		strip_path := strings.TrimRight(p.path_tiles, "/")
		pmtiles_handler = http.StripPrefix(strip_path, pmtiles_handler)

		log.Println("ADD", p.path_tiles)
		mux.Handle(p.path_tiles, pmtiles_handler)

		// Because inevitably I will forget...
		protomaps_tiles_database := strings.Replace(p.database, ".pmtiles", "", 1)

		pm_tile_url, err := url.JoinPath(p.path_tiles, protomaps_tiles_database)

		if err != nil {
			return fmt.Errorf("Failed to join path to derive Protomaps tile URL, %w", err)
		}

		pm_tile_url = fmt.Sprintf("%s/{z}/{x}/{y}.mvt", pm_tile_url)

		p.protomapsOptions.TileURL = pm_tile_url
	}

	return nil
}

func (p *ProtomapsProvider) SetLogger(logger *log.Logger) error {
	p.logger = logger
	return nil
}
