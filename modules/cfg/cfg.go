// Used to parse config file that was provided in [os.Args] as path.
// Then can be used globally to access constants (e.g. file paths).
package cfg

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

import (
	"flag"
	"fmt"
	"sync"

	"github.com/BurntSushi/toml"
)

var cfgPath *string = flag.String("config", "", "path to the config file, use 'go doc VladOS/modules/cfg.Config' for more details")

// Stores all necessary constants as struct, can  be accessed via [Get] function.
type Config struct {
	LogFileName string
	LogMaxSize  int // Megabytes
	Verbose     bool

	DBfileName string
}

// Parses config flag once via [flag] and [sync] packages,
// and decode given config file into [Config] struct via [github.com/BurntSushi/toml.DecodeFile] function.
//
// Panic if invalid config path is non existent or config file cannot be read.
func Get() Config {
	return sync.OnceValue(func() Config {
		if !flag.Parsed() {
			flag.Parse()
		}
		var cfg Config
		if _, err := toml.DecodeFile(*cfgPath, &cfg); err != nil {
			panic(fmt.Errorf("cannot read the config file: %w", err))
		}
		return cfg
	})()
}

// Acts fully as [Get] function but with difference of requiring a path to root directory.
// Used for testing purposes to still have configs in the same directory.
//
// Must be used with same pathToRoot values, since changing it after first call will not change the result.
func GetTest(pathToRoot string) Config {
	return sync.OnceValue(func() Config {
		if !flag.Parsed() {
			flag.Parse()
		}
		var cfg Config
		if _, err := toml.DecodeFile(pathToRoot+(*cfgPath), &cfg); err != nil {
			panic(fmt.Errorf("cannot read the config file: %w", err))
		}
		return cfg
	})()
}
