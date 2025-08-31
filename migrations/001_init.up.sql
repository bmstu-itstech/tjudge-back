CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(64) UNIQUE NOT NULL,
    fullname VARCHAR(64) NOT NULL,
    password BYTEA NOT NULL,
    isadmin BOOLEAN DEFAULT FALSE
);

CREATE TABLE teams (
    code VARCHAR(64) PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    leader_id INTEGER NOT NULL,
    contest_id INTEGER NOT NULL,
    max_size INTEGER NOT NULL,
    'CONSTRAINT fk_leader FOREIGN KEY (leader_id) REFERENCES users(id) ON DELETE RESTRICT,'
    CONSTRAINT unique_team_name UNIQUE (name)
);

CREATE TABLE team_members (
    user_id INTEGER NOT NULL,
    team_code VARCHAR(64) NOT NULL,
    PRIMARY KEY (user_id, team_code),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_team FOREIGN KEY (team_code) REFERENCES teams(code) ON DELETE CASCADE
);

CREATE INDEX idx_team_members_team_code ON team_members(team_code);