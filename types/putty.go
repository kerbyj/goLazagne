package types

type PuttyData struct {
	HostName string `json:"hostname"`
	UserName string `json:"username"`
	Key      string `json:"key"`
}
