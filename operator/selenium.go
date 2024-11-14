package operrator

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tebeka/selenium"
)

const (
	seleniumPath      = "./selenium-server-standalone.jar"
	chromeDriverPath  = "././chromedriver"
	port              = 8080
	cookieFilePath    = "cookies.json"
	talentlinkURL     = "https://emea5.lumessetalentlink.com/"
	talentlinkDashURL = "https://emea5.lumessetalentlink.com/tlk/app/#/dashboards/generic"
)

func main() {
	loginCmd := flag.NewFlagSet("login-save-cookies", flag.ExitOnError)
	username := loginCmd.String("username", "", "Username for login")
	password := loginCmd.String("password", "", "Password for login")
	company := loginCmd.String("company", "", "Company identifier for login")

	formsCmd := flag.NewFlagSet("forms-with-cookies", flag.ExitOnError)
	orgName := formsCmd.String("organization_name", "", "Organization name")
	fieldName := formsCmd.String("field_name", "", "Field name")
	displayed := formsCmd.Bool("displayed", false, "Displayed flag")
	mss := formsCmd.Bool("mss", false, "MSS flag")
	defaultValue := formsCmd.String("default_value", "", "Default value")
	area := formsCmd.String("area", "", "Area")
	required := formsCmd.Bool("required", false, "Required flag")

	if len(os.Args) < 2 {
		fmt.Println("login-save-cookies or forms-with-cookies subcommand is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "login-save-cookies":
		loginCmd.Parse(os.Args[2:])
		loginAndSaveCookies(*username, *password, *company)
	case "forms-with-cookies":
		formsCmd.Parse(os.Args[2:])
		formsWithCookies(*orgName, *fieldName, *displayed, *mss, *defaultValue, *area, *required)
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}

func setupDriver() (selenium.WebDriver, error) {
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

func loginAndSaveCookies(username, password, company string) {
	fmt.Printf("Starting login with %s, %s\n", username, company)
	driver, err := setupDriver()
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

func formsWithCookies(orgName, fieldName string, displayed, mss bool, defaultValue, area string, required bool) {
	fmt.Println("Starting forms")

	if _, err := os.Stat(cookieFilePath); os.IsNotExist(err) {
		log.Fatalf("Cookies file not found")
	}

	driver, err := setupDriver()
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
