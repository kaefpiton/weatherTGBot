package repository

type TLLRepository interface {
	Get(any) (TTLItem, bool)
	Put(any, TTLItem)
}

type TTLItem struct {
	Value      any
	ExpireTime int64
}

func CreateItem(value any) TTLItem {
	return TTLItem{Value: value}
}
