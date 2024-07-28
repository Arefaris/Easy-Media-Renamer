import Button from '@mui/material/Button';
import { OpenDirectoryDialog } from '../wailsjs/go/main/App';
import { FilesInDirectoryHandlerGO, RenameAllGO,  GetEpisodesGO, RenameSelectedGO  } from "../wailsjs/go/main/App";

//user folder
let userDir

//selected items from both lists (for selective renameming)
let selectFromUserFiles = ""
let selectFromApi = ""

//default color for selecting item in a list
let rngcolor = "264653"

function RenameBtnSc(){
    return (
        <div className="rename-btn-section">
                    <Button variant="contained" className="btn-chooser" onClick={openDirectory}>Directory Chooser</Button>
                    <Button variant="contained" className="btn-chooser" onClick={CallRenameAllGo}>Rename all</Button>
                    <Button variant="contained" className="btn-chooser" onClick={renameSelected}>Rename selected</Button>
        </div>
    )
}


//Open directory dialog
async function openDirectory() {
    try {
        //this will basicly give us directoy path
        userDir = await OpenDirectoryDialog();
        //getting file list in directory
        let file_name_list = await FilesInDirectoryHandlerGO(userDir)
        //rendering file list in our html
        await renderList(file_name_list)
    } catch (error) {
        console.error("Failed to open directory dialog:", error);
    }
}


async function CallRenameAllGo(){
    //taking Rearranged list  
    let fileList = document.querySelector(".list-elem")
    let fileListReorder = []
    for (let i = 0; i < fileList.childNodes.length; i++) {
        fileListReorder.push(fileList.childNodes[i].textContent)
      }
    
    await RenameAllGO(fileListReorder)
    if (userDir){
        let file_name_list = await FilesInDirectoryHandlerGO(userDir)
        //rendering file list in our html
        await renderList(file_name_list)
    }
    
}

async function renameSelected(){
    if (selectFromApi&&selectFromUserFiles){
        await RenameSelectedGO(selectFromUserFiles, selectFromApi)
        let file_name_list = await FilesInDirectoryHandlerGO(userDir)
        await renderList(file_name_list)
        selectFromUserFiles = null
    }
}

//render user files to html
async function renderList(list, showId=null){
    let fileList;
    let itemclass;

    if(showId){
        list = await GetEpisodesGO(parseInt(showId))
        fileList = document.querySelector(".list-elem-re")
        fileList.innerHTML = ""
        itemclass = "file-from-api"
       }else{
        fileList = document.querySelector(".list-elem")
        fileList.innerHTML = ""
        itemclass = "user-file"
       }
    
    
    for (let i = 0; i < list.length; i++) {
        let liEl = document.createElement("li")
        liEl.classList.add(itemclass)
        liEl.append(list[i])
        fileList.append(liEl)
      }
    
      await addListenersForFiles("."+itemclass)
}



export async function addListenersForFiles(selector){
    let EpList = document.querySelector(".list-search")
    let userElements = document.querySelectorAll(selector)
    
    for (let i = 0; i < userElements.length; i++) {

        userElements[i].addEventListener("click", (e)=>{

            //beginning of the search, getting an id of a show
            //passing to the render episode func
            let ElemenClassList = e.target.classList
            if (selector == ".show-from-api"){
                let showid = ElemenClassList[1]
                EpList.style.display = "none"
                EpList.style.border = "none"
                renderList(null,showid)
            
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




export default RenameBtnSc