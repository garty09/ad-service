package ad

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"ad-service/internal/model"
	"ad-service/internal/pagination"
)

type Api struct {
	r *model.Repository
}

func New(db *sql.DB) *Api {
	return &Api{model.NewRepository(db)}
}

func (a *Api) RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/ad", a.list).Methods("GET")
	r.HandleFunc("/ad/{id}", a.get).Methods("GET")
	r.HandleFunc("/ad", a.create).Methods("POST")
}

func (a *Api) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	// the id type from string to int
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BadRequest(err.Error()))
		return
	}

	ad, err := a.r.Get(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(buildErrorResponse(err))
		return
	}

	var (
		viewDescription bool
		viewPhotos bool
	)

	fieldsStr := strings.TrimSpace(r.URL.Query().Get("fields"))
	fieldsStrs := strings.Split(fieldsStr, ",")
	for _, s := range fieldsStrs {
		st := strings.TrimSpace(s)
		if st != "" {
			if st == "desc"{
				viewDescription = true
			}
			if st == "photos"{
				viewPhotos = true
			}
		}
	}
	f := &Full{
		ID:          ad.ID,
		CreatedAt:   ad.CreatedAt,
		Title:       ad.Title,
		Price:       ad.Price,
	}
	if len(ad.PhotoLinks)> 0{
		f.PhotoMain = ad.PhotoLinks[0]
	}
	if viewDescription{
		f.Description = ad.Description
	}
	if viewPhotos{
		f.PhotoLinks = ad.PhotoLinks
	}

	res := &AdResponse{
		Success: true,
		Ad: f,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return
}

func (a *Api) list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	count, err := a.r.Count(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(buildErrorResponse(err))
		return
	}
	page := pagination.ParseInt(r.URL.Query().Get("page"), 1)
	p := pagination.New(page, count)

	var desc bool
	sortStr := strings.TrimSpace(r.URL.Query().Get("sort_by"))
	sortStr = strings.TrimSpace(sortStr)
	if strings.HasPrefix(sortStr, "-") {
		desc = true
		sortStr = sortStr[1:]
	}

	items, err := a.r.List(ctx, p.Offset(), p.PerPage, sortStr, desc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(buildErrorResponse(err))
		return
	}
	var result []Short
	for _, item := range items {
		a := Short{
			ID:          item.ID,
			Title:       item.Title,
			Price:       item.Price,
		}
		if len(item.PhotoLinks)> 0{
			a.PhotoMain = item.PhotoLinks[0]
		}
		result = append(result, a)
	}

	res := AdsResponse{
		Success: true,
		Pagination: &Pagination{
			Limit:          p.PerPage,
			TotalCount:     p.TotalCount,
			CurrentPage:    p.Page,
			TotalPageCount: p.PageCount,
		},
		Ads: result,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return
}

func (a *Api) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	var input CreateAdRequest

	bodyRead, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	err := json.Unmarshal(bodyRead, &input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(buildErrorResponse(err))
		return
	}

	if err := input.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(InvalidInput(err))
		return
	}

	now := time.Now()
	id, err := a.r.Create(ctx, model.AdDB{
		CreatedAt:   now,
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		PhotoLinks:  input.PhotoLinks,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(buildErrorResponse(err))
		return
	}

	res := &AdCreateResponse{
		Success: true,
		Id:      id,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
	return
}
