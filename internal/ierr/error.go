package ierr

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Error struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Metadata MD     `json:"metadata,omitempty"`
}

var _ fmt.Stringer = (*Error)(nil)

func (e Error) String() string {
	s := "code=" + strconv.Itoa(e.Code) + " message={" + e.Message + "}"
	if len(e.Metadata) > 0 {
		s += " metadata=" + e.Metadata.String()
	}
	return s
}

var _ error = (*Error)(nil)

func (e Error) Error() string {
	return e.String()
}

type MD map[string]string

var _ fmt.Stringer = (*MD)(nil)

func (m MD) String() string {
	var ks []string
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	sb := strings.Builder{}
	for i := 0; i < len(ks); i++ {
		if sb.Len() > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(ks[i])
		sb.WriteString("=")
		sb.WriteString(m[ks[i]])
	}
	return "{" + sb.String() + "}"
}

func New(v ...any) *Error {
	e := &Error{Metadata: MD{}}
	switch len(v) {
	case 0:
	case 1:
		e.Code = v[0].(int)
	case 2:
		e.Message = v[1].(string)
	default:
		e.Message = fmt.Sprintf(v[1].(string), v[2:]...)
	}
	return e
}

func (e *Error) copy(code int) *Error {
	nmd := MD{}
	for k, v := range e.Metadata {
		nmd[k] = v
	}
	if e.Code != 0 {
		code = e.Code
	}
	ne := &Error{Code: code, Message: e.Message, Metadata: nmd}
	return ne
}

func (e *Error) WithCode(code int) *Error {
	ne := e.copy(0)
	ne.Code = code
	return ne
}

func (e *Error) WithMessage(f string, a ...any) *Error {
	ne := e.copy(0)
	ne.Message = fmt.Sprintf(f, a...)
	return ne
}

func (e *Error) WithMetadata(md MD) *Error {
	ne := e.copy(0)
	for k, v := range md {
		ne.Metadata[k] = v
	}
	return ne
}

func (e *Error) SetMetadata(md MD) *Error {
	ne := e.copy(0)
	ne.Metadata = md
	return ne
}

func (e *Error) Status(status int) (int, *Error) {
	return status, e.copy(status)
}

func (e *Error) OK() (int, *Error) {
	return e.Status(http.StatusOK)
}

func (e *Error) Param() (int, *Error) {
	return e.Status(http.StatusBadRequest)
}

func (e *Error) NoAuth() (int, *Error) {
	return e.Status(http.StatusUnauthorized)
}

func (e *Error) Forbidden() (int, *Error) {
	return e.Status(http.StatusForbidden)
}

func (e *Error) NotFound() (int, *Error) {
	return e.Status(http.StatusNotFound)
}

func (e *Error) NotAllowed() (int, *Error) {
	return e.Status(http.StatusMethodNotAllowed)
}

func (e *Error) Conflict() (int, *Error) {
	return e.Status(http.StatusConflict)
}

func (e *Error) Internal() (int, *Error) {
	return e.Status(http.StatusInternalServerError)
}

func OK(v ...any) (int, *Error) {
	return New(v...).OK()
}

func NotFound(v ...any) (int, *Error) {
	return New(v...).NotFound()
}

func NotAllowed(v ...any) (int, *Error) {
	return New(v...).NotAllowed()
}
