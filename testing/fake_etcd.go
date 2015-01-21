package testing

import (
	"errors"
	"github.com/coreos/go-etcd/etcd"
	"strings"
)

type FakeEtcd struct {
	keys     map[string]*etcd.Response
	watchers map[string]chan *etcd.Response
}

func (c *FakeEtcd) Get(key string, sort, recur bool) (*etcd.Response, error) {
	stored, ok := c.keys[key]
	if !ok {
		return nil, errors.New("Missing Key")
	}
	return stored, nil
}

func (c *FakeEtcd) Set(key string, value string, ttl uint64) (*etcd.Response, error) {
	response := &etcd.Response{Node: &etcd.Node{Key: key, Value: value}}
	c.keys[key] = response
	c.alertWatchers(key, "set", response)
	return response, nil
}

func (c *FakeEtcd) Delete(key string, recursive bool) (*etcd.Response, error) {
	stored, err := c.Get(key, false, recursive)
	if err != nil {
		return nil, err
	}
	c.alertWatchers(key, "delete", stored)
	c.keys[key] = nil
	return stored, nil
}

func (c *FakeEtcd) AddChild(key string, value string, ttl uint64) (*etcd.Response, error) {
	var parent *etcd.Response
	parent, ok := c.keys[key]
	if !ok {
		parent, _ = c.Set(key, "A Dir", 0)
	}
	response := &etcd.Response{Node: &etcd.Node{Value: value}}
	parent.Node.Nodes = append(parent.Node.Nodes, response.Node)
	return response, nil
}

func (c *FakeEtcd) Watch(prefix string, waitIndex uint64, recursive bool,
	receiver chan *etcd.Response, stop chan bool) (*etcd.Response, error) {
	if receiver == nil {
		return nil, nil
	}
	pipe := make(chan *etcd.Response)
	c.watchers[prefix] = pipe
	for r := range pipe {
		receiver <- r
	}
	return nil, nil
}

/* Method purely for test purposes */

func NewFakeEtcd() *FakeEtcd {
	etcdClient := &FakeEtcd{}
	etcdClient.keys = make(map[string]*etcd.Response)
	etcdClient.watchers = make(map[string]chan *etcd.Response)
	return etcdClient
}

func (c *FakeEtcd) alertWatchers(key, action string, r *etcd.Response) {
	r.Action = action
	for path, watchChan := range c.watchers {
		if strings.HasPrefix(key, path) {
			watchChan <- r
		}
	}
}
