package repository

type StickersRepository interface {
	GetStickersCodesByType(stickerTypeName string) ([]string, error)
	IsStickerExist(stickerCode string) bool
	GetStickerTypes() []string
	CreateSticker(title, code, categoryTitle string) error
}
