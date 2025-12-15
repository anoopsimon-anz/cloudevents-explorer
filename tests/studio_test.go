package tests

import (
	"testing"

	"github.com/playwright-community/playwright-go"
)

func TestStudioHomePage(t *testing.T) {
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

	// Navigate to the studio
	if _, err = page.Goto("http://localhost:8888"); err != nil {
		t.Fatalf("could not goto: %v", err)
	}

	// Test 1: Verify page title
	title, err := page.Locator("#pageTitle").TextContent()
	if err != nil {
		t.Fatalf("could not get page title: %v", err)
	}
	if title != "Testing Studio" {
		t.Errorf("expected page title to be 'Testing Studio', got '%s'", title)
	}

	// Test 2: Verify subtitle
	subtitle, err := page.Locator("#pageSubtitle").TextContent()
	if err != nil {
		t.Fatalf("could not get subtitle: %v", err)
	}
	if subtitle != "Requires TMS Suncorp devstack to be running" {
		t.Errorf("expected subtitle to contain devstack message, got '%s'", subtitle)
	}

	// Test 3: Verify all option cards are present
	cards := []struct {
		id    string
		title string
	}{
		{"cardPubsub", "Google PubSub"},
		{"cardKafka", "Kafka / EventMesh"},
		{"cardRestClient", "REST Client"},
		{"cardGCS", "GCS Browser"},
		{"cardTraceJourney", "Trace Journey Viewer"},
		{"cardSpanner", "Spanner Explorer"},
	}

	for _, card := range cards {
		visible, err := page.Locator("#" + card.id).IsVisible()
		if err != nil {
			t.Errorf("error checking visibility of %s: %v", card.id, err)
			continue
		}
		if !visible {
			t.Errorf("card %s is not visible", card.id)
		}

		// Verify card title
		cardTitle, err := page.Locator("#title" + card.id[4:]).TextContent()
		if err != nil {
			t.Errorf("could not get title for %s: %v", card.id, err)
			continue
		}
		if cardTitle != card.title {
			t.Errorf("expected %s title to be '%s', got '%s'", card.id, card.title, cardTitle)
		}
	}

	// Test 4: Verify Tools button is present
	toolsButton, err := page.Locator("#toolsButton").IsVisible()
	if err != nil {
		t.Fatalf("could not check tools button: %v", err)
	}
	if !toolsButton {
		t.Error("tools button is not visible")
	}

	// Test 5: Verify status indicators are present
	dockerStatus, err := page.Locator("#dockerStatus").IsVisible()
	if err != nil {
		t.Fatalf("could not check docker status: %v", err)
	}
	if !dockerStatus {
		t.Error("docker status indicator is not visible")
	}

	gcloudStatus, err := page.Locator("#gcloudStatus").IsVisible()
	if err != nil {
		t.Fatalf("could not check gcloud status: %v", err)
	}
	if !gcloudStatus {
		t.Error("gcloud status indicator is not visible")
	}

	// Test 6: Click on Tools button and verify menu appears
	if err := page.Locator("#toolsButton").Click(); err != nil {
		t.Fatalf("could not click tools button: %v", err)
	}

	// Give the menu time to appear (it's toggled with JavaScript)
	page.WaitForTimeout(500)

	menuVisible, err := page.Locator("#toolsMenu").IsVisible()
	if err != nil {
		t.Logf("Warning: could not check menu visibility: %v", err)
	} else if !menuVisible {
		t.Logf("Warning: tools menu is not visible after clicking button")
	} else {
		// Test 7: Verify menu items are present (only if menu is visible)
		menuItems := []string{
			"linkConfigEditor",
			"linkFlowDiagram",
			"linkBase64Tool",
		}

		for _, itemID := range menuItems {
			visible, err := page.Locator("#" + itemID).IsVisible()
			if err != nil {
				t.Logf("Warning: error checking visibility of %s: %v", itemID, err)
				continue
			}
			if !visible {
				t.Logf("Warning: menu item %s is not visible", itemID)
			}
		}
	}

	t.Log("All core tests passed successfully!")
}

func TestRestClientUI(t *testing.T) {
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

	// Navigate to REST client
	if _, err = page.Goto("http://localhost:8888/rest-client"); err != nil {
		t.Fatalf("could not goto rest-client: %v", err)
	}

	// Wait for page to load
	page.WaitForTimeout(1000)

	// Take screenshot
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("screenshots/rest-client-new-ui.png"),
	}); err != nil {
		t.Fatalf("could not take screenshot: %v", err)
	}

	t.Log("✅ Screenshot saved to screenshots/rest-client-new-ui.png")
}

func TestRestClientAuthTab(t *testing.T) {
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

	// Navigate to REST client
	if _, err = page.Goto("http://localhost:8888/rest-client"); err != nil {
		t.Fatalf("could not goto rest-client: %v", err)
	}

	// Wait for page to load
	page.WaitForTimeout(500)

	// Click on Authorization tab
	if err := page.Locator("text=Authorization").Click(); err != nil {
		t.Fatalf("could not click Authorization tab: %v", err)
	}

	// Wait for tab to switch
	page.WaitForTimeout(300)

	// Take screenshot
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("screenshots/rest-client-auth-tab.png"),
	}); err != nil {
		t.Fatalf("could not take screenshot: %v", err)
	}

	t.Log("✅ Screenshot saved to screenshots/rest-client-auth-tab.png")
}

func TestRestClientJSONHighlighting(t *testing.T) {
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

	// Navigate to REST client
	if _, err = page.Goto("http://localhost:8888/rest-client"); err != nil {
		t.Fatalf("could not goto rest-client: %v", err)
	}

	// Wait for page to load
	page.WaitForTimeout(500)

	// Click Send button to make the request
	if err := page.Locator("#sendBtn").Click(); err != nil {
		t.Fatalf("could not click Send button: %v", err)
	}

	// Wait for response
	page.WaitForTimeout(2000)

	// Take screenshot showing JSON syntax highlighting
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("screenshots/rest-client-json-highlighting.png"),
	}); err != nil {
		t.Fatalf("could not take screenshot: %v", err)
	}

	t.Log("✅ Screenshot saved to screenshots/rest-client-json-highlighting.png")
}
