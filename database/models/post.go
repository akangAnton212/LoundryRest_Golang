package models

import (
	"github.com/jinzhu/gorm"
	"loundry_rest/lib/common"
)

// Post data model
type Post struct {
	gorm.Model
	Text   string `sql:"type:text;"`
	User   User   `gorm:"ForeignKey:UserID"`
	UserID uint	  
}


// Serialize serializes post data
func (p Post) Serialize() common.JSON {
	return common.JSON{
		"id":         p.ID,
		"text":       p.Text,
		// "user":       p.User.Serialize(),
		"user":       p.User.DisplayName,
		"created_at": p.CreatedAt,
	}
}
