package handlers

import (
	"fmt"
	"net/http"

	"cloudevents-explorer/internal/templates"
)

func HandlePubSub(w http.ResponseWriter, r *http.Request) {
	html := templates.GetBaseHTML("PubSub CloudEvents", templates.PubSubContent, templates.PubSubJS)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}
