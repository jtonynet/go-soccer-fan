package routes

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/dto"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/middleware"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/service"
)

/*
	TODO:
	  - Segregar os métodos resolvidos nas rotas para controllers específicas.
*/

type ValidationErrorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ginRoutes struct {
	engine *gin.Engine
}

func NewGinRoutes(
	userService *service.User,
	competitionService *service.Competition,
	fanService *service.Fan,
	broadcastService *service.Broadcast,
) *ginRoutes {
	validate := validator.New()
	e := gin.Default()

	e.POST("/torcedores", func(c *gin.Context) {
		var fReq dto.FanCreateRequest
		if err := c.ShouldBindJSON(&fReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": "corpo da requisição inválido",
			})
			return
		}

		err := validate.Struct(fReq)
		if err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				c.JSON(http.StatusBadRequest, gin.H{
					"erro": formatValidationErrors(validationErrors, &dto.FanCreateRequest{}),
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": "corpo da requisição inválido!",
			})
			return
		}

		fResp, err := fanService.Create(&fReq)
		if err != nil {
			switch err.Error() {
			case "duplicated":
				c.JSON(http.StatusBadRequest, gin.H{
					"erro": "torcedor indisponível para criação.",
				})
				return
			case "not found":
				c.JSON(http.StatusNotFound, gin.H{
					"erro": "time não encontrado.",
				})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"erro": "erro interno, tente novamente mais tarde.",
				})
				return
			}
		}

		c.JSON(http.StatusAccepted, fResp)
	})

	e.POST("/user", func(c *gin.Context) {
		var uReq dto.UserCreateRequest
		if err := c.ShouldBindJSON(&uReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": "corpo da requisição inválido",
			})
			return
		}

		err := validate.Struct(uReq)
		if err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				c.JSON(http.StatusBadRequest, gin.H{
					"erro": formatValidationErrors(validationErrors, &dto.UserCreateRequest{}),
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": "corpo da requisição inválido!",
			})
			return
		}

		uResp, err := userService.Create(&uReq)
		if err != nil {
			if err.Error() == "duplicated" {
				c.JSON(http.StatusBadRequest, gin.H{
					"erro": "usuário indisponível para criação.",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"erro": "erro interno, tente novamente mais tarde",
			})
			return
		}

		c.JSON(http.StatusAccepted, uResp)

	})

	e.POST("/auth/login", func(c *gin.Context) {
		var uReq dto.UserLoginRequest
		if err := c.ShouldBindJSON(&uReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": "corpo da requisição inválido",
			})
			return
		}

		err := validate.Struct(uReq)
		if err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				c.JSON(http.StatusBadRequest, gin.H{
					"erro": formatValidationErrors(validationErrors, &dto.UserLoginRequest{}),
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": "corpo da requisição inválido!",
			})
			return
		}

		uResp, err := userService.Login(&uReq)
		if err != nil {
			switch err.Error() {
			case "not found":
				c.JSON(http.StatusNotFound, gin.H{
					"erro": "usuário não existente.",
				})
				return
			case "incorrect password":
				c.JSON(http.StatusBadRequest, gin.H{
					"erro": "senha incorreta.",
				})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"erro": "erro interno, tente novamente mais tarde",
				})
				return
			}
		}

		c.JSON(http.StatusAccepted, uResp)

	})

	e.GET("/campeonatos", middleware.JwtAuthMiddleware(), func(c *gin.Context) {
		result, err := competitionService.FindAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"erro": "erro interno, tente novamente mais tarde",
			})
			return
		}
		c.JSON(http.StatusOK, result)
	})

	e.GET("/campeonatos/:uid/partidas", middleware.JwtAuthMiddleware(), func(c *gin.Context) {
		uid, err := uuid.Parse(c.Param("uid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": "campeonato com URI ID inválido",
			})
			return
		}

		result, err := competitionService.FindMatchsByCompetitionUID(uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"erro": "erro interno, tente novamente mais tarde",
			})
		}
		c.JSON(http.StatusOK, result)
	})

	e.POST("/broadcast", middleware.JwtAuthMiddleware(), func(c *gin.Context) {
		var bReq dto.BroadcastSendRequest
		if err := c.ShouldBindJSON(&bReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": "corpo da requisição inválido",
			})
			return
		}

		err := validate.Struct(bReq)
		if err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				c.JSON(http.StatusBadRequest, gin.H{
					"erro": formatValidationErrors(validationErrors, &dto.BroadcastSendRequest{}),
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": "corpo da requisição inválido!",
			})
			return
		}

		bResp, err := broadcastService.Publish(&bReq)
		if err != nil {
			if err.Error() == "not found" {
				c.JSON(http.StatusNotFound, gin.H{
					"erro": "time não encontrado.",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"erro": "erro interno, tente novamente mais tarde",
			})
			return
		}

		c.JSON(http.StatusAccepted, bResp)

	})

	return &ginRoutes{
		engine: e,
	}
}

func (gr *ginRoutes) Run(port int) error {
	return gr.engine.Run(fmt.Sprintf(":%v", port))
}

func formatValidationErrors(validationErrors validator.ValidationErrors, obj interface{}) []ValidationErrorMessage {
	var errors []ValidationErrorMessage

	objType := reflect.TypeOf(obj)
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}
	if objType.Kind() != reflect.Struct {
		fmt.Println("obj não é uma struct ou um ponteiro para uma struct")
		return errors
	}

	for _, e := range validationErrors {
		fieldName := e.Field()

		field, ok := objType.FieldByName(fieldName)
		var jsonTag string
		if ok {
			jsonTag = field.Tag.Get("json")
		}
		if jsonTag == "" {
			jsonTag = fieldName
		}

		var message string
		switch e.Tag() {
		case "required":
			message = fmt.Sprintf("%s é obrigatório", jsonTag)
		case "min":
			message = fmt.Sprintf("%s deve ter no mínimo %s caracteres", jsonTag, e.Param())
		case "max":
			message = fmt.Sprintf("%s deve ter no máximo %s caracteres", jsonTag, e.Param())
		case "oneof":
			message = fmt.Sprintf("%s deve ser um dos seguintes valores: %s", jsonTag, strings.Join(strings.Split(e.Param(), " "), ", "))
		default:
			message = fmt.Sprintf("%s não é válido", jsonTag)
		}

		errors = append(errors, ValidationErrorMessage{
			Field:   jsonTag,
			Message: message,
		})
	}
	return errors
}
