import React from 'react'
import {createRoot} from 'react-dom/client'
import './style.css'
import App from './App'
import { createTheme, ThemeProvider } from '@mui/material/styles';
import {WindowSetTitle, WindowSetSize} from "../wailsjs/runtime/runtime"

const container = document.getElementById('root')

const root = createRoot(container)

// Theme of the app
const theme = createTheme({
    palette: {
      primary: {
        main: '#669bbc', // blue
      },
      secondary: {
        main: '#dc004e', // red
      },
    },
  });

//Setting tittle
WindowSetTitle("Easy Media Renamer")
//Setting default window size
WindowSetSize(640, 480)
root.render(
    <React.StrictMode>
        <ThemeProvider theme={theme}>
            <App/>
        </ThemeProvider>
    </React.StrictMode>
)
