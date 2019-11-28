package poller

import (
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
)

type PodNetworkStats struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	RxBytes   int64  `json:"rxBytes"`
	TxBytes   int64  `json:"txBytes"`
}

func GetMetrics(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error occured trying to get metrics from Kubelet:", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error occured reading HTTP body:", err)
	}
	return body
}

func ParseNetworkMetrics(json []byte) []PodNetworkStats {
	pods := gjson.GetBytes(json, `pods`)
	podCount := gjson.GetBytes(json, `pods.#`).Int()
	networkMetrics := make([]PodNetworkStats, podCount)
	i := 0
	pods.ForEach(func(key, value gjson.Result) bool {
		s := PodNetworkStats{
			Name:      value.Get("podRef.name").String(),
			Namespace: value.Get("podRef.namespace").String(),
			RxBytes:   value.Get("network.rxBytes").Int(),
			TxBytes:   value.Get("network.txBytes").Int()}
		networkMetrics[i] = s
		i++
		return true // keep iterating
	})
	return networkMetrics
}
