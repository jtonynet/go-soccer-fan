package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/dto"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/middleware"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/service"
)

var Validator *validator.Validate

type ginRoutes struct {
	engine *gin.Engine
}

/*
	TODO:
	  - Segregar os métodos resolvidos nas rotas para controllers específicas.
	  - Melhorar o tratamento de erros do validador.
*/

func NewGinRoutes(
	userService *service.User,
	competitionService *service.Competition,
	fanService *service.Fan,
	broadcastService *service.Broadcast,
) *ginRoutes {
	e := gin.Default()

	initValidator()

	e.POST("/torcedores", func(c *gin.Context) {
		var fReq dto.FanCreateRequest
		if err := c.ShouldBindJSON(&fReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": "corpo da requisição inválido",
			})
			return
		}

		if err := validateStruct(&fReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": getValidationErrorMessages(err),
			})
			return
		}

		fResp, err := fanService.Create(&fReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"erro": "erro interno, tente novamente mais tarde",
			})
			return
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

		if err := validateStruct(&uReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": getValidationErrorMessages(err),
			})
			return
		}

		uResp, err := userService.Create(&uReq)
		if err != nil {
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

		if err := validateStruct(&uReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": getValidationErrorMessages(err),
			})
			return
		}

		uResp, err := userService.Login(&uReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"erro": "erro interno, tente novamente mais tarde",
			})
			return
		}

		c.JSON(http.StatusAccepted, uResp)

	})

	e.GET("/campeonatos", middleware.JwtAuthMiddleware(), func(c *gin.Context) {
		result, _ := competitionService.FindAll()
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

		if err := validateStruct(&bReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": getValidationErrorMessages(err),
			})
			return
		}

		bResp, err := broadcastService.Publish(&bReq)
		if err != nil {
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

func initValidator() {
	Validator = validator.New()
}

func validateStruct(s interface{}) error {
	if err := Validator.Struct(s); err != nil {
		return err
	}
	return nil
}

func getValidationErrorMessages(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		var messages string
		for _, e := range errs {
			messages += fmt.Sprintf("Erro no campo '%s': %s; ", e.Field(), e.Tag())
		}
		return messages
	}
	return "validação falhou"
}
