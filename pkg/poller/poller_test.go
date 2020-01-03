package poller

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

var containerID string
var CLI client.Client

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

func setup() {
	ctx := context.Background()
	CLI, err := client.NewClientWithOpts(client.FromEnv)
	check(err)
	CLI.NegotiateAPIVersion(ctx)

	reader, err := CLI.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	check(err)

	io.Copy(os.Stdout, reader)

	resp, err := CLI.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"sleep", "60"},
	}, nil, nil, "test-container")
	check(err)

	containerID = resp.ID

	if err := CLI.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

}

func cleanup() {
	ctx := context.Background()
	CLI, err := client.NewClientWithOpts(client.FromEnv)
	check(err)
	CLI.NegotiateAPIVersion(ctx)

	var second time.Duration = time.Second
	if err := CLI.ContainerStop(ctx, containerID, &second); err != nil {
		panic(err)
	}
	if err := CLI.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true}); err != nil {
		panic(err)
	}

}

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	cleanup()
	os.Exit(ret)
}

func TestParseNetworkMetrics(t *testing.T) {
	json, err := ioutil.ReadFile("../../test/stats-summary-1.json")
	check(err)
	s := ParseNetworkMetrics(json)
	equals(t, s[0].Name, "network-metering-agent-bljqp")
	fmt.Println(s)
}
func TestParsePods(t *testing.T) {
	json, err := ioutil.ReadFile("../../test/pods.json")
	check(err)
	s := ParsePods(json)
	equals(t, s[0].Name, "webhook-6cbdc8b54-d5fq7")
	fmt.Println(s[0])
	equals(t, s[0].Containers[0].Name, "webhook")
}

func TestGetVeth(t *testing.T) {
	containers := []Container{Container{ID: containerID}}
	path, err := filepath.Abs("../../tools")
	check(err)
	fmt.Println("Path:", path)

	os.Setenv("PATH", os.Getenv("PATH")+":"+path)
	fmt.Println("Path:", os.Getenv("PATH"))
	p := Pod{Name: "test", Containers: containers}
	p.GetVeth()
	fmt.Println(p)

}
