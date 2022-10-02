package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"task-5-vix-btpns-HaiqalRamanizarAlFajri/app"
	"task-5-vix-btpns-HaiqalRamanizarAlFajri/models"

	"github.com/asaskevich/govalidator"
	"github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePhoto(ctx *gin.Context) {
	// Setup database
	db := ctx.MustGet("db").(*gorm.DB)
	tokenString := ctx.GetHeader("Authorization")

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	// Create Request Body
	requestBody := app.PhotoCreated{}
	photo := models.Photo{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	user := models.User{}
	err = db.Model(models.User{}).Where("id = ?", requestBody.UserId).First(&user).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "F", "message": "user id not found", "data": nil})
		return
	}
	data, err := app.GetData(strings.Split(tokenString, "Bearer ")[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "F", "message": "Token not valid", "data": nil})
		return
	}
	if data.ID != user.ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "F", "message": "User unauthorized", "data": nil})
		return
	}

	// Validate Request Body
	_, err = govalidator.ValidateStruct(requestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	result := db.Model(models.Photo{}).Where("user_id = ?", requestBody.UserId).First(&photo)
	if result.RowsAffected > 0 {
		ctx.JSON(http.StatusConflict, gin.H{"status": "F", "message": "photo exist", "data": nil})
		return
	}

	// Converting to model
	mapper := dto.Mapper{}
	mapper.Map(&photo, requestBody)

	// Add to database
	err = db.Create(&photo).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	// Response
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "S",
		"message": "Photo Created",
		"data": gin.H{
			"id":        photo.ID,
			"title":     photo.Title,
			"caption":   photo.Caption,
			"photo_url": photo.PhotoUrl,
			"user_id":   photo.UserId,
		},
	})
}

func GetPhotos(ctx *gin.Context) {
	// Setup database
	db := ctx.MustGet("db").(*gorm.DB)

	var photos []models.Photo
	err := db.Find(&photos).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "S",
		"message": "Success",
		"data":    photos,
	})
}

func UpdatePhoto(ctx *gin.Context) {
	// Setup database
	db := ctx.MustGet("db").(*gorm.DB)

	photoId := ctx.Param("photoId")
	tokenString := ctx.GetHeader("Authorization")

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	// Create Request Body
	requestBody := app.PhotoUpdate{}
	photo := models.Photo{}
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

	user := models.User{}
	err = db.Model(models.User{}).Where("id = ?", requestBody.UserId).First(&user).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "F", "message": "user id not found", "data": nil})
		return
	}
	data, err := app.GetData(strings.Split(tokenString, "Bearer ")[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "F", "message": "Token not valid", "data": nil})
		return
	}
	if data.ID != user.ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "F", "message": "User unauthorized", "data": nil})
		return
	}

	err = db.Model(models.User{}).Where("id = ?", photoId).First(&photo).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "F", "message": "photo not found", "data": nil})
		return
	}

	err = db.Model(&photo).Updates(models.Photo{
		Title:    requestBody.Title,
		Caption:  requestBody.Caption,
		PhotoUrl: requestBody.PhotoUrl,
		UserId:   requestBody.UserId,
	}).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "S",
		"message": "Photo Updated",
		"data":    photo,
	})

}

func DeletePhoto(ctx *gin.Context) {
	// Setup database
	db := ctx.MustGet("db").(*gorm.DB)

	photoId := ctx.Param("photoId")
	tokenString := ctx.GetHeader("Authorization")

	photo := models.Photo{}
	err := db.Model(models.User{}).Where("id = ?", photoId).First(&photo).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "F", "message": "photo not found", "data": nil})
		return
	}
	data, err := app.GetData(strings.Split(tokenString, "Bearer ")[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "F", "message": "Token not valid", "data": nil})
		return
	}
	if data.ID != photo.UserId {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "F", "message": "User unauthorized", "data": nil})
		return
	}

	err = db.Delete(&photo).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Response succes
	ctx.JSON(http.StatusOK, gin.H{"status": "S", "message": "delete photo success", "data": nil})
}
