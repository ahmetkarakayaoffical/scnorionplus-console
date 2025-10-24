package commands

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"https://github.com/ahmetkarakayaoffical/scnorionplus-console/internal/common"

	"github.com/ahmetkarakayaoffical/utils"
	"github.com/go-co-op/gocron/v2"
	"github.com/urfave/cli/v2"
)

func StartConsole() *cli.Command {
	return &cli.Command{
		Name:   "start",
		Usage:  "Start the scnorionplus console",
		Action: startConsole,
		Flags:  StartConsoleFlags(),
	}
}

func startConsole(cCtx *cli.Context) error {
	worker := common.NewWorker("")

	if err := worker.GenerateConsoleConfigFromCLI(cCtx); err != nil {
		log.Fatalf("[FATAL]: could not generate config for scnorionplus Console: %v", err)
	}

	// Get working directory
	cwd, err := utils.GetWd()
	if err != nil {
		log.Fatal("[FATAL]: could not get working directory")
	}

	// Create temp directory for downloads
	worker.DownloadDir = filepath.Join(cwd, "tmp", "download")
	if strings.HasSuffix(cwd, "tmp") {
		worker.DownloadDir = filepath.Join(cwd, "download")
	}
	if err := worker.CreateDowloadTempDir(); err != nil {
		log.Fatalf("[ERROR]: could not create download temp dir: %v", err)
	}

	// Create server releases directory
	worker.ServerReleasesFolder = filepath.Join(cwd, "tmp", "server-releases")
	if strings.HasSuffix(cwd, "tmp") {
		worker.ServerReleasesFolder = filepath.Join(cwd, "server-releases")
	}
	if err := worker.CreateServerReleasesDir(); err != nil {
		log.Fatalf("[FATAL]: could not create server releases temp dir: %v", err)
	}

	// Create winget directory
	worker.WinGetDBFolder = filepath.Join(cwd, "tmp", "winget")
	if strings.HasSuffix(cwd, "tmp") {
		worker.WinGetDBFolder = filepath.Join(cwd, "winget")
	}
	if err := worker.CreateWingetDBDir(); err != nil {
		log.Fatalf("[FATAL]: could not create winget temp dir: %v", err)
	}

	// Create flatpak directory
	worker.FlatpakDBFolder = filepath.Join(cwd, "tmp", "flatpak")
	if strings.HasSuffix(cwd, "tmp") {
		worker.FlatpakDBFolder = filepath.Join(cwd, "flatpak")
	}
	if err := worker.CreateFlatpakDBDir(); err != nil {
		log.Fatalf("[FATAL]: could not create flatpak temp dir: %v", err)
	}

	// Create brew directory
	worker.BrewDBFolder = filepath.Join(cwd, "tmp", "brew")
	if strings.HasSuffix(cwd, "tmp") {
		worker.BrewDBFolder = filepath.Join(cwd, "brew")
	}
	if err := worker.CreateBrewDBDir(); err != nil {
		log.Fatalf("[FATAL]: could not create brew temp dir: %v", err)
	}

	// Create common software directory
	worker.CommonSoftwareDBFolder = filepath.Join(cwd, "tmp", "commondb")
	if strings.HasSuffix(cwd, "tmp") {
		worker.CommonSoftwareDBFolder = filepath.Join(cwd, "commondb")
	}
	if err := worker.CreateCommonSoftwareDBDir(); err != nil {
		log.Fatalf("[FATAL]: could not create commondb temp dir: %v", err)
	}

	// Save pid to PIDFILE
	if err := os.WriteFile("PIDFILE", []byte(strconv.Itoa(os.Getpid())), 0666); err != nil {
		return err
	}

	// Start Task Scheduler
	worker.TaskScheduler, err = gocron.NewScheduler()
	if err != nil {
		log.Fatalf("[FATAL]: could not create task scheduler, reason: %s", err.Error())
	}
	worker.TaskScheduler.Start()
	log.Println("[INFO]: task scheduler has been started")

	// Start worker
	worker.StartWorker()

	// Keep the connection alive
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-done

	worker.StopWorker()

	return nil
}
