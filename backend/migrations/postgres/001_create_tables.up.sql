-- Create the users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email_hash TEXT NOT NULL UNIQUE,
    encryption_key TEXT NOT NULL UNIQUE,
    encryption_key_hash TEXT NOT NULL UNIQUE,
    reminder_time TEXT NOT NULL,
    timezone TEXT NOT NULL,
    telegram_bot_api_key TEXT NOT NULL,
    telegram_bot_api_key_hash TEXT NOT NULL,
    telegram_user_id TEXT NOT NULL,
    telegram_user_id_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create the birthdays table
CREATE TABLE IF NOT EXISTS birthdays (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    name TEXT NOT NULL,
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Trigger function to update updated_at column
CREATE TRIGGER update_users_updated_at
AFTER UPDATE ON users
FOR EACH ROW
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER update_birthdays_updated_at
AFTER UPDATE ON birthdays
FOR EACH ROW
BEGIN
    UPDATE birthdays SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
