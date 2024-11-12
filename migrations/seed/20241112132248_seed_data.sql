-- +goose Up
-- +goose StatementBegin
INSERT INTO
    users (id, name)
VALUES
    (1, 'Alex'),
    (2, 'Simon'),
    (3, 'Albert'),
    (4, 'Carl'),
    (5, 'James') ON CONFLICT DO NOTHING;

INSERT INTO
    balances (id, user_id, amount)
VALUES
    (1, 1, '0'),
    (2, 2, '0'),
    (3, 3, '0'),
    (4, 4, '0'),
    (5, 5, '0') ON CONFLICT DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM balances WHERE id BETWEEN 1 AND 5;
DELETE FROM users WHERE id BETWEEN 1 AND 5;

-- +goose StatementEnd