package helpers

import (
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func Render(c *gin.Context, status int, template templ.Component) error {
	c.Status(status)
	c.Header("Content-Type", "text/html")
	return template.Render(c.Request.Context(), c.Writer)
}
