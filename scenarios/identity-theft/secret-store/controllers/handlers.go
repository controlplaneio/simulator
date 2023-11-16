package controllers

import (
	"encoding/base64"
	"net/http"
	"simulation-scripts/scenario/identity-theft/_app/secret-store/database"
	"simulation-scripts/scenario/identity-theft/_app/secret-store/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @BasePath /api/v1

// GetUsers godoc
// @Summary get all users
// @Schemes
// @Description get all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {string} GetUsers
// @Router /users [get]
func GetUsers(c *gin.Context) {
	// Get all product stored
	users, err := database.FindAll()
	if err != nil {
		c.JSON(http.StatusNotFound, "User Not Found")
	}

	// Send response back
	c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary get user by id
// @Schemes
// @Description get user by id
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {string} GetUserByID
// @Router /users/:userId [get]
func GetUserByID(c *gin.Context) {

	// Set requested parameter
	reqID := c.Param("userId")

	user, err := database.FindByID(reqID)
	if err != nil {
		c.JSON(http.StatusNotFound, "User Not Found")
	} else {
		c.JSON(http.StatusOK, user)
	}

}

// NewUser godoc
// @Summary add a new user
// @Schemes
// @Description add a new user
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {string} NewUser
// @Router /users [post]
func NewUser(c *gin.Context) {

	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid Request")
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "New User Could not be Created")
		return
	}

	secret := base64.StdEncoding.EncodeToString([]byte(user.Secret))

	newUser := models.User{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  string(password),
		Secret:    secret,
		CreatedAt: time.Now(),
	}

	new, err := database.Insert(&newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "New User Could not be Created")
	} else {
		c.JSON(http.StatusCreated, new)
	}
}

// UpdateUser godoc
// @Summary update a user
// @Schemes
// @Description update a user (put)
// @Tags users
// @Accept json
// @Produce json
// @Success 204 {string} UpdateUser
// @Router /users/:userId [put]
func UpdateUser(c *gin.Context) {
	reqID := c.Param("userId")
	var userModel models.User

	if err := c.BindJSON(&userModel); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid Request")
		return
	}

	user, err := database.FindByID(reqID)
	if err != nil {
		c.JSON(http.StatusNotFound, "User Not Found")
		return
	}

	user.FirstName = userModel.FirstName
	user.LastName = userModel.LastName
	user.Email = userModel.Email

	if userModel.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(userModel.Password), 12)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Error: User Failed to Update")
			return
		}

		user.Password = string(password)
	}

	if userModel.Secret != "" {
		secret := base64.StdEncoding.EncodeToString([]byte(userModel.Secret))
		user.Secret = secret
	}

	updatedUser, err := database.Update(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error: User Failed to Update")
	} else {
		c.JSON(http.StatusNoContent, updatedUser)
	}

}

// DeleteUser godoc
// @Summary delete a user
// @Schemes
// @Description delete a user
// @Tags users
// @Accept json
// @Produce json
// @Success 204 {string} DeleteUser
// @Router /users/:userId [delete]
func DeleteUser(c *gin.Context) {
	reqID := c.Param("userId")
	var userModel models.User

	if err := c.BindJSON(&userModel); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid Request")
		return
	}

	user, err := database.FindByID(reqID)
	if err != nil {
		c.JSON(http.StatusNotFound, "User Not Found")
		return
	}

	delete, err := database.Delete(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, "User Not Found")
	} else {
		c.JSON(http.StatusNoContent, delete)
	}
}
