package templates

const FlowDiagram = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TMS Event Flow - Testing Studio</title>
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
            background: #1a73e8;
            border: 2px solid #000;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 9px;
            font-weight: bold;
            color: white;
            font-family: 'Courier New', monospace;
            padding: 8px;
            border-radius: 4px;
        }
        .message-incoming {
            animation: moveIncoming 5s infinite;
            left: -120px;
            top: 30%;
            width: 80px;
            height: 50px;
        }
        @keyframes moveIncoming {
            0% { left: -120px; opacity: 1; }
            35% { left: 50%; transform: translateX(-50%) scale(1); }
            40% { left: 50%; transform: translateX(-50%) scale(0.7); opacity: 0.7; }
            45% { left: 50%; transform: translateX(-50%) scale(0); opacity: 0; }
            100% { left: 50%; transform: translateX(-50%) scale(0); opacity: 0; }
        }
        .message-outgoing {
            animation: moveOutgoing 5s infinite;
            right: -120px;
            bottom: 30%;
            background: #188038;
            width: 70px;
            height: 50px;
        }
        @keyframes moveOutgoing {
            0%, 45% { right: -120px; opacity: 0; transform: scale(0); }
            50% { right: 50%; transform: translateX(50%) scale(0.7); opacity: 0.7; }
            55% { right: 50%; transform: translateX(50%) scale(1); opacity: 1; }
            100% { right: -120px; opacity: 1; }
        }
        .status-text {
            position: absolute;
            bottom: 20px;
            left: 50%;
            transform: translateX(-50%);
            font-size: 14px;
            font-weight: bold;
            text-align: center;
            animation: statusBlink 5s infinite;
        }
        @keyframes statusBlink {
            0%, 35% { opacity: 0; }
            40%, 90% { opacity: 1; }
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
                    <div>📬</div>
                    <div>UNICA</div>
                    <div style="font-size: 10px;">Marketing Events</div>
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

                <div class="message-flow message-incoming" style="line-height: 1.2;">{...}<br>Avro<br>Event</div>
                <div class="message-flow message-outgoing" style="line-height: 1.3;">⚙️<br>WF<br>Trigger</div>

                <div class="label-temporal">
                    <div>⚙️</div>
                    <div>TEMPORAL</div>
                    <div style="font-size: 10px;">Workflow</div>
                </div>

                <div class="status-text">⚡ Processing Avro Event → Triggering Workflow</div>
            </div>

            <div style="margin-top: 20px; padding: 15px; background: #f5f5f5; border-left: 3px solid #000;">
                <strong>What's happening:</strong>
                <ul style="margin: 10px 0 0 20px; font-size: 13px; line-height: 1.8;">
                    <li><strong>Blue box ({...} Avro Event)</strong>: Marketing event arriving from Unica topic in Avro format</li>
                    <li><strong>Worker</strong>: Event Handler consuming the message (arms waving = working!)</li>
                    <li><strong>Green box (⚙️ WF Trigger)</strong>: Temporal workflow being triggered with workflow icon</li>
                    <li><strong>Continuous cycle</strong>: This happens 24x7, never stops! (5-second cycle for visibility)</li>
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
`
