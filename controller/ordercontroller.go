package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/WhileCodingDoLearn/gpt_project/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OrderService interface {
	Find(ctx *gin.Context)
	Save(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type orderService struct {
	validate   *validator.Validate
	repository repository.RepositoryHandler
}

type IOrderInteface = repository.IOrder

func NewOrderService(repository repository.RepositoryHandler) OrderService {

	return &orderService{validate: validator.New(validator.WithRequiredStructEnabled()), repository: repository}
}

func (os *orderService) Find(ctx *gin.Context) {

}

/*
1. JSON-Daten an die Struct binden
2. Validierung von jedem Element im Array. Ggf : Fehler sammeln und als Antwort zurückgeben
3. Daten in interface umwandeln
4. An DB handler übergeben
5. Resonse 200 OK
*/
func (os *orderService) Save(ctx *gin.Context) {
	var input PostOrderInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	for i, order := range input.Orders {
		if err := os.validate.Struct(order); err != nil {
			validationErrors := make([]string, 0)
			for _, fieldErr := range err.(validator.ValidationErrors) {
				validationErrors = append(validationErrors, fmt.Sprintf("Order[%d] %s: %s", i, fieldErr.Field(), fieldErr.Tag()))
			}
			ctx.JSON(http.StatusBadRequest, gin.H{"error": strings.Join(validationErrors, ",")})
			return
		}
	}

	var orders []IOrderInteface
	for _, o := range input.Orders {
		orders = append(orders, o)
	}

	vals, err := os.repository.AddItems(orders)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	ctx.JSON(http.StatusAccepted, gin.H{"msg": "succress", "data": vals})
}

func (os *orderService) Update(ctx *gin.Context) {

}

func (os *orderService) Delete(ctx *gin.Context) {

}
