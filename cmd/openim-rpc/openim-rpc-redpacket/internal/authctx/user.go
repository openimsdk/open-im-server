package authctx

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

const opUserIDKey = "opUserID"

type userIDContextKey struct{}

func WithCurrentUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey{}, strings.TrimSpace(userID))
}

func CurrentUserID(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", fmt.Errorf("request context is nil")
	}
	if userID, ok := ctx.Value(userIDContextKey{}).(string); ok && strings.TrimSpace(userID) != "" {
		return strings.TrimSpace(userID), nil
	}
	if userID, ok := ctx.Value(opUserIDKey).(string); ok && strings.TrimSpace(userID) != "" {
		return strings.TrimSpace(userID), nil
	}
	return "", fmt.Errorf("op user id missing in context")
}

func BindCurrentUserID(c *gin.Context) error {
	if c == nil {
		return fmt.Errorf("gin context is nil")
	}
	userID := strings.TrimSpace(c.GetString(opUserIDKey))
	if userID == "" {
		if value := c.Request.Context().Value(opUserIDKey); value != nil {
			if fromCtx, ok := value.(string); ok {
				userID = strings.TrimSpace(fromCtx)
			}
		}
	}
	if userID == "" {
		return fmt.Errorf("op user id missing in context")
	}
	c.Request = c.Request.WithContext(WithCurrentUserID(c.Request.Context(), userID))
	return nil
}
