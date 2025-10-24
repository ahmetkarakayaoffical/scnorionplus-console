//go:build windows

package main

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/ahmetkarakayaoffical/scnorionplus-console/internal/common"
	"github.com/ahmetkarakayaoffical/scnorionplus-console/utils"
	"github.com/go-co-op/gocron/v2"
	"golang.org/x/sys/windows/svc"
)

func main() {
	var err error

	w := common.NewWorker("scnorionplus-console-service.txt")

	// Start Task Scheduler
	w.TaskScheduler, err = gocron.NewScheduler()
	if err != nil {
		log.Fatalf("[FATAL]: could not create task scheduler, reason: %s", err.Error())
		return
	}
	w.TaskScheduler.Start()
	log.Println("[INFO]: task scheduler has been started")

	if err := w.GenerateConsoleConfig(); err != nil {
		log.Printf("[ERROR]: could not generate config for scnorionplus console: %v", err)
		if err := w.StartGenerateConsoleConfigJob(); err != nil {
			log.Fatalf("[FATAL]: could not start job to generate config for scnorionplus console: %v", err)
		}
	}

	// Get working directory
	cwd, err := utils.GetWd()
	if err != nil {
		log.Fatal("[FATAL]: could not get working directory")
	}

	// Create temp directory for downloads
	w.DownloadDir = filepath.Join(cwd, "tmp", "download")
	if strings.HasSuffix(cwd, "tmp") {
		w.DownloadDir = filepath.Join(cwd, "download")
	}

	if err := w.CreateDowloadTempDir(); err != nil {
		log.Fatalf("[FATAL]: could not create download temp dir: %v", err)
	}

	// Create winget directory for index.db
	w.WinGetDBFolder = filepath.Join(cwd, "tmp", "winget")
	if err := w.CreateWingetDBDir(); err != nil {
		log.Fatalf("[FATAL]: could not create winget temp dir: %v", err)
	}

	// Create flatpak directory for flatpak.db
	w.FlatpakDBFolder = filepath.Join(cwd, "tmp", "flatpak")
	if err := w.CreateFlatpakDBDir(); err != nil {
		log.Fatalf("[FATAL]: could not create flatpak temp dir: %v", err)
	}

	// Create brew directory for brew.db
	w.BrewDBFolder = filepath.Join(cwd, "tmp", "brew")
	if err := w.CreateBrewDBDir(); err != nil {
		log.Fatalf("[FATAL]: could not create brew temp dir: %v", err)
	}

	// Create winget directory for flatpak.db
	w.CommonSoftwareDBFolder = filepath.Join(cwd, "tmp", "commondb")
	if err := w.CreateCommonSoftwareDBDir(); err != nil {
		log.Fatalf("[FATAL]: could not create commondb temp dir: %v", err)
	}

	// Create server releases directory
	w.ServerReleasesFolder = filepath.Join(cwd, "tmp", "server-releases")
	if err := w.CreateServerReleasesDir(); err != nil {
		log.Fatalf("[FATAL]: could not create server releases temp dir: %v", err)
	}

	// Configure the windows service
	s := utils.NewscnorionplusWindowsService()
	s.ServiceStart = w.StartWorker
	s.ServiceStop = w.StopWorker

	// Run service
	if err := svc.Run("scnorionplus-console-service", s); err != nil {
		log.Printf("[ERROR]: could not run service: %v", err)
	}
}
