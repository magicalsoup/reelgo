
CREATE TABLE users (
    uid SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    salt TEXT NOT NULL,
    instagram_id TEXT UNIQUE,
    verified BOOLEAN NOT NULL
);

CREATE TABLE verification_codes (
    vid SERIAL PRIMARY KEY,
    uid INTEGER UNIQUE NOT NULL,
    instagram_id TEXT NOT NULL,
    code TEXT NOT NULL,
    FOREIGN KEY (uid) REFERENCES "users" (uid)
);

CREATE TABLE tokens (
    id SERIAL PRIMARY KEY,
    bearer_token TEXT UNIQUE NOT NULL,
    expiry_time BIGINT NOT NULL,
    uid INTEGER NOT NULL,
    FOREIGN KEY(uid) REFERENCES "users" (uid)
);

CREATE TABLE trips (
    uid INTEGER NOT NULL,
    tid BIGSERIAL,
    trip_name TEXT NOT NULL,
    PRIMARY KEY(uid, tid),
    FOREIGN KEY(uid) REFERENCES "users" (uid)
);

Create Table attractions (
    uid INTEGER NOT NULL,
    tid INTEGER NOT NULL,
    aid BIGSERIAL,
    attraction_name TEXT NOT NULL,
    attraction_location TEXT NOT NULL,
    PRIMARY KEY(uid, tid, aid),
    FOREIGN KEY(uid, tid) REFERENCES "trips" (uid, tid)
);
