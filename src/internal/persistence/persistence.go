package persistence

import (
	"os"
	"sync"
	"time"

	marshaller "github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/marshaller/persistence"
	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/utils"
	"github.com/gomodule/redigo/redis"
)

// Persistor is a struct designed to persist things to a temporary file to ensure durability
type Persistor struct {
	FileName string
	count    int
	lock     sync.Mutex
}

func writeData(data string, file *os.File) {
	file.Truncate(0)
	file.Seek(0, 0)
	bytes := []byte(data)
	i := 0
	for i < len(bytes) {
		wrote, err := file.Write(bytes[i:])
		if err != nil {
			panic(err)
		}
		i += wrote
	}
	file.Sync()
}

func (p *Persistor) getFileHandle() *os.File {
	file, err := os.OpenFile(p.FileName, os.O_RDWR, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(p.FileName)
			if err != nil {
				panic(err)
			}
			writeData("{}\n", file)
		} else {
			panic(err)
		}
	}
	return file
}

func (p *Persistor) writeFile(content map[int]string) {
	file := p.getFileHandle()
	defer file.Close()

	serialized, err := marshaller.MarshallEntries(content)
	if err != nil {
		panic(err)
	}

	writeData(serialized, file)
}

func (p *Persistor) getEntries() map[int]string {
	file := p.getFileHandle()
	defer file.Close()

	content, err := marshaller.UnmarshallEntries(file)
	if err != nil {
		panic(err)
	}

	return content
}

// PublishEntriesToRedis processes each entry sending them to Redis
func (p *Persistor) PublishEntriesToRedis(conn redis.Conn, redisChannel string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	entries := p.getEntries()

	for id, entry := range entries {
		utils.KeepRetryingAfter(func() (interface{}, error) {
			return conn.Do("PUBLISH", redisChannel, entry)
		}, time.Second*3)
		p.removeEntry(id)
	}
}

// GetEntries returns the content of the persistence file as a map
func (p *Persistor) GetEntries() map[int]string {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.getEntries()
}

func (p *Persistor) addEntry(entry string) int {
	defer func() { p.count++ }()
	content := p.getEntries()
	content[p.count] = entry
	p.writeFile(content)
	return p.count
}

// AddEntry is a function for adding an entry in the persistance file and returns an id for removal later
func (p *Persistor) AddEntry(entry string) int {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.addEntry(entry)
}

func (p *Persistor) removeEntry(id int) {
	content := p.getEntries()
	delete(content, id)
	p.writeFile(content)
}

// RemoveEntry is a function for removing an entry of the persistance file from a given id
func (p *Persistor) RemoveEntry(id int) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.removeEntry(id)
}
