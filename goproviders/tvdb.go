/*
* @Author: Bartuccio Antoine
* @Date:   2018-09-18 21:58:32
* @Last Modified by:   klmp200
* @Last Modified time: 2018-09-20 00:46:38
 */

package goproviders

import (
	"github.com/danesparza/tvdb"
	"gitlab.com/GT-RIMi/RIMo/media"
	"gitlab.com/GT-RIMi/RIMo/query"
)

type TvdbProvider struct {
	client tvdb.Client
}

func (db *TvdbProvider) Init() {
	db.client = tvdb.Client{}
}

func (db *TvdbProvider) IsCompatible(mediaType media.MediaType) bool {
	return mediaType == media.Serie
}

func (db *TvdbProvider) Search(q query.QueryStruct) query.QueryResponse {
	return query.QueryResponse{}
}
