package configurator

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/russross/meddler"
	"io/ioutil"
	"log"
)

func (conn *DbConf) Connect(earlyDebug bool) (db *sql.DB) {
	var err error
	if conn.Type == "sqlite3" {
		if earlyDebug == true {
			fmt.Println("LOADCONFIG:DBCONNECT: DEBUG sqlite3 connection, opening")
		}
		db, err = sql.Open("sqlite3", conn.Server)
	} else if conn.Type == "MySQL" {
		if earlyDebug == true {
			fmt.Println("LOADCONFIG:DBCONNECT: DEBUG MySQL connection, opening")
		}
		param := ""
		if conn.UseTLS == true && conn.TLSSkipVerify == false && conn.ssl_ca == "" {
			if earlyDebug == true {
				fmt.Println("LOADCONFIG:DBCONNECT: DEBUG MySQL tls true")
			}
			param = "?tls=true"
		} else if conn.UseTLS == true && conn.TLSSkipVerify == true {
			if earlyDebug == true {
				fmt.Println("LOADCONFIG:DBCONNECT: DEBUG MySQL tls skip-verify")
			}
			param = "?tls=skip-verify"
		} else if conn.UseTLS == true && conn.ssl_ca != "" {
			if earlyDebug == true {
				fmt.Println("LOADCONFIG:DBCONNECT: DEBUG MySQL ssl_ca voodoo starting")
			}
			// this voodoo is explained here: https://godoc.org/github.com/go-sql-driver/mysql#RegisterTLSConfig
			rootCertPool := x509.NewCertPool()
			pem, err := ioutil.ReadFile(conn.ssl_ca)
			if err != nil {
				log.Fatal(err)
			}
			if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
				log.Fatal("Failed to append PEM.")
			}
			var tt = new(tls.Config)
			tt.RootCAs = rootCertPool
			mysql.RegisterTLSConfig("custom", tt)
			param = "?tls=custom"
			if earlyDebug == true {
				fmt.Println("LOADCONFIG:DBCONNECT: DEBUG MySQL ssl_ca voodoo done, tls custom")
			}
		}
		if earlyDebug == true {
			fmt.Println("LOADCONFIG:DBCONNECT: DEBUG MySQL opening connection")
		}
		db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s%s", conn.User, conn.Password, conn.Server, conn.DbName, param))
	} else {
		log.Fatalf("ERROR: The only supported database types are: MySQL | sqlite3. Currently set: %s\n", conn.Type)
	}
	if err != nil {
		log.Fatal(err)
	}
	if earlyDebug == true {
		fmt.Println("LOADCONFIG:DBCONNECT: DEBUG configuring meddler defaults")
	}
	if conn.Type == "sqlite3" {
		meddler.Default = meddler.SQLite
	} else if conn.Type == "MySQL" {
		meddler.Default = meddler.MySQL
	}
	return db
}
