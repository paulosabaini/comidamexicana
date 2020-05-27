package model

import (
	"html/template"

	"github.com/gomarkdown/markdown"
	"github.com/jinzhu/gorm"
)

type Address struct {
	ID           int    `json:"id"`
	Cep          string `json:"cep"`
	State        string `json:"state"`
	StateSlug    string `json:"state_slug"`
	City         string `json:"city"`
	CitySlug     string `json:"city_slug"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Number       int    `json:"number"`
}

type OpeningHour struct {
	ID           int    `json:"id"`
	RestaurantID int    `json:"restaurant_id"`
	Weekday      string `json:"weekday"`
	OpeningTime  string `json:"opening_time"`
	ClosingTime  string `json:"closing_time"`
}

type Option struct {
	ID           int  `json:"id"`
	Reservations bool `json:"reservations"`
	Parking      bool `json:"parking"`
	Takeout      bool `json:"takeout"`
	Delivery     bool `json:"delivery"`
	Wifi         bool `json:"wifi"`
	Wheelchair   bool `json:"wheelchair"`
	CreditCard   bool `json:"credit_card"`
	Outdoor      bool `json:"outdoor"`
	Alcohol      bool `json:"alcohol"`
}

type Image struct {
	ID           int    `json:"id"`
	RestaurantID int    `json:"restaurant_id"`
	Filename     string `json:"filename"`
	Url          string `json:"url"`
	Main         bool   `json:"main"`
}

type Plan struct {
	gorm.Model
	Description    string  `json:"description"`
	MonthlyPrice   float64 `json:"monthly_price"`
	AnnualPrice    float64 `json:"annual_rice"`
	NumberOfImages int     `json:"number_of_images"`
	Number         int     `json:"number"`
}

type Restaurant struct {
	gorm.Model
	Name                string        `json:"name"`
	NameSlug            string        `json:"name_slug"`
	Email               string        `json:"email"`
	ContactName         string        `json:"contactName"`
	Description         string        `json:"description"`
	Website             string        `json:"website"`
	Phone               string        `json:"phone"`
	MobilePhone         string        `json:"mobilePhone"`
	Whatsapp            bool          `json:"whatsapp"`
	AddressID           int           `json:"address_id"`
	Address             Address       `json:"address"`
	OpeningHours        []OpeningHour `json:"opening_hours"`
	OptionID            int           `json:"option_id"`
	Option              Option        `json:"options"`
	AveragePricePerMeal float64       `json:"average_price_per_meal"`
	Images              []Image       `json:"images"`
	PlanID              int           `json:"plan_id"`
	Plan                Plan          `json:"plan" gorm:"association_autoupdate:false"`
	Published           bool          `json:"published"`
	VideoLink           string        `json:"video_link"`
	InstagramLink       string        `json:"instagram_link"`
	FacebookLink        string        `json:"facebook_link"`
	TwitterLink         string        `json:"twitter_link"`
	OpinionPhrase       string        `json:"opinion_phrase"`
	OpitionAuthor       string        `json:"opinion_author"`
	IFoodLink           string        `json:"ifood_link"`
}

type NewsletterEmail struct {
	gorm.Model
	Email string `json:"email"`
}

type Post struct {
	gorm.Model
	Title      string `json:"title"`
	TitleSlug  string `json:"title_slug"`
	Author     string `json:"author"`
	CoverImage string `json:"cover_image"`
	Abstract   string `json:"abstract"`
	Content    string `json:"content"`
}

func (p *Post) CreatedDate() string {
	return p.CreatedAt.Format("02/01/2006")
}

func (p *Post) CreatedDateTime() string {
	return p.CreatedAt.Format("02/01/2006 15:04:05")
}

func (p *Post) FormattedContent() template.HTML {
	md := []byte(p.Content)
	return template.HTML(markdown.ToHTML(md, nil, nil))
}

func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Restaurant{}, &Address{}, &OpeningHour{}, &Option{}, &Image{}, &Plan{}, &Restaurant{}, &NewsletterEmail{}, &Post{})
	db.Model(&OpeningHour{}).AddForeignKey("restaurant_id", "restaurants(id)", "CASCADE", "CASCADE")
	db.Model(&Image{}).AddForeignKey("restaurant_id", "restaurants(id)", "CASCADE", "CASCADE")
	db.Model(&Restaurant{}).AddForeignKey("address_id", "addresses(id)", "CASCADE", "CASCADE")
	db.Model(&Restaurant{}).AddForeignKey("option_id", "options(id)", "CASCADE", "CASCADE")
	db.Model(&Restaurant{}).AddForeignKey("plan_id", "plans(id)", "CASCADE", "CASCADE")

	plan := &Plan{}

	db.First(plan)

	if *plan == (Plan{}) {
		gratis := Plan{Description: "gratis", MonthlyPrice: 0, AnnualPrice: 0, NumberOfImages: 3, Number: 0}
		interm := Plan{Description: "intermediario", MonthlyPrice: 30, AnnualPrice: 360, NumberOfImages: 6, Number: 1}
		comp := Plan{Description: "completo", MonthlyPrice: 60, AnnualPrice: 620, NumberOfImages: 12, Number: 2}

		db.Save(&gratis)
		db.Save(&interm)
		db.Save(&comp)
	}

	return db
}
