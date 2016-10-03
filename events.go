package trace

import xtr "golang.org/x/net/trace"

// An EventLog provides a log of events associated with a specific object.
type EventLog struct {
	family string
	title  string
	el     xtr.EventLog
}

// NewEventLog returns an initialized EventLog with the given family and title.
func NewEventLog(family, title string) xtr.EventLog {
	e := &EventLog{
		family: family,
		title:  title,
		el:     xtr.NewEventLog(family, title),
	}
	return e
}

// Printf formats its arguments with fmt.Sprintf and adds the
// result to the event log.
func (e *EventLog) Printf(format string, a ...interface{}) {
	newfmt, newvals := addEvent(e, format, a...)
	Log.Printf(newfmt, newvals...)
	e.el.Printf(format, a...)
}

// Errorf is like Printf, but it marks this event as an error.
func (e *EventLog) Errorf(format string, a ...interface{}) {
	Log.Printf("[ERROR] "+format, a...)
	e.el.Errorf(format, a...)
}

// Finish declares that this event log is complete.
// The event log should not be used after calling this method.
func (e *EventLog) Finish() {
	e.el.Finish()
}
