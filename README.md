# MessengerStats

[![Build status](https://github.com/psidex/MessengerStats/workflows/CI/badge.svg)](https://github.com/psidex/MessengerStats/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/psidex/MessengerStats)](https://goreportcard.com/report/github.com/psidex/MessengerStats)
[![license](https://img.shields.io/github/license/psidex/MessengerStats.svg)](./LICENSE)
[![Ko-fi donate link](https://img.shields.io/badge/Support%20Me-Ko--fi-orange.svg?style=flat&colorA=35383d)](https://ko-fi.com/M4M18XB1)

View statistics about your Messenger conversations

## Example

(Names censored)

![example](example.png)

## Architecture

The app is a basic http web server but, for obvious reasons, isn't hosted on the internet. To use it, download and run
one of the pre-built binaries (or build it yourself) and then go to the URL it gives you.

It uses basic HTML forms for data transfer, and the calculations for even lots of files shouldn't take very long; on my
machine it takes ~200ms to upload, parse, and calculate statistics for a 20.8 MB conversation split over 11 files
(Almost all of that time is from the calls to `ParseMultipartForm` and `UnmarshalJSON`, the actual statistics
calculations currently take ~2ms per MB).

## Credits

- Main website CSS: [water.css](https://watercss.kognise.dev/)
- Charting library: [Highcharts](https://www.highcharts.com/)
- Nice error notifications: [notyf](https://github.com/caroso1222/notyf)
- Favicon icon: [Feather Icons](https://feathericons.com/)
