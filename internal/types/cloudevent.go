package types

type CloudEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Subject   string                 `json:"subject"`
	Source    string                 `json:"source"`
	Schema    string                 `json:"schema"`
	Published string                 `json:"published"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp int64                  `json:"timestamp"`
	RawData   string                 `json:"rawData,omitempty"`
}
