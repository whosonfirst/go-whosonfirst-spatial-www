package pip

import (
	"context"
	"fmt"
	spatial_app "github.com/whosonfirst/go-whosonfirst-spatial/app"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

func QueryPointInPolygon(ctx context.Context, app *spatial_app.SpatialApplication, req *PointInPolygonRequest) (spr.StandardPlacesResults, error) {

	c, err := geo.NewCoordinate(req.Longitude, req.Latitude)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new coordinate, %v", err)
	}

	f, err := NewSPRFilterFromPointInPolygonRequest(req)

	if err != nil {
		return nil, err
	}

	db := app.SpatialDatabase
	return db.PointInPolygon(ctx, c, f)
}
