package database

import (
	"simulation-scripts/scenario/identity-theft/_app/secret-store/config"
	"simulation-scripts/scenario/identity-theft/_app/secret-store/models"
)

func GetUserByEmail(email string) (models.User, error) {

	user := models.User{}

	if result := config.DB.Db.Select("id", "first_name", "last_name", "email", "created_at", "updated_at", "password").Where("email = ?", email).Find(&user); result.Error != nil {
		return user, result.Error
	}

	return user, nil

}

// All method will be used to get all records from the users table.
func FindAll() ([]models.User, error) {
	var users []models.User

	if result := config.DB.Db.Select("id", "first_name", "last_name", "email", "created_at", "updated_at", "password", "secret").Find(&users); result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

// FindByID will be used to find a new user by id
func FindByID(id string) (models.User, error) {
	user := models.User{}

	if result := config.DB.Db.Where("id = ?", id).First(&user); result.Error != nil {
		return user, result.Error
	}

	return user, nil

}

// Insert will be used to insert a new user
func Insert(user *models.User) (*models.User, error) {
	if err := config.DB.Db.Create(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

// Update a user by id
func Update(user *models.User) (*models.User, error) {

	if err := config.DB.Db.Save(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

// Delete will be used to delete a user by id
func Delete(user *models.User) (*models.User, error) {
	if err := config.DB.Db.Delete(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}
