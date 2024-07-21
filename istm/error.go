package istm

import "github.com/go-errors/errors"

const RUNTIME_ERROR = "runtime_error"

func RuntimeErrorWrapper(err error) error {
	return errors.WrapPrefix(err, RUNTIME_ERROR, 1)
}
