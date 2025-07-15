package cfg_test

import (
	"reflect"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/cfg"
)

const pathToRoot string = "../../"

// Expects the given cfg file to match expected values:
/*
	Verbose      = true

	LogFilePath  = "./logs/test.log"
	LogMaxSize   = 10

	DBfilePath   = "./database/test.db"

	WebStaticDir = "./static"

	BotTokenPath = "./configs/bot.key"
*/
func TestConfig(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			t.Fatal("panic", x)
		}
	}()
	expected := cfg.Config{
		Verbose: true,

		LogFilePath: "./logs/test.log",
		LogMaxSize:  10,

		DBfilePath: "./database/test.db",

		WebStaticDir: "./static",

		BotTokenPath: "./configs/bot.key",
	}
	recieved := cfg.GetTest(pathToRoot)
	if !reflect.DeepEqual(expected, recieved) {
		t.Fatalf("recieved config is not what was expected; want: %v  |  got: %v", expected, recieved)
	}
}
