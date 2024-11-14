package operator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/tebeka/selenium"
)

const (
	cookieFilePath    = "cookies.json"
	talentlinkDashURL = "http://example.com"
)

func SetupDriver() (selenium.WebDriver, error) {
	// Example setup for the WebDriver
	const (
		seleniumPath    = "path/to/selenium-server-standalone.jar"
		geckoDriverPath = "path/to/geckodriver"
		port            = 8080
	)
	opts := []selenium.ServiceOption{}
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		return nil, fmt.Errorf("error starting the Selenium server: %v", err)
	}

	caps := selenium.Capabilities{"browserName": "firefox"}
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		service.Stop()
		return nil, fmt.Errorf("error creating new WebDriver session: %v", err)
	}

	return driver, nil
}

func FormsWithCookies(orgName, fieldName string, displayed, mss bool, defaultValue, area string, required bool) {
	fmt.Println("Starting forms")

	if _, err := os.Stat(cookieFilePath); os.IsNotExist(err) {
		log.Fatalf("Cookies file not found")
	}

	driver, err := SetupDriver()
	if err != nil {
		log.Fatalf("Error setting up driver: %v", err)
	}
	defer driver.Quit()

	err = driver.Get(talentlinkDashURL)
	if err != nil {
		log.Fatalf("Failed to load dashboard: %v", err)
	}

	cookieFile, err := os.Open(cookieFilePath)
	if err != nil {
		log.Fatalf("Failed to open cookie file: %v", err)
	}
	defer cookieFile.Close()

	var cookies []selenium.Cookie
	if err := json.NewDecoder(cookieFile).Decode(&cookies); err != nil {
		log.Fatalf("Failed to decode cookies: %v", err)
	}

	for _, cookie := range cookies {
		driver.AddCookie(&cookie)
	}
	err = driver.Refresh()
	if err != nil {
		log.Fatalf("Failed to refresh: %v", err)
	}
	fmt.Println("Level one created successfully")
}
