package config

////////////////////////////////////////////////////////////////////////////////////////////////
/// Imports
////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"os"
    "fmt"
    "log"
	"strconv"
    "github.com/jinzhu/gorm"
    "github.com/helyx-io/gtfs-api/database"
    "github.com/helyx-io/gtfs-api/database/mysql"
)

////////////////////////////////////////////////////////////////////////////////////////////////
/// Variables
////////////////////////////////////////////////////////////////////////////////////////////////

var (
	DB            *gorm.DB
	Http          *HttpConfig
	ConnectInfos  *database.DBConnectInfos
    RedisInfos    *RedisConfig
	BaseURL       string
)

////////////////////////////////////////////////////////////////////////////////////////////////
/// Structures
////////////////////////////////////////////////////////////////////////////////////////////////

type HttpConfig struct {
    Port int
}

type RedisConfig struct {
    Host string
    Port int
}

////////////////////////////////////////////////////////////////////////////////////////////////
/// Functions
////////////////////////////////////////////////////////////////////////////////////////////////

func Init() error {

	var err error

	dbDialect := "mysql"

	dbUsername := os.Getenv("DB_USERNAME")
	if dbUsername == "" {
		dbUsername = "gtfs"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "gtfs"
	}

	dbDatabase := os.Getenv("DB_DATABASE")
	if dbDatabase == "" {
		dbDatabase = "gtfs"
	}

	log.Printf("[CONFIG] DB infos - Database : '%s'", dbDatabase)
	log.Printf("[CONFIG] DB infos - Username : '%s'", dbUsername)
	log.Printf("[CONFIG] DB infos - Password : '%s'", "********")

	dbURL := fmt.Sprintf("%v:%v@/%v?charset=utf8mb4,utf8&parseTime=true", dbUsername, dbPassword, dbDatabase)

	dbMinCnx, _ := strconv.Atoi(os.Getenv("DB_MIN_CNX"))
	if dbMinCnx == 0 {
		dbMinCnx = 2
	}

	dbMaxCnx, _ := strconv.Atoi(os.Getenv("DB_MAX_CNX"))
	if dbMaxCnx == 0 {
		dbMaxCnx = 100
	}

	log.Printf("[CONFIG] DB infos - Min Connections : %d", dbMinCnx)
	log.Printf("[CONFIG] DB infos - Max Connections : %d", dbMaxCnx)

	ConnectInfos = &database.DBConnectInfos{dbDialect, dbURL, dbMinCnx, dbMaxCnx}

	// Init Gorm
	if DB, err = mysql.InitDB(ConnectInfos); err != nil {
		return err
	}

    redisHost := os.Getenv("REDIS_HOST")

    if redisHost == "" {
        redisHost = "localhost"
    }

    redisPort, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
    if redisPort == 0 {
        redisPort = 8888
    }

    RedisInfos = &RedisConfig{redisHost, redisPort}

    log.Printf("[CONFIG] Redis infos - Host : '%s'", RedisInfos.Host)
    log.Printf("[CONFIG] Redis infos - Port : '%d'", RedisInfos.Port)

	BaseURL = os.Getenv("HTTP_BASE_URL")

	log.Printf("[CONFIG] Application - Base URL : '%s'", BaseURL)

	httpPort, _ := strconv.Atoi(os.Getenv("HTTP_PORT"))
	if httpPort == 0 {
		httpPort = 4000
	}

	Http = &HttpConfig{httpPort}

	log.Printf("[CONFIG] Application - HTTP Port : %d", Http.Port)

	return nil
}

func Close() {
	if DB != nil {
		defer DB.Close()
	}
}
