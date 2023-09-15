package download

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

type ChunkFetchStat struct {
	metaPath string
	metaFile *os.File
	chunkMap map[string]int64
	lock     sync.RWMutex
	bw       *bufio.Writer
}

func NewChunkFetchStat() *ChunkFetchStat {
	return &ChunkFetchStat{
		chunkMap: map[string]int64{},
	}
}

func (cs *ChunkFetchStat) LoadChunkMeta(metaPath string) error {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	content, _ := os.ReadFile(metaPath)
	if content != nil {
		files := strings.Split(string(content), "\n")
		for _, v := range files {
			cs.chunkMap[v] = 1
		}
	}

	file, err := os.OpenFile(metaPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	cs.metaPath = metaPath
	cs.metaFile = file

	cs.bw = bufio.NewWriter(file)

	return nil
}

func (cs *ChunkFetchStat) Close() {
	if cs.metaFile != nil {
		cs.metaFile.Close()
		cs.metaFile = nil
	}
}

func (cs *ChunkFetchStat) IsDone(name string) bool {
	cs.lock.RLock()
	defer cs.lock.RUnlock()
	_, exist := cs.chunkMap[name]
	return exist
}

func (cs *ChunkFetchStat) SetDone(name string) {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	cs.chunkMap[name] = 1

	cs.bw.WriteString(name + "\n")
	cs.bw.Flush()
}
