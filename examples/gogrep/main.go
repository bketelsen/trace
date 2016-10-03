package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bketelsen/trace"
)

func main() {
	// optional: set metrics Namespace
	trace.Namespace = "bketelsen"
	// optional: set metrics Subsystem
	trace.Subsystem = "example"

	t, ctx := trace.NewContext(context.Background(), "gogrep", "main")

	// "gogrep" is the name of this application, and used as the "task" name
	// in metrics reporting
	// uncomment to push to a prometheus gateway
	//defer trace.PushMetrics(ctx, "gogrep", "http://mypushgateway:9091")

	duration := flag.Duration("timeout", 500*time.Millisecond, "timeout in milliseconds")
	flag.Usage = func() {
		fmt.Printf("%s by Brian Ketelsen\n", os.Args[0])
		fmt.Println("Usage:")
		fmt.Printf("	gogrep [flags] path pattern \n")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		t.SetError()
		t.Finish()
		fmt.Println(trace.DumpMetrics(ctx, "gogrep"))
		os.Exit(-1)
	}
	path := flag.Arg(0)
	pattern := flag.Arg(1)
	ctx, cf := context.WithTimeout(ctx, *duration)
	defer cf()
	m, err := search(ctx, path, pattern)
	if err != nil {
		t.SetError()
		t.Finish()
		fmt.Println(trace.DumpMetrics(ctx, "gogrep"))
		log.Fatal(err)
	}
	for _, name := range m {
		t.LazyLog(trace.LogMessage("found", trace.KeyValue("file", name)), false)
	}
	t.LazyPrintf("hit count %d", len(m))
	t.LazyLog(trace.LogMessage("finished", trace.KeyValue("hits", len(m))), false)

	// instead of calling in defer at the top of the function,
	// call here so it shows up in metrics output
	t.Finish()
	fmt.Println("\n\nMetrics Output:")
	fmt.Println(trace.DumpMetrics(ctx, "gogrep"))

}
