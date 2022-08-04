package utils

import (
	"log"
	"testing"
)

const (
	src = "/home/alex/sync/dir1/1.txt"
	dst = "/home/alex/sync/dir2/1.txt"
)

func BenchmarkCopyFilesWithIoCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := CopyFilesWithIoCopy(src, dst)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkCopyFilesWithIoutil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := CopyFilesWithIoutil(src, dst)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkCopyFilesWithOsRW(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := CopyFilesWithOsRW(src, dst)
		if err != nil {
			log.Fatal(err)
		}
	}
}
