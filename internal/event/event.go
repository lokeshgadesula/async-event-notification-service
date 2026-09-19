package event

type Event struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
	Attempts  int    `json:"attempts"`
}
