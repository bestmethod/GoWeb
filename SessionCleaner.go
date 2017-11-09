package goweb

import (
	"./multiLogger"
	"database/sql"
	"fmt"
	"github.com/russross/meddler"
	"time"
)

func SessionCleaner(logger *multiLogger.LogHandler, db *sql.DB) {
	Session := new(SessionStruct)
	err := meddler.QueryRow(db, Session, "delete from session where expires < ?", time.Now().Unix())
	if err != nil && err.Error() != "sql: no rows in result set" {
		logger.Error(fmt.Sprintf(LOG_CLEANER_FART, err))
	} else if err != nil {
		logger.Debug(LOG_CLEANER_DONE)
	} else {
		logger.Debug(LOG_CLEANER_DONE)
	}
}
