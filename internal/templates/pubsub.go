package templates

const PubSubContent = `<div class="panel">
    <div class="panel-header">
        <div class="panel-title">PubSub Connection Settings</div>
    </div>
    <div class="panel-body">
        <div class="form-row">
            <div class="form-group">
                <label>Saved Configurations</label>
                <select id="configSelect" onchange="loadSelectedConfig()">
                    <option value="">-- Select Configuration --</option>
                </select>
            </div>
            <div class="form-group">
                <label>Configuration Name</label>
                <input type="text" id="configName" placeholder="e.g., TMS Local">
            </div>
            <div class="form-group">
                <label>Emulator Host</label>
                <input type="text" id="emulatorHost" placeholder="localhost:8086">
            </div>
            <div class="form-group">
                <label>Project ID</label>
                <input type="text" id="projectId" placeholder="project-id">
            </div>
            <div class="form-group">
                <label>Subscription ID</label>
                <input type="text" id="subscriptionId" placeholder="subscription-name">
            </div>
            <div class="form-group">
                <label>Max Messages</label>
                <input type="number" id="maxMessages" value="20" min="1" max="100">
            </div>
        </div>
        <div class="button-group">
            <button class="btn-primary" onclick="pullMessages()">Pull Messages</button>
            <button class="btn-secondary" onclick="saveConfiguration()">Save Config</button>
            <button class="btn-secondary" onclick="refreshConfigs()">Refresh</button>
            <button class="btn-danger" onclick="clearAllMessages()">Clear All</button>
        </div>
    </div>
</div>`

const PubSubJS = `async function refreshConfigs() {
    const response = await fetch('/api/configs');
    const data = await response.json();
    const select = document.getElementById('configSelect');
    select.innerHTML = '<option value="">-- Select Configuration --</option>';
    data.pubsubConfigs.forEach((config, index) => {
        const option = document.createElement('option');
        option.value = index;
        option.textContent = config.name;
        select.appendChild(option);
    });
}

function loadSelectedConfig() {
    const select = document.getElementById('configSelect');
    if (select.value === '') return;
    fetch('/api/configs')
        .then(res => res.json())
        .then(data => {
            const config = data.pubsubConfigs[parseInt(select.value)];
            document.getElementById('configName').value = config.name;
            document.getElementById('emulatorHost').value = config.emulatorHost;
            document.getElementById('projectId').value = config.projectId;
            document.getElementById('subscriptionId').value = config.subscriptionId;
        });
}

async function saveConfiguration() {
    const config = {
        name: document.getElementById('configName').value,
        emulatorHost: document.getElementById('emulatorHost').value,
        projectId: document.getElementById('projectId').value,
        subscriptionId: document.getElementById('subscriptionId').value
    };
    if (!config.name || !config.emulatorHost || !config.projectId || !config.subscriptionId) {
        showStatus('Please fill in all configuration fields', true);
        return;
    }
    const response = await fetch('/api/pubsub/configs', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(config)
    });
    if (response.ok) {
        showStatus('Configuration saved successfully');
        refreshConfigs();
    } else {
        showStatus('Failed to save configuration', true);
    }
}

async function pullMessages() {
    const messagesDiv = document.getElementById('messages');
    messagesDiv.innerHTML = '<div class="loading"><div class="spinner"></div>Pulling messages...</div>';
    const params = {
        emulatorHost: document.getElementById('emulatorHost').value,
        projectId: document.getElementById('projectId').value,
        subscriptionId: document.getElementById('subscriptionId').value,
        maxMessages: parseInt(document.getElementById('maxMessages').value)
    };
    if (!params.emulatorHost || !params.projectId || !params.subscriptionId) {
        messagesDiv.innerHTML = '<div class="empty-state"><div>Please fill in all connection fields</div></div>';
        return;
    }
    try {
        const response = await fetch('/api/pubsub/pull', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(params)
        });
        const data = await response.json();
        if (!response.ok) throw new Error(data.error || 'Failed to pull messages');
        messagesData = data.messages.concat(messagesData);
        renderMessages();
        showStatus('Pulled ' + data.messages.length + ' new message(s)');
    } catch (error) {
        messagesDiv.innerHTML = '<div class="empty-state"><div>Error: ' + error.message + '</div></div>';
        showStatus('Failed to pull messages: ' + error.message, true);
    }
}

refreshConfigs();`
