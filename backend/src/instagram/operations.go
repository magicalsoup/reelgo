package instagram

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/magicalsoup/reelgo/.gen/reelgo/public/model"
	. "github.com/magicalsoup/reelgo/.gen/reelgo/public/table"
)

func addVerificationCode(db *sql.DB, code string, hashed_id string, igsid string) error {
	stmt := VerificationCodes.INSERT()
	
	return nil
}