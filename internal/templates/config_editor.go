package templates

const ConfigEditor = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Config Editor - Testing Studio</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Google Sans', 'Product Sans', -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", sans-serif;
            background: #f5f5f5;
            color: #202124;
        }
        .header {
            background: white;
            border-bottom: 1px solid #dadce0;
            padding: 16px 24px;
            display: flex;
            align-items: center;
            justify-content: space-between;
        }
        .header h1 {
            font-size: 20px;
            font-weight: 500;
            color: #202124;
        }
        .back-link {
            color: #1a73e8;
            text-decoration: none;
            font-size: 14px;
        }
        .back-link:hover {
            text-decoration: underline;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 24px;
        }
        .section {
            background: white;
            border: 1px solid #dadce0;
            border-radius: 8px;
            padding: 24px;
            margin-bottom: 24px;
        }
        .section-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 20px;
        }
        .section-title {
            font-size: 18px;
            font-weight: 500;
            color: #202124;
        }
        .add-btn {
            background: #1a73e8;
            color: white;
            border: none;
            padding: 8px 16px;
            border-radius: 4px;
            font-size: 14px;
            font-weight: 500;
            cursor: pointer;
            transition: background 0.2s;
        }
        .add-btn:hover {
            background: #1765cc;
        }
        .config-item {
            border: 1px solid #e8eaed;
            border-radius: 8px;
            padding: 16px;
            margin-bottom: 16px;
            background: #f8f9fa;
        }
        .config-item-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 12px;
        }
        .config-name {
            font-size: 16px;
            font-weight: 500;
            color: #202124;
        }
        .config-actions {
            display: flex;
            gap: 8px;
        }
        .edit-btn, .delete-btn, .save-btn, .cancel-btn {
            padding: 6px 12px;
            border-radius: 4px;
            font-size: 13px;
            font-weight: 500;
            cursor: pointer;
            border: none;
            transition: all 0.2s;
        }
        .edit-btn {
            background: #e8f0fe;
            color: #1967d2;
        }
        .edit-btn:hover {
            background: #d2e3fc;
        }
        .delete-btn {
            background: #fce8e6;
            color: #d93025;
        }
        .delete-btn:hover {
            background: #fad2cf;
        }
        .save-btn {
            background: #1a73e8;
            color: white;
        }
        .save-btn:hover {
            background: #1765cc;
        }
        .cancel-btn {
            background: #f1f3f4;
            color: #5f6368;
        }
        .cancel-btn:hover {
            background: #e8eaed;
        }
        .config-details {
            font-size: 13px;
            color: #5f6368;
            line-height: 1.6;
        }
        .config-details div {
            margin-bottom: 4px;
        }
        .form-group {
            margin-bottom: 16px;
        }
        .form-label {
            display: block;
            font-size: 13px;
            font-weight: 500;
            color: #5f6368;
            margin-bottom: 6px;
        }
        .form-input {
            width: 100%;
            padding: 10px 12px;
            border: 1px solid #dadce0;
            border-radius: 4px;
            font-size: 14px;
            font-family: monospace;
        }
        .form-input:focus {
            outline: none;
            border-color: #1a73e8;
        }
        .hidden {
            display: none;
        }
        .new-config-form {
            display: none;
            border: 2px dashed #dadce0;
            border-radius: 8px;
            padding: 20px;
            margin-top: 16px;
            background: white;
        }
        .new-config-form.visible {
            display: block;
        }
        .success-message {
            background: #e6f4ea;
            border: 1px solid #34a853;
            color: #137333;
            padding: 12px 16px;
            border-radius: 4px;
            margin-bottom: 16px;
            font-size: 14px;
            display: none;
        }
        .success-message.visible {
            display: block;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>⚙️ Configuration Editor</h1>
        <a href="/" class="back-link">← Back to Home</a>
    </div>

    <div class="container">
        <div id="successMessage" class="success-message"></div>

        <!-- PubSub Configs -->
        <div class="section">
            <div class="section-header">
                <div class="section-title">Google PubSub Configurations</div>
                <button class="add-btn" onclick="showNewConfigForm('pubsub')">+ Add New</button>
            </div>
            <div id="pubsubConfigs"></div>
            <div id="newPubSubForm" class="new-config-form">
                <div class="form-group">
                    <label class="form-label">Configuration Name</label>
                    <input type="text" class="form-input" id="newPubSubName" placeholder="e.g., TMS PubSub">
                </div>
                <div class="form-group">
                    <label class="form-label">Emulator Host</label>
                    <input type="text" class="form-input" id="newPubSubHost" placeholder="e.g., localhost:8086">
                </div>
                <div class="form-group">
                    <label class="form-label">Project ID</label>
                    <input type="text" class="form-input" id="newPubSubProject" placeholder="e.g., tms-suncorp-local">
                </div>
                <div class="form-group">
                    <label class="form-label">Subscription ID</label>
                    <input type="text" class="form-input" id="newPubSubSub" placeholder="e.g., cloudevents.subscription">
                </div>
                <div style="display: flex; gap: 8px; justify-content: flex-end;">
                    <button class="cancel-btn" onclick="hideNewConfigForm('pubsub')">Cancel</button>
                    <button class="save-btn" onclick="saveNewPubSubConfig()">Save Configuration</button>
                </div>
            </div>
        </div>

        <!-- Kafka Configs -->
        <div class="section">
            <div class="section-header">
                <div class="section-title">Kafka / EventMesh Configurations</div>
                <button class="add-btn" onclick="showNewConfigForm('kafka')">+ Add New</button>
            </div>
            <div id="kafkaConfigs"></div>
            <div id="newKafkaForm" class="new-config-form">
                <div class="form-group">
                    <label class="form-label">Configuration Name</label>
                    <input type="text" class="form-input" id="newKafkaName" placeholder="e.g., Unica Events">
                </div>
                <div class="form-group">
                    <label class="form-label">Brokers</label>
                    <input type="text" class="form-input" id="newKafkaBrokers" placeholder="e.g., localhost:19092">
                </div>
                <div class="form-group">
                    <label class="form-label">Topic</label>
                    <input type="text" class="form-input" id="newKafkaTopic" placeholder="e.g., unica.marketing.response.events">
                </div>
                <div class="form-group">
                    <label class="form-label">Consumer Group</label>
                    <input type="text" class="form-input" id="newKafkaGroup" placeholder="e.g., testing-studio-consumer">
                </div>
                <div class="form-group">
                    <label class="form-label">Schema Registry URL</label>
                    <input type="text" class="form-input" id="newKafkaSchema" placeholder="e.g., http://localhost:18081">
                </div>
                <div style="display: flex; gap: 8px; justify-content: flex-end;">
                    <button class="cancel-btn" onclick="hideNewConfigForm('kafka')">Cancel</button>
                    <button class="save-btn" onclick="saveNewKafkaConfig()">Save Configuration</button>
                </div>
            </div>
        </div>
    </div>

    <script>
        let configs = { pubsubConfigs: [], kafkaConfigs: [] };

        async function loadConfigs() {
            const response = await fetch('/api/configs');
            configs = await response.json();
            renderConfigs();
        }

        function renderConfigs() {
            renderPubSubConfigs();
            renderKafkaConfigs();
        }

        function renderPubSubConfigs() {
            const container = document.getElementById('pubsubConfigs');
            container.innerHTML = configs.pubsubConfigs.map((config, index) => ` + "`" + `
                <div class="config-item" id="pubsub-${index}">
                    <div class="config-item-header">
                        <div class="config-name">${config.name}</div>
                        <div class="config-actions">
                            <button class="edit-btn" onclick="editPubSubConfig(${index})">Edit</button>
                            <button class="delete-btn" onclick="deletePubSubConfig(${index})">Delete</button>
                        </div>
                    </div>
                    <div class="config-details" id="pubsub-details-${index}">
                        <div><strong>Emulator Host:</strong> ${config.emulatorHost}</div>
                        <div><strong>Project ID:</strong> ${config.projectId}</div>
                        <div><strong>Subscription ID:</strong> ${config.subscriptionId}</div>
                    </div>
                    <div class="hidden" id="pubsub-form-${index}">
                        <div class="form-group">
                            <label class="form-label">Configuration Name</label>
                            <input type="text" class="form-input" id="edit-pubsub-name-${index}" value="${config.name}">
                        </div>
                        <div class="form-group">
                            <label class="form-label">Emulator Host</label>
                            <input type="text" class="form-input" id="edit-pubsub-host-${index}" value="${config.emulatorHost}">
                        </div>
                        <div class="form-group">
                            <label class="form-label">Project ID</label>
                            <input type="text" class="form-input" id="edit-pubsub-project-${index}" value="${config.projectId}">
                        </div>
                        <div class="form-group">
                            <label class="form-label">Subscription ID</label>
                            <input type="text" class="form-input" id="edit-pubsub-sub-${index}" value="${config.subscriptionId}">
                        </div>
                        <div style="display: flex; gap: 8px; justify-content: flex-end;">
                            <button class="cancel-btn" onclick="cancelEditPubSub(${index})">Cancel</button>
                            <button class="save-btn" onclick="savePubSubConfig(${index})">Save</button>
                        </div>
                    </div>
                </div>
            ` + "`" + `).join('');
        }

        function renderKafkaConfigs() {
            const container = document.getElementById('kafkaConfigs');
            container.innerHTML = configs.kafkaConfigs.map((config, index) => ` + "`" + `
                <div class="config-item" id="kafka-${index}">
                    <div class="config-item-header">
                        <div class="config-name">${config.name}</div>
                        <div class="config-actions">
                            <button class="edit-btn" onclick="editKafkaConfig(${index})">Edit</button>
                            <button class="delete-btn" onclick="deleteKafkaConfig(${index})">Delete</button>
                        </div>
                    </div>
                    <div class="config-details" id="kafka-details-${index}">
                        <div><strong>Brokers:</strong> ${config.brokers}</div>
                        <div><strong>Topic:</strong> ${config.topic}</div>
                        <div><strong>Consumer Group:</strong> ${config.consumerGroup}</div>
                        <div><strong>Schema Registry:</strong> ${config.schemaRegistry}</div>
                    </div>
                    <div class="hidden" id="kafka-form-${index}">
                        <div class="form-group">
                            <label class="form-label">Configuration Name</label>
                            <input type="text" class="form-input" id="edit-kafka-name-${index}" value="${config.name}">
                        </div>
                        <div class="form-group">
                            <label class="form-label">Brokers</label>
                            <input type="text" class="form-input" id="edit-kafka-brokers-${index}" value="${config.brokers}">
                        </div>
                        <div class="form-group">
                            <label class="form-label">Topic</label>
                            <input type="text" class="form-input" id="edit-kafka-topic-${index}" value="${config.topic}">
                        </div>
                        <div class="form-group">
                            <label class="form-label">Consumer Group</label>
                            <input type="text" class="form-input" id="edit-kafka-group-${index}" value="${config.consumerGroup}">
                        </div>
                        <div class="form-group">
                            <label class="form-label">Schema Registry URL</label>
                            <input type="text" class="form-input" id="edit-kafka-schema-${index}" value="${config.schemaRegistry}">
                        </div>
                        <div style="display: flex; gap: 8px; justify-content: flex-end;">
                            <button class="cancel-btn" onclick="cancelEditKafka(${index})">Cancel</button>
                            <button class="save-btn" onclick="saveKafkaConfig(${index})">Save</button>
                        </div>
                    </div>
                </div>
            ` + "`" + `).join('');
        }

        function editPubSubConfig(index) {
            document.getElementById(` + "`pubsub-details-${index}`" + `).classList.add('hidden');
            document.getElementById(` + "`pubsub-form-${index}`" + `).classList.remove('hidden');
        }

        function cancelEditPubSub(index) {
            document.getElementById(` + "`pubsub-details-${index}`" + `).classList.remove('hidden');
            document.getElementById(` + "`pubsub-form-${index}`" + `).classList.add('hidden');
        }

        async function savePubSubConfig(index) {
            const updatedConfig = {
                name: document.getElementById(` + "`edit-pubsub-name-${index}`" + `).value,
                emulatorHost: document.getElementById(` + "`edit-pubsub-host-${index}`" + `).value,
                projectId: document.getElementById(` + "`edit-pubsub-project-${index}`" + `).value,
                subscriptionId: document.getElementById(` + "`edit-pubsub-sub-${index}`" + `).value
            };

            configs.pubsubConfigs[index] = updatedConfig;
            await saveConfigs();
            showSuccess('PubSub configuration updated successfully!');
        }

        function editKafkaConfig(index) {
            document.getElementById(` + "`kafka-details-${index}`" + `).classList.add('hidden');
            document.getElementById(` + "`kafka-form-${index}`" + `).classList.remove('hidden');
        }

        function cancelEditKafka(index) {
            document.getElementById(` + "`kafka-details-${index}`" + `).classList.remove('hidden');
            document.getElementById(` + "`kafka-form-${index}`" + `).classList.add('hidden');
        }

        async function saveKafkaConfig(index) {
            const updatedConfig = {
                name: document.getElementById(` + "`edit-kafka-name-${index}`" + `).value,
                brokers: document.getElementById(` + "`edit-kafka-brokers-${index}`" + `).value,
                topic: document.getElementById(` + "`edit-kafka-topic-${index}`" + `).value,
                consumerGroup: document.getElementById(` + "`edit-kafka-group-${index}`" + `).value,
                schemaRegistry: document.getElementById(` + "`edit-kafka-schema-${index}`" + `).value
            };

            configs.kafkaConfigs[index] = updatedConfig;
            await saveConfigs();
            showSuccess('Kafka configuration updated successfully!');
        }

        async function deletePubSubConfig(index) {
            if (!confirm('Are you sure you want to delete this configuration?')) return;
            configs.pubsubConfigs.splice(index, 1);
            await saveConfigs();
            showSuccess('PubSub configuration deleted successfully!');
        }

        async function deleteKafkaConfig(index) {
            if (!confirm('Are you sure you want to delete this configuration?')) return;
            configs.kafkaConfigs.splice(index, 1);
            await saveConfigs();
            showSuccess('Kafka configuration deleted successfully!');
        }

        function showNewConfigForm(type) {
            if (type === 'pubsub') {
                document.getElementById('newPubSubForm').classList.add('visible');
            } else {
                document.getElementById('newKafkaForm').classList.add('visible');
            }
        }

        function hideNewConfigForm(type) {
            if (type === 'pubsub') {
                document.getElementById('newPubSubForm').classList.remove('visible');
                clearPubSubForm();
            } else {
                document.getElementById('newKafkaForm').classList.remove('visible');
                clearKafkaForm();
            }
        }

        function clearPubSubForm() {
            document.getElementById('newPubSubName').value = '';
            document.getElementById('newPubSubHost').value = '';
            document.getElementById('newPubSubProject').value = '';
            document.getElementById('newPubSubSub').value = '';
        }

        function clearKafkaForm() {
            document.getElementById('newKafkaName').value = '';
            document.getElementById('newKafkaBrokers').value = '';
            document.getElementById('newKafkaTopic').value = '';
            document.getElementById('newKafkaGroup').value = '';
            document.getElementById('newKafkaSchema').value = '';
        }

        async function saveNewPubSubConfig() {
            const newConfig = {
                name: document.getElementById('newPubSubName').value,
                emulatorHost: document.getElementById('newPubSubHost').value,
                projectId: document.getElementById('newPubSubProject').value,
                subscriptionId: document.getElementById('newPubSubSub').value
            };

            if (!newConfig.name || !newConfig.emulatorHost || !newConfig.projectId || !newConfig.subscriptionId) {
                alert('Please fill in all fields');
                return;
            }

            configs.pubsubConfigs.push(newConfig);
            await saveConfigs();
            hideNewConfigForm('pubsub');
            showSuccess('New PubSub configuration added successfully!');
        }

        async function saveNewKafkaConfig() {
            const newConfig = {
                name: document.getElementById('newKafkaName').value,
                brokers: document.getElementById('newKafkaBrokers').value,
                topic: document.getElementById('newKafkaTopic').value,
                consumerGroup: document.getElementById('newKafkaGroup').value,
                schemaRegistry: document.getElementById('newKafkaSchema').value
            };

            if (!newConfig.name || !newConfig.brokers || !newConfig.topic || !newConfig.consumerGroup || !newConfig.schemaRegistry) {
                alert('Please fill in all fields');
                return;
            }

            configs.kafkaConfigs.push(newConfig);
            await saveConfigs();
            hideNewConfigForm('kafka');
            showSuccess('New Kafka configuration added successfully!');
        }

        async function saveConfigs() {
            const response = await fetch('/api/configs/save', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(configs)
            });

            if (response.ok) {
                await loadConfigs();
            } else {
                alert('Failed to save configurations');
            }
        }

        function showSuccess(message) {
            const msgEl = document.getElementById('successMessage');
            msgEl.textContent = message;
            msgEl.classList.add('visible');
            setTimeout(() => {
                msgEl.classList.remove('visible');
            }, 3000);
        }

        loadConfigs();
    </script>
</body>
</html>`
