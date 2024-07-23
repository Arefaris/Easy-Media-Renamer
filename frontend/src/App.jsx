import './App.css';
import { OpenDirectoryDialog } from '../wailsjs/go/main/App';
import { FilesInDirectoryHandler, JsPrintLn, SearchShow, RenameAll } from "../wailsjs/go/main/App";
let EpList = document.querySelector(".list-elem-re")
let result

function App() {

    return (
        <div id="App">
            <div>
            <div className="rename-btn-section">
                    <button className="btn-chooser" onClick={openDirectory}>Directory Chooser</button>
                    <button className="btn-chooser" onClick={RenameAll}>Rename all</button>
                    <button className="btn-chooser" onClick={openDirectory}>Rename selected</button>
            </div>

            <div className="input-wrapper">
                <input className="user-input-search" type="text" placeholder="Show name"></input>
                <button className="srch-btn" onClick={CallSeacrhShowGo}><i className="fa-solid fa-magnifying-glass"></i></button>
            </div>

            <div className="main-content-wrapper">
                    <div className="user-files">
                        <ul className="list-elem">
                            <li></li>
                            <li></li>
                            <li></li>
                        </ul>
            </div>

                <div className="arrow-gui-pointer"><i className="fa-solid fa-arrow-right"></i></div>
                <div className="rename-files">
                    <ul className="list-elem-re">
                        <li></li>
                        <li></li>
                        <li></li>
                </ul>
                    
                </div>
            </div>
            </div>
        </div>
    )
}
//Function to open directory dialog
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

async function CallSeacrhShowGo(){
    let showTitle = document.querySelector(".user-input-search")
    let episodeList = await SearchShow(showTitle.value)
    
    
    await renderEpList(episodeList)
}

async function renderFileList(list){
    let fileList = document.querySelector(".list-elem")
    fileList.innerHTML = ""
    
    for (let i = 0; i < list.length; i++) {
        let liEl = document.createElement("li")
        liEl.append(i+1+". "+list[i])
        fileList.append(liEl)
      }
}

async function renderEpList(list){
    let EpList = document.querySelector(".list-elem-re")
    
    EpList.innerHTML = ""
    for (let i = 0; i < list.length; i++) {
        let liEl = document.createElement("li")
        liEl.append(i+1+". "+list[i])
        EpList.append(liEl)
      }
}



export default App
