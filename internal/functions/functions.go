package internal

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	apcfg "github.com/AlexGrtsnk/go-prof-diploma-1/internal/app_config"
	ath "github.com/AlexGrtsnk/go-prof-diploma-1/internal/authentification"

	cks "github.com/AlexGrtsnk/go-prof-diploma-1/internal/cookies"
	db "github.com/AlexGrtsnk/go-prof-diploma-1/internal/db"
	gzp "github.com/AlexGrtsnk/go-prof-diploma-1/internal/gzp"
	flw "github.com/AlexGrtsnk/go-prof-diploma-1/internal/json_parser"
	lg "github.com/AlexGrtsnk/go-prof-diploma-1/internal/logger"

	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// аутентификация доне в теории
func registrateNewUserPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("we are here")
	if r.Method == http.MethodPost {
		reader, err := gzp.GzipFormatHandlerJSON(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		var cookiesTmp *http.Cookie
		_, err = cks.GetCookieHandler(w, r)
		if err != nil {
			token, err := ath.BuildJWTString()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, err = io.WriteString(w, "Error on the side")
				if err != nil {
					log.Fatal(err)
				}
			}
			cookiesTmp = cks.SetCookieHandler(w, r, token)
		} else {
			cookiesTmp, err = r.Cookie("exampleCookie")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, err = io.WriteString(w, "Error on the side")
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		fmt.Println("we are here1")
		http.SetCookie(w, cookiesTmp)
		fmt.Println("we are here2")
		var ath flw.Auth
		var buf bytes.Buffer
		_, err = buf.ReadFrom(reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &ath); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("we are here3")
		flag, _, err := db.DataBaseCheckUserExistance(ath.Login, ath.Password)
		fmt.Println("we are here4", flag)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if flag == 1 {
			w.WriteHeader(http.StatusConflict)
			http.SetCookie(w, nil)
		} else {
			err = db.DataBasePostUser(ath.Login, ath.Password, cookiesTmp.Value)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				_, err = io.WriteString(w, "Error on the side")
				if err != nil {
					log.Fatal(err)
				}
				return
			}
			fmt.Println("every thing went good ", ath.Login, ath.Password)
			w.WriteHeader(http.StatusOK)
		}

	}
}

// v teorii authen done
func authentificateUserPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		reader, err := gzp.GzipFormatHandlerJSON(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		var ath flw.Auth
		var buf bytes.Buffer
		_, err = buf.ReadFrom(reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &ath); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		flag, token, err := db.DataBaseCheckUserExistance(ath.Login, ath.Password)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if flag == 0 {
			w.WriteHeader(http.StatusBadRequest)
			http.SetCookie(w, nil)
		} else {
			cks.SetCookieHandler(w, r, token)
			w.WriteHeader(http.StatusOK)
			fmt.Println("webt good")
		}

	}

}

// dopisat proverky algoritm luna
func uploadNewOrderPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		token, err := cks.GetCookieHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println("want smth", token)
		flag, err := db.DataBaseCheckAuth(token)
		fmt.Println("want smth", flag, err)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if flag == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		fmt.Println("want smth1")
		orderTmp, _ := io.ReadAll(r.Body)
		orderNumber := string(orderTmp)
		fmt.Println("i want to see uploaded number ", orderNumber)
		flag, err = db.DataBaseCheckOrderExistance(orderNumber, token)
		fmt.Println("checking ", flag, err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if flag == 3 {
			w.WriteHeader(http.StatusAccepted)
			fmt.Println("help1")
			err = db.DataBasePostOrder(orderNumber, token)
			fmt.Println("help2")
			if err != nil {
				w.WriteHeader(http.StatusGatewayTimeout)
				_, err = io.WriteString(w, "Error on the side")
				if err != nil {
					log.Fatal(err)
				}
				return
			}
			fmt.Println("help3")
			return
		}

		if flag == 1 {
			w.WriteHeader(http.StatusOK)
			return
		}
		if flag == 2 {
			w.WriteHeader(http.StatusConflict)
			return
		}
		//w.WriteHeader(http.StatusProcessing)
	}
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		token, err := cks.GetCookieHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		flag, err := db.DataBaseCheckAuth(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if flag == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		//var newOrderItems flw.Ord
		fmt.Println("starting processing1")
		newOrdersList, err := db.DataBaseGetOrders(token)
		fmt.Println("starting processing2", newOrdersList[0], err)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(newOrdersList); err != nil {
			log.Panic(err)
		}
	}

}

/*
// v teorii done
func getOrdersNumbersPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		token, err := cks.GetCookieHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		flag, err := db.DataBaseCheckAuth(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if flag == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		//var newOrderItems flw.Ord
		newOrdersList, err := db.DataBaseGetOrders(token)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(newOrdersList); err != nil {
			log.Panic(err)
		}

	}

}
*/
// vrode done
func GetUserBalancePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		token, err := cks.GetCookieHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		flag, err := db.DataBaseCheckAuth(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if flag == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		accural, withdraw, err := db.DataBaseGetUserAccural(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		var answ flw.UserAnsw
		answ.Accural = accural
		answ.Withdraw = withdraw
		resp, err := json.Marshal(answ)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(resp)
		if err != nil {
			log.Fatal(err)
		}

	}
}

func dropBalancePage(w http.ResponseWriter, r *http.Request) {
	reader, err := gzp.GzipFormatHandlerJSON(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = io.WriteString(w, "Error on the side")
		if err != nil {
			log.Fatal(err)
		}
	}
	token, err := cks.GetCookieHandler(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err = io.WriteString(w, "Error on the side")
		if err != nil {
			log.Fatal(err)
		}
	}
	flag, err := db.DataBaseCheckAuth(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err = io.WriteString(w, "Error on the side")
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	if flag == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		_, err = io.WriteString(w, "Error on the side")
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	var balanceBatch flw.BalanceAnsw
	var buf bytes.Buffer
	_, err = buf.ReadFrom(reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &balanceBatch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	flag, err = db.DataBaseCheckOrderExistance(balanceBatch.Order, token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err = io.WriteString(w, "Error on the side")
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	if flag == 2 {
		w.WriteHeader(http.StatusFailedDependency)
		return
	}
	flag, err = db.DataBaseUserGetBalance(token, balanceBatch.Sum)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err = io.WriteString(w, "Error on the side")
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	if flag == 3 {
		w.WriteHeader(http.StatusLengthRequired)
		return
	}
	err = db.DataBaseUserSumBalance(token, balanceBatch.Sum, balanceBatch.Order)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err = io.WriteString(w, "Error on the side")
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)

}

func GetAllOrdersBalanceDropPage(w http.ResponseWriter, r *http.Request) {
	token, err := cks.GetCookieHandler(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err = io.WriteString(w, "Error on the side")
		if err != nil {
			log.Fatal(err)
		}
	}
	flag, err := db.DataBaseCheckAuth(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err = io.WriteString(w, "Error on the side")
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	if flag == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		_, err = io.WriteString(w, "Error on the side")
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	answBatch, err := db.DataBaseOrdersDropBalance(token)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = io.WriteString(w, "Error on the side")
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(answBatch); err != nil {
		log.Panic(err)
	}

}

func GetAllUsersBallsOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		token, err := cks.GetCookieHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
		}
		flag, err := db.DataBaseCheckAuth(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if flag == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		vars := mux.Vars(r)
		orderNumber, ok := vars["id"]
		if !ok {
			fmt.Println("id is missing in parameters")
			w.WriteHeader(http.StatusInternalServerError)
			_, err := io.WriteString(w, "bad request")
			if err != nil {
				log.Fatal(err)
			}
		}

		flag, err = db.DataBaseCheckOrderExistance(orderNumber, token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if flag == 2 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		answ, err := db.DataBaseOrdersBalls(orderNumber, token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		resp, err := json.Marshal(answ)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(resp)
		if err != nil {
			log.Fatal(err)
		}

	}
}

func Run() error {
	var cfg apcfg.Config
	err := env.Parse(&cfg)
	flagRunAddr, fileName, databaseDSN := apcfg.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.RunAddress != "" {
		flagRunAddr = "localhost:8080"
	}
	if cfg.AccrualSystemAddress != "" {
		fileName = cfg.AccrualSystemAddress
	}
	if cfg.DatabaseURI != "" {
		databaseDSN = cfg.DatabaseURI
	}
	log.Println(cfg)
	err = db.DataBaseStartConfig(databaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	if databaseDSN != "localhost" {
		err = db.DataBasePingHandler()
		if err != nil {
			log.Fatal(err)
		}
	}
	err = db.DataBaseCfg(flagRunAddr, fileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("where postgres is hosted:", databaseDSN)
	fmt.Println("where accural system is held", fileName)
	fmt.Println("Running server on", flagRunAddr)
	mux1 := mux.NewRouter()
	mux1.HandleFunc(`/api/user/register`, lg.WithLogging(restrationHandler()))
	mux1.HandleFunc(`/api/user/login`, lg.WithLogging(loginHandler()))
	mux1.HandleFunc(`/api/user/orders`, lg.WithLogging(OrdersHandler()))
	mux1.HandleFunc(`/api/user/balance/`, lg.WithLogging(UserBalanceHandler()))
	mux1.HandleFunc(`/api/user/balance/withdrawals`, lg.WithLogging(UserDropBalanceHandler()))
	mux1.HandleFunc(`/api/user/withdrawals`, lg.WithLogging(UserDroppedBalanceStatsHandler()))
	mux1.HandleFunc(`/api/orders/{id}`, lg.WithLogging(GetAllUsersBallsOrdersHandler()))
	return http.ListenAndServe(flagRunAddr, gzp.GzipHandle(mux1))
}

func restrationHandler() http.Handler {
	fn := registrateNewUserPage
	return http.HandlerFunc(fn)
}

func loginHandler() http.Handler {
	fn := authentificateUserPage
	return http.HandlerFunc(fn)
}

func OrdersHandler() http.Handler {
	fn := uploadNewOrderPage
	return http.HandlerFunc(fn)
}

func UserBalanceHandler() http.Handler {
	fn := GetUserBalancePage
	return http.HandlerFunc(fn)
}

func UserDropBalanceHandler() http.Handler {
	fn := dropBalancePage
	return http.HandlerFunc(fn)
}

func UserDroppedBalanceStatsHandler() http.Handler {
	fn := GetAllOrdersBalanceDropPage
	return http.HandlerFunc(fn)
}

func GetAllUsersBallsOrdersHandler() http.Handler {
	fn := GetAllUsersBallsOrders
	return http.HandlerFunc(fn)
}
