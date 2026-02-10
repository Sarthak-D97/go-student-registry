package controller

import (
	"net/http"
	"strconv"

	"github.com/Sarthak-D97/go_stuAPI/internal/types"
	"github.com/Sarthak-D97/go_stuAPI/internal/utils/response"
	"github.com/Sarthak-D97/go_stuAPI/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type StudentController interface {
	Create(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	GetList(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type studentController struct {
	service  service.StudentService
	validate *validator.Validate
}

func NewStudentController(service service.StudentService) StudentController {
	return &studentController{
		service:  service,
		validate: validator.New(),
	}
}

func (c *studentController) Create(ctx *gin.Context) {
	var student types.Student
	if err := ctx.ShouldBindJSON(&student); err != nil {
		_ = response.WriteJson(ctx.Writer, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	if err := c.validate.Struct(student); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			_ = response.WriteJson(ctx.Writer, http.StatusBadRequest, response.ValidationError(ve))
			return
		}
		_ = response.WriteJson(ctx.Writer, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	created, err := c.service.CreateStudent(student)
	if err != nil {
		_ = response.WriteJson(ctx.Writer, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	_ = response.WriteJson(ctx.Writer, http.StatusCreated, map[string]interface{}{
		"status":  response.StatusCreated,
		"student": created,
	})
}

func (c *studentController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	intID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		_ = response.WriteJson(ctx.Writer, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	student, err := c.service.GetStudentByID(intID)
	if err != nil {
		_ = response.WriteJson(ctx.Writer, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	_ = response.WriteJson(ctx.Writer, http.StatusOK, map[string]interface{}{
		"status":  response.StatusOK,
		"student": student,
	})
}

func (c *studentController) GetList(ctx *gin.Context) {
	students, err := c.service.GetAllStudents()
	if err != nil {
		_ = response.WriteJson(ctx.Writer, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	_ = response.WriteJson(ctx.Writer, http.StatusOK, map[string]interface{}{
		"status":   response.StatusOK,
		"students": students,
	})
}

func (c *studentController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	intID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		_ = response.WriteJson(ctx.Writer, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	var student types.Student
	if err := ctx.ShouldBindJSON(&student); err != nil {
		_ = response.WriteJson(ctx.Writer, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	if err := c.validate.Struct(student); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			_ = response.WriteJson(ctx.Writer, http.StatusBadRequest, response.ValidationError(ve))
			return
		}
		_ = response.WriteJson(ctx.Writer, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	if err := c.service.UpdateStudent(intID, student); err != nil {
		_ = response.WriteJson(ctx.Writer, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	_ = response.WriteJson(ctx.Writer, http.StatusOK, map[string]interface{}{
		"status":     response.StatusOK,
		"student_id": intID,
	})
}

func (c *studentController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	intID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		_ = response.WriteJson(ctx.Writer, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	if err := c.service.DeleteStudent(intID); err != nil {
		_ = response.WriteJson(ctx.Writer, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	_ = response.WriteJson(ctx.Writer, http.StatusOK, map[string]interface{}{
		"status": response.StatusOK,
		"msg":    "student deleted successfully",
	})
}

