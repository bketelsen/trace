// Package trace is a wrapper for [the net/trace package](https://github.com/golang/net/tree/master/trace) that adds logging and metrics to de-clutter your functions.
//
//`trace` wraps all of the functionality of net/trace, but also replicates the logs to a structured logger built on Go's standard logging library.
//
//Metrics are exported to Prometheus of trace durations by name/family, trace counts by name/family, and errors by trace name/family.
//
//the net/trace#EventLog is also implemented in the same manner, minus metrics exposition, which doesn't make sense there.
//
//examples/gogrep has an example command-line application that shows usage of the `trace` functionality to both capture trace information and logs with a single tool.
//
//examples/service has an example web application that shows usage of the `trace` functionality to both capture trace information and logs with a single tool, combined with the `trace.EventLog` which serves as a single logging and event source for your application.
package trace
