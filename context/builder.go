package ctx

import (
	"context"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type (
	ContextBuilder struct {
		ctx context.Context
	}
)

func NewContextBuilder(ctx context.Context) *ContextBuilder {
	return &ContextBuilder{ctx: ctx}
}

func (c *ContextBuilder) SetRequestId(requestId string) *ContextBuilder {
	c.ctx = context.WithValue(c.ctx, ContextRequestId{}, requestId)
	return c
}

func (c *ContextBuilder) SetLogger(entry *log.Entry) *ContextBuilder {
	c.ctx = context.WithValue(c.ctx, ContextLogger{}, entry)
	return c
}

func (c *ContextBuilder) SetSession(s Session) *ContextBuilder {
	c.ctx = context.WithValue(c.ctx, SessionContext{}, s)
	return c
}

func (c *ContextBuilder) GetSession() Session {
	if s, ok := c.ctx.Value(SessionContext{}).(Session); ok {
		return s
	}

	return nil
}

func (c *ContextBuilder) GetLogger() *log.Entry {
	if logEntry, ok := c.ctx.Value(ContextLogger{}).(*log.Entry); ok {
		return logEntry
	}

	return log.NewEntry(CustomLogger)
}

func (c *ContextBuilder) GetRequestId() string {
	value := c.ctx.Value(ContextRequestId{})
	if value == nil {
		value = uuid.New().String()
	}

	return value.(string)
}

func (c *ContextBuilder) Context() context.Context {
	return c.ctx
}
