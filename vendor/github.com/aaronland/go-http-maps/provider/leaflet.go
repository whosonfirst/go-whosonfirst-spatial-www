package provider

import (
	"fmt"
	"github.com/aaronland/go-http-leaflet"
	"net/url"
	"strconv"
)

func LeafletOptionsFromURL(u *url.URL) (*leaflet.LeafletOptions, error) {

	opts := leaflet.DefaultLeafletOptions()

	q := u.Query()

	q_enable_hash := q.Get("leaflet-enable-hash")

	if q_enable_hash != "" {

		v, err := strconv.ParseBool(q_enable_hash)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?leaflet-enable-hash= parameter, %w", err)
		}

		if v == true {
			opts.EnableHash()
		}
	}

	q_enable_fullscreen := q.Get("leaflet-enable-fullscreen")

	if q_enable_fullscreen != "" {

		v, err := strconv.ParseBool(q_enable_fullscreen)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?leaflet-enable-fullscreen= parameter, %w", err)
		}

		if v == true {
			opts.EnableFullscreen()
		}
	}

	q_enable_draw := q.Get("leaflet-enable-draw")

	if q_enable_draw != "" {

		v, err := strconv.ParseBool(q_enable_draw)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?leaflet-enable-draw= parameter, %w", err)
		}

		if v == true {
			opts.EnableDraw()
		}
	}

	return opts, nil
}
