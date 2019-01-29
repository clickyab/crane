package entities

import (
	openrtb "clickyab.com/crane/openrtb/v2.5"
)

//
//
// list: bamilo
// url: https://www.bamilo.com/product/active-%D9%BE%D9%88%D8%AF%D8%B1-%D9%84%D8%A8%D8%A7%D8%B3%D8%B4%D9%88%DB%8C%DB%8C-%D9%85%D8%A7%D8%B4%DB%8C%D9%86%DB%8C-500-%DA%AF%D8%B1%D9%85%DB%8C-9395631/
// img: //media.bamilo.com/p/active-1843-1365939-1-zoom.jpg
// title: پودر لباسشویی ماشینی 500 گرمی
// price: 5355
// discount: 10
// sku: AC696OT084RNIALIYUN
// isavailable: true
// category: سوپرمارکت,بهداشت منزل,شوینده لباس
// brand: Active

// Asset for product
type Asset struct {
	List        string   `json:"list"`
	URL         string   `json:"url"`
	Img         string   `json:"img"`
	Title       string   `json:"title"`
	Price       int64    `json:"price"`
	Discount    int64    `json:"discount"`
	SKU         string   `json:"sku"`
	IsAvailable bool     `json:"is_available"`
	Category    []string `json:"category"`
	Brand       string   `json:"brand"`

	User *openrtb.User `json:"-"`
}
