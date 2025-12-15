# Spanner Explorer Manual Test Plan

## Test Environment Setup

### Prerequisites
1. Spanner emulator running on `localhost:9010`
   ```bash
   docker run -p 9010:9010 gcr.io/cloud-spanner-emulator/emulator
   ```
2. Testing Studio server running
   ```bash
   ./testing-studio
   ```
3. Browser opened to http://localhost:8888

## Test Cases

### TC1: Homepage Integration
**Objective**: Verify Spanner Explorer card appears on homepage

**Steps**:
1. Navigate to http://localhost:8888
2. Locate the "Spanner Explorer" card in the grid

**Expected Results**:
- Spanner Explorer card visible with title "Spanner Explorer"
- Description reads: "Browse tables, run SQL queries, test local Spanner emulator"
- Badge shows "Database"
- Card is clickable

**Status**: READY FOR TESTING

---

### TC2: Page Load and UI Rendering
**Objective**: Verify Spanner Explorer page loads correctly

**Steps**:
1. Click on "Spanner Explorer" card from homepage
2. Wait for page to load

**Expected Results**:
- Page loads at http://localhost:8888/spanner
- Title shows "Spanner Explorer"
- Three main sections visible:
  - Connection Settings panel (top)
  - Table Browser sidebar (left)
  - SQL Editor and Results panels (right)
- All form fields present:
  - Configuration Profile dropdown
  - Profile Name input
  - Emulator Host input
  - Project ID input
  - Instance ID input
  - Database ID input
- Buttons present:
  - Connect
  - Save Configuration
  - Load Tables
  - Run Query
- Example Queries dropdown present

**Status**: READY FOR TESTING

---

### TC3: Default Configuration Load
**Objective**: Verify default configuration loads on page load

**Steps**:
1. Navigate to /spanner page
2. Observe form fields

**Expected Results**:
- Configuration dropdown shows "TMS Local"
- Form fields populated with default values:
  - Emulator Host: `localhost:9010`
  - Project ID: `tms-suncorp-local`
  - Instance ID: `tms-suncorp-local`
  - Database ID: `tms-suncorp-db`

**Status**: READY FOR TESTING

---

### TC4: Test Connection - Success
**Objective**: Verify successful connection to Spanner emulator

**Prerequisites**: Spanner emulator running with test database

**Steps**:
1. Ensure default configuration is loaded
2. Click "Connect" button
3. Observe connection status

**Expected Results**:
- Connection status div appears
- Status shows green background (#e8f5e9)
- Message shows "✓ Successfully connected to projects/tms-suncorp-local/instances/tms-suncorp-local/databases/tms-suncorp-db"
- No error messages

**Status**: READY FOR TESTING

---

### TC5: Test Connection - Failure
**Objective**: Verify error handling when emulator is not available

**Prerequisites**: Spanner emulator NOT running

**Steps**:
1. Ensure connection settings point to localhost:9010
2. Click "Connect" button
3. Observe connection status

**Expected Results**:
- Connection status div appears
- Status shows red background (#fce8e6)
- Error message displayed with "✗" prefix
- Error describes connection failure

**Status**: READY FOR TESTING

---

### TC6: Save Configuration
**Objective**: Verify configuration can be saved

**Steps**:
1. Enter new profile name: "Test Config"
2. Modify connection settings
3. Click "Save Configuration" button
4. Refresh page

**Expected Results**:
- Success message appears
- Configuration dropdown updates to include "Test Config"
- After refresh, "Test Config" appears in dropdown
- Selecting "Test Config" from dropdown loads saved values

**Status**: READY FOR TESTING

---

### TC7: Load Tables
**Objective**: Verify table list loads from database

**Prerequisites**: Spanner emulator running with tables in database

**Steps**:
1. Connect to emulator successfully
2. Click "Load Tables" button
3. Observe table list sidebar

**Expected Results**:
- Loading indicator appears briefly
- Table names populate in sidebar
- Success message shows count: "Loaded X tables"
- Tables are sorted alphabetically
- Each table name is clickable

**Status**: READY FOR TESTING

---

### TC8: Table Search Filter
**Objective**: Verify table search functionality

**Prerequisites**: Multiple tables loaded

**Steps**:
1. Load tables successfully
2. Enter search term in "Search tables..." input
3. Observe filtered results

**Expected Results**:
- Table list filters in real-time
- Only matching table names shown
- Search is case-insensitive
- Clearing search shows all tables again

**Status**: READY FOR TESTING

---

### TC9: Select Table
**Objective**: Verify table selection functionality

**Steps**:
1. Load tables successfully
2. Click on a table name in sidebar

**Expected Results**:
- Selected table highlighted with blue background (#e8f0fe)
- SQL editor populates with: `SELECT * FROM [TableName] LIMIT 10;`
- Table remains selected when hovering over other tables

**Status**: READY FOR TESTING

---

### TC10: Execute SELECT Query
**Objective**: Verify SELECT query execution

**Prerequisites**: Connected to emulator with data in tables

**Steps**:
1. Enter or select a SELECT query
2. Click "Run Query" button
3. Observe results panel

**Expected Results**:
- Query stats appear in green: "✓ Query executed successfully. Rows: X | Time: Yms"
- Results table displays with:
  - Column headers matching query columns
  - Rows with data
  - Alternating row colors for readability
  - NULL values shown in italic gray
- Proper formatting for different data types

**Status**: READY FOR TESTING

---

### TC11: Execute DML Query
**Objective**: Verify DML query execution (INSERT/UPDATE/DELETE)

**Prerequisites**: Connected to emulator

**Steps**:
1. Enter INSERT, UPDATE, or DELETE statement
2. Click "Run Query" button
3. Observe results

**Expected Results**:
- Success stats appear
- Results show "DML executed successfully"
- No error messages
- Execution time displayed

**Status**: READY FOR TESTING

---

### TC12: Query Error Handling
**Objective**: Verify error handling for invalid queries

**Steps**:
1. Enter invalid SQL: `SELECT * FROM NonExistentTable;`
2. Click "Run Query" button
3. Observe results

**Expected Results**:
- Error panel appears in red (#fce8e6)
- Error message describes the issue
- No results table displayed
- Previous results cleared

**Status**: READY FOR TESTING

---

### TC13: Example Queries
**Objective**: Verify example query templates work

**Steps**:
1. Select "Show all tables" from Example Queries dropdown
2. Observe SQL editor
3. Clear editor
4. Select a table from sidebar
5. Select "SELECT * FROM (selected table)" from dropdown
6. Observe SQL editor

**Expected Results**:
- "Show all tables" inserts correct information_schema query
- Table-specific queries use selected table name
- Dropdown resets to placeholder after selection

**Status**: READY FOR TESTING

---

### TC14: Configuration Switching
**Objective**: Verify switching between saved configurations

**Prerequisites**: Multiple configurations saved

**Steps**:
1. Select different configuration from dropdown
2. Observe form fields

**Expected Results**:
- All form fields update to selected configuration values
- Previous connection state preserved

**Status**: READY FOR TESTING

---

### TC15: Empty Results Handling
**Objective**: Verify handling of queries with no results

**Steps**:
1. Execute query that returns no rows: `SELECT * FROM TableName WHERE 1=0;`
2. Observe results

**Expected Results**:
- Success stats show "Rows: 0"
- Results panel shows "Query returned no rows"
- No error displayed

**Status**: READY FOR TESTING

---

## File Verification Checklist

### Core Implementation Files
- ✅ internal/types/spanner.go - Type definitions
- ✅ internal/spanner/spanner.go - Spanner client operations
- ✅ internal/handlers/spanner.go - HTTP handlers
- ✅ internal/templates/spanner.go - HTML/JS template
- ✅ internal/config/config.go - Configuration types (modified)
- ✅ cmd/server/main.go - Routes registered (modified)
- ✅ internal/templates/index.go - Homepage card (modified)

### Dependencies
- ✅ go.mod - Spanner dependencies added
- ✅ go.sum - Dependency checksums

### Binary
- ✅ testing-studio - Compiled binary (54MB)

### Documentation
- ✅ README.md - User documentation (to be updated in Task 11)
- ✅ docs/plans/2025-12-15-spanner-explorer.md - Implementation plan

## Test Execution Summary

**Total Test Cases**: 15
**Status**: All test cases ready for manual execution
**Blockers**: Requires Spanner emulator and server runtime

## Notes

- All UI tests require visual inspection
- API endpoint tests can also be performed via curl/Postman
- Performance testing not included in this plan
- Security testing not included in this plan
- Browser compatibility testing not included (assume modern Chrome/Firefox/Safari)

## API Endpoint Testing (Alternative)

For automated/scripted testing, the following curl commands can be used:

```bash
# Test Connection
curl -X POST http://localhost:8888/api/spanner/connect \
  -H "Content-Type: application/json" \
  -d '{"emulatorHost":"localhost:9010","projectId":"tms-suncorp-local","instanceId":"tms-suncorp-local","databaseId":"tms-suncorp-db"}'

# List Tables
curl -X POST http://localhost:8888/api/spanner/tables \
  -H "Content-Type: application/json" \
  -d '{"emulatorHost":"localhost:9010","projectId":"tms-suncorp-local","instanceId":"tms-suncorp-local","databaseId":"tms-suncorp-db"}'

# Execute Query
curl -X POST http://localhost:8888/api/spanner/query \
  -H "Content-Type: application/json" \
  -d '{"emulatorHost":"localhost:9010","projectId":"tms-suncorp-local","instanceId":"tms-suncorp-local","databaseId":"tms-suncorp-db","query":"SELECT 1"}'

# Save Configuration
curl -X POST http://localhost:8888/api/spanner/configs \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","emulatorHost":"localhost:9010","projectId":"test","instanceId":"test","databaseId":"test"}'

# Get Table Schema
curl -X POST http://localhost:8888/api/spanner/schema \
  -H "Content-Type: application/json" \
  -d '{"emulatorHost":"localhost:9010","projectId":"tms-suncorp-local","instanceId":"tms-suncorp-local","databaseId":"tms-suncorp-db","tableName":"Users"}'
```
