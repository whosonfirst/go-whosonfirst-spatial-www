package flags

import (
	"errors"
	"flag"
	"fmt"
	"github.com/aaronland/go-http-tangramjs"
	"github.com/whosonfirst/go-whosonfirst-index"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
	"log"
	"os"
	"sort"
	"strings"
)

func Parse(fs *flag.FlagSet) {

	args := os.Args[1:]

	if len(args) > 0 && args[0] == "-h" {
		fs.Usage()
		os.Exit(0)
	}

	if len(args) > 0 && args[0] == "-setenv" {
		SetFromEnv(fs)
	}

	fs.Parse(args)
}

func SetFromEnv(fs *flag.FlagSet) {

	fs.VisitAll(func(fl *flag.Flag) {

		name := fl.Name
		env := name

		env = strings.ToUpper(env)
		env = strings.Replace(env, "-", "_", -1)
		env = fmt.Sprintf("WOF_%s", env)

		val, ok := os.LookupEnv(env)

		if ok {
			log.Printf("set -%s flag (%s) from %s environment variable\n", name, val, env)
			fs.Set(name, val)
		}

	})
}

func ValidateCommonFlags(fs *flag.FlagSet) error {

	_, err := StringVar(fs, "mode")

	if err != nil {
		return err
	}

	_, err = StringVar(fs, "spatial-database-uri")

	if err != nil {
		return err
	}

	enable_properties, err := BoolVar(fs, "enable-properties")

	if err != nil {
		return err
	}

	if enable_properties {

		properties_reader_uri, err := StringVar(fs, "properties-reader-uri")

		if err != nil {
			return err
		}

		if properties_reader_uri == "" {
			return errors.New("Invalid or missing -properties-reader-uri flag")
		}
	}

	return nil
}

func ValidateWWWFlags(fs *flag.FlagSet) error {

	enable_www, err := BoolVar(fs, "enable-www")

	if err != nil {
		return err
	}

	if !enable_www {
		return nil
	}

	log.Println("-enable-www flag is true causing the following flags to also be true: -enable-geojson -enable-candidates -enable-properties")

	fs.Set("enable-geojson", "true")
	fs.Set("enable-properties", "true")
	fs.Set("enable-candidates", "true")

	properties_reader_uri, err := StringVar(fs, "properties-reader-uri")

	if err != nil {
		return err
	}

	if properties_reader_uri == "" {
		return errors.New("Invalid or missing -properties-reader-uri flag")
	}

	_, err = StringVar(fs, "data-endpoint")

	if err != nil {
		return err
	}

	init_lat, err := Float64Var(fs, "initial-latitude")

	if err != nil {
		return err
	}

	if !geo.IsValidLatitude(init_lat) {
		return errors.New("Invalid latitude")
	}

	init_lon, err := Float64Var(fs, "initial-longitude")

	if err != nil {
		return err
	}

	if !geo.IsValidLongitude(init_lon) {
		return errors.New("Invalid longitude")
	}

	init_zoom, err := IntVar(fs, "initial-zoom")

	if err != nil {
		return err
	}

	if init_zoom < 1 {
		return errors.New("Invalid zoom")
	}

	return nil
}

func Lookup(fl *flag.FlagSet, k string) (interface{}, error) {

	v := fl.Lookup(k)

	if v == nil {
		msg := fmt.Sprintf("Unknown flag '%s'", k)
		return nil, errors.New(msg)
	}

	// Go is weird...
	return v.Value.(flag.Getter).Get(), nil
}

func StringVar(fl *flag.FlagSet, k string) (string, error) {

	i, err := Lookup(fl, k)

	if err != nil {
		return "", err
	}

	return i.(string), nil
}

func IntVar(fl *flag.FlagSet, k string) (int, error) {

	i, err := Lookup(fl, k)

	if err != nil {
		return 0, err
	}

	return i.(int), nil
}

func Float64Var(fl *flag.FlagSet, k string) (float64, error) {

	i, err := Lookup(fl, k)

	if err != nil {
		return 0, err
	}

	return i.(float64), nil
}

func BoolVar(fl *flag.FlagSet, k string) (bool, error) {

	i, err := Lookup(fl, k)

	if err != nil {
		return false, err
	}

	return i.(bool), nil
}

func NewFlagSet(name string) *flag.FlagSet {

	fs := flag.NewFlagSet(name, flag.ExitOnError)

	fs.Usage = func() {
		fs.PrintDefaults()
	}

	return fs
}

func CommonFlags() (*flag.FlagSet, error) {

	fs := NewFlagSet("common")

	fs.String("spatial-database-uri", "rtree://", "Valid options are: rtree://")

	fs.Bool("enable-properties", false, "Enable support for 'properties' parameters in queries.")
	fs.String("properties-reader-uri", "", "...")

	modes := index.Modes()
	modes = append(modes, "spatialite")

	sort.Strings(modes)

	valid_modes := strings.Join(modes, ", ")
	desc_modes := fmt.Sprintf("Valid modes are: %s.", valid_modes)

	fs.String("mode", "files", desc_modes)

	fs.Bool("is-wof", true, "Input data is WOF-flavoured GeoJSON. (Pass a value of '0' or 'false' if you need to index non-WOF documents.")

	fs.Bool("enable-custom-placetypes", false, "...")
	fs.String("custom-placetypes-source", "", "...")
	fs.String("custom-placetypes", "", "...")

	// this is invoked/used in app/indexer.go but for the life of me I can't
	// figure out how to make the code in flags/exclude.go implement the
	// correct inferface wah wah so that flag.Lookup("exclude").Value returns
	// something we can loop over... so instead we just strings.Split() on
	// flag.Lookup("exclude").String() which is dumb but works...
	// (20180301/thisisaaronland)

	var exclude Exclude
	fs.Var(&exclude, "exclude", "Exclude (WOF) records based on their existential flags. Valid options are: ceased, deprecated, not-current, superseded.")

	fs.Bool("setenv", false, "Set flags from environment variables.")
	fs.Bool("verbose", false, "Be chatty.")
	fs.Bool("strict", false, "Be strict about flags and fail if any are missing or deprecated flags are used.")

	return fs, nil
}

func AppendWWWFlags(fs *flag.FlagSet) error {

	fs.String("host", "localhost", "The hostname to listen for requests on.")
	fs.Int("port", 8080, "The port number to listen for requests on.")

	fs.Bool("enable-geojson", false, "Allow users to request GeoJSON FeatureCollection formatted responses.")
	fs.Bool("enable-candidates", false, "Enable the /candidates endpoint to return candidate bounding boxes (as GeoJSON) for requests.")
	fs.Bool("enable-www", false, "Enable the interactive /debug endpoint to query points and display results.")

	fs.String("www-path", "/debug", "The URL path for the interactive debug endpoint.")

	fs.String("static-prefix", "", "Prepend this prefix to URLs for static assets.")

	fs.String("nextzen-apikey", "", "A valid Nextzen API key")
	fs.String("nextzen-style-url", "/tangram/refill-style.zip", "...")
	fs.String("nextzen-tile-url", tangramjs.NEXTZEN_MVT_ENDPOINT, "...")

	fs.String("templates", "", "An optional string for local templates. This is anything that can be read by the 'templates.ParseGlob' method.")

	fs.Float64("initial-latitude", 37.616906, "...")
	fs.Float64("initial-longitude", -122.386665, "...")
	fs.Int("initial-zoom", 13, "...")

	fs.String("data-endpoint", "https://data.whosonfirst.org", "...")

	return nil
}
