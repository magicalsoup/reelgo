package users

import (
	"database/sql"
	"errors"
	"time"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/magicalsoup/reelgo/.gen/reelgo/public/model"
	. "github.com/magicalsoup/reelgo/.gen/reelgo/public/table"
)

func getUser(db *sql.DB, email string) (*model.Users, error) {
	stmt := Users.SELECT(Users.AllColumns).WHERE(Users.Email.EQ(String(email)))
	var users []model.Users

	err := stmt.Query(db, &users)

	if err != nil {
		return nil, err
	}

	if len(users) != 1 {
		return nil, nil
	}

	return &users[0], nil
}

func createUser(db *sql.DB, name string, email string, password string) (*model.Users, error) {
	salt, hashed_secret := generateHashedPassword(password)
	stmt := Users.INSERT(Users.Name, Users.Email, Users.HashedPassword, Users.Salt, Users.Verified).VALUES(name, email, hashed_secret, salt, false).RETURNING(Users.AllColumns)

	user := &model.Users{}

	err := stmt.Query(db, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func createSessionToken(db *sql.DB, uid int32) (*model.Tokens, error) {
	sessionToken := uuid.New().String()

	// 8 days from now
	expiry_time := time.Now().Unix() + 60*60*24*8

	stmt := Tokens.INSERT(Tokens.BearerToken, Tokens.ExpiryTime, Tokens.UID).VALUES(sessionToken, expiry_time, uid).RETURNING(Tokens.AllColumns)

	token := &model.Tokens{}

	err := stmt.Query(db, token)

	if err != nil {
		return nil, err
	}
	return token, nil
}

func getTokenByUserId(db *sql.DB, uid int32) (*model.Tokens, error) {
	stmt := Tokens.SELECT(Tokens.AllColumns).WHERE(Tokens.UID.EQ(Int32(uid)))

	token := &model.Tokens{}
	err := stmt.Query(db, token)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func refreshSessionToken(db *sql.DB, uid int32) (*model.Tokens, error) {

	oldToken, err := getTokenByUserId(db, uid)

	if err != nil {
		return nil, err
	}

	expiry_time := time.Now().Unix() + 60*60*24*8
	stmt := Tokens.UPDATE(Tokens.ExpiryTime).SET(expiry_time).WHERE(Tokens.ID.EQ(Int32(oldToken.ID))).RETURNING(Tokens.AllColumns)

	newToken := &model.Tokens{}

	err = stmt.Query(db, newToken)

	if err != nil {
		return nil, err
	}
	return newToken, nil
}

func invalidateSessionToken(db *sql.DB, bearer_token string) error {
	stmt := Tokens.UPDATE(Tokens.ExpiryTime).SET(0).WHERE(Tokens.BearerToken.EQ(String(bearer_token)))

	_, err := stmt.Exec(db)

	if err != nil {
		return err
	}

	return err
}

func getUserByToken(db *sql.DB, bearer_token string) (*model.Users, error) {
	get_token_stmt := Tokens.SELECT(Tokens.AllColumns).WHERE(Tokens.BearerToken.EQ(String(bearer_token)))

	token := &model.Tokens{}

	err := get_token_stmt.Query(db, token)

	// should return no rows as an err
	if err != nil {
		return nil, err
	}

	if token.ExpiryTime < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	get_user_stmt := Users.SELECT(Users.AllColumns).WHERE(Users.UID.EQ(Int32(token.UID)))

	user := &model.Users{}

	err = get_user_stmt.Query(db, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
