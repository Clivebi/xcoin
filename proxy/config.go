package proxy

type AppConfig struct {
	ChannelID      string `json:"channnel"`
	OrgName        string `json:"orgname"`
	OrgAdmin       string `json:"orgadmin"`
	OrdererOrgName string `json:"orderorgname"`
	ChainCode      string `json:"chaincode"`
	ConfigFile     string `json:"config"`
}
