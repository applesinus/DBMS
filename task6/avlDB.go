package task6

func CreateDB() AVL {
	return AVL{root: nil}
}

func (db *AVL) CreatePool(settings map[string]string, name string) string {
	if db.root == nil {
		db.root = &nodeAVL{
			key:   name,
			value: nil,
		}
		return "ok"
	}
	return "error"
}
