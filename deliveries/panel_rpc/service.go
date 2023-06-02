package panel_rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/glatika/providence/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type PanelRPC struct {
	TaskUsecase     model.TaskUsecase
	StockVarUsecase model.StockVariantUsecase
	StockUsecase    model.StockUsecase
}

// TODO: make RPC available to panel

// getting dashboard executive summaries
func (u PanelRPC) RpcDashboardSummaries(w http.ResponseWriter, r *http.Request) {

}

// get list of recently beaconing client
func (u PanelRPC) RpcCurrentOnlineClient(w http.ResponseWriter, r *http.Request) {

}

// get list of tasks of specified client id
func (u PanelRPC) RpcClientTasks(w http.ResponseWriter, r *http.Request) {

}

// get client details
func (u PanelRPC) RpcClientDetails(w http.ResponseWriter, r *http.Request) {

}

// get page of tasks
func (u PanelRPC) RpcTasksPage(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(
		chi.URLParam(r, "page"),
	)
	if err != nil {
		w.WriteHeader(422)
		return
	}
	size, err := strconv.Atoi(
		chi.URLParam(r, "size"),
	)
	if err != nil {
		w.WriteHeader(422)
		return
	}

	tasks, err := u.TaskUsecase.GetAllTask(page, size)
	if err != nil {
		w.WriteHeader(502)
		return
	}

	jsonContent, err := json.Marshal(*tasks)
	if err != nil {
		w.WriteHeader(502)
		return
	}

	// setup header
	w.Header().Set("Content-Type", "application/json")
	// setup body
	w.Write(jsonContent)
	// write header will flush the http setup
	// make user put it at last
}

// get page of clients
func (u PanelRPC) RpcClientPage(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(
		chi.URLParam(r, "page"),
	)
	if err != nil {
		w.WriteHeader(422)
		return
	}
	size, err := strconv.Atoi(
		chi.URLParam(r, "size"),
	)
	if err != nil {
		w.WriteHeader(422)
		return
	}

	stocks, err := u.StockUsecase.GetAllStocks(page, size)
	if err != nil {
		w.WriteHeader(502)
		return
	}

	jsonContent, err := json.Marshal(stocks)
	if err != nil {
		w.WriteHeader(502)
		return
	}

	// setup header
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonContent)
	// setup body
}

func (u PanelRPC) Run(portNumber int) error {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Providence Pannel"))
	})
	r.Group(func(rpc chi.Router) {

		rpc.Get("/rpc/authenticate", func(w http.ResponseWriter, r *http.Request) {})
		rpc.Get("/rpc/dashboard-summaries", func(w http.ResponseWriter, r *http.Request) {})
		rpc.Get("/rpc/current-online-client", func(w http.ResponseWriter, r *http.Request) {})
		rpc.Get("/rpc/client-tasks", func(w http.ResponseWriter, r *http.Request) {})
		rpc.Get("/rpc/client-details", func(w http.ResponseWriter, r *http.Request) {})
		rpc.Get("/rpc/tasks-page-{page}-{size}", u.RpcTasksPage)
		rpc.Get("/rpc/client-page-{page}-{size}", u.RpcClientPage)
	})
	fmt.Printf("Pannel runnig on %d \n", *&portNumber)

	return http.ListenAndServe(":"+strconv.FormatInt(int64(portNumber), 10), r)
}
