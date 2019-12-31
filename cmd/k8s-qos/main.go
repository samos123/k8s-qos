package main

import (
	"flag"
	"fmt"
	"github.com/samos123/k8s-qos/pkg/poller"
	"time"
)

func main() {
	interval := flag.Int("interval", 10, "Interval to pull metrics from kubelet")
	url := "http://localhost:10255/stats/summary"
	kubeletUrl := flag.String("url", url, "URL to kubelet stats/summary endpoint")
	flag.Parse()
	c := time.Tick(time.Second * time.Duration(*interval))
	for _ = range c {
		fmt.Println("Getting metrics from URL:", *kubeletUrl)
		json := poller.GetURL(*kubeletUrl)
		metrics := poller.ParseNetworkMetrics(json)
		fmt.Println(metrics)
	}
	fmt.Println("Exitting")
}
