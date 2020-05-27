package handler

import (
	"math/rand"
	"net/http"
	"net/smtp"
	"net/url"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
	"github.com/paulosabaini/comidamexicana/app/model"
	"github.com/paulosabaini/comidamexicana/config"
)

func GetRestaurantsByCity(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	state := vars["state"]
	city := vars["city"]

	restaurants := []model.Restaurant{}

	db.Preload("Address").Preload("Option").Preload("OpeningHours").Preload("Images").Preload("Plan").Joins("Join addresses on restaurants.address_id = addresses.id").Joins("left join plans on restaurants.plan_id = plans.id").Where("addresses.state_slug = ? AND addresses.city_slug = ? AND restaurants.published = true", state, city).Order("plans.annual_price desc").Find(&restaurants)

	respondJSON(w, http.StatusOK, restaurants)
}

func GetRestaurantsByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	restaurant := model.Restaurant{}
	db.Preload("Address").Preload("Option").Preload("OpeningHours").Preload("Images").Preload("Plan").First(&restaurant, id)
	if restaurant.ID != 0 {
		respondJSON(w, http.StatusOK, restaurant)
	} else {
		respondJSON(w, http.StatusOK, "")
	}
}

func GetRestaurantByName(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	state := vars["state"]
	city := vars["city"]
	slug := vars["name"]

	restaurant := model.Restaurant{}
	db.Preload("Address").Preload("Option").Preload("OpeningHours").Preload("Images").Preload("Plan").Joins("Join addresses on restaurants.address_id = addresses.id").Where("addresses.state_slug = ? AND addresses.city_slug = ? AND restaurants.published = true AND restaurants.name_slug = ?", state, city, slug).First(&restaurant)
	if restaurant.ID != 0 {
		respondJSON(w, http.StatusOK, restaurant)
	} else {
		respondJSON(w, http.StatusOK, "")
	}
}

func CreateRestaurant(db *gorm.DB, session *session.Session, w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(100000)

	restaurant := model.Restaurant{}

	restaurant.ContactName = r.FormValue("nome")
	restaurant.Email = r.FormValue("email")
	restaurant.Name = r.FormValue("nomeEmpresa")
	restaurant.NameSlug = slug.Make(r.FormValue("nomeEmpresa"))
	restaurant.Description = r.FormValue("descricao")
	website, _ := url.Parse(r.FormValue("website"))
	website.Scheme = "https"
	restaurant.Website = website.String()
	restaurant.Phone = r.FormValue("telefoneFixo")
	restaurant.MobilePhone = r.FormValue("telefoneCel")
	videoLink, _ := url.Parse(r.FormValue("youtubeLink"))
	videoLink.Scheme = "https"
	restaurant.VideoLink = videoLink.String()
	instagramLink, _ := url.Parse(r.FormValue("instagramLink"))
	instagramLink.Scheme = "https"
	restaurant.InstagramLink = instagramLink.String()
	facebookLink, _ := url.Parse(r.FormValue("facebookLink"))
	facebookLink.Scheme = "https"
	restaurant.FacebookLink = facebookLink.String()
	twitterLink, _ := url.Parse(r.FormValue("twitterLink"))
	twitterLink.Scheme = "https"
	restaurant.TwitterLink = twitterLink.String()
	ifoodLink, _ := url.Parse(r.FormValue("ifoodLink"))
	ifoodLink.Scheme = "https"
	restaurant.IFoodLink = ifoodLink.String()
	restaurant.Whatsapp = false
	if r.FormValue("whatsapp") == "on" {
		restaurant.Whatsapp = true
	}

	address := model.Address{}
	address.Cep = r.FormValue("cep")
	address.State = r.FormValue("uf")
	address.StateSlug = slug.Make(r.FormValue("uf"))
	address.City = r.FormValue("cidade")
	address.CitySlug = slug.Make(r.FormValue("cidade"))
	address.Neighborhood = r.FormValue("bairro")
	address.Street = r.FormValue("rua")
	address.Number, _ = strconv.Atoi(r.FormValue("numero"))

	restaurant.Address = address

	openingHours := []model.OpeningHour{}
	if r.FormValue("seg") == "on" {
		if r.FormValue("aseg") != "-1" && r.FormValue("fseg") != "-1" {
			seg := model.OpeningHour{Weekday: "seg", OpeningTime: r.FormValue("aseg"), ClosingTime: r.FormValue("fseg")}
			openingHours = append(openingHours, seg)
		}
		if r.FormValue("aseg2") != "-1" && r.FormValue("fseg2") != "-1" {
			seg2 := model.OpeningHour{Weekday: "seg", OpeningTime: r.FormValue("aseg2"), ClosingTime: r.FormValue("fseg2")}
			openingHours = append(openingHours, seg2)
		}
	}
	if r.FormValue("ter") == "on" {
		if r.FormValue("ater") != "-1" && r.FormValue("fter") != "-1" {
			ter := model.OpeningHour{Weekday: "ter", OpeningTime: r.FormValue("ater"), ClosingTime: r.FormValue("fter")}
			openingHours = append(openingHours, ter)
		}
		if r.FormValue("ater2") != "-1" && r.FormValue("fter2") != "-1" {
			ter2 := model.OpeningHour{Weekday: "ter", OpeningTime: r.FormValue("ater2"), ClosingTime: r.FormValue("fter2")}
			openingHours = append(openingHours, ter2)
		}
	}
	if r.FormValue("qua") == "on" {
		if r.FormValue("aqua") != "-1" && r.FormValue("fqua") != "-1" {
			qua := model.OpeningHour{Weekday: "qua", OpeningTime: r.FormValue("aqua"), ClosingTime: r.FormValue("fqua")}
			openingHours = append(openingHours, qua)
		}
		if r.FormValue("aqua2") != "-1" && r.FormValue("fqua2") != "-1" {
			qua2 := model.OpeningHour{Weekday: "qua", OpeningTime: r.FormValue("aqua2"), ClosingTime: r.FormValue("fqua2")}
			openingHours = append(openingHours, qua2)
		}
	}
	if r.FormValue("qui") == "on" {
		if r.FormValue("aqui") != "-1" && r.FormValue("fqui") != "-1" {
			qui := model.OpeningHour{Weekday: "qui", OpeningTime: r.FormValue("aqui"), ClosingTime: r.FormValue("fqui")}
			openingHours = append(openingHours, qui)
		}
		if r.FormValue("aqui2") != "-1" && r.FormValue("fqui2") != "-1" {
			qui2 := model.OpeningHour{Weekday: "qui", OpeningTime: r.FormValue("aqui2"), ClosingTime: r.FormValue("fqui2")}
			openingHours = append(openingHours, qui2)
		}
	}
	if r.FormValue("sex") == "on" {
		if r.FormValue("asex") != "-1" && r.FormValue("fsex") != "-1" {
			sex := model.OpeningHour{Weekday: "sex", OpeningTime: r.FormValue("asex"), ClosingTime: r.FormValue("fsex")}
			openingHours = append(openingHours, sex)
		}
		if r.FormValue("asex2") != "-1" && r.FormValue("fsex2") != "-1" {
			sex2 := model.OpeningHour{Weekday: "sex", OpeningTime: r.FormValue("asex2"), ClosingTime: r.FormValue("fsex2")}
			openingHours = append(openingHours, sex2)
		}
	}
	if r.FormValue("sab") == "on" {
		if r.FormValue("asab") != "-1" && r.FormValue("fsab") != "-1" {
			sab := model.OpeningHour{Weekday: "sab", OpeningTime: r.FormValue("asab"), ClosingTime: r.FormValue("fsab")}
			openingHours = append(openingHours, sab)
		}
		if r.FormValue("asab2") != "-1" && r.FormValue("fsab2") != "-1" {
			sab2 := model.OpeningHour{Weekday: "sab", OpeningTime: r.FormValue("asab2"), ClosingTime: r.FormValue("fsab2")}
			openingHours = append(openingHours, sab2)
		}
	}
	if r.FormValue("dom") == "on" {
		if r.FormValue("adom") != "-1" && r.FormValue("fdom") != "-1" {
			dom := model.OpeningHour{Weekday: "dom", OpeningTime: r.FormValue("adom"), ClosingTime: r.FormValue("fdom")}
			openingHours = append(openingHours, dom)
		}
		if r.FormValue("adom2") != "-1" && r.FormValue("fdom2") != "-1" {
			dom2 := model.OpeningHour{Weekday: "dom", OpeningTime: r.FormValue("adom2"), ClosingTime: r.FormValue("fdom2")}
			openingHours = append(openingHours, dom2)
		}
	}

	restaurant.OpeningHours = openingHours

	option := model.Option{}

	option.Reservations = false
	if r.FormValue("reserva") == "1" {
		option.Reservations = true
	}
	option.Parking = false
	if r.FormValue("estacionamento") == "1" {
		option.Parking = true
	}
	option.Takeout = false
	if r.FormValue("embalagem") == "1" {
		option.Takeout = true
	}
	option.Delivery = false
	if r.FormValue("entrega") == "1" {
		option.Delivery = true
	}
	option.Wifi = false
	if r.FormValue("wifi") == "1" {
		option.Wifi = true
	}
	option.Wheelchair = false
	if r.FormValue("cadeirantes") == "1" {
		option.Wheelchair = true
	}
	option.CreditCard = false
	if r.FormValue("cartaocredito") == "1" {
		option.CreditCard = true
	}
	option.Outdoor = false
	if r.FormValue("exterior") == "1" {
		option.Outdoor = true
	}
	option.Alcohol = false
	if r.FormValue("bebidas") == "1" {
		option.Alcohol = true
	}

	restaurant.Option = option

	restaurant.AveragePricePerMeal, _ = strconv.ParseFloat(r.FormValue("precoMedio"), 64)

	images := []model.Image{}

	// // File upload - http://sanatgersappa.blogspot.com/2013/03/handling-multiple-file-uploads-in-go.html
	m := r.MultipartForm

	files := m.File["fotos"]
	config := config.GetConfig()
	for i, _ := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		filename := restaurant.NameSlug + "-" + strconv.Itoa(rand.Intn(10000)) + ".jpg"

		// Upload the file to S3.
		uploader := s3manager.NewUploader(session)

		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(config.Image.AwsS3Bucket),
			Key:    aws.String(filename),
			ACL:    aws.String("public-read"),
			Body:   file,
		})

		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		image := model.Image{}
		image.Filename = filename
		image.Url = config.Image.Url + filename
		image.Main = false
		images = append(images, image)
	}

	restaurant.Images = images

	plan := &model.Plan{}

	if r.FormValue("plano") == "0" {
		db.Where(&model.Plan{Description: "gratis"}).First(plan)
	} else if r.FormValue("plano") == "1" {
		db.Where(&model.Plan{Description: "intermediario"}).First(plan)
	} else if r.FormValue("plano") == "2" {
		db.Where(&model.Plan{Description: "completo"}).First(plan)
	}

	restaurant.Plan = *plan
	restaurant.Published = false

	if err := db.Save(&restaurant).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go sendEmail(restaurant.Name)

	http.Redirect(w, r, "/sucesso-envio", http.StatusSeeOther)
}

func GetRestaurantsByPlan(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	planNumber := vars["plan"]
	restaurants := []model.Restaurant{}
	db.Preload("Images").Preload("Plan").Preload("Address").Joins("Join plans on restaurants.plan_id = plans.id").Where("plans.number = ? AND restaurants.published = true", planNumber).Find(&restaurants)
	respondJSON(w, http.StatusOK, restaurants)

}

func GetCities(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	cities := []model.Address{}
	db.Select("distinct city, city_slug, state, state_slug").Find(&cities)
	respondJSON(w, http.StatusOK, cities)
}

func FilterRestaurants(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	where, opening, options, price := "", "", "", ""

	for key, value := range r.Form {
		if key == "seg" || key == "ter" || key == "qua" || key == "qui" || key == "sex" || key == "sab" || key == "dom" {
			if opening != "" {
				opening += " AND"
			}
			opening += " opening_hours.weekday = " + "'" + key + "'"
		}

		if key == "reservations" || key == "parking" || key == "takeout" || key == "delivery" || key == "wifi" || key == "wheelchair" || key == "credit_card" || key == "outdoor" || key == "alcohol" {
			if options != "" {
				options += " AND"
			}
			options += " options." + key + " = true"
		}

		if key == "precoMedio" {
			if value[0] == "" {
				price = ""
			} else {
				price += " restaurants.average_price_per_meal = " + value[0]
			}
		}
	}

	if opening != "" {
		openingTime := r.FormValue("openingTime")
		closingTime := r.FormValue("closingTime")
		opening += " AND opening_hours.opening_time <= " + "'" + openingTime + "'" + " AND opening_hours.closing_time >= " + "'" + closingTime + "'"
	}

	if opening != "" {
		if options != "" || price != "" {
			where += opening + " AND"
		} else {
			where += opening
		}
	}
	if options != "" {
		if price != "" {
			where += options + " AND"
		} else {
			where += options
		}
	}
	if price != "" {
		where += price
	}

	vars := mux.Vars(r)
	state := vars["state"]
	city := vars["city"]

	and := ""

	if opening == "" && options == "" && price == "" {
		and = ""
	} else {
		and = " AND "
	}

	where += and + "addresses.state_slug = " + "'" + state + "'" + " AND addresses.city_slug = " + "'" + city + "'" + " AND restaurants.published = true"

	restaurants := []model.Restaurant{}

	options_join, opening_join := "", ""
	if options != "" {
		options_join = "Join options on restaurants.option_id = options.id"
	}
	if opening != "" {
		opening_join = "Join opening_hours on restaurants.id = opening_hours.restaurant_id"
	}
	db.Preload("Address").Preload("Option").Preload("OpeningHours").Preload("Images").Preload("Plan").Joins("Join addresses on restaurants.address_id = addresses.id").Joins("Join plans on restaurants.plan_id = plans.id").Joins(options_join).Joins(opening_join).Where(where).Order("plans.annual_price desc").Find(&restaurants)

	respondJSON(w, http.StatusOK, restaurants)
}

func FilterRestaurantsByName(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("nome")

	vars := mux.Vars(r)
	state := vars["state"]
	city := vars["city"]

	restaurants := []model.Restaurant{}

	db.Preload("Address").Preload("Option").Preload("OpeningHours").Preload("Images").Preload("Plan").Joins("Join addresses on restaurants.address_id = addresses.id").Joins("Join plans on restaurants.plan_id = plans.id").Where("addresses.state_slug = ? AND addresses.city_slug = ? AND restaurants.published = true AND UPPER(restaurants.name) LIKE UPPER(?)", state, city, "%"+name+"%").Order("plans.annual_price desc").Find(&restaurants)

	respondJSON(w, http.StatusOK, restaurants)
}

func sendEmail(restaurant string) {
	config := config.GetConfig()

	auth := smtp.PlainAuth("", config.Email.Username, config.Email.Password, config.Email.Server)

	to := []string{config.Email.To}
	msg := []byte("To: " + config.Email.To + "\r\n" +
		"Subject: Novo restaurante cadastrado comidamexicana.com.br\r\n" +
		"\r\n" +
		"Nome: " + restaurant)
	smtp.SendMail(config.Email.Server+":"+config.Email.Port, auth, config.Email.Sender, to, msg)

}
