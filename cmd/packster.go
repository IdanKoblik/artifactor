package main

import (
	"fmt"
	"strings"
	"os"
	"math/rand/v2"

	internalConfig "packster/internal/config"
	"packster/internal/logging"
	"packster/internal/endpoints"
	"packster/internal/endpoints/gitlab"
	"packster/internal/repository"
	"packster/internal/sql"
	"packster/internal/ui"
	"packster/internal/utils"
	"packster/pkg/types"

	"github.com/gin-gonic/gin"
)

const BANNER = `
██████╗  █████╗  ██████╗██╗  ██╗███████╗████████╗███████╗██████╗
██╔══██╗██╔══██╗██╔════╝██║ ██╔╝██╔════╝╚══██╔══╝██╔════╝██╔══██╗
██████╔╝███████║██║     █████╔╝ ███████╗   ██║   █████╗  ██████╔╝
██╔═══╝ ██╔══██║██║     ██╔═██╗ ╚════██║   ██║   ██╔══╝  ██╔══██╗
██║     ██║  ██║╚██████╗██║  ██╗███████║   ██║   ███████╗██║  ██║
╚═╝     ╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝
`

const (
	MAINTAINER = "Idan Koblik"

	PURPLE = "\033[38;2;87;87;232m"
	RESET = "\033[0m"
)

var BUILD_TIME string

func main() {
	logging.SetupLogger()
	printBanner()

	cfg, err := internalConfig.ParseConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		logging.Log.Error(err)
		os.Exit(1)
	}

	logging.Log.Debugf("Max file size that can be uploaded: %d MB\n", cfg.FileUploadLimit)

	logging.Log.Info("Connecting to pgsql db:")
	logging.Log.Infof("Host: %s", cfg.Sql.Host)
	logging.Log.Infof("Port: %d", cfg.Sql.Port)
	logging.Log.Infof("Username: %s", cfg.Sql.User)
	logging.Log.Infof("Password: %s", generateMask())

	_, err = sql.OpenPgsqlConnection(&cfg.Sql)
	if err != nil {
		logging.Log.Error(err)
		os.Exit(1)
	}

	defer sql.PgsqlConn.Close()
	logging.Log.Info("Successfully connected to pgsql db")

	logging.Log.Info("Setting up rest api")
	router := gin.Default()

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "0.0.0.0:8080"
	}

	logging.Log.Debugf("Addr: %s", addr)

	repo := repository.NewAccountRepo(sql.PgsqlConn)
	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context){
			endpoints.HandleHealth(c, sql.PgsqlConn)
		})
		api.GET("/hosts", func(c *gin.Context){
			endpoints.HandleHosts(c, sql.PgsqlConn)
		})
	}

	// TODO log
	hosts, err := utils.GetHosts(sql.PgsqlConn)
	if err != nil {
		logging.Log.Error(err)
		os.Exit(1)
	}

	utils.Hosts = hosts
	auth := api.Group("/auth")
	setupGitlabEndpoints(auth, repo, hosts)

	ui.RegisterRoutes(router)

	err = router.Run(addr)
	if err != nil {
		logging.Log.Error(err)
		os.Exit(1)
	}

	logging.Log.Info("Packster is up and running!")
}

func setupGitlabEndpoints(authGroup *gin.RouterGroup, repo *repository.AccountRepo, hosts map[string]types.Host) {
	hasGitlab := false
	for _, h := range hosts {
		if h.Type == types.Gitlab {
			hasGitlab = true
			logging.Log.Infof("Gitlab instance detected: %s (id: %d)", h.Url, h.Id)
		}
	}

	if !hasGitlab {
		logging.Log.Info("No gitlab hosts found")
		return
	}

	handler := gitlab.NewGitlabHandler(repo, sql.PgsqlConn)
	group := authGroup.Group("/gitlab")
	{
		group.GET("/redirect", handler.HandleRedirect)
		group.GET("/callback", handler.HandleCallback)
	}
}

func printBanner() {
	fmt.Print(PURPLE)
	fmt.Print(BANNER)
	fmt.Print(RESET)
	fmt.Println()

	buildTime := BUILD_TIME
	if buildTime == "" {
		buildTime = "unknown"
	}

	fmt.Printf("\t\t%s • %s\n\n", MAINTAINER, buildTime)
}

func generateMask() string {
	n := rand.N(18) + 5
	return strings.Repeat("*", n)
}
