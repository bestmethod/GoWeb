package goweb

import (
	"./configurator"
	"./multiLogger"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/russross/meddler"
	"io/ioutil"
)

func Connect(conn *DbConf, logger *LogHandler) (db *sql.DB, errRet error) {
	var err error
	if conn.Type == "sqlite3" {
		logger.Debug(LOG_SQLITE3_OPEN)
		db, err = sql.Open("sqlite3", conn.Server)
	} else if conn.Type == "MySQL" {
		logger.Debug(LOG_MYSQL_OPEN)
		param := ""
		if conn.UseTLS == true && conn.TLSSkipVerify == false && conn.Ssl_ca == "" {
			logger.Debug(LOG_MYSQL_TLS_TRUE)
			param = "?tls=true"
		} else if conn.UseTLS == true && conn.TLSSkipVerify == true {
			logger.Debug(LOG_MYSQL_TLS_SKIP)
			param = "?tls=skip-verify"
		} else if conn.UseTLS == true && conn.Ssl_ca != "" {
			logger.Debug(LOG_MYSQL_VOODOO_START)
			// this voodoo is explained here: https://godoc.org/github.com/go-sql-driver/mysql#RegisterTLSConfig
			rootCertPool := x509.NewCertPool()
			pem, err := ioutil.ReadFile(conn.Ssl_ca)
			if err != nil {
				logger.Error(fmt.Sprint(err))
				return nil, err
			}
			if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
				logger.Error(LOG_MYSQL_PEMFAIL)
				return nil, errors.New(LOG_MYSQL_PEMFAIL)
			}
			var tt = new(tls.Config)
			tt.RootCAs = rootCertPool
			mysql.RegisterTLSConfig("custom", tt)
			param = "?tls=custom"
			logger.Debug(LOG_MYSQL_VOODOO_DONE)
		}
		logger.Debug(LOG_MYSQL_CONN)
		db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s%s", conn.User, conn.Password, conn.Server, conn.DbName, param))
	} else {
		logger.Error(fmt.Sprintf(LOG_WRONG_DBTYPE, conn.Type))
		return nil, errors.New(fmt.Sprintf(LOG_WRONG_DBTYPE, conn.Type))
	}
	if err != nil {
		logger.Error(fmt.Sprint(err))
		return nil, err
	}
	logger.Debug(LOG_MEDDLER_DEFAULTS)
	if conn.Type == "sqlite3" {
		meddler.Default = meddler.SQLite
	} else if conn.Type == "MySQL" {
		meddler.Default = meddler.MySQL
	}
	return db, nil
}
