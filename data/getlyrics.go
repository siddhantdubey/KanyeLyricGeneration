package main

import (
	"context"
	"fmt"
	"log"
	"os"

	spotifyauth "github.com/zmb3/spotify/v2/auth"

	godotenv "github.com/joho/godotenv"
	lyrics "github.com/rhnvrm/lyric-api-go"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

func getArtistSongs(artist string, client spotify.Client) []string {
	results, err := client.Search(artist, spotify.SearchTypeArtist)
	if err != nil {
		panic(err)
	}
	var artist_id spotify.ID
	var songs []string
	var songids []spotify.ID
	if results.Artists != nil {
		for _, item := range results.Artists.Artists {
			if item.Name == artist {
				artist_id = item.ID
				fmt.Println(item.Name)
				break
			}
		}
	}
	//get artist's songs
	if artist_id != "" {
		artist_albums, err := client.GetArtistAlbums(artist_id)
		if err != nil {
			panic(err)
		}

		var album_ids []spotify.ID
		for _, album := range artist_albums.Albums {
			album_ids = append(album_ids, album.ID)
		}
		//get albums' songs
		if len(album_ids) > 0 {
			for _, album_id := range album_ids {
				album_tracks, err := client.GetAlbumTracks(album_id)
				if err != nil {
					panic(err)
				}
				for _, track := range album_tracks.Tracks {
					//if song isn't in list, add it
					if !contains(songids, track.ID) {
						fmt.Printf("%s\n", track.Name)
						songs = append(songs, track.Name)
						songids = append(songids, track.ID)
					}
				}
			}
		}
	}
	return songs
}

func contains(s []spotify.ID, e spotify.ID) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {
	var err = godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	var clientId = os.Getenv("CLIENT_ID")
	var clientSecret = os.Getenv("CLIENT_SECRET")

	var ctx = context.Background()
	var config = &clientcredentials.Config{
		//get from .env file
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	var token, loginErr = config.Token(ctx)
	if loginErr != nil {
		panic(loginErr)
	}
	var httpClient = spotifyauth.New().Client(ctx, token)
	var client = spotify.NewClient(httpClient)

	artist := "Kanye West"
	song_list := getArtistSongs(artist, client)

	f, err := os.Create("lyrics/alllyrics.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for _, song := range song_list {
		fmt.Println(song)
		l := lyrics.New()

		lyric, err := l.Search(artist, song)
		if err != nil {
			fmt.Printf("Lyrics for %v-%v were not found", artist, song)
		}
		fmt.Println(lyric)
		_, err2 := f.WriteString(lyric)
		if err2 != nil {
			log.Fatal(err2)
		}
	}
}
