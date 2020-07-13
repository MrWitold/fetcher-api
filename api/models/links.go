package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Link struct representing collected links by the user
type Link struct {
	ID       uint64    `gorm:"primary_key;auto_increment" json:"id"`
	URL      string    `gorm:"size:255;not null;unique" json:"url"`
	Interval uint64    `gorm:"size:100;not null;" json:"interval"`
	History  []History `gorm:"foreignkey:LinkID"`
}

// FindAllLinks return list of all links
func (l *Link) FindAllLinks(db *gorm.DB) (*[]Link, error) {
	var err error
	links := []Link{}
	err = db.Debug().Model(&Link{}).Find(&links).Error
	if err != nil {
		return &[]Link{}, err
	}
	return &links, err

}

// FindLinkByID return link with specified ID
func (l *Link) FindLinkByID(db *gorm.DB, uid uint64) (*Link, error) {
	var err error
	err = db.Debug().Model(&Link{}).Where("id = ?", uid).Take(&l).Error
	if err != nil {
		return &Link{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Link{}, errors.New("Link not found")
	}
	return l, err
}

// CreateOrUpdateLink adds or update new record to the database
func (l *Link) CreateOrUpdateLink(db *gorm.DB, url string, interval uint64) (*Link, error) {

	var err error
	err = db.Debug().Model(&Link{}).Where("url = ?", url).Take(&l).Error
	if gorm.IsRecordNotFoundError(err) {
		db.Save(&Link{URL: url, Interval: interval})
		err = db.Debug().Model(&Link{}).Where("url = ?", url).Take(&l).Error
		if err != nil {
			return &Link{}, err
		}

		return l, err
	}
	l.Interval = interval
	db.Save(&l)

	return l, err
}

func (l *Link) DeleteLink(db *gorm.DB, uid uint64) (int64, error) {
	db = db.Debug().Model(&Link{}).Where("id = ? ", uid).Take(&Link{}).Delete(&Link{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Link not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
