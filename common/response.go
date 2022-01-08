package common

import "github.com/gin-gonic/gin"

func FailedResponse(err error, data interface{}) gin.H {
	result := gin.H{
		"message": err.Error(),
		"status":  "fail",
		"result":  data,
	}

	return result
}

func SucceedResponse(data interface{}) gin.H {
	result := gin.H{
		"message": "",
		"status":  "ok",
		"result":  data,
	}

	return result
}
