package logger

import (
	"io"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func WithLogging(h http.Handler) func(w http.ResponseWriter, r *http.Request) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		logEr := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err = io.WriteString(w, "Error on the logger side")
			if err != nil {
				log.Fatal(err)
			}
		}
		return logEr
	}
	sugar := *logger.Sugar()
	logFn := func(w http.ResponseWriter, r *http.Request) {
		// функция Now() возвращает текущее время
		start := time.Now()

		// эндпоинт /ping
		uri := r.RequestURI
		// метод запроса
		method := r.Method

		// точка, где выполняется хендлер pingHandler
		h.ServeHTTP(w, r) // обслуживание оригинального запроса

		// Since возвращает разницу во времени между start
		// и моментом вызова Since. Таким образом можно посчитать
		// время выполнения запроса.
		duration := time.Since(start)

		// отправляем сведения о запросе в zap
		sugar.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
		)

	}
	// возвращаем функционально расширенный хендлер
	return logFn
}
