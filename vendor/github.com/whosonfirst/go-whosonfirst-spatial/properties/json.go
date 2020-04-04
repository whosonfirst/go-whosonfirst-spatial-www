package properties

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	_ "log"
	"strings"
)

func AppendPropertiesWithJSON(ctx context.Context, source []byte, target []byte, extras []string, prefix string) ([]byte, error) {

	var err error

	for _, e := range extras {

		paths := make([]string, 0)

		if strings.HasSuffix(e, "*") || strings.HasSuffix(e, ":") {

			e = strings.Replace(e, "*", "", -1)

			props := gjson.GetBytes(source, "properties")

			for k, _ := range props.Map() {

				if strings.HasPrefix(k, e) {
					paths = append(paths, k)
				}
			}

		} else {
			paths = append(paths, e)
		}

		for _, p := range paths {

			get_path := fmt.Sprintf("properties.%s", p)
			set_path := p

			if prefix != "" {
				set_path = fmt.Sprintf("%s.%s", prefix, p)
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
