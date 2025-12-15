package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"cloudevents-explorer/internal/config"
	"cloudevents-explorer/internal/spanner"
	"cloudevents-explorer/internal/templates"
	"cloudevents-explorer/internal/types"
)

// HandleSpanner renders the main Spanner explorer page
func HandleSpanner(w http.ResponseWriter, r *http.Request) {
	html := templates.GetSpannerHTML()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// HandleSpannerConnect tests connection to Spanner
func HandleSpannerConnect(w http.ResponseWriter, r *http.Request) {
	var req types.ConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Try environment variables if not provided
	if req.EmulatorHost == "" {
		req.EmulatorHost = os.Getenv("SPANNER_EMULATOR_HOST")
	}
	if req.ProjectID == "" {
		req.ProjectID = os.Getenv("SPANNER_PROJECT")
	}
	if req.InstanceID == "" {
		req.InstanceID = os.Getenv("SPANNER_INSTANCE")
	}
	if req.DatabaseID == "" {
		req.DatabaseID = os.Getenv("SPANNER_DATABASE")
	}

	resp := spanner.TestConnection(req)

	w.Header().Set("Content-Type", "application/json")
	if !resp.Success {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(resp)
}

// HandleSpannerTables returns list of tables
func HandleSpannerTables(w http.ResponseWriter, r *http.Request) {
	var req types.ConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tables, err := spanner.ListTables(req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tables)
}

// HandleSpannerQuery executes a SQL query
func HandleSpannerQuery(w http.ResponseWriter, r *http.Request) {
	var req types.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := spanner.ExecuteQuery(req)

	w.Header().Set("Content-Type", "application/json")
	if resp.Error != "" {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(resp)
}

// HandleSaveSpannerConfig saves a Spanner configuration
func HandleSaveSpannerConfig(w http.ResponseWriter, r *http.Request) {
	var newConfig config.SpannerConfig
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := config.AddOrUpdateSpannerConfig(newConfig); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// HandleSpannerSchema returns schema for a specific table
func HandleSpannerSchema(w http.ResponseWriter, r *http.Request) {
	var req struct {
		types.ConnectionRequest
		TableName string `json:"tableName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	schema, err := spanner.GetTableSchema(req.ConnectionRequest, req.TableName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schema)
}
