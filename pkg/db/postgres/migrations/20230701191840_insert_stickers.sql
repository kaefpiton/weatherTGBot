-- +goose Up
-- +goose StatementBegin
-- temperature set
INSERT INTO stickers (title, code, sticker_type, type_id) VALUES('hot pepper','CAACAgIAAxkBAAMxYUSvqEM2N9kh3CxIi3aB3vFvdFcAAiQ4AALpVQUY4gY1Rhbj6zYgBA','high temperature',10);
INSERT INTO stickers (title, code, sticker_type, type_id) VALUES('ok Peach','CAACAgIAAxkBAAMzYUSwX6aCuOQ2viVi4IuJydC_07wAAmk8AALpVQUYKm9sk-lnnEggBA','normal temperature',11);
INSERT INTO stickers (title, code, sticker_type, type_id) VALUES('coquettish kitty','CAACAgIAAxkBAAMvYUDutIbJCmDgh8MgG3R9yNtjCTAAAjsNAAIHAAEhSlDQHnMITvqzIAQ','cold temperature',12);
INSERT INTO stickers (title, code, sticker_type, type_id) VALUES('cold thermometer','xWxjKo_XdRUryen7IAAm4DAAJji0YMsPWPsDU5IhsgBA-lnnEggBA','frost temperature',13);

-- pressure set
INSERT INTO stickers (title, code, sticker_type, type_id) VALUES('panicking crab','CAACAgIAAxkBAAMnYUDsTQABlTCRAAHMrhr1oE2THN_JGAACTwADWbv8JXAeFS_YqOxqIAQ','pressure high',14);
INSERT INTO stickers (title, code, sticker_type, type_id) VALUES('ok crab','CAACAgIAAxkBAAMXYUC8vhd01veJTcwJeBFFtH83nPEAAqcAAxZCawrZczRKlDRzKSAE','pressure normal',15);

-- wind set
INSERT INTO stickers (title, code, sticker_type, type_id) VALUES('panicking lama','CAACAgIAAxkBAAMfYUDDPZJqpvnGhXS0IEjsn8hBcv4AApYAAztgJBStCvJmBZEf3yAE','high wind',16);
INSERT INTO stickers (title, code, sticker_type, type_id) VALUES('ok tiger','CAACAgIAAxkBAAMfYUDDPZJqpvnGhXS0IEjsn8hBcv4AApYAAztgJBStCvJmBZEf3yAE','normal wind',17);
INSERT INTO stickers (title, code, sticker_type, type_id) VALUES('calm kitty','CAACAgIAAxkBAAMnYUDsTQABlTCRAAHMrhr1oE2THN_JGAACTwADWbv8JXAeFS_YqOxqIAQ','low wind',18);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE stickers;
-- +goose StatementEnd