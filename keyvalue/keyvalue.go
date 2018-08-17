package keyvalue

//go:generate mockery -name=KeyValue
type KeyValue interface {
	Delete(key string) error
	GetString(key string) (string, error)
	PutString(key string, value string) error
	GetInt64(key string) (int64, error)
	PutInt64(key string, value int64) error
}
