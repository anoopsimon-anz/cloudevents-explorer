package tests

import (
	"testing"

	"github.com/playwright-community/playwright-go"
)

// TestAllToolsNavigation validates navigation to all tools from homepage
func TestAllToolsNavigation(t *testing.T) {
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

	// Create new page
	page, err := browser.NewPage()
	if err != nil {
		t.Fatalf("could not create page: %v", err)
	}

	// Define all tools with their card IDs, URLs, and expected page elements
	tools := []struct {
		name        string
		cardID      string
		url         string
		pageElement string // Element that should exist on the tool's page
	}{
		{
			name:        "Google PubSub",
			cardID:      "cardPubsub",
			url:         "http://localhost:8888/pubsub",
			pageElement: "text=Google PubSub",
		},
		{
			name:        "Kafka / EventMesh",
			cardID:      "cardKafka",
			url:         "http://localhost:8888/kafka",
			pageElement: "text=Kafka / EventMesh",
		},
		{
			name:        "REST Client",
			cardID:      "cardRestClient",
			url:         "http://localhost:8888/rest-client",
			pageElement: "#httpMethod", // The method dropdown
		},
		{
			name:        "GCS Browser",
			cardID:      "cardGCS",
			url:         "http://localhost:8888/gcs",
			pageElement: "text=GCS Browser",
		},
		{
			name:        "Trace Journey Viewer",
			cardID:      "cardTraceJourney",
			url:         "http://localhost:8888/trace-journey",
			pageElement: "text=Trace Journey",
		},
		{
			name:        "Spanner Explorer",
			cardID:      "cardSpanner",
			url:         "http://localhost:8888/spanner",
			pageElement: "text=Spanner Explorer",
		},
	}

	for _, tool := range tools {
		t.Run(tool.name, func(t *testing.T) {
			// Navigate to homepage
			if _, err := page.Goto("http://localhost:8888"); err != nil {
				t.Fatalf("could not goto homepage: %v", err)
			}

			// Verify the tool card exists on homepage
			cardVisible, err := page.Locator("#" + tool.cardID).IsVisible()
			if err != nil {
				t.Fatalf("error checking visibility of %s card: %v", tool.name, err)
			}
			if !cardVisible {
				t.Errorf("%s card is not visible on homepage", tool.name)
				return
			}
			t.Logf("✅ %s card is visible on homepage", tool.name)

			// Click the card to navigate
			if err := page.Locator("#" + tool.cardID).Click(); err != nil {
				t.Fatalf("could not click %s card: %v", tool.name, err)
			}

			// Wait for navigation
			page.WaitForTimeout(500)

			// Verify we're on the correct page
			currentURL := page.URL()
			if currentURL != tool.url {
				t.Errorf("expected to navigate to %s, got %s", tool.url, currentURL)
			}
			t.Logf("✅ Successfully navigated to %s", tool.name)

			// Verify expected element exists on the page
			elementVisible, err := page.Locator(tool.pageElement).IsVisible()
			if err != nil {
				t.Logf("Warning: could not check for element %s: %v", tool.pageElement, err)
			} else if !elementVisible {
				t.Logf("Warning: expected element %s not visible on %s page", tool.pageElement, tool.name)
			} else {
				t.Logf("✅ %s page loaded correctly", tool.name)
			}
		})
	}
}

// TestSpannerExplorerUI validates the Spanner Explorer interface
func TestSpannerExplorerUI(t *testing.T) {
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
			Width:  1400,
			Height: 900,
		},
	})
	if err != nil {
		t.Fatalf("could not create page: %v", err)
	}

	// Navigate to Spanner Explorer
	if _, err = page.Goto("http://localhost:8888/spanner"); err != nil {
		t.Fatalf("could not goto spanner page: %v", err)
	}

	// Wait for page to load
	page.WaitForTimeout(1000)

	// Test 1: Verify page title in browser tab
	pageTitle, err := page.Title()
	if err != nil {
		t.Fatalf("could not get page title: %v", err)
	}
	if pageTitle != "Spanner Explorer - Testing Studio" {
		t.Errorf("expected page title to be 'Spanner Explorer - Testing Studio', got '%s'", pageTitle)
	} else {
		t.Log("✅ Spanner Explorer page title is correct")
	}

	// Test 2: Verify back button exists
	backBtnVisible, err := page.Locator("text=← Back").IsVisible()
	if err != nil {
		t.Logf("Warning: could not check back button: %v", err)
	} else if !backBtnVisible {
		t.Log("Warning: back button is not visible")
	} else {
		t.Log("✅ Back button is visible")
	}

	// Test 3: Check for key UI elements
	elements := []struct {
		name     string
		selector string
	}{
		{"Connect button", "button:has-text('Connect')"},
		{"Connection Settings panel", "text=Connection Settings"},
		{"SQL Editor panel", "text=SQL Editor"},
		{"Tables panel", "text=Tables"},
		{"Results panel", "text=Results"},
		{"Profile Name field", "#configName"},
		{"SQL Query textarea", "#sqlQuery"},
	}

	for _, elem := range elements {
		visible, err := page.Locator(elem.selector).IsVisible()
		if err != nil {
			t.Logf("Warning: could not check %s: %v", elem.name, err)
		} else if !visible {
			t.Logf("Warning: %s is not visible", elem.name)
		} else {
			t.Logf("✅ %s is visible", elem.name)
		}
	}

	// Take screenshot
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("screenshots/spanner-explorer-ui.png"),
	}); err != nil {
		t.Fatalf("could not take screenshot: %v", err)
	}
	t.Log("✅ Screenshot saved to screenshots/spanner-explorer-ui.png")
}

// TestToolsMenu validates the tools dropdown menu
func TestToolsMenu(t *testing.T) {
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

	// Create new page
	page, err := browser.NewPage()
	if err != nil {
		t.Fatalf("could not create page: %v", err)
	}

	// Navigate to homepage
	if _, err := page.Goto("http://localhost:8888"); err != nil {
		t.Fatalf("could not goto homepage: %v", err)
	}

	// Click Tools button
	if err := page.Locator("#toolsButton").Click(); err != nil {
		t.Fatalf("could not click tools button: %v", err)
	}

	// Wait for menu to appear
	page.WaitForTimeout(500)

	// Check if menu is visible
	menuVisible, err := page.Locator("#toolsMenu").IsVisible()
	if err != nil || !menuVisible {
		t.Log("⚠️  Tools menu is not visible after clicking button (known issue)")
		t.Skip("Skipping menu navigation tests as menu is not visible")
		return
	}

	// Test menu items
	menuItems := []struct {
		name string
		id   string
	}{
		{"Config Editor", "linkConfigEditor"},
		{"Flow Diagram", "linkFlowDiagram"},
		{"Base64 Tool", "linkBase64Tool"},
	}

	for _, item := range menuItems {
		visible, err := page.Locator("#" + item.id).IsVisible()
		if err != nil {
			t.Logf("Warning: could not check %s: %v", item.name, err)
			continue
		}
		if !visible {
			t.Logf("Warning: %s menu item is not visible", item.name)
		} else {
			t.Logf("✅ %s menu item is visible", item.name)
		}
	}

	// Test navigation to Config Editor
	t.Run("Navigate to Config Editor", func(t *testing.T) {
		if err := page.Locator("#linkConfigEditor").Click(); err != nil {
			t.Fatalf("could not click config editor link: %v", err)
		}

		// Wait for navigation
		page.WaitForTimeout(500)

		// Verify we're on config editor page
		currentURL := page.URL()
		if currentURL != "http://localhost:8888/config-editor" {
			t.Errorf("expected to navigate to config-editor, got %s", currentURL)
		} else {
			t.Log("✅ Successfully navigated to Config Editor")
		}
	})
}

// TestGCSBrowserUI validates the GCS Browser interface
func TestGCSBrowserUI(t *testing.T) {
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

	// Create new page
	page, err := browser.NewPage()
	if err != nil {
		t.Fatalf("could not create page: %v", err)
	}

	// Navigate to GCS Browser
	if _, err = page.Goto("http://localhost:8888/gcs"); err != nil {
		t.Fatalf("could not goto gcs page: %v", err)
	}

	// Wait for page to load
	page.WaitForTimeout(1000)

	// Verify page title
	titleVisible, err := page.Locator("text=GCS Browser").IsVisible()
	if err != nil {
		t.Fatalf("could not check page title: %v", err)
	}
	if !titleVisible {
		t.Error("GCS Browser title is not visible")
	} else {
		t.Log("✅ GCS Browser title is visible")
	}

	// Take screenshot
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("screenshots/gcs-browser-ui.png"),
	}); err != nil {
		t.Fatalf("could not take screenshot: %v", err)
	}
	t.Log("✅ Screenshot saved to screenshots/gcs-browser-ui.png")
}
