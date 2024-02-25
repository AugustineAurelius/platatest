package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/platatest/internal/config"
	"github.com/platatest/internal/repository"
)

type CurrencyHandlers interface {
	findCurrencyHandler(w http.ResponseWriter, r *http.Request)
	newCurrencyHandler(w http.ResponseWriter, r *http.Request)
	healthHandler(w http.ResponseWriter, r *http.Request)
	SetServer(server *http.Server)
	Server() *http.Server
	SetAvailableCurrencies(currencies map[string]string)
	IsAvailableCurrency(currency string) bool
	SetRequestURL(url string)
	RequestURL() string
	Close()
	UpdateCourses()
}

type currencyServiceHandler struct {
	dbhandler           repository.DatabaseHandler
	server              *http.Server
	availableCurrencies map[string]string
	requestURL          string
}

func Serve(config config.Config, db repository.DatabaseHandler) {

	handler := newCurrencyHandler(db)
	r := mux.NewRouter()
	curryncyRouter := r.PathPrefix("/currency").Subrouter()
	curryncyRouter.Methods("GET").Path("/{SearchCriteria}/{search:.*}").HandlerFunc(handler.findCurrencyHandler)
	curryncyRouter.Methods("POST").Path("/create").HandlerFunc(handler.newCurrencyHandler)
	r.Methods("GET").Path("/health").HandlerFunc(handler.healthHandler)

	handler.SetRequestURL(config.URL)

	go handler.UpdateCourses()
	handler.SetServer(&http.Server{
		Addr:         config.Address(),
		WriteTimeout: config.WriteTimeout(),
		ReadTimeout:  config.ReadTimeout(),
		IdleTimeout:  config.IdleTimeout(),
		Handler:      r,
	})

	resp, err := http.Get(config.AllCurrenciesURL)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	var mapOfAllCurrencies map[string]string
	err = json.NewDecoder(resp.Body).Decode(&mapOfAllCurrencies)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	handler.SetAvailableCurrencies(mapOfAllCurrencies)

	go func() {
		log.Println("start server")
		if err := handler.Server().ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	handler.Close()
}

func newCurrencyHandler(database repository.DatabaseHandler) CurrencyHandlers {
	return &currencyServiceHandler{dbhandler: database}
}

func (h *currencyServiceHandler) findCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	criteria, ok := vars["SearchCriteria"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{error: %s}`, CriteriaErr)
		return
	}

	searchkey, ok := vars["search"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{error: %s}`, CriteriaErr)
		return
	}

	var price repository.Price
	var err error

	switch strings.ToLower(criteria) {
	case "code":
		price, err = h.dbhandler.GetByName(searchkey)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "{error: %s}", IdErr)
			return
		}
	case "id":
		id, err := strconv.Atoi(searchkey)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "{error: %s}", UnmarshalIdErr)
			return
		}
		price, err = h.dbhandler.GetById(id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "{error: %s}", IdErr)
			return
		}
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "{error: %s}", CriteriaErr)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	err = json.NewEncoder(w).Encode(&price)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *currencyServiceHandler) newCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	var currency repository.Currency

	err := json.NewDecoder(r.Body).Decode(&currency)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{error : %s}`, CreateErr)
		return
	}

	if !h.IsAvailableCurrency(currency.Name) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{error : %s}`, NotSupportedCurrencyErr)
		return
	}

	id, err := h.dbhandler.Create(currency.Name)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{error : %s}`, CreateErr)
		return
	}
	w.Header().Set("Content-Type", "application/json:charset=utf8")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(repository.Currency{Id: id})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *currencyServiceHandler) healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("got health request")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json:charset=utf8")
	err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (h *currencyServiceHandler) UpdateCourses() {

	for range time.NewTicker(time.Second * 15).C {

		currecies, err := h.dbhandler.Fetch()
		if err != nil {
			log.Println(err)
			continue
		}
		for _, cur := range currecies {

			requestURL := fmt.Sprintf(h.requestURL, strings.ToLower(cur.Name))
			log.Println("send request to", requestURL)
			res, err := http.Get(requestURL)
			if err != nil {
				log.Println("error : ", err)
				return
			}
			var data map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&data)
			if err != nil {
				log.Println("error : ", err)
				return
			}
			if err = h.dbhandler.Update(
				data[strings.
					ToLower(strings.
						Split(cur.Name, "/")[1])].(float64),
				cur.Id); err != nil {
				log.Println(err)
				return
			}

		}
	}
}

func (h *currencyServiceHandler) SetAvailableCurrencies(currencies map[string]string) {
	h.availableCurrencies = currencies
}

func (h *currencyServiceHandler) IsAvailableCurrency(currency string) bool {
	if !strings.Contains(currency, "/") {
		return false
	}

	arr := strings.Split(currency, "/")
	if len(arr) != 2 {
		return false
	}
	cur1 := strings.ToLower(arr[0])
	cur2 := strings.ToLower(arr[1])

	_, ok1 := h.availableCurrencies[cur1]
	_, ok2 := h.availableCurrencies[cur2]
	return ok1 && ok2
}

func (h *currencyServiceHandler) SetServer(server *http.Server) {
	h.server = server
}

func (h *currencyServiceHandler) Server() *http.Server {
	return h.server
}

func (h *currencyServiceHandler) SetRequestURL(url string) {
	h.requestURL = url
}

func (h *currencyServiceHandler) RequestURL() string {
	return h.requestURL
}

func (h *currencyServiceHandler) Close() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	<-sigChannel

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	h.server.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
