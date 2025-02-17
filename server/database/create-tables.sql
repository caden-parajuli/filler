BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS players (
  id TEXT PRIMARY KEY,
  gameId INTEGER,
  FOREIGN KEY(gameId) REFERENCES games(id)
) WITHOUT ROWID;

CREATE TABLE IF NOT EXISTS games (
  id INTEGER PRIMARY KEY,
  player1Id TEXT,
  player2Id TEXT,
  board TEXT,
  turn TEXT
);

COMMIT;
