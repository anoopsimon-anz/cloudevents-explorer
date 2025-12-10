package templates

const RestClient = `
	<div class="panel">
		<div class="panel-header">
			<div class="panel-title">HTTP Request Configuration</div>
		</div>
		<div class="panel-body">
			<div class="form-row">
				<div class="form-group" style="grid-column: 1 / -1;">
					<label>Environment</label>
					<select id="envSelect" onchange="loadEnvironment()">
						<option value="">No Environment</option>
					</select>
				</div>
			</div>

			<div class="form-row">
				<div class="form-group" style="flex: 0 0 120px;">
					<label>Method</label>
					<select id="httpMethod" onchange="toggleBodyField()">
						<option value="GET">GET</option>
						<option value="POST">POST</option>
						<option value="PUT">PUT</option>
						<option value="PATCH">PATCH</option>
						<option value="DELETE">DELETE</option>
						<option value="HEAD">HEAD</option>
						<option value="OPTIONS">OPTIONS</option>
					</select>
				</div>
				<div class="form-group" style="flex: 1;">
					<label>URL</label>
					<input type="text" id="requestUrl" placeholder="https://api.example.com/endpoint" />
				</div>
			</div>

			<div class="form-row">
				<div class="form-group" style="grid-column: 1 / -1;">
					<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">
						<label style="margin: 0;">Headers</label>
						<button type="button" onclick="toggleHeadersView()" class="btn-secondary" style="padding: 4px 12px; font-size: 12px;">Switch to JSON View</button>
					</div>

					<!-- Key-Value Headers View -->
					<div id="headersKeyValue" style="display: block;">
						<div id="headersList"></div>
						<button type="button" onclick="addHeaderRow()" class="btn-secondary" style="margin-top: 8px; padding: 6px 12px; font-size: 13px;">+ Add Header</button>
					</div>

					<!-- JSON Headers View -->
					<textarea id="requestHeaders" rows="4" style="display: none; font-family: 'Monaco', 'Menlo', 'Consolas', monospace; font-size: 13px; width: 100%; padding: 8px; border: 1px solid #dadce0; border-radius: 4px;">{
  "Content-Type": "application/json",
  "Accept": "application/json"
}</textarea>
				</div>
			</div>

			<div class="form-row" id="bodySection">
				<div class="form-group" style="grid-column: 1 / -1;">
					<label>Body (JSON format)</label>
					<textarea id="requestBody" rows="10" style="font-family: 'Monaco', 'Menlo', 'Consolas', monospace; font-size: 13px; width: 100%; padding: 8px; border: 1px solid #dadce0; border-radius: 4px;">{
  "key": "value"
}</textarea>
				</div>
			</div>

			<div class="form-row">
				<div class="form-group" style="grid-column: 1 / -1;">
					<button type="button" onclick="toggleCertificateSection()" class="btn-secondary" style="width: 100%; text-align: left; display: flex; justify-content: space-between; align-items: center;">
						<span>🔒 Client Certificates (Optional)</span>
						<span id="certToggleIcon">▼</span>
					</button>
				</div>
			</div>

			<div id="certificateSection" style="display: none;">
				<div class="form-row">
					<div class="form-group" style="grid-column: 1 / -1;">
						<label>TLS Certificate (PEM format)</label>
						<textarea id="tlsCert" rows="6" placeholder="-----BEGIN CERTIFICATE-----
MIIDXTCCAkWgAwIBAgIJAKZ...
-----END CERTIFICATE-----" style="font-family: 'Monaco', 'Menlo', 'Consolas', monospace; font-size: 13px; width: 100%; padding: 8px; border: 1px solid #dadce0; border-radius: 4px;"></textarea>
					</div>
				</div>

				<div class="form-row">
					<div class="form-group" style="grid-column: 1 / -1;">
						<label>TLS Private Key (PEM format)</label>
						<textarea id="tlsKey" rows="6" placeholder="-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0B...
-----END PRIVATE KEY-----" style="font-family: 'Monaco', 'Menlo', 'Consolas', monospace; font-size: 13px; width: 100%; padding: 8px; border: 1px solid #dadce0; border-radius: 4px;"></textarea>
					</div>
				</div>
			</div>

			<div class="button-group">
				<button class="btn-primary" onclick="sendRequest()">Send Request</button>
				<button class="btn-secondary" onclick="clearRequest()">Clear</button>
				<button class="btn-secondary" onclick="saveEnvironment()">Save as Environment</button>
				<button class="btn-secondary" onclick="manageEnvironments()">Manage Environments</button>
			</div>
		</div>
	</div>

	<div class="panel">
		<div class="panel-header">
			<div class="panel-title">Response</div>
		</div>
		<div class="stats-bar" id="responseStats" style="display: none;">
			<div class="stat">
				<span class="stat-label">Status:</span>
				<span class="stat-value" id="responseStatus">-</span>
			</div>
			<div class="stat">
				<span class="stat-label">Time:</span>
				<span class="stat-value" id="responseTime">-</span>
			</div>
			<div class="stat">
				<span class="stat-label">Size:</span>
				<span class="stat-value" id="responseSize">-</span>
			</div>
		</div>
		<div class="panel-body">
			<div id="responseContainer"></div>
		</div>
	</div>

	<!-- Environment Management Modal -->
	<div id="envModal" style="display: none; position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0,0,0,0.5); z-index: 2000; align-items: center; justify-content: center;">
		<div style="background: white; border-radius: 8px; max-width: 600px; width: 90%; max-height: 80vh; overflow-y: auto; padding: 24px;">
			<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
				<h2 style="font-size: 20px; font-weight: 500; margin: 0;">Manage Environments</h2>
				<button onclick="closeEnvModal()" style="background: none; border: none; font-size: 24px; cursor: pointer; color: #5f6368;">&times;</button>
			</div>

			<div style="margin-bottom: 20px;">
				<label style="font-size: 13px; color: #5f6368; font-weight: 500; display: block; margin-bottom: 6px;">Environment Name</label>
				<input type="text" id="newEnvName" placeholder="Development" style="width: 100%; padding: 8px 12px; border: 1px solid #dadce0; border-radius: 4px; font-size: 14px; margin-bottom: 8px;" />
				<label style="font-size: 13px; color: #5f6368; font-weight: 500; display: block; margin-bottom: 6px;">Variables (JSON format)</label>
				<textarea id="newEnvVars" rows="8" placeholder='{
  "BASE_URL": "https://api-dev.example.com",
  "API_KEY": "dev-key-123",
  "TIMEOUT": "30"
}' style="font-family: 'Monaco', 'Menlo', 'Consolas', monospace; font-size: 13px; width: 100%; padding: 8px; border: 1px solid #dadce0; border-radius: 4px;"></textarea>
				<button class="btn-primary" onclick="addEnvironment()" style="margin-top: 8px;">Add Environment</button>
			</div>

			<div id="envList" style="margin-top: 20px;">
				<!-- Environment list will be populated here -->
			</div>
		</div>
	</div>
`

const RestClientJS = `
let environments = JSON.parse(localStorage.getItem('restClientEnvironments') || '{}');
let currentEnv = '';
let headersViewMode = 'keyvalue'; // 'keyvalue' or 'json'
let headerRows = [];

// Header row management
function addHeaderRow(key = '', value = '', enabled = true) {
	const id = 'header_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
	headerRows.push({ id, key, value, enabled });
	renderHeaderRows();
}

function removeHeaderRow(id) {
	headerRows = headerRows.filter(h => h.id !== id);
	renderHeaderRows();
}

function updateHeaderRow(id, field, value) {
	const header = headerRows.find(h => h.id === id);
	if (header) {
		header[field] = value;
	}
}

function renderHeaderRows() {
	const container = document.getElementById('headersList');
	let html = '';

	headerRows.forEach(header => {
		html += '<div style="display: flex; gap: 8px; margin-bottom: 8px; align-items: center;">';
		html += '<input type="checkbox" ' + (header.enabled ? 'checked' : '') + ' onchange="updateHeaderRow(\'' + header.id + '\', \'enabled\', this.checked)" style="width: 20px; height: 20px; cursor: pointer;" />';
		html += '<input type="text" value="' + header.key + '" placeholder="Header Key" onchange="updateHeaderRow(\'' + header.id + '\', \'key\', this.value)" style="flex: 1; padding: 8px; border: 1px solid #dadce0; border-radius: 4px; font-size: 14px;" />';
		html += '<input type="text" value="' + header.value + '" placeholder="Header Value" onchange="updateHeaderRow(\'' + header.id + '\', \'value\', this.value)" style="flex: 1; padding: 8px; border: 1px solid #dadce0; border-radius: 4px; font-size: 14px;" />';
		html += '<button onclick="removeHeaderRow(\'' + header.id + '\')" class="btn-secondary" style="padding: 8px 12px; font-size: 12px;">×</button>';
		html += '</div>';
	});

	container.innerHTML = html;
}

function toggleHeadersView() {
	const keyValueDiv = document.getElementById('headersKeyValue');
	const jsonTextarea = document.getElementById('requestHeaders');
	const toggleButton = event.target;

	if (headersViewMode === 'keyvalue') {
		// Switch to JSON view
		// First, sync key-value to JSON
		const headers = {};
		headerRows.filter(h => h.enabled && h.key).forEach(h => {
			headers[h.key] = h.value;
		});
		jsonTextarea.value = JSON.stringify(headers, null, 2);

		keyValueDiv.style.display = 'none';
		jsonTextarea.style.display = 'block';
		toggleButton.textContent = 'Switch to Key-Value View';
		headersViewMode = 'json';
	} else {
		// Switch to key-value view
		// First, sync JSON to key-value
		try {
			const headers = JSON.parse(jsonTextarea.value || '{}');
			headerRows = [];
			Object.entries(headers).forEach(([key, value]) => {
				addHeaderRow(key, value, true);
			});
		} catch (e) {
			showStatus('Invalid JSON in headers, keeping current view', true);
			return;
		}

		keyValueDiv.style.display = 'block';
		jsonTextarea.style.display = 'none';
		toggleButton.textContent = 'Switch to JSON View';
		headersViewMode = 'keyvalue';
	}
}

function toggleCertificateSection() {
	const section = document.getElementById('certificateSection');
	const icon = document.getElementById('certToggleIcon');

	if (section.style.display === 'none') {
		section.style.display = 'block';
		icon.textContent = '▲';
	} else {
		section.style.display = 'none';
		icon.textContent = '▼';
	}
}

function toggleBodyField() {
	const method = document.getElementById('httpMethod').value;
	const bodySection = document.getElementById('bodySection');
	const bodyTextarea = document.getElementById('requestBody');

	// Methods that support request body
	const methodsWithBody = ['POST', 'PUT', 'PATCH'];

	if (methodsWithBody.includes(method)) {
		bodySection.style.display = 'flex';
		bodyTextarea.disabled = false;
	} else {
		bodySection.style.display = 'none';
		bodyTextarea.disabled = true;
	}
}

function getHeadersObject() {
	if (headersViewMode === 'keyvalue') {
		const headers = {};
		headerRows.filter(h => h.enabled && h.key).forEach(h => {
			headers[h.key] = h.value;
		});
		return headers;
	} else {
		try {
			return JSON.parse(document.getElementById('requestHeaders').value || '{}');
		} catch (e) {
			throw new Error('Invalid JSON in headers');
		}
	}
}

function loadEnvironments() {
	const select = document.getElementById('envSelect');
	select.innerHTML = '<option value="">No Environment</option>';
	Object.keys(environments).forEach(name => {
		const option = document.createElement('option');
		option.value = name;
		option.textContent = name;
		select.appendChild(option);
	});
}

function loadEnvironment() {
	const select = document.getElementById('envSelect');
	currentEnv = select.value;
}

function replaceVariables(text) {
	if (!currentEnv || !environments[currentEnv]) return text;

	const vars = environments[currentEnv];
	let result = text;

	Object.keys(vars).forEach(key => {
		const regex = new RegExp('\\{\\{' + key + '\\}\\}', 'g');
		result = result.replace(regex, vars[key]);
	});

	return result;
}

async function sendRequest() {
	const method = document.getElementById('httpMethod').value;
	const url = replaceVariables(document.getElementById('requestUrl').value.trim());
	const bodyText = replaceVariables(document.getElementById('requestBody').value.trim());
	const tlsCert = document.getElementById('tlsCert').value.trim();
	const tlsKey = document.getElementById('tlsKey').value.trim();

	if (!url) {
		showStatus('Please enter a URL', true);
		return;
	}

	let headers = {};
	try {
		headers = getHeadersObject();
		// Apply variable replacement to header values
		Object.keys(headers).forEach(key => {
			headers[key] = replaceVariables(headers[key]);
		});
	} catch (e) {
		showStatus('Invalid headers: ' + e.message, true);
		return;
	}

	let body = null;
	if (['POST', 'PUT', 'PATCH'].includes(method) && bodyText) {
		try {
			body = JSON.parse(bodyText);
		} catch (e) {
			showStatus('Invalid JSON in body', true);
			return;
		}
	}

	const requestData = {
		method: method,
		url: url,
		headers: headers,
		body: body,
		tlsCert: tlsCert || null,
		tlsKey: tlsKey || null
	};

	const startTime = Date.now();

	try {
		const response = await fetch('/api/rest/send', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(requestData)
		});

		const endTime = Date.now();
		const duration = endTime - startTime;

		const result = await response.json();

		if (!response.ok) {
			showStatus('Request failed: ' + (result.error || 'Unknown error'), true);
			displayResponse({
				error: result.error || 'Unknown error',
				statusCode: result.statusCode || 0,
				duration: duration
			}, true);
			return;
		}

		showStatus('Request completed successfully');
		displayResponse({
			statusCode: result.statusCode,
			headers: result.headers,
			body: result.body,
			duration: duration,
			size: JSON.stringify(result.body).length
		}, false);

	} catch (error) {
		const endTime = Date.now();
		const duration = endTime - startTime;

		showStatus('Request failed: ' + error.message, true);
		displayResponse({
			error: error.message,
			duration: duration
		}, true);
	}
}

function displayResponse(data, isError) {
	const container = document.getElementById('responseContainer');
	const statsBar = document.getElementById('responseStats');

	statsBar.style.display = 'flex';

	if (isError) {
		document.getElementById('responseStatus').textContent = data.statusCode || 'Error';
		document.getElementById('responseStatus').style.color = '#d93025';
		document.getElementById('responseTime').textContent = data.duration + 'ms';
		document.getElementById('responseSize').textContent = '-';

		container.innerHTML = '<div class="json-viewer"><pre>' +
			JSON.stringify({ error: data.error }, null, 2) +
			'</pre></div>';
		return;
	}

	const statusColor = data.statusCode >= 200 && data.statusCode < 300 ? '#188038' :
	                     data.statusCode >= 400 ? '#d93025' : '#5f6368';

	document.getElementById('responseStatus').textContent = data.statusCode;
	document.getElementById('responseStatus').style.color = statusColor;
	document.getElementById('responseTime').textContent = data.duration + 'ms';
	document.getElementById('responseSize').textContent = formatBytes(data.size);

	let html = '';

	if (data.headers) {
		html += '<div style="margin-bottom: 16px;">';
		html += '<div style="font-weight: 500; margin-bottom: 8px; color: #202124;">Response Headers</div>';
		html += '<div class="json-viewer"><pre>' + syntaxHighlightJSON(data.headers) + '</pre></div>';
		html += '</div>';
	}

	if (data.body) {
		html += '<div>';
		html += '<div style="font-weight: 500; margin-bottom: 8px; color: #202124; display: flex; justify-content: space-between; align-items: center;">';
		html += '<span>Response Body</span>';
		html += '<button onclick="copyResponseBody()" class="btn-secondary" style="padding: 4px 12px; font-size: 12px;">Copy</button>';
		html += '</div>';
		html += '<div class="json-viewer"><pre id="responseBody">' + syntaxHighlightJSON(data.body) + '</pre></div>';
		html += '</div>';
	}

	container.innerHTML = html;
}

function copyResponseBody() {
	const bodyElement = document.getElementById('responseBody');
	const text = bodyElement.textContent;

	const textarea = document.createElement('textarea');
	textarea.value = text;
	document.body.appendChild(textarea);
	textarea.select();
	document.execCommand('copy');
	document.body.removeChild(textarea);

	showStatus('Response body copied to clipboard!');
}

function formatBytes(bytes) {
	if (bytes === 0) return '0 Bytes';
	const k = 1024;
	const sizes = ['Bytes', 'KB', 'MB'];
	const i = Math.floor(Math.log(bytes) / Math.log(k));
	return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

function clearRequest() {
	document.getElementById('requestUrl').value = '';
	document.getElementById('requestBody').value = '{\n  "key": "value"\n}';
	document.getElementById('tlsCert').value = '';
	document.getElementById('tlsKey').value = '';
	document.getElementById('httpMethod').value = 'GET';
	document.getElementById('envSelect').value = '';
	currentEnv = '';

	// Reset headers to defaults
	headerRows = [];
	addHeaderRow('Content-Type', 'application/json', true);
	addHeaderRow('Accept', 'application/json', true);

	// Reset to key-value view
	if (headersViewMode === 'json') {
		toggleHeadersView();
	}

	// Hide certificate section
	document.getElementById('certificateSection').style.display = 'none';
	document.getElementById('certToggleIcon').textContent = '▼';

	document.getElementById('responseContainer').innerHTML = '<div class="empty-state"><div>No response yet. Send a request to get started.</div></div>';
	document.getElementById('responseStats').style.display = 'none';
}

function saveEnvironment() {
	const name = prompt('Environment name:');
	if (!name) return;

	const url = document.getElementById('requestUrl').value.trim();
	const headers = document.getElementById('requestHeaders').value.trim();

	try {
		const varsToSave = {};

		// Extract variables from URL
		const urlMatches = url.match(/\{\{(\w+)\}\}/g);
		if (urlMatches) {
			urlMatches.forEach(match => {
				const varName = match.replace(/\{\{|\}\}/g, '');
				const value = prompt('Value for ' + varName + ':');
				if (value !== null) {
					varsToSave[varName] = value;
				}
			});
		}

		// Extract variables from headers
		const headerMatches = headers.match(/\{\{(\w+)\}\}/g);
		if (headerMatches) {
			headerMatches.forEach(match => {
				const varName = match.replace(/\{\{|\}\}/g, '');
				if (!varsToSave[varName]) {
					const value = prompt('Value for ' + varName + ':');
					if (value !== null) {
						varsToSave[varName] = value;
					}
				}
			});
		}

		environments[name] = varsToSave;
		localStorage.setItem('restClientEnvironments', JSON.stringify(environments));
		loadEnvironments();
		showStatus('Environment "' + name + '" saved');
	} catch (e) {
		showStatus('Failed to save environment: ' + e.message, true);
	}
}

function manageEnvironments() {
	const modal = document.getElementById('envModal');
	modal.style.display = 'flex';
	renderEnvironmentList();
}

function closeEnvModal() {
	const modal = document.getElementById('envModal');
	modal.style.display = 'none';
}

function renderEnvironmentList() {
	const container = document.getElementById('envList');
	const envNames = Object.keys(environments);

	if (envNames.length === 0) {
		container.innerHTML = '<div style="text-align: center; color: #5f6368; padding: 20px;">No environments saved yet</div>';
		return;
	}

	let html = '<div style="display: flex; flex-direction: column; gap: 12px;">';
	envNames.forEach(name => {
		html += '<div style="border: 1px solid #dadce0; border-radius: 4px; padding: 12px;">';
		html += '<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">';
		html += '<span style="font-weight: 500; color: #202124;">' + name + '</span>';
		html += '<button onclick="deleteEnvironment(\'' + name + '\')" class="btn-danger" style="padding: 4px 12px; font-size: 12px;">Delete</button>';
		html += '</div>';
		html += '<pre style="margin: 0; font-size: 12px; color: #5f6368; background: #f8f9fa; padding: 8px; border-radius: 4px; overflow-x: auto;">' +
			JSON.stringify(environments[name], null, 2) + '</pre>';
		html += '</div>';
	});
	html += '</div>';

	container.innerHTML = html;
}

function addEnvironment() {
	const name = document.getElementById('newEnvName').value.trim();
	const varsText = document.getElementById('newEnvVars').value.trim();

	if (!name) {
		showStatus('Please enter an environment name', true);
		return;
	}

	try {
		const vars = varsText ? JSON.parse(varsText) : {};
		environments[name] = vars;
		localStorage.setItem('restClientEnvironments', JSON.stringify(environments));

		document.getElementById('newEnvName').value = '';
		document.getElementById('newEnvVars').value = '';

		loadEnvironments();
		renderEnvironmentList();
		showStatus('Environment "' + name + '" added');
	} catch (e) {
		showStatus('Invalid JSON format: ' + e.message, true);
	}
}

function deleteEnvironment(name) {
	if (!confirm('Delete environment "' + name + '"?')) return;

	delete environments[name];
	localStorage.setItem('restClientEnvironments', JSON.stringify(environments));

	loadEnvironments();
	renderEnvironmentList();
	showStatus('Environment "' + name + '" deleted');

	if (currentEnv === name) {
		currentEnv = '';
		document.getElementById('envSelect').value = '';
	}
}

// Initialize
loadEnvironments();

// Set default headers
addHeaderRow('Content-Type', 'application/json', true);
addHeaderRow('Accept', 'application/json', true);

// Hide body field for GET requests (default)
toggleBodyField();

document.getElementById('responseContainer').innerHTML = '<div class="empty-state"><div>No response yet. Send a request to get started.</div></div>';
`
