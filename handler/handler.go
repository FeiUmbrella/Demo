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
		u.PUT("/update", h.Update)
		u.DELETE("/delete", h.Delete)
		//u.GET("/info", h.MGet)
		//u.POST("/search", h.Search)
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	users := make([]*model.UserEs, 0)

	if err := c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "Invalid argument"})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "Users is empty"})
		return
	}

	if err := h.service.BatchAdd(c, users); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}

func (h *UserHandler) Update(c *gin.Context) {
	users := make([]*model.UserEs, 0)

	if err := c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "Invalid argument"})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "There is no users to update"})
		return
	}

	if err := h.service.BatchUpdate(c, users); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}

func (h *UserHandler) Delete(c *gin.Context) {
	users := make([]*model.UserEs, 0)

	if err := c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "Invalid argument"})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "There is no users to delete"})
		return
	}

	if err := h.service.BatchDel(c, users); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1000, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
}
