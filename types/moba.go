package types

type MobaData struct {
	HostName    string `json:"hostname"`
	User        string `json:"user"`
	KeyLocation string `json:"keylocation"`
	Port        string `json:"port"`
	Key         string `json:"key"`
}
