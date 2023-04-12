package output

import (
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EventEmitter", func() {
	var (
		emitter *EventEmitter
	)

	BeforeEach(func() {
		emitter = NewEventEmitter()
	})

	Describe("On", func() {
		It("should register a listener for an event", func() {
			ch := emitter.On("myEvent")
			Expect(ch).NotTo(BeNil())
		})
	})

	Describe("Emit", func() {
		It("should send an event to all registered listeners", func() {
			done := make(chan interface{})
			go func() {
				// test code to run asynchronously
				var wg sync.WaitGroup
				wg.Add(2)
	
				ch1 := emitter.On("myEvent")
				ch2 := emitter.On("myEvent")
	
				go func() {
					defer wg.Done()
					data := <-ch1
					Expect(data).To(Equal("hello"))
				}()
	
				go func() {
					defer wg.Done()
					data := <-ch2
					Expect(data).To(Equal("hello"))
				}()
	
				emitter.Emit("myEvent", "hello")
	
				wg.Wait()
				close(done)
			}()

			Eventually(done, 2).Should(BeClosed())
		})
	})

	Describe("Off", func() {
		It("should remove the listener for an event", func() {
			ch1 := emitter.On("myEvent")
			ch2 := emitter.On("myEvent") // an extra listener that should remain

			emitter.Off("myEvent", ch1)

			Expect(len(emitter.listeners["myEvent"])).To(Equal(1))
			Expect(emitter.listeners["myEvent"][0]).To(Equal(ch2))
		})

		It("should close the channel of the removed listener", func() {
			go func() {
				ch := emitter.On("myEvent")
				emitter.Off("myEvent", ch)
				Expect(ch).To(BeClosed())
			}()
		})
	})
})
