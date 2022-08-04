package scanner

type FileScanner interface {
	Run()
	Wait()
	Close() error
	CopyFile(fileName string) error
	ScanDir() error
}
