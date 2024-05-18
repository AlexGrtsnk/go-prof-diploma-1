package databaseshortener

import (
	"database/sql"
	"fmt"
	"log"

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

func DataBaseFileNameSelect() (flnm string, err error) {
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
	fmt.Println("inside1")
	if err := db.QueryRow("SELECT token FROM users WHERE lgn = '" + string(login) + "' and psw = '" + string(password) + "';").Scan(&token); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("inside2")
			return 0, "", nil
		}
		fmt.Println("inside3")
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
	fmt.Println("inside11")
	quer := `INSERT INTO users(lgn, psw, token) VALUES ('` + string(login) + `', '` + string(password) + `', '` + token + `')`
	fmt.Println("inside12")
	_, err = db.Exec(quer)
	fmt.Println("inside13 ", err)
	if err != nil {
		return err
	}
	fmt.Println("inside14")
	return nil
}

func DataBaseCheckAuth(token string) (falg int, err error) {
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
	flag = 1
	flag_tmp := 0
	fmt.Println("Before1")
	if err := db.QueryRow("SELECT sts FROM orders WHERE nmb = '" + string(orderNumber) + "' and token = '" + string(token) + "';").Scan(&sts); err != nil {
		if err == sql.ErrNoRows {
			flag = 2
		} else {
			return 0, err
		}
	}
	fmt.Println("Before2", sts)
	if err := db.QueryRow("SELECT sts FROM orders WHERE nmb = '" + string(orderNumber) + "';").Scan(&sts); err != nil {
		if err == sql.ErrNoRows {
			flag_tmp = 3
			if flag == 2 {
				return 3, nil
			}
		} else {
			return 0, err
		}
	}
	if flag_tmp != 0 {
		return 1, err
	} else {
		return 2, err
	}
	//return flag, nil

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
	quer := `INSERT INTO orders(nmb, sts, token) VALUES ('` + string(orderNumber) + `', '` + string("Prinyat") + `', '` + token + `')`
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	return nil
}

func DataBaseGetOrders(token string) (answb flw.OrdersList, err error) {
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
		answ := new(flw.Orders)
		err = rows.Scan(&answ.Number, &answ.Status, &answ.Accrual, &answ.Uoloaded)
		if err != nil {
			return nil, err
		}
		if rows.Err() != nil {
			return nil, rows.Err()
		}
		//answ.ShortURL = apiRunAddr + "/" + answ.ShortURL
		answb = append(answb, *answ)
		flag = 1
	}
	if flag == 0 {
		return nil, sql.ErrNoRows
	}
	return answb, nil
}

func DataBaseGetUserAccural(token string) (ac int, wd int, err error) {
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
	var accuralSum int
	if err := db.QueryRow("SELECT balance FROM users WHERE token = '" + string(token) + "';").Scan(&accuralSum); err != nil {
		if err == sql.ErrNoRows {
			return -1, -1, err
		}
		return -1, -1, err
	}
	var witdhraw int
	if err := db.QueryRow("SELECT wtdh FROM users WHERE token = '" + string(token) + "';").Scan(&witdhraw); err != nil {
		if err == sql.ErrNoRows {
			return -1, -1, err
		}
		return -1, -1, err
	}

	return accuralSum, witdhraw, nil
}

func DataBaseUserSumBalance(token string, balls int, ordernum string) (err error) {
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
	var wtdh int
	if err := db.QueryRow("SELECT wtdh FROM users WHERE token = '" + string(token) + "';").Scan(&wtdh); err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return err
	}
	var balance int
	if err := db.QueryRow("SELECT balance FROM users WHERE token = '" + string(token) + "';").Scan(&balance); err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return err
	}
	ballSum := wtdh + balls
	newBalance := balance - balls
	quer := "UPDATE users SET wtdh='" + fmt.Sprint(ballSum) + "' WHERE token = '" + token + "';"
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	quer = "UPDATE users SET balance='" + fmt.Sprint(newBalance) + "' WHERE token = '" + token + "';"
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	quer = "UPDATE orders SET sumbals='" + fmt.Sprint(balls) + "' WHERE nmb = '" + ordernum + "';"
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}

	return nil
}

func DataBaseUserGetBalance(token string, balance int) (flag int, err error) {
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
	var blnc int
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

/*
func DataBaseUpdateUserBalance(token string) (err error) {
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
	quer := "SELECT accural from orders where token = '" + token + "' and status = 'PROCESSED';"
	rows, err := db.Query(quer)
	if err != nil {
		return err
	}
	//flag := 0
	accuralSum := 0
	for rows.Next() {
		//answ := new(flw.Orders)
		var tmp int
		err = rows.Scan(&tmp)
		if err != nil {
			return err
		}
		if rows.Err() != nil {
			return rows.Err()
		}
		//answ.ShortURL = apiRunAddr + "/" + answ.ShortURL
		//answb = append(answb, *answ)
		if tmp == -1 {
			continue
		}
		accuralSum += tmp
	}

}
*/

func DataBaseOrdersAllBalance(token string) (err error) {
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return err
	}
	db, err := sql.Open(dbms, dbName)
	if err != nil {
		return err
	}
	defer db.Close()
	quer := "SELECT accural, nmb from orders where token = '" + token + "' and status = 'PROCESSED' ;"
	rows, err := db.Query(quer)
	if err != nil {
		return err
	}
	flag := 0
	accuralSum := 0
	for rows.Next() {
		//answ := new(flw.Orders)
		var tmp int
		var nmb string
		err = rows.Scan(&tmp, &nmb)
		if err != nil {
			return err
		}
		if rows.Err() != nil {
			return rows.Err()
		}
		//answ.ShortURL = apiRunAddr + "/" + answ.ShortURL
		//answb = append(answb, *answ)
		if tmp == -1 {
			continue
		}
		flag = 1
		accuralSum += tmp
		quer = "UPDATE orders SET status='COMPLETED' WHERE nmb = '" + nmb + "';"
		_, err = db.Exec(quer)
		if err != nil {
			return err
		}
	}
	if flag == 0 {
		return sql.ErrNoRows
	}
	quer = "UPDATE users SET balance='" + fmt.Sprint(accuralSum) + "' WHERE token = '" + token + "';"
	_, err = db.Exec(quer)
	if err != nil {
		return err
	}
	return nil
}

func DataBaseOrdersDropBalance(token string) (answ flw.DrawAnswList, err error) {
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
	accuralSum := 0

	for rows.Next() {
		var sumbals int
		var nmb string
		var ts string
		var tmp flw.DrawAnsw
		err = rows.Scan(&nmb, &sumbals, &ts)
		if err != nil {
			return nil, err
		}
		if rows.Err() != nil {
			return nil, rows.Err()
		}
		if sumbals > 0 {
			accuralSum += sumbals
			tmp.Number = nmb
			tmp.Sum = sumbals
			tmp.ProccessedAt = ts
			answ = append(answ, tmp)
		}
	}
	if accuralSum == 0 {
		return nil, sql.ErrNoRows
	}
	return answ, nil
}

func DataBaseOrdersBalls(orderNumber string, token string) (answ flw.WithAnsw, err error) {
	var db *sql.DB
	dbName, dbms, err := DataBaseSelfConfigGet()
	if err != nil {
		return answ, err
	}
	db, err = sql.Open(dbms, dbName)
	if err != nil {
		return answ, err
	}
	defer db.Close()
	quer := "SELECT nmb, sts, accural from orders where token = '" + token + "';"
	rows, err := db.Query(quer)
	if err != nil {
		return answ, err
	}
	for rows.Next() {
		err = rows.Scan(&answ.Order, &answ.Status, &answ.Accrual)
		if err != nil {
			return answ, err
		}
		if rows.Err() != nil {
			return answ, rows.Err()
		}
	}
	return answ, nil

}
