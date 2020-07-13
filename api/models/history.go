package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// History struct represents collected links form websites
type History struct {
	Response  string    ` json:"response"`
	Duration  float64   `json:"duration"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	LinkID    uint64    `json:"linkid"`
}

// FindHistoryByID list history of specfied linkID
func (h *History) FindHistoryByID(db *gorm.DB, uid uint64) (*[]History, error) {
	var err error
	historylist := []History{}
	db.Debug().Model(&History{}).Where("link_id = ?", uid).Find(&historylist)
	if err != nil {
		return &[]History{}, err
	}

	return &historylist, err
}

// AddNewHistoryRecord adds new history to database
func (h *History) AddNewHistoryRecord(db *gorm.DB, uid uint64, response string, duration float64) (uint32, error) {
	db.Debug().Save(&History{LinkID: uid, Response: response, Duration: duration, CreatedAt: time.Now()})

	return 1, nil
}
