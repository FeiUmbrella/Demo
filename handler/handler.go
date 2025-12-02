package handler

import (
	"es_test/Service"
	"es_test/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	engine  *gin.Engine
	service *Service.UserService
}

func NewUserHandler(engine *gin.Engine, service *Service.UserService) *UserHandler {
	return &UserHandler{
		engine:  engine,
		service: service,
	}
}

func (h *UserHandler) Run() {
	// Force log's color
	gin.ForceConsoleColor()
	h.engine.Use(gin.Logger())
	h.engine.Use(gin.Recovery())
	h.registerRouter()

	err := h.engine.Run()
	if err != nil {
		log.Fatalln("server start failed")
	}
}

func (h *UserHandler) registerRouter() {
	u := h.engine.Group("api/user")
	{
		u.POST("/create", h.Create)
		//u.PUT("/update", h.Update)
		//u.DELETE("/delete", h.Delete)
		//u.GET("/info", h.MGet)
		//u.POST("/search", h.Search)
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	users := make([]*model.UserEs, 0)
	user := model.UserEs{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "Invalid argument"})
		return
	}
	users = append(users, &user)
	if err := h.service.BatchAdd(c, users); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}
