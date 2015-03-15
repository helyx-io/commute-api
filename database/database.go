package database

type DBConnectInfos struct {
    Dialect string
    URL string
    MaxIdelConns int
    MaxOpenConns int
}
