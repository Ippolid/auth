-- +goose Up
-- Добавляем ограничение уникальности для поля name
ALTER TABLE users_table
    ADD CONSTRAINT unique_name UNIQUE (name);

-- +goose Down
-- Удаляем ограничение уникальности для отката миграции
ALTER TABLE users_table
DROP CONSTRAINT unique_name;
