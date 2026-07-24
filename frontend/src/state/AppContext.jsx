import {createContext, useContext, useReducer} from 'react'

const initial = {providers:[], provider:'TVmaze', query:'', shows:[], show:null, episodes:[], dir:'', files:[], preview:[], loading:false, error:'', settings:false, config:null}
const Context = createContext(null)
function reducer(state, action) {
  if (action.type === 'patch') return {...state, ...action.value}
  if (action.type === 'error') return {...state, loading:false, error:action.value}
  return state
}
export function AppProvider({children}) { const [state, dispatch] = useReducer(reducer, initial); return <Context.Provider value={{state,dispatch}}>{children}</Context.Provider> }
export const useApp = () => useContext(Context)
