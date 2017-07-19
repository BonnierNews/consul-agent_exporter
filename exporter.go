package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

const (
	namespace = "consul"
)

func main() {
	var (
		showVersion   = flag.Bool("version", false, "Print version information.")
		listenAddress = flag.String("telemetry.addr", ":9217", "host:port for exporter.")
		metricsPath   = flag.String("telemetry.path", "/metrics", "URL path for surfacing collected metrics.")
		consulAddr    = flag.String("consul.addr", "http://localhost:8500", "Address to Consul agent")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf(version.Print("consul_agent_exporter"))
		os.Exit(0)
	}

	consulConfig := &api.Config{
		Address: *consulAddr,
	}
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		// All errors from NewClient are unrecoverable
		log.Fatalf("Error configuring Consul client: %v", err)
	}

	collector := newAgentStatsCollector(consulClient)
	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK! (... No check against consul being done yet)")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	log.Infoln("Starting Consul agent exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	log.Infoln("Starting server on", *listenAddress)
	log.Fatalf("Cannot start Consul agent exporter: %s", http.ListenAndServe(*listenAddress, nil))
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
