package main

import (
	"github.com/ucl-isd/go_operator_cli_selenium/operator"
)

func main() {
	operator.FormsWithCookies("orgName", "fieldName", true, false, "defaultValue", "area", true)
}
