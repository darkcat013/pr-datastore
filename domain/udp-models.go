package domain

type UdpAction struct {
	Action   string          `json:"action"`
	Name     string          `json:"name"`
	Port     string          `json:"port"`
	IsLeader bool            `json:"isLeader"`
	DataIds  map[string]bool `json:"dataIds"`
}
