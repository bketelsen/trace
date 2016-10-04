## Trace

[![](https://godoc.org/github.com/bketelsen/trace?status.svg)](http://godoc.org/github.com/bketelsen/trace)  [![Go Report Card](https://goreportcard.com/badge/github.com/bketelsen/trace)](https://goreportcard.com/report/github.com/bketelsen/trace)

`trace` is a wrapper for [the net/trace package](https://github.com/golang/net/tree/master/trace) that adds logging and metrics to de-clutter your functions.

`trace` wraps all of the functionality of net/trace, but also replicates the logs to a structured logger built on Go's standard logging library. 

Metrics are exported to Prometheus with trace duration histograms by name/family, trace counts by name/family, and errors by trace name/family.

the net/trace#EventLog is also implemented in the same manner, minus metrics exposition, which doesn't make sense there.

examples/gogrep has an example command-line application that shows usage of the `trace` functionality to both capture trace information and logs with a single tool.

examples/service has an example web application that shows usage of the `trace` functionality to both capture trace information and logs with a single tool, combined with the `trace.EventLog` which serves as a single logging and event source for your application.


### Log Output - trace

	2016/10/03 00:27:34 message=found file=../../events.go trace=main
	2016/10/03 00:27:34 message=found file=../../examples/gogrep/main.go trace=main
	2016/10/03 00:27:34 message=found file=../../log.go trace=main
	2016/10/03 00:27:34 message=found file=../../trace.go trace=main
	2016/10/03 00:27:34 message=found file=../../metrics.go trace=main
	2016/10/03 00:27:34 trace=main : hit count 5
	2016/10/03 00:27:34 message=finished hits=5 trace=main

### Log Output - EventLog

	2016/10/03 00:34:05 name=http - Listening on :3000

### Metrics Output

`trace` offers two useful and one fun way to expose your metrics.

`trace.ServeMetrics()` will serve the metrics in Prometheus text format.  Use this for long-running apps/services.

`trace.PushMetrics()` will push the metrics to a Prometheus push server.  Use this for command-line utilities.

`trace.DumpMetrics()` will return a string with the metrics that Prometheus would serve, suitable for inspection, printing, tests.

### /debug endpoints

`trace` exposes the underlying net/trace /debug/requests and /debug/events endpoints for handy visual representation of the traces, their timing/histograms, and the event log of your application.

![Requests](/examples/images/requests.png?raw=true "Requests")
![Events](/examples/images/events.png?raw=true "Events")
