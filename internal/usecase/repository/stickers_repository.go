package repository

type StickerType struct {
	Title string
	Alias string
}

type StickersRepository interface {
	GetStickersCodesByType(stickerTypeName string) ([]string, error)
	IsStickerExist(stickerCode string) bool
	GetStickerTypes() []StickerType
	CreateSticker(title, code, categoryTitle string) error
}
