package flags

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-index"
	"sort"
	"strings"
)

func AppendIndexingFlags(fs *flag.FlagSet) error {

	modes := index.Modes()
	sort.Strings(modes)

	valid_modes := strings.Join(modes, ", ")
	desc_modes := fmt.Sprintf("Valid modes are: %s.", valid_modes)

	fs.String("mode", "repo://", desc_modes)

	return nil
}

func ValidateIndexingFlags(fs *flag.FlagSet) error {

	_, err := StringVar(fs, "mode")

	if err != nil {
		return err
	}

	return nil
}
