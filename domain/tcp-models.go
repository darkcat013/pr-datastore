package domain

type TcpAction struct {
	Action string `json:"action"`
	Id     string `json:"id"`
	Value  []byte `json:"value"`
}
