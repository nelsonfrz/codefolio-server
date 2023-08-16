CREATE TYPE visibility AS ENUM ('private', 'public');
CREATE TABLE projects (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  thumbnail_url TEXT,
  description TEXT,
  content TEXT,
  visibility visibility NOT NULL,
  source_code_url TEXT,
  deployment_url TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);