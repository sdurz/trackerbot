package main

const helpMarkup = `
Bot commands:
*/getgpx* - get current gps data as .gpx file (if available)
*/pause* - pause tracking
*/resume* - resume tracking
*/end* - end tracking (and get .gpx file)
*/setprofile* - sets map mapping profile (car, bike or hike)

_*NOTE*_ only _hike_ profile is actually useful, Telegram clients update shared positions at an insufficient rate for higher speeds.
`
