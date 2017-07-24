# Metrics
Below is a short listing of the metrics exposed by the exporter, by collector.

## Checks and Services
This collector exposes the status of services and checks registered on the agent.

Name | Description | Type | Labels
-----|-------------|------|-------
`consul_checks_node_checks` | Number of node checks | gauge | - |
`consul_checks_service_checks` | Number of service checks | gauge | service_name,service_id,tags |
`consul_checks_service_checks_failing` | Number of failing service checks | gauge | service_name,service_id,tags |
`consul_checks_services` | Number of services | gauge | - |

## Stats
This collector exposes the equivalent of what `consul info` returns. There are different metrics depending on if the agent is a server node or not.

### Always
_(Note that a server node will have two values of the `ring` label: "lan" and "wan", while an agent will only have "lan")_

Name | Description | Type | Labels
-----|-------------|------|-------
`consul_stats_agent_check_monitors` | Number of monitors on checks | gauge | - |
`consul_stats_agent_check_ttls` | Number of checks with TTLs | gauge | - |
`consul_stats_agent_checks` | Number of checks registered on the agent | gauge | - |
`consul_stats_agent_services` | Number of services registered on the agent | gauge | - |
`consul_stats_cluster_is_server` | 1 if the agent is a server | gauge | - |
`consul_stats_goroutines` | Number of goroutines | gauge | - |
`consul_stats_serf_encrypted` | 1 if connections are encrypted | gauge | ring |
`consul_stats_serf_event_queue_length` | Length of the events queue | gauge | ring |
`consul_stats_serf_health_score` | The health score of the node | gauge | ring |
`consul_stats_serf_intent_queue_length` | Length of the intent queue | gauge | ring |
`consul_stats_serf_last_event_clock` | Last event clock seen | counter | ring |
`consul_stats_serf_last_membership_event_clock` | Last membership event clock seen | counter | ring |
`consul_stats_serf_last_query_clock` | Last query clock seen | counter | ring |
`consul_stats_serf_members` | Number of members in the cluster | gauge | ring |
`consul_stats_serf_members_failed` | Number of failing members in cluster | gauge | ring |
`consul_stats_serf_members_left` | Number of members having left the cluster | gauge | ring |
`consul_stats_serf_query_queue_length` | Length of the query queue | gauge | ring |
`consul_version` | Version of the Consul agent | gauge | revision,version |

### Server only
Name | Description | Type | Labels
-----|-------------|------|-------
`consul_stats_cluster_is_leader` | 1 if the agent is the cluster leader | gauge | - |
`consul_stats_cluster_known_datacenters` | Datacenters known to the agent | gauge | - |
`consul_stats_raft_configuration_index` | Current configuration index | gauge | - |
`consul_stats_raft_fsm_pending` | Pending FSM operations | gauge | - |
`consul_stats_raft_index_applied` | Last index applied | gauge | - |
`consul_stats_raft_index_committed` | Last index committed | gauge | - |
`consul_stats_raft_last_index_in_log` | Last index in log | gauge | - |
`consul_stats_raft_last_index_in_snapshot` | Last index in a snapshot | gauge | - |
`consul_stats_raft_last_term_in_log` | Last leader term in log | gauge | - |
`consul_stats_raft_last_term_in_snapshot` | Last leader term in a snapshot | gauge | - |
`consul_stats_raft_peers` | Number of peers | gauge | - |
`consul_stats_raft_state` | The agent's current role (leader/follower) | gauge | role |
`consul_stats_raft_term` | Current term | gauge | - |

### Non-server only
Name | Description | Type | Labels
-----|-------------|------|-------
`consul_stats_cluster_known_servers` | Servers known to the agent | gauge | - |


## Meta
Data about the scrape itself.

Name | Description | Type | Labels
-----|-------------|------|-------
`consul_exporter_collector_duration_seconds` | Duration of a collection. | gauge | collector |
`consul_exporter_collector_success` | Whether the collector was successful. | gauge | collector |

