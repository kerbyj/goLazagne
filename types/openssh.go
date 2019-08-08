package types

type OpensshData struct {
	Hosts []string `json:"hosts"`
	Keys  []string `json:"keys"`
}
