package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis"
	"github.com/rubenwo/home-automation/food-service/pkg/database"
	"log"
	"net/http"
	"time"
)

type api struct {
	db database.Database
}

func main() {
	db, err := database.Factory("redis")
	if err != nil {
		log.Fatal(err)
	}
	a := api{db: db}

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

	router.Get("/recipes", a.getRecipes)
	router.Post("/recipes", a.addRecipe)
	router.Delete("/recipes/{id}", a.deleteRecipe)
	fmt.Println(router)
	log.Println("food-service is online")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Println(err)
	}
	log.Println("food-service is offline")
}

func (a *api) getRecipes(w http.ResponseWriter, r *http.Request) {
	v, err := a.db.Get("food-recipes")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	raw, ok := v.(string)
	if !ok {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var recipes []Recipe
	if err := json.Unmarshal([]byte(raw), &recipes); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var resp struct {
		Recipes []Recipe `json:"recipes"`
	}
	resp.Recipes = recipes
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending getDevices: %s\n", err.Error())
	}
}
func (a *api) addRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	v, err := a.db.Get("food-recipes")
	if err != nil && err != redis.Nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(v)

	var recipes []Recipe
	if v != "" && v != nil && err != redis.Nil {
		if err := json.Unmarshal([]byte(v.(string)), &recipes); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	recipes = append(recipes, recipe)
	jsonData, err := json.Marshal(recipes)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.db.Set("food-recipes", jsonData); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var resp struct {
		Recipes []Recipe `json:"recipes"`
	}
	resp.Recipes = recipes
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending getDevices: %s\n", err.Error())
	}
}
func (a *api) deleteRecipe(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	v, err := a.db.Get("food-recipes")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var recipes []Recipe
	if err := json.Unmarshal([]byte(v.(string)), &recipes); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, recipe := range recipes {
		if recipe.ID == id {
			recipes = append(recipes[:i], recipes[i+1:]...)
		}
	}

	jsonData, err := json.Marshal(recipes)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.db.Set("food-recipes", jsonData); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var resp struct {
		Msg string `json:"msg"`
	}
	resp.Msg = fmt.Sprintf("deleted recipe with id: %s successfully", id)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Printf("error sending getDevices: %s\n", err.Error())
	}
}
