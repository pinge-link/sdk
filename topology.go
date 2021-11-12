package pinge

type TopologyConfig struct {
	Regions []TopologyRegion
}

type TopologyRegion struct {
	Id       string         `json:"id"`
	PingHost string         `json:"ping_host"`
	Gates    []TopologyGate `json:"gates"`
}

type TopologyGate struct {
	SecondaryAddress string `json:"secondary_address"`
	PrimaryAddress   string `json:"primary_address"`
	Busy             bool   `json:"-"`
}
