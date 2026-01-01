# packrat #

Packrat is a simple go app for downloading podcast audio to a local drive and providing an RSS feed to access it from the local network.

## Configuration ##

These environment variables need to be set:
- `LISTEN_HOST` the host that the server will listen on (use `0.0.0.0` to listen on all interfaces)
- `LISTEN_PORT` the tcp port to use
- `STATIC_DIR` the directory where podcast audio files and RSS xml files will be saved locally
- `SERVER_HOST` the hostname or IP address of the system hosting the server. This is the host that will be used in the feed data. May be different from `LISTEN_HOST` e.g. if running on docker.

## Usage ##

Podcast feeds and episodes can be accessed at `http://{server}/feeds`

To add a new podcast or check for new episodes of an existing one, send a POST request to the packrat server with the URL of the podcast RSS feed:
`curl -X POST http://127.0.0.1:5000/api/feeds -H "Content-Type: application/json" -d '{"url":"https://feeds.simplecast.com/wjQvYtdl"}'
`

The `docker-compose.yml` file is written for deploying the app on a Synology NAS, which is how I'm using it.

## Improvement ideas ##

- Provide a way to see number of pending downloads
- Cron job to automatically update existing feeds on some regular cadence
- DB to save data about what feeds have been added, etc.