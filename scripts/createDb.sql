CREATE TABLE IF NOT EXISTS movies (
    id INTEGER PRIMARY KEY,
    title TEXT,
    filepath TEXT NOT NULL,
    length REAL,
    isIndexed INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS shows (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS episodes (
    id INTEGER PRIMARY KEY,
    showId INTEGER,
    title TEXT,
    episodeNumber INTEGER,
    season INTEGER,
    filepath TEXT NOT NULL,
    length REAL,
    isIndexed INTEGER NOT NULL,
    FOREIGN KEY(showId) REFERENCES show(id)
);
CREATE TABLE IF NOT EXISTS session (
    id INTEGER PRIMARY KEY,
    movieId INTEGER,
    episodeId INTEGER,
    isPaused INTEGER NOT NULL,
    isPlaying INTEGER NOT NULL
);