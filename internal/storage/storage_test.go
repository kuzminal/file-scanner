package storage

import (
	"github.com/sirupsen/logrus"
	"reflect"
	"sync"
	"testing"
	"time"
)

var (
	filesToSync     = make(chan string, 2)
	dstDir          = "/home/alex/Dev/test/s2"
	path            = "/home/alex/Dev/test/s1/1.txt"
	pathNotExist    = "/home/alex/Dev/test/s1/2.txt"
	pathNotExistSrc = "/home/alex/Dev/test/s1/3.txt"
	hash            = "d41d8cd98f00b204e9800998ecf8427e"
	file            = FilesInfo{
		Hash:         hash,
		FilePath:     path,
		FileName:     "1.txt",
		LastModified: time.Now(),
	}
	fileNotExist = FilesInfo{
		Hash:         hash,
		FilePath:     pathNotExist,
		FileName:     "2.txt",
		LastModified: time.Now(),
	}
	fileNotExistSrc = FilesInfo{
		Hash:         hash,
		FilePath:     pathNotExistSrc,
		FileName:     "3.txt",
		LastModified: time.Now(),
	}
	logger = logrus.NewEntry(logrus.New())
)

func TestFiles_AddFileToSync(t *testing.T) {
	type fields struct {
		m map[string]FilesInfo
	}
	type args struct {
		file        FilesInfo
		wg          *sync.WaitGroup
		filesToSync chan string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "simple test",
			fields: fields{m: map[string]FilesInfo{}},
			args:   args{file: file, wg: &sync.WaitGroup{}, filesToSync: filesToSync},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFileStorage(tt.fields.m, logger)
			tt.args.wg.Add(1)
			f.AddFileToSync(tt.args.file, tt.args.wg, tt.args.filesToSync)
		})
	}
}

func TestFiles_ChangeStatusToSync(t *testing.T) {
	type fields struct {
		M map[string]FilesInfo
	}
	type args struct {
		fileName string
		wg       *sync.WaitGroup
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "simple test", args: args{fileName: "1.txt", wg: &sync.WaitGroup{}}, fields: fields{M: map[string]FilesInfo{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFileStorage(tt.fields.M, logger)
			tt.args.wg.Add(1)
			f.ChangeStatusToSync(tt.args.fileName, tt.args.wg)
		})
	}
}

func TestFiles_IsFileChanged(t *testing.T) {
	type fields struct {
		m map[string]FilesInfo
	}
	type args struct {
		path         string
		lastModified time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  string
	}{
		{name: "simple test", want: false, want1: hash, fields: fields{m: map[string]FilesInfo{file.FileName: file}}, args: args{path: path, lastModified: time.Now()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFileStorage(tt.fields.m, logger)

			got, got1 := f.IsFileChanged(tt.args.path, tt.args.lastModified)
			if got != tt.want {
				t.Errorf("IsFileChanged() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("IsFileChanged() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestFiles_GetFile(t *testing.T) {
	type fields struct {
		RWMutex sync.RWMutex
		m       map[string]FilesInfo
	}
	type args struct {
		fileName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   FilesInfo
		want1  bool
	}{
		{name: "simple test", want: file, want1: true, fields: fields{m: map[string]FilesInfo{file.FileName: file}}, args: args{fileName: file.FileName}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Files{
				RWMutex: tt.fields.RWMutex,
				m:       tt.fields.m,
			}
			got, got1 := f.GetFile(tt.args.fileName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFile() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetFile() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestFiles_CheckIfExistAndRemove(t *testing.T) {
	type fields struct {
		m      map[string]FilesInfo
		logger *logrus.Entry
	}
	type args struct {
		dstDir string
		wg     *sync.WaitGroup
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "simple test", wantErr: true, fields: fields{map[string]FilesInfo{fileNotExistSrc.FileName: fileNotExistSrc}, logger}, args: args{dstDir: dstDir, wg: &sync.WaitGroup{}}},
		{name: "not exist", wantErr: false, fields: fields{map[string]FilesInfo{fileNotExist.FileName: fileNotExist}, logger}, args: args{dstDir: dstDir, wg: &sync.WaitGroup{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Files{
				m:      tt.fields.m,
				logger: tt.fields.logger,
			}
			tt.args.wg.Add(1)
			if err := f.CheckIfExistAndRemove(tt.args.dstDir, tt.args.wg); (err != nil) != tt.wantErr {
				t.Errorf("CheckIfExistAndRemove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
