package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
)

func main() {
	var (
		address        = flag.String("address", "127.0.0.1:8118", "address to listen on")
		rps            = flag.Int64("rps", 16, "request per second")
		numbers        = flag.String("numbers", "", "comma separated list of phone numbers")
		printRequests  = flag.Bool("print-requests", false, "print requests")
		printResponses = flag.Bool("print-responses", false, "print responses")
		printProgress  = flag.Bool("print-progress", false, "print progress")
	)
	flag.Parse()

	if *numbers == "" {
		flag.Usage()
		fmt.Fprint(os.Stderr, "\n  Missing parameter: -numbers\n")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt)
	go func() {
		<-sigC
		cancel()
	}()

	startServer(
		ctx,
		*address,
		*rps,
		strings.Split(*numbers, ","),
		&printOptions{
			Requests:  *printRequests,
			Responses: *printResponses,
			Progress:  *printProgress,
		},
	)
}
