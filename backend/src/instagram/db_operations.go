package instagram

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/magicalsoup/reelgo/.gen/reelgo/public/model"
	. "github.com/magicalsoup/reelgo/.gen/reelgo/public/table"
)


func addVerificationCodeToDB(db *sql.DB, code string, hashed_id string, ig_id string) error {
	verification_code := model.VerificationCodes{
		Huid: hashed_id,
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