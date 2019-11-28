package main

import (
	"fmt"
	"github.com/samos123/k8s-dynamic-bandwidth-qos/pkg/poller"
	"time"
)

func main() {
	c := time.Tick(10 * time.Second)
	for _ = range c {
		url := "http://localhost:10255/stats/summary"
		fmt.Println("Getting metrics from URL:", url)
		json := poller.GetMetrics(url)
		metrics := poller.ParseNetworkMetrics(json)
		fmt.Println(metrics)
	}
	fmt.Println("Exitting")
}
