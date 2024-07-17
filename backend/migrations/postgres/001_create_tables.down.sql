-- Drop the trigger on the users table
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop the trigger on the birthdays table
DROP TRIGGER IF EXISTS update_birthdays_updated_at ON birthdays;

-- Drop the trigger function for UTC time enforcement
DROP TRIGGER IF EXISTS enforce_utc_time_trigger ON users;

-- Drop the trigger function for UTC time enforcement
DROP FUNCTION IF EXISTS enforce_utc_time;

-- Drop the trigger function if it's no longer needed
DROP FUNCTION IF EXISTS update_updated_at_column;

-- Drop the birthdays table first
DROP TABLE IF EXISTS birthdays;

-- Drop the users table after dropping the birthdays table
DROP TABLE IF EXISTS users;
