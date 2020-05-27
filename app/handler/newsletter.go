package handler

import (
	"encoding/json"
	"net/http"

	"github.com/paulosabaini/comidamexicana/app/model"

	"github.com/jinzhu/gorm"
)

func InsertNewsletterEmail(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	email := model.NewsletterEmail{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&email); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	}

	defer r.Body.Close()

	if err := db.Save(&email).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}

	return
}
