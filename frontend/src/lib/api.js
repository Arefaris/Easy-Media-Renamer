import * as App from '../../wailsjs/go/main/App'

const call = async (fn, ...args) => {
  try { return await fn(...args) }
  catch (error) { throw new Error(typeof error === 'string' ? error : error?.message || 'Unknown application error') }
}

export const api = {
  providers: () => call(App.ListProviders),
  search: (provider, query) => call(App.SearchShows, provider, query),
  episodes: (provider, id) => call(App.ListEpisodes, provider, id),
  chooseDirectory: () => call(App.OpenDirectoryDialog),
  scan: dir => call(App.ScanDirectory, dir),
  preview: (dir, files, episodes, show) => call(App.PreviewRename, dir, files, episodes, show),
  apply: (dir, ops) => call(App.ApplyRename, dir, ops),
  undo: () => call(App.UndoLastRename),
  getConfig: () => call(App.GetConfig),
  saveConfig: cfg => call(App.SaveConfig, cfg),
}
