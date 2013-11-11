package fileobj

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type FileObj struct {
	Path string
	Obj  interface{}
}

func New(filePath string, obj interface{}) (*FileObj, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	// try reading from file
	file, err := os.OpenFile(absPath, os.O_RDONLY, 0644)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	} else if err == nil {
		z, err := gzip.NewReader(file)
		if err != nil {
			file.Close()
			return nil, err
		}
		defer func() {
			z.Close()
			file.Close()
		}()
		err = gob.NewDecoder(z).Decode(obj)
		if err != nil {
			return nil, err
		}
	}

	ret := &FileObj{
		Path: absPath,
		Obj:  obj,
	}
	return ret, nil
}

func (self *FileObj) Save() error {
	tmpFilePath := self.Path + fmt.Sprintf(".%d%08d", time.Now().UnixNano(), rand.Int31())
	f, err := os.OpenFile(tmpFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	z := gzip.NewWriter(f)
	defer func() {
		z.Close()
		f.Close()
	}()
	err = gob.NewEncoder(z).Encode(self.Obj)
	if err != nil {
		return err
	}
	err = os.Rename(tmpFilePath, self.Path)
	if err != nil {
		return err
	}
	return nil
}
