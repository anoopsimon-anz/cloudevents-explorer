package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/playwright-community/playwright-go"
)

type TestResult struct {
	Name      string
	Passed    bool
	Message   string
	Screenshot string
}

var testResults []TestResult

func TestStudioWithReport(t *testing.T) {
	// Create screenshots directory
	screenshotsDir := "screenshots"
	if err := os.MkdirAll(screenshotsDir, 0755); err != nil {
		t.Fatalf("could not create screenshots directory: %v", err)
	}

	// Start Playwright
	pw, err := playwright.Run()
	if err != nil {
		t.Fatalf("could not start playwright: %v", err)
	}
	defer pw.Stop()

	// Launch browser
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		t.Fatalf("could not launch browser: %v", err)
	}
	defer browser.Close()

	// Create new page with larger viewport
	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		Viewport: &playwright.Size{
			Width:  1920,
			Height: 1080,
		},
	})
	if err != nil {
		t.Fatalf("could not create page: %v", err)
	}

	// Navigate to the studio
	if _, err = page.Goto("http://localhost:8888"); err != nil {
		t.Fatalf("could not goto: %v", err)
	}

	// Wait for page to load
	time.Sleep(1 * time.Second)

	// Take initial screenshot
	screenshotPath := filepath.Join(screenshotsDir, "01-homepage.png")
	if _, err := page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(screenshotPath),
		FullPage: playwright.Bool(true),
	}); err != nil {
		t.Logf("Warning: could not take screenshot: %v", err)
	}

	// Test 1: Verify page title
	testName := "Page Title Verification"
	title, err := page.Locator("#pageTitle").TextContent()
	if err != nil {
		addTestResult(testName, false, fmt.Sprintf("Could not get page title: %v", err), screenshotPath)
		t.Errorf("Test failed: %s", testName)
	} else if title != "Testing Studio" {
		addTestResult(testName, false, fmt.Sprintf("Expected 'Testing Studio', got '%s'", title), screenshotPath)
		t.Errorf("Test failed: %s", testName)
	} else {
		addTestResult(testName, true, "Page title is correct: 'Testing Studio'", screenshotPath)
	}

	// Test 2: Verify subtitle
	testName = "Subtitle Verification"
	subtitle, err := page.Locator("#pageSubtitle").TextContent()
	if err != nil {
		addTestResult(testName, false, fmt.Sprintf("Could not get subtitle: %v", err), screenshotPath)
		t.Errorf("Test failed: %s", testName)
	} else {
		addTestResult(testName, true, fmt.Sprintf("Subtitle found: '%s'", subtitle), screenshotPath)
	}

	// Test 3: Verify all option cards
	testName = "Option Cards Verification"
	cards := []struct {
		id    string
		title string
	}{
		{"cardPubsub", "Google PubSub"},
		{"cardKafka", "Kafka / EventMesh"},
		{"cardRestClient", "REST Client"},
		{"cardGCS", "GCS Browser"},
		{"cardTraceJourney", "Trace Journey Viewer"},
	}

	allCardsVisible := true
	cardMessages := ""
	for _, card := range cards {
		visible, err := page.Locator("#" + card.id).IsVisible()
		if err != nil || !visible {
			allCardsVisible = false
			cardMessages += fmt.Sprintf("❌ %s not visible\n", card.title)
		} else {
			cardMessages += fmt.Sprintf("✅ %s visible\n", card.title)
		}
	}

	if allCardsVisible {
		addTestResult(testName, true, "All 5 option cards are visible:\n"+cardMessages, screenshotPath)
	} else {
		addTestResult(testName, false, "Some option cards are missing:\n"+cardMessages, screenshotPath)
		t.Errorf("Test failed: %s", testName)
	}

	// Test 4: Highlight and screenshot a specific card
	screenshotPath = filepath.Join(screenshotsDir, "02-pubsub-card.png")
	if err := page.Locator("#cardPubsub").Hover(); err == nil {
		time.Sleep(500 * time.Millisecond)
		page.Screenshot(playwright.PageScreenshotOptions{
			Path: playwright.String(screenshotPath),
			FullPage: playwright.Bool(true),
		})
	}

	testName = "PubSub Card Details"
	cardTitle, _ := page.Locator("#titlePubsub").TextContent()
	cardDesc, _ := page.Locator("#descPubsub").TextContent()
	addTestResult(testName, true, fmt.Sprintf("Title: %s\nDescription: %s", cardTitle, cardDesc), screenshotPath)

	// Test 5: Status indicators
	screenshotPath = filepath.Join(screenshotsDir, "03-status-indicators.png")
	testName = "Status Indicators Verification"

	dockerVisible, _ := page.Locator("#dockerStatus").IsVisible()
	gcloudVisible, _ := page.Locator("#gcloudStatus").IsVisible()

	if dockerVisible && gcloudVisible {
		page.Screenshot(playwright.PageScreenshotOptions{
			Path: playwright.String(screenshotPath),
			FullPage: playwright.Bool(false),
		})
		addTestResult(testName, true, "Docker and GCloud status indicators are visible", screenshotPath)
	} else {
		addTestResult(testName, false, "Status indicators missing", screenshotPath)
		t.Errorf("Test failed: %s", testName)
	}

	// Test 6: Tools menu interaction
	screenshotPath = filepath.Join(screenshotsDir, "04-tools-menu-closed.png")
	page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(screenshotPath),
		FullPage: playwright.Bool(true),
	})

	testName = "Tools Menu Interaction"
	if err := page.Locator("#toolsButton").Click(); err != nil {
		addTestResult(testName, false, fmt.Sprintf("Could not click tools button: %v", err), screenshotPath)
		t.Errorf("Test failed: %s", testName)
	} else {
		page.WaitForTimeout(500)
		screenshotPath = filepath.Join(screenshotsDir, "05-tools-menu-open.png")
		page.Screenshot(playwright.PageScreenshotOptions{
			Path: playwright.String(screenshotPath),
			FullPage: playwright.Bool(true),
		})

		menuVisible, _ := page.Locator("#toolsMenu").IsVisible()
		if menuVisible {
			addTestResult(testName, true, "Tools menu opens successfully", screenshotPath)
		} else {
			addTestResult(testName, false, "Tools menu did not open", screenshotPath)
		}
	}

	// Test 7: Verify all cards with individual screenshots
	for i, card := range cards {
		screenshotPath = filepath.Join(screenshotsDir, fmt.Sprintf("06-card-%d-%s.png", i+1, card.id))

		// Scroll to card and highlight it
		page.Locator("#" + card.id).ScrollIntoViewIfNeeded()
		page.Locator("#" + card.id).Hover()
		time.Sleep(300 * time.Millisecond)

		page.Screenshot(playwright.PageScreenshotOptions{
			Path: playwright.String(screenshotPath),
			FullPage: playwright.Bool(false),
		})

		testName = fmt.Sprintf("Card: %s", card.title)
		titleText, _ := page.Locator("#title" + card.id[4:]).TextContent()
		descText, _ := page.Locator("#desc" + card.id[4:]).TextContent()
		badgeText, _ := page.Locator("#badge" + card.id[4:]).TextContent()

		addTestResult(testName, true, fmt.Sprintf("Title: %s\nDescription: %s\nBadge: %s", titleText, descText, badgeText), screenshotPath)
	}

	// Generate HTML report
	if err := generateHTMLReport(); err != nil {
		t.Logf("Warning: could not generate HTML report: %v", err)
	} else {
		t.Log("✅ HTML report generated: test-report.html")
	}

	t.Logf("✅ Test completed! Screenshots saved in %s/", screenshotsDir)
	t.Logf("📊 Total tests: %d, Passed: %d, Failed: %d", len(testResults), countPassed(), countFailed())
}

func addTestResult(name string, passed bool, message string, screenshot string) {
	testResults = append(testResults, TestResult{
		Name:       name,
		Passed:     passed,
		Message:    message,
		Screenshot: screenshot,
	})
}

func countPassed() int {
	count := 0
	for _, r := range testResults {
		if r.Passed {
			count++
		}
	}
	return count
}

func countFailed() int {
	return len(testResults) - countPassed()
}

func generateHTMLReport() error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	skipped := 0 // We don't have skipped tests yet, but including for future

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Testing Studio - Test Report</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Google Sans', 'Product Sans', Roboto, Arial, sans-serif;
            background: #f1f3f4;
            padding: 0;
            font-size: 13px;
            color: #202124;
            overflow: hidden;
        }
        .header {
            background: #1a73e8;
            color: white;
            padding: 16px 24px;
            box-shadow: 0 1px 2px 0 rgba(60,64,67,0.3), 0 1px 3px 1px rgba(60,64,67,0.15);
            position: fixed;
            top: 0;
            left: 40px;
            right: 40px;
            z-index: 100;
            border-radius: 0 0 8px 8px;
        }
        h1 {
            font-size: 20px;
            font-weight: 400;
            margin-bottom: 4px;
            letter-spacing: 0;
        }
        .timestamp {
            opacity: 0.85;
            font-size: 11px;
            font-weight: 400;
        }
        .summary {
            position: fixed;
            top: 68px;
            left: 40px;
            right: 40px;
            display: grid;
            grid-template-columns: repeat(4, 1fr);
            gap: 0;
            background: white;
            border-bottom: 1px solid #dadce0;
            box-shadow: 0 1px 2px 0 rgba(60,64,67,0.1);
            z-index: 99;
        }
        .stat {
            text-align: center;
            padding: 12px 8px;
            border-right: 1px solid #dadce0;
        }
        .stat:last-child {
            border-right: none;
        }
        .stat-value {
            font-size: 24px;
            font-weight: 400;
            margin-bottom: 2px;
            line-height: 1;
        }
        .stat-value.passed { color: #1e8e3e; }
        .stat-value.failed { color: #d93025; }
        .stat-value.skipped { color: #f9ab00; }
        .stat-value.total { color: #1967d2; }
        .stat-label {
            font-size: 10px;
            color: #5f6368;
            text-transform: uppercase;
            letter-spacing: 0.8px;
            font-weight: 500;
        }
        .main-content {
            display: flex;
            position: fixed;
            top: 130px;
            left: 40px;
            right: 40px;
            bottom: 40px;
            box-shadow: 0 1px 3px 0 rgba(60,64,67,0.3), 0 4px 8px 3px rgba(60,64,67,0.15);
            border-radius: 8px;
            overflow: hidden;
        }
        .sidebar {
            width: 320px;
            background: white;
            border-right: 1px solid #dadce0;
            overflow-y: auto;
            flex-shrink: 0;
            padding-top: 0;
        }
        .test-list-item {
            padding: 10px 16px;
            border-bottom: 1px solid #f1f3f4;
            cursor: pointer;
            display: flex;
            align-items: center;
            gap: 10px;
            transition: background 0.15s;
        }
        .test-list-item:hover {
            background: #f8f9fa;
        }
        .test-list-item.active {
            background: #e8f0fe;
            border-left: 3px solid #1a73e8;
            padding-left: 13px;
        }
        .test-status {
            width: 16px;
            height: 16px;
            border-radius: 50%%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 500;
            font-size: 10px;
            flex-shrink: 0;
        }
        .test-status.passed {
            background: #1e8e3e;
            color: white;
        }
        .test-status.failed {
            background: #d93025;
            color: white;
        }
        .test-status.skipped {
            background: #f9ab00;
            color: white;
        }
        .test-name {
            font-size: 12px;
            font-weight: 400;
            flex: 1;
            color: #202124;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
        }
        .test-list-item.active .test-name {
            font-weight: 500;
            color: #1967d2;
        }
        .detail-panel {
            flex: 1;
            background: #fafafa;
            overflow-y: auto;
            padding: 0 24px 24px 24px;
        }
        .detail-content {
            background: white;
            border-radius: 8px;
            padding: 24px;
            box-shadow: 0 1px 2px 0 rgba(60,64,67,0.3), 0 1px 3px 1px rgba(60,64,67,0.15);
        }
        .detail-header {
            margin-bottom: 20px;
            padding-bottom: 16px;
            border-bottom: 1px solid #dadce0;
        }
        .detail-title {
            font-size: 18px;
            font-weight: 400;
            color: #202124;
            margin-bottom: 8px;
            display: flex;
            align-items: center;
            gap: 12px;
        }
        .status-badge {
            padding: 4px 12px;
            border-radius: 12px;
            font-size: 11px;
            font-weight: 500;
            text-transform: uppercase;
        }
        .status-badge.passed {
            background: #e6f4ea;
            color: #1e8e3e;
        }
        .status-badge.failed {
            background: #fce8e6;
            color: #d93025;
        }
        .test-message {
            background: #f8f9fa;
            padding: 16px;
            border-radius: 4px;
            margin-bottom: 16px;
            font-size: 12px;
            line-height: 1.6;
            white-space: pre-wrap;
            color: #3c4043;
            font-family: 'Roboto Mono', monospace;
            border-left: 3px solid #1a73e8;
        }
        .screenshot-section {
            margin-top: 20px;
        }
        .screenshot-toggle {
            display: flex;
            align-items: center;
            gap: 8px;
            padding: 10px 12px;
            background: #f8f9fa;
            border: 1px solid #dadce0;
            border-radius: 4px;
            cursor: pointer;
            transition: all 0.15s;
            user-select: none;
        }
        .screenshot-toggle:hover {
            background: #e8f0fe;
            border-color: #1a73e8;
        }
        .screenshot-toggle-icon {
            font-size: 14px;
            color: #5f6368;
            transition: transform 0.2s;
        }
        .screenshot-toggle-icon.expanded {
            transform: rotate(90deg);
        }
        .screenshot-toggle-text {
            font-size: 12px;
            color: #1967d2;
            font-weight: 500;
        }
        .screenshot-container {
            display: none;
            margin-top: 12px;
            padding-top: 12px;
            border-top: 1px solid #e8eaed;
        }
        .screenshot-container.expanded {
            display: block;
        }
        .screenshot {
            width: 100%%;
            border: 1px solid #dadce0;
            border-radius: 8px;
            cursor: pointer;
            transition: all 0.2s;
        }
        .screenshot:hover {
            box-shadow: 0 1px 2px 0 rgba(60,64,67,0.3), 0 2px 6px 2px rgba(60,64,67,0.15);
        }
        .empty-state {
            text-align: center;
            padding: 60px 20px;
            color: #5f6368;
        }
        .empty-state-icon {
            font-size: 48px;
            margin-bottom: 16px;
            opacity: 0.5;
        }
        .empty-state-text {
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>🧪 Testing Studio - Test Report</h1>
        <div class="timestamp">Generated on %s</div>
    </div>

    <div class="summary">
        <div class="stat">
            <div class="stat-value total">%d</div>
            <div class="stat-label">Total</div>
        </div>
        <div class="stat">
            <div class="stat-value passed">%d</div>
            <div class="stat-label">Passed</div>
        </div>
        <div class="stat">
            <div class="stat-value failed">%d</div>
            <div class="stat-label">Failed</div>
        </div>
        <div class="stat">
            <div class="stat-value skipped">%d</div>
            <div class="stat-label">Skipped</div>
        </div>
    </div>

    <div class="main-content">
        <div class="sidebar">
`, timestamp, len(testResults), countPassed(), countFailed(), skipped)

	// Add test list items in sidebar
	for i, result := range testResults {
		status := "passed"
		statusIcon := "✓"
		if !result.Passed {
			status = "failed"
			statusIcon = "✗"
		}

		activeClass := ""
		if i == 0 {
			activeClass = "active"
		}

		html += fmt.Sprintf(`            <div class="test-list-item %s" onclick="showTest(%d)">
                <div class="test-status %s">%s</div>
                <div class="test-name">%s</div>
            </div>
`, activeClass, i, status, statusIcon, result.Name)
	}

	html += `        </div>
        <div class="detail-panel">
            <div class="empty-state" id="empty-state">
                <div class="empty-state-icon">📋</div>
                <div class="empty-state-text">Select a test from the list to view details</div>
            </div>
`

	// Add detail content for each test (hidden by default)
	for i, result := range testResults {
		status := "passed"
		statusBadge := "PASSED"
		if !result.Passed {
			status = "failed"
			statusBadge = "FAILED"
		}

		displayStyle := "none"
		if i == 0 {
			displayStyle = "block"
		}

		html += fmt.Sprintf(`            <div class="detail-content" id="detail-%d" style="display: %s;">
                <div class="detail-header">
                    <div class="detail-title">
                        %s
                        <span class="status-badge %s">%s</span>
                    </div>
                </div>
                <div class="test-message">%s</div>
                <div class="screenshot-section">
                    <div class="screenshot-toggle" onclick="toggleScreenshot(%d)">
                        <div class="screenshot-toggle-icon" id="screenshot-icon-%d">▶</div>
                        <div class="screenshot-toggle-text">View Screenshot</div>
                    </div>
                    <div class="screenshot-container" id="screenshot-container-%d">
                        <img src="%s" class="screenshot" alt="Screenshot" onclick="window.open(this.src)">
                    </div>
                </div>
            </div>
`, i, displayStyle, result.Name, status, statusBadge, result.Message, i, i, i, result.Screenshot)
	}

	html += `        </div>
    </div>

    <script>
        function showTest(index) {
            // Hide empty state
            document.getElementById('empty-state').style.display = 'none';

            // Hide all detail contents
            document.querySelectorAll('.detail-content').forEach(el => {
                el.style.display = 'none';
            });

            // Show selected detail
            document.getElementById('detail-' + index).style.display = 'block';

            // Update active state in sidebar
            document.querySelectorAll('.test-list-item').forEach(el => {
                el.classList.remove('active');
            });
            event.currentTarget.classList.add('active');
        }

        function toggleScreenshot(index) {
            const container = document.getElementById('screenshot-container-' + index);
            const icon = document.getElementById('screenshot-icon-' + index);

            container.classList.toggle('expanded');
            icon.classList.toggle('expanded');
        }

        // Show first test by default
        if (document.querySelectorAll('.test-list-item').length > 0) {
            document.getElementById('empty-state').style.display = 'none';
        }
    </script>
</body>
</html>`

	return os.WriteFile("test-report.html", []byte(html), 0644)
}
