package event

// Serializer converts between Events and data
type Serializer interface {
	// MarshalEvent converts an Event to a byte array
	MarshalEvent(event Event) ([]byte, error)

	// UnmarshalEvent converts an Event backed into a blob of data
	UnmarshalEvent(data []byte) (Event, error)
}
