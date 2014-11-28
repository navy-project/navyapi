package api

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"strings"
)

type ContainerEvent struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type containerState struct {
	State string
}

func WatchContainerEvents(e EtcdClient, c chan *ContainerEvent) {
	stopChan := make(chan bool)
	responseChan := make(chan *etcd.Response)
	go e.Watch("/navy/containers", 0, true, responseChan, stopChan)
	go func() {
		for response := range responseChan {
			if strings.HasSuffix(response.Node.Key, "/actual") {
				if response.Action == "set" {
					state := &containerState{}
					err := json.Unmarshal([]byte(response.Node.Value), state)
					if err != nil {
						fmt.Println(err)
					}
					name := extractNameFromContainerKey(response.Node.Key)
					c <- &ContainerEvent{Name: name, Status: state.State}
				} else if response.Action == "delete" {
					name := extractNameFromContainerKey(response.Node.Key)
					c <- &ContainerEvent{Name: name, Status: "destroy"}
				}
			}
		}
	}()
}

func extractNameFromContainerKey(key string) string {
	return strings.TrimRight(strings.TrimPrefix(key, "/navy/containers/"), "/actual")
}
