package main

import (
	"flag"
	"github.com/samos123/k8s-qos/pkg/poller"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	log.SetLevel(log.DebugLevel)
	interval := flag.Int("interval", 10, "Interval to pull data from kubelet")
	url := "http://localhost:10255"
	kubeletUrl := flag.String("url", url, "URL to kubelet api endpoint")
	flag.Parse()
	c := time.Tick(time.Second * time.Duration(*interval))
	for _ = range c {
		json := poller.GetURL(*kubeletUrl + "/stats/summary")
		metrics := poller.ParseNetworkMetrics(json)
		log.WithFields(log.Fields{"metrics": metrics}).Info("Parsed metrics")
		json = poller.GetURL(*kubeletUrl + "/pods")
		pods := poller.ParsePods(json)
		for _, p := range pods {
			p.GetVeth()
		}
		log.WithFields(log.Fields{"pods": pods}).Info("Parsed pods")
	}
}
