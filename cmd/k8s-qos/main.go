package main

import (
	"flag"
	"fmt"
	"github.com/samos123/k8s-qos/pkg/poller"
	"time"
)

func main() {
	interval := flag.Int("interval", 10, "Interval to pull data from kubelet")
	url := "http://localhost:10255"
	kubeletUrl := flag.String("url", url, "URL to kubelet api endpoint")
	flag.Parse()
	controller := poller.NewContoller()
	c := time.Tick(time.Second * time.Duration(*interval))
	for _ = range c {
		json := poller.GetURL(*kubeletUrl + "/stats/summary")
		metrics := poller.ParseNetworkMetrics(json)
		fmt.Println(metrics)
		pods := controller.GetPods(*kubeletUrl + "/pods")
		for _, p := range pods {
			p.GetVeth()
		}
		fmt.Println(pods)
	}
	fmt.Println("Exitting")
}
