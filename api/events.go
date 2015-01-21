package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Events struct {
	Etcd   EtcdClient
	routes *mux.Router
}

func HandleEvents(r *mux.Router, e EtcdClient) {
	handler := &Events{Etcd: e}
	r.HandleFunc("/events", handler.streamEvents).Methods("GET")
}

func (e Events) streamEvents(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request /events")
	containerChan := make(chan *ContainerEvent)
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "keep-alive")
	encoder := json.NewEncoder(w)
	WatchContainerEvents(e.Etcd, containerChan)
	for event := range containerChan {
		encoder.Encode(event)
		flusher.Flush()
	}
}
