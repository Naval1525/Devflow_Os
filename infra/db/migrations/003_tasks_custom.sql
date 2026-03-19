-- Allow custom tasks: multiple tasks per day, any label (type = task title)
-- Drop one-task-per-type-per-day constraint and type enum check

ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_user_id_type_date_key;
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_type_check;

-- type column now stores free-form task title/label
COMMENT ON COLUMN tasks."type" IS 'Task label/title (e.g. Coding, Call mom, Review PR)';
