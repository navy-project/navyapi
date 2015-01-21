package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"

	//"encoding/json"

	"github.com/navy-project/navyapi/api"
	ct "github.com/navy-project/navyapi/testing"
)

var _ = Describe("Container Events", func() {
	var etcdClient *ct.FakeEtcd
	var eventChannel chan *api.ContainerEvent

	BeforeEach(func() {
		etcdClient = ct.NewFakeEtcd()
		eventChannel = make(chan *api.ContainerEvent)
		api.WatchContainerEvents(etcdClient, eventChannel)
		time.Sleep(2 * time.Microsecond)
	})

	It("Sends ACTUAL state changes down the channel", func() {
		etcdClient.Set("/navy/containers/some_container_name/actual", "{\"state\":\"a_status\"}", 0)
		select {
		case event := <-eventChannel:
			Expect(event.Status).To(Equal("a_status"))
			Expect(event.Name).To(Equal("some_container_name"))
		case <-time.After(1 * time.Second):
			Fail("Did Not Receive on eventChannel")
		}
	})

	It("Ingores none ACTUAL state changes", func() {

		etcdClient.Set("/navy/containers/some_container_name/desired", "{\"state\":\"a_status\"}", 0)
		select {
		case event := <-eventChannel:
			Expect(event).ShouldNot(HaveOccurred())
		case <-time.After(1 * time.Second):
			//OK
		}
	})

	It("Sends DELETE of actual down the channel", func() {
		etcdClient.Set("/navy/containers/some_container_name/actual", "{\"state\":\"a_status\"}", 0)
		<-eventChannel //Clear Channel
		etcdClient.Delete("/navy/containers/some_container_name/actual", false)
		select {
		case event := <-eventChannel:
			Expect(event.Name).To(Equal("some_container_name"))
			Expect(event.Status).To(Equal("destroy"))
		case <-time.After(1 * time.Second):
			Fail("Did Not Receive on eventChannel")
		}
	})
})
