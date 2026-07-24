import {useEffect} from 'react'
import {ArrowForward, AutoFixHigh, FactCheck, FolderOpen, History, Search, Settings} from '@mui/icons-material'
import {Alert, AppBar, Box, Button, Chip, CircularProgress, Container, FormControl, IconButton, InputLabel, List, ListItemButton, ListItemText, MenuItem, Paper, Select, Snackbar, TextField, Toolbar, Tooltip, Typography} from '@mui/material'
import {api} from './lib/api'
import {AppProvider, useApp} from './state/AppContext'
import SortableFiles from './components/SortableFiles'
import SettingsDialog from './components/SettingsDialog'
import PreviewDialog from './components/PreviewDialog'
import HistoryDialog from './components/HistoryDialog'
import ChecksumDialog from './components/ChecksumDialog'

function Workspace() {
  const {state:s,dispatch}=useApp()
  const patch=value=>dispatch({type:'patch',value})
  const fail=e=>dispatch({type:'error',value:e?.message||String(e)})
  useEffect(()=>{Promise.all([api.providers(),api.getConfig(),api.presets()]).then(([providers,config,presets])=>patch({providers,config,presets})).catch(fail)},[])
  const busy=async fn=>{patch({loading:true,error:''});try{await fn()}catch(e){fail(e)}finally{patch({loading:false})}}
  const rescan=async()=>s.dir?api.scan(s.dir):[]
  const choose=()=>busy(async()=>{const dir=await api.chooseDirectory();if(!dir)return;patch({dir,files:await api.scan(dir),preview:[],pairs:[]})})
  const search=()=>busy(async()=>patch({shows:await api.search(s.provider,s.query),show:null,episodes:[],pairs:[]}))
  const selectShow=show=>busy(async()=>{const episodes=await api.episodes(s.provider,show.id);patch({show,shows:[],episodes,pairs:[]})})
  const autoMatch=()=>busy(async()=>{
    const pairs=await api.autoMatch(s.files,s.episodes)
    patch({pairs})
  })
  const preview=()=>busy(async()=>patch({preview:s.pairs.length?await api.previewMatched(s.dir,s.files,s.episodes,s.pairs,s.show):await api.preview(s.dir,s.files.map(f=>f.rel_path||f.name),s.episodes,s.show)}))
  const apply=()=>busy(async()=>{const r=await api.apply(s.dir,s.preview);patch({preview:[],pairs:[],files:await rescan(),notice:`Completed: ${r.renamed}, skipped: ${r.skipped}`});if(r.failed?.length)throw new Error(`Completed ${r.renamed}; failed ${r.failed.length}`)})
  const loadHistory=()=>busy(async()=>patch({history:await api.history(),historyOpen:true}))
  const revert=id=>busy(async()=>{await api.revertHistory(id);patch({history:await api.history(),files:await rescan(),notice:'Operation reverted'})})
  const clearHistory=()=>busy(async()=>{await api.clearHistory();patch({history:[]})})
  const verifyCRC=()=>busy(async()=>{
    const candidates=s.files.filter(f=>/\[[0-9a-f]{8}\]/i.test(f.name))
    const results=[]
    for(const f of candidates){try{results.push(await api.verifyCRC(f.path))}catch(e){results.push({file:f.name,expected:'',actual:'',ok:false,error:e.message})}}
    patch({checksumResults:results,checksumOpen:true})
  })
  const save=cfg=>busy(async()=>{await api.saveConfig(cfg);const config=await api.getConfig();patch({config,providers:await api.providers(),settings:false,files:await rescan(),pairs:[]})})
  const current=s.providers.find(x=>x.name===s.provider)
  const sample={file:s.files[0],episode:s.episodes[0],show:s.show}
  return <><AppBar position="static" color="transparent" elevation={0}><Toolbar><Typography variant="h5" sx={{flexGrow:1,fontWeight:700}}>Easy Media Renamer</Typography><Tooltip title="Checksums"><IconButton onClick={()=>patch({checksumOpen:true})}><FactCheck/></IconButton></Tooltip><Tooltip title="Operation history"><IconButton onClick={loadHistory}><History/></IconButton></Tooltip><IconButton onClick={()=>patch({settings:true})}><Settings/></IconButton></Toolbar></AppBar>
    <Container maxWidth="xl" className="workspace">
      <Paper className="controls"><Button variant="contained" startIcon={<FolderOpen/>} onClick={choose}>Choose folder</Button><Typography className="path" title={s.dir}>{s.dir||'No directory selected'}</Typography><Chip size="small" label={s.config?.action||'rename'} color="primary" variant="outlined"/><Button startIcon={<AutoFixHigh/>} disabled={!s.files.length||!s.episodes.length||s.loading} onClick={autoMatch}>Auto-match</Button><Button variant="outlined" endIcon={<ArrowForward/>} disabled={!s.dir||!s.show||!s.files.length||!s.episodes.length||s.loading} onClick={preview}>Preview</Button></Paper>
      <Paper className="search"><FormControl size="small" sx={{minWidth:150}}><InputLabel>Provider</InputLabel><Select value={s.provider} label="Provider" onChange={e=>patch({provider:e.target.value,shows:[],episodes:[],show:null,pairs:[]})}>{s.providers.map(p=><MenuItem key={p.name} value={p.name}>{p.name}{!p.configured?' · setup required':''}</MenuItem>)}</Select></FormControl><TextField size="small" fullWidth label="Search show or movie" value={s.query} onChange={e=>patch({query:e.target.value})} onKeyDown={e=>e.key==='Enter'&&search()}/><Button startIcon={<Search/>} onClick={current?.configured?search:()=>patch({settings:true})} variant="contained">{current?.configured?'Search':'Configure'}</Button></Paper>
      {s.shows.length>0&&<Paper className="results"><List>{s.shows.map(show=><ListItemButton key={show.id} onClick={()=>selectShow(show)}><ListItemText primary={show.name} secondary={[show.year,show.kind,show.genres?.join(', ')].filter(Boolean).join(' · ')}/></ListItemButton>)}</List></Paper>}
      <Box className="columns"><SortableFiles files={s.files} onChange={files=>patch({files,pairs:[]})}/><Paper className="panel"><Typography variant="h6">{s.show?s.show.name:'Episodes'} <small>({s.episodes.length})</small></Typography>{s.episodes.length===0?<p className="empty">Search and select a show</p>:<List>{s.episodes.map((ep,i)=>{const p=s.pairs.find(x=>x.episode_index===i);return <ListItemButton key={`${ep.season}-${ep.number}-${i}`} sx={p?.confidence<.6?{borderLeft:'3px solid #f0b232'}:{}}><ListItemText primary={`s${String(ep.season).padStart(2,'0')}e${String(ep.number).padStart(2,'0')} · ${ep.title}`} secondary={p?`${s.files[p.file_index]?.name||'file'} · ${p.strategy} · ${Math.round(p.confidence*100)}%`:ep.airdate|| (ep.special?'Special':null)}/></ListItemButton>})}</List>}</Paper></Box>
    </Container>
    {s.loading&&<div className="loading"><CircularProgress/></div>}
    <PreviewDialog ops={s.preview} pairs={s.pairs} onClose={()=>patch({preview:[]})} onApply={apply}/>
    <SettingsDialog open={s.settings} config={s.config} presets={s.presets} sample={sample} onClose={()=>patch({settings:false})} onSave={save}/>
    <HistoryDialog open={s.historyOpen} entries={s.history} onClose={()=>patch({historyOpen:false})} onRefresh={loadHistory} onRevert={revert} onClear={clearHistory}/>
    <ChecksumDialog open={s.checksumOpen} results={s.checksumResults} onClose={()=>patch({checksumOpen:false})} onVerify={verifyCRC}/>
    <Snackbar open={!!s.error} autoHideDuration={7000} onClose={()=>patch({error:''})}><Alert severity="error" onClose={()=>patch({error:''})}>{s.error}</Alert></Snackbar><Snackbar open={!!s.notice} autoHideDuration={4000} onClose={()=>patch({notice:''})} message={s.notice}/>
  </>
}
export default function App(){return <AppProvider><Workspace/></AppProvider>}
