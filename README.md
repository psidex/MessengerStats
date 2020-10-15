# MessengerStats

View statistics about your Messenger conversations

## Example

(Names censored)

![example](example.png)

## Architecture

The app is a basic http web server but, for obvious reasons, isn't hosted on the internet. To use it, download and run
one of the pre-built binaries (or build it yourself) and then go to the URL it gives you.

It uses basic HTML forms for data transfer, and the calculations for even lots of files shouldn't take very long; on my
machine it takes ~600ms to upload, parse, and calculate statistics for 20.8MB split over 11 files. (If that seems slow,
most of that time is parsing the files, the actual statistics calculations currently take ~2ms per MB)
