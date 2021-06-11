package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type Gpx struct {
	Metadata Metadata    `xml:"metadata"`
	Trk      Trk         `xml:"trk"`
	XMLName  interface{} `xml:"gpx"`
}

type Metadata struct {
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
	Time    string      `xml:"time"`
	XMLName interface{} `xml:"trkpt"`
}

func makeGpx(positions []*Position, vehicle string) (result []byte, matchOk bool, err error) {
	if len(positions) == 0 {
		err = errors.New("empty positions in status")
		return
	}

	now := time.Now()
	gpx := Gpx{
		Metadata: Metadata{
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
			Time: el.when.Format(time.RFC3339),
			Lat:  el.latitude,
			Lon:  el.longitude,
		})
	}

	var raw []byte
	if raw, err = xml.MarshalIndent(gpx, "", " "); err != nil {
		return
	}
	if vehicle != "" {
		if result, err = mapMatch(raw, vehicle); err != nil {
			result = raw
		} else {
			matchOk = true
		}
	}
	return
}

func mapMatch(data []byte, vehicle string) (result []byte, err error) {
	var resp *http.Response
	if resp, err = http.Post(graphHopperUrl+"/match?vehicle="+vehicle+"&type=gpx", "application/gpx+xml", bytes.NewBuffer(data)); err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result, err = ioutil.ReadAll(resp.Body)
	} else {
		err = errors.New("graphhopper request outcome status: " + resp.Status)
	}
	return
}
