package collector

import (
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/mitchellh/mapstructure"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

func init() {
	factories["stats"] = newAgentStatsCollector
}

func newAgentStatsCollector(client *api.Client) Collector {
	const subsystem = "stats"
	return agentStatsCollector{
		ConsulClient: client,

		ConsulAgentCheckMonitors: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "agent_check_monitors"),
			"Number of monitors on checks",
			nil, nil),
		ConsulAgentCheckTTLs: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "agent_check_ttls"),
			"Number of checks with TTLs",
			nil, nil),
		ConsulAgentChecks: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "agent_checks"),
			"Number of checks registered on the agent",
			nil, nil),
		ConsulAgentServices: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "agent_services"),
			"Number of services registered on the agent",
			nil, nil),
		ConsulVersion: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "version"),
			"Version of the Consul agent",
			[]string{"version", "revision"}, nil),
		ConsulKnownDatacenters: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "cluster_known_datacenters"),
			"Datacenters known to the agent",
			nil, nil),
		ConsulKnownServers: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "cluster_known_servers"),
			"Servers known to the agent",
			nil, nil),
		ConsulLeader: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "cluster_is_leader"),
			"1 if the agent is the cluster leader",
			nil, nil),
		ConsulServer: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "cluster_is_server"),
			"1 if the agent is a server",
			nil, nil),
		ConsulRaftIndexApplied: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "raft_index_applied"),
			"Last index applied",
			nil, nil),
		ConsulRaftIndexCommitted: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "raft_index_committed"),
			"Last index committed",
			nil, nil),
		ConsulRaftFsmPending: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "raft_fsm_pending"),
			"Pending FSM operations",
			nil, nil),
		ConsulRaftLastIndexInLog: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "raft_last_index_in_log"),
			"Last index in log",
			nil, nil),
		ConsulRaftLastTermInLog: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "raft_last_term_in_log"),
			"Last leader term in log",
			nil, nil),
		ConsulRaftLastIndexInSnapshot: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "raft_last_index_in_snapshot"),
			"Last index in a snapshot",
			nil, nil),
		ConsulRaftLastTermInSnapshot: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "raft_last_term_in_snapshot"),
			"Last leader term in a snapshot",
			nil, nil),
		ConsulRaftConfigurationIndex: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "raft_configuration_index"),
			"Current configuration index",
			nil, nil),
		ConsulRaftPeers: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "raft_peers"),
			"Number of peers",
			nil, nil),
		ConsulRaftLastContact: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "last_contact_seconds"),
			"Time since the last contact with the leader",
			nil, nil),
		ConsulRaftState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "raft_state"),
			"The agent's current role (leader/follower)",
			[]string{"role"}, nil),
		ConsulRaftTerm: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "raft_term"),
			"Current term",
			nil, nil),
		ConsulGoroutines: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "goroutines"),
			"Number of goroutines",
			nil, nil),

		ConsulSerfEncrypted: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "serf_encrypted"),
			"1 if connections are encrypted",
			[]string{"ring"}, nil),
		ConsulSerfHealthScore: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "serf_health_score"),
			"The health score of the node",
			[]string{"ring"}, nil),

		ConsulSerfMembers: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "serf_members"),
			"Number of members in the cluster",
			[]string{"ring"}, nil),
		ConsulSerfFailed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "serf_members_failed"),
			"Number of failing members in cluster",
			[]string{"ring"}, nil),
		ConsulSerfLeft: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "serf_members_left"),
			"Number of members having left the cluster",
			[]string{"ring"}, nil),

		ConsulSerfEventQueue: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "serf_event_queue_length"),
			"Length of the events queue",
			[]string{"ring"}, nil),
		ConsulSerfIntentQueue: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "serf_intent_queue_length"),
			"Length of the intent queue",
			[]string{"ring"}, nil),
		ConsulSerfQueryQueue: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "serf_query_queue_length"),
			"Length of the query queue",
			[]string{"ring"}, nil),

		ConsulSerfLastEventsClock: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "serf_last_event_clock"),
			"Last event clock seen",
			[]string{"ring"}, nil),
		ConsulSerfLastMembershipEventClock: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "serf_last_membership_event_clock"),
			"Last membership event clock seen",
			[]string{"ring"}, nil),
		ConsulSerfLastQueryClock: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "serf_last_query_clock"),
			"Last query clock seen",
			[]string{"ring"}, nil),
	}
}

type Stats struct {
	Agent   AgentStats   `mapstructure:"agent"`
	Build   BuildInfo    `mapstructure:"build"`
	Consul  ConsulStats  `mapstructure:"consul"`
	Raft    RaftStats    `mapstructure:"raft,omitempty"`
	Runtime RuntimeStats `mapstructure:"runtime"`
	SerfLan SerfStats    `mapstructure:"serf_lan"`
	SerfWan SerfStats    `mapstructure:"serf_wan"`
}
type AgentStats struct {
	CheckMonitors uint64 `mapstructure:"check_monitors"`
	CheckTTLs     uint64 `mapstructure:"check_ttls"`
	Checks        uint64 `mapstructure:"checks"`
	Services      uint64 `mapstructure:"services"`
}
type BuildInfo struct {
	Revision string `mapstructure:"revision"`
	Version  string `mapstructure:"version"`
}
type ConsulStats struct {
	KnownDatacenters uint64 `mapstructure:"known_datacenters"`
	KnownServers     uint64 `mapstructure:"known_servers"`
	IsLeader         bool   `mapstructure:"leader"`
	IsServer         bool   `mapstructure:"server"`
}
type RaftStats struct {
	AppliedIndex             uint64 `mapstructure:"applied_index"`
	CommitIndex              uint64 `mapstructure:"commit_index"`
	FsmPending               uint64 `mapstructure:"fsm_pending"`
	LastContact              string `mapstructure:"last_contact"`
	LastLogIndex             uint64 `mapstructure:"last_log_index"`
	LastLogTerm              uint64 `mapstructure:"last_log_term"`
	LastSnapshotIndex        uint64 `mapstructure:"last_snapshot_index"`
	LastSnapshotTerm         uint64 `mapstructure:"last_snapshot_term"`
	LatestConfigurationIndex uint64 `mapstructure:"latest_configuration_index"`
	NumPeers                 uint64 `mapstructure:"num_peers"`
	State                    string `mapstructure:"state"`
	Term                     uint64 `mapstructure:"term"`
}
type RuntimeStats struct {
	NumberOfGoroutines uint64 `mapstructure:"goroutines"`
}
type SerfStats struct {
	Encrypted   bool   `mapstructure:"encrypted"`
	EventQueue  uint64 `mapstructure:"event_queue"`
	EventTime   uint64 `mapstructure:"event_time"`
	Failed      uint64 `mapstructure:"failed"`
	HealthScore uint64 `mapstructure:"health_score"`
	IntentQueue uint64 `mapstructure:"intent_queue"`
	Left        uint64 `mapstructure:"left"`
	MemberTime  uint64 `mapstructure:"member_time"`
	Members     uint64 `mapstructure:"members"`
	QueryQueue  uint64 `mapstructure:"query_queue"`
	QueryTime   uint64 `mapstructure:"query_time"`
}

type agentStatsCollector struct {
	ConsulClient *api.Client

	ConsulAgentCheckMonitors           *prometheus.Desc
	ConsulAgentCheckTTLs               *prometheus.Desc
	ConsulAgentChecks                  *prometheus.Desc
	ConsulAgentServices                *prometheus.Desc
	ConsulVersion                      *prometheus.Desc
	ConsulKnownDatacenters             *prometheus.Desc
	ConsulKnownServers                 *prometheus.Desc
	ConsulLeader                       *prometheus.Desc
	ConsulServer                       *prometheus.Desc
	ConsulRaftIndexApplied             *prometheus.Desc
	ConsulRaftIndexCommitted           *prometheus.Desc
	ConsulRaftFsmPending               *prometheus.Desc
	ConsulRaftLastIndexInLog           *prometheus.Desc
	ConsulRaftLastTermInLog            *prometheus.Desc
	ConsulRaftLastIndexInSnapshot      *prometheus.Desc
	ConsulRaftLastTermInSnapshot       *prometheus.Desc
	ConsulRaftConfigurationIndex       *prometheus.Desc
	ConsulRaftPeers                    *prometheus.Desc
	ConsulRaftLastContact              *prometheus.Desc
	ConsulRaftState                    *prometheus.Desc
	ConsulRaftTerm                     *prometheus.Desc
	ConsulGoroutines                   *prometheus.Desc
	ConsulSerfEncrypted                *prometheus.Desc
	ConsulSerfHealthScore              *prometheus.Desc
	ConsulSerfMembers                  *prometheus.Desc
	ConsulSerfFailed                   *prometheus.Desc
	ConsulSerfLeft                     *prometheus.Desc
	ConsulSerfEventQueue               *prometheus.Desc
	ConsulSerfIntentQueue              *prometheus.Desc
	ConsulSerfQueryQueue               *prometheus.Desc
	ConsulSerfLastEventsClock          *prometheus.Desc
	ConsulSerfLastMembershipEventClock *prometheus.Desc
	ConsulSerfLastQueryClock           *prometheus.Desc
}

func (c agentStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.ConsulAgentCheckMonitors
	ch <- c.ConsulAgentCheckTTLs
	ch <- c.ConsulAgentChecks
	ch <- c.ConsulAgentServices
	ch <- c.ConsulVersion
	ch <- c.ConsulKnownDatacenters
	ch <- c.ConsulKnownServers
	ch <- c.ConsulLeader
	ch <- c.ConsulServer
	ch <- c.ConsulRaftIndexApplied
	ch <- c.ConsulRaftIndexCommitted
	ch <- c.ConsulRaftFsmPending
	ch <- c.ConsulRaftLastIndexInLog
	ch <- c.ConsulRaftLastTermInLog
	ch <- c.ConsulRaftLastIndexInSnapshot
	ch <- c.ConsulRaftLastTermInSnapshot
	ch <- c.ConsulRaftConfigurationIndex
	ch <- c.ConsulRaftLastContact
	ch <- c.ConsulRaftPeers
	ch <- c.ConsulRaftState
	ch <- c.ConsulRaftTerm
	ch <- c.ConsulGoroutines
	ch <- c.ConsulSerfEncrypted
	ch <- c.ConsulSerfHealthScore
	ch <- c.ConsulSerfMembers
	ch <- c.ConsulSerfFailed
	ch <- c.ConsulSerfLeft
	ch <- c.ConsulSerfEventQueue
	ch <- c.ConsulSerfIntentQueue
	ch <- c.ConsulSerfQueryQueue
	ch <- c.ConsulSerfLastEventsClock
	ch <- c.ConsulSerfLastMembershipEventClock
	ch <- c.ConsulSerfLastQueryClock
}

func (c agentStatsCollector) Collect(ch chan<- prometheus.Metric) error {
	info, err := c.ConsulClient.Agent().Self()
	if err != nil {
		log.Errorf("Could not fetch stats from Consul: %v", err)
		return err
	}

	var s Stats
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &s,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		log.Errorf("Could not create decoder: %v", err)
		return err
	}
	err = decoder.Decode(info["stats"])
	if err != nil {
		log.Errorf("Could not decode stats from Consul: %v", err)
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.ConsulAgentCheckMonitors, prometheus.GaugeValue, float64(s.Agent.CheckMonitors))
	ch <- prometheus.MustNewConstMetric(c.ConsulAgentCheckTTLs, prometheus.GaugeValue, float64(s.Agent.CheckTTLs))
	ch <- prometheus.MustNewConstMetric(c.ConsulAgentChecks, prometheus.GaugeValue, float64(s.Agent.Checks))
	ch <- prometheus.MustNewConstMetric(c.ConsulAgentServices, prometheus.GaugeValue, float64(s.Agent.Services))

	ch <- prometheus.MustNewConstMetric(c.ConsulVersion, prometheus.GaugeValue, 1.0, s.Build.Version, s.Build.Revision)

	ch <- prometheus.MustNewConstMetric(c.ConsulServer, prometheus.GaugeValue, boolToFloat64(s.Consul.IsServer))

	ch <- prometheus.MustNewConstMetric(c.ConsulGoroutines, prometheus.GaugeValue, float64(s.Runtime.NumberOfGoroutines))

	ch <- prometheus.MustNewConstMetric(c.ConsulSerfEncrypted, prometheus.GaugeValue, boolToFloat64(s.SerfLan.Encrypted), "lan")
	ch <- prometheus.MustNewConstMetric(c.ConsulSerfHealthScore, prometheus.GaugeValue, float64(s.SerfLan.HealthScore), "lan")
	ch <- prometheus.MustNewConstMetric(c.ConsulSerfMembers, prometheus.GaugeValue, float64(s.SerfLan.Members), "lan")
	ch <- prometheus.MustNewConstMetric(c.ConsulSerfFailed, prometheus.GaugeValue, float64(s.SerfLan.Failed), "lan")
	ch <- prometheus.MustNewConstMetric(c.ConsulSerfLeft, prometheus.GaugeValue, float64(s.SerfLan.Left), "lan")
	ch <- prometheus.MustNewConstMetric(c.ConsulSerfEventQueue, prometheus.GaugeValue, float64(s.SerfLan.EventQueue), "lan")
	ch <- prometheus.MustNewConstMetric(c.ConsulSerfIntentQueue, prometheus.GaugeValue, float64(s.SerfLan.IntentQueue), "lan")
	ch <- prometheus.MustNewConstMetric(c.ConsulSerfQueryQueue, prometheus.GaugeValue, float64(s.SerfLan.QueryQueue), "lan")
	ch <- prometheus.MustNewConstMetric(c.ConsulSerfLastEventsClock, prometheus.CounterValue, float64(s.SerfLan.EventTime), "lan")
	ch <- prometheus.MustNewConstMetric(c.ConsulSerfLastMembershipEventClock, prometheus.CounterValue, float64(s.SerfLan.MemberTime), "lan")
	ch <- prometheus.MustNewConstMetric(c.ConsulSerfLastQueryClock, prometheus.CounterValue, float64(s.SerfLan.QueryTime), "lan")

	if s.Consul.IsServer {
		ch <- prometheus.MustNewConstMetric(c.ConsulLeader, prometheus.GaugeValue, boolToFloat64(s.Consul.IsLeader))
		ch <- prometheus.MustNewConstMetric(c.ConsulKnownDatacenters, prometheus.GaugeValue, float64(s.Consul.KnownDatacenters))

		ch <- prometheus.MustNewConstMetric(c.ConsulRaftIndexApplied, prometheus.GaugeValue, float64(s.Raft.AppliedIndex))
		ch <- prometheus.MustNewConstMetric(c.ConsulRaftIndexCommitted, prometheus.GaugeValue, float64(s.Raft.CommitIndex))
		ch <- prometheus.MustNewConstMetric(c.ConsulRaftFsmPending, prometheus.GaugeValue, float64(s.Raft.FsmPending))
		ch <- prometheus.MustNewConstMetric(c.ConsulRaftLastIndexInLog, prometheus.GaugeValue, float64(s.Raft.LastLogIndex))
		ch <- prometheus.MustNewConstMetric(c.ConsulRaftLastTermInLog, prometheus.GaugeValue, float64(s.Raft.LastLogTerm))
		ch <- prometheus.MustNewConstMetric(c.ConsulRaftLastIndexInSnapshot, prometheus.GaugeValue, float64(s.Raft.LastSnapshotIndex))
		ch <- prometheus.MustNewConstMetric(c.ConsulRaftLastTermInSnapshot, prometheus.GaugeValue, float64(s.Raft.LastSnapshotTerm))
		ch <- prometheus.MustNewConstMetric(c.ConsulRaftConfigurationIndex, prometheus.GaugeValue, float64(s.Raft.LatestConfigurationIndex))
		ch <- prometheus.MustNewConstMetric(c.ConsulRaftPeers, prometheus.GaugeValue, float64(s.Raft.NumPeers))
		ch <- prometheus.MustNewConstMetric(c.ConsulRaftState, prometheus.GaugeValue, 1.0, strings.ToLower(s.Raft.State))
		ch <- prometheus.MustNewConstMetric(c.ConsulRaftTerm, prometheus.GaugeValue, float64(s.Raft.Term))

		lastContact, err := time.ParseDuration(s.Raft.LastContact)
		if err != nil {
			log.Errorf("Could not parse %q as a duration for last contact", s.Raft.LastContact)
		} else {
			ch <- prometheus.MustNewConstMetric(c.ConsulRaftLastContact, prometheus.GaugeValue, float64(lastContact.Seconds()))
		}

		ch <- prometheus.MustNewConstMetric(c.ConsulSerfEncrypted, prometheus.GaugeValue, boolToFloat64(s.SerfWan.Encrypted), "wan")
		ch <- prometheus.MustNewConstMetric(c.ConsulSerfHealthScore, prometheus.GaugeValue, float64(s.SerfWan.HealthScore), "wan")
		ch <- prometheus.MustNewConstMetric(c.ConsulSerfMembers, prometheus.GaugeValue, float64(s.SerfWan.Members), "wan")
		ch <- prometheus.MustNewConstMetric(c.ConsulSerfFailed, prometheus.GaugeValue, float64(s.SerfWan.Failed), "wan")
		ch <- prometheus.MustNewConstMetric(c.ConsulSerfLeft, prometheus.GaugeValue, float64(s.SerfWan.Left), "wan")
		ch <- prometheus.MustNewConstMetric(c.ConsulSerfEventQueue, prometheus.GaugeValue, float64(s.SerfWan.EventQueue), "wan")
		ch <- prometheus.MustNewConstMetric(c.ConsulSerfIntentQueue, prometheus.GaugeValue, float64(s.SerfWan.IntentQueue), "wan")
		ch <- prometheus.MustNewConstMetric(c.ConsulSerfQueryQueue, prometheus.GaugeValue, float64(s.SerfWan.QueryQueue), "wan")
		ch <- prometheus.MustNewConstMetric(c.ConsulSerfLastEventsClock, prometheus.CounterValue, float64(s.SerfWan.EventTime), "wan")
		ch <- prometheus.MustNewConstMetric(c.ConsulSerfLastMembershipEventClock, prometheus.CounterValue, float64(s.SerfWan.MemberTime), "wan")
		ch <- prometheus.MustNewConstMetric(c.ConsulSerfLastQueryClock, prometheus.CounterValue, float64(s.SerfWan.QueryTime), "wan")
	} else {
		ch <- prometheus.MustNewConstMetric(c.ConsulKnownServers, prometheus.GaugeValue, float64(s.Consul.KnownServers))
	}

	return nil
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
