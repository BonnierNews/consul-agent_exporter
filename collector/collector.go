package collector

import (
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	namespace = "consul"
)

type Collector interface {
	Collect(ch chan<- prometheus.Metric) (err error)
	Describe(ch chan<- *prometheus.Desc)
}

var (
	factories = make(map[string]func(*api.Client) Collector)

	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "exporter", "collector_duration_seconds"),
		"consul_exporter: Duration of a collection.",
		[]string{"collector"},
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "exporter", "collector_success"),
		"consul_exporter: Whether the collector was successful.",
		[]string{"collector"},
		nil,
	)
)

type ConsulCollector struct {
	collectors map[string]Collector
}

func NewConsulCollector(client *api.Client) ConsulCollector {
	return ConsulCollector{
		collectors: getCollectors(client),
	}
}

func getCollectors(client *api.Client) map[string]Collector {
	collectors := make(map[string]Collector)
	for name, factory := range factories {
		collectors[name] = factory(client)
	}
	return collectors
}

func (c ConsulCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc

	for _, ctr := range c.collectors {
		ctr.Describe(ch)
	}
}

func (c ConsulCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(c.collectors))
	for name, ctr := range c.collectors {
		go func(name string, ctr Collector) {
			execute(name, ctr, ch)
			wg.Done()
		}(name, ctr)
	}
	wg.Wait()
}
func execute(name string, c Collector, ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.Collect(ch)
	duration := time.Since(begin)
	var success float64

	if err != nil {
		log.Errorf("Collector %s failed after %fs: %s", name, duration.Seconds(), err)
		success = 0
	} else {
		log.Debugf("Collector %s succeeded after %fs.", name, duration.Seconds())
		success = 1
	}
	ch <- prometheus.MustNewConstMetric(
		scrapeDurationDesc,
		prometheus.GaugeValue,
		duration.Seconds(),
		name,
	)
	ch <- prometheus.MustNewConstMetric(
		scrapeSuccessDesc,
		prometheus.GaugeValue,
		success,
		name,
	)
}
