package persistence

import (
	"os"

	marshaller "github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/marshaller/persistence"
)

// Persistor is a struct designed to persist things to a temporary file to ensure durability
type Persistor struct {
	FileName string
	count    int
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

// GetEntries returns the content of the persistence file as a map
func (p *Persistor) GetEntries() map[int]string {
	file := p.getFileHandle()
	defer file.Close()

	content, err := marshaller.UnmarshallEntries(file)
	if err != nil {
		panic(err)
	}

	return content
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

// AddEntry is a function for adding an entry in the persistance file and returns an id for removal later
func (p *Persistor) AddEntry(entry string) int {
	defer func() { p.count++ }()
	content := p.GetEntries()
	content[p.count] = entry
	p.writeFile(content)
	return p.count
}

// RemoveEntry is a function for removing an entry of the persistance file from a given id
func (p *Persistor) RemoveEntry(id int) {
	content := p.GetEntries()
	delete(content, id)
	p.writeFile(content)
}
