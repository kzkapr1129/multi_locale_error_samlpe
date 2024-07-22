package istm

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var data = struct {
	once sync.Once
	dict map[string]interface{}
}{}

// LoadDictは辞書データを読み込みます
func LoadDict() {
	data.once.Do(func() {
		fn := "config.yaml"
		file, err := os.Open(fn)
		if err != nil {
			panic("failed to load " + fn)
		}
		defer file.Close()

		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&data.dict); err != nil {
			panic("failed to decode " + fn)
		}
	})
}

// GetDictは辞書データから指定したキーの文字列を返却します
func GetDict(keys []string, arg ...interface{}) (string, bool) {
	dict := data.dict
	for _, key := range keys {
		if tmp, ok := dict[key]; !ok {
			return fmt.Sprintf("Invalid dict of '%v'", keys), false
		} else if v, ok := tmp.(map[string]interface{}); ok {
			dict = v
		}
	}

	locale := "jp"
	if d, ok := dict[locale]; !ok {
		return fmt.Sprintf("Invalid locale for '%v': %s", keys, locale), false
	} else if v, ok := d.(string); !ok {
		return fmt.Sprintf("Invalid type for '%v': %s", keys, locale), false
	} else if 0 < len(arg) {
		return fmt.Sprintf(v, arg...), false
	} else {
		return v, true
	}
}
