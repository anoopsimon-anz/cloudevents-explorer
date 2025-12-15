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

		// Convert row to map - use string representation for all values
		rowMap := make(map[string]interface{})

		for i, col := range columns {
			// Use GenericColumnValue and convert to string for display
			var val spanner.GenericColumnValue
			if err := row.Column(i, &val); err != nil {
				rowMap[col] = fmt.Sprintf("Error: %v", err)
				continue
			}

			// Convert value to string for JSON serialization
			// This handles all Spanner types (STRING, INT64, FLOAT64, BOOL, BYTES, etc.)
			rowMap[col] = fmt.Sprintf("%v", val.Value)
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
