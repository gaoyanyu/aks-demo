package router

import (
	"aks-demo/pkg/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	Router *gin.Engine
)

func init() {
	Router = gin.New()
	Router.UseRawPath = true
	Router.UnescapePathValues = true
	Router.Use(gin.Logger())
	Router.Use(gin.Recovery())

	v1 := Router.Group("")
	initV1Router(v1)
	err := Router.Run(fmt.Sprintf(":%d", 8080))
	if err != nil {
		panic(err)
	}
}

func initV1Router(router *gin.RouterGroup) {
	version := router.Group("/version")
	version.GET("", service.Version)

	k8s := router.Group("/aks")
	k8s.GET("/detail", service.GetAks)
	k8s.POST("", service.CreateAks)
	k8s.PUT("", service.UpdateAks)
	k8s.DELETE("", service.DeleteAks)
}
