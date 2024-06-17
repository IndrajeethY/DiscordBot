// database/models.go
package database

import "github.com/jinzhu/gorm"

type User struct {
    gorm.Model
    UserID     string `gorm:"unique_index"`
    Username   string
    Interests  string // Example: comma-separated interests
}

type Server struct {
    gorm.Model
    ServerID   string `gorm:"unique_index"`
    ServerName string
}
