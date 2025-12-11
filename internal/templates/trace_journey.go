package templates

const TraceJourney = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Trace Journey - Testing Studio</title>
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
            max-width: 1200px;
            margin: 0 auto;
            padding: 24px;
        }
        .search-panel {
            background: white;
            border: 1px solid #dadce0;
            border-radius: 8px;
            padding: 24px;
            margin-bottom: 24px;
        }
        .search-title {
            font-size: 16px;
            font-weight: 500;
            margin-bottom: 16px;
        }
        .input-group {
            margin-bottom: 16px;
        }
        .input-label {
            font-size: 14px;
            color: #5f6368;
            margin-bottom: 8px;
            display: block;
        }
        .trace-input {
            width: 100%;
            padding: 12px;
            border: 1px solid #dadce0;
            border-radius: 4px;
            font-size: 14px;
            font-family: monospace;
        }
        .trace-input:focus {
            outline: none;
            border-color: #1a73e8;
        }
        .containers-section {
            margin-bottom: 16px;
        }
        .containers-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
            gap: 8px;
            margin-top: 8px;
        }
        .container-checkbox {
            display: flex;
            align-items: center;
            padding: 8px;
            border: 1px solid #dadce0;
            border-radius: 4px;
            cursor: pointer;
            transition: background 0.2s;
        }
        .container-checkbox:hover {
            background: #f8f9fa;
        }
        .container-checkbox input {
            margin-right: 8px;
        }
        .container-checkbox label {
            font-size: 12px;
            cursor: pointer;
            flex: 1;
            font-family: monospace;
        }
        .search-btn {
            background: #1a73e8;
            color: white;
            border: none;
            padding: 12px 32px;
            border-radius: 4px;
            font-size: 14px;
            font-weight: 500;
            cursor: pointer;
            transition: background 0.2s;
        }
        .search-btn:hover {
            background: #1765cc;
        }
        .search-btn:disabled {
            background: #dadce0;
            cursor: not-allowed;
        }
        .results-panel {
            background: white;
            border: 1px solid #dadce0;
            border-radius: 8px;
            padding: 24px;
            display: none;
        }
        .results-panel.visible {
            display: block;
        }
        .results-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 24px;
        }
        .results-title {
            font-size: 16px;
            font-weight: 500;
        }
        .results-count {
            font-size: 14px;
            color: #5f6368;
        }
        .error-count {
            color: #d93025;
            font-weight: 500;
        }
        .container-filters {
            display: flex;
            gap: 8px;
            margin-bottom: 16px;
            flex-wrap: wrap;
            align-items: center;
        }
        .filter-label {
            font-size: 12px;
            color: #5f6368;
            font-weight: 500;
        }
        .filter-pill {
            display: inline-block;
            background: #e8f0fe;
            color: #1967d2;
            padding: 6px 12px;
            border-radius: 16px;
            font-size: 11px;
            font-weight: 500;
            font-family: monospace;
            cursor: pointer;
            border: 2px solid transparent;
            transition: all 0.2s;
        }
        .filter-pill:hover {
            border-color: #1967d2;
        }
        .filter-pill.inactive {
            background: #f1f3f4;
            color: #5f6368;
            opacity: 0.5;
        }
        .timeline {
            position: relative;
            padding-left: 40px;
        }
        .timeline::before {
            content: '';
            position: absolute;
            left: 19px;
            top: 0;
            bottom: 0;
            width: 2px;
            background: #dadce0;
        }
        .timeline-item {
            position: relative;
            margin-bottom: 24px;
        }
        .timeline-dot {
            position: absolute;
            left: -28px;
            top: 4px;
            width: 20px;
            height: 20px;
            border-radius: 50%;
            background: #1a73e8;
            border: 3px solid white;
            box-shadow: 0 0 0 2px #1a73e8;
        }
        .timeline-dot.error {
            background: #d93025;
            box-shadow: 0 0 0 2px #d93025;
        }
        .timeline-content {
            background: #f8f9fa;
            border-radius: 8px;
            padding: 12px 16px;
            border-left: 3px solid #1a73e8;
        }
        .timeline-content.error {
            background: #fce8e6;
            border-left-color: #d93025;
        }
        .timeline-time {
            font-size: 11px;
            color: #5f6368;
            font-family: monospace;
            margin-bottom: 4px;
        }
        .timeline-container {
            font-size: 12px;
            margin-bottom: 8px;
        }
        .container-pill {
            display: inline-block;
            background: #e8f0fe;
            color: #1967d2;
            padding: 4px 12px;
            border-radius: 12px;
            font-size: 11px;
            font-weight: 500;
            font-family: monospace;
        }
        .timeline-action {
            font-size: 13px;
            font-weight: 500;
            color: #202124;
            margin-bottom: 4px;
        }
        .timeline-message {
            font-size: 13px;
            color: #5f6368;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
            max-width: 100%;
            cursor: pointer;
            transition: all 0.2s;
        }
        .timeline-message.expanded {
            white-space: pre-wrap;
            word-break: break-word;
        }
        .timeline-message:hover {
            color: #1a73e8;
        }
        .timeline-span {
            font-size: 11px;
            color: #5f6368;
            font-family: monospace;
            margin-top: 4px;
        }
        .loading {
            text-align: center;
            padding: 64px;
            color: #5f6368;
        }
        .empty-state {
            text-align: center;
            padding: 64px;
            color: #5f6368;
        }
        .error-badge {
            display: inline-block;
            background: #d93025;
            color: white;
            padding: 2px 8px;
            border-radius: 4px;
            font-size: 11px;
            font-weight: 500;
            margin-left: 8px;
        }
        .select-all {
            font-size: 12px;
            color: #1a73e8;
            cursor: pointer;
            margin-left: 16px;
        }
        .select-all:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>🔍 Trace Journey Viewer</h1>
        <a href="/" class="back-link">← Back to Home</a>
    </div>

    <div class="container">
        <div class="search-panel">
            <div class="search-title">Search for Trace ID</div>

            <div class="input-group">
                <label class="input-label">Trace ID</label>
                <input type="text"
                       class="trace-input"
                       id="traceIdInput"
                       placeholder="e.g., a4a0f2b6b62757a9503c03d99d77d5b2"
                       value="a4a0f2b6b62757a9503c03d99d77d5b2">
            </div>

            <div class="containers-section">
                <label class="input-label">
                    Containers to Search
                    <span class="select-all" onclick="selectAll()">Select All</span>
                    <span class="select-all" onclick="deselectAll()">Deselect All</span>
                </label>
                <div class="containers-grid" id="containersGrid">
                    <!-- Containers will be loaded dynamically -->
                </div>
            </div>

            <button class="search-btn" onclick="searchTrace()">Search Trace Journey</button>
        </div>

        <div class="results-panel" id="resultsPanel">
            <div class="results-header">
                <div class="results-title">Trace Journey</div>
                <div class="results-count" id="resultsCount">0 events</div>
            </div>
            <div class="container-filters" id="containerFilters">
                <!-- Container filter pills will be added here -->
            </div>
            <div class="timeline" id="timeline">
                <!-- Timeline items will be added here -->
            </div>
        </div>
    </div>

    <script>
        // Known devstack containers
        const knownContainers = [
            'devstack-component_worker-1',
            'devstack-component_eventhandler-1',
            'devstack-dep_flimflam-1',
            'devstack-dep_recorder-1',
            'devstack-datagen-1',
            'devstack-deps_idp-1',
            'dep_redpanda',
            'devstack-dep_pubsub-1'
        ];

        function initContainers() {
            const grid = document.getElementById('containersGrid');
            knownContainers.forEach(container => {
                const div = document.createElement('div');
                div.className = 'container-checkbox';
                div.innerHTML = ` + "`" + `
                    <input type="checkbox" id="cb_${container}" value="${container}" checked>
                    <label for="cb_${container}">${container}</label>
                ` + "`" + `;
                grid.appendChild(div);
            });
        }

        function selectAll() {
            document.querySelectorAll('#containersGrid input[type="checkbox"]').forEach(cb => {
                cb.checked = true;
            });
        }

        function deselectAll() {
            document.querySelectorAll('#containersGrid input[type="checkbox"]').forEach(cb => {
                cb.checked = false;
            });
        }

        async function searchTrace() {
            const traceId = document.getElementById('traceIdInput').value.trim();
            if (!traceId) {
                alert('Please enter a trace ID');
                return;
            }

            const selectedContainers = Array.from(
                document.querySelectorAll('#containersGrid input[type="checkbox"]:checked')
            ).map(cb => cb.value);

            if (selectedContainers.length === 0) {
                alert('Please select at least one container');
                return;
            }

            const resultsPanel = document.getElementById('resultsPanel');
            const timeline = document.getElementById('timeline');

            resultsPanel.classList.add('visible');
            timeline.innerHTML = '<div class="loading">Searching containers...</div>';

            try {
                const response = await fetch('/api/trace/search', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        traceId: traceId,
                        containers: selectedContainers
                    })
                });

                const data = await response.json();

                if (data.events && data.events.length > 0) {
                    displayTimeline(data.events);
                } else {
                    timeline.innerHTML = '<div class="empty-state">No logs found for this trace ID</div>';
                }
            } catch (error) {
                console.error('Search failed:', error);
                timeline.innerHTML = '<div class="empty-state">Error searching logs</div>';
            }
        }

        let allEvents = [];
        let activeContainers = new Set();

        function displayTimeline(events) {
            allEvents = events;
            const timeline = document.getElementById('timeline');
            const resultsCount = document.getElementById('resultsCount');
            const containerFilters = document.getElementById('containerFilters');

            // Get unique containers
            const containers = [...new Set(events.map(e => e.container))];
            activeContainers = new Set(containers);

            // Create filter pills
            if (containers.length > 1) {
                containerFilters.innerHTML = '<span class="filter-label">Filter:</span>' +
                    containers.map(container => {
                        const shortName = container.replace('devstack-', '').replace('-1', '');
                        return ` + "`<span class=\"filter-pill\" data-container=\"${container}\" onclick=\"toggleContainerFilter('${container}')\">#${shortName}</span>`" + `;
                    }).join('');
            } else {
                containerFilters.innerHTML = '';
            }

            renderFilteredTimeline();
        }

        function toggleContainerFilter(container) {
            if (activeContainers.has(container)) {
                activeContainers.delete(container);
            } else {
                activeContainers.add(container);
            }

            // Update pill styles
            document.querySelectorAll('.filter-pill').forEach(pill => {
                const pillContainer = pill.getAttribute('data-container');
                if (activeContainers.has(pillContainer)) {
                    pill.classList.remove('inactive');
                } else {
                    pill.classList.add('inactive');
                }
            });

            renderFilteredTimeline();
        }

        function renderFilteredTimeline() {
            const timeline = document.getElementById('timeline');
            const resultsCount = document.getElementById('resultsCount');

            // Filter events by active containers
            const filteredEvents = allEvents.filter(event => activeContainers.has(event.container));

            // Count errors in filtered events
            const errorCount = filteredEvents.filter(event =>
                event.severity === 'error' ||
                event.body.toLowerCase().includes('error') ||
                event.body.toLowerCase().includes('failed')
            ).length;

            // Build count HTML
            let countHTML = filteredEvents.length + ' event' + (filteredEvents.length !== 1 ? 's' : '');
            if (errorCount > 0) {
                countHTML += ' • <span class="error-count">' + errorCount + ' error' + (errorCount !== 1 ? 's' : '') + '</span>';
            }
            resultsCount.innerHTML = countHTML;

            // Render timeline
            timeline.innerHTML = filteredEvents.map((event, index) => {
                const isError = event.severity === 'error' ||
                               event.body.toLowerCase().includes('error') ||
                               event.body.toLowerCase().includes('failed');

                const containerName = event.container.replace('devstack-', '').replace('-1', '');

                return ` + "`" + `
                    <div class="timeline-item">
                        <div class="timeline-dot ${isError ? 'error' : ''}"></div>
                        <div class="timeline-content ${isError ? 'error' : ''}">
                            <div class="timeline-time">${event.timestamp}</div>
                            <div class="timeline-container">
                                <span class="container-pill">#${containerName}</span>
                                ${isError ? '<span class="error-badge">ERROR</span>' : ''}
                            </div>
                            <div class="timeline-action">${event.name || 'Log Entry'}</div>
                            <div class="timeline-message" onclick="toggleMessage(this)">${event.body}</div>
                            <div class="timeline-span">span: ${event.span_id}</div>
                        </div>
                    </div>
                ` + "`" + `;
            }).join('');
        }

        function toggleMessage(element) {
            element.classList.toggle('expanded');
        }

        // Initialize on page load
        initContainers();
    </script>
</body>
</html>`
