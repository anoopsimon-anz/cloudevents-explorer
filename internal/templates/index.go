package templates

const Index = `<!DOCTYPE html>
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

        ` + Base64ModalJS + `
    </script>
</body>
</html>`
