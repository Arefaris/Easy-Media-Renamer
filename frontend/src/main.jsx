import React from 'react'
import {createRoot} from 'react-dom/client'
import './style.css'
import App from './App'
import { createTheme, ThemeProvider } from '@mui/material/styles';

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


root.render(
    <React.StrictMode>
        <ThemeProvider theme={theme}>
            <App/>
        </ThemeProvider>
    </React.StrictMode>
)
