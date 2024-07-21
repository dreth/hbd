-- Create the users table
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email_hash TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    reminder_time TEXT NOT NULL,
    timezone TEXT NOT NULL,
    telegram_bot_api_key TEXT NOT NULL,
    telegram_bot_api_key_hash TEXT NOT NULL,
    telegram_user_id TEXT NOT NULL,
    telegram_user_id_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create the birthdays table
CREATE TABLE birthdays (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    date DATE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Trigger to automatically update the updated_at column on users table update
CREATE TRIGGER update_users_updated_at
AFTER UPDATE ON users
FOR EACH ROW
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Trigger to automatically update the updated_at column on birthdays table update
CREATE TRIGGER update_birthdays_updated_at
AFTER UPDATE ON birthdays
FOR EACH ROW
BEGIN
    UPDATE birthdays SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
