package keyboards

// Общие кнопки
const ExitButton = "выйти"

const ShowWeatherButton = "показать погоду"

var UserMainMenuChoice = []string{
	ShowWeatherButton,
	ExitButton,
}

const ExitToAdminPanelButton = "выйти в панель админа"
const AddStickerButton = "добавить новый стикер"
const ExitToUserPanelButton = "выйти из админки"

var AdminSetSticker = []string{
	ExitToAdminPanelButton,
	AddStickerButton,
	ExitToUserPanelButton,
}

const SetStickerButton = "задать стикер"

var AdminMainMenuChoice = []string{
	SetStickerButton,
	ExitButton,
}

func GetCustomKeyboard(buttons map[string]string) []string {
	result := make([]string, 0)

	for button, _ := range buttons {
		result = append(result, button)
	}

	result = append(result, ExitButton)

	return result
}
