package storage

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
	"sync_dir/internal/utils"
	"time"
)

type Files struct {
	sync.RWMutex
	m      map[string]FilesInfo
	logger *logrus.Entry
}

func (f *Files) GetFile(fileName string) (FilesInfo, bool) {
	f.RLock()
	defer f.RUnlock()
	file, ok := f.m[fileName]
	return file, ok
}

func (f *Files) ChangeStatusToSync(fileName string, wg *sync.WaitGroup) {
	defer wg.Done()
	f.Lock()
	defer f.Unlock()
	file := f.m[fileName]
	file.Status = Sync
	f.m[fileName] = file
}

func (f *Files) IsFileChanged(path string, lastModified time.Time) (bool, string) {
	_, fileName := filepath.Split(path)
	f.Lock()
	defer f.Unlock()
	file := f.m[fileName]
	if file.LastModified.Before(lastModified) {
		hash := utils.FileMD5(path)
		if file.Hash == hash {
			return false, file.Hash
		} else {
			return true, hash
		}
	}
	return false, file.Hash
}

func (f *Files) AddFileToSync(file FilesInfo, wg *sync.WaitGroup, filesToSync chan string) {
	defer wg.Done()
	f.Lock()
	defer f.Unlock()
	file.Status = InSync
	f.m[file.FileName] = file
	filesToSync <- file.FileName
}

func (f *Files) CheckIfExistAndRemove(dstDir string, wg *sync.WaitGroup) error {
	f.Lock()
	defer f.Unlock()
	defer wg.Done()
	for _, file := range f.m {
		if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
			err = os.Remove(filepath.Join(dstDir, file.FileName))
			if err == nil {
				delete(f.m, file.FileName)
				f.logger.Infof("Delete file %v from destination directory", file.FileName)
			} else {
				return err
			}
		}
	}
	return nil
}

func NewFileStorage(m map[string]FilesInfo, logger *logrus.Entry) *Files {
	return &Files{
		m:      m,
		logger: logger,
	}
}
