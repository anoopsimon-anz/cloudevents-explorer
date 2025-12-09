package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/linkedin/goavro/v2"
	"google.golang.org/api/option"
)

type PubSubConfig struct {
	Name           string `json:"name"`
	EmulatorHost   string `json:"emulatorHost"`
	ProjectID      string `json:"projectId"`
	SubscriptionID string `json:"subscriptionId"`
}

type KafkaConfig struct {
	Name           string `json:"name"`
	Brokers        string `json:"brokers"`
	Topic          string `json:"topic"`
	ConsumerGroup  string `json:"consumerGroup"`
	SchemaRegistry string `json:"schemaRegistry"`
}

type Config struct {
	PubSubConfigs []PubSubConfig `json:"pubsubConfigs"`
	KafkaConfigs  []KafkaConfig  `json:"kafkaConfigs"`
}

type CloudEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Subject   string                 `json:"subject"`
	Source    string                 `json:"source"`
	Schema    string                 `json:"schema"`
	Published string                 `json:"published"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp int64                  `json:"timestamp"`
	RawData   string                 `json:"rawData,omitempty"`
}

var (
	configFile     = "configs.json"
	mu             sync.RWMutex
	config         Config
	messageStore   []CloudEvent
	messageStoreMu sync.RWMutex
)

// decodeAvroMessage decodes an Avro message using schema from registry
func decodeAvroMessage(data []byte, schemaRegistryURL string) (map[string]interface{}, error) {
	if len(data) < 5 {
		return nil, fmt.Errorf("message too short")
	}

	// First byte is magic byte (should be 0)
	// Next 4 bytes are schema ID (big-endian)
	schemaID := binary.BigEndian.Uint32(data[1:5])

	// Fetch schema from registry
	schemaURL := fmt.Sprintf("%s/schemas/ids/%d", schemaRegistryURL, schemaID)
	resp, err := http.Get(schemaURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schema: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("schema registry returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema response: %w", err)
	}

	var schemaResp struct {
		Schema string `json:"schema"`
	}
	if err := json.Unmarshal(body, &schemaResp); err != nil {
		return nil, fmt.Errorf("failed to parse schema response: %w", err)
	}

	// Create Avro codec
	codec, err := goavro.NewCodec(schemaResp.Schema)
	if err != nil {
		return nil, fmt.Errorf("failed to create codec: %w", err)
	}

	// Decode the message (skip first 5 bytes)
	native, _, err := codec.NativeFromBinary(data[5:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode avro: %w", err)
	}

	// Convert to map[string]interface{}
	if result, ok := native.(map[string]interface{}); ok {
		return result, nil
	}

	return nil, fmt.Errorf("decoded data is not a map")
}

func loadConfig() error {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			config = Config{
				PubSubConfigs: []PubSubConfig{
					{
						Name:           "TMS Local",
						EmulatorHost:   "localhost:8086",
						ProjectID:      "tms-suncorp-local",
						SubscriptionID: "cloudevents.subscription",
					},
				},
				KafkaConfigs: []KafkaConfig{
					{
						Name:           "TMS Unica Local",
						Brokers:        "localhost:19092",
						Topic:          "unica.marketing.response.events",
						ConsumerGroup:  "cloudevents-explorer",
						SchemaRegistry: "http://localhost:18081",
					},
				},
			}
			return saveConfigLocked()
		}
		return err
	}

	return json.Unmarshal(data, &config)
}

func saveConfigLocked() error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0644)
}

func saveConfig() error {
	mu.Lock()
	defer mu.Unlock()
	return saveConfigLocked()
}

func getConfig() Config {
	mu.RLock()
	defer mu.RUnlock()
	return config
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>CloudEvents Explorer</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", sans-serif;
            background: #f5f5f5;
            color: #202124;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .landing {
            text-align: center;
            max-width: 800px;
            padding: 40px;
        }
        h1 {
            font-size: 32px;
            color: #202124;
            margin-bottom: 8px;
            font-weight: 400;
        }
        .subtitle {
            font-size: 16px;
            color: #5f6368;
            margin-bottom: 48px;
        }
        .options {
            display: grid;
            grid-template-columns: repeat(2, 1fr);
            gap: 16px;
            margin-top: 32px;
        }
        .option-card {
            background: white;
            border: 1px solid #dadce0;
            border-radius: 8px;
            padding: 32px;
            cursor: pointer;
            transition: all 0.2s;
            text-decoration: none;
            color: inherit;
            display: block;
        }
        .option-card:hover {
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
            border-color: #1a73e8;
        }
        .option-title {
            font-size: 20px;
            font-weight: 500;
            color: #202124;
            margin-bottom: 8px;
        }
        .option-desc {
            font-size: 14px;
            color: #5f6368;
            line-height: 1.5;
        }
        .badge {
            display: inline-block;
            background: #e8f0fe;
            color: #1967d2;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 11px;
            font-weight: 500;
            margin-top: 12px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
    </style>
</head>
<body>
    <div class="landing">
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 48px;">
            <div>
                <h1>CloudEvents Explorer</h1>
                <p class="subtitle">Choose your message streaming platform</p>
            </div>
            <div style="position: relative;">
                <button onclick="toggleToolsMenu()" style="background: white; border: 1px solid #dadce0; padding: 8px 16px; border-radius: 4px; cursor: pointer; color: #5f6368; font-size: 14px; font-weight: 500;">
                    Tools ▼
                </button>
                <div id="toolsMenu" style="display: none; position: absolute; right: 0; top: 40px; background: white; border: 1px solid #dadce0; border-radius: 4px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); min-width: 200px; z-index: 100;">
                    <a href="/flow-diagram" style="display: block; padding: 12px 16px; color: #202124; text-decoration: none; border-bottom: 1px solid #dadce0;">Flow Diagram</a>
                    <a onclick="openBase64Tool(); toggleToolsMenu(); return false;" href="#" style="display: block; padding: 12px 16px; color: #202124; text-decoration: none;">Base64 Encoder/Decoder</a>
                </div>
            </div>
        </div>

        <div class="options">
            <a href="/pubsub" class="option-card">
                <div class="option-title">Google PubSub</div>
                <div class="option-desc">View CloudEvents from Google Cloud PubSub subscriptions</div>
                <span class="badge">CloudEvents</span>
            </a>

            <a href="/kafka" class="option-card">
                <div class="option-title">Kafka / EventMesh</div>
                <div class="option-desc">Consume Avro messages from Kafka topics</div>
                <span class="badge">Avro Schema</span>
            </a>
        </div>
    </div>

    <div id="base64Modal" style="display: none; position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); z-index: 1000; align-items: center; justify-content: center;">
        <div style="background: white; border-radius: 8px; max-width: 900px; width: 90%; max-height: 90vh; overflow: hidden; display: flex; flex-direction: column;">
            <div style="padding: 20px; border-bottom: 1px solid #dadce0; display: flex; justify-content: space-between; align-items: center;">
                <h2 style="font-size: 20px; font-weight: 500; color: #202124;">Base64 Encoder/Decoder</h2>
                <button onclick="closeBase64Tool()" style="background: none; border: none; font-size: 24px; cursor: pointer; color: #5f6368;">&times;</button>
            </div>
            <div style="display: flex; flex: 1; overflow: hidden;">
                <div style="flex: 1; padding: 20px; border-right: 1px solid #dadce0; display: flex; flex-direction: column;">
                    <label style="font-size: 13px; color: #5f6368; font-weight: 500; margin-bottom: 8px;">Input Text:</label>
                    <textarea id="base64Input" style="flex: 1; font-family: 'Monaco', monospace; font-size: 13px; border: 1px solid #dadce0; border-radius: 4px; padding: 12px; resize: none;" placeholder="Enter text or Base64 string here"></textarea>
                    <div style="display: flex; gap: 8px; margin-top: 12px;">
                        <button onclick="encodeBase64()" style="flex: 1; background: #1a73e8; color: white; border: none; padding: 10px 20px; border-radius: 4px; cursor: pointer; font-weight: 500;">Encode to Base64</button>
                        <button onclick="decodeBase64()" style="flex: 1; background: #188038; color: white; border: none; padding: 10px 20px; border-radius: 4px; cursor: pointer; font-weight: 500;">Decode from Base64</button>
                    </div>
                </div>
                <div style="flex: 1; padding: 20px; display: flex; flex-direction: column;">
                    <label style="font-size: 13px; color: #5f6368; font-weight: 500; margin-bottom: 8px;">Output:</label>
                    <textarea id="base64Output" readonly style="flex: 1; font-family: 'Monaco', monospace; font-size: 13px; border: 1px solid #dadce0; border-radius: 4px; padding: 12px; resize: none; background: #f8f9fa;"></textarea>
                    <button onclick="copyOutput()" style="margin-top: 12px; background: white; color: #5f6368; border: 1px solid #dadce0; padding: 10px 20px; border-radius: 4px; cursor: pointer; font-weight: 500;">Copy to Clipboard</button>
                </div>
            </div>
        </div>
    </div>

    <script>
        function toggleToolsMenu() {
            const menu = document.getElementById('toolsMenu');
            menu.style.display = menu.style.display === 'none' ? 'block' : 'none';
        }

        // Close menu when clicking outside
        document.addEventListener('click', function(e) {
            const menu = document.getElementById('toolsMenu');
            const button = e.target.closest('button');
            if (menu && !menu.contains(e.target) && (!button || button.textContent.indexOf('Tools') === -1)) {
                menu.style.display = 'none';
            }
        });

        function openBase64Tool() {
            document.getElementById('base64Modal').style.display = 'flex';
        }

        function closeBase64Tool() {
            document.getElementById('base64Modal').style.display = 'none';
            document.getElementById('base64Input').value = '';
            document.getElementById('base64Output').value = '';
        }

        function encodeBase64() {
            const input = document.getElementById('base64Input').value;
            const output = document.getElementById('base64Output');

            if (!input) {
                output.value = 'Error: Please enter some text to encode';
                return;
            }

            try {
                const encoded = btoa(unescape(encodeURIComponent(input)));
                output.value = encoded;
            } catch (e) {
                output.value = 'Error: Failed to encode - ' + e.message;
            }
        }

        function decodeBase64() {
            const input = document.getElementById('base64Input').value;
            const output = document.getElementById('base64Output');

            if (!input) {
                output.value = 'Error: Please enter a Base64 string to decode';
                return;
            }

            try {
                const decoded = decodeURIComponent(escape(atob(input)));
                output.value = decoded;
            } catch (e) {
                output.value = 'Error: Invalid Base64 string - ' + e.message;
            }
        }

        function copyOutput() {
            const output = document.getElementById('base64Output');
            if (!output.value || output.value.startsWith('Error:')) {
                return;
            }
            output.select();
            document.execCommand('copy');

            const btn = event.target;
            const originalText = btn.textContent;
            btn.textContent = 'Copied!';
            btn.style.background = '#188038';
            btn.style.color = 'white';
            setTimeout(function() {
                btn.textContent = originalText;
                btn.style.background = 'white';
                btn.style.color = '#5f6368';
            }, 2000);
        }

        // Close modal on outside click
        document.getElementById('base64Modal')?.addEventListener('click', function(e) {
            if (e.target === this) {
                closeBase64Tool();
            }
        });
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func handlePubSub(w http.ResponseWriter, r *http.Request) {
	html := getBaseHTML("PubSub CloudEvents", `
        <div class="panel">
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
        </div>
`, `
        async function refreshConfigs() {
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

        refreshConfigs();
`)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func handleKafka(w http.ResponseWriter, r *http.Request) {
	html := getBaseHTML("Kafka EventMesh", `
        <div class="panel">
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
            <div style="background: white; border-radius: 8px; max-width: 700px; width: 90%; max-height: 90vh; overflow: hidden; display: flex; flex-direction: column;">
                <div style="padding: 20px; border-bottom: 1px solid #dadce0; display: flex; justify-content: space-between; align-items: center;">
                    <h2 style="font-size: 20px; font-weight: 500; color: #202124;">Publish Kafka Message</h2>
                    <button onclick="closePublishModal()" style="background: none; border: none; font-size: 24px; cursor: pointer; color: #5f6368;">&times;</button>
                </div>
                <div style="flex: 1; padding: 20px; display: flex; flex-direction: column; overflow: hidden;">
                    <label style="font-size: 13px; color: #5f6368; font-weight: 500; margin-bottom: 8px;">Message JSON:</label>
                    <textarea id="publishMessageJson" style="flex: 1; font-family: 'Monaco', monospace; font-size: 13px; border: 1px solid #dadce0; border-radius: 4px; padding: 12px; resize: none;" placeholder='{"header": {...}, "marketingResponse": {...}}'></textarea>
                    <div style="margin-top: 12px; padding: 12px; background: #f8f9fa; border-radius: 4px; font-size: 12px; color: #5f6368;">
                        <div>Topic: <strong id="publishTopic">-</strong></div>
                        <div>Schema Registry: <strong id="publishSchema">-</strong></div>
                    </div>
                    <button onclick="publishMessage()" style="margin-top: 12px; background: #188038; color: white; border: none; padding: 10px 20px; border-radius: 4px; cursor: pointer; font-weight: 500;">Publish to Kafka</button>
                </div>
            </div>
        </div>
`, `
        async function refreshConfigs() {
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

        refreshConfigs();
`)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func handleFlowDiagram(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TMS Event Flow - CloudEvents Explorer</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif;
            background: #ffffff;
            color: #000;
            min-height: 100vh;
            padding: 20px;
        }
        .topbar {
            background: #000;
            color: white;
            padding: 16px 24px;
            display: flex;
            align-items: center;
            gap: 16px;
            margin-bottom: 40px;
        }
        .back-btn {
            color: white;
            padding: 8px 16px;
            border: 2px solid white;
            text-decoration: none;
            font-size: 14px;
            font-weight: 600;
        }
        h1 {
            font-size: 28px;
            text-align: center;
            margin-bottom: 50px;
            color: #000;
            font-weight: 700;
            letter-spacing: -0.5px;
        }
        .flow-container {
            max-width: 1400px;
            margin: 0 auto;
        }
        .flow-row {
            display: flex;
            gap: 20px;
            margin-bottom: 40px;
            align-items: stretch;
        }
        .box {
            background: white;
            border: 2px solid #000;
            padding: 24px;
            position: relative;
            flex: 1;
            min-height: 180px;
            display: flex;
            flex-direction: column;
        }
        .box-external {
            background: #f5f5f5;
            border: 2px dashed #666;
        }
        .box-title {
            font-size: 16px;
            font-weight: 700;
            margin-bottom: 12px;
            color: #000;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        .box-desc {
            font-size: 13px;
            color: #333;
            line-height: 1.7;
            flex: 1;
        }
        .arrow-right {
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 32px;
            color: #000;
            font-weight: bold;
            width: 40px;
        }
        .arrow-down {
            text-align: center;
            font-size: 32px;
            color: #000;
            font-weight: bold;
            margin: 20px 0;
        }
        .step-number {
            position: absolute;
            top: -12px;
            left: -12px;
            background: #000;
            color: white;
            width: 32px;
            height: 32px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: bold;
            font-size: 14px;
        }
        .highlight {
            font-weight: 700;
            text-decoration: underline;
        }
        code {
            background: #f0f0f0;
            padding: 2px 6px;
            font-family: 'Courier New', monospace;
            font-size: 12px;
        }
        .topic-link {
            color: #1a73e8;
            text-decoration: underline;
            cursor: pointer;
            font-weight: 700;
        }
        .topic-link:hover {
            color: #1557b0;
        }
        .json-modal {
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: rgba(0,0,0,0.7);
            z-index: 1000;
            align-items: center;
            justify-content: center;
        }
        .json-modal-content {
            background: white;
            border: 3px solid #000;
            max-width: 800px;
            width: 90%;
            max-height: 90vh;
            overflow: auto;
            position: relative;
        }
        .json-modal-header {
            background: #000;
            color: white;
            padding: 16px 20px;
            display: flex;
            justify-content: space-between;
            align-items: center;
            position: sticky;
            top: 0;
            z-index: 1;
        }
        .json-modal-title {
            font-size: 16px;
            font-weight: 700;
        }
        .json-modal-close {
            background: white;
            color: #000;
            border: none;
            font-size: 24px;
            cursor: pointer;
            padding: 0 8px;
            font-weight: bold;
        }
        .json-modal-body {
            padding: 20px;
        }
        .json-display {
            background: #f5f5f5;
            border: 1px solid #ccc;
            padding: 16px;
            font-family: 'Courier New', monospace;
            font-size: 12px;
            line-height: 1.6;
            overflow-x: auto;
            white-space: pre;
        }
        .animation-modal {
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: rgba(0,0,0,0.85);
            z-index: 1000;
            align-items: center;
            justify-content: center;
        }
        .animation-content {
            background: white;
            border: 3px solid #000;
            max-width: 900px;
            width: 90%;
            padding: 30px;
            position: relative;
        }
        .animation-scene {
            border: 2px solid #ccc;
            padding: 40px;
            background: #fafafa;
            margin: 20px 0;
            min-height: 400px;
            position: relative;
            overflow: hidden;
        }
        .event-handler-worker {
            width: 80px;
            height: 100px;
            position: absolute;
            left: 50%;
            top: 50%;
            transform: translate(-50%, -50%);
            animation: workerPulse 2s infinite;
        }
        @keyframes workerPulse {
            0%, 100% { transform: translate(-50%, -50%) scale(1); }
            50% { transform: translate(-50%, -50%) scale(1.05); }
        }
        .worker-head {
            width: 40px;
            height: 40px;
            background: #000;
            border-radius: 50%;
            margin: 0 auto 5px;
        }
        .worker-body {
            width: 60px;
            height: 50px;
            background: #333;
            margin: 0 auto;
            position: relative;
        }
        .worker-arms {
            position: absolute;
            width: 100%;
            height: 100%;
        }
        .worker-arm {
            width: 30px;
            height: 8px;
            background: #333;
            position: absolute;
            top: 10px;
        }
        .worker-arm.left {
            left: -25px;
            transform-origin: right center;
            animation: armWave 1.5s infinite;
        }
        .worker-arm.right {
            right: -25px;
            transform-origin: left center;
            animation: armWave 1.5s infinite 0.75s;
        }
        @keyframes armWave {
            0%, 100% { transform: rotate(0deg); }
            50% { transform: rotate(-30deg); }
        }
        .message-flow {
            position: absolute;
            width: 60px;
            height: 40px;
            background: #1a73e8;
            border: 2px solid #000;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 10px;
            font-weight: bold;
            color: white;
        }
        .message-incoming {
            animation: moveIncoming 3s infinite;
            left: -100px;
            top: 30%;
        }
        @keyframes moveIncoming {
            0% { left: -100px; opacity: 1; }
            45% { left: 50%; transform: translateX(-50%) scale(1); }
            50% { left: 50%; transform: translateX(-50%) scale(0.5); opacity: 0.5; }
            55% { left: 50%; transform: translateX(-50%) scale(0); opacity: 0; }
            100% { left: 50%; transform: translateX(-50%) scale(0); opacity: 0; }
        }
        .message-outgoing {
            animation: moveOutgoing 3s infinite;
            right: -100px;
            bottom: 30%;
            background: #188038;
        }
        @keyframes moveOutgoing {
            0%, 55% { right: -100px; opacity: 0; transform: scale(0); }
            60% { right: 50%; transform: translateX(50%) scale(0.5); opacity: 0.5; }
            65% { right: 50%; transform: translateX(50%) scale(1); opacity: 1; }
            100% { right: -100px; opacity: 1; }
        }
        .status-text {
            position: absolute;
            bottom: 20px;
            left: 50%;
            transform: translateX(-50%);
            font-size: 14px;
            font-weight: bold;
            text-align: center;
            animation: statusBlink 3s infinite;
        }
        @keyframes statusBlink {
            0%, 45% { opacity: 0; }
            50%, 95% { opacity: 1; }
            100% { opacity: 0; }
        }
        .label-kafka {
            position: absolute;
            left: 20px;
            top: 50%;
            transform: translateY(-50%);
            font-size: 12px;
            font-weight: bold;
            text-align: center;
        }
        .label-temporal {
            position: absolute;
            right: 20px;
            top: 50%;
            transform: translateY(-50%);
            font-size: 12px;
            font-weight: bold;
            text-align: center;
        }
        .legend {
            max-width: 1400px;
            margin: 50px auto 30px;
            padding: 20px;
            border: 2px solid #000;
            background: #f9f9f9;
        }
        .legend-title {
            font-size: 14px;
            font-weight: 700;
            margin-bottom: 10px;
            text-transform: uppercase;
        }
        .legend-text {
            font-size: 13px;
            line-height: 1.6;
            color: #333;
        }
    </style>
</head>
<body>
    <div class="topbar">
        <a href="/" class="back-btn">← Back to Home</a>
        <h2 style="font-size: 18px; color: white; font-weight: 600;">COMMS EPIC Eventing Flow</h2>
    </div>

    <h1>COMMS EPIC EVENTING FLOW AND CAP INTEGRATION</h1>

    <div class="legend">
        <div class="legend-title">Overview</div>
        <div class="legend-text">
            This diagram illustrates the end-to-end event flow for customer communication activities,
            from marketing program events through to data analytics. The system ensures reliable processing
            and comprehensive reporting of all customer interactions.
        </div>
    </div>

    <div class="flow-container">
        <!-- Row 1: AMP → EventMesh → TMS Event Handler -->
        <div class="flow-row">
            <div class="box box-external">
                <div class="step-number">1</div>
                <div class="box-title">ANZ Marketing Program (AMP)</div>
                <div class="box-desc">
                    AMP generates customer communication events (email opens, clicks, campaign responses) and
                    <span class="highlight">publishes them to EventMesh</span> for downstream processing.
                    <br><br><em style="color: #666;">* External system - not available in test environment</em>
                </div>
            </div>
            <div class="arrow-right">→</div>
            <div class="box">
                <div class="step-number">2</div>
                <div class="box-title">EventMesh (Kafka Platform)</div>
                <div class="box-desc">
                    EventMesh is ANZ's enterprise <span class="highlight">eventing platform</span> built on Kafka.
                    Messages are stored in topic <code><span class="topic-link" onclick="showSampleMessage()">unica.marketing.response.events</span></code> with Avro schema validation.
                </div>
            </div>
            <div class="arrow-right">→</div>
            <div class="box" style="cursor: pointer;" onclick="showEventHandlerAnimation()">
                <div class="step-number">3</div>
                <div class="box-title">TMS Event Handler (24x7 Service) <span style="font-size: 12px; color: #1a73e8;">▶ Click to see animation</span></div>
                <div class="box-desc">
                    The TMS Event Handler continuously <span class="highlight">listens to the Kafka topic</span>.
                    Upon receiving a message, it immediately triggers a Temporal workflow for processing.
                </div>
            </div>
        </div>

        <div class="arrow-down">↓</div>

        <!-- Row 2: Temporal Workflow → CAP Diary → Success Response -->
        <div class="flow-row">
            <div class="box">
                <div class="step-number">4</div>
                <div class="box-title">Temporal Workflow Execution</div>
                <div class="box-desc">
                    The workflow orchestrates the communication event processing and makes a
                    <span class="highlight">REST HTTP call to CAP Diary</span> to record the customer interaction in the system of record.
                </div>
            </div>
            <div class="arrow-right">→</div>
            <div class="box">
                <div class="step-number">5</div>
                <div class="box-title">CAP Diary Integration</div>
                <div class="box-desc">
                    CAP (Customer Activity Platform) receives the diary entry request, validates the data,
                    and returns <span class="highlight">HTTP 200 OK</span> upon successful creation, confirming the transaction.
                </div>
            </div>
            <div class="arrow-right">→</div>
            <div class="box">
                <div class="step-number">6</div>
                <div class="box-title">Google PubSub (CloudEvents)</div>
                <div class="box-desc">
                    TMS publishes a <span class="topic-link" onclick="showCloudEvent()">CloudEvent</span> to Google PubSub confirming successful completion.
                    This event includes workflow execution metadata and transaction status.
                </div>
            </div>
        </div>

        <div class="arrow-down">↓</div>

        <!-- Row 3: BigQuery → Complete -->
        <div class="flow-row">
            <div class="box box-external">
                <div class="step-number">7</div>
                <div class="box-title">BigQuery Analytics Platform</div>
                <div class="box-desc">
                    CloudEvents are automatically <span class="highlight">streamed from PubSub to BigQuery</span> tables
                    on Google Cloud Platform. This data powers reporting dashboards and business intelligence analytics.
                    <br><br><em style="color: #666;">* External system - not available in test environment</em>
                </div>
            </div>
            <div class="arrow-right">→</div>
            <div class="box" style="border: 3px solid #000;">
                <div class="box-title">✓ Process Complete</div>
                <div class="box-desc">
                    End-to-end event processing complete: AMP → EventMesh → TMS → CAP → PubSub → BigQuery.
                    Customer interaction successfully recorded and available for reporting.
                </div>
            </div>
        </div>
    </div>

    <div class="legend" style="margin-top: 50px;">
        <div class="legend-title">Key Components</div>
        <div class="legend-text">
            <strong>AMP:</strong> ANZ Marketing Program |
            <strong>EventMesh:</strong> Kafka-based eventing platform |
            <strong>TMS:</strong> Technical Migration Service |
            <strong>CAP:</strong> Customer Activity Platform |
            <strong>PubSub:</strong> Google Cloud messaging service |
            <strong>BigQuery:</strong> Google Cloud data warehouse
        </div>
    </div>

    <!-- Sample Message Modal -->
    <div id="sampleMessageModal" class="json-modal">
        <div class="json-modal-content">
            <div class="json-modal-header">
                <div class="json-modal-title">Sample Kafka Message: unica.marketing.response.events</div>
                <button class="json-modal-close" onclick="closeSampleMessage()">×</button>
            </div>
            <div class="json-modal-body">
                <p style="margin-bottom: 16px; font-size: 14px; color: #333;">
                    This is an example of a marketing response event published to the Kafka topic by AMP:
                </p>
                <div class="json-display">{
    "header": {
        "eventUUID": "workflow-trigger-test-uuid",
        "hierarchy": "AU/RETAIL",
        "occurrenceDateTime": 1765250641805,
        "source": "unica-workflow-trigger-test",
        "spanId": "span-workflow-test",
        "subject": {
            "subjectId": "4015640641",
            "subjectType": "CRN"
        },
        "traceId": "trace-workflow-test"
    },
    "marketingResponse": {
        "accountNumber": "XX-DEFAULT-NO-NULLS",
        "actionCode": "Open",
        "browserName": "Chrome",
        "browserVersion": "109.0.0.0",
        "customerNumber": "4015640641",
        "customerType": "INDIVIDUAL",
        "deliveryStatus": "sent",
        "deviceName": "pc",
        "eventTimestamp": "2025-12-09T14:24:01+11:00",
        "failureDescription": "",
        "failureReason": "",
        "failureType": "",
        "metadataList": [
            {
                "name": "flex_01",
                "type": "string",
                "value": "01-Dec-2026"
            },
            {
                "name": "flex_02",
                "type": "string",
                "value": "T-21 MigRem"
            },
            {
                "name": "channel",
                "type": "string",
                "value": "Email"
            },
            {
                "name": "msgVersionId",
                "type": "string",
                "value": "SREM01"
            },
            {
                "name": "campaignCode",
                "type": "string",
                "value": "C005044"
            },
            {
                "name": "campaignLabel",
                "type": "string",
                "value": "Bank Led Migration"
            }
        ],
        "mobile": "false",
        "msgVersionName": "C005044_2EM1_InspireEducate2_V2_EDM_2025",
        "operatingSystem": "",
        "operatingSystemFamily": "",
        "optOutConfirmation": "",
        "optOutReason": "",
        "requestId": "",
        "requestType": "Batch",
        "sourceEventId": "",
        "trackingLabel": "",
        "trackingSource": "",
        "trackingType": "",
        "treatmentCode": "003531798",
        "uniqueMsgId": ""
    }
}</div>
            </div>
        </div>
    </div>

    <!-- Event Handler Animation Modal -->
    <div id="eventHandlerAnimation" class="animation-modal">
        <div class="animation-content">
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
                <h2 style="font-size: 20px; font-weight: 700; margin: 0;">TMS Event Handler - 24x7 Worker</h2>
                <button onclick="closeEventHandlerAnimation()" style="background: #000; color: white; border: none; font-size: 24px; cursor: pointer; padding: 0 12px; font-weight: bold;">×</button>
            </div>

            <p style="font-size: 14px; color: #333; margin-bottom: 20px;">
                The Event Handler is like a tireless worker that never sleeps. It continuously monitors the Kafka topic,
                instantly consumes incoming messages, processes them, and triggers Temporal workflows.
            </p>

            <div class="animation-scene">
                <div class="label-kafka">
                    <div>📡</div>
                    <div>KAFKA</div>
                    <div style="font-size: 10px;">Topic</div>
                </div>

                <div class="event-handler-worker">
                    <div class="worker-head"></div>
                    <div class="worker-body">
                        <div class="worker-arms">
                            <div class="worker-arm left"></div>
                            <div class="worker-arm right"></div>
                        </div>
                    </div>
                    <div style="text-align: center; margin-top: 10px; font-size: 11px; font-weight: bold;">
                        EVENT<br>HANDLER
                    </div>
                </div>

                <div class="message-flow message-incoming">MSG</div>
                <div class="message-flow message-outgoing">WF</div>

                <div class="label-temporal">
                    <div>⚙️</div>
                    <div>TEMPORAL</div>
                    <div style="font-size: 10px;">Workflow</div>
                </div>

                <div class="status-text">⚡ Processing Message → Triggering Workflow</div>
            </div>

            <div style="margin-top: 20px; padding: 15px; background: #f5f5f5; border-left: 3px solid #000;">
                <strong>What's happening:</strong>
                <ul style="margin: 10px 0 0 20px; font-size: 13px; line-height: 1.8;">
                    <li><strong>Blue box (MSG)</strong>: Marketing event arriving from Kafka</li>
                    <li><strong>Worker</strong>: Event Handler consuming the message (arms waving = working!)</li>
                    <li><strong>Green box (WF)</strong>: Temporal workflow being triggered</li>
                    <li><strong>Continuous cycle</strong>: This happens 24x7, never stops!</li>
                </ul>
            </div>
        </div>
    </div>

    <!-- CloudEvent Modal -->
    <div id="cloudEventModal" class="json-modal">
        <div class="json-modal-content">
            <div class="json-modal-header">
                <div class="json-modal-title">Sample CloudEvent: Migration Phase Completed</div>
                <button class="json-modal-close" onclick="closeCloudEvent()">×</button>
            </div>
            <div class="json-modal-body">
                <p style="margin-bottom: 16px; font-size: 14px; color: #333;">
                    This is an example CloudEvent published to Google PubSub when TMS completes a workflow successfully:
                </p>
                <div class="json-display">{
  "data": {
    "customers": {
      "4015645348": {
        "customerId": "4015645348",
        "groupId": "4015645348"
      }
    },
    "name": "name:\"migrations/4015645348\" phase:\"WriteProfileDiaryNote\" customer_group_id:\"4015645348\" customers:{customer_id:\"4015645348\"} status:STATUS_COMPLETED",
    "phase": "migrations/4015645348",
    "status": "STATUS_COMPLETED"
  }
}</div>
                <p style="margin-top: 16px; font-size: 13px; color: #666; line-height: 1.6;">
                    <strong>Key Fields:</strong><br>
                    • <strong>customerId</strong>: Customer identifier (4015645348)<br>
                    • <strong>phase</strong>: Migration phase identifier (migrations/4015645348)<br>
                    • <strong>status</strong>: Workflow completion status (STATUS_COMPLETED)<br>
                    • This event confirms the CAP Diary update was successful
                </p>
            </div>
        </div>
    </div>

    <script>
        function showSampleMessage() {
            document.getElementById('sampleMessageModal').style.display = 'flex';
        }

        function closeSampleMessage() {
            document.getElementById('sampleMessageModal').style.display = 'none';
        }

        function showCloudEvent() {
            document.getElementById('cloudEventModal').style.display = 'flex';
        }

        function closeCloudEvent() {
            document.getElementById('cloudEventModal').style.display = 'none';
        }

        function showEventHandlerAnimation() {
            document.getElementById('eventHandlerAnimation').style.display = 'flex';
        }

        function closeEventHandlerAnimation() {
            document.getElementById('eventHandlerAnimation').style.display = 'none';
        }

        // Close modals when clicking outside
        document.getElementById('sampleMessageModal').addEventListener('click', function(e) {
            if (e.target === this) {
                closeSampleMessage();
            }
        });

        document.getElementById('cloudEventModal').addEventListener('click', function(e) {
            if (e.target === this) {
                closeCloudEvent();
            }
        });

        document.getElementById('eventHandlerAnimation').addEventListener('click', function(e) {
            if (e.target === this) {
                closeEventHandlerAnimation();
            }
        });
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func getBaseHTML(title, content, extraJS string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - CloudEvents Explorer</title>
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
        .container { max-width: 1200px; margin: 0 auto; padding: 20px; }
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
        .btn-danger { background: #d93025; color: white; border-color: #d93025; }
        .btn-danger:hover { background: #c5221f; }
        .stats-bar { display: flex; gap: 24px; padding: 12px 20px; background: #f8f9fa; border-bottom: 1px solid #dadce0; font-size: 13px; }
        .stat { display: flex; align-items: center; gap: 6px; color: #5f6368; }
        .stat-value { color: #202124; font-weight: 600; }
        .message-list { display: flex; flex-direction: column; gap: 12px; }
        .message-card { background: white; border: 1px solid #dadce0; border-radius: 8px; overflow: hidden; transition: box-shadow 0.2s; }
        .message-card:hover { box-shadow: 0 1px 3px rgba(0,0,0,0.12), 0 1px 2px rgba(0,0,0,0.24); }
        .message-header {
            padding: 12px 16px;
            display: grid;
            grid-template-columns: auto 1fr auto auto;
            gap: 16px;
            align-items: center;
            cursor: pointer;
            user-select: none;
        }
        .message-header:hover { background: #f8f9fa; }
        .expand-icon { color: #5f6368; transition: transform 0.2s; font-size: 12px; }
        .expand-icon.expanded { transform: rotate(90deg); }
        .message-info { display: flex; flex-direction: column; gap: 4px; min-width: 0; }
        .message-type { font-size: 14px; font-weight: 500; color: #202124; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
        .message-subject { font-size: 12px; color: #5f6368; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
        .message-meta { display: flex; gap: 12px; font-size: 12px; color: #5f6368; }
        .message-time { font-size: 12px; color: #5f6368; white-space: nowrap; }
        .message-body { display: none; padding: 16px; border-top: 1px solid #dadce0; }
        .message-body.expanded { display: block; }
        .message-details {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 12px;
            padding: 12px;
            background: #f8f9fa;
            border-radius: 4px;
            margin-bottom: 12px;
            font-size: 13px;
        }
        .detail-item { display: flex; flex-direction: column; gap: 4px; }
        .detail-label { color: #5f6368; font-size: 11px; text-transform: uppercase; letter-spacing: 0.5px; }
        .detail-value { color: #202124; word-break: break-all; }
        .json-viewer { background: #f8f9fa; border: 1px solid #dadce0; border-radius: 4px; padding: 16px; overflow-x: auto; }
        .json-viewer pre { margin: 0; font-family: 'Monaco', 'Menlo', 'Consolas', monospace; font-size: 13px; line-height: 1.6; color: #202124; }
        .json-key { color: #1967d2; }
        .json-string { color: #188038; }
        .json-number { color: #1967d2; }
        .json-boolean { color: #d93025; }
        .json-null { color: #5f6368; }
        .empty-state { text-align: center; padding: 60px 20px; color: #5f6368; }
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
        .loading { text-align: center; padding: 40px; color: #5f6368; }
        .spinner {
            border: 3px solid #dadce0;
            border-top: 3px solid #1a73e8;
            border-radius: 50%%;
            width: 40px;
            height: 40px;
            animation: spin 1s linear infinite;
            margin: 0 auto 16px;
        }
        @keyframes spin {
            0%% { transform: rotate(0deg); }
            100%% { transform: rotate(360deg); }
        }
    </style>
</head>
<body>
    <div class="topbar">
        <a href="/" class="logo">CloudEvents Explorer</a>
        <a href="/" class="back-btn">← Back</a>
    </div>

    <div class="container">
        %s
        <div class="panel">
            <div class="stats-bar">
                <div class="stat">
                    <span class="stat-label">Total Messages:</span>
                    <span class="stat-value" id="totalMessages">0</span>
                </div>
                <div class="stat">
                    <span class="stat-label">Last Updated:</span>
                    <span class="stat-value" id="lastUpdated">Never</span>
                </div>
            </div>
            <div class="panel-body">
                <div id="messages"></div>
            </div>
        </div>
    </div>

    <div id="statusToast" class="status-toast"></div>

    <script>
        let messagesData = [];

        function showStatus(message, isError = false) {
            const toast = document.getElementById('statusToast');
            toast.textContent = message;
            toast.className = 'status-toast ' + (isError ? 'error' : 'success');
            toast.style.display = 'block';
            setTimeout(() => { toast.style.display = 'none'; }, 3000);
        }

        function syntaxHighlightJSON(json) {
            if (typeof json !== 'string') {
                json = JSON.stringify(json, null, 2);
            }
            json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
            return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function (match) {
                let cls = 'json-number';
                if (/^"/.test(match)) {
                    if (/:$/.test(match)) {
                        cls = 'json-key';
                    } else {
                        cls = 'json-string';
                    }
                } else if (/true|false/.test(match)) {
                    cls = 'json-boolean';
                } else if (/null/.test(match)) {
                    cls = 'json-null';
                }
                return '<span class="' + cls + '">' + match + '</span>';
            });
        }

        function toggleMessage(index) {
            const body = document.getElementById('msg-body-' + index);
            const icon = document.getElementById('msg-icon-' + index);
            if (body.classList.contains('expanded')) {
                body.classList.remove('expanded');
                icon.classList.remove('expanded');
            } else {
                body.classList.add('expanded');
                icon.classList.add('expanded');
            }
        }

        function renderMessages() {
            const container = document.getElementById('messages');
            if (messagesData.length === 0) {
                container.innerHTML = '<div class="empty-state"><div>No messages yet. Pull messages to get started.</div></div>';
                return;
            }
            let html = '<div class="message-list">';
            messagesData.forEach((msg, index) => {
                const time = new Date(msg.published).toLocaleString();
                const hasData = msg.data && Object.keys(msg.data).length > 0;
                const hasRawData = msg.rawData && msg.rawData.length > 0;

                html += '<div class="message-card">';
                html += '<div class="message-header" onclick="toggleMessage(' + index + ')">';
                html += '<span class="expand-icon" id="msg-icon-' + index + '">▶</span>';
                html += '<div class="message-info">';
                html += '<div class="message-type">' + (msg.type || msg.subject || 'Message') + '</div>';
                html += '<div class="message-subject">' + (msg.subject || msg.id || 'No subject') + '</div>';
                html += '</div>';
                html += '<div class="message-meta">';
                if (msg.id) html += '<span>ID: ' + msg.id + '</span>';
                if (msg.source) html += '<span>Source: ' + msg.source + '</span>';
                html += '</div>';
                html += '<div class="message-time">' + time + '</div>';
                html += '</div>';

                html += '<div class="message-body" id="msg-body-' + index + '">';
                html += '<div class="message-details">';
                if (msg.id) html += '<div class="detail-item"><div class="detail-label">Message ID</div><div class="detail-value">' + msg.id + '</div></div>';
                if (msg.type) html += '<div class="detail-item"><div class="detail-label">Type</div><div class="detail-value">' + msg.type + '</div></div>';
                if (msg.subject) html += '<div class="detail-item"><div class="detail-label">Subject</div><div class="detail-value">' + msg.subject + '</div></div>';
                if (msg.source) html += '<div class="detail-item"><div class="detail-label">Source</div><div class="detail-value">' + msg.source + '</div></div>';
                if (msg.schema) html += '<div class="detail-item"><div class="detail-label">Schema</div><div class="detail-value">' + msg.schema + '</div></div>';
                html += '<div class="detail-item"><div class="detail-label">Published</div><div class="detail-value">' + msg.published + '</div></div>';
                html += '</div>';

                if (hasData) {
                    html += '<div style="position: relative;">';
                    html += '<button onclick="copyMessageData(' + index + ')" style="position: absolute; top: 8px; right: 8px; background: #1a73e8; color: white; border: none; padding: 6px 12px; border-radius: 4px; cursor: pointer; font-size: 12px; font-weight: 500;">Copy JSON</button>';
                    html += '<div class="json-viewer"><pre>' + syntaxHighlightJSON(msg.data) + '</pre></div>';
                    html += '</div>';
                } else if (hasRawData) {
                    html += '<div style="position: relative;">';
                    html += '<button onclick="copyMessageData(' + index + ')" style="position: absolute; top: 8px; right: 8px; background: #1a73e8; color: white; border: none; padding: 6px 12px; border-radius: 4px; cursor: pointer; font-size: 12px; font-weight: 500;">Copy JSON</button>';
                    html += '<div class="json-viewer"><pre>' + syntaxHighlightJSON(msg.rawData) + '</pre></div>';
                    html += '</div>';
                }

                html += '</div>';
                html += '</div>';
            });
            html += '</div>';

            container.innerHTML = html;
            document.getElementById('totalMessages').textContent = messagesData.length;
            document.getElementById('lastUpdated').textContent = new Date().toLocaleTimeString();
        }

        function clearAllMessages() {
            if (confirm('Are you sure you want to clear all messages?')) {
                messagesData = [];
                renderMessages();
                showStatus('All messages cleared');
            }
        }

        function copyMessageData(index) {
            const msg = messagesData[index];
            const jsonData = msg.data || msg.rawData;
            const jsonString = JSON.stringify(jsonData, null, 2);

            const textarea = document.createElement('textarea');
            textarea.value = jsonString;
            document.body.appendChild(textarea);
            textarea.select();
            document.execCommand('copy');
            document.body.removeChild(textarea);

            showStatus('Message JSON copied to clipboard!');
        }

        %s

        renderMessages();
    </script>
</body>
</html>`, title, content, extraJS)
}

func handleGetConfigs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getConfig())
}

func handleSavePubSubConfig(w http.ResponseWriter, r *http.Request) {
	var newConfig PubSubConfig
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	found := false
	for i, cfg := range config.PubSubConfigs {
		if cfg.Name == newConfig.Name {
			config.PubSubConfigs[i] = newConfig
			found = true
			break
		}
	}
	if !found {
		config.PubSubConfigs = append(config.PubSubConfigs, newConfig)
	}
	mu.Unlock()

	if err := saveConfig(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func handleSaveKafkaConfig(w http.ResponseWriter, r *http.Request) {
	var newConfig KafkaConfig
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	found := false
	for i, cfg := range config.KafkaConfigs {
		if cfg.Name == newConfig.Name {
			config.KafkaConfigs[i] = newConfig
			found = true
			break
		}
	}
	if !found {
		config.KafkaConfigs = append(config.KafkaConfigs, newConfig)
	}
	mu.Unlock()

	if err := saveConfig(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func handlePullPubSub(w http.ResponseWriter, r *http.Request) {
	var params struct {
		EmulatorHost   string `json:"emulatorHost"`
		ProjectID      string `json:"projectId"`
		SubscriptionID string `json:"subscriptionId"`
		MaxMessages    int    `json:"maxMessages"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	os.Setenv("PUBSUB_EMULATOR_HOST", params.EmulatorHost)

	client, err := pubsub.NewClient(ctx, params.ProjectID,
		option.WithEndpoint(params.EmulatorHost),
		option.WithoutAuthentication(),
	)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to create client: %v", err)})
		return
	}
	defer client.Close()

	subscription := client.Subscription(params.SubscriptionID)

	messages := []CloudEvent{}
	var msgMu sync.Mutex

	receiveCtx, receiveCancel := context.WithTimeout(ctx, 5*time.Second)
	defer receiveCancel()

	err = subscription.Receive(receiveCtx, func(ctx context.Context, msg *pubsub.Message) {
		event := CloudEvent{
			ID:        msg.ID,
			Type:      msg.Attributes["ce-type"],
			Subject:   msg.Attributes["ce-subject"],
			Source:    msg.Attributes["ce-source"],
			Schema:    msg.Attributes["ce-dataschema"],
			Published: msg.PublishTime.Format(time.RFC3339),
			Timestamp: msg.PublishTime.Unix(),
		}

		if len(msg.Data) > 0 {
			var data map[string]interface{}
			if err := json.Unmarshal(msg.Data, &data); err == nil {
				event.Data = data
			}
		}

		msgMu.Lock()
		messages = append(messages, event)
		msgMu.Unlock()

		msg.Ack()

		if len(messages) >= params.MaxMessages {
			receiveCancel()
		}
	})

	if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to receive: %v", err)})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"messages": messages,
		"count":    len(messages),
	})
}

func handlePullKafka(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Brokers        string `json:"brokers"`
		Topic          string `json:"topic"`
		ConsumerGroup  string `json:"consumerGroup"`
		SchemaRegistry string `json:"schemaRegistry"`
		MaxMessages    int    `json:"maxMessages"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": params.Brokers,
		"group.id":          params.ConsumerGroup,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to create consumer: %v", err)})
		return
	}
	defer c.Close()

	err = c.Subscribe(params.Topic, nil)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to subscribe: %v", err)})
		return
	}


	messages := []CloudEvent{}
	timeout := time.After(5 * time.Second)

	for len(messages) < params.MaxMessages {
		select {
		case <-timeout:
			goto done
		default:
			msg, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				continue
			}

			event := CloudEvent{
				ID:        fmt.Sprintf("%s-%d-%d", params.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset),
				Subject:   params.Topic,
				Published: msg.Timestamp.Format(time.RFC3339),
				Timestamp: msg.Timestamp.Unix(),
			}

			// Decode Avro message using schema registry
			if params.SchemaRegistry != "" && len(msg.Value) > 5 {
				decodedData, err := decodeAvroMessage(msg.Value, params.SchemaRegistry)
				if err == nil && decodedData != nil {
					event.Data = decodedData
				} else {
					// Fall back to raw data if decoding fails
					event.RawData = string(msg.Value)
				}
			} else {
				// Try plain JSON first
				var data map[string]interface{}
				if err := json.Unmarshal(msg.Value, &data); err == nil {
					event.Data = data
				} else {
					event.RawData = string(msg.Value)
				}
			}

			messages = append(messages, event)
		}
	}

done:
	// Reverse messages so newest appears first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"messages": messages,
		"count":    len(messages),
	})
}

func handlePublishKafka(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Brokers        string                 `json:"brokers"`
		Topic          string                 `json:"topic"`
		SchemaRegistry string                 `json:"schemaRegistry"`
		Message        map[string]interface{} `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert message to JSON bytes
	messageJSON, err := json.Marshal(params.Message)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to marshal message: %v", err)})
		return
	}

	// Create Kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": params.Brokers,
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to create producer: %v", err)})
		return
	}
	defer p.Close()

	// Encode to Avro if schema registry is configured
	var messageBytes []byte
	if params.SchemaRegistry != "" {
		// Fetch schema from registry (using subject for Unica events)
		schemaURL := fmt.Sprintf("%s/subjects/au.data.unica.comms.event/versions/latest", params.SchemaRegistry)
		resp, err := http.Get(schemaURL)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to fetch schema: %v", err)})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Schema registry returned status %d", resp.StatusCode)})
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to read schema response: %v", err)})
			return
		}

		var schemaResp struct {
			Schema string `json:"schema"`
			ID     int    `json:"id"`
		}
		if err := json.Unmarshal(body, &schemaResp); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to parse schema response: %v", err)})
			return
		}

		// Create Avro codec
		codec, err := goavro.NewCodec(schemaResp.Schema)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to create codec: %v", err)})
			return
		}

		// Encode message to Avro binary
		avroBinary, err := codec.BinaryFromNative(nil, params.Message)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to encode to Avro: %v", err)})
			return
		}

		// Prepend magic byte (0) and schema ID (4 bytes, big-endian)
		messageBytes = make([]byte, 5+len(avroBinary))
		messageBytes[0] = 0 // Magic byte
		binary.BigEndian.PutUint32(messageBytes[1:5], uint32(schemaResp.ID))
		copy(messageBytes[5:], avroBinary)
	} else {
		// Plain JSON
		messageBytes = messageJSON
	}

	// Publish to Kafka
	deliveryChan := make(chan kafka.Event)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &params.Topic, Partition: kafka.PartitionAny},
		Value:          messageBytes,
	}, deliveryChan)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to produce message: %v", err)})
		return
	}

	// Wait for delivery report
	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Delivery failed: %v", m.TopicPartition.Error)})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"partition": m.TopicPartition.Partition,
		"offset":    m.TopicPartition.Offset,
	})
}

func main() {
	if err := loadConfig(); err != nil {
		log.Printf("Warning: Failed to load config: %v", err)
	}

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/pubsub", handlePubSub)
	http.HandleFunc("/kafka", handleKafka)
	http.HandleFunc("/flow-diagram", handleFlowDiagram)
	http.HandleFunc("/api/configs", handleGetConfigs)
	http.HandleFunc("/api/pubsub/configs", handleSavePubSubConfig)
	http.HandleFunc("/api/kafka/configs", handleSaveKafkaConfig)
	http.HandleFunc("/api/pubsub/pull", handlePullPubSub)
	http.HandleFunc("/api/kafka/pull", handlePullKafka)
	http.HandleFunc("/api/kafka/publish", handlePublishKafka)

	port := "8888"
	log.Printf("🚀 CloudEvents Explorer starting on http://localhost:%s", port)
	log.Printf("📝 Configuration file: %s", configFile)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}