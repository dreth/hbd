CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    encryption_key TEXT NOT NULL,
    birthday_list_id UUID,
    reminder_time TIME WITH TIME ZONE NOT NULL,
    telegram_bot_api_key TEXT NOT NULL,
    telegram_user_id TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
