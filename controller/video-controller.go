package controller

import (
	"net/http"
	"strconv"

	"github.com/Sarthak-D97/go_stuAPI/entity"
	"github.com/Sarthak-D97/go_stuAPI/service"
	"github.com/Sarthak-D97/go_stuAPI/validators"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type VideoController interface {
	FindAll(ctx *gin.Context)
	Save(ctx *gin.Context) error
	Update(ctx *gin.Context) error
	Delete(ctx *gin.Context) error
	ShowAll(ctx *gin.Context)
}

type controller struct {
	videoService service.VideoService
}

var validate *validator.Validate

func New(videoService service.VideoService) VideoController {
	validate = validator.New()
	validate.RegisterValidation("is-cool", validators.ValidateCoolTitle)
	return &controller{
		videoService: videoService,
	}
}

func (c *controller) FindAll(ctx *gin.Context) {
	videos, err := c.videoService.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, videos)
}

func (c *controller) Save(ctx *gin.Context) error {
	var video entity.Video
	err := ctx.ShouldBindJSON(&video)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return err
	}
	err = validate.Struct(video)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return err
	}
	c.videoService.Save(video)
	return nil
}

func (c *controller) Update(ctx *gin.Context) error {
	var video entity.Video
	err := ctx.ShouldBindJSON(&video)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return err
	}
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		return err
	}
	video.ID = id
	err = validate.Struct(video)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return err
	}
	c.videoService.Update(video)
	return nil

}
func (c *controller) Delete(ctx *gin.Context) error {
	var video entity.Video
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		return err
	}
	video.ID = id
	c.videoService.Delete(video)
	return nil
}
func (c *controller) ShowAll(ctx *gin.Context) {
	videos, err := c.videoService.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := gin.H{
		"title":  "Video Page",
		"videos": videos,
	}
	ctx.HTML(http.StatusOK, "index.html", data)
}
