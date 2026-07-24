import {useEffect, useState} from 'react'
import {Button, Checkbox, Dialog, DialogActions, DialogContent, DialogTitle, FormControl, FormControlLabel, InputLabel, MenuItem, Select, Stack, TextField, Typography} from '@mui/material'
import {api} from '../lib/api'

const tokens = '{n} {s} {e} {s00e00} {sxe} {t} {y} {airdate} {abs} {ext} {fn} {vf} {vc} {ac} {channels} {hd} {resolution} {duration} {source} {group} {crc32}'

export default function SettingsDialog({open, config, presets, sample, onClose, onSave}) {
  const [draft, setDraft] = useState(config)
  const [extensionsText, setExtensionsText] = useState('')
  const [rendered, setRendered] = useState('')

  useEffect(() => { setDraft(config); setExtensionsText((config?.media_extensions || []).join(' ')) }, [config])
  useEffect(() => {
    if (!open || !draft?.name_template || !sample?.file || !sample?.episode || !sample?.show) { setRendered(''); return }
    const timer=setTimeout(()=>api.renderTemplate(draft.name_template,sample.file,sample.episode,sample.show).then(setRendered).catch(e=>setRendered(e.message)),180)
    return ()=>clearTimeout(timer)
  }, [open,draft?.name_template,sample])

  if (!draft) return null
  const set = (key, value) => setDraft({...draft, [key]: value})
  const preset = name => { const p=presets.find(x=>x.name===name); setDraft({...draft,preset:name,name_template:p?.template||draft.name_template}) }
  const save = () => onSave({...draft, media_extensions: extensionsText.split(/[,\s]+/).filter(Boolean)})

  return <Dialog open={open} onClose={onClose} fullWidth maxWidth="md"><DialogTitle>Settings</DialogTitle><DialogContent><Stack spacing={2} sx={{mt:1}}>
    <TextField label="TMDB API key" type="password" value={draft.tmdb_api_key||''} onChange={e=>set('tmdb_api_key',e.target.value)}/>
    <FormControl><InputLabel>Template preset</InputLabel><Select label="Template preset" value={draft.preset||'Simple'} onChange={e=>preset(e.target.value)}>{presets.map(p=><MenuItem key={p.name} value={p.name}>{p.name}</MenuItem>)}</Select></FormControl>
    <TextField label="Name template" value={draft.name_template||''} onChange={e=>set('name_template',e.target.value)} helperText={`Tokens: ${tokens}. [ - {group}] is removed when group is empty. Use / for folders and \\[...\\] for literal brackets.`}/>
    {rendered&&<Typography className="template-preview">Preview: {rendered}</Typography>}
    <TextField label="Media extensions" multiline value={extensionsText} onChange={e=>setExtensionsText(e.target.value)} helperText="Separate extensions with spaces or commas, for example: .mkv .mp4 .txt"/>
    <Stack direction="row" spacing={2}><FormControl sx={{minWidth:180}}><InputLabel>File action</InputLabel><Select label="File action" value={draft.action||'rename'} onChange={e=>set('action',e.target.value)}>{['rename','move','copy','hardlink','symlink','test'].map(x=><MenuItem key={x} value={x}>{x}</MenuItem>)}</Select></FormControl><TextField label="Maximum depth" type="number" value={draft.max_depth||5} onChange={e=>set('max_depth',Number(e.target.value))} inputProps={{min:1,max:99}}/><TextField label="Metadata language" value={draft.language||'en-US'} onChange={e=>set('language',e.target.value)}/></Stack>
    <Stack direction="row"><FormControlLabel control={<Checkbox checked={!!draft.recursive} onChange={e=>set('recursive',e.target.checked)}/>} label="Scan subfolders"/><FormControlLabel control={<Checkbox checked={!!draft.auto_match} onChange={e=>set('auto_match',e.target.checked)}/>} label="Auto-match by default"/><FormControlLabel control={<Checkbox checked={!!draft.include_specials} onChange={e=>set('include_specials',e.target.checked)}/>} label="Include specials"/></Stack>
  </Stack></DialogContent><DialogActions><Button onClick={onClose}>Cancel</Button><Button variant="contained" onClick={save}>Save</Button></DialogActions></Dialog>
}
