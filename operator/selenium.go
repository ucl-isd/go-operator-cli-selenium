package operator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tebeka/selenium"
)

const (
	seleniumPath      = "./selenium-server-standalone.jar"
	chromeDriverPath  = "./chromedriver"
	port              = 8080
	cookieFilePath    = "cookies.json"
	talentlinkURL     = "https://emea5.lumessetalentlink.com/"
	talentlinkDashURL = "https://emea5.lumessetalentlink.com/tlk/app/#/dashboards/generic"
)

func SetupDriver() (selenium.WebDriver, error) {
	opts := []selenium.ServiceOption{}
	selenium.SetDebug(false)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		return nil, fmt.Errorf("error starting the Selenium server: %v", err)
	}

	caps := selenium.Capabilities{"browserName": "chrome"}
	caps["goog:chromeOptions"] = map[string]interface{}{
		"args": []string{"--disable-extensions", "--disable-javascript", "--start-maximized"},
	}
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		service.Stop()
		return nil, fmt.Errorf("error connecting to WebDriver: %v", err)
	}
	return driver, nil
}

func LoginAndSaveCookies(username, password, company string) {
	fmt.Printf("Starting login with %s, %s\n", username, company)
	driver, err := SetupDriver()
	if err != nil {
		log.Fatalf("Error setting up driver: %v", err)
	}
	defer driver.Quit()

	err = driver.Get(talentlinkURL)
	if err != nil {
		log.Fatalf("Failed to load page: %v", err)
	}

	time.Sleep(time.Second)

	// Locate and send keys to login field
	loginElement, err := driver.FindElement(selenium.ByID, "login")
	if err != nil {
		log.Fatalf("Failed to find login input: %v", err)
	}
	loginElement.SendKeys(username)

	time.Sleep(time.Second)
	// Locate and send keys to password field
	passwordElement, err := driver.FindElement(selenium.ByName, "password")
	if err != nil {
		log.Fatalf("Failed to find password input: %v", err)
	}
	passwordElement.SendKeys(password)

	time.Sleep(time.Second)

	// Locate and send keys to company field
	companyElement, err := driver.FindElement(selenium.ByID, "company")
	if err != nil {
		log.Fatalf("Failed to find company input: %v", err)
	}
	companyElement.SendKeys(company)

	time.Sleep(time.Second)

	// Locate and click the login button
	loginButton, err := driver.FindElement(selenium.ByCSSSelector, ".MuiButton-label-297")
	if err != nil {
		log.Fatalf("Failed to find login button: %v", err)
	}
	loginButton.Click()

	time.Sleep(5 * time.Second)

	cookies, err := driver.GetCookies()
	if err != nil {
		log.Fatalf("Failed to get cookies: %v", err)
	}

	cookieFile, err := os.Create(cookieFilePath)
	if err != nil {
		log.Fatalf("Failed to create cookie file: %v", err)
	}
	defer cookieFile.Close()

	json.NewEncoder(cookieFile).Encode(cookies)
	fmt.Println("Login successful, cookies saved to cookies.json")
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
	driver.Refresh()
	fmt.Println("Level one created successfully")
}
