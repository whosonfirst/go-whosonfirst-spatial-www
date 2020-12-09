package flags

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func NewFlagSet(name string) *flag.FlagSet {

	fs := flag.NewFlagSet(name, flag.ExitOnError)

	fs.Usage = func() {
		fs.PrintDefaults()
	}

	return fs
}

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
