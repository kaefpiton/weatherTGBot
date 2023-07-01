package postgres

type Sticker struct {
	ID          int
	title       string
	code        string
	stickerType string
}

func (db *DB) GetStickersCodesByType(stickerType string) ([]string, error) {
	var stickerCodes []string

	query := `SELECT code FROM stickers JOIN sticker_types st on stickers.type_id = st.id where st.title= $1`
	rows, err := db.Query(query, stickerType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stickerCode string
		err = rows.Scan(&stickerCode)
		if err != nil {
			return nil, err
		}
		stickerCodes = append(stickerCodes, stickerCode)
	}

	return stickerCodes, nil
}
