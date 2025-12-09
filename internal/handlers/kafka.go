package handlers

import (
	"fmt"
	"net/http"

	"cloudevents-explorer/internal/templates"
)

func HandleKafka(w http.ResponseWriter, r *http.Request) {
	html := templates.GetBaseHTML("Kafka EventMesh", templates.KafkaContent, templates.KafkaJS)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}
