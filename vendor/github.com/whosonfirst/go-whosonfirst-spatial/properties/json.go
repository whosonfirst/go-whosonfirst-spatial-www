package properties

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	_ "log"
	"strings"
)

type AppendPropertiesOptions struct {
	SourcePrefix string
	TargetPrefix string
	Keys         []string
}

func AppendPropertiesWithJSON(ctx context.Context, opts *AppendPropertiesOptions, source []byte, target []byte) ([]byte, error) {

	var err error

	for _, e := range opts.Keys {

		paths := make([]string, 0)

		if strings.HasSuffix(e, "*") || strings.HasSuffix(e, ":") {

			e = strings.Replace(e, "*", "", -1)

			var props gjson.Result

			if opts.SourcePrefix != "" {
				props = gjson.GetBytes(source, opts.SourcePrefix)
			} else {
				props = gjson.ParseBytes(source)
			}

			for k, _ := range props.Map() {

				if strings.HasPrefix(k, e) {
					paths = append(paths, k)
				}
			}

		} else {
			paths = append(paths, e)
		}

		for _, p := range paths {

			get_path := p
			set_path := p

			if opts.SourcePrefix != "" {
				get_path = fmt.Sprintf("%s.%s", opts.SourcePrefix, get_path)
			}

			if opts.TargetPrefix != "" {
				set_path = fmt.Sprintf("%s.%s", opts.TargetPrefix, p)
			}

			v := gjson.GetBytes(source, get_path)

			/*
				log.Println("GET", get_path)
				log.Println("SET", set_path)
				log.Println("VALUE", v.Value())
			*/

			if !v.Exists() {
				continue
			}

			target, err = sjson.SetBytes(target, set_path, v.Value())

			if err != nil {
				return nil, err
			}
		}
	}

	return target, nil
}
