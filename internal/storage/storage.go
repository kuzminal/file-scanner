package storage

import (
	"sync"
	"time"
)

type Status int

const (
	InSync Status = iota
	Sync
)

type FilesInfo struct {
	FileName     string
	FilePath     string
	Hash         string
	LastModified time.Time
	Status       Status
}

type Storage interface {
	ChangeStatusToSync(fileName string, wg *sync.WaitGroup)
	AddFileToSync(file FilesInfo, wg *sync.WaitGroup, fileToSync chan string)
	IsFileChanged(path string, lastModified time.Time) (bool, string)
	GetFile(fileName string) (FilesInfo, bool)
	CheckIfExistAndRemove(dstDir string, wg *sync.WaitGroup) error
}
