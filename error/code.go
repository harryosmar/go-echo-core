package error

import (
	"fmt"
	"net/http"
)

type (
	CodeErr int

	CodeErrEntity struct {
		Code         string
		ErrorMessage string // detailed error message
		Message      string // displayed to user
		Status       int
		Args         []interface{}
	}

	TypeOfErrEntity interface {
		GetCodeErrEntity() CodeErrEntity
	}
)

func (c CodeErrEntity) Error() string {
	return c.ErrorMessage
}

// Convention Name Rule for Error Code
// eg :ERR400000 => ERRXXX0YY
// XXX : 3 digits status code 200, 201, 400, 401, 403, 404, 500, ...
// 0 : to identify error coming from `core library`
// YY: increment start from 00, 01, 02, ...
var (
	codeErrMap = map[CodeErr]CodeErrEntity{
		ErrGeneral:            {Code: "ERR500000", Status: http.StatusInternalServerError, Message: "Internal Server Error"},
		ErrUnauthorizedAccess: {Code: "ERR401001", Status: http.StatusUnauthorized, Message: "Unauthorized"},
		ErrUnauthorizedAccessSessionContextNotFound:  {Code: "ERR401002", Status: http.StatusUnauthorized, Message: "Unauthorized"},
		ErrUnauthorizedAccessSub404:                  {Code: "ERR401003", Status: http.StatusUnauthorized, Message: "Unauthorized"},
		ErrUnauthorizedAccessJti404:                  {Code: "ERR401004", Status: http.StatusUnauthorized, Message: "Unauthorized"},
		ErrUnauthorizedAccessPrivileges404:           {Code: "ERR401005", Status: http.StatusUnauthorized, Message: "Unauthorized"},
		ErrUnauthorizedAccessPrivilegesInvalidFormat: {Code: "ERR401006", Status: http.StatusUnauthorized, Message: "Unauthorized"},
		ErrNotFound:                            {Code: "ERR404007", Status: http.StatusNotFound, Message: "Not Found"},
		ErrValidation:                          {Code: "ERR400008", Status: http.StatusBadRequest, Message: "Bad Request"},
		ErrInvalidCredentials:                  {Code: "ERR400009", Status: http.StatusBadRequest, Message: "Invalid Credentials"},
		ErrForbidden:                           {Code: "ERR403010", Status: http.StatusForbidden, Message: "Forbidden"},
		ErrForbiddenPrivilege:                  {Code: "ERR403011", Status: http.StatusForbidden, Message: "Forbidden"},
		ErrForbiddenStatus:                     {Code: "ERR403012", Status: http.StatusForbidden, Message: "Forbidden"},
		ErrTooManyRequests:                     {Code: "ERR429013", Status: http.StatusTooManyRequests, Message: "Too Many Requests"},
		ErrUnauthorizedAccessRoleInvalidFormat: {Code: "ERR401014", Status: http.StatusUnauthorized, Message: "Unauthorized"},
		ErrForbiddenMultipleSessionDetected:    {Code: "ERR403015", Status: http.StatusUnauthorized, Message: "Multiple Session detected."},
		ErrUnauthorizedSessionNotFound:         {Code: "ERR401016", Status: http.StatusUnauthorized, Message: "Session not found."},
		ErrForbiddenInvalidRole:                {Code: "ERR403017", Status: http.StatusUnauthorized, Message: "Forbidden role not allowed."},
		ErrUnauthorizedAccessPlatform404:       {Code: "ERR401018", Status: http.StatusUnauthorized, Message: "Platform not found."},
	}
)

func AppendCodeErrMap(err CodeErr, entity CodeErrEntity) {
	codeErrMap[err] = entity
}

const (
	ErrGeneral CodeErr = iota
	ErrUnauthorizedAccess
	ErrUnauthorizedAccessSessionContextNotFound
	ErrUnauthorizedAccessSub404
	ErrUnauthorizedAccessJti404
	ErrUnauthorizedAccessPrivileges404
	ErrUnauthorizedAccessPrivilegesInvalidFormat
	ErrUnauthorizedAccessRoleInvalidFormat
	ErrUnauthorizedAccessPlatform404
	ErrNotFound
	ErrValidation
	ErrForbidden
	ErrForbiddenPrivilege
	ErrTooManyRequests
	ErrInvalidCredentials
	ErrForbiddenStatus
	ErrForbiddenMultipleSessionDetected
	ErrUnauthorizedSessionNotFound
	ErrForbiddenInvalidRole
)

func (c CodeErr) Error() string {
	return codeErrMap[c].ErrorMessage
}

func (c CodeErr) String() string {
	errEntity := codeErrMap[c]
	if errEntity.Args == nil {
		return errEntity.ErrorMessage
	}
	return fmt.Sprintf(errEntity.Message, errEntity.Args...)
}

func (c CodeErr) Code() string {
	return codeErrMap[c].Code
}

func (c CodeErr) Status() int {
	return codeErrMap[c].Status
}

// GetCodeErrEntity must be instanced of TypeOfErrEntity interface
func (c CodeErr) GetCodeErrEntity() CodeErrEntity {
	entity := codeErrMap[c]
	return entity
}

func (c CodeErr) WithArgs(args ...interface{}) CodeErrEntity {
	entity := codeErrMap[c]
	entity.Args = args
	return entity
}

func (c CodeErr) WithError(err error) CodeErrEntity {
	codeErrEntity := codeErrMap[c]
	codeErrEntity.ErrorMessage = err.Error()
	return codeErrEntity
}
