package rtree

import (
	"context"
	"errors"
	"fmt"
	"github.com/dhconnelly/rtreego"
	gocache "github.com/patrickmn/go-cache"
	pm_geojson "github.com/paulmach/go.geojson"
	"github.com/skelterjohn/geom"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/geometry"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-spatial"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
	"github.com/whosonfirst/go-whosonfirst-spatial/timer"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	"net/url"
	"strconv"
	"sync"
	"time"
)

func init() {
	ctx := context.Background()
	database.RegisterSpatialDatabase(ctx, "rtree", NewRTreeSpatialDatabase)
}

type RTreeCache struct {
	Geometry *pm_geojson.Geometry     `json:"geometry"`
	SPR      spr.StandardPlacesResult `json:"properties"`
}

// PLEASE DISCUSS WHY patrickm/go-cache AND NOT whosonfirst/go-cache HERE

type RTreeSpatialDatabase struct {
	database.SpatialDatabase
	Logger          *log.WOFLogger
	Timer           *timer.Timer
	index_alt_files bool
	rtree           *rtreego.Rtree
	gocache         *gocache.Cache
	mu              *sync.RWMutex
	strict          bool
}

type RTreeSpatialIndex struct {
	Rect      *rtreego.Rect
	Id        string
	FeatureId string
	IsAlt     bool
	AltLabel  string
}

func (i *RTreeSpatialIndex) Bounds() *rtreego.Rect {
	return i.Rect
}

type RTreeResults struct {
	spr.StandardPlacesResults `json:",omitempty"`
	Places                    []spr.StandardPlacesResult `json:"places"`
}

func (r *RTreeResults) Results() []spr.StandardPlacesResult {
	return r.Places
}

func NewRTreeSpatialDatabase(ctx context.Context, uri string) (database.SpatialDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	strict := true

	if q.Get("strict") == "false" {
		strict = false
	}

	expires := 0 * time.Second
	cleanup := 0 * time.Second

	str_exp := q.Get("default_expiration")
	str_cleanup := q.Get("cleanup_interval")

	if str_exp != "" {

		int_expires, err := strconv.Atoi(str_exp)

		if err != nil {
			return nil, err
		}

		expires = time.Duration(int_expires) * time.Second
	}

	if str_cleanup != "" {

		int_cleanup, err := strconv.Atoi(str_cleanup)

		if err != nil {
			return nil, err
		}

		cleanup = time.Duration(int_cleanup) * time.Second
	}

	index_alt_files := false

	str_index_alt := q.Get("index_alt_files")

	if str_index_alt != "" {

		index_alt, err := strconv.ParseBool(str_index_alt)

		if err != nil {
			return nil, err
		}

		index_alt_files = index_alt
	}

	gc := gocache.New(expires, cleanup)

	logger := log.SimpleWOFLogger("index")

	rtree := rtreego.NewTree(2, 25, 50)

	mu := new(sync.RWMutex)

	t := timer.NewTimer()

	db := &RTreeSpatialDatabase{
		Logger:          logger,
		Timer:           t,
		rtree:           rtree,
		index_alt_files: index_alt_files,
		gocache:         gc,
		strict:          strict,
		mu:              mu,
	}

	return db, nil
}

func (r *RTreeSpatialDatabase) Close(ctx context.Context) error {
	return nil
}

func (r *RTreeSpatialDatabase) IndexFeature(ctx context.Context, f wof_geojson.Feature) error {

	err := r.setCache(ctx, f)

	if err != nil {
		return err
	}

	is_alt := whosonfirst.IsAlt(f)
	alt_label := whosonfirst.AltLabel(f)

	if is_alt && !r.index_alt_files {
		return nil
	}

	if is_alt && alt_label == "" {
		return errors.New("Invalid alt label")
	}

	feature_id := f.Id()

	bboxes, err := f.BoundingBoxes()

	if err != nil {
		return err
	}

	for i, bbox := range bboxes.Bounds() {

		sp_id, err := spatial.SpatialIdWithFeature(f, i)

		if err != nil {
			return err
		}

		sw := bbox.Min
		ne := bbox.Max

		llat := ne.Y - sw.Y
		llon := ne.X - sw.X

		pt := rtreego.Point{sw.X, sw.Y}
		rect, err := rtreego.NewRect(pt, []float64{llon, llat})

		if err != nil {

			if r.strict {
				return err
			}

			r.Logger.Error("%s failed indexing, (%v). Strict mode is disabled, so skipping.", sp_id, err)
			return nil
		}

		r.Logger.Status("index %s %v", sp_id, rect)

		sp := RTreeSpatialIndex{
			Rect:      rect,
			Id:        sp_id,
			FeatureId: feature_id,
			IsAlt:     is_alt,
			AltLabel:  alt_label,
		}

		r.mu.Lock()
		r.rtree.Insert(&sp)
		r.mu.Unlock()
	}

	return nil
}

func (r *RTreeSpatialDatabase) PointInPolygon(ctx context.Context, coord *geom.Coord, filters ...spatial.Filter) (spr.StandardPlacesResults, error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	rsp_ch := make(chan spr.StandardPlacesResult)
	err_ch := make(chan error)
	done_ch := make(chan bool)

	results := make([]spr.StandardPlacesResult, 0)
	working := true

	go r.PointInPolygonWithChannels(ctx, rsp_ch, err_ch, done_ch, coord, filters...)

	for {
		select {
		case <-ctx.Done():
			return nil, nil
		case <-done_ch:
			working = false
		case rsp := <-rsp_ch:
			results = append(results, rsp)
		case err := <-err_ch:
			return nil, err
		default:
			// pass
		}

		if !working {
			break
		}
	}

	spr_results := &RTreeResults{
		Places: results,
	}

	return spr_results, nil
}

func (r *RTreeSpatialDatabase) PointInPolygonWithChannels(ctx context.Context, rsp_ch chan spr.StandardPlacesResult, err_ch chan error, done_ch chan bool, coord *geom.Coord, filters ...spatial.Filter) {

	defer func() {
		done_ch <- true
	}()

	rows, err := r.getIntersectsByCoord(coord)

	if err != nil {
		err_ch <- err
		return
	}

	r.inflateResultsWithChannels(ctx, rsp_ch, err_ch, rows, coord, filters...)
	return
}

func (r *RTreeSpatialDatabase) PointInPolygonCandidates(ctx context.Context, coord *geom.Coord, filters ...spatial.Filter) ([]*spatial.PointInPolygonCandidate, error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	rsp_ch := make(chan *spatial.PointInPolygonCandidate)
	err_ch := make(chan error)
	done_ch := make(chan bool)

	candidates := make([]*spatial.PointInPolygonCandidate, 0)
	working := true

	go r.PointInPolygonCandidatesWithChannels(ctx, rsp_ch, err_ch, done_ch, coord, filters...)

	for {
		select {
		case <-ctx.Done():
			return nil, nil
		case <-done_ch:
			working = false
		case rsp := <-rsp_ch:
			candidates = append(candidates, rsp)
		case err := <-err_ch:
			return nil, err
		default:
			// pass
		}

		if !working {
			break
		}
	}

	return candidates, nil
}

func (r *RTreeSpatialDatabase) PointInPolygonCandidatesWithChannels(ctx context.Context, rsp_ch chan *spatial.PointInPolygonCandidate, err_ch chan error, done_ch chan bool, coord *geom.Coord, filters ...spatial.Filter) {

	defer func() {
		done_ch <- true
	}()

	intersects, err := r.getIntersectsByCoord(coord)

	if err != nil {
		err_ch <- err
		return
	}

	for _, raw := range intersects {

		sp := raw.(*RTreeSpatialIndex)

		// bounds := sp.Bounds()

		c := &spatial.PointInPolygonCandidate{
			Id:        sp.Id,
			FeatureId: sp.FeatureId,
			AltLabel:  sp.AltLabel,
			// FIX ME
			// Bounds:   bounds,
		}

		rsp_ch <- c
	}

	return
}

func (r *RTreeSpatialDatabase) getIntersectsByCoord(coord *geom.Coord) ([]rtreego.Spatial, error) {

	lat := coord.Y
	lon := coord.X

	pt := rtreego.Point{lon, lat}
	rect, err := rtreego.NewRect(pt, []float64{0.0001, 0.0001}) // how small can I make this?

	if err != nil {
		return nil, err
	}

	return r.getIntersectsByRect(rect)
}

func (r *RTreeSpatialDatabase) getIntersectsByRect(rect *rtreego.Rect) ([]rtreego.Spatial, error) {

	// to do: timings that don't slow everything down the way
	// go-whosonfirst-timer does now (20170915/thisisaaronland)

	results := r.rtree.SearchIntersect(rect)
	return results, nil
}

func (r *RTreeSpatialDatabase) inflateResultsWithChannels(ctx context.Context, rsp_ch chan spr.StandardPlacesResult, err_ch chan error, possible []rtreego.Spatial, c *geom.Coord, filters ...spatial.Filter) {

	seen := make(map[string]bool)

	mu := new(sync.RWMutex)
	wg := new(sync.WaitGroup)

	for _, row := range possible {

		sp := row.(*RTreeSpatialIndex)
		wg.Add(1)

		go func(sp *RTreeSpatialIndex) {

			sp_id := sp.Id
			feature_id := sp.FeatureId

			t1 := time.Now()

			defer func() {
				r.Timer.Add(ctx, sp_id, "time to inflate", time.Since(t1))
			}()

			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
				// pass
			}

			mu.RLock()
			_, ok := seen[feature_id]
			mu.RUnlock()

			if ok {
				return
			}

			mu.Lock()
			seen[feature_id] = true
			mu.Unlock()

			t2 := time.Now()

			cache_item, err := r.retrieveCache(ctx, sp)

			r.Timer.Add(ctx, sp_id, "time to retrieve cache", time.Since(t2))

			if err != nil {
				r.Logger.Error("Failed to retrieve cache for %s, %v", sp_id, err)
				return
			}

			s := cache_item.SPR

			t3 := time.Now()

			for _, f := range filters {

				err = filter.FilterSPR(f, s)

				if err != nil {
					r.Logger.Debug("SKIP %s because filter error %s", sp_id, err)
					return
				}
			}

			r.Timer.Add(ctx, sp_id, "time to filter", time.Since(t3))

			t4 := time.Now()

			geom := cache_item.Geometry

			contains := false

			switch geom.Type {
			case "Polygon":
				contains = geo.PolygonContainsCoord(geom.Polygon, c)
			case "MultiPolygon":
				contains = geo.MultiPolygonContainsCoord(geom.MultiPolygon, c)
			default:
				r.Logger.Warning("Geometry has unsupported geometry type '%s'", geom.Type)
			}

			r.Timer.Add(ctx, sp_id, "time to test geometry", time.Since(t4))

			if !contains {
				r.Logger.Debug("SKIP %s because does not contain coord (%v)", sp_id, c)
				return
			}

			rsp_ch <- s
		}(sp)
	}

	wg.Wait()
}

func (r *RTreeSpatialDatabase) setCache(ctx context.Context, f wof_geojson.Feature) error {

	s, err := f.SPR()

	if err != nil {
		return err
	}

	geom, err := geometry.GeometryForFeature(f)

	if err != nil {
		return err
	}

	alt_label := whosonfirst.AltLabel(f)

	feature_id := f.Id()

	cache_key := fmt.Sprintf("%s:%s", feature_id, alt_label)

	cache_item := &RTreeCache{
		Geometry: geom,
		SPR:      s,
	}

	r.gocache.Set(cache_key, cache_item, -1)
	return nil
}

func (r *RTreeSpatialDatabase) retrieveCache(ctx context.Context, sp *RTreeSpatialIndex) (*RTreeCache, error) {

	cache_key := fmt.Sprintf("%s:%s", sp.FeatureId, sp.AltLabel)

	cache_item, ok := r.gocache.Get(cache_key)

	if !ok {
		return nil, errors.New("Invalid cache ID")
	}

	return cache_item.(*RTreeCache), nil
}
