package storage

import (
	"os"

	"github.com/sameer-bishnoi/counter/app/infra/algo"
)

type DataStore interface {
	Load() error
	Store() error
}

type File struct {
	filePath string
	queue    algo.SlidingWindower
}

func NewFile() *File {
	return &File{}
}

func (f *File) WithFilePath(filePath string) {
	f.filePath = filePath
}

func (f *File) WithQueue(queue algo.SlidingWindower) {
	f.queue = queue
}

func (f *File) Load() error {
	file, err := os.OpenFile(f.filePath, os.O_RDWR, os.ModeAppend)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = f.queue.Load(file); err != nil {
		return err
	}

	return os.Truncate(f.filePath, 0)
}

func (f *File) Store() error {
	file, err := os.OpenFile(f.filePath, os.O_RDWR, os.ModeAppend)
	if err != nil {
		return err
	}
	defer file.Close()

	return f.queue.Store(file)
}
