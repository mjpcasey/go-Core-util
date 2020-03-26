package gcoordinatorTypes

const (
	// EventChanged node data property changed event
	EventChanged = 1 << iota
	// EventCreated node point created event
	EventCreated
	// EventDeleted node deleted event
	EventDeleted
	// EventChildrenChanged node children changed event
	EventChildrenChanged
)

// Event Coordinator node watch callback event object
type Event struct {
	Type int
	Data interface{}
}

// EventCallback node module watch callback function
type EventCallback func(evt Event) error