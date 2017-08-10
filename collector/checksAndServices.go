package collector

import (
	"fmt"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

func init() {
	factories["checksAndServices"] = newAgentChecksAndServicesCollector
}

func newAgentChecksAndServicesCollector(client *api.Client) Collector {
	const subsystem = "checks"
	return agentChecksAndServicesCollector{
		ConsulClient: client,

		NodeChecks: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "node_checks"),
			"Number of node checks",
			nil, nil),
		NodeCheckStatus: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "node_checks_status"),
			"Status of node checks, 1 if passing, 0 otherwise",
			[]string{"check"}, nil),

		Services: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "services"),
			"Number of services",
			nil, nil),
		ServiceChecks: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "service_checks"),
			"Number of service checks",
			[]string{"service_name", "service_id", "tags"}, nil),
		ServiceChecksFailing: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "service_checks_failing"),
			"Number of failing service checks",
			[]string{"service_name", "service_id", "tags"}, nil),
	}
}

type agentChecksAndServicesCollector struct {
	ConsulClient *api.Client

	NodeChecks           *prometheus.Desc
	NodeCheckStatus      *prometheus.Desc
	Services             *prometheus.Desc
	ServiceChecks        *prometheus.Desc
	ServiceChecksFailing *prometheus.Desc
}

func (c agentChecksAndServicesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.NodeChecks
	ch <- c.NodeCheckStatus
	ch <- c.ServiceChecks
	ch <- c.ServiceChecksFailing
}

func (c agentChecksAndServicesCollector) Collect(ch chan<- prometheus.Metric) error {
	checks, err := c.ConsulClient.Agent().Checks()
	if err != nil {
		log.Errorf("Could not fetch checks from Consul: %v", err)
		return err
	}
	services, err := c.ConsulClient.Agent().Services()
	if err != nil {
		log.Errorf("Could not fetch checks from Consul: %v", err)
		return err
	}

	serviceById := make(map[string]*api.AgentService)
	// Initialize this with all services, so we do not miss any services without checks
	serviceChecks := make(map[string][]*api.AgentCheck)
	for _, service := range services {
		serviceChecks[service.ID] = make([]*api.AgentCheck, 0)
		serviceById[service.ID] = service
	}

	ch <- prometheus.MustNewConstMetric(c.Services, prometheus.GaugeValue, float64(len(services)))

	nodeChecks := make([]*api.AgentCheck, 0)
	for _, check := range checks {
		if check.ServiceID == "" {
			nodeChecks = append(nodeChecks, check)
		} else if _, ok := serviceChecks[check.ServiceID]; ok {
			serviceChecks[check.ServiceID] = append(serviceChecks[check.ServiceID], check)
		} else {
			// During service registration, checks can be orphaned for a short while before they are cleaned up
			log.Debugf("Found check %s registered on service %s, but that service does not exist!", check.CheckID, check.ServiceID)
		}
	}

	ch <- prometheus.MustNewConstMetric(c.NodeChecks, prometheus.GaugeValue, float64(len(nodeChecks)))
	for _, check := range nodeChecks {
		ch <- prometheus.MustNewConstMetric(c.NodeCheckStatus, prometheus.GaugeValue, isPassingFloat(check.Status), check.Name)
	}

	for serviceId, checkList := range serviceChecks {
		service := serviceById[serviceId]
		ch <- prometheus.MustNewConstMetric(c.ServiceChecks, prometheus.GaugeValue, float64(len(checkList)), service.Service, service.ID, createTagsString(service.Tags))
		failing := 0.0
		for _, check := range checkList {
			failing += 1.0 - isPassingFloat(check.Status)
		}
		ch <- prometheus.MustNewConstMetric(c.ServiceChecksFailing, prometheus.GaugeValue, failing, service.Service, service.ID, createTagsString(service.Tags))
	}

	return nil
}

func isPassingFloat(status string) float64 {
	if status == "passing" {
		return 1.0
	}
	return 0.0
}

// Create a string of tags separated by commas, as well as having leading and
// trailing commas for easier regexp matching
func createTagsString(tags []string) string {
	return fmt.Sprintf(",%s,", strings.Join(tags, ","))
}
