package pip

import (
	"context"
	"fmt"
	spatial_app "github.com/whosonfirst/go-whosonfirst-spatial/app"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

func QueryPointInPolygon(ctx context.Context, app *spatial_app.SpatialApplication, req *PointInPolygonRequest) (interface{}, error) {

	c, err := geo.NewCoordinate(req.Longitude, req.Latitude)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new coordinate, %v", err)
	}

	f, err := NewSPRFilterFromPointInPolygonRequest(req)

	if err != nil {
		return nil, err
	}

	db := app.SpatialDatabase
	pr := app.PropertiesReader

	var rsp interface{}

	r, err := db.PointInPolygon(ctx, c, f)

	if err != nil {
		return nil, fmt.Errorf("Failed to query database with coord %v, %v", c, err)
	}

	rsp = r

	if pr != nil && len(req.Properties) > 0 {

		r, err := pr.PropertiesResponseResultsWithStandardPlacesResults(ctx, rsp.(spr.StandardPlacesResults), req.Properties)

		if err != nil {
			return nil, fmt.Errorf("Failed to generate properties response, %v", err)
		}

		rsp = r
	}

	return r, nil
}
