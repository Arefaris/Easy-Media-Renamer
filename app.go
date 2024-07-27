package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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
}

// NewApp creates a new App application struct
func NewApp() *App {
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


func (a *App) SearchShow(show string) ([]Show){
	url := "https://api.tvmaze.com/search/shows?q="+show 
	
	
	
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Ошибка при чтении ответа: %v", err)
	}
	

	var shows []Show
	if err := json.Unmarshal(body, &shows); err != nil {
		log.Fatalf("Ошибка при декодировании JSON: %v", err)
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
	
		episodeList = append(episodeList, cleanEpisode)
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

func (a *App)RenameAllGO(fileNamelist[]string) {
	
	if fileNamelist != nil && episodeList != nil {
		for index, file := range fileNamelist {
			ext := filepath.Ext(file)
			episodeNumber := index+1
			enS := strconv.Itoa(episodeNumber)
			
			newfile := userDIR+"\\"+enS+". "+episodeList[index]+ext
			oldfile := userDIR+"\\"+fileNamelist[index]
			
			e := os.Rename(oldfile, newfile) 
	
			if e != nil { 
				log.Fatal(e) 
			}
			fmt.Println(oldfile, "changed to: ", newfile)
		}
	}
	
}

func (a *App)RenameSelectedGO(originalFileName string, newFileName string){

	originalFilePath := userDIR+"\\"+originalFileName
	ext := filepath.Ext(originalFilePath)
	newFilePath := userDIR+"\\"+newFileName+ext
	
	e := os.Rename(originalFilePath, newFilePath) 
	
			if e != nil { 
				log.Fatal(e) 
			} 
}