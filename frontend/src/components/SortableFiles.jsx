import {DndContext, closestCenter} from '@dnd-kit/core'
import {arrayMove, SortableContext, useSortable, verticalListSortingStrategy} from '@dnd-kit/sortable'
import {CSS} from '@dnd-kit/utilities'
import {DragIndicator} from '@mui/icons-material'
import {List, ListItem, ListItemIcon, ListItemText, Paper, Typography} from '@mui/material'

function Row({file}) { const {attributes,listeners,setNodeRef,transform,transition}=useSortable({id:file.path}); return <ListItem ref={setNodeRef} style={{transform:CSS.Transform.toString(transform),transition}} {...attributes} {...listeners}><ListItemIcon><DragIndicator/></ListItemIcon><ListItemText primary={file.name} secondary={`${(file.size/1048576).toFixed(1)} MB`}/></ListItem> }
export default function SortableFiles({files,onChange}) { const end=({active,over})=>{if(over&&active.id!==over.id){const from=files.findIndex(x=>x.path===active.id),to=files.findIndex(x=>x.path===over.id);onChange(arrayMove(files,from,to))}}; return <Paper className="panel"><Typography variant="h6">Original files <small>({files.length})</small></Typography>{files.length===0?<p className="empty">Choose a folder with media files</p>:<DndContext collisionDetection={closestCenter} onDragEnd={end}><SortableContext items={files.map(x=>x.path)} strategy={verticalListSortingStrategy}><List>{files.map(f=><Row key={f.path} file={f}/>)}</List></SortableContext></DndContext>}</Paper> }
