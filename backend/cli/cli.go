package cli

import (
	"fmt"
	"log"
	"os"
)

func PrintHelp() {
	fmt.Println(`
FaDO - Usage Information

The function and data orchestrator for multi-serverless platforms.

Command line arguments:

    --config, -c          Path to initial configuration JSON file. Falls back to value from environment.
	--database, -d        Database connection string. Falls back to value from environment.
	--server-url          Server URL. Falls back to value from environment.
	--caddy-admin-url     Caddy admin management endpoint. Falls back to value from environment.
	--lb-endpoint         Load balancer endpoint. Falls back to value from environment.
    --help, -h            Display this information.

Environment variables:

    FADO_CONFIG           Path to initial configuration JSON file. Falls back to "./config.json".
	FADO_DATABASE         Database connection string.
	FADO_SERVER_URL       Server URL.
	FADO_CADDY_ADMIN_URL  Caddy admin managment endpoint.
	FADO_LB_ENDPOINT      Load balander endpoint.`)
}

type CliInput struct {
	ConfigFilePath, DatabaseConnectionString, ServerURL, CaddyAdminURL, LBDomain, LBPort string
}

var Input CliInput

func ReadInput() (i CliInput) {
	argsWithoutProg := os.Args[1:]
	var nextVal string
	for _, a := range argsWithoutProg {
		if a == "--config" || a == "-c" {
			nextVal = "config"
		} else if a == "--database" || a == "-d" {
			nextVal = "database"
		} else if a == "--server-url" {
			nextVal = "server-url"
		} else if a == "--caddy-admin-url" {
			nextVal = "caddy-admin-url"
		} else if a == "--lb-domain" {
			nextVal = "lb-domain"
		} else if a == "--lb-port" {
			nextVal = "lb-port"
		} else if a == "--help" || a == "-h" {
			PrintHelp()
			os.Exit(0)
		} else if nextVal == "config" {
			i.ConfigFilePath = a
			nextVal = ""
		} else if nextVal == "database" {
			i.DatabaseConnectionString = a
			nextVal = ""
		} else if nextVal == "server-url" {
			i.ServerURL = a
			nextVal = ""
		} else if nextVal == "caddy-admin-url" {
			i.CaddyAdminURL = a
			nextVal = ""
		} else if nextVal == "lb-domain" {
			i.LBDomain = a
			nextVal = ""
		} else if nextVal == "lb-port" {
			i.LBPort = a
			nextVal = ""
		} else {
			problemArgument := a
			if nextVal != "" { problemArgument = nextVal }

			log.Printf("ERROR: Encountered unexpected argument '%v'.", problemArgument)
			PrintHelp()
			os.Exit(1)
		}
	}

	if nextVal != "" {
		log.Printf("ERROR: Encountered unexpected argument '%v'.", nextVal)
		PrintHelp()
		os.Exit(1)
	}

	if i.ConfigFilePath == "" { i.ConfigFilePath = os.Getenv("FADO_CONFIG") }
	if i.ConfigFilePath == "" { i.ConfigFilePath = "./config.json" }

	if i.DatabaseConnectionString == "" { i.DatabaseConnectionString = os.Getenv("FADO_DATABASE") }
	if i.DatabaseConnectionString == "" { i.DatabaseConnectionString = "postgres://fado:password@localhost:5454/fado_db" }

	if i.ServerURL == "" { i.ServerURL = os.Getenv("FADO_SERVER_URL") }
	if i.ServerURL == "" { i.ServerURL = "https://server.fado" }

	if i.CaddyAdminURL == "" { i.CaddyAdminURL = os.Getenv("FADO_CADDY_ADMIN_URL") }
	if i.CaddyAdminURL == "" { i.CaddyAdminURL = "http://caddy-admin.fado" }

	if i.LBDomain == "" { i.LBDomain = os.Getenv("FADO_LB_DOMAIN") }
	if i.LBDomain == "" { i.LBDomain = "" }

	if i.LBPort == "" { i.LBPort = os.Getenv("FADO_LB_PORT") }
	if i.LBPort == "" { i.LBPort = "443" }

	Input = i

	return
}