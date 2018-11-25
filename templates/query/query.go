/*
* @Author: Bartuccio Antoine
* @Date:   2018-09-16 17:38:07
* @Last Modified by:   Bartuccio Antoine
* @Last Modified time: 2018-11-24 23:12:51
 */

package query

import (
	"strings"

	"gitlab.com/GT-RIMi/RIMo/media"
)

type QueryStatus int

const (
	Err      QueryStatus = 0
	Timeout  QueryStatus = 1
	Found    QueryStatus = 2
	NotFound QueryStatus = 3
)

// QueryStore represent a couple query/response for sessions
type QueryStore struct {
	Query    QueryStruct
	Response QueryResponse
}

// QueryDiscoverStruct format of a query to the discover backend
type QueryDiscoverStruct struct {
	MediaType media.MediaType `json:"media_type"`
	Media     media.Media     `json:"media"`
}

// QueryStruct format of the query
type QueryStruct struct {
	MediaType media.MediaType `json:"media_type"`
	Genres    []string        `json:"genres"`
	Title     string          `json:"title"`
	Year      int             `json:"year"`
	Adult     bool            `json:"adult"`
	Page      uint            `json:"page"`
	ID        string          `json:"id"`
}

// QueryResponse format of response
type QueryResponse struct {
	Medias []media.Media `json:"medias"`
	Status QueryStatus   `json:"status"`
}

// Search query using registered providers
func Search(query *QueryStruct, lastQuery *QueryStore) QueryResponse {
	if len(providers) == 0 {
		return QueryResponse{
			Medias: []media.Media{},
			Status: Err,
		}
	}
	for _, provider := range providers {
		if provider.IsCompatible(query.MediaType) {
			resp := provider.Search(*query)
			resp.removeLastQuery(query, lastQuery)
			return resp
		}
	}
	return QueryResponse{
		Medias: []media.Media{},
		Status: Err,
	}
}

// Discover search for related items
func Discover(query *QueryDiscoverStruct) QueryResponse {
	if len(providers) == 0 {
		return QueryResponse{
			Medias: []media.Media{},
			Status: Err,
		}
	}
	for _, provider := range providers {
		if provider.IsCompatible(query.MediaType) {
			resp := provider.Discover(*query)
			return resp
		}
	}
	return QueryResponse{
		Medias: []media.Media{},
		Status: Err,
	}
}

func (query *QueryStruct) compareExceptPage(q *QueryStruct) bool {
	if query.MediaType != q.MediaType || query.Year != q.Year || query.Adult != q.Adult {
		return false
	}

	if !strings.EqualFold(q.Title, query.Title) {
		return false
	}

	if len(q.Genres) != len(query.Genres) {
		return false
	}

	for i, genre := range query.Genres {
		if !strings.EqualFold(genre, q.Genres[i]) {
			return false
		}
	}

	return true
}

func (response *QueryResponse) removeLastQuery(query *QueryStruct, lastQuery *QueryStore) {
	if lastQuery == nil {
		return
	}
	if !(query.compareExceptPage(&lastQuery.Query) &&
		query.Page != 1 &&
		(lastQuery.Query.Page == query.Page+1 || lastQuery.Query.Page == query.Page-1)) {
		return
	}
	keys := make(map[string]bool)
	list := []media.Media{}
	for _, entry := range lastQuery.Response.Medias {
		keys[entry.ID] = true
	}
	for _, entry := range response.Medias {
		if _, value := keys[entry.ID]; !value {
			list = append(list, entry)
		}
	}
	response.Medias = list
}

// RemoveDuplicates Remove duplicates in a QueryResponse
func (response *QueryResponse) RemoveDuplicates() {
	keys := make(map[string]bool)
	list := []media.Media{}
	for _, entry := range response.Medias {
		if _, value := keys[entry.ID]; !value {
			list = append(list, entry)
		}
	}
	response.Medias = list
}
