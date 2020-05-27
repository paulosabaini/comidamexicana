package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/paulosabaini/comidamexicana/app/handler"
	"github.com/paulosabaini/comidamexicana/app/model"
	"github.com/paulosabaini/comidamexicana/config"
)

type App struct {
	Router     *mux.Router
	DB         *gorm.DB
	AwsSession *session.Session
}

func (a *App) Initialize(config *config.Config) {
	dbURI := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		config.DB.Host,
		config.DB.Port,
		config.DB.Username,
		config.DB.Name,
		config.DB.Password)

	db, err := gorm.Open(config.DB.Dialect, dbURI)

	if err != nil {
		log.Println(err.Error())
		log.Fatal("Could not connect to database")
	}

	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Image.AwsS3Region),
		Credentials: credentials.NewStaticCredentials(
			config.Image.AwsId,
			config.Image.AwsSecretKey,
			""),
	})
	if err != nil {
		panic(err)
	}

	a.AwsSession = awsSession

	a.DB = model.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.setRouters()
	fs := http.FileServer(neuteredFileSystem{http.Dir(config.Assets.Static)})
	a.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	//a.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(config.Assets.Static))))
}

func (a *App) setRouters() {
	a.Get("/", a.handleRequest(handler.GetIndexPage))
	a.Get("/contato", a.handleRequest(handler.GetContactPage))
	a.Get("/sobre", a.handleRequest(handler.GetAboutPage))
	a.Get("/listar-negocio", a.handleRequest(handler.GetPlansPage))
	a.Get("/envio-negocio/{plano}", a.handleRequest(handler.GetSubmitPage))
	a.Get("/sucesso-envio", a.handleRequest(handler.GetSubmitSuccessPage))
	a.Get("/sucesso-contato", a.handleRequest(handler.GetContactSuccessPage))
	a.Get("/restaurantes/{uf}/{cidade}", a.handleRequest(handler.GetRestaurantsPage))
	a.Get("/restaurantes/{state}/{city}/{name}", a.handleRequest(handler.GetRestaurantPage))
	a.Get("/blog", a.handleRequest(handler.GetBlogPage))
	a.Get("/blog/{post}", a.handleRequest(handler.GetBlogPostPage))
	a.Post("/blog/pesquisar", a.handleRequest(handler.GetBlogSearchPage))
	a.Post("/restaurantes/create", a.handleRequestSession(handler.CreateRestaurant))
	a.Post("/newsletter/create", a.handleRequest(handler.InsertNewsletterEmail))
	a.Post("/contact/create", a.handleRequest(handler.SendContactEmail))
	a.Get("/api/cities", a.handleRequest(handler.GetCities))
	a.Get("/api/restaurantes/{state}/{city}", a.handleRequest(handler.GetRestaurantsByCity))
	a.Post("/api/restaurantes/filtrar/{state}/{city}", a.handleRequest(handler.FilterRestaurants))
	a.Post("/api/restaurantes/filtrar/nome/{state}/{city}", a.handleRequest(handler.FilterRestaurantsByName))
	a.Get("/api/restaurantes/{state}/{city}/{name}", a.handleRequest(handler.GetRestaurantByName))
	a.Get("/api/restaurantes/{plan}", a.handleRequest(handler.GetRestaurantsByPlan))
	a.Router.NotFoundHandler = a.handleRequest(handler.Get404Page)
}

func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

type RequestHandlerFunction func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, w, r)
	}
}

type RequestHandlerFunctionSession func(db *gorm.DB, s *session.Session, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequestSession(handler RequestHandlerFunctionSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, a.AwsSession, w, r)
	}
}

// neuteredFileSystem is used to prevent directory listing of static assets
type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	// Check if path exists
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	// If path exists, check if is a file or a directory.
	// If is a directory, stop here with an error saying that file
	// does not exist. So user will get a 404 error code for a file or directory
	// that does not exist, and for directories that exist.
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		return nil, os.ErrNotExist
	}

	// If file exists and the path is not a directory, let's return the file
	return f, nil
}
