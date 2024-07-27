import Button from '@mui/material/Button';
import Input from '@mui/material/Input';
import {SearchShow} from "../wailsjs/go/main/App";
import React, {useRef, useState, useEffect} from 'react'
import { addListenersForFiles } from './OpenRename';
import Sortable from 'sortablejs';

const Search = () => {
        const [inputValue, setInputValue] = useState('');

        let listElemRef = useRef(null);

        useEffect(() => {
            if (listElemRef.current) {
                //making list sortable
                Sortable.create(listElemRef.current)
            }
        }, []);

      
        const handleChange = (event) => {
          setInputValue(event.target.value);
        };
      
        const handleSubmit = async () => {
          await CallSearchShowGo(inputValue);
        };

    return(
        <>
        <div className="input-wrapper">
                <Input className="user-input-search" type="text" placeholder="Show name" value={inputValue} onChange={handleChange} onKeyUp={event => {
                if (event.key === 'Enter') {
                  handleSubmit()
                }
              }}></Input>
                <Button variant="outlined" className="srch-btn" onClick={handleSubmit}><i className="fa-solid fa-magnifying-glass"></i></Button>
            </div>

            <div className="search-modal">
                <ul className="list-search">
                </ul>
            </div>

            <div className="main-content-wrapper">
                    <div className="user-files">
                    <div>Original files</div>
                        <ul className="list-elem" ref={listElemRef}>
                        </ul>
                    </div>

                <div className="arrow-gui-pointer"><i className="fa-solid fa-arrow-right"></i></div>
                    <div className="rename-files">
                    <div>New files</div>
                        <ul className="list-elem-re">
                        </ul>
                    </div>
            </div>


        </>
    )
}



//Calling an api
const CallSearchShowGo = async (input) => {
    if (input) {
        let showList = await SearchShow(input)
        await renderShowList(showList)
    }
    
}


const renderShowList = async (list) =>{
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


export default Search