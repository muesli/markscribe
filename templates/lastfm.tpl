### Hi there 👋

## Favourite albums of all time 🎶

{{range lastfmFavouriteAlbums 5}}
- {{.Artist.Name}} - {{.Name}}
{{- end}}


## Favourite artists of all time 👨‍🎤

{{range lastfmFavouriteArtists 5}}
- {{.Name}} ({{.PlayCount}})
{{- end}}


## Favourite tracks of all time 💿

{{range lastfmFavouriteTracks 5}}
- {{.Artist.Name}} - {{.Name}} ({{.PlayCount}})
{{- end}}


## Most recent tracks 🎺

{{range lastfmRecentTracks 10}}
- {{.Artist.Name}} - {{.Name}}
{{- end}}
