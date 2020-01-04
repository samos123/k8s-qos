package poller

import (
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"strings"
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
	Name       string
	UID        string
	IPAddress  net.IP
	Veth       string
	BWLimit    int64
	Containers []Container
}

type Container struct {
	Name string
	ID   string
	Veth string
}

// Applies bandwidth limits based on the number of pods. The parameter totalBW
// is in Mbit/s
func PodCountLimiter(pods []Pod, totalBW int) {
	n := len(pods)
	var limit int64
	if n == 1 {
		return
	} else if n >= 2 && n < 4 {
		limit = int64(float64(totalBW) * float64(0.8))
	} else if n >= 4 && n < 8 {
		limit = int64(float64(totalBW) * float64(0.6))
	} else {
		limit = int64(float64(totalBW) * float64(0.4))
	}
	for _, pod := range pods {
		pod.Limit(limit, 20)
	}
}

// Calculate total BW based on the amount of CPUs. This currently assumes
// that n1 machine type is used.
// TODO: Have a more dynamic way of calculating total bandwidth
func TotalBWonGKE(cpus int) int {
	switch {
	case cpus == 1:
		return 2 * 1000
	case cpus >= 2 && cpus <= 4:
		return 10 * 1000
	case cpus >= 8 && cpus <= 15:
		return 16 * 1000
	case cpus >= 16:
		return 32 * 1000
	}
	return 32 * 1000
}

func (p *Pod) GetVeth() {
	// TODO use Docker golang client
	for i := range p.Containers {
		c := p.Containers[i]
		out, err := exec.Command("getveth.sh", c.ID).Output()
		if err != nil {
			log.WithFields(log.Fields{"err": err, "out": out}).Warn("error running getveth.sh")
			continue
		}
		log.WithFields(log.Fields{"out": out, "container": c, "pod": p}).Info("ran getveth.sh")
		c.Veth = string(out)
		if strings.HasPrefix(c.Veth, "veth") {
			split := strings.Split(c.Veth, "@")
			c.Veth = split[0]
			p.Veth = c.Veth
			log.WithFields(log.Fields{"container": c, "pod": p}).Info("found veth")
			break
		}
	}
}

// Apply a bandwidth limit on the pod
func (p *Pod) Limit(rate int64, latency int) {
	log.WithFields(log.Fields{"pod": p, "limit": rate}).Info("Applying limit to pod")
	p.GetVeth()
	p.BWLimit = rate
	TcLimit(p.Veth, string(rate)+"mbit", string(latency)+"ms")
}

// Apply a bandwdith limit using the tc linux command
func TcLimit(netinterface, rate, latency string) {
	cmd := exec.Command("tc", "qdisc", "change", "dev", netinterface,
		"root", "tbf", "rate", rate, "latency", latency, "burst", "1540")
	if err := cmd.Run(); err != nil {
		log.Error("Error occured executing tc command", err)
	}
}

func GetURL(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Warn("Error occured trying to get URL from Kubelet:", url, err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warn("Error occured reading HTTP body:", err)
		return nil
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
			c.ID = strings.TrimPrefix(cResult.Get("containerID").String(), "docker://")
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
