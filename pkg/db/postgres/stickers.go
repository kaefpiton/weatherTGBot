package postgres

type Sticker struct {
	ID          int
	title       string
	code        string
	StickerType string
}

func (db *DB) GetStickersCodesByType(stickerType string) ([]string, error) {
	var stickerCodes []string

	query := `SELECT code FROM stickers  where sticker_type = $1`
	rows, err := db.Query(query, stickerType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stickers []*Sticker

	for rows.Next() {
		sticker := new(Sticker)
		err = rows.Scan(&sticker.code)
		if err != nil {
			return nil, err
		}
		stickers = append(stickers, sticker)
	}

	for _, s := range stickers {
		stickerCodes = append(stickerCodes, s.code)
	}

	return stickerCodes, nil
}
