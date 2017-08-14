package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

const (
	namespace = "twemproxy"
)

type Exporter struct {
	mutex   sync.RWMutex
	address string
	timeout time.Duration
	up      prometheus.Gauge
}

func NewExporter(address string, timeout time.Duration) (*Exporter, error) {

	return &Exporter{
		address: address,
		timeout: timeout,
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "up",
			Help: "Current health status of the server",
		},
		),
	}, nil
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.scrape()

	ch <- e.up
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up.Desc()
}

func (e *Exporter) scrape() {
	conn, err := net.DialTimeout("tcp", e.address, e.timeout)
	if err != nil {
		e.up.Set(0)
		log.Errorf("Connection to %v failed: %v", e.address, err)
		return
	}
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		e.up.Set(0)
		log.Errorln(err)
		return
	}

	e.up.Set(1)
}

func main() {

	var (
		twemproxyAddress = flag.String("twemproxy.address", "localhost:22222", "Address to query Twemproxy statistics")
		twemproxyTimeout = flag.Duration("twemproxy.timeout", 10*time.Second, "Timeout when connecting to Twemproxy")
		listenAddress    = flag.String("web.listen-address", ":9284", "Address to listen on for web interface and telemetry.")
		metricsPath      = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		showVersion      = flag.Bool("version", false, "Print version information.")
	)
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("twemproxy_exporter"))
		os.Exit(0)
	}

	log.Infoln("Starting twemproxy_exporter", version.Info())

	exporter, err := NewExporter(*twemproxyAddress, *twemproxyTimeout)
	if err != nil {
		log.Fatal(err)
	}

	prometheus.MustRegister(exporter)
	prometheus.MustRegister(version.NewCollector("twemproxy_exporter"))

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
