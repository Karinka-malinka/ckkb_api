package router

import (
	"net/http"

	"github.com/ckkb_api/internal/controller/handlers"
	"github.com/labstack/echo/v4"
)

type Router struct {
	Echo *echo.Echo
}

func NewEchoRouter(handlers []handlers.Handler) *Router {

	router := &Router{Echo: echo.New()}

	router.Echo.Use(CORSMiddleware())

	ckkbAPIGroup := router.Echo.Group("/ckkbAPI")

	for _, handler := range handlers {
		handler.RegisterHandler(ckkbAPIGroup)
	}

	return router
}

func CORSMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Устанавливаем CORS‑заголовки
			c.Response().Header().Set("Access-Control-Allow-Origin", "https://am-check.ru")
			c.Response().Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")

			// Обрабатываем preflight‑запрос
			if c.Request().Method == "OPTIONS" {
				return c.NoContent(http.StatusNoContent)
			}

			// Передаём управление следующему обработчику
			return next(c)
		}
	}
}
