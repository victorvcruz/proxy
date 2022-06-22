package cache

type CacheClient interface {
	ConnectToDatabase() error
	InsertInDatabase(key string, value string) error
	FindInDatabase(key string) (string, error)
}
