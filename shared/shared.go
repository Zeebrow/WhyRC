package shared

type Message struct {
	Code    int    `json:"code"`
	From    string `json:"from"`
	Message string `json:"message"`
}
