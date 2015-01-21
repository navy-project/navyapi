package main

import (
	"log"
	"os"

	"github.com/coreos/go-etcd/etcd"
	navy "github.com/navy-project/navyapi/server"
)

func main() {
	etcdClient := setupEtcd()
	server := navy.NewServer(etcdClient)
	server.ListenAndServe()
}

func setupEtcd() *etcd.Client {
	etcdserver := "http://" + os.Getenv("ETCD_PORT_4001_TCP_ADDR") + ":" + os.Getenv("ETCD_PORT_4001_TCP_PORT")
	log.Println("Connected To Etcd: ", etcdserver)
	return etcd.NewClient([]string{etcdserver})
}
