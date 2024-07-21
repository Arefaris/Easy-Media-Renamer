import './App.css';
import { OpenDirectoryDialog } from '../wailsjs/go/main/App';
import { FilesInDirectoryHandler } from "../wailsjs/go/main/App"
function App() {

    return (
        <div id="App">
            <div>

                <header>
                    <p1>MEDIA RENAMER</p1>
                </header>

                <div className="main-content-wrapper">
                    <div className="user-files">
                    <p2>Choose folder with your files</p2>
                        <button className="btn-chooser" onClick={openDirectory}>Directory Chooser</button>
                        <ul className="list-elem">
                            <li>jkson</li>
                            <li>jkson</li>
                            <li>jkson</li>
                        </ul>
                    </div>

                <div className="rename-files">
                    <p2>Type a show name and press enter</p2>
                    <input className="user-input-search" type="text"></input>
                    <ul className="list-elem-re">
                        <li>fast and furious.mkv</li>
                        <li>fast and furious.mkv</li>
                        <li>fast and furious.mkv</li>
                    </ul>
                    
                </div>
            </div>
            <div className="rename-btn-section">
                    <button className="btn-chooser" onClick={openDirectory}>Rename all</button>
                    <button className="btn-chooser" onClick={openDirectory}>Rename selected</button>
                </div>
            </div>
        </div>
    )
}

async function openDirectory() {
    try {
        const result = await OpenDirectoryDialog();
        console.log("Selected directory:", result);
        await FilesInDirectoryHandler(result)
    } catch (error) {
        console.error("Failed to open directory dialog:", error);
    }
}

export default App
