-- Drop the trigger on the users table
DROP TRIGGER IF EXISTS update_users_updated_at;

-- Drop the trigger on the birthdays table
DROP TRIGGER IF EXISTS update_birthdays_updated_at;

-- Drop the birthdays table first
DROP TABLE IF EXISTS birthdays;

-- Drop the users table after dropping the birthdays table
DROP TABLE IF EXISTS users;
