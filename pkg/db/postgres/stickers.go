package postgres

type Sticker struct {
	Stickers_id      int
	Stickers_name    string
	Stickers_code    string
	Stickers_type_id string
}

func (db *DB) GetStickersCodesByType(stickerTypeName string) ([]string, error) {
	var stickerCodes []string

	rows, err := db.Query(
		`SELECT stickers_code FROM stickers INNER JOIN stickertype on stickertype_id = stickers.stickers_type_id where stickerType_name = $1`,
		stickerTypeName)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stiskers := []*Sticker{}

	for rows.Next() {
		sticker := new(Sticker)
		err := rows.Scan(&sticker.Stickers_code)
		if err != nil {
			return nil, err
		}
		stiskers = append(stiskers, sticker)
	}

	for _, s := range stiskers {
		stickerCodes = append(stickerCodes, s.Stickers_code)
	}

	return stickerCodes, nil
}
