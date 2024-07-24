import './App.css';
import { OpenDirectoryDialog } from '../wailsjs/go/main/App';
import { FilesInDirectoryHandler, JsPrintLn, SearchShow, RenameAll,  GetEpisodes, RenameSelected  } from "../wailsjs/go/main/App";
import {WindowSetTitle} from "../wailsjs/runtime/runtime"

let EpList
let userDir
let episodeList
let showList
let selectFromUserFiles = ""
let selectFromApi = ""
let rngcolor = "264653"

//Setting windows tittle
WindowSetTitle("Easy Media Renamer")

//html rendering
function App() {

    return (
        <div id="App">
            <div>

            <div className="rename-btn-section">
                    <button className="btn-chooser" onClick={openDirectory}>Directory Chooser</button>
                    <button className="btn-chooser" onClick={CallRenameAllGo}>Rename all</button>
                    <button className="btn-chooser" onClick={renameSelected}>Rename selected</button>
            </div>

            <div className="input-wrapper">
                <input className="user-input-search" type="text" placeholder="Show name"></input>
                <button className="srch-btn" onClick={CallSeacrhShowGo}><i className="fa-solid fa-magnifying-glass"></i></button>
            </div>

            <div className="search-modal">
                <ul className="list-search">
                </ul>
            </div>
            
            <div className="main-content-wrapper">
                    <div className="user-files">
                    <div>Original files</div>
                        <ul className="list-elem">
                        </ul>
                    </div>

                <div className="arrow-gui-pointer"><i className="fa-solid fa-arrow-right"></i></div>
                    <div className="rename-files">
                    <div>New files</div>
                        <ul className="list-elem-re">
                        </ul>
                        
                    </div>
            </div>
            </div>
        </div>
    )
}


async function renameSelected(){
    if (selectFromApi&&selectFromUserFiles){
        await RenameSelected(selectFromUserFiles, selectFromApi)
        let file_name_list = await FilesInDirectoryHandler(userDir)
        await renderFileList(file_name_list)
    }
}

//Open directory dialog
async function openDirectory() {
    try {
        //this will basicly give us directoy path
        userDir = await OpenDirectoryDialog();
        //getting file list in directory
        let file_name_list = await FilesInDirectoryHandler(userDir)
        //rendering file list in our html
        await renderFileList(file_name_list)
    } catch (error) {
        console.error("Failed to open directory dialog:", error);
    }
}
//Calling an api
async function CallSeacrhShowGo(){

    let showTitle = document.querySelector(".user-input-search")

    if (showTitle.value) {
        showList = await SearchShow(showTitle.value)
        
        await renderShowList(showList)
    }
    
}

async function CallRenameAllGo(){
    
   
        await RenameAll()
        if (userDir){
            let file_name_list = await FilesInDirectoryHandler(userDir)
            //rendering file list in our html
            await renderFileList(file_name_list)
        }
        
}

async function renderFileList(list){

    fileList = document.querySelector(".list-elem")
    fileList.innerHTML = ""
    
    for (let i = 0; i < list.length; i++) {
        let liEl = document.createElement("li")
        liEl.classList.add("user-file")
        liEl.append(list[i])
        fileList.append(liEl)
      }
    
      await addListenersForFiles(".user-file")
}

async function renderEpList(showid){

    let episodelist = await GetEpisodes(parseInt(showid))
    let EpList = document.querySelector(".list-elem-re")
    EpList.innerHTML = ""
    
    for (let i = 0; i < episodelist.length; i++) {
        let liEl = document.createElement("li")
        liEl.classList.add("file-from-api")
        liEl.append(i+1+". "+episodelist[i])
        EpList.append(liEl)
      }

      await addListenersForFiles(".file-from-api")
}

async function renderShowList(list){
    let EpList = document.querySelector(".list-search")
    EpList.style.display = "flex"
    EpList.style.border = "solid"
    EpList.innerHTML = ""
 
    for (let i = 0; i < list.length; i++) {
        let liEl = document.createElement("li")
        liEl.classList.add("show-from-api")
        liEl.classList.add(list[i].show.id)
        liEl.append(list[i].show.name)
        EpList.append(liEl)
      }

      await addListenersForFiles(".show-from-api")
}


async function addListenersForFiles(selector){
    let EpList = document.querySelector(".list-search")
    let userElements = document.querySelectorAll(selector)
    
    for (let i = 0; i < userElements.length; i++) {

        userElements[i].addEventListener("click", (e)=>{

            //beggining of the search, getting an id of a show
            //passing to the render episode func
            let ElemenClassList = e.target.classList
            if (selector == ".show-from-api"){
                let showid = ElemenClassList[1]
                EpList.style.display = "none"
                EpList.style.border = "none"
                renderEpList(showid)
            
            }else if (ElemenClassList[1] == ("selected")&&ElemenClassList[0] == "file-from-api"){
                ElemenClassList.remove("selected")
                e.target.style.backgroundColor = ""
                selectFromApi = null
            }else if (ElemenClassList[1] == ("selected")&&ElemenClassList[0] == "user-file"){
                ElemenClassList.remove("selected")
                e.target.style.backgroundColor = ""
                selectFromUserFiles = null
            }
        

            else { 
                    //no more then 2 items
                    if (selectFromApi&&selectFromUserFiles){
                        let sellist = document.querySelectorAll(".selected")
                        for (let i = 0; i < sellist.length; i++) {
                            sellist[i].classList.remove("selected")
                            sellist[i].style.backgroundColor = ""
                            rngcolor = Math.floor(100000 + Math.random() * 900000);
                            selectFromApi = null
                            selectFromUserFiles = null
                        }
                    }
            } 

            //making sure its not possible to choose multiple items from a same column
            if (ElemenClassList[0] == "file-from-api"&&!selectFromApi){
               selectFromApi = e.target.textContent
               ElemenClassList.add("selected")
               e.target.style.backgroundColor = "#"+rngcolor
            }

            if (ElemenClassList[0] == "user-file"&&!selectFromUserFiles){
                selectFromUserFiles = e.target.textContent
                ElemenClassList.add("selected")
                e.target.style.backgroundColor = "#"+rngcolor
            }
})
        
      }
    
}



export default App
