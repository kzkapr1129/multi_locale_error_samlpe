package main

import (
	"fmt"
	"test/istm"
)

func main() {

	err := istm.NewIstmError("E1236", "dict.word.sbom-form-name", 2)
	newError := istm.RuntimeErrorWrapper(err)

	fmt.Println(istm.Unwrap(err), istm.Unwrap(newError), newError)
}
