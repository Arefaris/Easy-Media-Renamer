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
  autoMatch: (files, episodes) => call(App.AutoMatch, files, episodes),
  presets: () => call(App.ListPresets),
  renderName: (file, episode, show) => call(App.RenderNamePreview, file, episode, show),
  renderTemplate: (template, file, episode, show) => call(App.RenderTemplatePreview, template, file, episode, show),
  preview: (dir, files, episodes, show) => call(App.PreviewRename, dir, files, episodes, show),
  previewMatched: (dir, files, episodes, pairs, show) => call(App.PreviewMatched, dir, files, episodes, pairs, show),
  apply: (dir, ops) => call(App.ApplyRename, dir, ops),
  undo: () => call(App.UndoLastRename),
  history: () => call(App.ListHistory),
  revertHistory: id => call(App.RevertHistory, id),
  clearHistory: () => call(App.ClearHistory),
  checksum: (path, algo) => call(App.ComputeChecksum, path, algo),
  verifyCRC: path => call(App.VerifyEmbeddedCRC, path),
  getConfig: () => call(App.GetConfig),
  saveConfig: cfg => call(App.SaveConfig, cfg),
}
