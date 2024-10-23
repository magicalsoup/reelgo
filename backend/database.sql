CREATE TABLE users (
    uid INTEGER PRIMARY KEY,
    name TEXT,
    email TEXT,
    hashed_password TEXT,
    instagram_id TEXT,
    verified BOOLEAN
);

CREATE TABLE verification_codes (
    uid INTEGER,
    code TEXT,
    PRIMARY KEY (uid, code),
    FOREIGN KEY (uid) REFERENCES "users" (uid)
);

CREATE TABLE trips (
    uid INTEGER,
    tid INTEGER,
    trip_name TEXT,
    PRIMARY KEY(uid, tid),
    FOREIGN KEY(uid) REFERENCES "users" (uid)
);

Create Table attractions (
    uid INTEGER,
    tid INTEGER,
    aid INTEGER,
    attraction_name TEXT,
    attraction_location TEXT,
    PRIMARY KEY(uid, tid, aid),
    FOREIGN KEY(uid, tid) REFERENCES "trips" (uid, tid)
);
