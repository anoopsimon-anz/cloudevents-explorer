package handlers

import (
	"fmt"
	"net/http"

	"cloudevents-explorer/internal/templates"
)

func HandleFlowDiagram(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, templates.FlowDiagram)
}
