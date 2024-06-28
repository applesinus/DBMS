package interfaces

// collectionInterface - basic collection interface
type CollectionInterface interface {
	Set(key string, secondaryKey string, value string) string
	Update(key string, value string) string
	Get(key string) (string, string)
	GetBySecondaryKey(key string) (string, string)
	GetRange(leftBound string, rightBound string) (*map[string]string, string)
	GetRangeBySecondaryKey(leftBound string, rightBound string) (*map[string]string, string)
	GetAll() (*[]string, *[]string, *[]string, string)
	Delete(key string) string
	Copy() CollectionInterface
}
