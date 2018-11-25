/*
* @Author: Bartuccio Antoine
* @Date:   2018-09-16 18:09:47
* @Last Modified by:   Bartuccio Antoine
* @Last Modified time: 2018-11-22 23:03:38
 */

package media

import (
	"fmt"
	"time"
)

// MediaType enum specifying the media type
type MediaType int

// JSONTime wrapper around time.Time to have a custom marshaller
type JSONTime struct {
	time.Time
}

// MarshalJSON custom marshaller for the time
func (t *JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.Format("02 January 2006"))), nil
}

// UnmarshalJSON custom unmarshaller for the time
func (t *JSONTime) UnmarshalJSON(data []byte) error {
	ut, err := time.Parse("02 January 2006", string(data))
	if err == nil {
		*t = JSONTime{ut}
	}
	// I had some issues with sessions so I'm gonna pass around this error
	*t = JSONTime{time.Now()}
	return nil
}

const (
	Movie        MediaType = 0
	Book         MediaType = 1
	GraphicNovel MediaType = 2
	VideoGame    MediaType = 3
	Serie        MediaType = 4
)

// Media rerpesent a media
type Media struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	ReleaseDate JSONTime `json:"release_date"`
	Genres      []string `json:"genres"`
	Image       string   `json:"image"`
	Adult       bool     `json:"adult"`
	Suggested   bool     `json:"suggested"`
}
