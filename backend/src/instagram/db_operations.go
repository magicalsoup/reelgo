package instagram

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/magicalsoup/reelgo/.gen/reelgo/public/model"
	. "github.com/magicalsoup/reelgo/.gen/reelgo/public/table"
)


func addVerificationCodeToDB(db *sql.DB, code string, huid string, ig_id string) error {
	verification_code := model.VerificationCodes{
		Huid: huid,
		InstagramID: &ig_id,
		Code: &code,
	}

	stmt := VerificationCodes.INSERT(VerificationCodes.AllColumns).MODEL(verification_code).ON_CONFLICT(VerificationCodes.Huid).DO_UPDATE(
		SET(
			VerificationCodes.Code.SET(String(code)),
		),
	)

	_, err := stmt.Exec(db)

	if err != nil {
		return err
	}

	return nil
}

func getVerificationStatus(db *sql.DB, ig_id string) (bool, error) {
	stmt := Users.SELECT(Users.Verified).WHERE(Users.InstagramID.EQ(String(ig_id)))

	verified := false

	err := stmt.Query(db, verified)

	if err != nil {
		return verified, err
	}

	return verified, nil
}
