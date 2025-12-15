# Spanner Explorer Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add Spanner database explorer to test local Spanner emulator connections, browse tables, run SQL queries (SELECT and DML)

**Architecture:** Follow existing Testing Studio patterns - config management in internal/config, business logic in internal/spanner, HTTP handlers in internal/handlers, HTML templates in internal/templates. UI follows Spanner Studio layout with table browser sidebar and SQL editor.

**Tech Stack:** Go 1.21+, cloud.google.com/go/spanner, HTML/CSS/JavaScript (vanilla), Material Design styling

---

## Task 1: Add Spanner Dependencies

**Files:**
- Modify: `go.mod`
- Modify: `go.sum`

**Step 1: Add Spanner client dependency**

Run:
```bash
go get cloud.google.com/go/spanner@latest
```

Expected: go.mod and go.sum updated with Spanner SDK

**Step 2: Add Google API iterator dependency**

Run:
```bash
go get google.golang.org/api/iterator@latest
```

Expected: go.mod and go.sum updated

**Step 3: Verify dependencies**

Run:
```bash
go mod tidy
```

Expected: All dependencies clean

**Step 4: Commit**

```bash
git add go.mod go.sum
git commit -m "feat: add Spanner client dependencies"
```

---

## Task 2: Add Spanner Configuration Types

**Files:**
- Modify: `internal/config/config.go`

**Step 1: Add SpannerConfig struct**

Add after line 24 (after KafkaConfig struct):

```go
type SpannerConfig struct {
	Name         string `json:"name"`
	EmulatorHost string `json:"emulatorHost"`
	ProjectID    string `json:"projectId"`
	InstanceID   string `json:"instanceId"`
	DatabaseID   string `json:"databaseId"`
}
```

**Step 2: Add SpannerConfigs field to Config struct**

Modify Config struct (around line 26-29):

```go
type Config struct {
	PubSubConfigs []PubSubConfig `json:"pubsubConfigs"`
	KafkaConfigs  []KafkaConfig  `json:"kafkaConfigs"`
	SpannerConfigs []SpannerConfig `json:"spannerConfigs"`
}
```

**Step 3: Add default Spanner config in Load function**

Modify Load() function to include default Spanner config (around line 60):

```go
SpannerConfigs: []SpannerConfig{
	{
		Name:         "TMS Local",
		EmulatorHost: "localhost:9010",
		ProjectID:    "tms-suncorp-local",
		InstanceID:   "tms-suncorp-local",
		DatabaseID:   "tms-suncorp-db",
	},
},
```

**Step 4: Add AddOrUpdateSpannerConfig function**

Add after AddOrUpdateKafkaConfig function (around line 125):

```go
func AddOrUpdateSpannerConfig(newConfig SpannerConfig) error {
	mu.Lock()
	found := false
	for i, cfg := range config.SpannerConfigs {
		if cfg.Name == newConfig.Name {
			config.SpannerConfigs[i] = newConfig
			found = true
			break
		}
	}
	if !found {
		config.SpannerConfigs = append(config.SpannerConfigs, newConfig)
	}
	mu.Unlock()

	return Save()
}
```

**Step 5: Verify code compiles**

Run:
```bash
go build ./...
```

Expected: No errors

**Step 6: Commit**

```bash
git add internal/config/config.go
git commit -m "feat: add Spanner configuration types"
```

---

## Task 3: Create Spanner Types

**Files:**
- Create: `internal/types/spanner.go`

**Step 1: Create spanner types file**

Create file `internal/types/spanner.go`:

```go
package types

// QueryRequest represents a SQL query request
type QueryRequest struct {
	EmulatorHost string `json:"emulatorHost"`
	ProjectID    string `json:"projectId"`
	InstanceID   string `json:"instanceId"`
	DatabaseID   string `json:"databaseId"`
	Query        string `json:"query"`
}

// QueryResponse represents the result of a SQL query
type QueryResponse struct {
	Columns      []string                 `json:"columns"`
	Rows         []map[string]interface{} `json:"rows"`
	RowCount     int                      `json:"rowCount"`
	ExecutionTime string                   `json:"executionTime"`
	Error        string                   `json:"error,omitempty"`
}

// TableInfo represents metadata about a table
type TableInfo struct {
	Name       string       `json:"name"`
	RowCount   int64        `json:"rowCount,omitempty"`
	Columns    []ColumnInfo `json:"columns,omitempty"`
}

// ColumnInfo represents metadata about a column
type ColumnInfo struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	IsNullable bool   `json:"isNullable"`
	IsPrimaryKey bool `json:"isPrimaryKey"`
}

// ConnectionRequest represents a connection test request
type ConnectionRequest struct {
	EmulatorHost string `json:"emulatorHost"`
	ProjectID    string `json:"projectId"`
	InstanceID   string `json:"instanceId"`
	DatabaseID   string `json:"databaseId"`
}

// ConnectionResponse represents the result of a connection test
type ConnectionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
```

**Step 2: Verify code compiles**

Run:
```bash
go build ./...
```

Expected: No errors

**Step 3: Commit**

```bash
git add internal/types/spanner.go
git commit -m "feat: add Spanner types for API requests/responses"
```

---

## Task 4: Create Spanner Client Operations

**Files:**
- Create: `internal/spanner/spanner.go`

**Step 1: Create spanner package directory**

Run:
```bash
mkdir -p internal/spanner
```

**Step 2: Create spanner client file**

Create file `internal/spanner/spanner.go`:

```go
package spanner

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/spanner"
	"cloudevents-explorer/internal/types"
	"google.golang.org/api/iterator"
)

// TestConnection tests the connection to Spanner emulator
func TestConnection(req types.ConnectionRequest) types.ConnectionResponse {
	// Set emulator host environment variable
	if req.EmulatorHost != "" {
		os.Setenv("SPANNER_EMULATOR_HOST", req.EmulatorHost)
	}

	ctx := context.Background()
	dbPath := fmt.Sprintf("projects/%s/instances/%s/databases/%s",
		req.ProjectID, req.InstanceID, req.DatabaseID)

	client, err := spanner.NewClient(ctx, dbPath)
	if err != nil {
		return types.ConnectionResponse{
			Success: false,
			Message: "Failed to connect",
			Error:   err.Error(),
		}
	}
	defer client.Close()

	// Try a simple query to verify connection
	stmt := spanner.Statement{SQL: "SELECT 1"}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	_, err = iter.Next()
	if err != nil && err != iterator.Done {
		return types.ConnectionResponse{
			Success: false,
			Message: "Connection failed",
			Error:   err.Error(),
		}
	}

	return types.ConnectionResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully connected to %s", dbPath),
	}
}

// ListTables returns all tables in the database
func ListTables(req types.ConnectionRequest) ([]types.TableInfo, error) {
	if req.EmulatorHost != "" {
		os.Setenv("SPANNER_EMULATOR_HOST", req.EmulatorHost)
	}

	ctx := context.Background()
	dbPath := fmt.Sprintf("projects/%s/instances/%s/databases/%s",
		req.ProjectID, req.InstanceID, req.DatabaseID)

	client, err := spanner.NewClient(ctx, dbPath)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = ''
		ORDER BY table_name
	`

	stmt := spanner.Statement{SQL: query}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	var tables []types.TableInfo
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var tableName string
		if err := row.Columns(&tableName); err != nil {
			return nil, err
		}

		tables = append(tables, types.TableInfo{Name: tableName})
	}

	return tables, nil
}

// GetTableSchema returns the schema for a specific table
func GetTableSchema(req types.ConnectionRequest, tableName string) (*types.TableInfo, error) {
	if req.EmulatorHost != "" {
		os.Setenv("SPANNER_EMULATOR_HOST", req.EmulatorHost)
	}

	ctx := context.Background()
	dbPath := fmt.Sprintf("projects/%s/instances/%s/databases/%s",
		req.ProjectID, req.InstanceID, req.DatabaseID)

	client, err := spanner.NewClient(ctx, dbPath)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	query := `
		SELECT
			column_name,
			spanner_type,
			is_nullable
		FROM information_schema.columns
		WHERE table_name = @tableName
		ORDER BY ordinal_position
	`

	stmt := spanner.Statement{
		SQL: query,
		Params: map[string]interface{}{
			"tableName": tableName,
		},
	}

	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	var columns []types.ColumnInfo
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var colName, colType, nullable string
		if err := row.Columns(&colName, &colType, &nullable); err != nil {
			return nil, err
		}

		columns = append(columns, types.ColumnInfo{
			Name:       colName,
			Type:       colType,
			IsNullable: nullable == "YES",
		})
	}

	return &types.TableInfo{
		Name:    tableName,
		Columns: columns,
	}, nil
}

// ExecuteQuery executes a SQL query and returns results
func ExecuteQuery(req types.QueryRequest) types.QueryResponse {
	if req.EmulatorHost != "" {
		os.Setenv("SPANNER_EMULATOR_HOST", req.EmulatorHost)
	}

	ctx := context.Background()
	dbPath := fmt.Sprintf("projects/%s/instances/%s/databases/%s",
		req.ProjectID, req.InstanceID, req.DatabaseID)

	client, err := spanner.NewClient(ctx, dbPath)
	if err != nil {
		return types.QueryResponse{
			Error: fmt.Sprintf("Failed to create client: %v", err),
		}
	}
	defer client.Close()

	startTime := time.Now()

	// Detect if this is a DML statement
	queryUpper := strings.ToUpper(strings.TrimSpace(req.Query))
	isDML := strings.HasPrefix(queryUpper, "INSERT") ||
		strings.HasPrefix(queryUpper, "UPDATE") ||
		strings.HasPrefix(queryUpper, "DELETE")

	if isDML {
		// Execute DML
		_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
			stmt := spanner.Statement{SQL: req.Query}
			_, err := txn.Update(ctx, stmt)
			return err
		})

		executionTime := time.Since(startTime).String()

		if err != nil {
			return types.QueryResponse{
				Error:         err.Error(),
				ExecutionTime: executionTime,
			}
		}

		return types.QueryResponse{
			Columns:       []string{"Status"},
			Rows:          []map[string]interface{}{{"Status": "DML executed successfully"}},
			RowCount:      1,
			ExecutionTime: executionTime,
		}
	}

	// Execute SELECT query
	stmt := spanner.Statement{SQL: req.Query}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	var columns []string
	var rows []map[string]interface{}

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return types.QueryResponse{
				Error:         err.Error(),
				ExecutionTime: time.Since(startTime).String(),
			}
		}

		// Get column names from first row
		if columns == nil {
			columns = row.ColumnNames()
		}

		// Convert row to map
		rowMap := make(map[string]interface{})
		values := make([]interface{}, len(columns))
		ptrs := make([]interface{}, len(columns))
		for i := range values {
			ptrs[i] = &values[i]
		}

		if err := row.Columns(ptrs...); err != nil {
			return types.QueryResponse{
				Error:         err.Error(),
				ExecutionTime: time.Since(startTime).String(),
			}
		}

		for i, col := range columns {
			rowMap[col] = values[i]
		}

		rows = append(rows, rowMap)
	}

	executionTime := time.Since(startTime).String()

	return types.QueryResponse{
		Columns:       columns,
		Rows:          rows,
		RowCount:      len(rows),
		ExecutionTime: executionTime,
	}
}
```

**Step 3: Verify code compiles**

Run:
```bash
go build ./...
```

Expected: No errors

**Step 4: Commit**

```bash
git add internal/spanner/spanner.go
git commit -m "feat: add Spanner client operations for connection, tables, queries"
```

---

## Task 5: Create Spanner HTTP Handlers

**Files:**
- Create: `internal/handlers/spanner.go`

**Step 1: Create spanner handlers file**

Create file `internal/handlers/spanner.go`:

```go
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
```

**Step 2: Verify code compiles**

Run:
```bash
go build ./...
```

Expected: No errors

**Step 3: Commit**

```bash
git add internal/handlers/spanner.go
git commit -m "feat: add Spanner HTTP handlers for API endpoints"
```

---

## Task 6: Create Spanner HTML Template (Part 1 - Structure)

**Files:**
- Create: `internal/templates/spanner.go`

**Step 1: Create spanner template file with base structure**

Create file `internal/templates/spanner.go` with the HTML structure:

```go
package templates

const SpannerContent = `
<div class="panel">
    <div class="panel-header">
        <div class="panel-title">Connection Settings</div>
    </div>
    <div class="panel-body">
        <div class="form-row">
            <div class="form-group">
                <label for="configSelect">Configuration Profile</label>
                <select id="configSelect" onchange="loadConfig()">
                    <option value="">-- New Configuration --</option>
                </select>
            </div>
            <div class="form-group">
                <label for="configName">Profile Name</label>
                <input type="text" id="configName" placeholder="TMS Local">
            </div>
        </div>
        <div class="form-row">
            <div class="form-group">
                <label for="emulatorHost">Emulator Host</label>
                <input type="text" id="emulatorHost" placeholder="localhost:9010">
            </div>
            <div class="form-group">
                <label for="projectId">Project ID</label>
                <input type="text" id="projectId" placeholder="tms-suncorp-local">
            </div>
        </div>
        <div class="form-row">
            <div class="form-group">
                <label for="instanceId">Instance ID</label>
                <input type="text" id="instanceId" placeholder="tms-suncorp-local">
            </div>
            <div class="form-group">
                <label for="databaseId">Database ID</label>
                <input type="text" id="databaseId" placeholder="tms-suncorp-db">
            </div>
        </div>
        <div class="button-group">
            <button class="btn-primary" onclick="testConnection()">Connect</button>
            <button class="btn-secondary" onclick="saveConfig()">Save Configuration</button>
            <button class="btn-secondary" onclick="loadTables()">Load Tables</button>
        </div>
        <div id="connectionStatus" style="margin-top: 12px; padding: 8px; border-radius: 4px; display: none;"></div>
    </div>
</div>

<div style="display: grid; grid-template-columns: 250px 1fr; gap: 16px; height: calc(100vh - 400px); min-height: 600px;">
    <!-- Table Browser Sidebar -->
    <div class="panel" style="height: 100%; display: flex; flex-direction: column;">
        <div class="panel-header" style="flex-shrink: 0;">
            <div class="panel-title">Tables</div>
        </div>
        <div style="padding: 12px; flex-shrink: 0;">
            <input type="text" id="tableSearch" placeholder="Search tables..."
                   onkeyup="filterTables()"
                   style="width: 100%; padding: 6px 8px; font-size: 13px;">
        </div>
        <div style="flex: 1; overflow-y: auto; padding: 0 12px 12px 12px;">
            <div id="tableList" style="display: flex; flex-direction: column; gap: 4px;">
                <div style="color: #5f6368; font-size: 13px; padding: 20px; text-align: center;">
                    Click "Load Tables" to view tables
                </div>
            </div>
        </div>
    </div>

    <!-- Main Content Area -->
    <div style="display: flex; flex-direction: column; gap: 16px; height: 100%;">
        <!-- SQL Editor -->
        <div class="panel" style="flex: 1; min-height: 250px; display: flex; flex-direction: column;">
            <div class="panel-header" style="flex-shrink: 0;">
                <div class="panel-title">SQL Editor</div>
            </div>
            <div class="panel-body" style="flex: 1; display: flex; flex-direction: column;">
                <div style="flex: 1; display: flex; flex-direction: column; gap: 8px;">
                    <textarea id="sqlQuery"
                              style="flex: 1; min-height: 150px; font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
                                     font-size: 13px; resize: vertical; padding: 12px;"
                              placeholder="-- Enter SQL query here&#10;SELECT * FROM TableName LIMIT 10;"></textarea>
                    <div class="button-group">
                        <button class="btn-primary" onclick="executeQuery()">Run Query</button>
                        <select id="exampleQueries" onchange="loadExampleQuery()" style="padding: 8px 12px;">
                            <option value="">-- Example Queries --</option>
                            <option value="SHOW_TABLES">Show all tables</option>
                            <option value="SELECT_ALL">SELECT * FROM (selected table)</option>
                            <option value="COUNT">Count rows in (selected table)</option>
                        </select>
                    </div>
                </div>
            </div>
        </div>

        <!-- Results Panel -->
        <div class="panel" style="flex: 1; min-height: 250px; display: flex; flex-direction: column;">
            <div class="panel-header" style="flex-shrink: 0;">
                <div class="panel-title">Results</div>
            </div>
            <div id="queryStats" style="display: none; padding: 8px 20px; background: #e8f5e9; border-bottom: 1px solid #dadce0; font-size: 13px; color: #188038; flex-shrink: 0;"></div>
            <div id="queryError" style="display: none; padding: 12px 20px; background: #fce8e6; border-bottom: 1px solid #dadce0; font-size: 13px; color: #d93025; flex-shrink: 0;"></div>
            <div class="panel-body" style="flex: 1; overflow: auto; padding: 0;">
                <div id="queryResults" style="padding: 20px; color: #5f6368; font-size: 13px;">
                    Run a query to see results here
                </div>
            </div>
        </div>
    </div>
</div>
`

func GetSpannerHTML() string {
	return GetBaseHTML("Spanner Explorer", SpannerContent, SpannerJS)
}
```

**Step 2: Verify code compiles**

Run:
```bash
go build ./...
```

Expected: Compilation error - SpannerJS not defined (expected, we'll add it next)

---

## Task 7: Create Spanner HTML Template (Part 2 - JavaScript)

**Files:**
- Modify: `internal/templates/spanner.go`

**Step 1: Add JavaScript code to spanner.go**

Add this JavaScript constant before the GetSpannerHTML function:

```go
const SpannerJS = `
let currentConfig = {};
let allTables = [];
let selectedTable = '';

// Load configurations on page load
async function loadConfigurations() {
    try {
        const response = await fetch('/api/configs');
        const data = await response.json();

        const select = document.getElementById('configSelect');
        select.innerHTML = '<option value="">-- New Configuration --</option>';

        if (data.spannerConfigs && data.spannerConfigs.length > 0) {
            data.spannerConfigs.forEach(cfg => {
                const option = document.createElement('option');
                option.value = cfg.name;
                option.textContent = cfg.name;
                select.appendChild(option);
            });

            // Auto-load first config
            loadConfigByName(data.spannerConfigs[0].name, data.spannerConfigs);
        } else {
            // Load from environment variables if no configs
            loadFromEnvironment();
        }
    } catch (error) {
        console.error('Failed to load configurations:', error);
        loadFromEnvironment();
    }
}

function loadFromEnvironment() {
    // These will be empty in browser, but useful for documentation
    document.getElementById('emulatorHost').value = 'localhost:9010';
    document.getElementById('projectId').value = 'tms-suncorp-local';
    document.getElementById('instanceId').value = 'tms-suncorp-local';
    document.getElementById('databaseId').value = 'tms-suncorp-db';
}

function loadConfigByName(name, configs) {
    const config = configs.find(c => c.name === name);
    if (config) {
        currentConfig = config;
        document.getElementById('configName').value = config.name;
        document.getElementById('emulatorHost').value = config.emulatorHost;
        document.getElementById('projectId').value = config.projectId;
        document.getElementById('instanceId').value = config.instanceId;
        document.getElementById('databaseId').value = config.databaseId;
        document.getElementById('configSelect').value = name;
    }
}

async function loadConfig() {
    const select = document.getElementById('configSelect');
    const selectedName = select.value;

    if (!selectedName) {
        // Clear form for new config
        document.getElementById('configName').value = '';
        document.getElementById('emulatorHost').value = '';
        document.getElementById('projectId').value = '';
        document.getElementById('instanceId').value = '';
        document.getElementById('databaseId').value = '';
        return;
    }

    try {
        const response = await fetch('/api/configs');
        const data = await response.json();
        loadConfigByName(selectedName, data.spannerConfigs);
    } catch (error) {
        showStatus('Failed to load configuration: ' + error.message, true);
    }
}

async function testConnection() {
    const connectionReq = {
        emulatorHost: document.getElementById('emulatorHost').value,
        projectId: document.getElementById('projectId').value,
        instanceId: document.getElementById('instanceId').value,
        databaseId: document.getElementById('databaseId').value
    };

    const statusDiv = document.getElementById('connectionStatus');
    statusDiv.style.display = 'block';
    statusDiv.style.background = '#e8f5e9';
    statusDiv.style.color = '#188038';
    statusDiv.textContent = 'Connecting...';

    try {
        const response = await fetch('/api/spanner/connect', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(connectionReq)
        });

        const result = await response.json();

        if (result.success) {
            statusDiv.style.background = '#e8f5e9';
            statusDiv.style.color = '#188038';
            statusDiv.textContent = '✓ ' + result.message;
            showStatus('Connection successful!');
        } else {
            statusDiv.style.background = '#fce8e6';
            statusDiv.style.color = '#d93025';
            statusDiv.textContent = '✗ ' + (result.error || result.message);
            showStatus('Connection failed', true);
        }
    } catch (error) {
        statusDiv.style.background = '#fce8e6';
        statusDiv.style.color = '#d93025';
        statusDiv.textContent = '✗ Error: ' + error.message;
        showStatus('Connection error: ' + error.message, true);
    }
}

async function saveConfig() {
    const config = {
        name: document.getElementById('configName').value,
        emulatorHost: document.getElementById('emulatorHost').value,
        projectId: document.getElementById('projectId').value,
        instanceId: document.getElementById('instanceId').value,
        databaseId: document.getElementById('databaseId').value
    };

    if (!config.name) {
        showStatus('Please enter a profile name', true);
        return;
    }

    try {
        const response = await fetch('/api/spanner/configs', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(config)
        });

        if (response.ok) {
            showStatus('Configuration saved successfully!');
            await loadConfigurations();
            document.getElementById('configSelect').value = config.name;
        } else {
            showStatus('Failed to save configuration', true);
        }
    } catch (error) {
        showStatus('Error saving configuration: ' + error.message, true);
    }
}

async function loadTables() {
    const connectionReq = {
        emulatorHost: document.getElementById('emulatorHost').value,
        projectId: document.getElementById('projectId').value,
        instanceId: document.getElementById('instanceId').value,
        databaseId: document.getElementById('databaseId').value
    };

    const tableList = document.getElementById('tableList');
    tableList.innerHTML = '<div style="color: #5f6368; font-size: 13px; padding: 20px; text-align: center;">Loading tables...</div>';

    try {
        const response = await fetch('/api/spanner/tables', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(connectionReq)
        });

        const tables = await response.json();

        if (tables.error) {
            tableList.innerHTML = '<div style="color: #d93025; font-size: 13px; padding: 20px; text-align: center;">Error: ' + tables.error + '</div>';
            return;
        }

        allTables = tables;
        renderTables(tables);
        showStatus('Loaded ' + tables.length + ' tables');
    } catch (error) {
        tableList.innerHTML = '<div style="color: #d93025; font-size: 13px; padding: 20px; text-align: center;">Error: ' + error.message + '</div>';
        showStatus('Failed to load tables: ' + error.message, true);
    }
}

function renderTables(tables) {
    const tableList = document.getElementById('tableList');

    if (tables.length === 0) {
        tableList.innerHTML = '<div style="color: #5f6368; font-size: 13px; padding: 20px; text-align: center;">No tables found</div>';
        return;
    }

    tableList.innerHTML = '';
    tables.forEach(table => {
        const div = document.createElement('div');
        div.textContent = table.name;
        div.style.cssText = 'padding: 8px 12px; cursor: pointer; border-radius: 4px; font-size: 13px; transition: background 0.2s;';
        div.onmouseover = () => div.style.background = '#f1f3f4';
        div.onmouseout = () => {
            div.style.background = selectedTable === table.name ? '#e8f0fe' : '';
        };
        div.onclick = () => selectTable(table.name);

        if (selectedTable === table.name) {
            div.style.background = '#e8f0fe';
        }

        tableList.appendChild(div);
    });
}

function filterTables() {
    const searchText = document.getElementById('tableSearch').value.toLowerCase();
    const filtered = allTables.filter(t => t.name.toLowerCase().includes(searchText));
    renderTables(filtered);
}

function selectTable(tableName) {
    selectedTable = tableName;
    renderTables(allTables);

    // Auto-fill query
    document.getElementById('sqlQuery').value = 'SELECT * FROM ' + tableName + ' LIMIT 10;';
}

function loadExampleQuery() {
    const select = document.getElementById('exampleQueries');
    const value = select.value;
    const queryArea = document.getElementById('sqlQuery');

    if (value === 'SHOW_TABLES') {
        queryArea.value = "SELECT table_name FROM information_schema.tables WHERE table_schema = '' ORDER BY table_name;";
    } else if (value === 'SELECT_ALL' && selectedTable) {
        queryArea.value = 'SELECT * FROM ' + selectedTable + ' LIMIT 10;';
    } else if (value === 'COUNT' && selectedTable) {
        queryArea.value = 'SELECT COUNT(*) as row_count FROM ' + selectedTable + ';';
    }

    select.value = '';
}

async function executeQuery() {
    const query = document.getElementById('sqlQuery').value.trim();

    if (!query) {
        showStatus('Please enter a SQL query', true);
        return;
    }

    const queryReq = {
        emulatorHost: document.getElementById('emulatorHost').value,
        projectId: document.getElementById('projectId').value,
        instanceId: document.getElementById('instanceId').value,
        databaseId: document.getElementById('databaseId').value,
        query: query
    };

    // Hide previous results/errors
    document.getElementById('queryStats').style.display = 'none';
    document.getElementById('queryError').style.display = 'none';
    document.getElementById('queryResults').innerHTML = '<div style="padding: 20px; color: #5f6368; font-size: 13px;">Executing query...</div>';

    try {
        const response = await fetch('/api/spanner/query', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(queryReq)
        });

        const result = await response.json();

        if (result.error) {
            document.getElementById('queryError').style.display = 'block';
            document.getElementById('queryError').textContent = 'Error: ' + result.error;
            document.getElementById('queryResults').innerHTML = '';
            showStatus('Query failed', true);
            return;
        }

        // Show success stats
        const statsDiv = document.getElementById('queryStats');
        statsDiv.style.display = 'block';
        statsDiv.textContent = '✓ Query executed successfully. Rows: ' + result.rowCount + ' | Time: ' + result.executionTime;

        // Render results table
        if (result.rows && result.rows.length > 0) {
            renderResultsTable(result.columns, result.rows);
        } else {
            document.getElementById('queryResults').innerHTML = '<div style="padding: 20px; color: #5f6368; font-size: 13px;">Query returned no rows</div>';
        }

        showStatus('Query executed successfully');
    } catch (error) {
        document.getElementById('queryError').style.display = 'block';
        document.getElementById('queryError').textContent = 'Error: ' + error.message;
        document.getElementById('queryResults').innerHTML = '';
        showStatus('Query execution error: ' + error.message, true);
    }
}

function renderResultsTable(columns, rows) {
    const resultsDiv = document.getElementById('queryResults');

    let html = '<table style="width: 100%; border-collapse: collapse; font-size: 13px;">';

    // Header
    html += '<thead><tr style="background: #f8f9fa; border-bottom: 2px solid #dadce0;">';
    columns.forEach(col => {
        html += '<th style="padding: 12px; text-align: left; font-weight: 600; color: #5f6368;">' + col + '</th>';
    });
    html += '</tr></thead>';

    // Rows
    html += '<tbody>';
    rows.forEach((row, idx) => {
        const bgColor = idx % 2 === 0 ? 'white' : '#f8f9fa';
        html += '<tr style="background: ' + bgColor + '; border-bottom: 1px solid #e8eaed;">';
        columns.forEach(col => {
            let value = row[col];
            if (value === null || value === undefined) {
                value = '<span style="color: #5f6368; font-style: italic;">NULL</span>';
            } else if (typeof value === 'object') {
                value = JSON.stringify(value);
            }
            html += '<td style="padding: 10px; color: #202124;">' + value + '</td>';
        });
        html += '</tr>';
    });
    html += '</tbody></table>';

    resultsDiv.innerHTML = html;
}

// Load configurations on page load
loadConfigurations();
`
```

**Step 2: Verify code compiles**

Run:
```bash
go build ./...
```

Expected: No errors

**Step 3: Commit**

```bash
git add internal/templates/spanner.go
git commit -m "feat: add Spanner HTML template with UI and JavaScript"
```

---

## Task 8: Register Spanner Routes in Main

**Files:**
- Modify: `cmd/server/main.go`

**Step 1: Add Spanner page route**

Add after line 25 (after trace-journey route):

```go
http.HandleFunc("/spanner", handlers.HandleSpanner)
```

**Step 2: Add Spanner API routes**

Add after line 42 (after trace search API route):

```go
http.HandleFunc("/api/spanner/connect", handlers.HandleSpannerConnect)
http.HandleFunc("/api/spanner/tables", handlers.HandleSpannerTables)
http.HandleFunc("/api/spanner/query", handlers.HandleSpannerQuery)
http.HandleFunc("/api/spanner/configs", handlers.HandleSaveSpannerConfig)
http.HandleFunc("/api/spanner/schema", handlers.HandleSpannerSchema)
```

**Step 3: Verify code compiles**

Run:
```bash
go build ./...
```

Expected: No errors

**Step 4: Build the binary**

Run:
```bash
go build -o testing-studio cmd/server/main.go
```

Expected: Binary created successfully

**Step 5: Commit**

```bash
git add cmd/server/main.go
git commit -m "feat: register Spanner routes in main server"
```

---

## Task 9: Add Spanner Card to Homepage

**Files:**
- Modify: `internal/templates/index.go`

**Step 1: Add Spanner card HTML**

Add after the Trace Journey card (around line 233, before the closing </div> of optionsGrid):

```go
            <a href="/spanner" class="option-card" id="cardSpanner">
                <div class="option-title" id="titleSpanner">Spanner Explorer</div>
                <div class="option-desc" id="descSpanner">Browse tables, run SQL queries, test local Spanner emulator</div>
                <span class="badge" id="badgeSpanner">Database</span>
            </a>
```

**Step 2: Verify code compiles**

Run:
```bash
go build ./...
```

Expected: No errors

**Step 3: Test build**

Run:
```bash
go build -o testing-studio cmd/server/main.go
```

Expected: Binary created successfully

**Step 4: Commit**

```bash
git add internal/templates/index.go
git commit -m "feat: add Spanner Explorer card to homepage"
```

---

## Task 10: Manual Testing and Verification

**Files:**
- None (manual testing)

**Step 1: Start Spanner emulator (if available)**

Optional - only if you have Spanner emulator:
```bash
docker run -p 9010:9010 gcr.io/cloud-spanner-emulator/emulator
```

**Step 2: Start Testing Studio**

Run:
```bash
./testing-studio
```

Expected: Server starts on http://localhost:8888

**Step 3: Verify homepage shows Spanner card**

Open browser to http://localhost:8888
Expected: See "Spanner Explorer" card in the grid

**Step 4: Click Spanner card and verify page loads**

Click the Spanner Explorer card
Expected: Spanner page loads with connection form, table browser, and SQL editor

**Step 5: Test connection (if emulator running)**

Fill in connection details and click "Connect"
Expected: Connection status shows success or appropriate error

**Step 6: Test configuration save**

Enter a profile name and click "Save Configuration"
Expected: Configuration saved, appears in dropdown

**Step 7: Stop server and commit**

Press Ctrl+C to stop server

```bash
git add -A
git commit -m "test: verify Spanner Explorer functionality"
```

---

## Task 11: Update Documentation

**Files:**
- Modify: `README.md`

**Step 1: Add Spanner to features list**

Add after the Flow Diagram Tool section (around line 38):

```markdown
### 🗄️ Spanner Explorer
- Connect to local Spanner emulator
- Browse tables and view schemas
- Execute SQL queries (SELECT and DML)
- Save multiple connection profiles
- Spanner Studio-style interface
```

**Step 2: Add Spanner to usage guide**

Add after the REST Client section (around line 121):

```markdown
### Spanner Explorer
1. Navigate to **http://localhost:8888/spanner**
2. Configure connection settings:
   - **Emulator Host**: `localhost:9010`
   - **Project ID**: `tms-suncorp-local`
   - **Instance ID**: `tms-suncorp-local`
   - **Database ID**: `tms-suncorp-db`
3. Click **Connect** to test connection
4. Click **Load Tables** to browse database tables
5. Click a table to select it
6. Enter SQL query in editor or use example queries
7. Click **Run Query** to execute

**Local Development**: Requires Spanner emulator running on `localhost:9010`
```

**Step 3: Add Spanner emulator to running services**

Add after fake-gcs-server section (around line 165):

```markdown
**Spanner Emulator**:
```bash
docker run -d --name spanner-emulator \
  -p 9010:9010 \
  gcr.io/cloud-spanner-emulator/emulator
```
```

**Step 4: Add Spanner to configuration example**

Update configs.json example (around line 188) to include:

```json
  "spanner": [
    {
      "name": "TMS Local",
      "emulatorHost": "localhost:9010",
      "projectId": "tms-suncorp-local",
      "instanceId": "tms-suncorp-local",
      "databaseId": "tms-suncorp-db"
    }
  ]
```

**Step 5: Add Spanner troubleshooting section**

Add after GCS Browser troubleshooting (around line 231):

```markdown
### Can't connect to Spanner emulator
- Verify emulator is running: `docker ps | grep spanner`
- Check emulator host in configuration
- Ensure database exists in emulator
- Check emulator logs: `docker logs spanner-emulator`
```

**Step 6: Commit documentation**

```bash
git add README.md
git commit -m "docs: add Spanner Explorer documentation to README"
```

---

## Task 12: Final Build and Verification

**Files:**
- None (verification only)

**Step 1: Clean build**

Run:
```bash
go clean
go build -o testing-studio cmd/server/main.go
```

Expected: Clean build with no errors

**Step 2: Run final tests**

Run:
```bash
go test ./... -v
```

Expected: Tests run (integration tests will fail without server, expected)

**Step 3: Verify all files committed**

Run:
```bash
git status
```

Expected: Working tree clean

**Step 4: Review commit history**

Run:
```bash
git log --oneline -12
```

Expected: See all 12 commits from this implementation

**Step 5: Final commit (if needed)**

If any uncommitted changes:
```bash
git add -A
git commit -m "chore: final cleanup and verification"
```

---

## Success Criteria

✅ Spanner dependencies added to go.mod
✅ Configuration types for Spanner created
✅ Spanner client operations implemented (connect, list tables, query, schema)
✅ HTTP handlers created for all Spanner API endpoints
✅ HTML template created with Spanner Studio-style UI
✅ Routes registered in main.go
✅ Homepage card added
✅ README documentation updated
✅ All code compiles without errors
✅ Manual testing successful (connection, tables, queries)

## Notes

- This implementation follows existing Testing Studio patterns exactly
- No breaking changes to existing functionality
- All code is production-ready for local testing
- UI matches Google Cloud Console Material Design style
- Supports both SELECT and DML queries
- Configuration persists across sessions
- Environment variables as fallback for first-time setup
