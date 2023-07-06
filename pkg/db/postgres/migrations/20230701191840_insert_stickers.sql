-- +goose Up
-- +goose StatementBegin
-- temperature set
INSERT INTO stickers (title, code, type_id)select 'hot pepper','CAACAgIAAxkBAAMxYUSvqEM2N9kh3CxIi3aB3vFvdFcAAiQ4AALpVQUY4gY1Rhbj6zYgBA',id from sticker_types where title = 'high temperature';

INSERT INTO stickers (title, code, type_id) select'ok Peach','CAACAgIAAxkBAAMzYUSwX6aCuOQ2viVi4IuJydC_07wAAmk8AALpVQUYKm9sk-lnnEggBA',id from sticker_types where title = 'normal temperature';

INSERT INTO stickers (title, code, type_id) select'coquettish kitty','CAACAgIAAxkBAAMvYUDutIbJCmDgh8MgG3R9yNtjCTAAAjsNAAIHAAEhSlDQHnMITvqzIAQ',id from sticker_types where title = 'cold temperature';

INSERT INTO stickers (title, code, type_id) select'cold thermometer','CAACAgIAAxkBAAM1YUS4Az_9i-xWxjKo_XdRUryen7IAAm4DAAJji0YMsPWPsDU5IhsgBA',id from sticker_types where title = 'frost temperature';

-- pressure set
INSERT INTO stickers (title, code, type_id) select'panicking crab','CAACAgIAAxkBAAMZYUC9EJbzQAABsY2j_ZT66RHqWtnMAAKgAAMWQmsKBAICCOb4TIYgBA',id from sticker_types where title = 'pressure high';


INSERT INTO stickers (title, code, type_id) select'ok crab','CAACAgIAAxkBAAMXYUC8vhd01veJTcwJeBFFtH83nPEAAqcAAxZCawrZczRKlDRzKSAE',id from sticker_types where title = 'pressure normal';


-- wind set
INSERT INTO stickers (title, code, type_id) select'panicking lama','CAACAgIAAxkBAAMfYUDDPZJqpvnGhXS0IEjsn8hBcv4AApYAAztgJBStCvJmBZEf3yAE',id from sticker_types where title = 'high wind';

INSERT INTO stickers (title, code, type_id) select'o  k tiger','CAACAgIAAxkBAAMjYUDoxXrQ6XLQA0kHcRE0LinmqR8AAmYAA1m7_CWcvJiSv7eHFyAE',id from sticker_types where title = 'normal wind';

INSERT INTO stickers (title, code, type_id) select'calm kitty','CAACAgIAAxkBAAMnYUDsTQABlTCRAAHMrhr1oE2THN_JGAACTwADWbv8JXAeFS_YqOxqIAQ',id from sticker_types where title = 'low wind';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE stickers;
-- +goose StatementEnd