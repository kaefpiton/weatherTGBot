package repository

type City struct {
	Title string
	Alias string
}

type CitiesRepository interface {
	GetCites() ([]City, error)
}
