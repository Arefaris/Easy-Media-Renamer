import './App.css';
import React from 'react'
import RenameBtnSc from "./OpenRename"
import Search from "./HandleSearch"


//rendering
function App() {
    return (
        <div id="App">
            <div>
                <RenameBtnSc />
                <Search />
            </div>
        </div>
    
    )
}


export default App
