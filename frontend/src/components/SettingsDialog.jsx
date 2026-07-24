import {useEffect, useState} from 'react'
import {Button, Checkbox, Dialog, DialogActions, DialogContent, DialogTitle, FormControlLabel, Stack, TextField} from '@mui/material'

export default function SettingsDialog({open, config, onClose, onSave}) {
  const [draft, setDraft] = useState(config)
  const [extensionsText, setExtensionsText] = useState('')

  useEffect(() => {
    setDraft(config)
    setExtensionsText((config?.media_extensions || []).join(' '))
  }, [config])

  if (!draft) return null
  const set = (key, value) => setDraft({...draft, [key]: value})
  const save = () => onSave({
    ...draft,
    media_extensions: extensionsText.split(/[,\s]+/).filter(Boolean),
  })

  return <Dialog open={open} onClose={onClose} fullWidth maxWidth="sm">
    <DialogTitle>Settings</DialogTitle>
    <DialogContent>
      <Stack spacing={2} sx={{mt: 1}}>
        <TextField label="TMDB API key" type="password" value={draft.tmdb_api_key || ''} onChange={e => set('tmdb_api_key', e.target.value)}/>
        <TextField label="Name template" value={draft.name_template || ''} onChange={e => set('name_template', e.target.value)} helperText="Tokens: {show} {season} {episode} {title} {year} {ext}"/>
        <TextField label="Media extensions" multiline value={extensionsText} onChange={e => setExtensionsText(e.target.value)} helperText="Separate extensions with spaces or commas, for example: .mkv .mp4 .txt"/>
        <FormControlLabel control={<Checkbox checked={!!draft.include_specials} onChange={e => set('include_specials', e.target.checked)}/>} label="Include specials"/>
      </Stack>
    </DialogContent>
    <DialogActions><Button onClick={onClose}>Cancel</Button><Button variant="contained" onClick={save}>Save</Button></DialogActions>
  </Dialog>
}
