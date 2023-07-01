-- +goose Up
-- +goose StatementBegin
-- temperature set
INSERT INTO stickers (title, code, sticker_type) VALUES('hot pepper','CAACAgIAAxkBAAMxYUSvqEM2N9kh3CxIi3aB3vFvdFcAAiQ4AALpVQUY4gY1Rhbj6zYgBA','high temperature');
INSERT INTO stickers (title, code, sticker_type) VALUES('ok Peach','CAACAgIAAxkBAAMzYUSwX6aCuOQ2viVi4IuJydC_07wAAmk8AALpVQUYKm9sk-lnnEggBA','normal temperature');
INSERT INTO stickers (title, code, sticker_type) VALUES('coquettish kitty','CAACAgIAAxkBAAMvYUDutIbJCmDgh8MgG3R9yNtjCTAAAjsNAAIHAAEhSlDQHnMITvqzIAQ','cold temperature');
INSERT INTO stickers (title, code, sticker_type) VALUES('cold thermometer','xWxjKo_XdRUryen7IAAm4DAAJji0YMsPWPsDU5IhsgBA-lnnEggBA','frost temperature');

-- pressure set
INSERT INTO stickers (title, code, sticker_type) VALUES('panicking crab','CAACAgIAAxkBAAMnYUDsTQABlTCRAAHMrhr1oE2THN_JGAACTwADWbv8JXAeFS_YqOxqIAQ','pressure high');
INSERT INTO stickers (title, code, sticker_type) VALUES('ok crab','CAACAgIAAxkBAAMXYUC8vhd01veJTcwJeBFFtH83nPEAAqcAAxZCawrZczRKlDRzKSAE','pressure normal');

-- wind set
INSERT INTO stickers (title, code, sticker_type) VALUES('panicking lama','CAACAgIAAxkBAAMfYUDDPZJqpvnGhXS0IEjsn8hBcv4AApYAAztgJBStCvJmBZEf3yAE','high wind');
INSERT INTO stickers (title, code, sticker_type) VALUES('ok tiger','CAACAgIAAxkBAAMfYUDDPZJqpvnGhXS0IEjsn8hBcv4AApYAAztgJBStCvJmBZEf3yAE','normal wind');
INSERT INTO stickers (title, code, sticker_type) VALUES('calm kitty','CAACAgIAAxkBAAMnYUDsTQABlTCRAAHMrhr1oE2THN_JGAACTwADWbv8JXAeFS_YqOxqIAQ','low wind');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE stickers;
-- +goose StatementEnd
