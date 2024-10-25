package auth

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/magicalsoup/reelgo/.gen/reelgo/public/model"
	. "github.com/magicalsoup/reelgo/.gen/reelgo/public/table"
)


func getVerificationCode(db *sql.DB, huid string) (*model.VerificationCodes, error) {

	stmt := VerificationCodes.SELECT(VerificationCodes.AllColumns).WHERE(VerificationCodes.Huid.EQ(String(huid))).LIMIT(1)

	verification_code := &model.VerificationCodes{}
	
	err := stmt.Query(db, verification_code)

	if err != nil {
		return verification_code, err
	}

	return verification_code, nil
}

func storeIGIDToUser(db *sql.DB, uid int32, ig_id string) error {

	// we avoid verified users
	stmt := Users.UPDATE().SET(
		Users.InstagramID.SET(String(ig_id)),
	).WHERE(Users.UID.EQ(Int32(uid)).AND(
		Users.Verified.NOT_EQ(Bool(false))),
	)

	_, err := stmt.Exec(db)

	if err != nil {
		return err
	}
	return nil
}