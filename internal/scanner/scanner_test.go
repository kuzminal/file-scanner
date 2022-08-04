package scanner

import (
	"context"
	"github.com/sirupsen/logrus"
	"reflect"
	"sync"
	"sync_dir/internal/storage"
	"sync_dir/internal/storage/generated_storage"
	"testing"
	"time"
)

var (
	ctx, cancel  = context.WithTimeout(context.Background(), 20*time.Second)
	filesToSync  = make(chan string, 2)
	syncDone     = make(chan string, 2)
	logger       = logrus.NewEntry(logrus.New())
	storageFiles = generated_storage.NewStorageWithLogrus(storage.NewFileStorage(map[string]storage.FilesInfo{}, logger), logger)
	wg           = &sync.WaitGroup{}
	path         = "/home/alex/Dev/test/s1/1.txt"
	//hash         = "d41d8cd98f00b204e9800998ecf8427e"
	hashChanged = "d41d8cd98f00b204e9800998ecf84271"
	file        = storage.FilesInfo{
		Hash:         hashChanged,
		FilePath:     path,
		FileName:     "1.txt",
		LastModified: time.Now().Add(-24 * time.Hour),
		Status:       storage.Sync,
	}
	storageWithOneFile = generated_storage.NewStorageWithLogrus(storage.NewFileStorage(map[string]storage.FilesInfo{file.FileName: file}, logger), logger)
)

func TestDirScanner_Close(t *testing.T) {
	defer cancel()
	type fields struct {
		wg        *sync.WaitGroup
		ctx       context.Context
		sourceDir string
		destDir   string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "simple test", fields: fields{wg: &sync.WaitGroup{}, ctx: ctx, destDir: "/dest", sourceDir: "/src"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DirScanner{
				wg:        tt.fields.wg,
				ctx:       tt.fields.ctx,
				sourceDir: tt.fields.sourceDir,
				destDir:   tt.fields.destDir,
				logger:    logger,
				storage:   storageFiles,
			}
			if err := d.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDirScanner_CopyFile(t *testing.T) {
	defer cancel()
	type fields struct {
		wg        *sync.WaitGroup
		ctx       context.Context
		sourceDir string
		destDir   string
	}
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "simple test", fields: fields{wg: &sync.WaitGroup{}, ctx: ctx, destDir: "/home/alex/Dev/test/s2", sourceDir: "/home/alex/Dev/test/s1"}, args: args{fileName: "1.txt"}, wantErr: false},
		{name: "with source error", fields: fields{wg: &sync.WaitGroup{}, ctx: ctx, destDir: "/home/alex/Dev/test/s2", sourceDir: "/home/alex/Dev/test/s1"}, args: args{fileName: "2.txt"}, wantErr: true},
		//{name: "with dst exist", fields: fields{wg: &sync.WaitGroup{}, ctx: ctx, destDir: "/home/alex/Dev/test/s2", sourceDir: "/home/alex/Dev/test/s1"}, args: args{fileName: "3.txt"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DirScanner{
				wg:        tt.fields.wg,
				ctx:       tt.fields.ctx,
				sourceDir: tt.fields.sourceDir,
				destDir:   tt.fields.destDir,
				logger:    logger,
				storage:   storageFiles,
				syncDone:  syncDone,
			}
			d.wg.Add(1)
			if err := d.CopyFile(tt.args.fileName); (err != nil) != tt.wantErr {
				t.Errorf("CopyFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDirScanner_Run(t *testing.T) {
	defer cancel()
	type fields struct {
		wg        *sync.WaitGroup
		ctx       context.Context
		sourceDir string
		destDir   string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "simple test", fields: fields{wg: &sync.WaitGroup{}, ctx: ctx, destDir: "/home/alex/Dev/test/s2", sourceDir: "/home/alex/Dev/test/s1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DirScanner{
				wg:        tt.fields.wg,
				ctx:       tt.fields.ctx,
				sourceDir: tt.fields.sourceDir,
				destDir:   tt.fields.destDir,
				logger:    logger,
				storage:   storageFiles,
			}
			d.Run()
		})
	}
}

func TestDirScanner_ScanDir(t *testing.T) {
	defer cancel()
	type fields struct {
		wg        *sync.WaitGroup
		ctx       context.Context
		storage   generated_storage.StorageWithLogrus
		sourceDir string
		destDir   string
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "simple test", fields: fields{wg: &sync.WaitGroup{}, ctx: ctx, destDir: "/home/alex/Dev/test/s2", sourceDir: "/home/alex/Dev/test/s1", storage: storageFiles}, wantErr: false},
		{name: "with error", fields: fields{wg: &sync.WaitGroup{}, ctx: ctx, destDir: "/home/alex/Dev/test/s2", sourceDir: "/home/alex/Dev/test/s3", storage: storageWithOneFile}, wantErr: true},
		{name: "exist in file list", fields: fields{wg: &sync.WaitGroup{}, ctx: ctx, destDir: "/home/alex/Dev/test/s2", sourceDir: "/home/alex/Dev/test/s1", storage: storageWithOneFile}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DirScanner{
				wg:        tt.fields.wg,
				ctx:       tt.fields.ctx,
				sourceDir: tt.fields.sourceDir,
				destDir:   tt.fields.destDir,
				logger:    logger,
				storage:   tt.fields.storage,
			}
			d.wg.Add(1)
			if err := d.ScanDir(); (err != nil) != tt.wantErr {
				t.Errorf("ScanDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDirScanner_Wait(t *testing.T) {
	defer cancel()
	type fields struct {
		wg        *sync.WaitGroup
		ctx       context.Context
		sourceDir string
		destDir   string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "simple test", fields: fields{wg: &sync.WaitGroup{}, ctx: ctx, destDir: "/home/alex/Dev/test/s2", sourceDir: "/home/alex/Dev/test/s1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DirScanner{
				wg:        tt.fields.wg,
				ctx:       tt.fields.ctx,
				sourceDir: tt.fields.sourceDir,
				destDir:   tt.fields.destDir,
				logger:    logger,
				storage:   storageFiles,
			}
			d.Wait()
		})
	}
}

func TestNewDirScanner(t *testing.T) {
	defer cancel()
	type args struct {
		srcDir string
		dstDir string
		ctx    context.Context
	}
	tests := []struct {
		name string
		args args
		want *DirScanner
	}{
		{name: "simple test",
			args: args{
				ctx:    ctx,
				dstDir: "/home/alex/Dev/test/s2",
				srcDir: "/home/alex/Dev/test/s1"},
			want: &DirScanner{
				wg:          wg,
				storage:     storageFiles,
				ctx:         ctx,
				destDir:     "/home/alex/Dev/test/s2",
				sourceDir:   "/home/alex/Dev/test/s1",
				filesToSync: filesToSync,
				syncDone:    syncDone,
				logger:      logger}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDirScanner(tt.args.srcDir, tt.args.dstDir, tt.args.ctx, filesToSync, syncDone, logger, storageFiles, wg, 15); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDirScanner() = %v, want %v", got, tt.want)
			}
		})
	}
}
