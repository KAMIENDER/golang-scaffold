package bizerror

import (
	"fmt"

	"github.com/pkg/errors"
)

// e *BizError error
type BizError struct {
	msg      string
	oriError error
	code     string
}

func NewBizErrorWithStack(msg, code string, oriError error) *BizError {
	return &BizError{
		msg:      msg,
		code:     code,
		oriError: errors.Wrap(oriError, ""),
	}
}

func NewBizError(msg, code string, oriError error) *BizError {
	return &BizError{
		msg:      msg,
		code:     code,
		oriError: oriError,
	}
}

func (e *BizError) Error() string {
	return fmt.Sprintf("msg: \n%s\nori error:%s\n", e.msg, e.oriError)
}

func (e *BizError) Code() string {
	return e.code
}

func (e *BizError) Msg() string {
	return e.msg
}
