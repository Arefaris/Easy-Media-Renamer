import {Button, Dialog, DialogActions, DialogContent, DialogTitle, List, ListItem, ListItemText} from '@mui/material'

export default function HistoryDialog({open, entries, onClose, onRefresh, onRevert, onClear}) {
  return <Dialog open={open} onClose={onClose} fullWidth maxWidth="md">
    <DialogTitle>Operation history</DialogTitle>
    <DialogContent>{entries.length===0?<p className="empty">History is empty</p>:<List>{entries.map(h=><ListItem key={h.id} secondaryAction={<Button color="warning" onClick={()=>onRevert(h.id)}>Revert</Button>}><ListItemText primary={`${h.action} · ${h.count} file${h.count===1?'':'s'}`} secondary={`${new Date(h.time).toLocaleString()} · ${h.directory}`}/></ListItem>)}</List>}</DialogContent>
    <DialogActions><Button onClick={onRefresh}>Refresh</Button><Button color="error" disabled={!entries.length} onClick={onClear}>Clear</Button><Button onClick={onClose}>Close</Button></DialogActions>
  </Dialog>
}
