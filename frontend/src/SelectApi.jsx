import React from 'react'
import Box from '@mui/material/Box';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import FormControl from '@mui/material/FormControl';
import Select from '@mui/material/Select';


export let selectedApi

export default function SelectApi() {

  let [api, setApi] = React.useState('TVmaze');
  
  const handleChange = (event) => {
    setApi(event.target.value);
    selectedApi = event.target.value

  };

  return (
    <Box sx={{ minWidth: 80}}>
      <FormControl variant="standard" fullWidth>
        
        <InputLabel id="demo-simple-select-label"></InputLabel>
        <Select
          defaultValue="TVmaze"
          labelId="demo-simple-select-label"
          id="demo-simple-select"
          value={api}
          label="api"
          color="primary"
          onChange={handleChange}
          sx={{
            '&:before': {
            borderColor: '#669bbc',
        },
        '&:after': {
            borderColor: '#669bbc',
        },
            color: '#e5e5e5',
            '& .MuiOutlinedInput-notchedOutline': {
                borderColor: '#669bbc'
            },
            '& .MuiSvgIcon-root': {
                color: '#669bbc'
            }, 
            '&:not(.Mui-disabled):hover::before': {
            borderColor: '#669bbc',
            },
            root: {
                color: '#669bbc',
            },
        }}
    
        >
          <MenuItem value={"TVmaze"}>TVmaze</MenuItem>
          <MenuItem value={"AniDB"}>AniDB</MenuItem>
          <MenuItem value={"TMDB"}>TMDB</MenuItem>
        </Select>
      </FormControl>
    </Box>
  );
}

