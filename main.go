package main

import (
	"flag"
	"fmt"
	"json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/rokett/citrix-netscaler-exporter/netscaler"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Configuration struct {
	url string
	username string
	password string
	bindPort int
	ignoreCert bool
}

var (
	app        = "Citrix-NetScaler-Exporter"
	version    string
	build      string
	config     = flag.String("config", "", "Configuration file. E.g. citrix-netscaler-exporter.json")
	url        = flag.String("url", "", "Base URL of the NetScaler management interface.  Normally something like https://my-netscaler.something.x")
	username   = flag.String("username", "", "Username with which to connect to the NetScaler API")
	password   = flag.String("password", "", "Password with which to connect to the NetScaler API")
	bindPort   = flag.Int("bind_port", 9280, "Port to bind the exporter endpoint to")
	versionFlg = flag.Bool("version", false, "Display application version")
	ignoreCert = flag.Bool("ignore-cert", false, "Ignore certificate errors; use with caution")
	logger     log.Logger

	nsInstance string

//filename is the path to the json config file
if config != nil {
  file, err := os.Open(filename) if err != nil {  return err }
  decoder := json.NewDecoder(file)
  err = decoder.Decode(&configuration)
  if err != nil {  return err }
}

	modelID = prometheus.NewDesc(
		"model_id",
		"NetScaler model - reflects the bandwidth available; for example VPX 10 would report as 10.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	mgmtCPUUsage = prometheus.NewDesc(
		"mgmt_cpu_usage",
		"Current CPU utilisation for management",
		[]string{
			"ns_instance",
		},
		nil,
	)

	pktCPUUsage = prometheus.NewDesc(
		"pkt_cpu_usage",
		"Current CPU utilisation for packet engines, excluding management",
		[]string{
			"ns_instance",
		},
		nil,
	)

	memUsage = prometheus.NewDesc(
		"mem_usage",
		"Current memory utilisation",
		[]string{
			"ns_instance",
		},
		nil,
	)

	flashPartitionUsage = prometheus.NewDesc(
		"flash_partition_usage",
		"Used space in /flash partition of the disk, as a percentage.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	varPartitionUsage = prometheus.NewDesc(
		"var_partition_usage",
		"Used space in /var partition of the disk, as a percentage. ",
		[]string{
			"ns_instance",
		},
		nil,
	)

	totRxMB = prometheus.NewDesc(
		"total_received_mb",
		"Total number of Megabytes received by the NetScaler appliance",
		[]string{
			"ns_instance",
		},
		nil,
	)

	totTxMB = prometheus.NewDesc(
		"total_transmit_mb",
		"Total number of Megabytes transmitted by the NetScaler appliance",
		[]string{
			"ns_instance",
		},
		nil,
	)

	httpRequests = prometheus.NewDesc(
		"http_requests",
		"Total number of HTTP requests received",
		[]string{
			"ns_instance",
		},
		nil,
	)

	httpResponses = prometheus.NewDesc(
		"http_responses",
		"Total number of HTTP responses sent",
		[]string{
			"ns_instance",
		},
		nil,
	)

	tcpCurrentClientConnections = prometheus.NewDesc(
		"tcp_current_client_connections",
		"Client connections, including connections in the Opening, Established, and Closing state.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	tcpCurrentClientConnectionsEstablished = prometheus.NewDesc(
		"tcp_current_client_connections_established",
		"Current client connections in the Established state, which indicates that data transfer can occur between the NetScaler and the client.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	tcpCurrentServerConnections = prometheus.NewDesc(
		"tcp_current_server_connections",
		"Server connections, including connections in the Opening, Established, and Closing state.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	tcpCurrentServerConnectionsEstablished = prometheus.NewDesc(
		"tcp_current_server_connections_established",
		"Current server connections in the Established state, which indicates that data transfer can occur between the NetScaler and the server.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	interfacesRxBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_received_bytes",
			Help: "Number of bytes received by specific interfaces.",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesTxBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_transmitted_bytes",
			Help: "Number of bytes transmitted by specific interfaces.",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesRxPackets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_received_packets",
			Help: "Number of packets received by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesTxPackets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_transmitted_packets",
			Help: "Number of packets transmitted by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesJumboPacketsRx = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_jumbo_packets_received",
			Help: "Number of bytes received by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesJumboPacketsTx = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_jumbo_packets_transmitted",
			Help: "Number of jumbo packets transmitted by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesErrorPacketsRx = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_error_packets_received",
			Help: "Number of error packets received by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	virtualServersWaitingRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_waiting_requests",
			Help: "Number of requests waiting on a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersHealth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_health",
			Help: "Percentage of UP services bound to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersInactiveServices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_inactive_services",
			Help: "Number of inactive services bound to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersActiveServices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_active_services",
			Help: "Number of active services bound to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersTotalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_hits",
			Help: "Total virtual server hits",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_requests",
			Help: "Total virtual server requests",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_responses",
			Help: "Total virtual server responses",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_request_bytes",
			Help: "Total virtual server request bytes",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)
	virtualServersTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_response_bytes",
			Help: "Total virtual server response bytes",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersCurrentClientConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_current_client_connections",
			Help: "Number of current client connections on a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersCurrentServerConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_current_server_connections",
			Help: "Number of current connections to the actual servers behind the specific virtual server.",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	servicesThroughput = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_throughput",
			Help: "Number of bytes received or sent by this service (Mbps)",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesAvgTTFB = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_average_time_to_first_byte",
			Help: "Average TTFB between the NetScaler appliance and the server. TTFB is the time interval between sending the request packet to a service and receiving the first response from the service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_state",
			Help: "Current state of the service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_total_requests",
			Help: "Total number of requests received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_total_responses",
			Help: "Total number of responses received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_total_request_bytes",
			Help: "Total number of request bytes received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_total_response_bytes",
			Help: "Total number of response bytes received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesCurrentClientConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_current_client_connections",
			Help: "Number of current client connections",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesSurgeCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_surge_count",
			Help: "Number of requests in the surge queue",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesCurrentServerConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_current_server_connections",
			Help: "Number of current connections to the actual servers",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesServerEstablishedConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_server_established_connections",
			Help: "Number of server connections in ESTABLISHED state",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesCurrentReusePool = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_current_reuse_pool",
			Help: "Number of requests in the idle queue/reuse pool.",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesMaxClients = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_max_clients",
			Help: "Maximum open connections allowed on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesCurrentLoad = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_current_load",
			Help: "Load on the service that is calculated from the bound load based monitor",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesVirtualServerServiceHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_virtual_server_service_hits",
			Help: "Number of times that the service has been provided",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesActiveTransactions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_active_transactions",
			Help: "Number of active transactions handled by this service. (Including those in the surge queue.) Active Transaction means number of transactions currently served by the server including those waiting in the SurgeQ",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	serviceGroupsState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_state",
			Help: "Current state of the server",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsAvgTTFB = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_average_time_to_first_byte",
			Help: "Average TTFB between the NetScaler appliance and the server.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "servicegroup_total_requests",
			Help: "Total number of requests received on this service",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "servicegroup_total_responses",
			Help: "Number of responses received on this service.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "servicegroup_total_request_bytes",
			Help: "Total number of request bytes received on this service",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "servicegroup_total_response_bytes",
			Help: "Number of response bytes received by this service",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsCurrentClientConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_current_client_connections",
			Help: "Number of current client connections.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsSurgeCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_surge_count",
			Help: "Number of requests in the surge queue.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsCurrentServerConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_current_server_connections",
			Help: "Number of current connections to the actual servers behind the virtual server.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsServerEstablishedConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_server_established_connections",
			Help: "Number of server connections in ESTABLISHED state.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsCurrentReusePool = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_current_reuse_pool",
			Help: "Number of requests in the idle queue/reuse pool.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsMaxClients = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_max_clients",
			Help: "Maximum open connections allowed on this service.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	gslbServicesState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_state",
			Help: "Current state of the service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	gslbServicesTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_service_total_requests",
			Help: "Total number of requests received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	gslbServicesTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_service_total_responses",
			Help: "Total number of responses received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	gslbServicesTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_service_total_request_bytes",
			Help: "Total number of request bytes received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	gslbServicesTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_service_total_response_bytes",
			Help: "Total number of response bytes received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	gslbServicesCurrentClientConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_current_client_connections",
			Help: "Number of current client connections",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	gslbServicesCurrentServerConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_current_server_connections",
			Help: "Number of current connections to the actual servers",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	gslbServicesEstablishedConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_established_connections",
			Help: "Number of server connections in ESTABLISHED state",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	gslbServicesCurrentLoad = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_current_load",
			Help: "Load on the service that is calculated from the bound load based monitor",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	gslbServicesVirtualServerServiceHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_service_virtual_server_service_hits",
			Help: "Number of times that the service has been provided",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	gslbVirtualServersHealth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_health",
			Help: "Percentage of UP services bound to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	gslbVirtualServersInactiveServices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_inactive_services",
			Help: "Number of inactive services bound to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	gslbVirtualServersActiveServices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_active_services",
			Help: "Number of active services bound to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	gslbVirtualServersTotalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_virtual_servers_total_hits",
			Help: "Total virtual server hits",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	gslbVirtualServersTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_virtual_servers_total_requests",
			Help: "Total virtual server requests",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	gslbVirtualServersTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_virtual_servers_total_responses",
			Help: "Total virtual server responses",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	gslbVirtualServersTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_virtual_servers_total_request_bytes",
			Help: "Total virtual server request bytes",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)
	gslbVirtualServersTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_virtual_servers_total_response_bytes",
			Help: "Total virtual server response bytes",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	gslbVirtualServersCurrentClientConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_current_client_connections",
			Help: "Number of current client connections on a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	gslbVirtualServersCurrentServerConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_current_server_connections",
			Help: "Number of current connections to the actual servers behind the specific virtual server.",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_state",
			Help: "Current state of the server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersTotalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_total_hits",
			Help: "Total virtual server hits",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_total_requests",
			Help: "Total virtual server requests",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_total_responses",
			Help: "Total virtual server responses",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_total_request_bytes",
			Help: "Total virtual server request bytes",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_total_response_bytes",
			Help: "Total virtual server response bytes",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersCurrentClientConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_current_client_connections",
			Help: "Number of current client connections on a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersCurrentServerConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_current_server_connections",
			Help: "Number of current connections to the actual servers behind the specific virtual server.",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersEstablishedConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_established_connections",
			Help: "Number of client connections in ESTABLISHED state.",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersTotalPacketsReceived = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_total_packets_received",
			Help: "Total number of packets received",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersTotalPacketsSent = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_total_packets_sent",
			Help: "Total number of packets sent.",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersTotalSpillovers = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_total_spillovers",
			Help: "Number of times vserver experienced spill over.",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersDeferredRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_deferred_requests",
			Help: "Number of deferred request on this vserver",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersNumberInvalidRequestResponse = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_number_invalid_request_response",
			Help: "Number invalid requests/responses on this vserver",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersNumberInvalidRequestResponseDropped = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_number_invalid_request_response_dropped",
			Help: "Number invalid requests/responses dropped on this vserver",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersTotalVServerDownBackupHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_total_vserver_down_backup_hits",
			Help: "Number of times traffic was diverted to backup vserver since primary vserver was DOWN.",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersCurrentMultipathSessions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_current_multipath_sessions",
			Help: "Current Multipath TCP sessions",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersCurrentMultipathSubflows = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_current_multipath_subflows",
			Help: "Current Multipath TCP subflows",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)
)

// Exporter represents the metrics exported to Prometheus
type Exporter struct {
	modelID                                             *prometheus.Desc
	mgmtCPUUsage                                        *prometheus.Desc
	memUsage                                            *prometheus.Desc
	pktCPUUsage                                         *prometheus.Desc
	flashPartitionUsage                                 *prometheus.Desc
	varPartitionUsage                                   *prometheus.Desc
	totRxMB                                             *prometheus.Desc
	totTxMB                                             *prometheus.Desc
	httpRequests                                        *prometheus.Desc
	httpResponses                                       *prometheus.Desc
	tcpCurrentClientConnections                         *prometheus.Desc
	tcpCurrentClientConnectionsEstablished              *prometheus.Desc
	tcpCurrentServerConnections                         *prometheus.Desc
	tcpCurrentServerConnectionsEstablished              *prometheus.Desc
	interfacesRxBytes                                   *prometheus.GaugeVec
	interfacesTxBytes                                   *prometheus.GaugeVec
	interfacesRxPackets                                 *prometheus.GaugeVec
	interfacesTxPackets                                 *prometheus.GaugeVec
	interfacesJumboPacketsRx                            *prometheus.GaugeVec
	interfacesJumboPacketsTx                            *prometheus.GaugeVec
	interfacesErrorPacketsRx                            *prometheus.GaugeVec
	virtualServersWaitingRequests                       *prometheus.GaugeVec
	virtualServersHealth                                *prometheus.GaugeVec
	virtualServersInactiveServices                      *prometheus.GaugeVec
	virtualServersActiveServices                        *prometheus.GaugeVec
	virtualServersTotalHits                             *prometheus.CounterVec
	virtualServersTotalRequests                         *prometheus.CounterVec
	virtualServersTotalResponses                        *prometheus.CounterVec
	virtualServersTotalRequestBytes                     *prometheus.CounterVec
	virtualServersTotalResponseBytes                    *prometheus.CounterVec
	virtualServersCurrentClientConnections              *prometheus.GaugeVec
	virtualServersCurrentServerConnections              *prometheus.GaugeVec
	servicesThroughput                                  *prometheus.CounterVec
	servicesAvgTTFB                                     *prometheus.GaugeVec
	servicesState                                       *prometheus.GaugeVec
	servicesTotalRequests                               *prometheus.CounterVec
	servicesTotalResponses                              *prometheus.CounterVec
	servicesTotalRequestBytes                           *prometheus.CounterVec
	servicesTotalResponseBytes                          *prometheus.CounterVec
	servicesCurrentClientConns                          *prometheus.GaugeVec
	servicesSurgeCount                                  *prometheus.GaugeVec
	servicesCurrentServerConns                          *prometheus.GaugeVec
	servicesServerEstablishedConnections                *prometheus.GaugeVec
	servicesCurrentReusePool                            *prometheus.GaugeVec
	servicesMaxClients                                  *prometheus.GaugeVec
	servicesCurrentLoad                                 *prometheus.GaugeVec
	servicesVirtualServerServiceHits                    *prometheus.CounterVec
	servicesActiveTransactions                          *prometheus.GaugeVec
	serviceGroupsState                                  *prometheus.GaugeVec
	serviceGroupsAvgTTFB                                *prometheus.GaugeVec
	serviceGroupsTotalRequests                          *prometheus.CounterVec
	serviceGroupsTotalResponses                         *prometheus.CounterVec
	serviceGroupsTotalRequestBytes                      *prometheus.CounterVec
	serviceGroupsTotalResponseBytes                     *prometheus.CounterVec
	serviceGroupsCurrentClientConnections               *prometheus.GaugeVec
	serviceGroupsSurgeCount                             *prometheus.GaugeVec
	serviceGroupsCurrentServerConnections               *prometheus.GaugeVec
	serviceGroupsServerEstablishedConnections           *prometheus.GaugeVec
	serviceGroupsCurrentReusePool                       *prometheus.GaugeVec
	serviceGroupsMaxClients                             *prometheus.GaugeVec
	gslbServicesState                                   *prometheus.GaugeVec
	gslbServicesTotalRequests                           *prometheus.CounterVec
	gslbServicesTotalResponses                          *prometheus.CounterVec
	gslbServicesTotalRequestBytes                       *prometheus.CounterVec
	gslbServicesTotalResponseBytes                      *prometheus.CounterVec
	gslbServicesCurrentClientConns                      *prometheus.GaugeVec
	gslbServicesCurrentServerConns                      *prometheus.GaugeVec
	gslbServicesCurrentLoad                             *prometheus.GaugeVec
	gslbServicesVirtualServerServiceHits                *prometheus.CounterVec
	gslbServicesEstablishedConnections                  *prometheus.GaugeVec
	gslbVirtualServersHealth                            *prometheus.GaugeVec
	gslbVirtualServersInactiveServices                  *prometheus.GaugeVec
	gslbVirtualServersActiveServices                    *prometheus.GaugeVec
	gslbVirtualServersTotalHits                         *prometheus.CounterVec
	gslbVirtualServersTotalRequests                     *prometheus.CounterVec
	gslbVirtualServersTotalResponses                    *prometheus.CounterVec
	gslbVirtualServersTotalRequestBytes                 *prometheus.CounterVec
	gslbVirtualServersTotalResponseBytes                *prometheus.CounterVec
	gslbVirtualServersCurrentClientConnections          *prometheus.GaugeVec
	gslbVirtualServersCurrentServerConnections          *prometheus.GaugeVec
	csVirtualServersState                               *prometheus.GaugeVec
	csVirtualServersTotalHits                           *prometheus.CounterVec
	csVirtualServersTotalRequests                       *prometheus.CounterVec
	csVirtualServersTotalResponses                      *prometheus.CounterVec
	csVirtualServersTotalRequestBytes                   *prometheus.CounterVec
	csVirtualServersTotalResponseBytes                  *prometheus.CounterVec
	csVirtualServersCurrentClientConnections            *prometheus.GaugeVec
	csVirtualServersCurrentServerConnections            *prometheus.GaugeVec
	csVirtualServersEstablishedConnections              *prometheus.GaugeVec
	csVirtualServersTotalPacketsReceived                *prometheus.CounterVec
	csVirtualServersTotalPacketsSent                    *prometheus.CounterVec
	csVirtualServersTotalSpillovers                     *prometheus.CounterVec
	csVirtualServersDeferredRequests                    *prometheus.CounterVec
	csVirtualServersNumberInvalidRequestResponse        *prometheus.CounterVec
	csVirtualServersNumberInvalidRequestResponseDropped *prometheus.CounterVec
	csVirtualServersTotalVServerDownBackupHits          *prometheus.CounterVec
	csVirtualServersCurrentMultipathSessions            *prometheus.GaugeVec
	csVirtualServersCurrentMultipathSubflows            *prometheus.GaugeVec
}

// NewExporter initialises the exporter
func NewExporter() (*Exporter, error) {
	return &Exporter{
		modelID:                                             modelID,
		mgmtCPUUsage:                                        mgmtCPUUsage,
		memUsage:                                            memUsage,
		pktCPUUsage:                                         pktCPUUsage,
		flashPartitionUsage:                                 flashPartitionUsage,
		varPartitionUsage:                                   varPartitionUsage,
		totRxMB:                                             totRxMB,
		totTxMB:                                             totTxMB,
		httpRequests:                                        httpRequests,
		httpResponses:                                       httpResponses,
		tcpCurrentClientConnections:                         tcpCurrentClientConnections,
		tcpCurrentClientConnectionsEstablished:              tcpCurrentClientConnectionsEstablished,
		tcpCurrentServerConnections:                         tcpCurrentServerConnections,
		tcpCurrentServerConnectionsEstablished:              tcpCurrentServerConnectionsEstablished,
		interfacesRxBytes:                                   interfacesRxBytes,
		interfacesTxBytes:                                   interfacesTxBytes,
		interfacesRxPackets:                                 interfacesRxPackets,
		interfacesTxPackets:                                 interfacesTxPackets,
		interfacesJumboPacketsRx:                            interfacesJumboPacketsRx,
		interfacesJumboPacketsTx:                            interfacesJumboPacketsTx,
		interfacesErrorPacketsRx:                            interfacesErrorPacketsRx,
		virtualServersWaitingRequests:                       virtualServersWaitingRequests,
		virtualServersHealth:                                virtualServersHealth,
		virtualServersInactiveServices:                      virtualServersInactiveServices,
		virtualServersActiveServices:                        virtualServersActiveServices,
		virtualServersTotalHits:                             virtualServersTotalHits,
		virtualServersTotalRequests:                         virtualServersTotalRequests,
		virtualServersTotalResponses:                        virtualServersTotalResponses,
		virtualServersTotalRequestBytes:                     virtualServersTotalRequestBytes,
		virtualServersTotalResponseBytes:                    virtualServersTotalResponseBytes,
		virtualServersCurrentClientConnections:              virtualServersCurrentClientConnections,
		virtualServersCurrentServerConnections:              virtualServersCurrentServerConnections,
		servicesThroughput:                                  servicesThroughput,
		servicesAvgTTFB:                                     servicesAvgTTFB,
		servicesState:                                       servicesState,
		servicesTotalRequests:                               servicesTotalRequests,
		servicesTotalResponses:                              servicesTotalResponses,
		servicesTotalRequestBytes:                           servicesTotalRequestBytes,
		servicesTotalResponseBytes:                          servicesTotalResponseBytes,
		servicesCurrentClientConns:                          servicesCurrentClientConns,
		servicesSurgeCount:                                  servicesSurgeCount,
		servicesCurrentServerConns:                          servicesCurrentServerConns,
		servicesServerEstablishedConnections:                servicesServerEstablishedConnections,
		servicesCurrentReusePool:                            servicesCurrentReusePool,
		servicesMaxClients:                                  servicesMaxClients,
		servicesCurrentLoad:                                 servicesCurrentLoad,
		servicesVirtualServerServiceHits:                    servicesVirtualServerServiceHits,
		servicesActiveTransactions:                          servicesActiveTransactions,
		serviceGroupsState:                                  serviceGroupsState,
		serviceGroupsAvgTTFB:                                serviceGroupsAvgTTFB,
		serviceGroupsTotalRequests:                          serviceGroupsTotalRequests,
		serviceGroupsTotalResponses:                         serviceGroupsTotalResponses,
		serviceGroupsTotalRequestBytes:                      serviceGroupsTotalRequestBytes,
		serviceGroupsTotalResponseBytes:                     serviceGroupsTotalResponseBytes,
		serviceGroupsCurrentClientConnections:               serviceGroupsCurrentClientConnections,
		serviceGroupsSurgeCount:                             serviceGroupsSurgeCount,
		serviceGroupsCurrentServerConnections:               serviceGroupsCurrentServerConnections,
		serviceGroupsServerEstablishedConnections:           serviceGroupsServerEstablishedConnections,
		serviceGroupsCurrentReusePool:                       serviceGroupsCurrentReusePool,
		serviceGroupsMaxClients:                             serviceGroupsMaxClients,
		gslbServicesState:                                   gslbServicesState,
		gslbServicesTotalRequests:                           gslbServicesTotalRequests,
		gslbServicesTotalResponses:                          gslbServicesTotalResponses,
		gslbServicesTotalRequestBytes:                       gslbServicesTotalRequestBytes,
		gslbServicesTotalResponseBytes:                      gslbServicesTotalResponseBytes,
		gslbServicesCurrentClientConns:                      gslbServicesCurrentClientConns,
		gslbServicesCurrentServerConns:                      gslbServicesCurrentServerConns,
		gslbServicesCurrentLoad:                             gslbServicesCurrentLoad,
		gslbServicesVirtualServerServiceHits:                gslbServicesVirtualServerServiceHits,
		gslbServicesEstablishedConnections:                  gslbServicesEstablishedConnections,
		gslbVirtualServersHealth:                            gslbVirtualServersHealth,
		gslbVirtualServersInactiveServices:                  gslbVirtualServersInactiveServices,
		gslbVirtualServersActiveServices:                    gslbVirtualServersActiveServices,
		gslbVirtualServersTotalHits:                         gslbVirtualServersTotalHits,
		gslbVirtualServersTotalRequests:                     gslbVirtualServersTotalRequests,
		gslbVirtualServersTotalResponses:                    gslbVirtualServersTotalResponses,
		gslbVirtualServersTotalRequestBytes:                 gslbVirtualServersTotalRequestBytes,
		gslbVirtualServersTotalResponseBytes:                gslbVirtualServersTotalResponseBytes,
		gslbVirtualServersCurrentClientConnections:          gslbVirtualServersCurrentClientConnections,
		gslbVirtualServersCurrentServerConnections:          gslbVirtualServersCurrentServerConnections,
		csVirtualServersState:                               csVirtualServersState,
		csVirtualServersTotalHits:                           csVirtualServersTotalHits,
		csVirtualServersTotalRequests:                       csVirtualServersTotalRequests,
		csVirtualServersTotalResponses:                      csVirtualServersTotalResponses,
		csVirtualServersTotalRequestBytes:                   csVirtualServersTotalRequestBytes,
		csVirtualServersTotalResponseBytes:                  csVirtualServersTotalResponseBytes,
		csVirtualServersCurrentClientConnections:            csVirtualServersCurrentClientConnections,
		csVirtualServersCurrentServerConnections:            csVirtualServersCurrentServerConnections,
		csVirtualServersEstablishedConnections:              csVirtualServersEstablishedConnections,
		csVirtualServersTotalPacketsReceived:                csVirtualServersTotalPacketsReceived,
		csVirtualServersTotalPacketsSent:                    csVirtualServersTotalPacketsSent,
		csVirtualServersTotalSpillovers:                     csVirtualServersTotalSpillovers,
		csVirtualServersDeferredRequests:                    csVirtualServersDeferredRequests,
		csVirtualServersNumberInvalidRequestResponse:        csVirtualServersNumberInvalidRequestResponse,
		csVirtualServersNumberInvalidRequestResponseDropped: csVirtualServersNumberInvalidRequestResponseDropped,
		csVirtualServersTotalVServerDownBackupHits:          csVirtualServersTotalVServerDownBackupHits,
		csVirtualServersCurrentMultipathSessions:            csVirtualServersCurrentMultipathSessions,
		csVirtualServersCurrentMultipathSubflows:            csVirtualServersCurrentMultipathSubflows,
	}, nil
}

// Describe implements Collector
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- modelID
	ch <- mgmtCPUUsage
	ch <- memUsage
	ch <- pktCPUUsage
	ch <- flashPartitionUsage
	ch <- varPartitionUsage
	ch <- totRxMB
	ch <- totTxMB
	ch <- httpResponses
	ch <- httpRequests
	ch <- tcpCurrentClientConnections
	ch <- tcpCurrentClientConnectionsEstablished
	ch <- tcpCurrentServerConnections
	ch <- tcpCurrentServerConnectionsEstablished

	e.interfacesRxBytes.Describe(ch)
	e.interfacesTxBytes.Describe(ch)
	e.interfacesRxPackets.Describe(ch)
	e.interfacesTxPackets.Describe(ch)
	e.interfacesJumboPacketsRx.Describe(ch)
	e.interfacesJumboPacketsTx.Describe(ch)
	e.interfacesErrorPacketsRx.Describe(ch)

	e.virtualServersWaitingRequests.Describe(ch)
	e.virtualServersHealth.Describe(ch)
	e.virtualServersInactiveServices.Describe(ch)
	e.virtualServersActiveServices.Describe(ch)
	e.virtualServersTotalHits.Describe(ch)
	e.virtualServersTotalRequests.Describe(ch)
	e.virtualServersTotalResponses.Describe(ch)
	e.virtualServersTotalRequestBytes.Describe(ch)
	e.virtualServersTotalResponseBytes.Describe(ch)
	e.virtualServersCurrentClientConnections.Describe(ch)
	e.virtualServersCurrentServerConnections.Describe(ch)

	e.servicesThroughput.Describe(ch)
	e.servicesAvgTTFB.Describe(ch)
	e.servicesState.Describe(ch)
	e.servicesTotalRequests.Describe(ch)
	e.servicesTotalResponses.Describe(ch)
	e.servicesTotalRequestBytes.Describe(ch)
	e.servicesTotalResponseBytes.Describe(ch)
	e.servicesCurrentClientConns.Describe(ch)
	e.servicesSurgeCount.Describe(ch)
	e.servicesCurrentServerConns.Describe(ch)
	e.servicesServerEstablishedConnections.Describe(ch)
	e.servicesCurrentReusePool.Describe(ch)
	e.servicesMaxClients.Describe(ch)
	e.servicesCurrentLoad.Describe(ch)
	e.servicesVirtualServerServiceHits.Describe(ch)
	e.servicesActiveTransactions.Describe(ch)

	e.serviceGroupsState.Describe(ch)
	e.serviceGroupsAvgTTFB.Describe(ch)
	e.serviceGroupsTotalRequests.Describe(ch)
	e.serviceGroupsTotalResponses.Describe(ch)
	e.serviceGroupsTotalRequestBytes.Describe(ch)
	e.serviceGroupsTotalResponseBytes.Describe(ch)
	e.serviceGroupsCurrentClientConnections.Describe(ch)
	e.serviceGroupsSurgeCount.Describe(ch)
	e.serviceGroupsCurrentServerConnections.Describe(ch)
	e.serviceGroupsServerEstablishedConnections.Describe(ch)
	e.serviceGroupsCurrentReusePool.Describe(ch)
	e.serviceGroupsMaxClients.Describe(ch)

	e.gslbServicesState.Describe(ch)
	e.gslbServicesTotalRequests.Describe(ch)
	e.gslbServicesTotalResponses.Describe(ch)
	e.gslbServicesTotalRequestBytes.Describe(ch)
	e.gslbServicesTotalResponseBytes.Describe(ch)
	e.gslbServicesCurrentClientConns.Describe(ch)
	e.gslbServicesCurrentServerConns.Describe(ch)
	e.gslbServicesCurrentLoad.Describe(ch)
	e.gslbServicesVirtualServerServiceHits.Describe(ch)
	e.gslbServicesEstablishedConnections.Describe(ch)

	e.gslbVirtualServersHealth.Describe(ch)
	e.gslbVirtualServersInactiveServices.Describe(ch)
	e.gslbVirtualServersActiveServices.Describe(ch)
	e.gslbVirtualServersTotalHits.Describe(ch)
	e.gslbVirtualServersTotalRequests.Describe(ch)
	e.gslbVirtualServersTotalResponses.Describe(ch)
	e.gslbVirtualServersTotalRequestBytes.Describe(ch)
	e.gslbVirtualServersTotalResponseBytes.Describe(ch)
	e.gslbVirtualServersCurrentClientConnections.Describe(ch)
	e.gslbVirtualServersCurrentServerConnections.Describe(ch)

	e.csVirtualServersState.Describe(ch)
	e.csVirtualServersTotalHits.Describe(ch)
	e.csVirtualServersTotalRequests.Describe(ch)
	e.csVirtualServersTotalResponses.Describe(ch)
	e.csVirtualServersTotalRequestBytes.Describe(ch)
	e.csVirtualServersTotalResponseBytes.Describe(ch)
	e.csVirtualServersCurrentClientConnections.Describe(ch)
	e.csVirtualServersCurrentServerConnections.Describe(ch)
	e.csVirtualServersEstablishedConnections.Describe(ch)
	e.csVirtualServersTotalPacketsReceived.Describe(ch)
	e.csVirtualServersTotalPacketsSent.Describe(ch)
	e.csVirtualServersTotalSpillovers.Describe(ch)
	e.csVirtualServersDeferredRequests.Describe(ch)
	e.csVirtualServersNumberInvalidRequestResponse.Describe(ch)
	e.csVirtualServersNumberInvalidRequestResponseDropped.Describe(ch)
	e.csVirtualServersTotalVServerDownBackupHits.Describe(ch)
	e.csVirtualServersCurrentMultipathSessions.Describe(ch)
	e.csVirtualServersCurrentMultipathSubflows.Describe(ch)
}

func (e *Exporter) collectInterfacesRxBytes(ns netscaler.NSAPIResponse) {
	e.interfacesRxBytes.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalReceivedBytes, 64)
		e.interfacesRxBytes.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesTxBytes(ns netscaler.NSAPIResponse) {
	e.interfacesTxBytes.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalTransmitBytes, 64)
		e.interfacesTxBytes.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesRxPackets(ns netscaler.NSAPIResponse) {
	e.interfacesRxPackets.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalReceivedPackets, 64)
		e.interfacesRxPackets.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesTxPackets(ns netscaler.NSAPIResponse) {
	e.interfacesTxPackets.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalTransmitPackets, 64)
		e.interfacesTxPackets.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsRx(ns netscaler.NSAPIResponse) {
	e.interfacesJumboPacketsRx.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.JumboPacketsReceived, 64)
		e.interfacesJumboPacketsRx.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsTx(ns netscaler.NSAPIResponse) {
	e.interfacesJumboPacketsTx.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.JumboPacketsTransmitted, 64)
		e.interfacesJumboPacketsTx.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesErrorPacketsRx(ns netscaler.NSAPIResponse) {
	e.interfacesErrorPacketsRx.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.ErrorPacketsReceived, 64)
		e.interfacesErrorPacketsRx.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectVirtualServerWaitingRequests(ns netscaler.NSAPIResponse) {
	e.virtualServersWaitingRequests.Reset()

	for _, vs := range ns.VirtualServerStats {
		waitingRequests, _ := strconv.ParseFloat(vs.WaitingRequests, 64)
		e.virtualServersWaitingRequests.WithLabelValues(nsInstance, vs.Name).Set(waitingRequests)
	}
}

func (e *Exporter) collectVirtualServerHealth(ns netscaler.NSAPIResponse) {
	e.virtualServersHealth.Reset()

	for _, vs := range ns.VirtualServerStats {
		health, _ := strconv.ParseFloat(vs.Health, 64)
		e.virtualServersHealth.WithLabelValues(nsInstance, vs.Name).Set(health)
	}
}

func (e *Exporter) collectVirtualServerInactiveServices(ns netscaler.NSAPIResponse) {
	e.virtualServersInactiveServices.Reset()

	for _, vs := range ns.VirtualServerStats {
		inactiveServices, _ := strconv.ParseFloat(vs.InactiveServices, 64)
		e.virtualServersInactiveServices.WithLabelValues(nsInstance, vs.Name).Set(inactiveServices)
	}
}

func (e *Exporter) collectVirtualServerActiveServices(ns netscaler.NSAPIResponse) {
	e.virtualServersActiveServices.Reset()

	for _, vs := range ns.VirtualServerStats {
		activeServices, _ := strconv.ParseFloat(vs.ActiveServices, 64)
		e.virtualServersActiveServices.WithLabelValues(nsInstance, vs.Name).Set(activeServices)
	}
}

func (e *Exporter) collectVirtualServerTotalHits(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalHits.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalHits, _ := strconv.ParseFloat(vs.TotalHits, 64)
		e.virtualServersTotalHits.WithLabelValues(nsInstance, vs.Name).Set(totalHits)
	}
}

func (e *Exporter) collectVirtualServerTotalRequests(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalRequests.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalRequests, _ := strconv.ParseFloat(vs.TotalRequests, 64)
		e.virtualServersTotalRequests.WithLabelValues(nsInstance, vs.Name).Set(totalRequests)
	}
}

func (e *Exporter) collectVirtualServerTotalResponses(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalResponses.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalResponses, _ := strconv.ParseFloat(vs.TotalResponses, 64)
		e.virtualServersTotalResponses.WithLabelValues(nsInstance, vs.Name).Set(totalResponses)
	}
}

func (e *Exporter) collectVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalRequestBytes.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalRequestBytes, _ := strconv.ParseFloat(vs.TotalRequestBytes, 64)
		e.virtualServersTotalRequestBytes.WithLabelValues(nsInstance, vs.Name).Set(totalRequestBytes)
	}
}

func (e *Exporter) collectVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalResponseBytes.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalResponseBytes, _ := strconv.ParseFloat(vs.TotalResponseBytes, 64)
		e.virtualServersTotalResponseBytes.WithLabelValues(nsInstance, vs.Name).Set(totalResponseBytes)
	}
}

func (e *Exporter) collectVirtualServerCurrentClientConnections(ns netscaler.NSAPIResponse) {
	e.virtualServersCurrentClientConnections.Reset()

	for _, vs := range ns.VirtualServerStats {
		currentClientConnections, _ := strconv.ParseFloat(vs.CurrentClientConnections, 64)
		e.virtualServersCurrentClientConnections.WithLabelValues(nsInstance, vs.Name).Set(currentClientConnections)
	}
}

func (e *Exporter) collectVirtualServerCurrentServerConnections(ns netscaler.NSAPIResponse) {
	e.virtualServersCurrentServerConnections.Reset()

	for _, vs := range ns.VirtualServerStats {
		currentServerConnections, _ := strconv.ParseFloat(vs.CurrentServerConnections, 64)
		e.virtualServersCurrentServerConnections.WithLabelValues(nsInstance, vs.Name).Set(currentServerConnections)
	}
}

func (e *Exporter) collectServicesThroughput(ns netscaler.NSAPIResponse) {
	e.servicesThroughput.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.Throughput, 64)
		e.servicesThroughput.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesAvgTTFB(ns netscaler.NSAPIResponse) {
	e.servicesAvgTTFB.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.AvgTimeToFirstByte, 64)
		e.servicesAvgTTFB.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesState(ns netscaler.NSAPIResponse) {
	e.servicesState.Reset()

	for _, service := range ns.ServiceStats {
		state := 0.0

		if service.State == "UP" {
			state = 1.0
		}

		e.servicesState.WithLabelValues(nsInstance, service.Name).Set(state)
	}
}

func (e *Exporter) collectServicesTotalRequests(ns netscaler.NSAPIResponse) {
	e.servicesTotalRequests.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequests, 64)
		e.servicesTotalRequests.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesTotalResponses(ns netscaler.NSAPIResponse) {
	e.servicesTotalResponses.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponses, 64)
		e.servicesTotalResponses.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.servicesTotalRequestBytes.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequestBytes, 64)
		e.servicesTotalRequestBytes.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.servicesTotalResponseBytes.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponseBytes, 64)
		e.servicesTotalResponseBytes.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentClientConns(ns netscaler.NSAPIResponse) {
	e.servicesCurrentClientConns.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentClientConnections, 64)
		e.servicesCurrentClientConns.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesSurgeCount(ns netscaler.NSAPIResponse) {
	e.servicesSurgeCount.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.SurgeCount, 64)
		e.servicesSurgeCount.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentServerConns(ns netscaler.NSAPIResponse) {
	e.servicesCurrentServerConns.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentServerConnections, 64)
		e.servicesCurrentServerConns.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesServerEstablishedConnections(ns netscaler.NSAPIResponse) {
	e.servicesServerEstablishedConnections.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ServerEstablishedConnections, 64)
		e.servicesServerEstablishedConnections.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentReusePool(ns netscaler.NSAPIResponse) {
	e.servicesCurrentReusePool.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentReusePool, 64)
		e.servicesCurrentReusePool.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesMaxClients(ns netscaler.NSAPIResponse) {
	e.servicesMaxClients.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.MaxClients, 64)
		e.servicesMaxClients.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentLoad(ns netscaler.NSAPIResponse) {
	e.servicesCurrentLoad.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentLoad, 64)
		e.servicesCurrentLoad.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesVirtualServerServiceHits(ns netscaler.NSAPIResponse) {
	e.servicesVirtualServerServiceHits.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ServiceHits, 64)
		e.servicesVirtualServerServiceHits.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesActiveTransactions(ns netscaler.NSAPIResponse) {
	e.servicesActiveTransactions.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ActiveTransactions, 64)
		e.servicesActiveTransactions.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsState(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsState.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		state := 0.0

		if sg.State == "UP" {
			state = 1.0
		}

		e.serviceGroupsState.WithLabelValues(nsInstance, sgName, servername).Set(state)
	}
}

func (e *Exporter) collectServiceGroupsAvgTTFB(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsAvgTTFB.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.AvgTimeToFirstByte, 64)
		e.serviceGroupsAvgTTFB.WithLabelValues(nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsTotalRequests(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsTotalRequests.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.TotalRequests, 64)
		e.serviceGroupsTotalRequests.WithLabelValues(nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsTotalResponses(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsTotalResponses.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.TotalResponses, 64)
		e.serviceGroupsTotalResponses.WithLabelValues(nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsTotalRequestBytes(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsTotalRequestBytes.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.TotalRequestBytes, 64)
		e.serviceGroupsTotalRequestBytes.WithLabelValues(nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsTotalResponseBytes(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsTotalResponseBytes.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.TotalResponseBytes, 64)
		e.serviceGroupsTotalResponseBytes.WithLabelValues(nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsCurrentClientConnections(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsCurrentClientConnections.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.CurrentClientConnections, 64)
		e.serviceGroupsCurrentClientConnections.WithLabelValues(nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsSurgeCount(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsSurgeCount.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.SurgeCount, 64)
		e.serviceGroupsSurgeCount.WithLabelValues(nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsCurrentServerConnections(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsCurrentServerConnections.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.CurrentServerConnections, 64)
		e.serviceGroupsCurrentServerConnections.WithLabelValues(nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsServerEstablishedConnections(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsServerEstablishedConnections.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.ServerEstablishedConnections, 64)
		e.serviceGroupsServerEstablishedConnections.WithLabelValues(nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsCurrentReusePool(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsCurrentReusePool.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.CurrentReusePool, 64)
		e.serviceGroupsCurrentReusePool.WithLabelValues(nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsMaxClients(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsMaxClients.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.MaxClients, 64)
		e.serviceGroupsMaxClients.WithLabelValues(nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesState(ns netscaler.NSAPIResponse) {
	e.gslbServicesState.Reset()

	for _, service := range ns.GSLBServiceStats {
		state := 0.0

		if service.State == "UP" {
			state = 1.0
		}

		e.gslbServicesState.WithLabelValues(nsInstance, service.Name).Set(state)
	}
}

func (e *Exporter) collectGSLBServicesTotalRequests(ns netscaler.NSAPIResponse) {
	e.gslbServicesTotalRequests.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequests, 64)
		e.gslbServicesTotalRequests.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesTotalResponses(ns netscaler.NSAPIResponse) {
	e.gslbServicesTotalResponses.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponses, 64)
		e.gslbServicesTotalResponses.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.gslbServicesTotalRequestBytes.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequestBytes, 64)
		e.gslbServicesTotalRequestBytes.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.gslbServicesTotalResponseBytes.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponseBytes, 64)
		e.gslbServicesTotalResponseBytes.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesCurrentClientConns(ns netscaler.NSAPIResponse) {
	e.gslbServicesCurrentClientConns.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentClientConnections, 64)
		e.gslbServicesCurrentClientConns.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesCurrentServerConns(ns netscaler.NSAPIResponse) {
	e.gslbServicesCurrentServerConns.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentServerConnections, 64)
		e.gslbServicesCurrentServerConns.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesEstablishedConnections(ns netscaler.NSAPIResponse) {
	e.gslbServicesEstablishedConnections.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.EstablishedConnections, 64)
		e.gslbServicesEstablishedConnections.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesCurrentLoad(ns netscaler.NSAPIResponse) {
	e.gslbServicesCurrentLoad.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentLoad, 64)
		e.gslbServicesCurrentLoad.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesVirtualServerServiceHits(ns netscaler.NSAPIResponse) {
	e.gslbServicesVirtualServerServiceHits.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.ServiceHits, 64)
		e.gslbServicesVirtualServerServiceHits.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBVirtualServerHealth(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersHealth.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		health, _ := strconv.ParseFloat(vs.Health, 64)
		e.gslbVirtualServersHealth.WithLabelValues(nsInstance, vs.Name).Set(health)
	}
}

func (e *Exporter) collectGSLBVirtualServerInactiveServices(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersInactiveServices.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		inactiveServices, _ := strconv.ParseFloat(vs.InactiveServices, 64)
		e.gslbVirtualServersInactiveServices.WithLabelValues(nsInstance, vs.Name).Set(inactiveServices)
	}
}

func (e *Exporter) collectGSLBVirtualServerActiveServices(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersActiveServices.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		activeServices, _ := strconv.ParseFloat(vs.ActiveServices, 64)
		e.gslbVirtualServersActiveServices.WithLabelValues(nsInstance, vs.Name).Set(activeServices)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalHits(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersTotalHits.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		totalHits, _ := strconv.ParseFloat(vs.TotalHits, 64)
		e.gslbVirtualServersTotalHits.WithLabelValues(nsInstance, vs.Name).Set(totalHits)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalRequests(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersTotalRequests.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		totalRequests, _ := strconv.ParseFloat(vs.TotalRequests, 64)
		e.gslbVirtualServersTotalRequests.WithLabelValues(nsInstance, vs.Name).Set(totalRequests)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalResponses(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersTotalResponses.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		totalResponses, _ := strconv.ParseFloat(vs.TotalResponses, 64)
		e.gslbVirtualServersTotalResponses.WithLabelValues(nsInstance, vs.Name).Set(totalResponses)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersTotalRequestBytes.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		totalRequestBytes, _ := strconv.ParseFloat(vs.TotalRequestBytes, 64)
		e.gslbVirtualServersTotalRequestBytes.WithLabelValues(nsInstance, vs.Name).Set(totalRequestBytes)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersTotalResponseBytes.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		totalResponseBytes, _ := strconv.ParseFloat(vs.TotalResponseBytes, 64)
		e.gslbVirtualServersTotalResponseBytes.WithLabelValues(nsInstance, vs.Name).Set(totalResponseBytes)
	}
}

func (e *Exporter) collectGSLBVirtualServerCurrentClientConnections(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersCurrentClientConnections.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		currentClientConnections, _ := strconv.ParseFloat(vs.CurrentClientConnections, 64)
		e.gslbVirtualServersCurrentClientConnections.WithLabelValues(nsInstance, vs.Name).Set(currentClientConnections)
	}
}

func (e *Exporter) collectGSLBVirtualServerCurrentServerConnections(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersCurrentServerConnections.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		currentServerConnections, _ := strconv.ParseFloat(vs.CurrentServerConnections, 64)
		e.gslbVirtualServersCurrentServerConnections.WithLabelValues(nsInstance, vs.Name).Set(currentServerConnections)
	}
}

func (e *Exporter) collectCSVirtualServerState(ns netscaler.NSAPIResponse) {
	e.csVirtualServersState.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		state := 0.0

		if vs.State == "UP" {
			state = 1.0
		}

		e.csVirtualServersState.WithLabelValues(nsInstance, vs.Name).Set(state)
	}
}

func (e *Exporter) collectCSVirtualServerTotalHits(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalHits.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalHits, _ := strconv.ParseFloat(vs.TotalHits, 64)
		e.csVirtualServersTotalHits.WithLabelValues(nsInstance, vs.Name).Set(totalHits)
	}
}

func (e *Exporter) collectCSVirtualServerTotalRequests(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalRequests.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalRequests, _ := strconv.ParseFloat(vs.TotalRequests, 64)
		e.csVirtualServersTotalRequests.WithLabelValues(nsInstance, vs.Name).Set(totalRequests)
	}
}

func (e *Exporter) collectCSVirtualServerTotalResponses(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalResponses.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalResponses, _ := strconv.ParseFloat(vs.TotalResponses, 64)
		e.csVirtualServersTotalResponses.WithLabelValues(nsInstance, vs.Name).Set(totalResponses)
	}
}

func (e *Exporter) collectCSVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalRequestBytes.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalRequestBytes, _ := strconv.ParseFloat(vs.TotalRequestBytes, 64)
		e.csVirtualServersTotalRequestBytes.WithLabelValues(nsInstance, vs.Name).Set(totalRequestBytes)
	}
}

func (e *Exporter) collectCSVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalResponseBytes.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalResponseBytes, _ := strconv.ParseFloat(vs.TotalResponseBytes, 64)
		e.csVirtualServersTotalResponseBytes.WithLabelValues(nsInstance, vs.Name).Set(totalResponseBytes)
	}
}

func (e *Exporter) collectCSVirtualServerCurrentClientConnections(ns netscaler.NSAPIResponse) {
	e.csVirtualServersCurrentClientConnections.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		currentClientConnections, _ := strconv.ParseFloat(vs.CurrentClientConnections, 64)
		e.csVirtualServersCurrentClientConnections.WithLabelValues(nsInstance, vs.Name).Set(currentClientConnections)
	}
}

func (e *Exporter) collectCSVirtualServerCurrentServerConnections(ns netscaler.NSAPIResponse) {
	e.csVirtualServersCurrentServerConnections.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		currentServerConnections, _ := strconv.ParseFloat(vs.CurrentServerConnections, 64)
		e.csVirtualServersCurrentServerConnections.WithLabelValues(nsInstance, vs.Name).Set(currentServerConnections)
	}
}

func (e *Exporter) collectCSVirtualServerEstablishedConnections(ns netscaler.NSAPIResponse) {
	e.csVirtualServersEstablishedConnections.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		EstablishedConnections, _ := strconv.ParseFloat(vs.EstablishedConnections, 64)
		e.csVirtualServersEstablishedConnections.WithLabelValues(nsInstance, vs.Name).Set(EstablishedConnections)
	}
}

func (e *Exporter) collectCSVirtualServerTotalPacketsReceived(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalPacketsReceived.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalPacketsReceived, _ := strconv.ParseFloat(vs.TotalPacketsReceived, 64)
		e.csVirtualServersTotalPacketsReceived.WithLabelValues(nsInstance, vs.Name).Set(totalPacketsReceived)
	}
}

func (e *Exporter) collectCSVirtualServerTotalPacketsSent(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalPacketsSent.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalPacketsSent, _ := strconv.ParseFloat(vs.TotalPacketsSent, 64)
		e.csVirtualServersTotalPacketsSent.WithLabelValues(nsInstance, vs.Name).Set(totalPacketsSent)
	}
}

func (e *Exporter) collectCSVirtualServerTotalSpillovers(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalSpillovers.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalSpillovers, _ := strconv.ParseFloat(vs.TotalSpillovers, 64)
		e.csVirtualServersTotalSpillovers.WithLabelValues(nsInstance, vs.Name).Set(totalSpillovers)
	}
}

func (e *Exporter) collectCSVirtualServerDeferredRequests(ns netscaler.NSAPIResponse) {
	e.csVirtualServersDeferredRequests.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		deferredRequests, _ := strconv.ParseFloat(vs.DeferredRequests, 64)
		e.csVirtualServersDeferredRequests.WithLabelValues(nsInstance, vs.Name).Set(deferredRequests)
	}
}

func (e *Exporter) collectCSVirtualServerNumberInvalidRequestResponse(ns netscaler.NSAPIResponse) {
	e.csVirtualServersNumberInvalidRequestResponse.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		numberInvalidRequestResponse, _ := strconv.ParseFloat(vs.InvalidRequestResponse, 64)
		e.csVirtualServersNumberInvalidRequestResponse.WithLabelValues(nsInstance, vs.Name).Set(numberInvalidRequestResponse)
	}
}

func (e *Exporter) collectCSVirtualServerNumberInvalidRequestResponseDropped(ns netscaler.NSAPIResponse) {
	e.csVirtualServersNumberInvalidRequestResponseDropped.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		numberInvalidRequestResponseDropped, _ := strconv.ParseFloat(vs.InvalidRequestResponseDropped, 64)
		e.csVirtualServersNumberInvalidRequestResponseDropped.WithLabelValues(nsInstance, vs.Name).Set(numberInvalidRequestResponseDropped)
	}
}

func (e *Exporter) collectCSVirtualServerTotalVServerDownBackupHits(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalVServerDownBackupHits.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalVServerDownBackupHits, _ := strconv.ParseFloat(vs.TotalVServerDownBackupHits, 64)
		e.csVirtualServersTotalVServerDownBackupHits.WithLabelValues(nsInstance, vs.Name).Set(totalVServerDownBackupHits)
	}
}

func (e *Exporter) collectCSVirtualServerCurrentMultipathSessions(ns netscaler.NSAPIResponse) {
	e.csVirtualServersCurrentMultipathSessions.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		currentMultipathSessions, _ := strconv.ParseFloat(vs.CurrentMultipathSessions, 64)
		e.csVirtualServersCurrentMultipathSessions.WithLabelValues(nsInstance, vs.Name).Set(currentMultipathSessions)
	}
}

func (e *Exporter) collectCSVirtualServerCurrentMultipathSubflows(ns netscaler.NSAPIResponse) {
	e.csVirtualServersCurrentMultipathSubflows.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		currentMultipathSubflows, _ := strconv.ParseFloat(vs.CurrentMultipathSubflows, 64)
		e.csVirtualServersCurrentMultipathSubflows.WithLabelValues(nsInstance, vs.Name).Set(currentMultipathSubflows)
	}
}

// Collect is initiated by the Prometheus handler and gathers the metrics
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	nsClient, err := netscaler.NewNitroClient(*url, *username, *password, *ignoreCert)
	if err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}

	err = netscaler.Connect(nsClient)
	if err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}

	nslicense, err := netscaler.GetNSLicense(nsClient, "")
	if err != nil {
		level.Error(logger).Log("msg", err)
	}

	ns, err := netscaler.GetNSStats(nsClient, "")
	if err != nil {
		level.Error(logger).Log("msg", err)
	}

	interfaces, err := netscaler.GetInterfaceStats(nsClient, "")
	if err != nil {
		level.Error(logger).Log("msg", err)
	}

	virtualServers, err := netscaler.GetVirtualServerStats(nsClient, "")
	if err != nil {
		level.Error(logger).Log("msg", err)
	}

	services, err := netscaler.GetServiceStats(nsClient, "")
	if err != nil {
		level.Error(logger).Log("msg", err)
	}

	gslbServices, err := netscaler.GetGSLBServiceStats(nsClient, "")
	if err != nil {
		level.Error(logger).Log("msg", err)
	}

	gslbVirtualServers, err := netscaler.GetGSLBVirtualServerStats(nsClient, "")
	if err != nil {
		level.Error(logger).Log("msg", err)
	}

	csVirtualServers, err := netscaler.GetCSVirtualServerStats(nsClient, "")
	if err != nil {
		level.Error(logger).Log("msg", err)
	}

	fltModelID, _ := strconv.ParseFloat(nslicense.NSLicense.ModelID, 64)

	fltTotRxMB, _ := strconv.ParseFloat(ns.NSStats.TotalReceivedMB, 64)
	fltTotTxMB, _ := strconv.ParseFloat(ns.NSStats.TotalTransmitMB, 64)
	fltHTTPRequests, _ := strconv.ParseFloat(ns.NSStats.HTTPRequests, 64)
	fltHTTPResponses, _ := strconv.ParseFloat(ns.NSStats.HTTPResponses, 64)

	fltTCPCurrentClientConnections, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentClientConnections, 64)
	fltTCPCurrentClientConnectionsEstablished, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentClientConnectionsEstablished, 64)
	fltTCPCurrentServerConnections, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentServerConnections, 64)
	fltTCPCurrentServerConnectionsEstablished, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentServerConnectionsEstablished, 64)

	ch <- prometheus.MustNewConstMetric(
		modelID, prometheus.GaugeValue, fltModelID, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		mgmtCPUUsage, prometheus.GaugeValue, ns.NSStats.MgmtCPUUsagePcnt, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		memUsage, prometheus.GaugeValue, ns.NSStats.MemUsagePcnt, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		pktCPUUsage, prometheus.GaugeValue, ns.NSStats.PktCPUUsagePcnt, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		flashPartitionUsage, prometheus.GaugeValue, ns.NSStats.FlashPartitionUsage, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		varPartitionUsage, prometheus.GaugeValue, ns.NSStats.VarPartitionUsage, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		totRxMB, prometheus.GaugeValue, fltTotRxMB, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		totTxMB, prometheus.GaugeValue, fltTotTxMB, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		httpRequests, prometheus.GaugeValue, fltHTTPRequests, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		httpResponses, prometheus.GaugeValue, fltHTTPResponses, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		tcpCurrentClientConnections, prometheus.GaugeValue, fltTCPCurrentClientConnections, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		tcpCurrentClientConnectionsEstablished, prometheus.GaugeValue, fltTCPCurrentClientConnectionsEstablished, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		tcpCurrentServerConnections, prometheus.GaugeValue, fltTCPCurrentServerConnections, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		tcpCurrentServerConnectionsEstablished, prometheus.GaugeValue, fltTCPCurrentServerConnectionsEstablished, nsInstance,
	)

	e.collectInterfacesRxBytes(interfaces)
	e.interfacesRxBytes.Collect(ch)

	e.collectInterfacesTxBytes(interfaces)
	e.interfacesTxBytes.Collect(ch)

	e.collectInterfacesRxPackets(interfaces)
	e.interfacesRxPackets.Collect(ch)

	e.collectInterfacesTxPackets(interfaces)
	e.interfacesTxPackets.Collect(ch)

	e.collectInterfacesJumboPacketsRx(interfaces)
	e.interfacesJumboPacketsRx.Collect(ch)

	e.collectInterfacesJumboPacketsTx(interfaces)
	e.interfacesJumboPacketsTx.Collect(ch)

	e.collectInterfacesErrorPacketsRx(interfaces)
	e.interfacesErrorPacketsRx.Collect(ch)

	e.collectVirtualServerWaitingRequests(virtualServers)
	e.virtualServersWaitingRequests.Collect(ch)

	e.collectVirtualServerHealth(virtualServers)
	e.virtualServersHealth.Collect(ch)

	e.collectVirtualServerInactiveServices(virtualServers)
	e.virtualServersInactiveServices.Collect(ch)

	e.collectVirtualServerActiveServices(virtualServers)
	e.virtualServersActiveServices.Collect(ch)

	e.collectVirtualServerTotalHits(virtualServers)
	e.virtualServersTotalHits.Collect(ch)

	e.collectVirtualServerTotalRequests(virtualServers)
	e.virtualServersTotalRequests.Collect(ch)

	e.collectVirtualServerTotalResponses(virtualServers)
	e.virtualServersTotalResponses.Collect(ch)

	e.collectVirtualServerTotalRequestBytes(virtualServers)
	e.virtualServersTotalRequestBytes.Collect(ch)

	e.collectVirtualServerTotalResponseBytes(virtualServers)
	e.virtualServersTotalResponseBytes.Collect(ch)

	e.collectVirtualServerCurrentClientConnections(virtualServers)
	e.virtualServersCurrentClientConnections.Collect(ch)

	e.collectVirtualServerCurrentServerConnections(virtualServers)
	e.virtualServersCurrentServerConnections.Collect(ch)

	e.collectServicesThroughput(services)
	e.servicesThroughput.Collect(ch)

	e.collectServicesAvgTTFB(services)
	e.servicesAvgTTFB.Collect(ch)

	e.collectServicesState(services)
	e.servicesState.Collect(ch)

	e.collectServicesTotalRequests(services)
	e.servicesTotalRequests.Collect(ch)

	e.collectServicesTotalResponses(services)
	e.servicesTotalResponses.Collect(ch)

	e.collectServicesTotalRequestBytes(services)
	e.servicesTotalRequestBytes.Collect(ch)

	e.collectServicesTotalResponseBytes(services)
	e.servicesTotalResponseBytes.Collect(ch)

	e.collectServicesCurrentClientConns(services)
	e.servicesCurrentClientConns.Collect(ch)

	e.collectServicesSurgeCount(services)
	e.servicesSurgeCount.Collect(ch)

	e.collectServicesCurrentServerConns(services)
	e.servicesCurrentServerConns.Collect(ch)

	e.collectServicesServerEstablishedConnections(services)
	e.servicesServerEstablishedConnections.Collect(ch)

	e.collectServicesCurrentReusePool(services)
	e.servicesCurrentReusePool.Collect(ch)

	e.collectServicesMaxClients(services)
	e.servicesMaxClients.Collect(ch)

	e.collectServicesCurrentLoad(services)
	e.servicesCurrentLoad.Collect(ch)

	e.collectServicesVirtualServerServiceHits(services)
	e.servicesVirtualServerServiceHits.Collect(ch)

	e.collectServicesActiveTransactions(services)
	e.servicesActiveTransactions.Collect(ch)

	e.collectGSLBServicesState(gslbServices)
	e.gslbServicesState.Collect(ch)

	e.collectGSLBServicesTotalRequests(gslbServices)
	e.gslbServicesTotalRequests.Collect(ch)

	e.collectGSLBServicesTotalResponses(gslbServices)
	e.gslbServicesTotalResponses.Collect(ch)

	e.collectGSLBServicesTotalRequestBytes(gslbServices)
	e.gslbServicesTotalRequestBytes.Collect(ch)

	e.collectGSLBServicesTotalResponseBytes(gslbServices)
	e.gslbServicesTotalResponseBytes.Collect(ch)

	e.collectGSLBServicesCurrentClientConns(gslbServices)
	e.gslbServicesCurrentClientConns.Collect(ch)

	e.collectGSLBServicesCurrentServerConns(gslbServices)
	e.gslbServicesCurrentServerConns.Collect(ch)

	e.collectGSLBServicesEstablishedConnections(gslbServices)
	e.gslbServicesEstablishedConnections.Collect(ch)

	e.collectGSLBServicesCurrentLoad(gslbServices)
	e.gslbServicesCurrentLoad.Collect(ch)

	e.collectGSLBServicesVirtualServerServiceHits(gslbServices)
	e.gslbServicesVirtualServerServiceHits.Collect(ch)

	e.collectGSLBVirtualServerHealth(gslbVirtualServers)
	e.gslbVirtualServersHealth.Collect(ch)

	e.collectGSLBVirtualServerInactiveServices(gslbVirtualServers)
	e.gslbVirtualServersInactiveServices.Collect(ch)

	e.collectGSLBVirtualServerActiveServices(gslbVirtualServers)
	e.gslbVirtualServersActiveServices.Collect(ch)

	e.collectGSLBVirtualServerTotalHits(gslbVirtualServers)
	e.gslbVirtualServersTotalHits.Collect(ch)

	e.collectGSLBVirtualServerTotalRequests(gslbVirtualServers)
	e.gslbVirtualServersTotalRequests.Collect(ch)

	e.collectGSLBVirtualServerTotalResponses(gslbVirtualServers)
	e.gslbVirtualServersTotalResponses.Collect(ch)

	e.collectGSLBVirtualServerTotalRequestBytes(gslbVirtualServers)
	e.gslbVirtualServersTotalRequestBytes.Collect(ch)

	e.collectGSLBVirtualServerTotalResponseBytes(gslbVirtualServers)
	e.gslbVirtualServersTotalResponseBytes.Collect(ch)

	e.collectGSLBVirtualServerCurrentClientConnections(gslbVirtualServers)
	e.gslbVirtualServersCurrentClientConnections.Collect(ch)

	e.collectGSLBVirtualServerCurrentServerConnections(gslbVirtualServers)
	e.gslbVirtualServersCurrentServerConnections.Collect(ch)

	e.collectCSVirtualServerState(csVirtualServers)
	e.csVirtualServersState.Collect(ch)

	e.collectCSVirtualServerTotalHits(csVirtualServers)
	e.csVirtualServersTotalHits.Collect(ch)

	e.collectCSVirtualServerTotalRequests(csVirtualServers)
	e.csVirtualServersTotalRequests.Collect(ch)

	e.collectCSVirtualServerTotalResponses(csVirtualServers)
	e.csVirtualServersTotalResponses.Collect(ch)

	e.collectCSVirtualServerTotalRequestBytes(csVirtualServers)
	e.csVirtualServersTotalRequestBytes.Collect(ch)

	e.collectCSVirtualServerTotalResponseBytes(csVirtualServers)
	e.csVirtualServersTotalResponseBytes.Collect(ch)

	e.collectCSVirtualServerCurrentClientConnections(csVirtualServers)
	e.csVirtualServersCurrentClientConnections.Collect(ch)

	e.collectCSVirtualServerCurrentServerConnections(csVirtualServers)
	e.csVirtualServersCurrentServerConnections.Collect(ch)

	e.collectCSVirtualServerEstablishedConnections(csVirtualServers)
	e.csVirtualServersEstablishedConnections.Collect(ch)

	e.collectCSVirtualServerTotalPacketsReceived(csVirtualServers)
	e.csVirtualServersTotalPacketsReceived.Collect(ch)

	e.collectCSVirtualServerTotalPacketsSent(csVirtualServers)
	e.csVirtualServersTotalPacketsSent.Collect(ch)

	e.collectCSVirtualServerTotalSpillovers(csVirtualServers)
	e.csVirtualServersTotalSpillovers.Collect(ch)

	e.collectCSVirtualServerDeferredRequests(csVirtualServers)
	e.csVirtualServersDeferredRequests.Collect(ch)

	e.collectCSVirtualServerNumberInvalidRequestResponse(csVirtualServers)
	e.csVirtualServersNumberInvalidRequestResponse.Collect(ch)

	e.collectCSVirtualServerNumberInvalidRequestResponseDropped(csVirtualServers)
	e.csVirtualServersNumberInvalidRequestResponseDropped.Collect(ch)

	e.collectCSVirtualServerTotalVServerDownBackupHits(csVirtualServers)
	e.csVirtualServersTotalVServerDownBackupHits.Collect(ch)

	e.collectCSVirtualServerCurrentMultipathSessions(csVirtualServers)
	e.csVirtualServersCurrentMultipathSessions.Collect(ch)

	e.collectCSVirtualServerCurrentMultipathSubflows(csVirtualServers)
	e.csVirtualServersCurrentMultipathSubflows.Collect(ch)

	servicegroups, err := netscaler.GetServiceGroups(nsClient, "attrs=servicegroupname")
	if err != nil {
		level.Error(logger).Log("msg", err)
	}

	for _, sg := range servicegroups.ServiceGroups {
		bindings, err2 := netscaler.GetServiceGroupMemberBindings(nsClient, sg.Name)
		if err2 != nil {
			level.Error(logger).Log("msg", err2)
		}

		for _, member := range bindings.ServiceGroupMemberBindings {
			// NetScaler API has a bug which means it throws errors if you try to retrieve stats for a wildcard port (* in GUI, 65535 in API and CLI).
			// Until Citrix resolve the issue we skip attempting to retrieve stats for those service groups.
			if member.Port != 65535 {
				port := strconv.FormatInt(member.Port, 10)

				qs := "args=servicegroupname:" + sg.Name + ",servername:" + member.ServerName + ",port:" + port
				stats, err2 := netscaler.GetServiceGroupMemberStats(nsClient, qs)
				if err2 != nil {
					level.Error(logger).Log("msg", err2)
				}

				e.collectServiceGroupsState(stats, sg.Name, member.ServerName)
				e.serviceGroupsState.Collect(ch)

				e.collectServiceGroupsAvgTTFB(stats, sg.Name, member.ServerName)
				e.serviceGroupsAvgTTFB.Collect(ch)

				e.collectServiceGroupsTotalRequests(stats, sg.Name, member.ServerName)
				e.serviceGroupsTotalRequests.Collect(ch)

				e.collectServiceGroupsTotalResponses(stats, sg.Name, member.ServerName)
				e.serviceGroupsTotalResponses.Collect(ch)

				e.collectServiceGroupsTotalRequestBytes(stats, sg.Name, member.ServerName)
				e.serviceGroupsTotalRequestBytes.Collect(ch)

				e.collectServiceGroupsTotalResponseBytes(stats, sg.Name, member.ServerName)
				e.serviceGroupsTotalResponseBytes.Collect(ch)

				e.collectServiceGroupsCurrentClientConnections(stats, sg.Name, member.ServerName)
				e.serviceGroupsCurrentClientConnections.Collect(ch)

				e.collectServiceGroupsSurgeCount(stats, sg.Name, member.ServerName)
				e.serviceGroupsSurgeCount.Collect(ch)

				e.collectServiceGroupsCurrentServerConnections(stats, sg.Name, member.ServerName)
				e.serviceGroupsCurrentServerConnections.Collect(ch)

				e.collectServiceGroupsServerEstablishedConnections(stats, sg.Name, member.ServerName)
				e.serviceGroupsServerEstablishedConnections.Collect(ch)

				e.collectServiceGroupsCurrentReusePool(stats, sg.Name, member.ServerName)
				e.serviceGroupsCurrentReusePool.Collect(ch)

				e.collectServiceGroupsMaxClients(stats, sg.Name, member.ServerName)
				e.serviceGroupsMaxClients.Collect(ch)
			}
		}
	}

	err = netscaler.Disconnect(nsClient)
	if err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	if *versionFlg {
		fmt.Println(app + " v" + version + " build " + build)
		os.Exit(0)
	}

	if *url == "" || *username == "" || *password == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	nsInstance = strings.TrimLeft(*url, "https://")
	nsInstance = strings.Trim(nsInstance, " /")

	logger = log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller, "app", app, "bind_port", *bindPort, "url", *url, "version", "v"+version, "build", build)

	exporter, _ := NewExporter()
	prometheus.MustRegister(exporter)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Citrix NetScaler Exporter</title></head>
			<body>
			<h1>Citrix NetScaler Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})

	http.Handle("/metrics", promhttp.Handler())

	listeningPort := ":" + strconv.Itoa(*bindPort)
	level.Info(logger).Log("msg", "Listening on port "+listeningPort)

	err := http.ListenAndServe(listeningPort, nil)
	if err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}
}
