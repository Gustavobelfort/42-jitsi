package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

// GetOption will modify the query to apply the wanted parameters to it.
type GetOption func(*gorm.DB) *gorm.DB

/*
 * ScaleTeam Get Options
 */

// ScaleTeamIDOption adds condition if the ScaleTeam id is `id`.
func ScaleTeamIDOption(id int) GetOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

// ScaleTeamNotifiedOption adds condition if the ScaleTeam is `notified`.
func ScaleTeamNotifiedOption(notified bool) GetOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("notified IS ?", notified)
	}
}

// ScaleTeamBeginAtAfterOption adds condition if the ScaleTeam begins at or before `beginAt`
func ScaleTeamBeginAtBeforeOption(beginAt time.Time) GetOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("begin_at <= ?", beginAt.Format(time.RFC3339))
	}
}

// ScaleTeamBeginAtAfterOption adds condition if the ScaleTeam begins at or after `beginAt`
func ScaleTeamBeginAtAfterOption(beginAt time.Time) GetOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("begin_at >= ?", beginAt.Format(time.RFC3339))
	}
}

// ScaleTeamBeginAtInOption adds condition if the ScaleTeam begins in `duration` or less.
func ScaleTeamBeginAtInOption(duration time.Duration) GetOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("begin_at - interval ? <= NOW()", duration.String())
	}
}

/*
 * User Get Options
 */

// ScaleTeamIDOption adds condition if the ScaleTeam id is `id`.
func UserIDOption(id int) GetOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

// UserStatusOption adds condition if User's status is `status`.
func UserStatusOption(status UserStatus) GetOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}

// UserStatusOption adds condition if User's login is `login`.
func UserLoginOption(login string) GetOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("login = ?", login)
	}
}

// UserStatusOption adds condition if User's ScaleTeam id is `scaleTeamId`.
func UserScaleTeamOption(scaleTeamId int) GetOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("scale_team_id = ?", scaleTeamId)
	}
}
