package persistence

import (
	"encoding/json"
	"io"
	"os"
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

	decoder := json.NewDecoder(file)
	content := make(map[int]string)
	err := decoder.Decode(&content)
	if err != nil {
		if err != io.EOF {
			panic(err)
		}
	}

	return content
}

func (p *Persistor) writeFile(content map[int]string) {
	file := p.getFileHandle()
	defer file.Close()

	bytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	serialized := string(bytes)
	writeData(serialized, file)
}

// AddEntry is a function for adding an entry in the persistance file and returns an id for removal later
func (p *Persistor) AddEntry(entry string) int {
	rawIn := json.RawMessage(entry)
	bytes, err := rawIn.MarshalJSON()
	if err != nil {
		panic(err)
	}
	defer func() { p.count++ }()
	serialized := string(bytes)
	content := p.GetEntries()
	content[p.count] = serialized
	p.writeFile(content)
	return p.count
}

// RemoveEntry is a function for removing an entry of the persistance file from a given id
func (p *Persistor) RemoveEntry(id int) {
	content := p.GetEntries()
	delete(content, id)
	p.writeFile(content)
}
