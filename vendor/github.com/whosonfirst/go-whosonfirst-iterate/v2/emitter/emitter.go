package emitter

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	"io"
	"net/url"
	"os"
	"sort"
	"strings"
)

const STDIN string = "STDIN"

type IndexerContextKey string

type EmitterInitializeFunc func(context.Context, string) (Emitter, error)

type EmitterCallbackFunc func(context.Context, string, io.ReadSeeker, ...interface{}) error

type Emitter interface {
	WalkURI(context.Context, EmitterCallbackFunc, string) error
}

var emitters roster.Roster

func ensureSpatialRoster() error {

	if emitters == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		emitters = r
	}

	return nil
}

func RegisterEmitter(ctx context.Context, scheme string, f EmitterInitializeFunc) error {

	err := ensureSpatialRoster()

	if err != nil {
		return err
	}

	return emitters.Register(ctx, scheme, f)
}

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureSpatialRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range emitters.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func NewEmitter(ctx context.Context, uri string) (Emitter, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := emitters.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	fn := i.(EmitterInitializeFunc)
	return fn(ctx, uri)
}

//

func ReaderWithPath(ctx context.Context, abs_path string) (io.ReadSeekCloser, error) {

	if abs_path == STDIN {
		return os.Stdin, nil
	}

	fh, err := os.Open(abs_path)

	if err != nil {
		return nil, err
	}

	return fh, nil
}
