package ws

type Response[ResponseT any] struct {
	CommandMeta
	OK       bool      `json:"ok"`
	Error    string    `json:"error,omitempty"`
	Response ResponseT `json:"response,omitempty"`
}
