package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"context"
	"strings"

	"artifactor/internal/redis"
	"artifactor/internal/sql"
	"artifactor/internal/config"
	"artifactor/internal/logging"
)

const BANNER = `
 █████╗ ██████╗ ████████╗██╗███████╗ █████╗  ██████╗████████╗ ██████╗ ██████╗
██╔══██╗██╔══██╗╚══██╔══╝██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝██╔═══██╗██╔══██╗
███████║██████╔╝   ██║   ██║█████╗  ███████║██║        ██║   ██║   ██║██████╔╝
██╔══██║██╔══██╗   ██║   ██║██╔══╝  ██╔══██║██║        ██║   ██║   ██║██╔══██╗
██║  ██║██║  ██║   ██║   ██║██║     ██║  ██║╚██████╗   ██║   ╚██████╔╝██║  ██║
╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝   ╚═╝╚═╝     ╚═╝  ╚═╝ ╚═════╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝
`

const MAINTAINER = "Idan Koblik"

const PURPLE = "\033[38;2;87;87;232m"
const RESET = "\033[0m"

var BUILD_TIME string

func main() {
	logging.SetupLogger()

	fmt.Print(PURPLE)
	fmt.Print(BANNER)
	fmt.Print(RESET)
	fmt.Println()

	buildTime := BUILD_TIME
	if buildTime == "" {
		buildTime = "unknown"
	}

	fmt.Printf("\t\t%s • %s\n\n", MAINTAINER, buildTime)

	cfg, err := config.ParseConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		logging.Log.Error(err)
		os.Exit(1)
	}

	logging.Log.Debugf("Max file size that can be uploaded: %d MB\n", cfg.FileUploadLimit)
	logging.Log.Info("Connecting to pgsql database.")
	logging.Log.Debugf("Username: %s", cfg.Sql.Username)
	logging.Log.Debugf("Password: %s", generatePasswordMask())
	logging.Log.Debugf("Addr: %s", cfg.Sql.Addr)
	logging.Log.Debugf("Database: %s\n", cfg.Sql.Database)

	err = sql.OpenConnection(&cfg.Sql)
	if err != nil {
		logging.Log.Error("Failed to connect to pgsql database\n", err)
		os.Exit(1)
	}

	defer sql.Conn.Close(context.Background())
	logging.Log.Info("Successfully connected to pgsql database!\n")

	logging.Log.Info("Connecting to redis database.")
	logging.Log.Debugf("Addr: %s", cfg.Redis.Addr)
	logging.Log.Debugf("Password: %s", generatePasswordMask())

	err = redis.OpenConnection(&cfg.Redis)
	if err != nil {
		logging.Log.Error("Failed to connect to redis database\n", err)
		os.Exit(1)
	}

	defer redis.Client.Close()
	logging.Log.Info("Successfully connected to redis database!\n")
}

func generatePasswordMask() string {
	n := rand.N(18) + 5
	return strings.Repeat("*", n)
}
