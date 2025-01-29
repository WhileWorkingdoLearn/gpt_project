package endpoint

import (
	"net/http"

	"github.com/WhileCodingDoLearn/gpt_project/controller"
	"github.com/WhileCodingDoLearn/gpt_project/repository"
	"github.com/gin-gonic/gin"
)

var orderService controller.OrderService

func init() {
	/*
		postVar = "dive,required" +
			",required,min=3" + // OrderId muss vorhanden und min. 3 Zeichen sein
			",required,min=3" + // Name muss vorhanden und min. 3 Zeichen sein
			",required,oneof='Neu' 'In Bearbeitung' 'Abgeschlossen'" + // Status muss in der Liste sein
			",required,min=1" // Data muss vorhanden und min. 1 Zeichen sein
	*/
}

func StartSever(port string, middleware []gin.HandlerFunc, handler repository.RepositoryHandler) error {

	orderService = controller.NewOrderService(handler)

	r := gin.Default()

	g := r.Group("/v1")

	for _, m := range middleware {
		g.Use(m)
	}

	g.GET("/orders", orderService.Find)

	g.POST("/orders", orderService.Save)

	g.POST("/import", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	})

	g.PATCH("/orders/{id}", orderService.Update)

	g.DELETE("/orders/{id}", orderService.Delete)

	err := r.Run(port)
	if err != nil {
		return err
	}
	return nil
}
