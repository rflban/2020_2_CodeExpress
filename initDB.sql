\c musicexpress;

DROP TABLE IF EXISTS artists, users, albums, tracks, genres, track_genre, user_track, session CASCADE;

CREATE TABLE artists (
    id serial NOT NULL PRIMARY KEY,
    name varchar(100) NOT NULL UNIQUE,
    poster varchar(100) DEFAULT ''
);

CREATE TABLE users (
    id serial PRIMARY KEY,
    name varchar(64) NOT NULL UNIQUE,
    email varchar(64) NOT NULL UNIQUE,
    password varchar(64) NOT NULL,
    avatar varchar(255) DEFAULT ''
);

CREATE TABLE albums (
    id serial NOT NULL PRIMARY KEY,
    artist_id int NOT NULL,
    title varchar(100) NOT NULL,
    poster varchar(100) DEFAULT '',
    FOREIGN KEY(artist_id) REFERENCES artists(id) ON DELETE CASCADE
);

CREATE TABLE tracks (
    id serial NOT NULL PRIMARY KEY,
    album_id int NOT NULL,
    title varchar(100) NOT NULL,
    duration int NOT NULL default 0,
    index int NOT NULL default 0,
    audio varchar(100) not NULL default '',
    FOREIGN KEY(album_id) REFERENCES albums(id) ON DELETE CASCADE
);

CREATE TABLE genres (
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL UNIQUE
);

CREATE TABLE track_genre (
    track_id int NOT NULL,
    genre_id int NOT NULL,
    PRIMARY KEY(track_id, genre_id),
    FOREIGN KEY(track_id) REFERENCES tracks(id) ON DELETE CASCADE,
    FOREIGN KEY(genre_id) REFERENCES genres(id) ON DELETE CASCADE
);

CREATE TABLE user_track (
    user_id int NOT NULL,
    track_id int NOT NULL,
    PRIMARY KEY(user_id, track_id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(track_id) REFERENCES tracks(id) ON DELETE CASCADE
);

CREATE TABLE session (
    id varchar(64) NOT NULL PRIMARY KEY,
    userID int NOT NULL,
    expire date NOT NULL,
    FOREIGN KEY(userID) REFERENCES users(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION  count_index() RETURNS trigger AS $emp_stamp$
    BEGIN
        new.index := count(*) + 1 from tracks where tracks.album_id = new.album_id;
        Return new;
    END;
$emp_stamp$ LANGUAGE plpgsql;

CREATE TRIGGER emp_stamp BEFORE INSERT ON tracks
    FOR EACH ROW EXECUTE PROCEDURE count_index();
