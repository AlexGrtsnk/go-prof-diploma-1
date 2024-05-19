package databaseshortener

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	bn "github.com/AlexGrtsnk/go-prof-diploma-1/internal/bindata"

	flw "github.com/AlexGrtsnk/go-prof-diploma-1/internal/json_parser"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/golang-migrate/migrate/database/sqlite3"
	bindata "github.com/golang-migrate/migrate/source/go_bindata"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

const drriver = "sqlite3"
const dbbName = "shortenerdbs.db"

func NewDB() (*sql.DB, error) {
	dbname, driverTemp, err := DataBaseSelfConfigGet()
	if err != nil {
		return nil, err
	}
	sqliteDB, err := sql.Open(driverTemp, dbname)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open sqlite DB")
	}

	return sqliteDB, nil
}

func RunMigrateScripts(db *sql.DB) error {
	var driver database.Driver
	var err error
	dbNameTemp, _, err := DataBaseSelfConfigGet()
	if err != nil {
		return err
	}
	if dbNameTemp == dbbName {
		driver, err = sqlite3.WithInstance(db, &sqlite3.Config{})
	} else {
		driver, err = postgres.WithInstance(db, &postgres.Config{})
	}
	if err != nil {
		return fmt.Errorf("creating db driver failed %s", err)
	}

	res := bindata.Resource(bn.AssetNames(),
		func(name string) ([]byte, error) {
			return bn.Asset(name)
		})

	d, _ := bindata.WithInstance(res)
	m, err := migrate.NewWithInstance("go-bindata", d, dbNameTemp, driver)
	if err != nil {
		return fmt.Errorf("initializing db migration failed %s", err)
	}
	if dbNameTemp == dbbName {
		_ = m.Steps(-1)
		err = m.Steps(1)
		if err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migrating database failed %s", err)
		}
	} else {
		_ = m.Down()
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migrating database failed %s", err)
		}
	}
	return nil
}

func DataBaseCfg(flagRunAddr string, fileName string) (err error) {
	db, err := NewDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	err = RunMigrateScripts(db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	quer := `INSERT INTO cfg(flagRunAddr, flnm) VALUES ('` + string(flagRunAddr) + `', '` + fileName + `')`
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	return nil
}

func DataBasePingHandler() (err error) {
	_, driverTemp, err := DataBaseSelfConfigGet()
	if err != nil {
		return err
	}
	dbName := fmt.Sprintf("host=%s port=%s  user=%s password=%s dbname=%s sslmode=disable",
		`postgres`, `5432`, `postgres`, `postgres`, `praktikum`)
	err = DataBasePing(dbName, driverTemp)
	if err != nil {
		err = DataBaseSelfConfigUpdate(dbbName, drriver)
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

func DataBasePing(dbbname string, driver string) (err error) {
	var db *sql.DB
	var res string
	db, err = sql.Open(driver, dbbname)
	if err != nil {
		return err
	}
	defer db.Close()
	quer := "SELECT 1;"
	rows, err := db.Query(quer)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return rows.Err()
	}
	rows.Next()
	err = rows.Scan(&res)
	if err != nil {
		return err
	}
	return nil
}

func DataBaseAPIAddressSelect() (apiAddress string, err error) {
	var db *sql.DB
	var apiRunAddr string
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return "", err
	}

	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return "", err
	}
	defer db.Close()
	quer := "SELECT flnm FROM cfg WHERE id = 1;"
	rows, err := db.Query(quer)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return "", rows.Err()
	}
	rows.Next()
	err = rows.Scan(&apiRunAddr)
	if err != nil {
		return "", err
	}
	return apiRunAddr, nil
}

func DataBaseStartConfig(dbName string) (err error) {
	var db *sql.DB
	db, err = sql.Open("sqlite3", "cfg.db")
	if err != nil {
		return err
	}
	defer db.Close()
	var driver string
	if dbName != "localhost" {
		driver = "pgx"
	} else {
		driver = "sqlite3"
		dbName = dbbName
	}
	sts1 := `
	DROP TABLE IF EXISTS cfg;
	CREATE TABLE cfg (id INTEGER PRIMARY KEY, dbbname TEXT, driver TEXT);
	INSERT INTO cfg(dbbname, driver) VALUES ('` + string(dbName) + `', '` + string(driver) + `');`
	_, err = db.Exec(sts1)

	if err != nil {
		return err
	}
	return nil
}

func DataBaseSelfConfigGet() (dbbname string, driver string, err error) {
	var db *sql.DB
	db, err = sql.Open("sqlite3", "cfg.db")
	if err != nil {
		return "", "", err
	}
	quer := "SELECT dbbname FROM cfg WHERE id = 1;"
	rows, err := db.Query(quer)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return "", "", rows.Err()
	}
	rows.Next()
	var dbNameTemp string
	err = rows.Scan(&dbNameTemp)
	if err != nil {
		return "", "", err
	}
	quer = "SELECT driver FROM cfg WHERE id = 1;"
	rows, err = db.Query(quer)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return "", "", rows.Err()
	}
	rows.Next()
	var driverTemp string
	err = rows.Scan(&driverTemp)
	if err != nil {
		return "", "", err
	}
	return dbNameTemp, driverTemp, nil
}
func DataBaseSelfConfigUpdate(dbbname string, driver string) (err error) {
	var db *sql.DB
	db, err = sql.Open("sqlite3", "cfg.db")
	if err != nil {
		return err
	}
	defer db.Close()
	quer := "UPDATE cfg SET dbbname='" + dbbname + "', '" + "driver='" + driver + "' WHERE id=1;"
	_, err = db.Exec(quer)

	if err != nil {
		return err
	}
	return nil
}

func DataBaseCheckUserExistance(login string, password string) (flag int, tknm string, err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return 0, "", err
	}

	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return 0, "", err
	}
	defer db.Close()
	var token string
	if err := db.QueryRow("SELECT token FROM users WHERE lgn = '" + string(login) + "' and psw = '" + string(password) + "';").Scan(&token); err != nil {
		if err == sql.ErrNoRows {
			return 0, "", nil
		}
		return 0, "", err
	}
	return 1, token, nil
}

func DataBasePostUser(login string, password string, token string) (err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return err
	}
	defer db.Close()
	quer := `INSERT INTO users(lgn, psw, token) VALUES ('` + string(login) + `', '` + string(password) + `', '` + token + `')`
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	return nil
}

func DataBaseCheckAuth(token string) (flag int, err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return 0, err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return 0, err
	}
	defer db.Close()
	var login string
	if err := db.QueryRow("SELECT lgn FROM users WHERE token = '" + string(token) + "';").Scan(&login); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return 1, nil

}

func DataBaseCheckOrderExistance(orderNumber string, token string) (flag int, err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return 0, err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return 0, err
	}
	defer db.Close()
	var sts string
	flagEQ := 0
	flagNEQ := 0
	if err := db.QueryRow("SELECT sts FROM orders WHERE nmb = '" + string(orderNumber) + "' and token = '" + string(token) + "';").Scan(&sts); err != nil {
		if err == sql.ErrNoRows {
			flagEQ = 1
		} else {
			return 0, err
		}
	}
	if err := db.QueryRow("SELECT sts FROM orders WHERE nmb = '" + string(orderNumber) + "' and token != '" + string(token) + "';").Scan(&sts); err != nil {
		if err == sql.ErrNoRows {
			flagNEQ = 1
		} else {
			return 0, err
		}
	}
	if flagEQ == 1 && flagNEQ == 1 {
		return 3, nil
	}
	if flagEQ == 1 && flagNEQ == 0 {
		return 2, nil
	}
	if flagEQ == 0 && flagNEQ == 1 {
		return 1, nil
	}
	if flagEQ == 0 && flagNEQ == 0 {
		return 2, nil
	}
	return flag, nil

}

func DataBasePostOrder(orderNumber string, token string) (err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return err
	}
	defer db.Close()
	quer := `INSERT INTO orders(nmb, sts, token) VALUES ('` + string(orderNumber) + `', '` + string("PROCESSING") + `', '` + token + `')`
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	return nil
}

func DataBaseGetOrders(token string) (ordersList flw.OrdersList, err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return nil, err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	quer := "SELECT nmb, sts, accural, ts from orders where token = '" + token + "';"
	rows, err := db.Query(quer)
	if err != nil {
		return nil, err
	}
	flag := 0
	for rows.Next() {
		orders := new(flw.Orders)
		err = rows.Scan(&orders.Number, &orders.Status, &orders.Accrual, &orders.Uoloaded)
		if err != nil {
			return nil, err
		}
		if rows.Err() != nil {
			return nil, rows.Err()
		}
		ordersList = append(ordersList, *orders)
		flag = 1
	}
	if flag == 0 {
		return nil, sql.ErrNoRows
	}
	return ordersList, nil
}

func DataBaseGetUserAccural(token string) (ac float64, wd float64, err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return -1, -1, err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return -1, -1, err
	}
	defer db.Close()
	var accuralSum float64
	if err := db.QueryRow("SELECT balance FROM users WHERE token = '" + string(token) + "';").Scan(&accuralSum); err != nil {
		if err == sql.ErrNoRows {
			return -1, -1, err
		}
		return -1, -1, err
	}
	var witdhraw float64
	if err := db.QueryRow("SELECT wtdh FROM users WHERE token = '" + string(token) + "';").Scan(&witdhraw); err != nil {
		if err == sql.ErrNoRows {
			return -1, -1, err
		}
		return -1, -1, err
	}
	return accuralSum, witdhraw, nil
}

func DataBaseUserSumBalance(token string, points float64, ordernum string) (err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return err
	}
	defer db.Close()
	var wtdh float64
	if err := db.QueryRow("SELECT wtdh FROM users WHERE token = '" + string(token) + "';").Scan(&wtdh); err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return err
	}
	var balance float64
	if err := db.QueryRow("SELECT balance FROM users WHERE token = '" + string(token) + "';").Scan(&balance); err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return err
	}
	pointSum := wtdh + points
	newBalance := balance - points
	quer := "UPDATE users SET wtdh='" + fmt.Sprint(pointSum) + "' WHERE token = '" + token + "';"
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	quer = "UPDATE users SET balance='" + fmt.Sprint(newBalance) + "' WHERE token = '" + token + "';"
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}

	quer = `INSERT INTO orders(nmb, sts, token) VALUES ('` + string(ordernum) + `', '` + string("PROCESSED") + `', '` + token + `')`
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	quer = "UPDATE orders SET sumbals='" + fmt.Sprint(points) + "' WHERE token = '" + token + "';"

	_, err = db.Exec(quer)
	if err != nil {
		return err
	}

	return nil
}

func DataBaseUserGetBalance(token string, balance float64) (flag int, err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return 0, err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return 0, err
	}
	defer db.Close()
	var blnc float64
	if err := db.QueryRow("SELECT balance FROM users WHERE token = '" + string(token) + "';").Scan(&blnc); err != nil {
		if err == sql.ErrNoRows {
			return 0, sql.ErrNoRows
		}
		return 0, err
	}
	if blnc-balance < 0 {
		return 3, nil
	}
	return 0, nil
}

func DataBaseOrdersAllBalance(token string, sts string, balance float64, nmb string) (err error) {
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return err
	}
	db, err := sql.Open(dbms, dbName)
	if err != nil {
		return err
	}
	defer db.Close()
	quer := "UPDATE orders SET sts='" + sts + "' WHERE token = '" + token + "';"
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	quer = "UPDATE orders SET accural='" + strconv.FormatFloat(balance, 'f', 3, 64) + "' WHERE token = '" + token + "';"
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	var blncOld float64
	if err := db.QueryRow("SELECT balance FROM users WHERE token = '" + string(token) + "';").Scan(&blncOld); err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return err
	}
	quer = "UPDATE users SET balance='" + strconv.FormatFloat(balance, 'f', 3, 64) + "' WHERE token = '" + token + "';"
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	return nil
}

func DataBaseOrdersDropBalance(token string) (drawAnswList flw.DrawAnswList, err error) {
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open(dbms, dbName)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	quer := "SELECT nmb, sumbals, ts from orders where token = '" + token + "';"
	rows, err := db.Query(quer)
	if err != nil {
		return nil, err
	}
	accuralSum := 0.0

	for rows.Next() {
		var sumPoints float64
		var nmb string
		var ts string
		var drawAnsw flw.DrawAnsw
		err = rows.Scan(&nmb, &sumPoints, &ts)
		if err != nil {
			return nil, err
		}
		if rows.Err() != nil {
			return nil, rows.Err()
		}
		if sumPoints > 0 {
			accuralSum += sumPoints
			drawAnsw.Order = nmb
			drawAnsw.Sum = sumPoints
			drawAnsw.ProccessedAt = ts
			drawAnswList = append(drawAnswList, drawAnsw)
		}
	}
	if accuralSum == 0 {
		return nil, sql.ErrNoRows
	}
	return drawAnswList, nil
}

func DataBaseOrdersPoints(orderNumber string, token string) (withAnsw flw.WithAnsw, err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return withAnsw, err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return withAnsw, err
	}
	defer db.Close()
	quer := "SELECT nmb, sts, accural from orders where token = '" + token + "';"
	rows, err := db.Query(quer)
	if err != nil {
		return withAnsw, err
	}
	for rows.Next() {
		err = rows.Scan(&withAnsw.Order, &withAnsw.Status, &withAnsw.Accrual)
		if err != nil {
			return withAnsw, err
		}
		if rows.Err() != nil {
			return withAnsw, rows.Err()
		}
	}
	return withAnsw, nil

}
