-- name: InsertUser :one
INSERT INTO users (
    full_name,
    username,
    password,
    bio,
    profile_picture_url,
    personal_website_url,
    github_url,
    location
  )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;
-- name: SelectAllUsers :many
SELECT *
FROM users
LIMIT $1 OFFSET $2;
;
-- name: SelectAllUsersWithoutPassword :many
SELECT id,
  full_name,
  username,
  bio,
  profile_picture_url,
  personal_website_url,
  github_url,
  location,
  created_at
FROM users
LIMIT $1 OFFSET $2;
-- name: SelectUser :one
SELECT *
FROM users
WHERE id = $1;
-- name: SelectUserByUsername :one
SELECT *
FROM users
WHERE username = $1;
-- name: UpdateUser :one
UPDATE users
SET full_name = $2,
  username = $3,
  password = $4,
  bio = $5,
  profile_picture_url = $5,
  personal_website_url = $6,
  github_url = $7,
  location = $8
WHERE id = $1
RETURNING *;
-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
-- name: InsertProject :one
INSERT INTO projects (
    user_id,
    name,
    thumbnail_url,
    description,
    content,
    visibility,
    source_code_url,
    deployment_url
  )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;
-- name: SelectAllProjects :many
SELECT *
FROM projects
WHERE visibility = 'public'
  OR (
    visibility = 'private'
    AND user_id = $3
  )
LIMIT $1 OFFSET $2;
-- name: SelectProject :one
SELECT *
FROM projects
WHERE id = $1;
-- name: UpdateProct :one
UPDATE projects
SET user_id = $2,
  name = $3,
  thumbnail_url = $4,
  description = $5,
  content = $6,
  visibility = $7,
  source_code_url = $8,
  deployment_url = $9
WHERE id = $1
RETURNING *;
-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1 AND user_id = $2;