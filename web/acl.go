package web

import (
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func (s *server) requireUsername(requiredUsername string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			contextUsername, ok := c.Get("username").(string)
			if !ok {
				s.logger.Error("username not found in the context")
				return c.String(401, "Unauthorized")
			}
			if contextUsername != requiredUsername {
				s.logger.Error("context username doesn't match the required one",
					zap.String("requiredUsername", requiredUsername),
					zap.String("contextUsername", contextUsername),
				)
				return c.String(403, "Forbidden")
			}
			return next(c)
		}
	}
}
