package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/prometheus/common/version"
)

func main() {

	var (
		showVersion = flag.Bool("version", false, "Print version information.")
	)
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("twemproxy_exporter"))
		os.Exit(0)
	}
}
