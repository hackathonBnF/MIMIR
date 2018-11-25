/*
* @Author: Bartuccio Antoine
* @Date:   2018-09-19 22:18:46
* @Last Modified by:   Bartuccio Antoine
* @Last Modified time: 2018-11-25 02:25:36
*/

Vue.component('media-card', {
	props: ['media', "mediaType"],
	template: `<div class="media-card card border-dark mb-3 flex-row flex-wrap" :class="{ 'bg-warning': media.suggested }">

				<div class="card-header-border-0">
					<img class="card-img-left rounded" :alt="media.title" :src="media.image">
				</div>
				<div class="card-block px-2">
					<h5 class="card-title">{{ media.title }}</h5>
					<div class="card-text">
						<p>Release: {{ media.release_date }}</p>
						<p>Genres: {{ media.genres }}</p>
						<button class="btn btn-outline-warning my-2 my-sm-0" v-on:click="app.discoverMedia(media, media.media_type)">Montrer du contenu relatif</button>
					</div>
				</div>
				</div>`
})

var mediaTypes = {
	"Movies": 0,
	"Books": 1,
	"GraphicNovels": 2,
	"VideoGames": 3,
	"Series": 4
}

function searchQuery(app, type){
	query = {
		media_type: mediaTypes[type],
		title: app.searchedTitle,
		page: app.page
	}
	if (Number(app.searchedYear) > 0) {
		query.year = Number(app.searchedYear)
	}
	axios.post('/search', query).then(function(response) {
		console.log(response)
		app.medias[type] = []
		for (var i = response.data.medias.length - 1; i >= 0; i--) {
			response.data.medias[i].media_type = mediaTypes[type]
			app.medias[type].push(response.data.medias[i])
		}
	}).catch(function(error) {
		console.log(error)
	})
}

var app = new Vue({
	el: '#app',
	data: {
		navbarActive: {
			"isMovies": true,
			"isSeries": false,
			"isBooks": false,
			"isBd": false,
			"isVideoGames": false,
			"isDiscover": false
		},
		medias: {
			"Movies": [],
			"Books": [],
			"GraphicNovels": [],
			"VideoGames": [],
			"Series": [],
			"Discovers": []
		},
		searchedTitle: "",
		searchedYear: "",
		advancedSearch: false,
		page: 1
	},
	methods: {
		navbarSelect: function(category) {
			for (var key in this.navbarActive){
				this.navbarActive[key] = false
			}
			this.navbarActive[category] = true
		},
		searchMedia: function(resetPagination) {
			if (resetPagination){
				this.page = 1
			}
			if (app.navbarActive["isDiscover"]){
				app.navbarSelect('isMovies')
			}
			for (var key in mediaTypes){
				searchQuery(this, key)
			}
		},
		showAdvancedSearch: function() {
			this.advancedSearch = !this.advancedSearch
		},
		discoverMedia: function(media, media_type) {
			app.navbarSelect('isDiscover')
			query = {
				media_type: media_type,
				media: media
			}
			axios.post('/discover', query).then(function(response) {
				console.log(response)
				app.medias["Discovers"] = response.data.medias
			}).catch(function(error) {
				console.log(error)
			})
		},
		goToPage: function(page) {
			if (page > 0 ){
				this.page = page
				this.searchMedia(false)
			}
		}
	}
})
