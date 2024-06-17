package interfaces

// collectionInterface - basic collection interface
type CollectionInterface interface {
	Set(key string, value string) string
	Update(key string, value string) string
	Get(key string) (string, string)
	GetRange(leftBound string, rightBound string) (*map[string]string, string)
	Delete(key string) string
	Copy() CollectionInterface
}
