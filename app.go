package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.felesatra.moe/anidb"
	tmdb "github.com/cyruzin/golang-tmdb"
)

//go:embed keys.json
var jsonData []byte

// App struct
type App struct {
	ctx context.Context
}


type Show struct {
	Score float64 `json:"score"`
	Show  struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Genre []string `json:"genres"`

	} `json:"show"`
}

type Episode struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Season int    `json:"season"`
	Number int `json:"number"`
}

type Api struct {
	NameAniDb string `json:"anidb"`
	Anidbv int `json:"anidbv"`
	Tmdb string `json:"tmdb"`
}


var clientAniDb anidb.Client
var api Api

func loadClients(){
	//opening our json

	json.Unmarshal(jsonData, &api)

	clientAniDb = anidb.Client{
		Name: api.NameAniDb, 
		Version: api.Anidbv,
	}
	

}


// NewApp creates a new App application struct
func NewApp() *App {
  	loadClients()
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}




func (a *App) OpenDirectoryDialog() (string, error) {
    return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
        Title: "Select Directory",
    })

}

var apitype string
func (a *App) SearchShow(show string, apis string) ([]Show){
	shows := []Show{}

	apitype = apis
	if apitype == "TVmaze"{
		shows = tvmazeApi(show)
	}else if apitype == "AniDB"{
		shows = aniDbApi(show)
		
	}else if apitype == "TMDB"{
		shows = tmdDbApi(show)
	}	
	return shows
}

func tmdDbApi(show string)([]Show){

	apishow := Show{}
	apishows := []Show{}

	tmdbClient, err := tmdb.Init(api.Tmdb)

	if err != nil {
		fmt.Println(err)
	}

	options := make(map[string]string)
	options["language"] = "EN"

	// Multi Search
	search, err := tmdbClient.GetSearchMulti(show, options)

	if err != nil {
		fmt.Println(err)
	}

	
	// Iterate
	for _, v := range search.Results {

		if v.MediaType == "movie" {
			
			if (len(v.Title) > 0){
				
				apishow.Show.ID = int(v.ID)
				apishow.Show.Name = v.Title 
				apishows = append(apishows, apishow)
				
			}
			
		} 
		 if v.MediaType == "tv" {
				
				
				apishow.Show.ID = int(v.ID)
				apishow.Show.Name = v.Name 
				apishows = append(apishows, apishow)
			
			
		} 
	}
	
	return apishows
}

func aniDbApi(show string)([]Show){
	
	apishow := Show{}
	apishows := []Show{}

	c, err := anidb.DefaultTitlesCache()
	if err != nil {
		panic(err)
	}
	defer c.SaveIfUpdated()

	titles, err := c.GetTitles()
	if err != nil {
		panic(err)
	}

	searchTerm := show
    for _, anime := range titles {
        for _, title := range anime.Titles {
            if strings.Contains(strings.ToLower(title.Name), strings.ToLower(searchTerm)) {
				apishow.Show.ID = anime.AID
				apishow.Show.Name = title.Name
				apishows = append(apishows, apishow)
            }
        }
    }
	return apishows

}

func tvmazeApi(show string)([]Show){
	url := "https://api.tvmaze.com/search/shows?q="+show 
	
	
	
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("%v", err)
	}
	

	var shows []Show
	if err := json.Unmarshal(body, &shows); err != nil {
		log.Fatalf("%v", err)
	}
	
	return shows
}




var episodes []Episode
var episodeList = []string{}

func (a *App) GetEpisodesGO(showID int) []string{
	
	if apitype == "TVmaze"{
		episodeList = a.EpisodesTvMaze(showID)
	}else if apitype == "AniDB"{
		episodeList = a.EpisodesAniDb(showID)
	}else if apitype == "TMDB"{
		episodeList = a.EpisodesTmDb(showID)
	}
	
	return episodeList
}

func (a *App) EpisodesTmDb(showId int) []string{
	episodeList = nil
	tmdbClient, err := tmdb.Init(api.Tmdb)

	
	if err != nil {
		fmt.Println(err)

	}
	

	options := make(map[string]string)
	options["language"] = "EN"
	
	details, _ := tmdbClient.GetTVDetails(showId, options)
	
	
	// Перебираем все сезоны и эпизоды
	for _, season := range details.Seasons {
		episodes, err := tmdbClient.GetTVSeasonDetails(showId, season.SeasonNumber, nil)
		if err != nil {
			log.Fatalf("Error getting season details: %v", err)
		}

		
		for _, episode := range episodes.Episodes {
			epseason := strconv.Itoa(season.SeasonNumber)
			epnumber := strconv.Itoa(episode.EpisodeNumber)
			cleanEpisode := a.cleanName(episode.Name)

			if (episode.EpisodeNumber < 10){
				epnumber = "0"+epnumber
			}
	
			if (episode.SeasonNumber < 10){
				epseason = "0"+epseason
			}
			episodeList = append(episodeList, cleanEpisode+" - s"+epseason+"e"+epnumber)
		}
	}
	return episodeList

}


func (a *App) EpisodesTvMaze(showID int) []string{ 
	episodeList = nil
	url := fmt.Sprintf("https://api.tvmaze.com/shows/%d/episodes", showID)
	
	resp, err := http.Get(url)

	if err != nil {
		log.Fatalf("%v", err)
	}
	defer resp.Body.Close()
	
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%v", err)
	}

	if err := json.Unmarshal(body, &episodes); err != nil {
		log.Printf("%v", err)
	}
	
	
	for _, episode := range episodes{
		cleanEpisode := a.cleanName(episode.Name)
		epnumber := strconv.Itoa(episode.Number)
		epseason := strconv.Itoa(episode.Season)

		if (episode.Number < 10){
			epnumber = "0"+epnumber
		}

		if (episode.Season < 10 && episode.Season != 0){
			epseason = "0"+epseason
		}
		
		
		
		episodeList = append(episodeList, cleanEpisode+" - s"+epseason+"e"+epnumber)
	}
	
	return episodeList

}

func (a *App) EpisodesAniDb(showID int) []string{

		anime, err := clientAniDb.RequestAnime(showID)
		
		if err != nil {
			panic(err)
		}
		namesNepisodes := map[int]string{}
		for _, Episode := range anime.Episodes {
			
			
			for _, title := range Episode.Titles {
				if (title.Lang == "en"){
					n , _ := strconv.Atoi(Episode.EpNo)
					
					
					namesNepisodes[n] = title.Title
					
			}
				}

				sortedList := []string{}
				for i := 1; i < len(namesNepisodes); i++ {
					epname := a.cleanName(namesNepisodes[i])
					
			
					epnumber := strconv.Itoa(i)

					if (i < 10){
						epnumber = "0"+epnumber
					}

					sortedList = append(sortedList, epname + " - " + epnumber)
				}
				
				episodeList = sortedList
				
		}
		return episodeList
}

//clean name from special characters.
func (a *App) cleanName(name string)string{
	specialChars := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*", "?"}
	for _, char := range specialChars {
		if strings.Contains(name, char) {
			name = strings.Replace(name, char, "", -1)
		}
	}
	return name
}

var file_names_list = []string{}
var file_path_list  = []string{}
var userDIR string

func (a *App) FilesInDirectoryHandlerGO(directory string)[]string{
	file_names_list = nil
	userDIR = directory
	files, err := os.ReadDir(directory)

    if err != nil {
        println(err)
    }
    
	for _, file := range files {
		file_names_list = append(file_names_list, file.Name())
		file_path_list = append(file_path_list, filepath.Join(directory, file.Name()))
	}
	
	return file_names_list

}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func (a *App)RenameAllGO(fileNamelist[]string) {
	if len(fileNamelist) > 0 && len(episodeList) > 0 {
		fmt.Println(episodeList)
		for i := 0; i < len(episodeList)&&i < len(fileNamelist); i++ {
			ext := filepath.Ext(fileNamelist[i])
			newfile := filepath.Join(userDIR, episodeList[i]+ext)
			oldfile := filepath.Join(userDIR, fileNamelist[i])

			//checking if name is already exist in a folder, if not, rename
			if !stringInSlice(episodeList[i]+ext, file_names_list){
				e := os.Rename(oldfile, newfile)
				if e != nil { 
					log.Fatal(e)
				}
				fmt.Println(oldfile, "changed to: ", newfile)
		}
	
		}
	}
	
}

func (a *App)RenameSelectedGO(originalFileName string, newFileName string){
	originalFilePath := filepath.Join(userDIR, originalFileName)
	ext := filepath.Ext(originalFilePath)
	newFilePath := filepath.Join(userDIR, newFileName+ext)
	
	
	
	if !stringInSlice(newFileName+ext, file_names_list){
		e := os.Rename(originalFilePath, newFilePath)
		if e != nil { 
			log.Fatal(e)
		}
	}
	
}