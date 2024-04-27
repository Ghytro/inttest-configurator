package mockservice

type CommonServiceResult struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Port int    `json:"port"`
}
