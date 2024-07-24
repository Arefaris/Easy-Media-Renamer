import './App.css';
import { OpenDirectoryDialog } from '../wailsjs/go/main/App';
import { FilesInDirectoryHandler, JsPrintLn, SearchShow, RenameAll,  GetEpisodes  } from "../wailsjs/go/main/App";
import {WindowSetTitle} from "../wailsjs/runtime/runtime"

let fileList
let EpList
let result
let episodeList
let showList

addEventListener("DOMContentLoaded", (event) => {});
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
                    <button className="btn-chooser" onClick={openDirectory}>Rename selected</button>
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

//Open directory dialog
async function openDirectory() {
    try {
        //this will basicly give us directoy path
        result = await OpenDirectoryDialog();
        //getting file list in directory
        let file_name_list = await FilesInDirectoryHandler(result)
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
        if (result){
            let file_name_list = await FilesInDirectoryHandler(result)
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
        liEl.append(i+1+". "+list[i])
        fileList.append(liEl)
      }
    
      await addListenersForFiles(".list-search")
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


            if (selector == ".show-from-api"){
                
                let showid = e.target.classList[1]
                EpList.style.display = "none"
                EpList.style.border = "none"
                renderEpList(showid)
                
            }else if (e.target.classList[1] == ("selected")) {
                e.target.style.backgroundColor = ""
                e.target.classList.remove("selected")
        
            }else {
                e.target.style.backgroundColor = "#264653"
                e.target.classList.add("selected")
            }


    })
        
      }
    
}



export default App
