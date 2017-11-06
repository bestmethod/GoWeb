package rpcListener

import (
	"../multiLogger"
	"database/sql"
	"fmt"
	"github.com/russross/meddler"
	"time"
)

func SessionCleaner(logger *multiLogger.LogHandler, db *sql.DB) {
	Session := new(SessionStruct)
	err := meddler.QueryRow(db, Session, "delete from session where expires < ?", time.Now().Unix())
	if err != nil && err.Error() != "sql: no rows in result set" {
		logger.Error(fmt.Sprintf("Could not cleanup session data in the DB: %s\n", err))
	} else if err != nil {
		logger.Debug("Clean complete")
	} else {
		logger.Debug("Clean complete")
	}
}
