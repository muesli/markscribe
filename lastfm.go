package main

import (
	"github.com/shkh/lastfm-go/lastfm"
)

func lastfmFavouriteAlbums(count int) interface{} {

	params := lastfm.P{
		"user":  lastfmUser,
		"limit": count,
	}
	albums, err := lastfmApi.User.GetTopAlbums(params)
	if err != nil {
		panic(err)
	}
	return albums.Albums
}

func lastfmFavouriteTracks(count int) interface{} {
	params := lastfm.P{
		"user":  lastfmUser,
		"limit": count,
	}
	tracks, err := lastfmApi.User.GetTopTracks(params)
	if err != nil {
		panic(err)
	}
	return tracks.Tracks
}

func lastfmFavouriteArtists(count int) interface{} {
	params := lastfm.P{
		"user":  lastfmUser,
		"limit": count,
	}
	artists, err := lastfmApi.User.GetTopArtists(params)
	if err != nil {
		panic(err)
	}
	return artists.Artists
}

func lastfmRecentTracks(count int) interface{} {
	params := lastfm.P{
		"user":  lastfmUser,
		"limit": count,
	}
	tracks, err := lastfmApi.User.GetRecentTracks(params)
	if err != nil {
		panic(err)
	}
	return tracks.Tracks
}
