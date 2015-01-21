package api

import (
	"github.com/coreos/go-etcd/etcd"
)

type EtcdClient interface {
	//Set(key string, value string, ttl uint64) (*etcd.Response, error)
	//Get(key string, sort, recur bool) (*etcd.Response, error)
	AddChild(key string, value string, ttl uint64) (*etcd.Response, error)
	Watch(prefix string, waitIndex uint64, recursive bool,
		receiver chan *etcd.Response, stop chan bool) (*etcd.Response, error)
}
