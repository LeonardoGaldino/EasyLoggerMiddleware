package persistence

import (
	"encoding/json"
	"os"
)

// Persistor is a struct designed to persist things to a temporary file to ensure durability
type Persistor struct {
	FileName string
	count    int
}

func writeData(data string, file *os.File) {
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
	file, err := os.Open(p.FileName)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(p.FileName)
			if err != nil {
				panic(err)
			}
			writeData("{}", file)
		} else {
			panic(err)
		}
	}
	return file
}

func (p *Persistor) loadFile() map[string]interface{} {
	file := p.getFileHandle()
	defer file.Close()

	decoder := json.NewDecoder(file)
	content := make(map[string]interface{})
	err := decoder.Decode(&content)
	if err != nil {
		panic(err)
	}

	return content
}

func (p *Persistor) writeFile(content map[string]interface{}) {
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
func (p *Persistor) AddEntry(entry interface{}) int {
	marshalled, err := json.Marshal(entry)
	if err != nil {
		panic(err)
	}
	defer func() { p.count++ }()
	serialized := string(marshalled)
	content := p.loadFile()
	content[string(p.count)] = serialized
	p.writeFile(content)
	return p.count
}

// RemoveEntry is a function for removing an entry of the persistance file from a given id
func (p *Persistor) RemoveEntry(id int) {
	content := p.loadFile()
	delete(content, string(id))
	p.writeFile(content)
}
