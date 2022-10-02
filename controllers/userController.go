package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"task-5-vix-btpns-HaiqalRamanizarAlFajri/app"
	"task-5-vix-btpns-HaiqalRamanizarAlFajri/helpers"
	"task-5-vix-btpns-HaiqalRamanizarAlFajri/models"

	"github.com/asaskevich/govalidator"
	"github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(ctx *gin.Context) {
	// Setup database
	db := ctx.MustGet("db").(*gorm.DB)

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	// Create Request Body
	requestBody := app.UserRegister{}
	user := models.User{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	// Validate Request Body
	_, err = govalidator.ValidateStruct(requestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	// Converting to model
	mapper := dto.Mapper{}
	mapper.Map(&user, requestBody)

	// Add to database
	err = db.Create(&user).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	// Response
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "S",
		"message": "User Registered",
		"data": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}

func Login(ctx *gin.Context) {
	// Setup database
	db := ctx.MustGet("db").(*gorm.DB)

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	// Create Request Body
	requestBody := app.UserLogin{}
	user := models.User{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	// Validate Request Body
	_, err = govalidator.ValidateStruct(requestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	// Email Checking
	err = db.Where("email = ?", requestBody.Email).First(&user).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "F", "message": "User not found", "data": nil})
		return
	}

	// Verify password
	err = helpers.VerifyPassword(user.Password, requestBody.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "F", "message": "User not valid", "data": nil})
		return
	}

	// Generate JWT
	token, err := app.GenerateJWT(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "S",
		"message": "Login Success",
		"data": gin.H{
			"token": token,
		},
	})
}

func UpdateUser(ctx *gin.Context) {
	// Setup database
	db := ctx.MustGet("db").(*gorm.DB)

	userId := ctx.Param("userId")

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	// Create Request Body
	requestBody := app.UserUpdate{}
	user := models.User{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	// Validate Request Body
	_, err = govalidator.ValidateStruct(requestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	err = db.Model(models.User{}).Where("id = ?", userId).First(&user).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "F", "message": "user not found", "data": nil})
		return
	}

	err = db.Model(&user).Updates(models.User{
		Username: requestBody.Username,
		Password: requestBody.Password,
	}).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "S",
		"message": "User Updated",
		"data": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})

}

func DeleteUser(ctx *gin.Context) {
	// Setup database
	db := ctx.MustGet("db").(*gorm.DB)

	userId := ctx.Param("userId")

	user := models.User{}
	err := db.Model(models.User{}).Where("id = ?", userId).First(&user).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "F", "message": "user not found", "data": nil})
		return
	}

	err = db.Delete(&user).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Response succes
	ctx.JSON(http.StatusOK, gin.H{"status": "T", "message": "delete user success", "data": nil})
}
