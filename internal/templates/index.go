package templates

const Index = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Testing Studio</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Google Sans', 'Product Sans', -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", sans-serif;
            background: #fafafa;
            color: #202124;
            min-height: 100vh;
            padding: 80px 24px 24px 24px;
        }
        .landing {
            max-width: 1000px;
            margin: 0 auto;
        }
        .hero {
            text-align: center;
            margin-bottom: 64px;
        }
        h1 {
            font-size: 48px;
            color: #202124;
            margin-bottom: 12px;
            font-weight: 400;
            letter-spacing: -0.5px;
        }
        .subtitle {
            font-size: 15px;
            color: #5f6368;
            font-weight: 400;
        }
        .options {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
            gap: 24px;
            margin-top: 48px;
        }
        .option-card {
            background: white;
            border: 1px solid #e8eaed;
            border-radius: 12px;
            padding: 32px 24px;
            cursor: pointer;
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
            text-decoration: none;
            color: inherit;
            display: flex;
            flex-direction: column;
            position: relative;
            overflow: hidden;
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
        .status-indicators {
            position: absolute;
            top: 20px;
            right: 20px;
            display: flex;
            flex-direction: column;
            gap: 8px;
        }
        .status-indicator {
            display: flex;
            align-items: center;
            gap: 6px;
            font-size: 12px;
            color: #5f6368;
            font-weight: 500;
        }
        .status-dot {
            width: 10px;
            height: 10px;
            border-radius: 50%;
            box-shadow: 0 0 0 2px rgba(0,0,0,0.1);
        }
        .status-dot.green {
            background: #34a853;
        }
        .status-dot.red {
            background: #ea4335;
        }
    </style>
</head>
<body>
    <div class="status-indicators">
        <div class="status-indicator" id="dockerStatus">
            <div class="status-dot red"></div>
            <span>Docker</span>
        </div>
        <div class="status-indicator" id="gcloudStatus">
            <div class="status-dot red"></div>
            <span id="gcloudText">GCloud</span>
        </div>
    </div>
    <div class="landing">
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 48px;">
            <div>
                <h1>Testing Studio</h1>
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

            <a href="/rest-client" class="option-card">
                <div class="option-title">REST Client</div>
                <div class="option-desc">Send HTTP requests with custom headers, body, and TLS certificates</div>
                <span class="badge">API Testing</span>
            </a>

            <a href="/gcs" class="option-card">
                <div class="option-title">GCS Browser</div>
                <div class="option-desc">Browse Google Cloud Storage buckets and files with preview and download</div>
                <span class="badge">Storage</span>
            </a>
        </div>
    </div>

    ` + Base64Modal + `

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

        // Check Docker status
        async function checkDockerStatus() {
            try {
                const response = await fetch('/api/docker/status');
                const data = await response.json();
                const statusIndicator = document.getElementById('dockerStatus');
                const dot = statusIndicator.querySelector('.status-dot');

                if (data.running) {
                    dot.classList.remove('red');
                    dot.classList.add('green');
                } else {
                    dot.classList.remove('green');
                    dot.classList.add('red');
                }
            } catch (error) {
                console.error('Failed to check Docker status:', error);
            }
        }

        // Check Docker status on page load
        checkDockerStatus();

        // Check Docker status every 10 seconds
        setInterval(checkDockerStatus, 10000);

        // Check GCloud authentication status
        async function checkGCloudStatus() {
            try {
                const response = await fetch('/api/gcloud/status');
                const data = await response.json();
                const statusIndicator = document.getElementById('gcloudStatus');
                const dot = statusIndicator.querySelector('.status-dot');
                const textSpan = document.getElementById('gcloudText');

                if (data.authenticated) {
                    dot.classList.remove('red');
                    dot.classList.add('green');

                    // Show last login time if available
                    if (data.lastLoginTime) {
                        textSpan.textContent = 'GCloud: ' + data.lastLoginTime;
                        textSpan.title = 'Account: ' + data.account;
                    } else {
                        textSpan.textContent = 'GCloud: Active';
                        textSpan.title = 'Account: ' + data.account;
                    }
                } else {
                    dot.classList.remove('green');
                    dot.classList.add('red');
                    textSpan.textContent = 'GCloud: Not authenticated';
                    textSpan.title = 'Run: gcloud auth login';
                }
            } catch (error) {
                console.error('Failed to check GCloud status:', error);
            }
        }

        // Check GCloud status on page load
        checkGCloudStatus();

        // Check GCloud status every 30 seconds
        setInterval(checkGCloudStatus, 30000);

        ` + Base64ModalJS + `
    </script>
</body>
</html>`
