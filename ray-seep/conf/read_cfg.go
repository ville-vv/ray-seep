package conf

import (
	"github.com/BurntSushi/toml"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"os"
	"strings"
	"vilgo/vutil"
)

type Reader interface {
	CnfRead(cnf interface{}) error
}

const (
	OriginTypeConfigEnv  = "env"
	OriginTypeConfigToml = "toml"
	OriginTypeConfigJson = "json"
	OriginTypeConfigYaml = "yaml"
	OriginTypeConfigDef  = "default"
)

type ReaderFactory struct {
	fileName string
}

func NewReaderFactory(fileName string) *ReaderFactory {
	return &ReaderFactory{fileName: fileName}
}

func (f *ReaderFactory) GetReader() Reader {
	origin := f.fileName
	format := f.getFileType(origin)
	switch format {
	case OriginTypeConfigJson:
		return NewJsonReader(origin)
	case OriginTypeConfigYaml:
	case OriginTypeConfigToml:
		return NewTomlReader(origin)
	}
	return &DefaultReader{}
}

// GetType 配置文件类型判断
// origin 配置文件路径
func (f *ReaderFactory) getFileType(origin string) string {
	origin = strings.Trim(origin, " ")
	if origin == "" {
		// 空文件路径使用默认配置
		return OriginTypeConfigDef
	}
	if in := vutil.PathExists(origin); !in {
		// 文件不存在使用默认配置
		return OriginTypeConfigDef
	}

	arr := strings.Split(origin, ".")
	if len(arr) > 1 {
		return arr[len(arr)-1]
	}
	return OriginTypeConfigDef
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
}

func (sel *DefaultReader) CnfRead(cnf interface{}) error {
	return nil
}

//-----------------------------------------------------------------------------------
type tomlReader struct {
	BaseReader
}

func NewTomlReader(fName string) *tomlReader {
	return &tomlReader{
		BaseReader{fileName: fName},
	}
}

func (t *tomlReader) CnfRead(cnf interface{}) error {
	return t.readFile(cnf)
}

func (t *tomlReader) readFile(obj interface{}) (err error) {
	buf, err := t.ReadFile(t.fileName)
	if err != nil {
		return
	}
	return toml.Unmarshal(buf, obj)
}

//-----------------------------------------------------------------------------------

type jsonReader struct {
	BaseReader
}

func NewJsonReader(fName string) *jsonReader {
	return &jsonReader{
		BaseReader: BaseReader{fileName: fName},
	}
}

func (t *jsonReader) CnfRead(cnf interface{}) error {
	return t.readFile(cnf)
}

func (t *jsonReader) readFile(obj interface{}) (err error) {
	buf, err := t.ReadFile(t.fileName)
	if err != nil {
		return
	}
	return jsoniter.Unmarshal(buf, obj)
}

//-----------------------------------------------------------------------------------
