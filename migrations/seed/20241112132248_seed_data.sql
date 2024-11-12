-- +goose Up
-- +goose StatementBegin
INSERT INTO
    users (name)
VALUES
    ('Alex'),
    ('Simon'),
    ('Albert'),
    ('Carl'),
    ('James');

INSERT INTO
    balances (user_id, amount)
VALUES
    (1, '0'),
    (2, '0'),
    (3, '0'),
    (4, '0'),
    (5, '0');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM balances WHERE id BETWEEN 1 AND 5;
DELETE FROM users WHERE id BETWEEN 1 AND 5;

-- +goose StatementEnd