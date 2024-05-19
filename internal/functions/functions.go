package internal

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

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

func Valid(number int) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}

func registrateNewUserPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		reader, err := gzp.GzipFormatHandlerJSON(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		var tokenTmp string
		_, err = cks.GetCookieHandler(w, r)
		if err != nil {
			token, err := ath.BuildJWTString()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, err = io.WriteString(w, "Error on the side")
				if err != nil {
					log.Fatal(err)
				}
				return
			}
			tokenTmp = token
		}
		var ath flw.Auth
		var buf bytes.Buffer
		_, err = buf.ReadFrom(reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &ath); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		flag, _, err := db.DataBaseCheckUserExistance(ath.Login, ath.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if flag == 1 {
			w.WriteHeader(http.StatusConflict)
		} else {
			err = db.DataBasePostUser(ath.Login, ath.Password, tokenTmp)
			qwe := cks.SetCookieHandler(w, r, tokenTmp)
			http.SetCookie(w, qwe)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, err = io.WriteString(w, "Error on the side")
				if err != nil {
					log.Fatal(err)
				}
				return
			}
			w.WriteHeader(http.StatusOK)
		}

	}
}

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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &ath); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		flag, token, err := db.DataBaseCheckUserExistance(ath.Login, ath.Password)
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
			return
		}
		cks.SetCookieHandler(w, r, token)
		w.WriteHeader(http.StatusOK)
	}

}

func uploadNewOrderPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		token, err := cks.GetCookieHandler(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
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
		orderTmp, _ := io.ReadAll(r.Body)
		orderNumber := string(orderTmp)
		tmpOrderNum, err := strconv.Atoi(orderNumber)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		OrderIsValid := Valid(tmpOrderNum)
		if !OrderIsValid {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
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
		if flag == 3 {
			w.WriteHeader(http.StatusAccepted)
			err = db.DataBasePostOrder(orderNumber, token)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, err = io.WriteString(w, "Error on the side")
				if err != nil {
					log.Fatal(err)
				}
				return
			}
			APIAddres, err := db.DataBaseAPIAddressSelect()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, err = io.WriteString(w, "Error on the side")
				if err != nil {
					log.Fatal(err)
				}
				return
			}
			kol := 0
			dataTmp := new(flw.WithAnsw)
			for kol < 60 {
				kol += 1
				res, err := http.Get(APIAddres + "/api/orders/" + orderNumber)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_, err = io.WriteString(w, "Error on the side")
					if err != nil {
						log.Fatal(err)
					}
					return
				}
				defer res.Body.Close()
				data := new(flw.WithAnsw)
				err = json.NewDecoder(res.Body).Decode(data)
				if err != nil {
					continue
				} else {
					dataTmp.Accrual = data.Accrual
					dataTmp.Order = data.Order
					dataTmp.Status = data.Status
					break
				}
			}
			if kol < 60 {
				err = db.DataBaseOrdersAllBalance(token, dataTmp.Status, dataTmp.Accrual, dataTmp.Order)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_, err = io.WriteString(w, "Error on the side")
					if err != nil {
						log.Fatal(err)
					}
					return
				}
			}

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
	}
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
		newOrdersList, err := db.DataBaseGetOrders(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(newOrdersList); err != nil {
			log.Panic(err)
		}
	}

}

func getUserBalancePage(w http.ResponseWriter, r *http.Request) {
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
		accural, withdraw, err := db.DataBaseGetUserAccural(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, "Error on the side")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		var userAnsw flw.UserAnsw
		userAnsw.Current = accural
		userAnsw.Withdraw = withdraw
		resp, err := json.Marshal(userAnsw)
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
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusInternalServerError)
		_, err = io.WriteString(w, "Error on the side")
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)

}

func getAllOrdersBalanceDropPage(w http.ResponseWriter, r *http.Request) {
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
	drawAnswList, err := db.DataBaseOrdersDropBalance(token)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	var tmpDrawAnswList flw.DrawAnswList
	tmpDrawAnswList = append(tmpDrawAnswList, drawAnswList[len(drawAnswList)-1])
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(tmpDrawAnswList); err != nil {
		log.Panic(err)
	}

}

func getAllUsersPointsOrders(w http.ResponseWriter, r *http.Request) {
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
		if flag == 3 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		answ, err := db.DataBaseOrdersPoints(orderNumber, token)
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
	mux1.HandleFunc(`/api/user/balance`, lg.WithLogging(UserBalanceHandler()))
	mux1.HandleFunc(`/api/user/balance/withdraw`, lg.WithLogging(UserDropBalanceHandler()))
	mux1.HandleFunc(`/api/user/withdrawals`, lg.WithLogging(UserDroppedBalanceStatsHandler()))
	mux1.HandleFunc(`/api/orders/{id}`, lg.WithLogging(GetAllUsersPointsOrdersHandler()))
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
	fn := getUserBalancePage
	return http.HandlerFunc(fn)
}

func UserDropBalanceHandler() http.Handler {
	fn := dropBalancePage
	return http.HandlerFunc(fn)
}

func UserDroppedBalanceStatsHandler() http.Handler {
	fn := getAllOrdersBalanceDropPage
	return http.HandlerFunc(fn)
}

func GetAllUsersPointsOrdersHandler() http.Handler {
	fn := getAllUsersPointsOrders
	return http.HandlerFunc(fn)
}
