package handler

import (
	"net/http"
	"net/smtp"

	"github.com/jinzhu/gorm"
	"github.com/paulosabaini/comidamexicana/config"
)

func SendContactEmail(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.FormValue("nome")
	email := r.FormValue("email")
	phone := r.FormValue("telefone")
	website := r.FormValue("website")
	message := r.FormValue("mensagem")

	config := config.GetConfig()

	auth := smtp.PlainAuth("", config.Email.Username, config.Email.Password, config.Email.Server)

	to := []string{config.Email.To}
	msg := []byte("To: " + config.Email.To + "\r\n" +
		"Subject: Nova mensagem comidamexicana.com.br\r\n" +
		"\r\n" +
		"Nome: " + name + " \r\n" +
		"Email: " + email + " \r\n" +
		"Telefone: " + phone + " \r\n" +
		"Website: " + website + " \r\n" +
		"Mensagem: " + message + " \r\n")
	err := smtp.SendMail(config.Email.Server+":"+config.Email.Port, auth, config.Email.Sender, to, msg)

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}

	http.Redirect(w, r, "/sucesso-contato", http.StatusSeeOther)
}
