package handler

import (
	"assignment-2/models"
	"assignment-2/utils"
	"fmt"
	"gorm.io/gorm"
	"html/template"
	"net/http"
)

func RouteIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		session, _ := store.Get(r, "session")
		val := session.Values["user"]
		user, ok := val.(models.User)

		if !ok {
			var tmpl = template.Must(template.New("auth").ParseFiles("layout/auth.html"))
			var err = tmpl.Execute(w, nil)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		var tmpl = template.Must(template.New("home").ParseFiles("layout/home.html"))
		var err = tmpl.Execute(w, user)

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

		var username = r.FormValue("username")
		var password = r.Form.Get("password")

		user := models.User{}

		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				var tmpl = template.Must(template.New("auth").ParseFiles("layout/auth.html"))
				var err = tmpl.Execute(w, models.Layout{
					Message: "Account not found, please check if the username and password are correct.",
				})

				if err != nil {
					http.Redirect(w, r, "/auth", http.StatusSeeOther)
				}

				return
			}
		}

		if ok := utils.CheckPasswordHash(password, user.Password); !ok {
			var tmpl = template.Must(template.New("auth").ParseFiles("layout/auth.html"))
			var err = tmpl.Execute(w, models.Layout{
				Message: "wrong password",
			})

			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			return
		}

		session, err := store.Get(r, "session")
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Values["user"] = user
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}
