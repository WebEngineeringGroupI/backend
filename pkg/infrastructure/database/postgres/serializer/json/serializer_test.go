package json_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres/serializer/json"
)

var _ = Describe("Domain / Event / Serializer / JSON", func() {
	var (
		serializer *json.Serializer
	)
	BeforeEach(func() {
		serializer = json.NewSerializer(&KnownEvent{})
	})

	It("serializes a known event into JSON", func() {
		data, err := serializer.MarshalEvent(&KnownEvent{
			Base: event.Base{
				ID:      "someId",
				Version: 0,
				At:      time.Time{},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(data).To(MatchJSON(`{"type":"KnownEvent","data":{"ID":"someId","Version":0,"At":"0001-01-01T00:00:00Z"}}`))
	})

	It("deserializes a known event into JSON", func() {
		eventPayload := []byte(`{"type":"KnownEvent","data":{"ID":"someId","Version":0,"At":"0001-01-01T00:00:00Z"}}`)

		unmarshalledEvent, err := serializer.UnmarshalEvent(eventPayload)
		Expect(err).ToNot(HaveOccurred())
		Expect(unmarshalledEvent).To(Equal(&KnownEvent{
			Base: event.Base{
				ID:      "someId",
				Version: 0,
				At:      time.Time{},
			},
		}))
	})

	When("the event is unknown for the marshaller", func() {
		It("marshals it correctly", func() {
			data, err := serializer.MarshalEvent(&UnknownEvent{
				Base: event.Base{
					ID:      "someId",
					Version: 0,
					At:      time.Time{},
				},
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(data).To(MatchJSON(`{"type":"UnknownEvent","data":{"ID":"someId","Version":0,"At":"0001-01-01T00:00:00Z"}}`))
		})

		It("is unable to unmarshal it correctly, because the type cannot be created", func() {
			eventPayload := []byte(`{"type":"UnknownEvent","data":{"ID":"someId","Version":0,"At":"0001-01-01T00:00:00Z"}}`)

			unmarshalledEvent, err := serializer.UnmarshalEvent(eventPayload)

			Expect(err).ToNot(MatchError("unknown event type"))
			Expect(unmarshalledEvent).To(BeNil())
		})

		When("the type is bound after the creation", func() {
			It("is able to unmarshal it correctly again, because the type is now known and can be created", func() {
				eventPayload := []byte(`{"type":"UnknownEvent","data":{"ID":"someId","Version":0,"At":"0001-01-01T00:00:00Z"}}`)
				serializer.Bind(&UnknownEvent{})

				unmarshalledEvent, err := serializer.UnmarshalEvent(eventPayload)

				Expect(err).ToNot(HaveOccurred())
				Expect(unmarshalledEvent).To(Equal(&UnknownEvent{
					Base: event.Base{
						ID:      "someId",
						Version: 0,
						At:      time.Time{},
					},
				}))
			})
		})
	})
})

type KnownEvent struct {
	event.Base
}

type UnknownEvent struct {
	event.Base
}
