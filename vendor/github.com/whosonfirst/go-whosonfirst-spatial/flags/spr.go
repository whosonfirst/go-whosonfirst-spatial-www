package flags

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"net/url"
)

func NewSPRFilterFromFlagSet(fs *flag.FlagSet) (filter.Filter, error) {

	placetypes, err := MultiStringVar(fs, "placetype")

	if err != nil {
		return nil, err
	}

	geometries, err := StringVar(fs, "geometries")

	if err != nil {
		return nil, err
	}

	alt_geoms, err := MultiStringVar(fs, "alternate-geometry")

	if err != nil {
		return nil, err
	}

	is_current, err := MultiStringVar(fs, "is-current")

	if err != nil {
		return nil, err
	}

	is_ceased, err := MultiStringVar(fs, "is-ceased")

	if err != nil {
		return nil, err
	}

	is_deprecated, err := MultiStringVar(fs, "is-deprecated")

	if err != nil {
		return nil, err
	}

	is_superseded, err := MultiStringVar(fs, "is-superseded")

	if err != nil {
		return nil, err
	}

	is_superseding, err := MultiStringVar(fs, "is-superseding")

	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Set("geometries", geometries)

	for _, v := range alt_geoms {
		q.Add("alternate_geometry", v)
	}

	for _, v := range placetypes {
		q.Add("placetype", v)
	}

	for _, v := range is_current {
		q.Add("is_current", v)
	}

	for _, v := range is_ceased {
		q.Add("is_ceased", v)
	}

	for _, v := range is_deprecated {
		q.Add("is_deprecated", v)
	}

	for _, v := range is_superseded {
		q.Add("is_superseded", v)
	}

	for _, v := range is_superseding {
		q.Add("is_superseding", v)
	}

	return filter.NewSPRFilterFromQuery(q)
}
