/*
* @Author: Bartuccio Antoine
* @Date:   2018-09-18 22:22:00
* @Last Modified by:   klmp200
* @Last Modified time: 2018-10-01 23:44:48
 */

package goproviders

import (
	"testing"

	"gitlab.com/GT-RIMi/RIMo/media"
	"gitlab.com/GT-RIMi/RIMo/query"
	"gitlab.com/GT-RIMi/RIMo/settings"
)

func TestTmdbFilms(t *testing.T) {
	settings.InitSettings("../settings.json", "../settings_custom.json")
	provider := TmdbProvider{}
	provider.Init()
	testResponse(t, provider.Search(query.QueryStruct{
		MediaType: media.Movie,
		Title:     "Gravity Falls",
	}))
}

func TestTmdbSeries(t *testing.T) {
	settings.InitSettings("../settings.json", "../settings_custom.json")
	provider := TmdbProvider{}
	provider.Init()
	testResponse(t, provider.Search(query.QueryStruct{
		MediaType: media.Serie,
		Title:     "Gravity Falls",
	}))
}

func testResponse(t *testing.T, resp query.QueryResponse) {
	if resp.Status == query.Err {
		t.Error("error status")
	}

	if resp.Status == query.NotFound {
		t.Error("no gravity falls found ?! are you fucking kidding me ?!")
	}

	for _, media := range resp.Medias {
		t.Log(media.Title)
	}
}
