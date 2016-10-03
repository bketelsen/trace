package trace

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/prometheus/common/expfmt"
)

func init() {
	prometheus.MustRegister(histograms)
	prometheus.MustRegister(counts)
	prometheus.MustRegister(errcounts)
}

// Namespace is used to differentiate metrics - and specifically used
// in prometheus reporting. It may safely be left blank.
var Namespace string

// Subsystem is used to differentiate metrics - and specifically used
// in prometheus reporting.  It may safely be left blank.
var Subsystem string

var histograms = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: Namespace,
	Subsystem: Subsystem,
	Name:      "latency_seconds",
	Help:      "The latency of the labeled function in seconds, partitioned by family and title",
}, []string{"family", "title"})

var counts = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: Subsystem,
		Name:      "traces",
		Help:      "How many times a trace is created, partitioned by family and title.",
	},
	[]string{"family", "title"})

var errcounts = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: Subsystem,
		Name:      "errors",
		Help:      "How many times an error condition was encountered, partitioned by family and title.",
	},
	[]string{"family", "title"})

func incr(t *trace) {
	//cm.Counter(t.name).Add()
	counts.With(prometheus.Labels{"family": t.family, "title": t.title}).Add(1)
}

func incrError(t *trace) {
	//cm.Counter(s.errorname).Add()
	errcounts.With(prometheus.Labels{"family": t.family, "title": t.title}).Add(1)
}

func duration(t *trace) {
	sec := time.Since(t.start).Seconds()
	histograms.With(prometheus.Labels{"family": t.family, "title": t.title}).Observe(sec)

}

// ServeMetrics serves Prometheus metrics endpoint on the
// provided net.Listener
// Use for long-running services.  For cli tools, use PushMetrics instead.
func ServeMetrics(ctx context.Context, l net.Listener) error {
	return http.Serve(l, promhttp.Handler())
}

// PushMetrics sends the metrics collected to the prometheus push gateway
// at the url in `gatewayURL` with the job name of `task`
// Should be called in a defer in main() of a cli application.
// For long-running services, use ServeMetrics instead
func PushMetrics(ctx context.Context, task string, gatewayURL string) {
	if err := push.Collectors(
		task, push.HostnameGroupingKey(),
		gatewayURL,
		histograms, counts, errcounts,
	); err != nil {
		fmt.Println("Could not push completion time to Pushgateway:", err)
	}
}

// DumpMetrics returns the metrics prometheus would return when collected
// as a string, for fun and testing
func DumpMetrics(ctx context.Context, task string) (string, error) {
	gatherer := prometheus.DefaultGatherer
	mfs, err := gatherer.Gather()
	if err != nil {
		return "", errors.Wrap(err, "gathering metrics")
	}

	buf := &bytes.Buffer{}
	enc := expfmt.NewEncoder(buf, expfmt.FmtText)

	for _, mf := range mfs {
		if err := enc.Encode(mf); err != nil {
			return buf.String(), errors.Wrap(err, "encoding metrics")
		}
	}
	return buf.String(), nil
}
