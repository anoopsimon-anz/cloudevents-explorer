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
        <div class="panel" style="min-height: 250px; max-height: 250px; display: flex; flex-direction: column;">
            <div class="panel-header" style="flex-shrink: 0;">
                <div class="panel-title">SQL Editor</div>
            </div>
            <div class="panel-body">
                <div style="display: flex; flex-direction: column; gap: 8px;">
                    <textarea id="sqlQuery"
                              style="height: 150px; width: 100%; font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
                                     font-size: 13px; resize: none; padding: 12px; box-sizing: border-box;"
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

    let html = '<table style="width: 100%; border-collapse: collapse; font-size: 13px;">';

    // Header
    html += '<thead><tr style="background: #f8f9fa;">';
    columns.forEach(col => {
        html += '<th style="padding: 12px; text-align: left; font-weight: 600; color: #5f6368; border: 1px solid #dadce0;">' + col + '</th>';
    });
    html += '</tr></thead>';

    // Rows
    html += '<tbody>';
    rows.forEach((row, idx) => {
        const bgColor = idx % 2 === 0 ? 'white' : '#f8f9fa';
        html += '<tr style="background: ' + bgColor + ';">';
        columns.forEach(col => {
            let value = row[col];
            if (value === null || value === undefined) {
                value = '<span style="color: #5f6368; font-style: italic;">NULL</span>';
            } else if (typeof value === 'object') {
                value = JSON.stringify(value);
            }
            html += '<td style="padding: 10px; color: #202124; border: 1px solid #dadce0;">' + value + '</td>';
        });
        html += '</tr>';
    });
    html += '</tbody></table>';

    resultsDiv.innerHTML = html;
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
        .container { max-width: 1400px; margin: 0 auto; padding: 20px; }
        .panel { background: white; border: 1px solid #dadce0; border-radius: 8px; margin-bottom: 16px; }
        .panel-header { padding: 16px 20px; border-bottom: 1px solid #dadce0; }
        .panel-title { font-size: 14px; font-weight: 500; color: #5f6368; text-transform: uppercase; letter-spacing: 0.5px; }
        .panel-body { padding: 20px; }
        .form-row { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 12px; margin-bottom: 16px; }
        .form-group { display: flex; flex-direction: column; gap: 6px; }
        label { font-size: 13px; color: #5f6368; font-weight: 500; }
        input, select {
            background: white;
            border: 1px solid #dadce0;
            color: #202124;
            padding: 8px 12px;
            border-radius: 4px;
            font-size: 14px;
        }
        input:focus, select:focus { outline: none; border-color: #1a73e8; box-shadow: 0 0 0 1px #1a73e8; }
        .button-group { display: flex; gap: 8px; flex-wrap: wrap; }
        button {
            padding: 8px 16px;
            border: 1px solid #dadce0;
            border-radius: 4px;
            font-size: 14px;
            font-weight: 500;
            cursor: pointer;
            transition: all 0.2s;
            background: white;
            color: #202124;
        }
        .btn-primary { background: #1a73e8; color: white; border-color: #1a73e8; }
        .btn-primary:hover { background: #1765cc; }
        .btn-secondary { background: white; color: #5f6368; }
        .btn-secondary:hover { background: #f1f3f4; }
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
