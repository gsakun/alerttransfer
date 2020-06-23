package datatype

type NetData struct {
	ServerIp    string
	ServerPort  int64
	ClientPort  int64
	ClientIp    string
	SourceBytes int64
	DestBytes   int64
	Time        string
	ProtoType   string
	Stat        bool
}

type FlowInfo struct {
	Final     bool   `json: "final"`
	Timestamp string `json: "@timestamp"`
	Type      string `json:"type"`
	Proc      string `json:"proc"`
	Dest      struct {
		Ip    string `json:"ip"`
		Port  int    `json:"port"`
		Stats *Stat  `json: "stats"`
	} `json: "dest"`
	Source struct {
		Ip    string `json:"ip"`
		Port  int    `json:"port"`
		Stats *Stat  `json: "stats"`
	} `json: "source"`
	Direction string `json: "direction"`
}

type DataInfo struct {
	Type         string  `json:"type"`
	Direction    string  `json: "direction"`
	Timestamp    string  `json: "@timestamp"`
	Ip           string  `json:"ip"`
	Port         int     `json:"port"`
	Client_ip    string  `json: "client_ip"`
	Client_port  string  `json: "client_port"`
	Bytes_in     float64 `json: "bytes_in"`
	Byte_out     float64 `json: "bytes_out"`
	Responsetime int     `json: "reponsetime"`
	Proc         string  `json: "proc"`
	Status       string  `json: "status"`
}
type Stat struct {
	Byte_Total float64 `json: "net_bytes_total"`
}
