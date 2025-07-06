package responses

import (
	"github.com/gin-gonic/gin"
)

func HTTP200(c *gin.Context, obj any) {
	c.IndentedJSON(200, obj)
}

func HTTP201(c *gin.Context, obj any) {
	c.JSON(201, obj)
}

func HTTP201WithEmpty(c *gin.Context) {
	c.JSON(201, gin.H{})
}

func HTTP400(c *gin.Context, obj any) {
	c.IndentedJSON(400, obj)
}

func HTTP401(c *gin.Context, obj any) {
	c.IndentedJSON(401, obj)
}

func HTTP404(c *gin.Context, obj any) {
	c.IndentedJSON(404, obj)
}

func HTTP500(c *gin.Context, obj any) {
	c.IndentedJSON(500, obj)
}
