package main

import (
	"context"
	"flag"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync_dir/internal/scanner"
	"sync_dir/internal/scanner/generated_scanner"
	"sync_dir/internal/storage"
	"sync_dir/internal/storage/generated_storage"
	"sync_dir/internal/utils"
	"syscall"
)

var (
	sourceDir, destDir, logLevel, logPath *string
	timeInterval                          *int
	logger                                *logrus.Entry
)

func init() {
	sourceDir = flag.String("sourceDir", ".", "Source directory to sync")
	destDir = flag.String("destDir", ".", "Destination directory to copy files")
	logLevel = flag.String("logLevel", "info", "Log level")
	logPath = flag.String("logPath", "log.txt", "Path to log file")
	timeInterval = flag.Int("scanInterval", 15, "Time interval for scanning in seconds")
}

func main() {
	flag.Parse()
	file, err := os.OpenFile(*logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	level, _ := logrus.ParseLevel(*logLevel)
	logger = utils.DefaultLogger(file, level)

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()
	fileToSync := make(chan string, 5)
	syncDone := make(chan string, 5)
	fileStorage := storage.NewFileStorage(map[string]storage.FilesInfo{}, logger)
	wrappedStorage := generated_storage.NewStorageWithLogrus(fileStorage, logger)
	dirScanner := generated_scanner.NewFileScannerWithLogrus(
		scanner.NewDirScanner(*sourceDir, *destDir, ctx, fileToSync, syncDone, logger, wrappedStorage, &sync.WaitGroup{}, *timeInterval),
		logger)
	dirScanner.Run()
	dirScanner.Wait()
}
