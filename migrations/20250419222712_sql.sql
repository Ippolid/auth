-- +goose Up
ALTER TABLE users_table
ALTER COLUMN name TYPE TEXT,
    ALTER COLUMN email TYPE TEXT,
    ALTER COLUMN password TYPE TEXT;

-- Добавляем ограничения отдельно
ALTER TABLE users_table
    ADD CONSTRAINT check_name_length CHECK (length(name) <= 100),
    ADD CONSTRAINT check_email_length CHECK (length(email) <= 100),
    ADD CONSTRAINT check_password_length CHECK (length(password) <= 100);

-- +goose Down
ALTER TABLE users_table
ALTER COLUMN name TYPE VARCHAR(100),
    ALTER COLUMN email TYPE VARCHAR(100),
    ALTER COLUMN password TYPE VARCHAR(100);

-- Удаляем ограничения
ALTER TABLE users_table
DROP CONSTRAINT check_name_length,
    DROP CONSTRAINT check_email_length,
    DROP CONSTRAINT check_password_length;
