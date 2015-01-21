package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http/httptest"

	"encoding/json"

	"github.com/navy-project/navyapi/api"
	nc "github.com/navy-project/navyapi/client"
	sv "github.com/navy-project/navyapi/server"
	ct "github.com/navy-project/navyapi/testing"
)

var _ = Describe("/convoys", func() {
	var etcdClient *ct.FakeEtcd
	var server *sv.Server
	var client *nc.Client

	BeforeEach(func() {
		etcdClient = ct.NewFakeEtcd()
		server = sv.NewServer(etcdClient)

		testServer := httptest.NewServer(server.Routes)
		client = nc.NewClient(testServer.URL)
	})

	Describe("Create Convoy", func() {
		It("Puts the convoy on the create queue", func() {
			err := client.CreateConvoy("the_name", "some_manifest_yaml")
			Expect(err).ShouldNot(HaveOccurred())

			dir, err := etcdClient.Get("/navy/queues/convoys", false, false)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(dir.Node.Nodes)).To(Equal(1))
			event := dir.Node.Nodes[0].Value

			item := &api.ConvoyQueueEvent{}
			err = json.Unmarshal([]byte(event), item)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(item.Request).To(Equal("create"))
			Expect(item.Name).To(Equal("the_name"))
			Expect(item.Manifest).To(Equal("some_manifest_yaml"))
		})
	})

	Describe("Delete Convoy", func() {
		It("Puts the convoy on the destroy queue", func() {
			err := client.DeleteConvoy("the_name")
			Expect(err).ShouldNot(HaveOccurred())

			dir, err := etcdClient.Get("/navy/queues/convoys", false, false)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(dir.Node.Nodes)).To(Equal(1))
			event := dir.Node.Nodes[0].Value

			item := &api.ConvoyQueueEvent{}
			err = json.Unmarshal([]byte(event), item)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(item.Request).To(Equal("destroy"))
			Expect(item.Name).To(Equal("the_name"))
		})
	})
})
