// database/database.go
package database

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "fmt"
)

var (
    DB *gorm.DB
)

// Init initializes the SQLite database
func Init(dbPath string) error {
    var err error
    DB, err = gorm.Open("sqlite3", dbPath)
    if err != nil {
        return fmt.Errorf("failed to connect database: %v", err)
    }

    // Automigrate the schema
    DB.AutoMigrate(&User{}, &Server{})

    return nil
}

// Close closes the database connection
func Close() {
    DB.Close()
}
