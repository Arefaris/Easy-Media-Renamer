import React from 'react'
import {createRoot} from 'react-dom/client'
import '@fontsource/roboto/400.css'
import '@fontsource/roboto/500.css'
import {CssBaseline, ThemeProvider} from '@mui/material'
import {createTheme} from '@mui/material/styles'
import './App.css'
import App from './App'

const theme = createTheme({
  palette: {
    mode: 'dark',
    primary: {main: '#5865f2', light: '#7983f5', dark: '#4752c4'},
    secondary: {main: '#2f9df4'},
    background: {default: '#1e1f22', paper: '#2b2d31'},
    text: {primary: '#f2f3f5', secondary: '#b5bac1'},
    divider: '#3f4147',
    error: {main: '#f23f42'},
    success: {main: '#23a559'},
    warning: {main: '#f0b232'},
  },
  shape: {borderRadius: 6},
  typography: {
    fontFamily: 'Roboto, "Segoe UI", sans-serif',
    button: {fontWeight: 500, textTransform: 'none', letterSpacing: 0},
    h5: {fontSize: '1rem', fontWeight: 500},
    h6: {fontSize: '.82rem', fontWeight: 500, letterSpacing: '.02em'},
  },
  components: {
    MuiCssBaseline: {styleOverrides: {body: {scrollbarColor: '#4e5058 #1e1f22'}}},
    MuiPaper: {styleOverrides: {root: {backgroundImage: 'none'}}},
    MuiButton: {defaultProps: {disableElevation: true}, styleOverrides: {root: {minHeight: 34, borderRadius: 4}, containedPrimary: {'&:hover': {backgroundColor: '#4752c4'}}}},
    MuiIconButton: {styleOverrides: {root: {borderRadius: 4, color: '#b5bac1', '&:hover': {color: '#f2f3f5', backgroundColor: '#35373c'}}}},
    MuiOutlinedInput: {styleOverrides: {root: {backgroundColor: '#1e1f22', '& fieldset': {borderColor: '#3f4147'}, '&:hover fieldset': {borderColor: '#5c5f66'}, '&.Mui-focused fieldset': {borderColor: '#5865f2'}}}},
    MuiInputLabel: {styleOverrides: {root: {color: '#949ba4'}}},
    MuiListItemButton: {styleOverrides: {root: {borderRadius: 4, margin: '1px 6px', width: 'auto', '&:hover': {backgroundColor: '#35373c'}, '&.Mui-selected': {backgroundColor: '#404249'}}}},
    MuiDialog: {styleOverrides: {paper: {border: '1px solid #3f4147', boxShadow: '0 18px 48px rgba(0,0,0,.45)'}}},
  },
})

createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <ThemeProvider theme={theme}>
      <CssBaseline/>
      <App/>
    </ThemeProvider>
  </React.StrictMode>,
)
