package controller

import (
	"net/http"
	"strconv"

	"github.com/Sarthak-D97/go_stuAPI/entity"
	"github.com/Sarthak-D97/go_stuAPI/service"
	"github.com/gin-gonic/gin"
)

// StudentController defines the interface for the controller
type StudentController interface {
	Create(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	GetList(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type studentController struct {
	service service.StudentService
}

// NewStudentController creates a new instance of the controller
func NewStudentController(service service.StudentService) StudentController {
	return &studentController{
		service: service,
	}
}

// Create - POST /api/students/
func (c *studentController) Create(ctx *gin.Context) {
	var student entity.Student

	// Bind JSON body to struct
	if err := ctx.ShouldBindJSON(&student); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	err := c.service.Create(&student)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, student)
}

// GetByID - GET /api/students/:id
func (c *studentController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Casting to uint because your Service expects uint (based on previous steps)
	student, err := c.service.FindByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	ctx.JSON(http.StatusOK, student)
}

// GetList - GET /api/students/
func (c *studentController) GetList(ctx *gin.Context) {
	students, err := c.service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, students)
}

// Update - PUT /api/students/:id
func (c *studentController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var student entity.Student
	if err := ctx.ShouldBindJSON(&student); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// FIX: Cast to int because entity.Student.ID is an int
	student.ID = int(id)

	err = c.service.Update(&student)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, student)
}

// Delete - DELETE /api/students/:id
func (c *studentController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Casting to uint for the Service call
	err = c.service.Delete(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}
