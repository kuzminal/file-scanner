package wrappers

//go:generate gowrap gen -p ../scanner -i FileScanner -t logrus -o ../scanner/generated_scanner/scanner_with_log.go
//go:generate gowrap gen -p ../storage -i Storage -t logrus -o ../storage/generated_storage/storage_with_log.go
