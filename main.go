package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ucl-isd/go-operator-cli-selenium/operator"
)

func main() {
	loginCmd := flag.NewFlagSet("login-save-cookies", flag.ExitOnError)
	username := loginCmd.String("username", "", "Username for login")
	password := loginCmd.String("password", "", "Password for login")
	company := loginCmd.String("company", "", "Company identifier for login")

	if len(os.Args) < 2 {
		fmt.Println("login-save-cookies subcommand is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "login-save-cookies":
		loginCmd.Parse(os.Args[2:])
		operator.LoginAndSaveCookies(*username, *password, *company)
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}
