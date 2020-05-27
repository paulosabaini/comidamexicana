package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/paulosabaini/comidamexicana/app/model"
	"github.com/paulosabaini/comidamexicana/config"
)

func GetIndexPage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	config := config.GetConfig()
	dir := config.Assets.Template
	tmpl := template.Must(template.ParseFiles(dir+"index.html", dir+"header.html", dir+"footer.html"))
	tmpl.ExecuteTemplate(w, "index", nil)
}

func GetContactPage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	config := config.GetConfig()
	dir := config.Assets.Template
	tmpl := template.Must(template.ParseFiles(dir+"contact.html", dir+"header.html", dir+"footer.html"))
	tmpl.ExecuteTemplate(w, "contact", nil)
}

func GetAboutPage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	config := config.GetConfig()
	dir := config.Assets.Template
	tmpl := template.Must(template.ParseFiles(dir+"about.html", dir+"header.html", dir+"footer.html"))
	tmpl.ExecuteTemplate(w, "about", nil)
}

func GetPlansPage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	config := config.GetConfig()
	dir := config.Assets.Template
	tmpl := template.Must(template.ParseFiles(dir+"plans.html", dir+"header.html", dir+"footer.html"))
	tmpl.ExecuteTemplate(w, "plans", nil)
}

func GetSubmitPage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	config := config.GetConfig()
	dir := config.Assets.Template
	tmpl := template.Must(template.ParseFiles(dir+"submit.html", dir+"header.html", dir+"footer.html"))
	tmpl.ExecuteTemplate(w, "submit", nil)
}

func GetSubmitSuccessPage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	config := config.GetConfig()
	dir := config.Assets.Template
	tmpl := template.Must(template.ParseFiles(dir+"submitsuccess.html", dir+"header.html", dir+"footer.html"))
	tmpl.ExecuteTemplate(w, "submitsuccess", nil)
}

func GetContactSuccessPage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	config := config.GetConfig()
	dir := config.Assets.Template
	tmpl := template.Must(template.ParseFiles(dir+"contactsuccess.html", dir+"header.html", dir+"footer.html"))
	tmpl.ExecuteTemplate(w, "contactsuccess", nil)
}

func GetRestaurantsPage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	config := config.GetConfig()
	dir := config.Assets.Template
	tmpl := template.Must(template.ParseFiles(dir+"search.html", dir+"header.html", dir+"footer.html"))
	tmpl.ExecuteTemplate(w, "search", nil)
}

func GetRestaurantPage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	config := config.GetConfig()
	dir := config.Assets.Template
	tmpl := template.Must(template.ParseFiles(dir+"restaurant.html", dir+"header.html", dir+"footer.html"))
	tmpl.ExecuteTemplate(w, "restaurant", nil)
}

func Get404Page(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	config := config.GetConfig()
	dir := config.Assets.Template
	tmpl := template.Must(template.ParseFiles(dir+"404.html", dir+"header.html", dir+"footer.html"))
	tmpl.ExecuteTemplate(w, "404", nil)
}

func GetBlogPage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	config := config.GetConfig()
	dir := config.Assets.Template
	tmpl := template.Must(template.ParseFiles(dir+"blog.html", dir+"header.html", dir+"footer.html", dir+"blog_widget.html"))
	posts := []model.Post{}
	db.Order("created_at desc").Find(&posts)
	tmpl.ExecuteTemplate(w, "blog", posts)
}

func GetBlogPostPage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	config := config.GetConfig()
	dir := config.Assets.Template
	tmpl := template.Must(template.ParseFiles(dir+"blog_post.html", dir+"header.html", dir+"footer.html", dir+"blog_widget.html"))
	vars := mux.Vars(r)
	postSlug := vars["post"]
	posts := []model.Post{}
	db.Where("title_slug = ?", postSlug).First(&posts)
	err := tmpl.ExecuteTemplate(w, "blogpost", posts)
	if err != nil {
		log.Println(err.Error())
	}
}

func GetBlogSearchPage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	config := config.GetConfig()
	dir := config.Assets.Template
	tmpl := template.Must(template.ParseFiles(dir+"blog.html", dir+"header.html", dir+"footer.html", dir+"blog_widget.html"))
	r.ParseForm()
	search := r.FormValue("search")
	posts := []model.Post{}
	db.Where("UPPER(content) LIKE UPPER(?) OR UPPER(title) LIKE UPPER(?)", "%"+search+"%", "%"+search+"%").Find(&posts)
	err := tmpl.ExecuteTemplate(w, "blog", posts)
	if err != nil {
		log.Println(err.Error())
	}
}
