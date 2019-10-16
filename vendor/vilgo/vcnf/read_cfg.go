package vcnf

import (
	"encoding/json"
	"errors"
	"github.com/BurntSushi/toml"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"vilgo/vfile"
)

// Reader
type Reader interface {
	CnfRead(cnf interface{}) error
}

// ReaderFactory
type ReaderFactory struct {
	fileName string
	readers  map[string]Reader
}

func NewReader(fileName string) Reader {
	return NewReaderFactory(fileName).GetReader()
}

// NewReaderFactory
func NewReaderFactory(fileName string) *ReaderFactory {
	return &ReaderFactory{fileName: fileName, readers: map[string]Reader{
		"json":    newJsonReader(fileName),
		"yml":     newYamlReader(fileName),
		"yaml":    newYamlReader(fileName),
		"toml":    newTomlReader(fileName),
		"default": NewDefaultReader(fileName),
	}}
}

// AddReader
func (f *ReaderFactory) AddReader(rdType string, rd Reader) {
	f.readers[rdType] = rd
}

// GetReader
func (f *ReaderFactory) GetReader() Reader {
	if rd, ok := f.readers[f.getFileType(f.fileName)]; ok {
		return rd
	}
	return f.readers["default"]
}

// getFileType 配置文件类型判断
// origin 配置文件路径
func (f *ReaderFactory) getFileType(origin string) string {
	origin = strings.Trim(origin, " ")
	if origin == "" {
		// 空文件路径使用默认配置
		return ""
	}
	if in := vfile.PathExists(origin); !in {
		// 文件不存在使用默认配置
		return ""
	}

	arr := strings.Split(origin, ".")
	if len(arr) > 1 {
		return arr[len(arr)-1]
	}
	return ""
}

//-----------------------------------------------------------------------------------
type BaseReader struct {
	fileName string
	fType    string
}

func (b *BaseReader) ReadFile(fileName string) (data []byte, err error) {
	var (
		file *os.File
	)
	if file, err = os.Open(fileName); err != nil {
		return
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

// -----------------------------------------------------------------------------------
type DefaultReader struct {
	BaseReader
	content []byte
}

func NewDefaultReader(fName string) *DefaultReader {
	return &DefaultReader{BaseReader: BaseReader{fileName: fName}}
}

func (sel *DefaultReader) CnfRead(cnf interface{}) error {
	switch sel.fType {
	case "toml":
		return toml.Unmarshal(sel.content, cnf)
	case "json":
		return json.Unmarshal(sel.content, cnf)
	case "yaml", "yml":
		return yaml.Unmarshal(sel.content, cnf)
	default:
		return errors.New("content format is not supported")
	}
}

func (sel *DefaultReader) SetInfo(content string, format string) {
	sel.content = []byte(content)
	sel.fType = strings.ToLower(format)
}

//-----------------------------------------------------------------------------------
type tomlReader struct {
	BaseReader
}

func newTomlReader(fName string) *tomlReader {
	return &tomlReader{
		BaseReader{fileName: fName},
	}
}

func (t *tomlReader) CnfRead(cnf interface{}) error {
	buf, err := t.ReadFile(t.fileName)
	if err != nil {
		return err
	}
	return toml.Unmarshal(buf, cnf)
}

//-----------------------------------------------------------------------------------

type jsonReader struct {
	BaseReader
}

func newJsonReader(fName string) *jsonReader {
	return &jsonReader{
		BaseReader: BaseReader{fileName: fName},
	}
}

func (t *jsonReader) CnfRead(cnf interface{}) error {
	buf, err := t.ReadFile(t.fileName)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(buf, cnf)
}

//-----------------------------------------------------------------------------------

type yamlReader struct {
	BaseReader
}

func newYamlReader(fName string) *yamlReader {
	return &yamlReader{
		BaseReader: BaseReader{fileName: fName},
	}
}

func (t *yamlReader) CnfRead(cnf interface{}) error {
	buf, err := t.ReadFile(t.fileName)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(buf, cnf)
}
