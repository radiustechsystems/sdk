package common

// Event represents an EVM contract event emitted during transaction execution.
// Contains decoded event data and the raw event payload.
type Event struct {
	// Name is the name of the event
	Name string

	// Data is the data of the event
	Data map[string]interface{}

	// Raw is the raw data of the event
	Raw []byte
}

// NewEvent creates a new Event with the given name, data, and raw bytes
// @param name The name of the event
// @param data The decoded data of the event as key-value pairs
// @param raw The raw bytes of the event
// @return A new Event instance
func NewEvent(name string, data map[string]interface{}, raw []byte) *Event {
	return &Event{
		Name: name,
		Data: data,
		Raw:  raw,
	}
}
