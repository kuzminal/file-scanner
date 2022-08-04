package scanner

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"sync_dir/internal/storage/generated_storage"

	"sync_dir/internal/storage"
	"sync_dir/internal/utils"
	"time"
)

type DirScanner struct {
	wg           *sync.WaitGroup
	ctx          context.Context
	sourceDir    string
	destDir      string
	logger       *logrus.Entry
	filesToSync  chan string
	syncDone     chan string
	storage      generated_storage.StorageWithLogrus
	timeInterval int
}

func (d *DirScanner) Run() {
	ticker := time.NewTicker(time.Duration(d.timeInterval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-d.ctx.Done():
			d.Close()
			return
		case <-ticker.C:
			d.wg.Add(2)
			go d.ScanDir()
			go d.storage.CheckIfExistAndRemove(d.destDir, d.wg)
		case fileName := <-d.filesToSync:
			d.wg.Add(1)
			go d.CopyFile(fileName)
		case fileName := <-d.syncDone:
			d.wg.Add(1)
			go d.storage.ChangeStatusToSync(fileName, d.wg)
		}
	}
}

func (d *DirScanner) Wait() {
	d.wg.Wait()
}

func (d *DirScanner) Close() error {
	d.logger.Println("Closing Scanner")
	for len(d.filesToSync) > 0 && len(d.syncDone) > 0 {
		time.Sleep(time.Millisecond * 10)
	}
	return nil
}

func (d *DirScanner) CopyFile(fileName string) error {
	defer d.wg.Done()
	dst := filepath.Join(d.destDir, fileName)
	src := filepath.Join(d.sourceDir, fileName)
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	d.logger.Printf("Copy file %v with size %v bytes\n", sourceFileStat.Name(), sourceFileStat.Size())
	err = utils.CopyFilesWithOsRW(src, dst)
	d.syncDone <- fileName
	return err
}

func (d *DirScanner) ScanDir() error {
	defer d.wg.Done()
	err := filepath.WalkDir(d.sourceDir, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		info, _ := dir.Info()
		if !info.IsDir() {
			file, ok := d.storage.GetFile(dir.Name())
			if !ok {
				file.Hash = utils.FileMD5(path)
				file.FilePath = path
				file.FileName = dir.Name()
				file.LastModified = info.ModTime()
				d.wg.Add(1)
				go d.storage.AddFileToSync(file, d.wg, d.filesToSync)
			} else if ok && file.Status != storage.InSync {
				res, hash := d.storage.IsFileChanged(path, info.ModTime())
				if res {
					file.Hash = hash
					d.wg.Add(1)
					go d.storage.AddFileToSync(file, d.wg, d.filesToSync)
				}
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking the path : %v\n", err)
	}
	return nil
}

func NewDirScanner(srcDir, dstDir string, ctx context.Context, fileToSync chan string, syncDone chan string, logger *logrus.Entry, storage generated_storage.StorageWithLogrus, wg *sync.WaitGroup, timeInterval int) *DirScanner {
	return &DirScanner{
		wg:           wg,
		ctx:          ctx,
		sourceDir:    srcDir,
		destDir:      dstDir,
		logger:       logger,
		filesToSync:  fileToSync,
		storage:      storage,
		syncDone:     syncDone,
		timeInterval: timeInterval,
	}
}
