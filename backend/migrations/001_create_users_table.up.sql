CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email_hash TEXT NOT NULL UNIQUE,
    encryption_key TEXT NOT NULL UNIQUE,
    encryption_key_hash TEXT NOT NULL UNIQUE,
    reminder_time TIME WITH TIME ZONE NOT NULL,
    timezone TEXT NOT NULL,
    telegram_bot_api_key TEXT NOT NULL,
    telegram_bot_api_key_hash TEXT NOT NULL,
    telegram_user_id TEXT NOT NULL,
    telegram_user_id_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
