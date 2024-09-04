package meta

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func getParameterStringToInt(ctx *gin.Context, Name string, defaultValue int) int {
	if val, ok := ctx.GetQuery(Name); ok {
		if val, err := strconv.Atoi(val); err == nil {
			return val
		}
	}
	return defaultValue
}
