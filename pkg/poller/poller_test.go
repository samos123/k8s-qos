package poller

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func TestParseNetworkMetrics(t *testing.T) {
	json, err := ioutil.ReadFile("../../test/stats-summary-1.json")
	check(err)
	s := ParseNetworkMetrics(json)
	equals(t, s[0].Name, "network-metering-agent-bljqp")
	fmt.Println(s)
}
