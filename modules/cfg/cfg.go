package cfg

import (
	"flag"
	"fmt"
	"sync"

	"github.com/BurntSushi/toml"
)

var cfgPath *string = flag.String("config", "", "path to the config file, check [vlados/modules/cfg.Config] for more config details")

type Config struct {
	LogFilePath string
	LogMaxSize  int // Megabytes
	Verbose     bool
}

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
