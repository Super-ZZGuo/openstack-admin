package main

import (
	"go-admin/cmd"
)

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	cmd.Execute()
}
