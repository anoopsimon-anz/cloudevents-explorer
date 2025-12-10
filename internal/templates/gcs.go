package templates

const GCS = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GCS Browser - Testing Studio</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", sans-serif;
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
            max-width: 1400px;
            margin: 0 auto;
            padding: 24px;
        }
        .breadcrumb {
            background: white;
            border: 1px solid #dadce0;
            border-radius: 8px;
            padding: 12px 16px;
            margin-bottom: 16px;
            font-size: 14px;
            color: #5f6368;
        }
        .breadcrumb-item {
            display: inline;
            cursor: pointer;
            color: #1a73e8;
        }
        .breadcrumb-item:hover {
            text-decoration: underline;
        }
        .breadcrumb-separator {
            margin: 0 8px;
            color: #5f6368;
        }
        .main-content {
            display: grid;
            grid-template-columns: 300px 1fr;
            gap: 16px;
        }
        .sidebar {
            background: white;
            border: 1px solid #dadce0;
            border-radius: 8px;
            padding: 16px;
            height: fit-content;
        }
        .sidebar-title {
            font-size: 14px;
            font-weight: 500;
            color: #202124;
            margin-bottom: 12px;
        }
        .bucket-list {
            list-style: none;
        }
        .bucket-item {
            padding: 8px 12px;
            margin-bottom: 4px;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
            color: #202124;
            transition: background 0.2s;
        }
        .bucket-item:hover {
            background: #f1f3f4;
        }
        .bucket-item.active {
            background: #e8f0fe;
            color: #1967d2;
            font-weight: 500;
        }
        .content-area {
            background: white;
            border: 1px solid #dadce0;
            border-radius: 8px;
            padding: 16px;
            min-height: 500px;
        }
        .empty-state {
            text-align: center;
            padding: 64px 24px;
            color: #5f6368;
        }
        .empty-state-icon {
            font-size: 48px;
            margin-bottom: 16px;
        }
        .file-list {
            list-style: none;
        }
        .file-item {
            display: flex;
            align-items: center;
            padding: 12px;
            border-bottom: 1px solid #f1f3f4;
            cursor: pointer;
            transition: background 0.2s;
        }
        .file-item:hover {
            background: #f8f9fa;
        }
        .file-icon {
            font-size: 20px;
            margin-right: 12px;
            width: 24px;
            text-align: center;
        }
        .file-info {
            flex: 1;
        }
        .file-name {
            font-size: 14px;
            color: #202124;
            margin-bottom: 2px;
        }
        .file-meta {
            font-size: 12px;
            color: #5f6368;
        }
        .file-actions {
            display: flex;
            gap: 8px;
        }
        .btn {
            padding: 6px 12px;
            border: 1px solid #dadce0;
            border-radius: 4px;
            background: white;
            color: #202124;
            font-size: 12px;
            cursor: pointer;
            transition: all 0.2s;
        }
        .btn:hover {
            background: #f8f9fa;
            border-color: #1a73e8;
        }
        .btn-primary {
            background: #1a73e8;
            color: white;
            border-color: #1a73e8;
        }
        .btn-primary:hover {
            background: #1765cc;
        }
        .modal {
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: rgba(0,0,0,0.5);
            z-index: 1000;
            align-items: center;
            justify-content: center;
        }
        .modal.active {
            display: flex;
        }
        .modal-content {
            background: white;
            border-radius: 8px;
            width: 90%;
            max-width: 800px;
            max-height: 80vh;
            display: flex;
            flex-direction: column;
        }
        .modal-header {
            padding: 16px 24px;
            border-bottom: 1px solid #dadce0;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .modal-title {
            font-size: 18px;
            font-weight: 500;
            color: #202124;
        }
        .modal-close {
            background: none;
            border: none;
            font-size: 24px;
            color: #5f6368;
            cursor: pointer;
            padding: 0;
            width: 32px;
            height: 32px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .modal-close:hover {
            background: #f1f3f4;
        }
        .modal-body {
            padding: 24px;
            overflow-y: auto;
        }
        .file-preview {
            background: #f8f9fa;
            border: 1px solid #dadce0;
            border-radius: 4px;
            padding: 16px;
            font-family: 'Monaco', 'Menlo', monospace;
            font-size: 12px;
            white-space: pre-wrap;
            word-wrap: break-word;
            max-height: 500px;
            overflow-y: auto;
        }
        .loading {
            text-align: center;
            padding: 32px;
            color: #5f6368;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>📦 GCS Browser</h1>
        <a href="/" class="back-link">← Back to Home</a>
    </div>

    <div class="container">
        <div class="breadcrumb" id="breadcrumb">
            <span class="breadcrumb-item" onclick="navigateToRoot()">Buckets</span>
        </div>

        <div class="main-content">
            <div class="sidebar">
                <div class="sidebar-title">Buckets</div>
                <ul class="bucket-list" id="bucketList">
                    <li class="loading">Loading buckets...</li>
                </ul>
            </div>

            <div class="content-area" id="contentArea">
                <div class="empty-state">
                    <div class="empty-state-icon">📦</div>
                    <div>Select a bucket to view its contents</div>
                </div>
            </div>
        </div>
    </div>

    <div class="modal" id="previewModal">
        <div class="modal-content">
            <div class="modal-header">
                <div class="modal-title" id="previewTitle">File Preview</div>
                <button class="modal-close" onclick="closePreview()">×</button>
            </div>
            <div class="modal-body">
                <div class="file-preview" id="previewContent">Loading...</div>
            </div>
        </div>
    </div>

    <script>
        let currentBucket = null;
        let currentPrefix = '';

        async function loadBuckets() {
            try {
                const response = await fetch('/api/gcs/buckets');
                const data = await response.json();

                const bucketList = document.getElementById('bucketList');
                if (data.buckets && data.buckets.length > 0) {
                    bucketList.innerHTML = data.buckets.map(bucket =>
                        '<li class="bucket-item" onclick="selectBucket(\'' + bucket.name + '\')">' + bucket.name + '</li>'
                    ).join('');
                } else {
                    bucketList.innerHTML = '<li style="color: #5f6368; padding: 8px;">No buckets found</li>';
                }
            } catch (error) {
                console.error('Failed to load buckets:', error);
                document.getElementById('bucketList').innerHTML = '<li style="color: #ea4335; padding: 8px;">Error loading buckets</li>';
            }
        }

        function selectBucket(bucketName) {
            currentBucket = bucketName;
            currentPrefix = '';

            // Update active state
            document.querySelectorAll('.bucket-item').forEach(item => {
                item.classList.remove('active');
                if (item.textContent === bucketName) {
                    item.classList.add('active');
                }
            });

            updateBreadcrumb();
            loadObjects();
        }

        async function loadObjects() {
            if (!currentBucket) return;

            document.getElementById('contentArea').innerHTML = '<div class="loading">Loading objects...</div>';

            try {
                const url = '/api/gcs/objects?bucket=' + encodeURIComponent(currentBucket) +
                           (currentPrefix ? '&prefix=' + encodeURIComponent(currentPrefix) : '');
                const response = await fetch(url);
                const data = await response.json();

                if (data.prefixes || (data.items && data.items.length > 0)) {
                    let html = '<ul class="file-list">';

                    // Show prefixes (folders)
                    if (data.prefixes) {
                        data.prefixes.forEach(prefix => {
                            const folderName = prefix.replace(currentPrefix, '').replace('/', '');
                            html += '<li class="file-item" onclick="navigateToPrefix(\'' + prefix + '\')">';
                            html += '<div class="file-icon">📁</div>';
                            html += '<div class="file-info">';
                            html += '<div class="file-name">' + folderName + '</div>';
                            html += '<div class="file-meta">Folder</div>';
                            html += '</div>';
                            html += '</li>';
                        });
                    }

                    // Show objects (files)
                    if (data.items) {
                        data.items.forEach(item => {
                            const fileName = item.name.replace(currentPrefix, '');
                            if (fileName) { // Skip if it's the prefix itself
                                const sizeKB = Math.round(item.size / 1024);
                                html += '<li class="file-item">';
                                html += '<div class="file-icon">📄</div>';
                                html += '<div class="file-info">';
                                html += '<div class="file-name">' + fileName + '</div>';
                                html += '<div class="file-meta">' + sizeKB + ' KB</div>';
                                html += '</div>';
                                html += '<div class="file-actions">';
                                html += '<button class="btn" onclick="previewFile(\'' + item.name + '\')">Preview</button>';
                                html += '<button class="btn btn-primary" onclick="downloadFile(\'' + item.name + '\')">Download</button>';
                                html += '</div>';
                                html += '</li>';
                            }
                        });
                    }

                    html += '</ul>';
                    document.getElementById('contentArea').innerHTML = html;
                } else {
                    document.getElementById('contentArea').innerHTML =
                        '<div class="empty-state"><div class="empty-state-icon">📭</div><div>This bucket is empty</div></div>';
                }
            } catch (error) {
                console.error('Failed to load objects:', error);
                document.getElementById('contentArea').innerHTML =
                    '<div class="empty-state"><div class="empty-state-icon">⚠️</div><div>Error loading objects</div></div>';
            }
        }

        function navigateToPrefix(prefix) {
            currentPrefix = prefix;
            updateBreadcrumb();
            loadObjects();
        }

        function navigateToRoot() {
            currentPrefix = '';
            updateBreadcrumb();
            if (currentBucket) {
                loadObjects();
            }
        }

        function updateBreadcrumb() {
            let breadcrumb = '<span class="breadcrumb-item" onclick="navigateToRoot()">Buckets</span>';

            if (currentBucket) {
                breadcrumb += '<span class="breadcrumb-separator">/</span>';
                breadcrumb += '<span class="breadcrumb-item" onclick="selectBucket(\'' + currentBucket + '\')">' + currentBucket + '</span>';
            }

            if (currentPrefix) {
                const parts = currentPrefix.split('/').filter(p => p);
                let path = '';
                parts.forEach(part => {
                    path += part + '/';
                    breadcrumb += '<span class="breadcrumb-separator">/</span>';
                    breadcrumb += '<span class="breadcrumb-item" onclick="navigateToPrefix(\'' + path + '\')">' + part + '</span>';
                });
            }

            document.getElementById('breadcrumb').innerHTML = breadcrumb;
        }

        async function previewFile(objectName) {
            document.getElementById('previewModal').classList.add('active');
            document.getElementById('previewTitle').textContent = objectName;
            document.getElementById('previewContent').textContent = 'Loading...';

            try {
                const url = '/api/gcs/object/content?bucket=' + encodeURIComponent(currentBucket) +
                           '&object=' + encodeURIComponent(objectName);
                const response = await fetch(url);
                const data = await response.json();

                document.getElementById('previewContent').textContent = data.content;
            } catch (error) {
                console.error('Failed to preview file:', error);
                document.getElementById('previewContent').textContent = 'Error loading file content';
            }
        }

        function closePreview() {
            document.getElementById('previewModal').classList.remove('active');
        }

        async function downloadFile(objectName) {
            try {
                const url = '/api/gcs/object/download?bucket=' + encodeURIComponent(currentBucket) +
                           '&object=' + encodeURIComponent(objectName);
                window.location.href = url;
            } catch (error) {
                console.error('Failed to download file:', error);
                alert('Error downloading file');
            }
        }

        // Close modal when clicking outside
        document.getElementById('previewModal').addEventListener('click', function(e) {
            if (e.target === this) {
                closePreview();
            }
        });

        // Load buckets on page load
        loadBuckets();
    </script>
</body>
</html>`
