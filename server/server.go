package navyapi

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/navy-project/navyapi/api"
)

type Server struct {
	Port   string
	Routes *mux.Router
}

func NewServer(e api.EtcdClient) *Server {
	server := &Server{Port: "4040"}
	r := mux.NewRouter()
	api.HandleConvoys(r, e)
	api.HandleEvents(r, e)
	server.Routes = r
	return server
}

func (s *Server) ListenAndServe() {
	err := http.ListenAndServe(":"+s.Port, s.Routes)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
