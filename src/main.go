package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Article représente un produit
type Article struct {
	ID          int
	Nom         string
	Description string
	Prix        float64
	Reduction   float64 // entre 0 et 1
	Stock       int
	Image       string
	CreatedAt   time.Time
}

// Variable globale des articles (au moins 5 articles comme demandé)
var articles = []Article{
	{ID: 1, Nom: "Sweat Vert", Description: "Sweat à capuche vert pastel Palace.", Prix: 90.00, Reduction: 0.15, Stock: 12, Image: "/static/img/16A.webp", CreatedAt: time.Now()},
	{ID: 2, Nom: "Sweat Noir", Description: "Sweat noir Palace avec logo London.", Prix: 95.00, Reduction: 0.0, Stock: 8, Image: "/static/img/18A.webp", CreatedAt: time.Now()},
	{ID: 3, Nom: "T-shirt Blanc", Description: "T-shirt blanc simple avec logo discret.", Prix: 25.00, Reduction: 0.10, Stock: 20, Image: "/static/img/1.png", CreatedAt: time.Now()},
	{ID: 4, Nom: "Casquette Beige", Description: "Casquette unisexe style urbain.", Prix: 29.99, Reduction: 0.0, Stock: 15, Image: "/static/img/2.png", CreatedAt: time.Now()},
	{ID: 5, Nom: "Pantalon Cargo", Description: "Pantalon cargo beige coupe large.", Prix: 80.00, Reduction: 0.20, Stock: 5, Image: "/static/img/3.png", CreatedAt: time.Now()},
	{ID: 6, Nom: "Sweat Bleu", Description: "Sweat bleu ciel confortable.", Prix: 70.00, Reduction: 0.10, Stock: 10, Image: "/static/img/4.png", CreatedAt: time.Now()},
	{ID: 7, Nom: "Bonnet Noir", Description: "Bonnet chaud 100% laine.", Prix: 19.99, Reduction: 0.0, Stock: 25, Image: "/static/img/5.png", CreatedAt: time.Now()},
	{ID: 8, Nom: "Veste Grise", Description: "Veste légère en coton.", Prix: 120.00, Reduction: 0.25, Stock: 6, Image: "/static/img/6.png", CreatedAt: time.Now()},
}

func main() {
	// Fonctions utilisables depuis les templates
	funcMap := template.FuncMap{
		"mul": func(a, b float64) float64 { return a * b },
		"sub": func(a, b float64) float64 { return a - b },
	}

	// Parse des templates avec funcMap
	temp, err := template.New("").Funcs(funcMap).ParseGlob("./templates/*.html")
	if err != nil {
		fmt.Println("Erreur de chargement des templates:", err)
		os.Exit(1)
	}

	// Serveur de fichiers statiques (css, images, js si besoin)
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Route racine -> index (menu)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := temp.ExecuteTemplate(w, "index", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Route : liste des articles (Challenge 01)
	http.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {
		if err := temp.ExecuteTemplate(w, "articles", articles); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Route : affichage d'un article par id (Challenge 02)
	http.HandleFunc("/article", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idStr)
		if id <= 0 {
			http.NotFound(w, r)
			return
		}
		for _, a := range articles {
			if a.ID == id {
				if err := temp.ExecuteTemplate(w, "detail", a); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
		http.NotFound(w, r)
	})

	// Route : afficher le formulaire d'ajout (GET) (Challenge 03)
	http.HandleFunc("/ajouter", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if err := temp.ExecuteTemplate(w, "ajouter", nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		// Si méthode non supportée, rediriger vers la page d'ajout
		http.Redirect(w, r, "/ajouter", http.StatusSeeOther)
	})

	// Route : traitement du formulaire (POST)
	http.HandleFunc("/ajouter/traitement", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/ajouter", http.StatusSeeOther)
			return
		}

		// Récupération et normalisation
		nom := strings.TrimSpace(r.FormValue("nom"))
		description := strings.TrimSpace(r.FormValue("description"))
		prixStr := strings.TrimSpace(r.FormValue("prix"))
		reducStr := strings.TrimSpace(r.FormValue("reduction"))
		stockStr := strings.TrimSpace(r.FormValue("stock"))

		// Validation minimale
		if nom == "" || description == "" || prixStr == "" {
			http.Error(w, "Nom, description et prix sont obligatoires.", http.StatusBadRequest)
			return
		}

		prix, err1 := strconv.ParseFloat(prixStr, 64)
		reduc := 0.0
		stock := 0
		var err2, err3 error
		if reducStr != "" {
			reduc, err2 = strconv.ParseFloat(reducStr, 64)
		}
		if stockStr != "" {
			stock, err3 = strconv.Atoi(stockStr)
		}

		if err1 != nil || (reducStr != "" && err2 != nil) || (stockStr != "" && err3 != nil) {
			http.Error(w, "Prix, réduction ou stock invalide.", http.StatusBadRequest)
			return
		}
		if prix <= 0 {
			http.Error(w, "Le prix doit être supérieur à 0.", http.StatusBadRequest)
			return
		}
		if reduc < 0 || reduc > 1 {
			http.Error(w, "La réduction doit être entre 0 et 1.", http.StatusBadRequest)
			return
		}
		if stock < 0 {
			http.Error(w, "Le stock ne peut pas être négatif.", http.StatusBadRequest)
			return
		}

		// Création et ajout du nouvel article
		newID := nextArticleID()
		newArticle := Article{
			ID:          newID,
			Nom:         nom,
			Description: description,
			Prix:        prix,
			Reduction:   reduc,
			Stock:       stock,
			Image:       "/static/img/placeholder.webp", // image par défaut
			CreatedAt:   time.Now(),
		}
		articles = append(articles, newArticle)

		// Redirection vers la page de détail du produit ajouté
		http.Redirect(w, r, fmt.Sprintf("/article?id=%d", newID), http.StatusSeeOther)
	})

	fmt.Println("✅ Serveur lancé sur http://localhost:8000")
	http.ListenAndServe("localhost:8000", nil)
}

// nextArticleID retourne un id unique simple (taille slice + 1)
func nextArticleID() int {
	max := 0
	for _, a := range articles {
		if a.ID > max {
			max = a.ID
		}
	}
	return max + 1
}
