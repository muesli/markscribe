### Hi there 👋

## Favourite albums of all time 🎶

{{range lastFmFavouriteAlbums 5}}
- {{.Artist.Name}} - {{.Name}}
{{- end}}


## Favourite artists of all time 👨‍🎤

{{range lastFmFavouriteArtists 5}}
- {{.Name}} ({{.PlayCount}})
{{- end}}


## Favourite tracks of all time 💿

{{range lastFmFavouriteTracks 5}}
- {{.Artist.Name}} - {{.Name}} ({{.PlayCount}})
{{- end}}


## Most recent tracks 🎺

{{range lastFmRecentTracks 10}}
- {{.Artist.Name}} - {{.Name}}
{{- end}}
