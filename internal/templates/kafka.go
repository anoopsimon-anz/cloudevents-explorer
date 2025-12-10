package templates

const KafkaContent = `<div class="panel">
    <div class="panel-header">
        <div class="panel-title">Kafka Connection Settings</div>
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
                <input type="text" id="configName" placeholder="e.g., TMS Unica Local">
            </div>
            <div class="form-group">
                <label>Brokers</label>
                <input type="text" id="brokers" placeholder="localhost:19092">
            </div>
            <div class="form-group">
                <label>Topic</label>
                <input type="text" id="topic" placeholder="unica.marketing.response.events">
            </div>
            <div class="form-group">
                <label>Consumer Group</label>
                <input type="text" id="consumerGroup" placeholder="cloudevents-explorer">
            </div>
            <div class="form-group">
                <label>Schema Registry (optional)</label>
                <input type="text" id="schemaRegistry" placeholder="http://localhost:18081">
            </div>
            <div class="form-group">
                <label>Max Messages</label>
                <input type="number" id="maxMessages" value="20" min="1" max="100">
            </div>
        </div>
        <div class="button-group">
            <button class="btn-primary" onclick="pullMessages()">Pull Messages</button>
            <button class="btn-primary" onclick="openPublishModal()" style="background: #188038; border-color: #188038;">Publish Message</button>
            <button class="btn-secondary" onclick="saveConfiguration()">Save Config</button>
            <button class="btn-secondary" onclick="refreshConfigs()">Refresh</button>
            <button class="btn-danger" onclick="clearAllMessages()">Clear All</button>
        </div>
    </div>
</div>

<div id="publishModal" style="display: none; position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); z-index: 1000; align-items: center; justify-content: center;">
    <div style="background: white; border-radius: 12px; max-width: 1000px; width: 90%; height: 85vh; overflow: hidden; display: flex; flex-direction: column; box-shadow: 0 8px 32px rgba(0,0,0,0.2);">
        <div style="padding: 24px; border-bottom: 1px solid #e8eaed; display: flex; justify-content: space-between; align-items: center;">
            <h2 style="font-size: 20px; font-weight: 500; color: #202124;">Publish Kafka Message</h2>
            <button onclick="closePublishModal()" style="background: none; border: none; font-size: 28px; cursor: pointer; color: #5f6368; line-height: 1; padding: 0; width: 32px; height: 32px;">&times;</button>
        </div>
        <div style="flex: 1; padding: 24px; display: flex; flex-direction: column; overflow: hidden;">
            <div style="margin-bottom: 16px; padding: 12px 16px; background: #e8f0fe; border-left: 4px solid #1a73e8; border-radius: 4px; font-size: 13px; color: #1967d2;">
                <div style="margin-bottom: 4px;"><strong>Topic:</strong> <span id="publishTopic" style="font-family: monospace;">-</span></div>
                <div><strong>Schema Registry:</strong> <span id="publishSchema" style="font-family: monospace;">-</span></div>
            </div>
            <label style="font-size: 14px; color: #202124; font-weight: 500; margin-bottom: 10px;">Message JSON:</label>
            <textarea id="publishMessageJson" style="flex: 1; font-family: 'Monaco', 'Menlo', 'Consolas', monospace; font-size: 14px; line-height: 1.6; border: 1px solid #dadce0; border-radius: 6px; padding: 16px; resize: none; background: #f8f9fa;" placeholder='Paste your JSON message here...

Example:
{
  "header": {
    "messageId": "123",
    "timestamp": "2024-01-01T00:00:00Z"
  },
  "marketingResponse": {
    "customerId": "456",
    "responseCode": "SUCCESS"
  }
}'></textarea>
            <div style="display: flex; gap: 12px; margin-top: 16px;">
                <button onclick="publishMessage()" style="flex: 1; background: #188038; color: white; border: none; padding: 12px 24px; border-radius: 6px; cursor: pointer; font-weight: 500; font-size: 14px; transition: background 0.2s;" onmouseover="this.style.background='#137333'" onmouseout="this.style.background='#188038'">Publish to Kafka</button>
                <button onclick="closePublishModal()" style="background: #f1f3f4; color: #5f6368; border: none; padding: 12px 24px; border-radius: 6px; cursor: pointer; font-weight: 500; font-size: 14px; transition: background 0.2s;" onmouseover="this.style.background='#e8eaed'" onmouseout="this.style.background='#f1f3f4'">Cancel</button>
            </div>
        </div>
    </div>
</div>`

const KafkaJS = `async function refreshConfigs() {
    const response = await fetch('/api/configs');
    const data = await response.json();
    const select = document.getElementById('configSelect');
    select.innerHTML = '<option value="">-- Select Configuration --</option>';
    data.kafkaConfigs.forEach((config, index) => {
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
            const config = data.kafkaConfigs[parseInt(select.value)];
            document.getElementById('configName').value = config.name;
            document.getElementById('brokers').value = config.brokers;
            document.getElementById('topic').value = config.topic;
            document.getElementById('consumerGroup').value = config.consumerGroup;
            document.getElementById('schemaRegistry').value = config.schemaRegistry || '';
        });
}

async function saveConfiguration() {
    const config = {
        name: document.getElementById('configName').value,
        brokers: document.getElementById('brokers').value,
        topic: document.getElementById('topic').value,
        consumerGroup: document.getElementById('consumerGroup').value,
        schemaRegistry: document.getElementById('schemaRegistry').value
    };
    if (!config.name || !config.brokers || !config.topic || !config.consumerGroup) {
        showStatus('Please fill in required fields', true);
        return;
    }
    const response = await fetch('/api/kafka/configs', {
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
    messagesDiv.innerHTML = '<div class="loading"><div class="spinner"></div>Pulling messages from Kafka...</div>';
    const baseConsumerGroup = document.getElementById('consumerGroup').value;
    const uniqueConsumerGroup = baseConsumerGroup + '-' + Date.now();
    const params = {
        brokers: document.getElementById('brokers').value,
        topic: document.getElementById('topic').value,
        consumerGroup: uniqueConsumerGroup,
        schemaRegistry: document.getElementById('schemaRegistry').value,
        maxMessages: parseInt(document.getElementById('maxMessages').value)
    };
    if (!params.brokers || !params.topic || !baseConsumerGroup) {
        messagesDiv.innerHTML = '<div class="empty-state"><div>Please fill in all required fields</div></div>';
        return;
    }
    try {
        const response = await fetch('/api/kafka/pull', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(params)
        });
        const data = await response.json();
        if (!response.ok) throw new Error(data.error || 'Failed to pull messages');
        messagesData = data.messages.concat(messagesData);
        renderMessages();
        showStatus('Pulled ' + data.messages.length + ' new message(s) from Kafka');
    } catch (error) {
        messagesDiv.innerHTML = '<div class="empty-state"><div>Error: ' + error.message + '</div></div>';
        showStatus('Failed to pull messages: ' + error.message, true);
    }
}

function openPublishModal() {
    const topic = document.getElementById('topic').value;
    const schemaRegistry = document.getElementById('schemaRegistry').value;

    if (!topic) {
        showStatus('Please configure topic first', true);
        return;
    }

    document.getElementById('publishTopic').textContent = topic;
    document.getElementById('publishSchema').textContent = schemaRegistry || 'Not configured';
    document.getElementById('publishModal').style.display = 'flex';
}

function closePublishModal() {
    document.getElementById('publishModal').style.display = 'none';
    document.getElementById('publishMessageJson').value = '';
}

async function publishMessage() {
    const jsonInput = document.getElementById('publishMessageJson').value.trim();

    if (!jsonInput) {
        showStatus('Please enter message JSON', true);
        return;
    }

    let messageData;
    try {
        messageData = JSON.parse(jsonInput);
    } catch (e) {
        showStatus('Invalid JSON: ' + e.message, true);
        return;
    }

    const params = {
        brokers: document.getElementById('brokers').value,
        topic: document.getElementById('topic').value,
        schemaRegistry: document.getElementById('schemaRegistry').value,
        message: messageData
    };

    if (!params.brokers || !params.topic) {
        showStatus('Please configure brokers and topic first', true);
        return;
    }

    try {
        const response = await fetch('/api/kafka/publish', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(params)
        });
        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Failed to publish message');
        }

        showStatus('Message published successfully!');
        closePublishModal();

        // Auto-pull to show the newly published message
        setTimeout(function() {
            pullMessages();
        }, 500);
    } catch (error) {
        showStatus('Failed to publish: ' + error.message, true);
    }
}

document.getElementById('publishModal')?.addEventListener('click', function(e) {
    if (e.target === this) {
        closePublishModal();
    }
});

refreshConfigs();`
