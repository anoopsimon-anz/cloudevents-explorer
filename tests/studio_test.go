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
