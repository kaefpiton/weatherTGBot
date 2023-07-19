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

//todo добавить клавиатуру городов и стикертайпов
