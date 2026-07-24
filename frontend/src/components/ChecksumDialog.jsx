import {Button, Chip, Dialog, DialogActions, DialogContent, DialogTitle, List, ListItem, ListItemText} from '@mui/material'

export default function ChecksumDialog({open, results, onClose, onVerify}) {
  return <Dialog open={open} onClose={onClose} fullWidth maxWidth="md"><DialogTitle>Embedded CRC32 verification</DialogTitle><DialogContent>{results.length===0?<p className="empty">Checks names containing [A1B2C3D4].</p>:<List>{results.map((r,i)=><ListItem key={i} secondaryAction={<Chip size="small" color={r.ok?'success':'error'} label={r.ok?'valid':'mismatch'}/>}><ListItemText primary={r.file} secondary={r.error||`expected ${r.expected||'—'} · actual ${r.actual||'—'}`}/></ListItem>)}</List>}</DialogContent><DialogActions><Button variant="contained" onClick={onVerify}>Verify files</Button><Button onClick={onClose}>Close</Button></DialogActions></Dialog>
}
