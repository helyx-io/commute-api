package config

////////////////////////////////////////////////////////////////////////////////////////////////
/// Imports
////////////////////////////////////////////////////////////////////////////////////////////////

import (
    "os"
    "fmt"
    "log"
    "strconv"
)


////////////////////////////////////////////////////////////////////////////////////////////////
/// Structures
////////////////////////////////////////////////////////////////////////////////////////////////

type Config struct {
    Http          *HttpConfig
    ConnectInfos  *DBConnectInfos
    RedisInfos    *RedisConfig
    LoggerInfos   *LoggerConfig
    TmpDir        string
    BaseURL       string
}

type DBConnectInfos struct {
    Dialect string
    URL string
    MaxIdelConns int
    MaxOpenConns int
}

type HttpConfig struct {
    Port int
}

type RedisConfig struct {
    Host string
    Port int
}

type LoggerConfig struct {
    Path string
}


////////////////////////////////////////////////////////////////////////////////////////////////
/// Functions
////////////////////////////////////////////////////////////////////////////////////////////////

func Init() *Config {

    connectInfos := createConnectInfos()
    redisInfos := createRedisConfig()
    http := createHttpConfig()
    tmpDir := getTmpDir()
    baseURL := getBaseURL()
    loggerInfos := createLoggerConfig()

    return &Config{http, connectInfos, redisInfos, loggerInfos, tmpDir, baseURL}
}

func getBaseURL() string {
    baseURL := os.Getenv("GTFS_BASE_URL")

    log.Printf("[CONFIG] Application - Base URL : '%s'", baseURL)

    return baseURL
}

func getTmpDir() string {
    tmpDir := os.Getenv("GTFS_TMP_DIR")

    log.Printf("[CONFIG] Application - Temp Directory : '%s'", tmpDir)

    return tmpDir
}

func createHttpConfig() *HttpConfig {

    httpPort, _ := strconv.Atoi(os.Getenv("HTTP_PORT"))
    if httpPort == 0 {
        httpPort = 3000
    }

    httpConfig := &HttpConfig{httpPort}

    log.Printf("[CONFIG] Application - HTTP Port : %d", httpConfig.Port)

    return httpConfig
}

func createLoggerConfig() *LoggerConfig {

    loggerFilePath := os.Getenv("LOGGER_FILE_PATH")
    if loggerFilePath == "" {
        loggerFilePath = "/var/log/gtfs-importer/access.log"
    }

    loggerConfig := &LoggerConfig{loggerFilePath}

    log.Printf("[CONFIG] Application - Logger File Path : '%s'", loggerConfig.Path)

    return loggerConfig
}

func createRedisConfig() *RedisConfig {

    redisHost := os.Getenv("REDIS_HOST")

    if redisHost == "" {
        redisHost = "localhost"
    }

    redisPort, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
    if redisPort == 0 {
        redisPort = 8888
    }

    redisInfos := &RedisConfig{redisHost, redisPort}

    log.Printf("[CONFIG] Redis infos - Host : '%s'", redisInfos.Host)
    log.Printf("[CONFIG] Redis infos - Port : '%d'", redisInfos.Port)

    return redisInfos
}

func createConnectInfos() *DBConnectInfos {

    dbDialect := os.Getenv("GTFS_DB_DIALECT")
    if dbDialect == "" {
        dbDialect = "mysql"
    }

    dbHostname := os.Getenv("GTFS_DB_HOSTNAME")
    if dbHostname == "" {
        dbHostname = "localhost"
    }

    dbPort := os.Getenv("GTFS_DB_PORT")
    if dbPort == "" {
        dbPort = "3306"
    }

    dbUsername := os.Getenv("GTFS_DB_USERNAME")
    if dbUsername == "" {
        dbUsername = "gtfs"
    }

    dbPassword := os.Getenv("GTFS_DB_PASSWORD")
    if dbPassword == "" {
        dbPassword = "gtfs"
    }

    dbDatabase := os.Getenv("GTFS_DB_DATABASE")
    if dbDatabase == "" {
        dbDatabase = "gtfs"
    }

    dbURL := os.Getenv("GTFS_DB_URL")
    if dbURL == "" {
        log.Printf("[CONFIG] DB infos - Dialect : '%s'", dbDialect)
        log.Printf("[CONFIG] DB infos - Hostname : '%s'", dbHostname)
        log.Printf("[CONFIG] DB infos - Port : %s", dbPort)
        log.Printf("[CONFIG] DB infos - Database : '%s'", dbDatabase)
        log.Printf("[CONFIG] DB infos - Username : '%s'", dbUsername)
        log.Printf("[CONFIG] DB infos - Password : '%s'", "********")

        if dbDialect == "mysql" {

            dbURL = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4,utf8&parseTime=true", dbUsername, dbPassword, dbHostname, dbPort, dbDatabase)
        } else if dbDialect == "postgres" {
            // postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full
            dbURL = fmt.Sprintf("%s://%v:%v@%v:%v/%v?sslmode=disable", dbDialect, dbUsername, dbDatabase, dbHostname, dbPort, dbDatabase)
        }
    }

    log.Printf("[CONFIG] DB infos - URL : '%s'", dbURL)

    dbMinCnx, _ := strconv.Atoi(os.Getenv("GTFS_DB_MIN_CNX"))
    if dbMinCnx == 0 {
        dbMinCnx = 128
    }

    dbMaxCnx, _ := strconv.Atoi(os.Getenv("GTFS_DB_MAX_CNX"))
    if dbMaxCnx == 0 {
        dbMaxCnx = 128
    }

    log.Printf("[CONFIG] DB infos - Min Connections : %d", dbMinCnx)
    log.Printf("[CONFIG] DB infos - Max Connections : %d", dbMaxCnx)

    return &DBConnectInfos{dbDialect, dbURL, dbMinCnx, dbMaxCnx}
}
