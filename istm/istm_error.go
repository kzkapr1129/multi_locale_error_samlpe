package istm

import (
	"errors"
	"strings"
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

func toString(errCode string, arg ...interface{}) string {
	LoadDict()
	convArg := []interface{}{}
	for _, a := range arg {
		if v, ok := a.(string); ok && strings.Contains(v, ".") {
			// 文字列型かつ文字列中にドット(.)を含んでいる場合は辞書のキーとみなす
			keys := strings.Split(v, ".")
			if str, ok := GetDict(keys); ok {
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
	str, _ := GetDict(keys, convArg...)
	// 辞書キーの不正指定を早期発見するため、エラーを無視して変換エラーが格納された文字列をそのまま使用する
	return str
}
