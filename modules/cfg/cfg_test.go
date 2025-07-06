package cfg_test

import (
	"reflect"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/cfg"
)

const pathToRoot string = "../../"

// Expects the given cfg file to match expected values:
/*
	LogFileName = "test.log"
	LogMaxSize  = 10
	Verbose     = true
	DBfileName  = "test.db"
*/
func TestConfig(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			t.Fatal("panic", x)
		}
	}()
	expected := cfg.Config{
		LogFileName: "test.log",
		LogMaxSize:  10,
		Verbose:     true,
		DBfileName:  "test.db",
	}
	recieved := cfg.GetTest(pathToRoot)
	if !reflect.DeepEqual(expected, recieved) {
		t.Fatalf("recieved config is not what was expected; want: %v  |  got: %v", expected, recieved)
	}
}
