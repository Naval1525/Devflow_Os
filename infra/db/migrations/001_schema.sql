-- DevFlow OS v1 — initial schema
-- All tables scoped by user_id

CREATE TABLE IF NOT EXISTS users (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email         TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tasks (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  "type"     TEXT NOT NULL CHECK ("type" IN ('coding', 'leetcode', 'content')),
  date       DATE NOT NULL,
  completed  BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (user_id, "type", date)
);

CREATE TABLE IF NOT EXISTS ideas (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  hook       TEXT NOT NULL,
  idea       TEXT,
  "type"     TEXT NOT NULL CHECK ("type" IN ('reel', 'tweet', 'thread')),
  status     TEXT NOT NULL DEFAULT 'idea' CHECK (status IN ('idea', 'ready', 'posted')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS leetcode_logs (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  problem_name TEXT NOT NULL,
  difficulty   TEXT NOT NULL CHECK (difficulty IN ('easy', 'medium', 'hard')),
  approach     TEXT,
  mistake      TEXT,
  time_taken   INTEGER,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sessions (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  start_time TIMESTAMPTZ NOT NULL,
  end_time   TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS opportunities (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name       TEXT NOT NULL,
  "type"     TEXT NOT NULL CHECK ("type" IN ('job', 'freelance')),
  stage      TEXT NOT NULL CHECK (stage IN ('applied', 'interview', 'closed')),
  source     TEXT,
  notes      TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS finances (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  amount     NUMERIC(12,2) NOT NULL,
  "type"     TEXT NOT NULL CHECK ("type" IN ('salary', 'freelance', 'insta_paid_collab', 'other', 'spend')),
  note       TEXT,
  date       DATE NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_tasks_user_date ON tasks(user_id, date);
CREATE INDEX idx_ideas_user ON ideas(user_id);
CREATE INDEX idx_leetcode_logs_user ON leetcode_logs(user_id);
CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_opportunities_user ON opportunities(user_id);
CREATE INDEX idx_finances_user_date ON finances(user_id, date);
