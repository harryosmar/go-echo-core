package ctx

import (
	"context"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	coreError "github.com/harryosmar/go-echo-core/error"
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
		Role       Role     `json:"role"`
		Platform   string   `json:"platform"`
	}

	Role struct {
		Id   int64  `json:"id"`
		Code string `json:"code"`
	}

	Session interface {
		GetUserId() string
		GetRole() Role
		GetUserIdInt64() (int64, error)
		GetPrivileges() ([]string, error)
		IsHasPrivilege(privilege string) error
		IsHasPrivileges(privileges []string) error
	}

	session struct {
		UserId     string            `json:"id"`
		Privileges map[string]string `json:"privileges"`
		Role       Role              `json:"role"`
	}
)

func (j JwtClaim) ToJwtClaim() *jwt.MapClaims {
	return &jwt.MapClaims{
		"sub":        j.Sub,
		"iss":        j.Iss,
		"iat":        j.Iat,
		"exp":        j.Exp,
		"jti":        j.Jti,
		"privileges": j.Privileges,
		"role":       j.Role,
		"platform":   j.Platform,
	}
}

func NewSession(jwt *JwtClaim) Session {
	s := &session{UserId: jwt.Sub}
	s.Privileges = s.privilegeListToMap(jwt.Privileges)
	s.Role = jwt.Role
	return s
}

func (s session) privilegeListToMap(privileges []string) map[string]string {
	m := map[string]string{}
	for _, v := range privileges {
		m[v] = v
	}

	return m
}

func (s session) GetRole() Role {
	return s.Role
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
	s := NewContextBuilder(ctx).GetSession()
	if s == nil {
		return ""
	}
	return s.GetUserId()
}

func GetUserIdInt64FromSession(ctx context.Context) (int64, error) {
	s := NewContextBuilder(ctx).GetSession()
	if s == nil {
		return 0, coreError.ErrUnauthorizedAccess
	}
	idStr := s.GetUserId()
	return strconv.ParseInt(idStr, 10, 64)
}

func GetRoleFromSession(ctx context.Context) Role {
	s := NewContextBuilder(ctx).GetSession()
	if s == nil {
		return Role{}
	}
	return s.GetRole()
}
