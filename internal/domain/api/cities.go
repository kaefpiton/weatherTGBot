package api

// todo при инициализации приложения брать из базы и добавлять в кеш, а пока так
var Cities = map[string]string{
	"Москва":    "Moscow",
	"Ростов":    "Rostov",
	"Агалатово": "Agalatovo",
	"Минск":     "Minsk",
}

// todo зарефакторить это
func GetCitiesKeys(cities map[string]string) []string {
	result := make([]string, 0, 0)
	for k, _ := range cities {
		result = append(result, k)
	}
	return result
}
