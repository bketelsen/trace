package trace

import (
	"bytes"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"strconv"
	"sync"
)

// Log is the logger used by all trace package functions
var Log = stdlog.New(os.Stderr, "", stdlog.LstdFlags)

// SetLogger replaces the default logger with a new
// one that writes to 'out', has 'prefix', and flags 'flag'
func SetLogger(out io.Writer, prefix string, flag int) {
	Log = stdlog.New(out, prefix, flag)
}

type logmessage struct {
	m  string
	kv []keyval
}

type keyval struct {
	key string
	val interface{}
}

// KeyValue creates a Key/Value pair that is suitable for use in the
// LogMessage() function
func KeyValue(key string, val interface{}) keyval {
	return keyval{key: key, val: val}
}

// LogMessage creates a message that complies with fmt.Stringer but also
// includes a message and key/value pairs for structured logging.
// Use it in trace.LazyLog:
//
// 		t.LazyLog(trace.LogMessage("found", trace.KeyValue("file", name)), false)
func LogMessage(message string, keyvals ...keyval) *logmessage {
	return &logmessage{m: message, kv: keyvals}
}

func logMessageWithTrace(t *trace, message string, keyvals ...keyval) string {
	keyvals = append(keyvals, KeyValue("trace", t.title))
	b := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(b)
	b.Reset()

	writeKV := func(kv []keyval) {
		for _, kv := range kv {
			b.WriteString(" ")
			b.WriteString(kv.key)
			b.WriteString("=")
			b.WriteString(toString(kv.val))
		}
	}

	b.WriteString(message)
	writeKV(keyvals)

	return string(b.Bytes())
}

func addTrace(t *trace, format string, a ...interface{}) (string, []interface{}) {
	newf := "trace=%s : " + format
	var b []interface{}
	b = append(b, t.title)
	b = append(b, a...)
	return newf, b
}

func addEvent(t *EventLog, format string, a ...interface{}) (string, []interface{}) {
	newf := "name=%s - " + format
	var b []interface{}
	b = append(b, t.title)
	b = append(b, a...)
	return newf, b
}
func (m *logmessage) String() string {

	b := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(b)
	b.Reset()

	writeKV := func(kv []keyval) {
		for _, kv := range kv {
			b.WriteString(" ")
			b.WriteString(kv.key)
			b.WriteString("=")
			b.WriteString(toString(kv.val))
		}
	}

	b.WriteString("message=")
	b.WriteString(m.m)
	writeKV(m.kv)

	return string(b.Bytes())
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func toString(x interface{}) string {
	switch v := x.(type) {
	case bool:
		if v {
			return "true"
		}
		return "false"
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', 6, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', 6, 64)
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("!(%T=%+v)", x, x)
	}
}
