package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

func main() {

	var (
		listenAddress = flag.String("web.listen-address", ":9284", "Address to listen on for web interface and telemetry.")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		showVersion   = flag.Bool("version", false, "Print version information.")
	)
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("twemproxy_exporter"))
		os.Exit(0)
	}

	log.Infoln("Starting twemproxy_exporter", version.Info())

	log.Infoln("Listening on", *listenAddress)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
	<head><title>Twemproxy Exporter</title></head>
	<body>
		<h1>Twemproxy Exporter</h1>
		<p><a href="` + *metricsPath + `">Metrics</a></p>
	</body>
</html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
