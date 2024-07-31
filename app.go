package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"strconv"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.felesatra.moe/anidb"
	"apikeys"
)



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


var clientAniDb anidb.Client
func loadENV(){
	apikeys.LoadApiKeys("123")
  }


func loadClients(){
	clientAniDb = anidb.Client{
		Name: os.Getenv("ANI_DB"), 
		Version: 1,
	}
}


// NewApp creates a new App application struct
func NewApp() *App {
  	loadENV()
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


func (a *App) SearchShow(show string, apitype string) ([]Show){
	shows := []Show{}

	//anime, err := clientAniDb.RequestAnime(17635)
	
	
	if apitype == "TVmaze"{
		shows = tvmazeApi(show)
	}else if apitype == "AniDB"{
		shows = aniDbApi(show)
	}
	return shows
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
                fmt.Printf("Found anime: %s (AID: %d)\n", title.Name, anime.AID)
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

		if (episode.Season < 10){
			epseason = "0"+epseason
		}
		
		
		
		episodeList = append(episodeList, cleanEpisode+" - s"+epseason+"e"+epnumber)
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

		for index, file := range fileNamelist {
			ext := filepath.Ext(file)
			newfile := userDIR+"\\"+episodeList[index]+ext
			oldfile := userDIR+"\\"+fileNamelist[index]

			//checking if name is already exist in a folder, if not, rename
			if !stringInSlice(episodeList[index]+ext, file_names_list){
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
	
	originalFilePath := userDIR+"\\"+originalFileName
	ext := filepath.Ext(originalFilePath)
	newFilePath := userDIR+"\\"+newFileName+ext
	
	
	
	if !stringInSlice(newFileName+ext, file_names_list){
		e := os.Rename(originalFilePath, newFilePath)
		if e != nil { 
			log.Fatal(e)
		}
	}
	
}