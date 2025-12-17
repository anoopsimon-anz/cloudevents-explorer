package templates

const SpannerContent = `
<div style="background: #e8f0fe; border: 1px solid #d2e3fc; border-radius: 4px; padding: 10px 16px; margin-bottom: 16px; font-size: 13px; color: #1967d2;">
    💡 <strong>Tip:</strong> You can manage all Spanner connection profiles in <a href="/config-editor" style="color: #1a73e8; text-decoration: underline; font-weight: 500;">Global Settings</a>
</div>

<div class="panel">
    <div class="panel-header" style="cursor: pointer; display: flex; justify-content: space-between; align-items: center;" onclick="toggleConnectionSettings()">
        <div class="panel-title">Connection Settings</div>
        <span id="connectionToggleIcon" style="color: #5f6368; font-size: 12px;">▼</span>
    </div>
    <div class="panel-body" id="connectionSettingsBody">
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
        </div>
        <div id="connectionStatus" style="margin-top: 12px; padding: 8px; border-radius: 4px; display: none;"></div>
    </div>
</div>

<div id="workspaceContainer" style="display: grid; grid-template-columns: 250px 1fr; gap: 16px; height: calc(100vh - 400px); min-height: 600px; transition: grid-template-columns 0.3s ease; max-width: 100%; overflow: hidden;">
    <!-- Table Browser Sidebar -->
    <div class="panel" id="tablesPanel" style="height: 100%; display: flex; flex-direction: column; overflow: hidden; min-width: 250px; transition: min-width 0.3s ease;">
        <div class="panel-header" style="flex-shrink: 0; cursor: pointer; display: flex; justify-content: center; align-items: center; white-space: nowrap; position: relative;" onclick="toggleTablesPanel()">
            <div id="tablesPanelTitle" style="font-size: 13px; font-weight: 500; color: #5f6368; transition: all 0.3s ease;">TABLES</div>
            <span id="tablesPanelToggleIcon" style="color: #5f6368; font-size: 12px; position: absolute; top: 8px; right: 8px;">◀</span>
        </div>
        <div id="tablesPanelContent" style="padding: 12px; flex-shrink: 0;">
            <input type="text" id="tableSearch" placeholder="Search tables..."
                   onkeyup="filterTables()"
                   style="width: 100%; padding: 6px 8px; font-size: 13px;">
        </div>
        <div id="tablesListContainer" style="flex: 1; overflow-y: auto; padding: 0 12px 12px 12px;">
            <div id="tableList" style="display: flex; flex-direction: column; gap: 4px;">
                <div style="color: #5f6368; font-size: 13px; padding: 20px; text-align: center;">
                    Click "Load Tables" to view tables
                </div>
            </div>
        </div>
    </div>

    <!-- Main Content Area -->
    <div style="display: flex; flex-direction: column; gap: 16px; height: 100%; overflow: hidden; min-width: 0;">
        <!-- SQL Editor -->
        <div class="panel" style="min-height: 250px; max-height: 250px; display: flex; flex-direction: column; overflow: hidden; min-width: 0;">
            <div class="panel-header" style="flex-shrink: 0;">
                <div class="panel-title">SQL Editor</div>
            </div>
            <div class="panel-body" style="overflow: hidden;">
                <div style="display: flex; flex-direction: column; gap: 8px;">
                    <textarea id="sqlQuery"
                              style="height: 150px; width: 100%; max-width: 100%; font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
                                     font-size: 13px; resize: none; padding: 10px; box-sizing: border-box; border: 1px solid #dadce0; border-radius: 3px; overflow-x: auto;"
                              placeholder="-- Enter SQL query here&#10;SELECT * FROM TableName LIMIT 10;"></textarea>
                    <div class="button-group">
                        <button class="btn-primary" onclick="executeQuery()">Run Query</button>
                        <select id="exampleQueries" onchange="loadExampleQuery()" style="padding: 6px 10px;">
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
        <div class="panel" style="flex: 1; min-height: 250px; display: flex; flex-direction: column; overflow: hidden; min-width: 0;">
            <div class="panel-header" style="flex-shrink: 0;">
                <div class="panel-title">Results</div>
            </div>
            <div id="queryStats" style="display: none; padding: 8px 20px; background: #e8f5e9; border-bottom: 1px solid #dadce0; font-size: 13px; color: #188038; flex-shrink: 0;"></div>
            <div id="queryError" style="display: none; padding: 12px 20px; background: #fce8e6; border-bottom: 1px solid #dadce0; font-size: 13px; color: #d93025; flex-shrink: 0;"></div>
            <div class="panel-body" style="flex: 1; overflow-x: auto; overflow-y: auto; padding: 0; max-width: 100%;">
                <div id="queryResults" style="padding: 20px; color: #5f6368; font-size: 13px; min-width: min-content;">
                    Run a query to see results here
                </div>
            </div>
        </div>
    </div>
</div>
`

const SpannerJS = `
let currentConfig = {};
let allTables = [];
let selectedTable = '';
let currentTableColumns = [];
let currentTableRows = [];

// Toggle connection settings panel
function toggleConnectionSettings() {
    const body = document.getElementById('connectionSettingsBody');
    const icon = document.getElementById('connectionToggleIcon');

    if (body.style.display === 'none') {
        body.style.display = 'block';
        icon.textContent = '▼';
    } else {
        body.style.display = 'none';
        icon.textContent = '▶';
    }
}

function collapseConnectionSettings() {
    const body = document.getElementById('connectionSettingsBody');
    const icon = document.getElementById('connectionToggleIcon');
    body.style.display = 'none';
    icon.textContent = '▶';
}

// Toggle tables panel
function toggleTablesPanel() {
    const panel = document.getElementById('tablesPanel');
    const container = document.getElementById('workspaceContainer');
    const icon = document.getElementById('tablesPanelToggleIcon');
    const title = document.getElementById('tablesPanelTitle');
    const content = document.getElementById('tablesPanelContent');
    const listContainer = document.getElementById('tablesListContainer');
    const header = title.parentElement;
    const currentCols = container.style.gridTemplateColumns;

    if (currentCols === '40px 1fr' || currentCols.startsWith('40px')) {
        // Expand
        container.style.gridTemplateColumns = '250px 1fr';
        icon.textContent = '◀';
        icon.style.top = '8px';
        icon.style.right = '8px';
        icon.style.transform = 'none';
        panel.style.minWidth = '250px';
        title.style.writingMode = 'horizontal-tb';
        title.textContent = 'TABLES';
        header.style.justifyContent = 'center';
        header.style.padding = '10px 16px';
        content.style.display = 'block';
        listContainer.style.display = 'block';
    } else {
        // Collapse
        container.style.gridTemplateColumns = '40px 1fr';
        icon.textContent = '▶';
        icon.style.top = '4px';
        icon.style.right = '50%';
        icon.style.transform = 'translateX(50%)';
        panel.style.minWidth = '40px';
        title.style.writingMode = 'vertical-rl';
        title.textContent = 'TABLES';
        header.style.justifyContent = 'center';
        header.style.padding = '16px 0';
        content.style.display = 'none';
        listContainer.style.display = 'none';
    }
}

function collapseTablesPanel() {
    const panel = document.getElementById('tablesPanel');
    const container = document.getElementById('workspaceContainer');
    const icon = document.getElementById('tablesPanelToggleIcon');
    const title = document.getElementById('tablesPanelTitle');
    const content = document.getElementById('tablesPanelContent');
    const listContainer = document.getElementById('tablesListContainer');
    const header = title.parentElement;
    container.style.gridTemplateColumns = '40px 1fr';
    icon.textContent = '▶';
    icon.style.top = '4px';
    icon.style.right = '50%';
    icon.style.transform = 'translateX(50%)';
    panel.style.minWidth = '40px';
    title.style.writingMode = 'vertical-rl';
    title.textContent = 'TABLES';
    header.style.justifyContent = 'center';
    header.style.padding = '16px 0';
    content.style.display = 'none';
    listContainer.style.display = 'none';
}

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
            // Auto-load tables on successful connection
            loadTables();
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
        div.style.cssText = 'padding: 8px 12px; cursor: pointer; border-radius: 4px; font-size: 13px; transition: all 0.2s; border: 1px solid #dadce0; margin-bottom: 4px; background: white;';
        div.onmouseover = () => {
            div.style.background = '#f1f3f4';
            div.style.borderColor = '#1a73e8';
        };
        div.onmouseout = () => {
            div.style.background = selectedTable === table.name ? '#e8f0fe' : 'white';
            div.style.borderColor = selectedTable === table.name ? '#1a73e8' : '#dadce0';
        };
        div.onclick = () => selectTable(table.name);

        if (selectedTable === table.name) {
            div.style.background = '#e8f0fe';
            div.style.borderColor = '#1a73e8';
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

    // Store current table data for copying
    currentTableColumns = columns;
    currentTableRows = rows;

    let html = '<table id="resultsTable" style="width: 100%; border-collapse: collapse; font-size: 13px; background: white;">';

    // Header
    html += '<thead><tr style="background: #f8f9fa; border-bottom: 2px solid #dadce0; height: 36px;">';
    // Row number header
    html += '<th style="padding: 8px 12px; text-align: center; font-weight: 500; color: #5f6368; border: 1px solid #dadce0; font-size: 12px; width: 50px; background: #f1f3f4;">#</th>';
    columns.forEach(col => {
        html += '<th style="padding: 8px 12px; text-align: left; font-weight: 500; color: #5f6368; border: 1px solid #dadce0; font-size: 12px; text-transform: none; white-space: nowrap; max-width: 250px; overflow: hidden; text-overflow: ellipsis;">' + col + '</th>';
    });
    html += '</tr></thead>';

    // Rows
    html += '<tbody>';
    rows.forEach((row, idx) => {
        html += '<tr data-row-index="' + idx + '" style="background: white; border-bottom: 1px solid #e0e0e0; height: 32px;">';
        // Row number cell
        html += '<td class="row-number" onclick="selectRow(' + idx + ')" oncontextmenu="showRowContextMenu(event, ' + idx + ')" style="padding: 6px 12px; text-align: center; color: #5f6368; border: 1px solid #dadce0; font-size: 12px; background: #f8f9fa; cursor: pointer; user-select: none; font-weight: 500;">' + (idx + 1) + '</td>';

        columns.forEach((col, colIdx) => {
            let rawValue = row[col];
            let displayValue = rawValue;
            let tooltipValue = '';

            if (rawValue === null || rawValue === undefined) {
                displayValue = '<span style="color: #9e9e9e; font-style: italic;">NULL</span>';
                tooltipValue = 'NULL';
            } else if (typeof rawValue === 'object') {
                displayValue = JSON.stringify(rawValue);
                tooltipValue = JSON.stringify(rawValue, null, 2);
            } else {
                displayValue = String(rawValue);
                tooltipValue = String(rawValue);
            }

            // Escape for HTML attribute
            const escapedTooltip = tooltipValue.replace(/"/g, '&quot;').replace(/'/g, '&#39;');
            const escapedData = tooltipValue.replace(/"/g, '&quot;').replace(/'/g, "\\'");

            html += '<td data-row="' + idx + '" data-col="' + colIdx + '" onclick="selectCell(' + idx + ',' + colIdx + ')" oncontextmenu="showCellContextMenu(event, ' + idx + ', ' + colIdx + ')" style="padding: 6px 12px; color: #424242; border: 1px solid #e0e0e0; font-size: 13px; white-space: nowrap; max-width: 250px; overflow: hidden; text-overflow: ellipsis; cursor: pointer; transition: background 0.15s;" title="' + escapedTooltip + '" ondblclick="showCellData(\'' + escapedData + '\')">' + displayValue + '</td>';
        });
        html += '</tr>';
    });
    html += '</tbody></table>';

    resultsDiv.innerHTML = html;
}

function selectRow(rowIndex) {
    // Clear all selections
    clearSelections();

    // Select all cells in this row
    const table = document.getElementById('resultsTable');
    const row = table.querySelector('tr[data-row-index="' + rowIndex + '"]');
    if (row) {
        row.querySelectorAll('td:not(.row-number)').forEach(cell => {
            cell.style.background = '#c2daff';
        });
        // Highlight row number
        row.querySelector('.row-number').style.background = '#1a73e8';
        row.querySelector('.row-number').style.color = 'white';
    }
}

function selectCell(rowIndex, colIndex) {
    // Don't trigger if double-clicking (let double-click handler work)
    if (event && event.detail === 2) return;

    // Clear all selections
    clearSelections();

    // Select this cell
    const table = document.getElementById('resultsTable');
    const cell = table.querySelector('td[data-row="' + rowIndex + '"][data-col="' + colIndex + '"]');
    if (cell) {
        cell.style.background = '#c2daff';
        cell.style.border = '2px solid #1a73e8';
    }
}

function clearSelections() {
    const table = document.getElementById('resultsTable');
    if (!table) return;

    // Reset all cells
    table.querySelectorAll('td:not(.row-number)').forEach(cell => {
        cell.style.background = '';
        cell.style.border = '1px solid #e0e0e0';
    });

    // Reset row numbers
    table.querySelectorAll('.row-number').forEach(cell => {
        cell.style.background = '#f8f9fa';
        cell.style.color = '#5f6368';
    });
}

// Context Menu Functions
function showCellContextMenu(event, rowIndex, colIndex) {
    event.preventDefault();
    hideContextMenu();

    const menu = createContextMenu(event.clientX, event.clientY);
    const menuItem = (text, handler) => {
        const item = document.createElement('div');
        item.textContent = text;
        item.style.cssText = 'padding: 8px 16px; cursor: pointer; font-size: 13px; color: #202124;';
        item.onmouseover = () => item.style.background = '#f1f3f4';
        item.onmouseout = () => item.style.background = '';
        item.onclick = () => { handler(); hideContextMenu(); };
        return item;
    };

    menu.appendChild(menuItem('Copy Cell', () => copyCellValue(rowIndex, colIndex)));
    document.body.appendChild(menu);
}

function showRowContextMenu(event, rowIndex) {
    event.preventDefault();
    hideContextMenu();

    const menu = createContextMenu(event.clientX, event.clientY);
    const menuItem = (text, handler) => {
        const item = document.createElement('div');
        item.textContent = text;
        item.style.cssText = 'padding: 8px 16px; cursor: pointer; font-size: 13px; color: #202124;';
        item.onmouseover = () => item.style.background = '#f1f3f4';
        item.onmouseout = () => item.style.background = '';
        item.onclick = () => { handler(); hideContextMenu(); };
        return item;
    };

    menu.appendChild(menuItem('Copy Row (with headers)', () => copyRow(rowIndex, true)));
    menu.appendChild(menuItem('Copy Row (without headers)', () => copyRow(rowIndex, false)));
    document.body.appendChild(menu);
}

function createContextMenu(x, y) {
    const menu = document.createElement('div');
    menu.id = 'contextMenu';
    menu.style.cssText = 'position: fixed; background: white; border: 1px solid #dadce0; border-radius: 4px; box-shadow: 0 2px 8px rgba(0,0,0,0.15); z-index: 10000; min-width: 200px;';
    menu.style.left = x + 'px';
    menu.style.top = y + 'px';
    return menu;
}

function hideContextMenu() {
    const existing = document.getElementById('contextMenu');
    if (existing) existing.remove();
}

function copyCellValue(rowIndex, colIndex) {
    const row = currentTableRows[rowIndex];
    const column = currentTableColumns[colIndex];
    let value = row[column];

    if (value === null || value === undefined) {
        value = 'NULL';
    } else if (typeof value === 'object') {
        value = JSON.stringify(value, null, 2);
    } else {
        value = String(value);
    }

    navigator.clipboard.writeText(value).then(() => {
        showStatus('Cell value copied to clipboard');
    }).catch(err => {
        showStatus('Failed to copy: ' + err.message, true);
    });
}

function copyRow(rowIndex, withHeaders) {
    const row = currentTableRows[rowIndex];
    let text = '';

    if (withHeaders) {
        text = currentTableColumns.join('\t') + '\n';
    }

    const values = currentTableColumns.map(col => {
        let value = row[col];
        if (value === null || value === undefined) {
            return 'NULL';
        } else if (typeof value === 'object') {
            return JSON.stringify(value);
        } else {
            return String(value);
        }
    });

    text += values.join('\t');

    navigator.clipboard.writeText(text).then(() => {
        showStatus('Row copied to clipboard' + (withHeaders ? ' (with headers)' : ''));
    }).catch(err => {
        showStatus('Failed to copy: ' + err.message, true);
    });
}

// Close context menu on click outside
document.addEventListener('click', hideContextMenu);

function showCellData(data) {
    // Create modal overlay
    const overlay = document.createElement('div');
    overlay.style.cssText = 'position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); z-index: 10000; display: flex; align-items: center; justify-content: center;';

    // Create modal content
    const modal = document.createElement('div');
    modal.style.cssText = 'background: white; border-radius: 8px; padding: 24px; max-width: 600px; max-height: 80vh; overflow: auto; box-shadow: 0 4px 16px rgba(0,0,0,0.2);';

    const title = document.createElement('div');
    title.style.cssText = 'font-size: 16px; font-weight: 500; color: #202124; margin-bottom: 16px; display: flex; justify-content: space-between; align-items: center;';
    title.innerHTML = '<span>Cell Data</span><button onclick="this.parentElement.parentElement.parentElement.remove()" style="background: none; border: none; font-size: 24px; color: #5f6368; cursor: pointer; padding: 0; width: 32px; height: 32px;">&times;</button>';

    const content = document.createElement('pre');
    content.style.cssText = 'background: #f8f9fa; padding: 16px; border-radius: 4px; font-family: Monaco, monospace; font-size: 13px; white-space: pre-wrap; word-wrap: break-word; margin: 0; color: #202124;';
    content.textContent = data;

    modal.appendChild(title);
    modal.appendChild(content);
    overlay.appendChild(modal);

    // Close on overlay click
    overlay.onclick = (e) => {
        if (e.target === overlay) overlay.remove();
    };

    document.body.appendChild(overlay);
}

// Load configurations on page load
loadConfigurations();
`

func GetSpannerHTML() string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Spanner Explorer - Testing Studio</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", sans-serif;
            background: #f5f5f5;
            color: #202124;
            min-height: 100vh;
        }
        .topbar {
            background: white;
            border-bottom: 1px solid #dadce0;
            padding: 16px 24px;
            display: flex;
            align-items: center;
            gap: 16px;
        }
        .logo {
            font-size: 18px;
            font-weight: 500;
            color: #202124;
            text-decoration: none;
        }
        .back-btn {
            color: #1a73e8;
            padding: 6px 12px;
            border-radius: 4px;
            text-decoration: none;
            font-size: 14px;
            transition: background 0.2s;
        }
        .back-btn:hover { background: #f1f3f4; }
        .container { max-width: calc(100vw - 2cm); margin: 0 auto; padding: 16px; overflow: hidden; }
        .panel { background: white; border: 1px solid #dadce0; border-radius: 4px; margin-bottom: 12px; box-shadow: 0 1px 2px rgba(0,0,0,0.05); }
        .panel-header { padding: 10px 16px; border-bottom: 1px solid #e0e0e0; background: #fafafa; }
        .panel-title { font-size: 13px; font-weight: 500; color: #5f6368; }
        .panel-body { padding: 16px; }
        .form-row { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 10px; margin-bottom: 12px; }
        .form-group { display: flex; flex-direction: column; gap: 4px; }
        label { font-size: 12px; color: #5f6368; font-weight: 500; }
        input, select {
            background: white;
            border: 1px solid #dadce0;
            color: #424242;
            padding: 6px 10px;
            border-radius: 3px;
            font-size: 13px;
        }
        input:focus, select:focus { outline: none; border-color: #1a73e8; box-shadow: 0 0 0 1px #1a73e8; }
        .button-group { display: flex; gap: 8px; flex-wrap: wrap; }
        button {
            padding: 6px 14px;
            border: 1px solid #dadce0;
            border-radius: 3px;
            font-size: 13px;
            font-weight: 500;
            cursor: pointer;
            transition: all 0.2s;
            background: white;
            color: #424242;
        }
        .btn-primary { background: #1a73e8; color: white; border-color: #1a73e8; }
        .btn-primary:hover { background: #1557b0; }
        .btn-secondary { background: white; color: #5f6368; }
        .btn-secondary:hover { background: #f5f5f5; }
        .status-toast {
            position: fixed;
            top: 80px;
            right: 24px;
            padding: 12px 20px;
            border-radius: 4px;
            font-size: 14px;
            display: none;
            z-index: 1000;
            animation: slideIn 0.3s ease;
            box-shadow: 0 2px 8px rgba(0,0,0,0.15);
        }
        @keyframes slideIn {
            from { transform: translateX(400px); opacity: 0; }
            to { transform: translateX(0); opacity: 1; }
        }
        .status-toast.success { background: #188038; color: white; }
        .status-toast.error { background: #d93025; color: white; }
    </style>
</head>
<body>
    <div class="topbar">
        <a href="/" class="logo">Testing Studio</a>
        <a href="/" class="back-btn">← Back</a>
    </div>

    <div class="container">
        ` + SpannerContent + `
    </div>

    <div id="statusToast" class="status-toast"></div>

    <script>
        function showStatus(message, isError = false) {
            const toast = document.getElementById('statusToast');
            toast.textContent = message;
            toast.className = 'status-toast ' + (isError ? 'error' : 'success');
            toast.style.display = 'block';
            setTimeout(() => { toast.style.display = 'none'; }, 3000);
        }

        ` + SpannerJS + `
    </script>
</body>
</html>`
}
