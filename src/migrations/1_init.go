package migrations

import (
	"log"

	"github.com/alielmi98/image-processing-service/constants"
	"github.com/alielmi98/image-processing-service/internal/auth/domain/models"
	userModels "github.com/alielmi98/image-processing-service/internal/auth/domain/models"
	imageModels "github.com/alielmi98/image-processing-service/internal/image/domain/models"
	"github.com/alielmi98/image-processing-service/pkg/db"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Up1() {
	database := db.GetDb()

	createTables(database)
	createDefaultUserInformation(database)

}

func createTables(database *gorm.DB) {
	tables := []interface{}{}

	// User
	tables = addNewTable(database, userModels.User{}, tables)
	tables = addNewTable(database, userModels.Role{}, tables)
	tables = addNewTable(database, userModels.UserRole{}, tables)

	// Image
	tables = addNewTable(database, imageModels.Image{}, tables)
	tables = addNewTable(database, imageModels.ProcessingJob{}, tables)
	tables = addNewTable(database, imageModels.ProcessingResult{}, tables)

	err := database.Migrator().CreateTable(tables...)
	if err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Migration, err.Error())
	}
	log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Migration, "tables created")
}

func addNewTable(database *gorm.DB, model interface{}, tables []interface{}) []interface{} {
	if !database.Migrator().HasTable(model) {
		tables = append(tables, model)
	}
	return tables
}

func createDefaultUserInformation(database *gorm.DB) {

	adminRole := models.Role{Name: constants.AdminRoleName}
	createRoleIfNotExists(database, &adminRole)

	defaultRole := models.Role{Name: constants.DefaultRoleName}
	createRoleIfNotExists(database, &defaultRole)

	u := models.User{Username: constants.DefaultUserName, FirstName: "Test", LastName: "Test",
		MobileNumber: "09111112222", Email: "admin@admin.com"}
	pass := "12345678"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)

	createAdminUserIfNotExists(database, &u, adminRole.Id)

}

func createRoleIfNotExists(database *gorm.DB, r *models.Role) {
	exists := 0
	database.
		Model(&models.Role{}).
		Select("1").
		Where("name = ?", r.Name).
		First(&exists)
	if exists == 0 {
		database.Create(r)
	}
}

func createAdminUserIfNotExists(database *gorm.DB, u *models.User, roleId int) {
	exists := 0
	database.
		Model(&models.User{}).
		Select("1").
		Where("username = ?", u.Username).
		First(&exists)
	if exists == 0 {
		database.Create(u)
		ur := models.UserRole{UserId: u.Id, RoleId: roleId}
		database.Create(&ur)
	}
}
