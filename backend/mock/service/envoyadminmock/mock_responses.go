package envoyadminmock

const clustersResponse = `{
 "cluster_statuses": [
  {
   "name": "statsd",
   "host_statuses": [
    {
     "address": {
      "socket_address": {
       "address": "127.0.0.1",
       "port_value": 8125
      }
     },
     "stats": [
      {
       "type": "GAUGE",
       "name": "cx_active"
      },
      {
       "type": "GAUGE",
       "name": "rq_active"
      }
     ],
     "health_status": {
      "failed_active_health_check": true,
      "eds_health_status": "HEALTHY"
     },
     "weight": 1,
     "locality": {}
    }
   ]
  },
  {
   "name": "local_service",
   "host_statuses": [
    {
     "address": {
      "pipe": {
       "path": "/run/gunicorn/gunicorn.sock"
      }
     },
     "stats": [
      {
       "type": "GAUGE",
       "name": "cx_active"
      },
      {
       "type": "GAUGE",
       "name": "rq_active"
      }
     ],
     "health_status": {
      "eds_health_status": "HEALTHY"
     },
     "weight": 1,
     "locality": {}
    }
   ]
  }
 ]
}
`

const configDumpResponse = `
{
 "configs": [
  {
   "@type": "type.googleapis.com/envoy.admin.v3.RoutesConfigDump",
   "static_route_configs": [
    {
     "route_config": {
      "@type": "type.googleapis.com/envoy.api.v2.RouteConfiguration"
     },
     "last_updated": "2020-05-28T00:25:48.062Z"
    }
   ]
  },
  {
   "@type": "type.googleapis.com/envoy.admin.v3.SecretsConfigDump"
  }
 ]
}
`

const listenersResponse = `{
 "listener_statuses": [
  {
   "name": "healthcheck_non_passthrough",
   "local_address": {
    "socket_address": {
     "address": "0.0.0.0",
     "port_value": 9212
    }
   }
  },
  {
   "name": "egress",
   "local_address": {
    "socket_address": {
     "address": "127.0.0.1",
     "port_value": 9009
    }
   }
  }
 ]
}
`

const serverInfoResponse = `
{
 "version": "dc76d44abceefd511d63c2d91d7439d5bd415639/1.15.0-dev/Modified/RELEASE/BoringSSL",
 "state": "PRE_INITIALIZING",
 "hot_restart_version": "11.104",
 "command_line_options": {
  "base_id": "0",
  "concurrency": 2,
  "config_path": "/etc/envoy/envoy.yaml",
  "config_yaml": "",
  "allow_unknown_static_fields": false,
  "reject_unknown_dynamic_fields": false,
  "ignore_unknown_dynamic_fields": false,
  "admin_address_path": "",
  "local_address_ip_version": "v4",
  "log_level": "info",
  "component_log_level": "",
  "log_format": "[%Y-%m-%d %T.%e][%t][%l][%n] [%g:%#] %v",
  "log_format_escaped": false,
  "log_path": "",
  "service_cluster": "example-staging-pdx",
  "service_node": "0c0c08040185ed0df",
  "service_zone": "us-east-1d",
  "mode": "Serve",
  "disable_hot_restart": false,
  "enable_mutex_tracing": false,
  "restart_epoch": 0,
  "cpuset_threads": false,
  "disabled_extensions": [],
  "bootstrap_version": 0,
  "hidden_envoy_deprecated_max_stats": "0",
  "hidden_envoy_deprecated_max_obj_name_len": "0",
  "file_flush_interval": "10s",
  "drain_time": "600s",
  "parent_shutdown_time": "900s"
 },
 "uptime_current_epoch": "24s",
 "uptime_all_epochs": "24s"
}
`
const runtimeResponse = `{
 "layers": [
  "root",
  "admin"
 ],
 "entries": {
  "routing.www2.brochure.conduce_pct": {
   "final_value": "100",
   "layer_values": [
    "100",
    ""
   ]
  },
  "circuit_breakers.local_service.default.max_requests": {
   "final_value": "265",
   "layer_values": [
    "265",
    ""
   ]
  }
 }
}
`

const statsResponse = `
cluster_manager.cds.version_text: ""
cluster_manager.cluster_updated: 0
filesystem.flushed_by_timer: 1
filesystem.write_completed: 1
http.admin.downstream_cx_active: 1
http.admin.downstream_cx_tx_bytes_total: 116030
listener_manager.lds.update_attempt: 0
overload.envoy.resource_monitors.fixed_heap.failed_updates: 0
runtime.admin_overrides_active: 0
runtime.num_keys: 362
runtime.num_layers: 3
runtime.override_dir_not_exists: 1
server.memory_allocated: 9981232
cluster.local_service.upstream_cx_connect_ms: No recorded values
server.initialization_time_ms: No recorded values
`
