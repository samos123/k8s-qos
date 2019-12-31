package poller

import (
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type PodNetworkStats struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	RxBytes   int64     `json:"rxBytes"`
	TxBytes   int64     `json:"txBytes"`
	Time      time.Time `json:"time"`
}

type Pod struct {
	Name         string
	UID          string
	IPAddress    net.IP
	Veth         string
	IngressLimit int64
	EgressLimit  int64
	Containers   []Container
}

type Container struct {
	Name string
	ID   string
	Veth string
}

func GetURL(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error occured trying to get URL from Kubelet:", url, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error occured reading HTTP body:", err)
	}
	return body
}

func ParsePods(json []byte) []Pod {
	pods := gjson.GetBytes(json, `items`)
	podCount := gjson.GetBytes(json, `items.#`).Int()

	podsArr := make([]Pod, podCount)
	i := 0

	// items[status.containerStatuses[name, imageID, containerID]]
	pods.ForEach(func(key, podResult gjson.Result) bool {
		p := Pod{}
		p.Name = podResult.Get("metadata.name").String()
		p.UID = podResult.Get("metadata.uid").String()
		p.IPAddress = net.ParseIP(podResult.Get("status.podIP").String())

		containers := podResult.Get("status.containerStatuses")
		containerCount := podResult.Get("status.containerStatuses.#").Int()
		p.Containers = make([]Container, containerCount)
		j := 0
		containers.ForEach(func(key2, cResult gjson.Result) bool {
			c := Container{}
			c.Name = cResult.Get("name").String()
			c.ID = cResult.Get("containerID").String()
			p.Containers[j] = c
			j++
			return true
		})

		podsArr[i] = p
		i++
		return true
	})
	return podsArr
}

func ParseNetworkMetrics(json []byte) []PodNetworkStats {
	pods := gjson.GetBytes(json, `pods`)
	podCount := gjson.GetBytes(json, `pods.#`).Int()
	networkMetrics := make([]PodNetworkStats, podCount)
	i := 0
	pods.ForEach(func(key, value gjson.Result) bool {
		networkTime, _ := time.Parse(time.RFC3339, value.Get("network.time").String())
		s := PodNetworkStats{
			Name:      value.Get("podRef.name").String(),
			Namespace: value.Get("podRef.namespace").String(),
			RxBytes:   value.Get("network.rxBytes").Int(),
			TxBytes:   value.Get("network.txBytes").Int(),
			Time:      networkTime}
		networkMetrics[i] = s
		i++
		return true // keep iterating
	})
	return networkMetrics
}

// NetworkThoughput calculates the ingress and egress thoughput using 2 datapoints using the formula
// (latest - old) / (delta in time between latest and old)
func NetworkThroughput(old, latest PodNetworkStats) (rx, tx float64) {
	timeDelta := (latest.Time.Sub(old.Time).Seconds())
	rx = float64(latest.RxBytes-old.RxBytes) / timeDelta
	tx = float64(latest.TxBytes-old.TxBytes) / timeDelta
	return rx, tx
}
