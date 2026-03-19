-- Add LinkedIn as an idea type
ALTER TABLE ideas DROP CONSTRAINT IF EXISTS ideas_type_check;
ALTER TABLE ideas ADD CONSTRAINT ideas_type_check CHECK ("type" IN ('reel', 'tweet', 'thread', 'linkedin'));
