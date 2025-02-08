package ctx

import (
	"context"
	coreError "github.com/harryosmar/go-echo-core/error"
	"strconv"
)

type (
	SessionContext struct{}

	JwtClaim struct {
		Iss        string   `json:"iss"`
		Sub        string   `json:"sub"`
		Exp        int64    `json:"exp"`
		Iat        int64    `json:"iat"`
		Jti        string   `json:"jti"`
		Privileges []string `json:"privileges"`
	}

	Session interface {
		GetUserId() string
		GetUserIdInt64() (int64, error)
		GetPrivileges() ([]string, error)
		IsHasPrivilege(privilege string) error
		IsHasPrivileges(privileges []string) error
	}

	session struct {
		UserId     string            `json:"id"`
		Privileges map[string]string `json:"privileges"`
	}
)

func NewSession(jwt *JwtClaim) Session {
	s := &session{UserId: jwt.Sub}
	s.Privileges = s.privilegeListToMap(jwt.Privileges)
	return s
}

func (s session) privilegeListToMap(privileges []string) map[string]string {
	m := map[string]string{}
	for _, v := range privileges {
		m[v] = v
	}

	return m
}

func (s session) GetUserId() string {
	return s.UserId
}

func (s session) GetUserIdInt64() (int64, error) {
	return strconv.ParseInt(s.UserId, 10, 64)
}

func (s session) GetPrivileges() ([]string, error) {
	list := []string{}
	for k, _ := range s.Privileges {
		list = append(list, k)
	}
	return list, nil
}

func (s session) IsHasPrivilege(privilege string) error {
	if _, ok := s.Privileges[privilege]; !ok {
		return coreError.ErrForbiddenPrivilege
	}
	return nil
}

func (s session) IsHasPrivileges(privileges []string) error {
	for _, privilege := range privileges {
		if _, ok := s.Privileges[privilege]; !ok {
			return coreError.ErrForbiddenPrivilege
		}
	}
	return nil
}

func GetUserIdFromSession(ctx context.Context) string {
	return NewContextBuilder(ctx).GetSession().GetUserId()
}

func GetUserIdInt64FromSession(ctx context.Context) (int64, error) {
	idStr := NewContextBuilder(ctx).GetSession().GetUserId()
	return strconv.ParseInt(idStr, 10, 64)
}
