package istm

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

/*
IstmErrorはエラーコードから多言語化対応されたエラーメッセージを生成するエラーオブジェクトです。
*/
type IstmError struct {
	message string
}

/*
IstmErrorを生成します。設定ファイルでエラーコードを定義しコード例を参考に使用してください。

設定ファイルの例:

	dict:
	  word:
	    sbom-form-name:
	      jp: "名前"
	      en: "Name"
	  error:
	    E1234:
	      jp: テストエラー
	      en: the test error
	    E1235:
	      jp: "'%s'の型が不正です"
	      en: "'%s' is invalid type"
	    E1236:
	      jp: "'%s'の数値が不正です: %d"
	      en: "The number of '%s' is invalid: %d"

コード例:

	func main(){
	    // エラーコードだけを指定する場合
	    err := NewIstmError("E1234")
	    // エラーコードとパラメータを指定する場合 (パラメータは辞書のキーを指定)
	    err := NewIstmError("E1235", "dict.word.sbom-form-name")
	    // エラーコードとパラメータを指定する場合 (パラメータは辞書を使用しない)
	    err := NewIstmError("E1235", "名前")

	    // 多言語化対応されたエラーメッセージを受け取る
	    errMessage := err.Error()

	    // ラップされたエラーからIstmErrorを取り出す
	    werr := fmt.Errorf("wrapped %w", err)
	    err := Unwrap(werr)
	}
*/
func NewIstmError(errCode string, arg ...interface{}) error {
	message := toString(errCode, arg...)
	return &IstmError{
		message: message,
	}
}

/*
Errorは多言語化対応されたエラーメッセージを返却します。
*/
func (i *IstmError) Error() string {
	return i.message
}

/*
UnwrapはラップされたエラーからIstmErrorを取り出します。
IstmErrorがエラーツリー上に存在しない場合は引数に指定されたエラーをそのまま返します。
*/
func Unwrap(err error) error {
	var target *IstmError
	if errors.As(err, &target) {
		return target
	} else {
		return err
	}
}

var data = struct {
	once sync.Once
	dict map[string]interface{}
}{}

func toString(errCode string, arg ...interface{}) string {
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

	convArg := []interface{}{}
	for _, a := range arg {
		if v, ok := a.(string); ok && strings.Contains(v, ".") {
			// 文字列型かつ文字列中にドット(.)を含んでいる場合は辞書のキーとみなす
			keys := strings.Split(v, ".")
			if str, ok := toStringWithKeys(keys); ok {
				// 辞書の取り出しに成功した場合
				convArg = append(convArg, str) // パラメータをキーとして辞書引きして使用する
			} else {
				// 辞書の取り出しに失敗した場合はキー名をそのまま使用する
				convArg = append(convArg, a) // パラメータをそのまま使用する
			}
		} else {
			// 辞書キー以外の場合
			convArg = append(convArg, a) // パラメータをそのまま使用する
		}
	}

	keys := []string{
		"dict",
		"error",
		errCode,
	}
	str, _ := toStringWithKeys(keys, convArg...)
	// 辞書キーの不正指定を早期発見するため、エラーを無視して変換エラーが格納された文字列をそのまま使用する
	return str
}

func toStringWithKeys(keys []string, arg ...interface{}) (string, bool) {
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
