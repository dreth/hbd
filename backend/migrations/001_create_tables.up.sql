-- Create the users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email_hash TEXT NOT NULL UNIQUE,
    encryption_key TEXT NOT NULL UNIQUE,
    encryption_key_hash TEXT NOT NULL UNIQUE,
    reminder_time TIME NOT NULL,
    timezone TEXT NOT NULL,
    telegram_bot_api_key TEXT NOT NULL,
    telegram_bot_api_key_hash TEXT NOT NULL,
    telegram_user_id TEXT NOT NULL,
    telegram_user_id_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Trigger function to enforce UTC time on reminder_time column
CREATE OR REPLACE FUNCTION enforce_utc_time()
RETURNS TRIGGER AS $$
BEGIN
    NEW.reminder_time := NEW.reminder_time AT TIME ZONE 'UTC';
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to enforce UTC time on reminder_time column
CREATE TRIGGER enforce_utc_time_trigger
BEFORE INSERT OR UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION enforce_utc_time();

-- Create the birthdays table
CREATE TABLE birthdays (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    name TEXT NOT NULL,
    date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Trigger function to update updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to automatically update the updated_at column on users table update
CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();
