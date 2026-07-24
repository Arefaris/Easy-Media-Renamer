import {useEffect} from 'react'
import {ArrowForward, FolderOpen, Search, Settings, Undo} from '@mui/icons-material'
import {Alert, AppBar, Box, Button, CircularProgress, Container, FormControl, IconButton, InputLabel, List, ListItemButton, ListItemText, MenuItem, Paper, Select, Snackbar, TextField, Toolbar, Tooltip, Typography} from '@mui/material'
import {api} from './lib/api'
import {AppProvider, useApp} from './state/AppContext'
import SortableFiles from './components/SortableFiles'
import SettingsDialog from './components/SettingsDialog'
import PreviewDialog from './components/PreviewDialog'

function Workspace() {
  const {state:s,dispatch} = useApp(); const patch=value=>dispatch({type:'patch',value}); const fail=e=>dispatch({type:'error',value:e.message})
  useEffect(()=>{Promise.all([api.providers(),api.getConfig()]).then(([providers,config])=>patch({providers,config})).catch(fail)},[])
  const busy=async fn=>{patch({loading:true,error:''});try{await fn()}catch(e){fail(e)}finally{patch({loading:false})}}
  const choose=()=>busy(async()=>{const dir=await api.chooseDirectory();if(!dir)return;patch({dir,files:await api.scan(dir),preview:[]})})
  const search=()=>busy(async()=>patch({shows:await api.search(s.provider,s.query),show:null,episodes:[]}))
  const selectShow=show=>busy(async()=>patch({show,shows:[],episodes:await api.episodes(s.provider,show.id)}))
  const preview=()=>busy(async()=>patch({preview:await api.preview(s.dir,s.files.map(f=>f.name),s.episodes,s.show)}))
  const apply=()=>busy(async()=>{const r=await api.apply(s.dir,s.preview);patch({preview:[],files:await api.scan(s.dir)});if(r.failed?.length)throw new Error(`Renamed ${r.renamed}; failed ${r.failed.length}`)})
  const undo=()=>busy(async()=>{await api.undo();if(s.dir)patch({files:await api.scan(s.dir)})})
  const save=cfg=>busy(async()=>{await api.saveConfig(cfg);patch({config:await api.getConfig(),providers:await api.providers(),settings:false})})
  const current=s.providers.find(x=>x.name===s.provider)
  return <><AppBar position="static" color="transparent" elevation={0}><Toolbar><Typography variant="h5" sx={{flexGrow:1,fontWeight:700}}>Easy Media Renamer</Typography><Tooltip title="Undo last rename"><span><IconButton onClick={undo} disabled={s.loading}><Undo/></IconButton></span></Tooltip><IconButton onClick={()=>patch({settings:true})}><Settings/></IconButton></Toolbar></AppBar>
    <Container maxWidth="xl" className="workspace">
      <Paper className="controls"><Button variant="contained" startIcon={<FolderOpen/>} onClick={choose}>Choose folder</Button><Typography className="path" title={s.dir}>{s.dir||'No directory selected'}</Typography><Button variant="outlined" endIcon={<ArrowForward/>} disabled={!s.dir||!s.show||!s.files.length||!s.episodes.length||s.loading} onClick={preview}>Preview rename</Button></Paper>
      <Paper className="search"><FormControl size="small" sx={{minWidth:150}}><InputLabel>Provider</InputLabel><Select value={s.provider} label="Provider" onChange={e=>patch({provider:e.target.value,shows:[],episodes:[],show:null})}>{s.providers.map(p=><MenuItem key={p.name} value={p.name}>{p.name}{!p.configured?' · setup required':''}</MenuItem>)}</Select></FormControl><TextField size="small" fullWidth label="Search show or movie" value={s.query} onChange={e=>patch({query:e.target.value})} onKeyDown={e=>e.key==='Enter'&&search()}/><Button startIcon={<Search/>} onClick={current?.configured?search:()=>patch({settings:true})} variant="contained">{current?.configured?'Search':'Configure'}</Button></Paper>
      {s.shows.length>0&&<Paper className="results"><List>{s.shows.map(show=><ListItemButton key={show.id} onClick={()=>selectShow(show)}><ListItemText primary={show.name} secondary={[show.year,show.kind,show.genres?.join(', ')].filter(Boolean).join(' · ')}/></ListItemButton>)}</List></Paper>}
      <Box className="columns"><SortableFiles files={s.files} onChange={files=>patch({files})}/><Paper className="panel"><Typography variant="h6">{s.show?s.show.name:'Episodes'} <small>({s.episodes.length})</small></Typography>{s.episodes.length===0?<p className="empty">Search and select a show</p>:<List>{s.episodes.map((ep,i)=><ListItemButton key={`${ep.season}-${ep.number}-${i}`}><ListItemText primary={`s${String(ep.season).padStart(2,'0')}e${String(ep.number).padStart(2,'0')} · ${ep.title}`} secondary={ep.special?'Special':null}/></ListItemButton>)}</List>}</Paper></Box>
    </Container>
    {s.loading&&<div className="loading"><CircularProgress/></div>}
    <PreviewDialog ops={s.preview} onClose={()=>patch({preview:[]})} onApply={apply}/><SettingsDialog open={s.settings} config={s.config} onClose={()=>patch({settings:false})} onSave={save}/><Snackbar open={!!s.error} autoHideDuration={7000} onClose={()=>patch({error:''})}><Alert severity="error" onClose={()=>patch({error:''})}>{s.error}</Alert></Snackbar></>
}
export default function App(){return <AppProvider><Workspace/></AppProvider>}
