package handler

import (
	"assignment-2/config"
	"assignment-2/models"
	"assignment-2/utils"
	"fmt"
	"gorm.io/gorm"
	"html/template"
	"net/http"
)

var (
	store = config.PgStore()
	db    = config.GetDB()
)

func RouteRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var tmpl = template.Must(template.New("register").ParseFiles("layout/register.html"))
		var err = tmpl.Execute(w, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var firstname = r.FormValue("firstname")
		var lastname = r.Form.Get("lastname")
		var username = r.Form.Get("username")
		var password = r.Form.Get("password")

		hashedPassword, _ := utils.HashPassword(password)

		user := models.User{}

		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				http.Redirect(w, r, "/register", http.StatusSeeOther)
				return
			}
		}

		if user.ID > 0 {
			var tmpl = template.Must(template.New("register").ParseFiles("layout/register.html"))
			var err = tmpl.Execute(w, models.Layout{
				Message: "Account Already Exists",
			})

			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}

			return
		}

		user = models.User{
			FirstName: firstname,
			LastName:  lastname,
			Password:  hashedPassword,
			Username:  username,
		}

		if err := db.Create(&user).Error; err != nil {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func RouteLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		session, err := store.Get(r, "session")
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Options.MaxAge = -1
		session.Save(r, w)

		var tmpl = template.Must(template.New("login").ParseFiles("layout/login.html"))
		err = tmpl.Execute(w, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}
