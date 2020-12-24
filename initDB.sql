DROP TABLE IF EXISTS artists, users, albums, tracks, genres, track_genre, user_track, playlists, track_playlist,
    user_track_like, session CASCADE;

CREATE TABLE artists (
    id serial NOT NULL PRIMARY KEY,
    name varchar(100) NOT NULL UNIQUE,
    description text NOT NULL DEFAULT '',
    poster varchar(100) NOT NULL DEFAULT '',
    avatar varchar(100) NOT NULL DEFAULT ''
);

CREATE TABLE users (
    id serial NOT NULL PRIMARY KEY,
    name varchar(64) NOT NULL UNIQUE CHECK (length(name) > 2),
    email varchar(64) NOT NULL UNIQUE,
    password varchar(64) NOT NULL,
    avatar varchar(255) NOT NULL DEFAULT '',
    is_admin bool NOT NULL DEFAULT false
);

CREATE TABLE albums (
    id serial NOT NULL PRIMARY KEY,
    artist_id int NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
    title varchar(100) NOT NULL,
    poster varchar(100) NOT NULL DEFAULT ''
);

CREATE TABLE tracks (
    id serial NOT NULL PRIMARY KEY,
    album_id int NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
    title varchar(100) NOT NULL,
    duration int NOT NULL DEFAULT 0,
    index int NOT NULL DEFAULT 0,
    audio varchar(100) NOT NULL DEFAULT '',
    likes_count int NOT NULL DEFAULT 0
);

CREATE TABLE genres (
    id serial NOT NULL PRIMARY KEY,
    name varchar(100) NOT NULL UNIQUE
);

CREATE TABLE track_genre (
    track_id int NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    genre_id int NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY(track_id, genre_id)
);

CREATE TABLE user_track (
    user_id int NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id int NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    PRIMARY KEY(user_id, track_id)
);

CREATE TABLE user_track_like (
    user_id int NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id int NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    PRIMARY KEY(user_id, track_id)
);

CREATE TABLE playlists (
    id serial NOT NULL PRIMARY KEY,
    user_id int NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title varchar(100) NOT NULL,
    poster varchar(100) DEFAULT '',
    is_public bool NOT NULL DEFAULT false
);

CREATE TABLE track_playlist (
    track_id int NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    playlist_id int NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    PRIMARY KEY(track_id, playlist_id)
);

CREATE TABLE user_subscriber (
    user_subscriber_id int NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_id int NOT NULL REFERENCES users(id) ON DELETE CASCADE CHECK (user_subscriber_id != user_id),
    PRIMARY KEY(user_subscriber_id, user_id)
);

CREATE TABLE session (
    id varchar(64) NOT NULL PRIMARY KEY,
    userID int NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expire date NOT NULL
);

CREATE OR REPLACE FUNCTION count_index() RETURNS TRIGGER LANGUAGE plpgsql AS $emp_stamp$
    BEGIN
    new.index := COUNT(*) + 1 FROM tracks WHERE tracks.album_id = new.album_id;
    RETURN new;
    END;
$emp_stamp$;

CREATE TRIGGER emp_stamp BEFORE INSERT
    ON tracks
    FOR EACH ROW
    EXECUTE PROCEDURE count_index();

CREATE OR REPLACE FUNCTION track_like() RETURNS TRIGGER LANGUAGE plpgsql AS $track_like$
    BEGIN
    UPDATE tracks SET likes_count = likes_count + 1 WHERE track_id = new.track_id;
    RETURN new;
    END;
$track_like$;

CREATE TRIGGER track_like AFTER INSERT
    ON user_track_like
    FOR EACH ROW
    EXECUTE PROCEDURE track_like();

CREATE OR REPLACE FUNCTION track_dislike() RETURNS TRIGGER LANGUAGE plpgsql AS $track_dislike$
    BEGIN
    UPDATE tracks SET likes_count = likes_count - 1 WHERE track_id = old.track_id;
    RETURN old;
    END;
$track_dislike$;

CREATE TRIGGER track_dislike AFTER DELETE
    ON user_track_like
    FOR EACH ROW
    EXECUTE PROCEDURE track_dislike();

CREATE INDEX IF NOT EXISTS albums_title_index ON albums (title);

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO meuser;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO meuser;
