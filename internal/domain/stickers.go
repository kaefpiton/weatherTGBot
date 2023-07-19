package domain

type Sticker struct {
	Name string
	Code string
}

func NewSticker(name string, code string) *Sticker {
	return &Sticker{
		Name: name,
		Code: code,
	}
}
