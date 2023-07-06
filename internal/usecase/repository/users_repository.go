package repository

type StickersRepository interface {
	GetStickersCodesByType(stickerTypeName string) ([]string, error)
}
