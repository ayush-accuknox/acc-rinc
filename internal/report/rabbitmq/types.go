package rabbitmq

const (
	healthCheckEndpoint = "/api/health/checks/virtual-hosts"
	overviewEndpoint    = "/api/overview"
	nodesEndpoint       = "/api/nodes"
	queuesEndpoint      = "/api/queues"
	exchangesEndpoint   = "/api/exchanges"
	consumersEndpoint   = "/api/consumers"
)

// Metrics contains a set of important RabbitMQ metrics that need to be
// included in the report.
type Metrics struct {
	Overview  Overview
	Nodes     Nodes
	Queues    Queues
	Consumers Consumers
	Exchanges Exchanges
}

type Overview struct {
	Version      string       `json:"rabbitmq_version"`
	QueueTotals  queueTotals  `json:"queue_totals"`
	ObjectTotals objectTotals `json:"object_totals"`
}

type Nodes []Node

type Node struct {
	Name           string   `json:"name"`
	Running        bool     `json:"running"`
	CPUCount       uint     `json:"processors"`
	MemUsed        float64  `json:"mem_used"`
	FreeDisk       float64  `json:"disk_free"`
	ProcUsed       uint     `json:"proc_used"`
	SocketsUsed    uint     `json:"sockets_used"`
	FDUsed         uint     `json:"fd_used"`
	Uptime         uint64   `json:"uptime"`
	EnabledPlugins []string `json:"enabled_plugins"`
}

type Queues []Queue

type Queue struct {
	Durable                bool   `json:"durable"`
	Messages               uint   `json:"messages"`
	UnacknowledgedMessages uint   `json:"messages_unacknowledged"`
	ReadyMessages          uint   `json:"messages_ready"`
	Name                   string `json:"name"`
	State                  string `json:"state"`
}

type Consumers []Consumer

type Consumer struct {
	Active        bool              `json:"active"`
	Tag           string            `json:"consumer_tag"`
	PrefetchCount uint              `json:"prefetch_count"`
	Queue         consumerQueueInfo `json:"queue"`
}

type Exchanges []Exchange

type Exchange struct {
	Durable bool   `json:"durable"`
	Name    string `json:"name"`
	Typ     string `json:"type"`
}

type queueTotals struct {
	Messages               uint `json:"messages"`
	ReadyMessages          uint `json:"messages_ready"`
	UnacknowledgedMessages uint `json:"messages_unacknowledged"`
}

type objectTotals struct {
	Channels    uint `json:"channels"`
	Connections uint `json:"connections"`
	Consumers   uint `json:"consumers"`
	Exchanges   uint `json:"exchanges"`
	Queues      uint `json:"queues"`
}

type consumerQueueInfo struct {
	Name string `json:"name"`
}
