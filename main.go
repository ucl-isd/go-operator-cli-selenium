package main

import (
	"github.com/ucl-isd/go-operator-cli-selenium/operator"
)

func main() {
	operator.FormsWithCookies("orgName", "fieldName", true, false, "defaultValue", "area", true)
}
