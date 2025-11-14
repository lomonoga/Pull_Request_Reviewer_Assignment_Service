CREATE TABLE IF NOT EXISTS "user" (
    "user_id" VARCHAR(256) NOT NULL PRIMARY KEY,
    "username" VARCHAR(256) NOT NULL,
    "is_active" BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS "team" (
    "team_name" VARCHAR(256) NOT NULL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS "team_member" (
    "team_name" VARCHAR(256) NOT NULL REFERENCES "team"("team_name") ON DELETE CASCADE,
    "user_id" VARCHAR(256) NOT NULL REFERENCES "user"("user_id") ON DELETE CASCADE,
    PRIMARY KEY ("team_name", "user_id")
);

CREATE TYPE "pull_request_status" AS ENUM ('OPEN', 'MERGED');

CREATE TABLE IF NOT EXISTS "pull_request" (
    "pull_request_id" VARCHAR(256) NOT NULL PRIMARY KEY,
    "pull_request_name" VARCHAR(256) NOT NULL,
    "author_id" VARCHAR(256) NOT NULL REFERENCES "user"("user_id") ON DELETE CASCADE,
    status "pull_request_status" NOT NULL DEFAULT 'OPEN',
    "reviewer_1" VARCHAR(256) REFERENCES "user"("user_id") ON DELETE CASCADE DEFAULT NULL,
    "reviewer_2" VARCHAR(256) REFERENCES "user"("user_id") ON DELETE CASCADE DEFAULT NULL,
    "mergedAt" TIMESTAMP DEFAULT NULL
);