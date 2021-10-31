package foodservice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-pg/pg/v10"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type api struct {
	db *pg.DB
}

func Run(cfg *Config) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("error validating config: %w", err)
	}

	db := pg.Connect(&pg.Options{
		Addr:     cfg.DatabaseAddr,
		User:     cfg.DatabaseUser,
		Password: cfg.DatabasePassword,
		Database: cfg.DatabaseName,
	})

	if err := db.Ping(context.Background()); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	if err := createSchema(db); err != nil {
		return fmt.Errorf("couldn't create schema: %w", err)
	}

	a := &api{
		db: db,
	}

	router := chi.NewRouter()

	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/healthz", a.healthz)

	router.Get("/recipes", a.getAllRecipes)
	router.Post("/recipes", a.addRecipe)
	router.Get("/recipes/{id}", a.getRecipe)
	router.Put("/recipes/{id}", a.updateRecipe)
	router.Delete("/recipes/{id}", a.deleteRecipe)

	router.Get("/recipes/suggestions", a.getRecipeSuggestions)

	if err := http.ListenAndServe(cfg.ApiAddr, router); err != nil {
		return fmt.Errorf("http.ListenAndServe returned error: %w", err)
	}

	return nil
}

func (a *api) healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&HealthzModel{
		IsHealthy:    true,
		ErrorMessage: "",
	}); err != nil {
		log.Printf("error sending healthz: %s\n", err.Error())
	}
}

func (a *api) getAllRecipes(w http.ResponseWriter, r *http.Request) {
	var recipes []Recipe
	if err := a.db.Model(&recipes).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if recipes == nil {
		recipes = []Recipe{}
	}
	var resp struct {
		Recipes []Recipe `json:"recipes"`
	}
	resp.Recipes = recipes

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending healthz: %s\n", err.Error())
	}
}

func (a *api) getRecipe(w http.ResponseWriter, r *http.Request) {
	rawId := chi.URLParam(r, "id")
	if rawId == "" {
		http.Error(w, "no id was provided in the request", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(rawId)
	if err != nil {
		http.Error(w, "provided id was not a number, thus couldn't be parsed", http.StatusBadRequest)
		return
	}

	var recipe Recipe
	if err := a.db.Model(&recipe).Where("recipe.Id = ?", id).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var resp struct {
		Recipe Recipe `json:"recipe"`
	}
	resp.Recipe = recipe

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending healthz: %s\n", err.Error())
	}
}

func (a *api) addRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, fmt.Sprintf("unable to decode recipe: %s", err.Error()), http.StatusBadRequest)
		return
	}

	result, err := a.db.Model(&recipe).Insert()
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't insert model into database: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	// return all recipes now
	a.getAllRecipes(w, r)
}

func (a *api) updateRecipe(w http.ResponseWriter, r *http.Request) {
	rawId := chi.URLParam(r, "id")
	if rawId == "" {
		http.Error(w, "no id was provided in the request", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(rawId)
	if err != nil {
		http.Error(w, "provided id was not a number, thus couldn't be parsed", http.StatusBadRequest)
		return
	}
	var recipe Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, fmt.Sprintf("couldn't decode body: %s", err.Error()), http.StatusBadRequest)
		return
	}
	if recipe.Id != int64(id) {
		http.Error(w, "the id in the request path does not equal the id in the request body", http.StatusBadRequest)
		return
	}

	result, err := a.db.Model(&recipe).Where("recipe.id = ?", id).Update()
	if err != nil {
		http.Error(w, fmt.Sprintf("an error occured when updating item with id: %d, error: %s", id, err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	a.getAllRecipes(w, r)
}

func (a *api) deleteRecipe(w http.ResponseWriter, r *http.Request) {
	rawId := chi.URLParam(r, "id")
	if rawId == "" {
		http.Error(w, "no id was provided in the request", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(rawId)
	if err != nil {
		http.Error(w, "provided id was not a number, thus couldn't be parsed", http.StatusBadRequest)
		return
	}

	result, err := a.db.Model(&Recipe{Id: int64(id)}).Where("recipe.id = ?", id).Delete()
	if err != nil {
		http.Error(w, fmt.Sprintf("an error occured when deleteing item with id: %d, error: %s", id, err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	a.getAllRecipes(w, r)
}

func (a *api) getRecipeSuggestions(w http.ResponseWriter, r *http.Request) {
	amountStr := r.URL.Query().Get("amount")
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		http.Error(w, "could not convert the amount param", http.StatusBadRequest)
		return
	}

	var recipes []Recipe
	if err := a.db.Model(&recipes).Select(); err != nil {
		http.Error(w, fmt.Sprintf("couldn't load item from database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	if len(recipes) <= amount {
		if err := json.NewEncoder(w).Encode(&recipes); err != nil {
			log.Printf("error sending suggestions: %s\n", err.Error())
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(recipes), func(i, j int) { recipes[i], recipes[j] = recipes[j], recipes[i] })
	recipes = recipes[:amount]
	if err := json.NewEncoder(w).Encode(&recipes); err != nil {
		log.Printf("error sending suggestions: %s\n", err.Error())
	}
}
