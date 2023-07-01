package postgres

type Sticker struct {
	Stickers_id      int
	Stickers_name    string
	Stickers_code    string
	Stickers_type_id string
}

func (db *DB) GetStickersCodesByType(stickerType string) ([]string, error) {
	var stickerCodes []string

	query := `SELECT code FROM stickers  where sticker_type = $1`
	rows, err := db.Query(query, stickerType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stiskers []*Sticker

	for rows.Next() {
		sticker := new(Sticker)
		err = rows.Scan(&sticker.Stickers_code)
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
