package pages

import (
	"clickyab.com/crane/models/internal/entities"
	"clickyab.com/crane/workers/models"
	"github.com/clickyab/services/pool"
)

var pagesPool pool.Interface

//GetPages get total seats of all creatives in network per type
//TODO: should use seat interface instead of structure
func GetPages() []entities.PublisherPage {
	data := pagesPool.All()
	all := make([]entities.PublisherPage, len(data))

	var c int
	for i := range data {
		all[c] = data[i].(entities.PublisherPage)
		c++
	}

	return all
}

// GetPageByKeys try to get network seats based on its creative type
//TODO: should use seat interface instead of structure
func GetPageByKeys(publisherDomain, URLKey string) *entities.PublisherPage {
	data := pagesPool.All()

	key := entities.GenPagePoolKey(publisherDomain, URLKey)

	page := data[key]
	if page == nil {
		return nil
	}

	return page.(*entities.PublisherPage)
}

// GetByURLAndDomain try to get network seats based on its creative type
func GetByURLAndDomain(publisherDomain, URL string) *entities.PublisherPage {
	ukey := entities.GenerateURLKey(URL)
	return GetPageByKeys(publisherDomain, ukey)
}

// AddAndGetPublisherPage get page by key in pool if not found select on db and if not found again inser it
func AddAndGetPublisherPage(im models.Impression) (*entities.PublisherPage, error) {
	ukey := entities.GenerateURLKey(im.ParentURL)
	page := GetPageByKeys(im.Publisher, ukey)

	if page != nil {
		return page, nil
	}

	return entities.AddAndGetPublisherPage(im)
}
