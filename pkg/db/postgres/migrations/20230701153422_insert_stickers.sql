-- +goose Up
-- +goose StatementBegin
INSERT INTO stickers (title, code, sticker_type) VALUES('sticker1','CAACAgIAAxkBAAMnYUDsTQABlTCRAAHMrhr1oE2THN_JGAACTwADWbv8JXAeFS_YqOxqIAQ','high temperature');
INSERT INTO stickers (title, code, sticker_type) VALUES('sticker2','CAACAgIAAxkBAAMnYUDsTQABlTCRAAHMrhr1oE2THN_JGAACTwADWbv8JXAeFS_YqOxqIAQ','normal temperature');
INSERT INTO stickers (title, code, sticker_type) VALUES('sticker3','CAACAgIAAxkBAAMnYUDsTQABlTCRAAHMrhr1oE2THN_JGAACTwADWbv8JXAeFS_YqOxqIAQ','cold temperature');
INSERT INTO stickers (title, code, sticker_type) VALUES('sticker4','CAACAgIAAxkBAAMnYUDsTQABlTCRAAHMrhr1oE2THN_JGAACTwADWbv8JXAeFS_YqOxqIAQ','frost temperature');
INSERT INTO stickers (title, code, sticker_type) VALUES('sticker5','CAACAgIAAxkBAAMnYUDsTQABlTCRAAHMrhr1oE2THN_JGAACTwADWbv8JXAeFS_YqOxqIAQ','pressure high');
INSERT INTO stickers (title, code, sticker_type) VALUES('sticker6','CAACAgIAAxkBAAMfYUDDPZJqpvnGhXS0IEjsn8hBcv4AApYAAztgJBStCvJmBZEf3yAE','pressure normal');
INSERT INTO stickers (title, code, sticker_type) VALUES('sticker7','CAACAgIAAxkBAAMfYUDDPZJqpvnGhXS0IEjsn8hBcv4AApYAAztgJBStCvJmBZEf3yAE','high wind');
INSERT INTO stickers (title, code, sticker_type) VALUES('sticker8','CAACAgIAAxkBAAMfYUDDPZJqpvnGhXS0IEjsn8hBcv4AApYAAztgJBStCvJmBZEf3yAE','high wind');
INSERT INTO stickers (title, code, sticker_type) VALUES('sticker9','CAACAgIAAxkBAAMfYUDDPZJqpvnGhXS0IEjsn8hBcv4AApYAAztgJBStCvJmBZEf3yAE','normal wind');
INSERT INTO stickers (title, code, sticker_type) VALUES('sticker10','CAACAgIAAxkBAAMfYUDDPZJqpvnGhXS0IEjsn8hBcv4AApYAAztgJBStCvJmBZEf3yAE','low wind');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE stickers;
-- +goose StatementEnd
