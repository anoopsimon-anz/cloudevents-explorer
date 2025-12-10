package templates

import "fmt"

func GetBaseHTML(title, content, extraJS string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - Testing Studio</title>
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
        <a href="/" class="logo">Testing Studio</a>
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
