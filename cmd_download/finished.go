package cmd_download

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

type FinishedFiles struct {
	File_path         string
	File              *os.File
	Lock              sync.RWMutex
	Finished_file_map map[string]int64
	f_writer          *bufio.Writer
}

func NewFinishedFiles() *FinishedFiles {
	return &FinishedFiles{
		Finished_file_map: map[string]int64{},
	}
}

func (f *FinishedFiles) ReadFinishFileList(finishedFilePath string) error {
	f.Lock.Lock()
	defer f.Lock.Unlock()

	content, _ := os.ReadFile(finishedFilePath)
	if content != nil {
		files := strings.Split(string(content), "\n")
		for _, v := range files {
			f.Finished_file_map[v] = 1
		}
	}

	file, err := os.OpenFile(finishedFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	f.File_path = finishedFilePath
	f.File = file

	f.f_writer = bufio.NewWriter(file)

	return nil
}

func (f *FinishedFiles) Close() {
	if f.File != nil {
		f.File.Close()
		f.File = nil
	}
}

func (f *FinishedFiles) IsFinished(fileName string) bool {
	f.Lock.RLock()
	defer f.Lock.RUnlock()
	_, exist := f.Finished_file_map[fileName]
	return exist
}

func (f *FinishedFiles) DownloadFinish(fileName string) {
	f.Lock.Lock()
	defer f.Lock.Unlock()
	f.Finished_file_map[fileName] = 1

	f.f_writer.WriteString(fileName + "\n")
	f.f_writer.Flush()
}
