package trace

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net/http"
	"time"

	xtr "golang.org/x/net/trace"
)

type trace struct {
	family string
	title  string
	start  time.Time
	trace  xtr.Trace
	id     string
	err    bool
}

// SetAuthRequest sets the AuthRequest function for the underlying trace HTTP listener, which
// determines whether a specific request is permitted to load the
// /debug/requests or /debug/events pages.
func SetAuthRequest(f func(req *http.Request) (any, sensitive bool)) {
	xtr.AuthRequest = f
}

// New returns a new Trace with the specified family and title.
func New(family, title string) xtr.Trace {
	return newTrace(family, title)
}

func newTrace(family, title string) *trace {
	tr := xtr.New(family, title)
	s := &trace{
		family: family,
		title:  title,
		trace:  tr,
		start:  time.Now(),
		id:     randomID(),
	}
	return s
}

func (t *trace) child(title string) *trace {
	tr := xtr.New(t.family, title)
	s := &trace{
		family: t.family,
		title:  title,
		trace:  tr,
		start:  time.Now(),
		id:     randomID(),
	}
	return s
}

// LazyLog adds x to the event log. It will be evaluated each time the
// /debug/requests page is rendered. Any memory referenced by x will be
// pinned until the trace is finished and later discarded.
func (t *trace) LazyLog(x fmt.Stringer, sensitive bool) {
	Log.Println(logMessageWithTrace(t, x.String()))
	t.trace.LazyLog(x, sensitive)
}

// LazyPrintf evaluates its arguments with fmt.Sprintf each time the
// /debug/requests page is rendered. Any memory referenced by a will be
// pinned until the trace is finished and later discarded.
func (t *trace) LazyPrintf(format string, a ...interface{}) {
	newfmt, newvals := addTrace(t, format, a...)
	Log.Printf(newfmt, newvals...)
	t.trace.LazyPrintf(format, a...)
}

// SetError declares that this trace resulted in an error.
func (t *trace) SetError() {
	t.err = true
	t.trace.SetError()
}

// SetRecycler sets a recycler for the trace.
// f will be called for each event passed to LazyLog at a time when
// it is no longer required, whether while the trace is still active
// and the event is discarded, or when a completed trace is discarded.
func (t *trace) SetRecycler(f func(interface{})) {
	t.trace.SetRecycler(f)
}

// SetTraceInfo sets the trace info for the trace.
// This is currently unused.
func (t *trace) SetTraceInfo(traceID uint64, spanID uint64) {
	t.trace.SetTraceInfo(traceID, spanID)
}

// SetMaxEvents sets the maximum number of events that will be stored
// in the trace. This has no effect if any events have already been
// added to the trace.
func (t *trace) SetMaxEvents(m int) {
	t.trace.SetMaxEvents(m)
}

// Finish declares that this trace is complete.
// The trace should not be used after calling this method.
func (t *trace) Finish() {
	if t.err {
		incrError(t)
	}
	incr(t)
	duration(t)
	if t.err {
		incrError(t)
	}
	t.trace.Finish()
}
func randomID() string {
	var v [8]byte
	b := v[:8]
	rand.Read(b)
	u := binary.BigEndian.Uint64(b)
	return fmt.Sprintf("%x", u)
}
