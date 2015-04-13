package database

import (
    "fmt"
    "log"
    "database/sql"
    "github.com/jinzhu/gorm"
    "github.com/helyx-io/gtfs-api/data"
    "github.com/helyx-io/gtfs-api/utils"
    "github.com/helyx-io/gtfs-api/config"

    _ "github.com/lib/pq"
    _ "github.com/go-sql-driver/mysql"
)

type DBConnectInfos struct {
    Dialect string
    URL string
    MaxIdelConns int
    MaxOpenConns int
}

func InitDB(dbInfos *config.DBConnectInfos) (*gorm.DB, error) {
    db, err := gorm.Open(dbInfos.Dialect, dbInfos.URL)

    if err != nil {
        return nil, err
    }

    db.DB()

    // Then you could invoke `*sql.DB`'s functions with it
    db.DB().Ping()

    db.DB().SetMaxIdleConns(dbInfos.MaxIdelConns)
    db.DB().SetMaxOpenConns(dbInfos.MaxOpenConns)

    db.SingularTable(true)

    return &db, nil
}

func Rows(db *gorm.DB, connectInfos *config.DBConnectInfos, filename string, params ...interface{}) (*sql.Rows, error) {
    filePath := fmt.Sprintf("resources/ddl/%s/%s.sql", connectInfos.Dialect, filename)
    log.Printf("Executing query query from file path: '%s'", filePath)

    dml, err := data.Asset(filePath)
    utils.FailOnError(err, fmt.Sprintf("Could get dml resource at path '%s' for exec", filePath))
    execStmt := fmt.Sprintf(string(dml), params...)
    log.Printf("Exec Stmt: '%s' - Params: %v", execStmt, params)

    return db.Raw(execStmt).Rows()
}
