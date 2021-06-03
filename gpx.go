package main

import (
	"encoding/xml"
	"errors"
	"time"
)

type Gpx struct {
	Metadata Metadata    `xml:"metadata"`
	Trk      Trk         `xml:"trk"`
	XMLName  interface{} `xml:"gpx"`
}

type Metadata struct {
	Link Link   `xml:"link"`
	Time string `xml:"time"`
}

type Link struct {
	Text string `xml:"text"`
}

type Trk struct {
	Name   string `xml:"name"`
	Trkseg Trkseg `xml:"trkseg"`
}

type Trkseg struct {
	Trkpt []TrkptEl
}

type TrkptEl struct {
	Lat     float64     `xml:"lat,attr"`
	Lon     float64     `xml:"lon,attr"`
	Ele     float64     `xml:"ele"`
	Time    string      `xml:"time"`
	XMLName interface{} `xml:"trkpt"`
}

func makeGpx(positions []*Position) (result []byte, err error) {
	if len(positions) == 0 {
		err = errors.New("empty positions in status")
		return
	}

	now := time.Now()
	gpx := Gpx{
		Metadata: Metadata{
			Link: Link{
				Text: "Garmin International",
			},
			Time: now.Format(time.RFC3339),
		},
		Trk: Trk{
			Name: "atrack",
			Trkseg: Trkseg{
				Trkpt: []TrkptEl{},
			},
		},
	}

	for _, el := range positions {
		gpx.Trk.Trkseg.Trkpt = append(gpx.Trk.Trkseg.Trkpt, TrkptEl{
			Ele:  0.,
			Time: el.when.Format(time.RFC3339),
			Lat:  el.latitude,
			Lon:  el.longitude,
		})
	}

	result, err = xml.MarshalIndent(gpx, "", " ")
	return
}
