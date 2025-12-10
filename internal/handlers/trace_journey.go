package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"sort"
	"strings"
	"time"

	"cloudevents-explorer/internal/templates"
)

type TraceSearchRequest struct {
	TraceID    string   `json:"traceId"`
	Containers []string `json:"containers"`
}

type TraceEvent struct {
	Container string `json:"container"`
	Timestamp string `json:"timestamp"`
	Body      string `json:"body"`
	Name      string `json:"name"`
	SpanID    string `json:"span_id"`
	TraceID   string `json:"trace_id"`
	Severity  string `json:"severity"`
}

type TraceSearchResponse struct {
	Events []TraceEvent `json:"events"`
}

func HandleTraceJourney(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, templates.TraceJourney)
}

func HandleTraceSearch(w http.ResponseWriter, r *http.Request) {
	var req TraceSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	if req.TraceID == "" {
		http.Error(w, "trace_id is required", http.StatusBadRequest)
		return
	}

	if len(req.Containers) == 0 {
		http.Error(w, "at least one container is required", http.StatusBadRequest)
		return
	}

	// Search logs in parallel
	events := searchLogsInContainers(req.TraceID, req.Containers)

	// Sort by timestamp
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp < events[j].Timestamp
	})

	response := TraceSearchResponse{Events: events}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func searchLogsInContainers(traceID string, containers []string) []TraceEvent {
	eventsChan := make(chan []TraceEvent, len(containers))

	// Search each container in parallel
	for _, container := range containers {
		go func(cont string) {
			eventsChan <- searchContainerLogs(cont, traceID)
		}(container)
	}

	// Collect results
	var allEvents []TraceEvent
	for i := 0; i < len(containers); i++ {
		events := <-eventsChan
		allEvents = append(allEvents, events...)
	}

	return allEvents
}

func searchContainerLogs(container, traceID string) []TraceEvent {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "logs", "--timestamps", container)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}

	lines := strings.Split(string(output), "\n")
	var events []TraceEvent

	for _, line := range lines {
		if !strings.Contains(line, traceID) {
			continue
		}

		// Try to parse as JSON first
		var logJSON map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logJSON); err == nil {
			// JSON log format
			timestamp, hasTimestamp := logJSON["timestamp"].(string)
			logTraceID, hasTraceID := logJSON["trace_id"].(string)

			if !hasTimestamp || !hasTraceID {
				continue
			}

			event := TraceEvent{
				Container: container,
				Timestamp: timestamp,
				TraceID:   logTraceID,
			}

			// Extract optional fields
			if body, ok := logJSON["body"].(string); ok {
				event.Body = body
			}

			if name, ok := logJSON["name"].(string); ok {
				event.Name = name
			}

			if spanID, ok := logJSON["span_id"].(string); ok {
				event.SpanID = spanID
			}

			if severity, ok := logJSON["severity"].(string); ok {
				event.Severity = severity
			} else if severity, ok := logJSON["severity_text"].(string); ok {
				event.Severity = severity
			}

			events = append(events, event)
		} else {
			// Plain text log format - any log line containing the trace ID
			// Since trace ID is a GUID, it's unlikely to match unrelated logs

			// Extract timestamp from docker logs --timestamps prefix
			parts := strings.SplitN(line, " ", 2)
			if len(parts) < 2 {
				continue
			}

			dockerTimestamp := parts[0]
			logContent := parts[1]

			event := TraceEvent{
				Container: container,
				Timestamp: dockerTimestamp,
				TraceID:   traceID,
				Body:      logContent,
				Name:      "Log Entry",
			}

			// Try to extract more context from the log line
			if strings.Contains(strings.ToLower(logContent), "error") ||
			   strings.Contains(strings.ToLower(logContent), "failed") {
				event.Severity = "error"
			}

			events = append(events, event)
		}
	}

	return events
}
