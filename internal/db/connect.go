package db

import (
	"fmt"

	"github.com/gustavobelfort/42-jitsi/internal/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Init database environment. Creates the database connection and initiates the models managers.
//
// Returns err and do not initiate anything on error.
func Init() error {
	pgConf := config.Conf.Postgres
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		pgConf.User,
		pgConf.Password,
		pgConf.Host,
		pgConf.Port,
		pgConf.DB,
	)
	db, err := gorm.Open("postgres", url)
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&userModel{}, &scaleTeamModel{}).Error; err != nil {
		return err
	}
	if err := db.Model(&userModel{}).AddForeignKey("scale_team_id", "scale_teams(id)", "CASCADE", "CASCADE").Error; err != nil {
		return err
	}
	ScaleTeamManager = &scaleTeamManager{db: db}
	UserManager = &userManager{db: db}
	return nil
}

var (
	ScaleTeamManager ScaleTeamManagerInterface = nil
	UserManager      UserManagerInterface      = nil
)
