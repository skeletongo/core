package basic

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"path"

	"github.com/skeletongo/core/log"
)

type Package interface {
	Name() string
	Init() error
	io.Closer
}

type Encrypt interface {
	IsCipherText([]byte) bool
	Encrypt([]byte) []byte
	Decode([]byte) []byte
}

var packages = make(map[string]Package)

var configEncrypt Encrypt

func RegisterPackage(p Package) {
	packages[p.Name()] = p
}

func RegisterEncrypt(h Encrypt) {
	configEncrypt = h
}

func LoadPackages(filename string) {
	switch path.Ext(filename) {
	case ".json":
		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			_ = log.Logger.Errorf("Error while reading config file %s: %s", filename, err)
			break
		}
		if configEncrypt != nil {
			if configEncrypt.IsCipherText(bytes) {
				bytes = configEncrypt.Decode(bytes)
			}
		}
		var data interface{}
		if err = json.Unmarshal(bytes, &data); err != nil {
			_ = log.Errorf("Error while Unmarshal data failed %s: %s", filename, err)
			break
		}
		configs := data.(map[string]interface{})
		for name, pkg := range packages {
			cfg, ok := configs[name]
			if !ok {
				_ = log.Warnf("Package %v init data not exist.", pkg.Name())
				continue
			}
			bytes, err := json.Marshal(cfg)
			if err != nil {
				_ = log.Warnf("Package %v marshal data failed.", pkg.Name())
				continue
			}
			if err = json.Unmarshal(bytes, &pkg); err != nil {
				_ = log.Errorf("Error while unmarshalling JSON from config file %s: %s", filename, err)
				continue
			}
			if err = pkg.Init(); err != nil {
				_ = log.Errorf("Error while initializing package %s: %s", pkg.Name(), err)
				continue
			}
			log.Infof("package [%16s] load success", pkg.Name())
		}
	default:
		panic("Unsupported config file: " + filename)
	}
}

func ClosePackages() {
	for _, v := range packages {
		if err := v.Close(); err != nil {
			_ = log.Errorf("Error while closing package %s: %s", v.Name(), err)
		} else {
			log.Infof("package [%16s] close success", v.Name())
		}
	}
}
