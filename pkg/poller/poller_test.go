package poller

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func TestParseNetworkMetrics(t *testing.T) {
	json, err := ioutil.ReadFile("../../test/stats-summary-1.json")
	check(err)
	s := ParseNetworkMetrics(json)
	fmt.Println(s)
}
