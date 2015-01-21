package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Convoys struct {
	Etcd   EtcdClient
	routes *mux.Router
}

type ConvoyRequest struct {
	Name     string `json:"name"`
	Manifest string `json:"manifest"`
}

type ConvoyQueueEvent struct {
	Request  string `json:"request"`
	Name     string `json:"name"`
	Manifest string `json:"manifest"`
}

func HandleConvoys(r *mux.Router, e EtcdClient) {
	handler := &Convoys{Etcd: e}
	r.HandleFunc("/convoys", handler.createConvoy).Methods("POST")
	r.HandleFunc("/convoys/{name}", handler.deleteConvoy).Methods("DELETE")
}

func (c Convoys) createConvoy(w http.ResponseWriter, r *http.Request) {
	apirequest := &ConvoyRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(apirequest)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	println("Create Convoy: " + apirequest.Name)
	item := ConvoyQueueEvent{"create", apirequest.Name, apirequest.Manifest}
	err = c.queueEvent(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}

func (c Convoys) deleteConvoy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	item := ConvoyQueueEvent{Request: "destroy", Name: name}
	err := c.queueEvent(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}

}

func (c Convoys) queueEvent(item ConvoyQueueEvent) error {
	b, err := json.Marshal(item)
	if err != nil {
		return err
	}
	c.Etcd.AddChild("/navy/queues/convoys", string(b), 0)
	return nil
}
