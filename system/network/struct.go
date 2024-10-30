package network

type NetWorkIndicators struct {
	Rxpck float64
	Txpck float64
	Rxkb  float64
	Txkb  float64
}

type NetWorkCard struct {
	Name         string
	Time         []string
	NetworkTotal []NetWorkIndicators
}
type NetworkInformation struct {
	NetCard []NetWorkCard
}
