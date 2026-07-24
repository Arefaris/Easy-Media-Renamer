# Easy Media Renamer

Desktop media renamer built with Go, Wails v2 and React. It searches TVmaze, TMDB and AniDB, pairs naturally sorted local files with episodes, shows a safe preview, then applies the rename. The last operation can be undone.

## Features

- TVmaze works without credentials.
- TMDB supports TV shows and movies.
- AniDB uses the HTTPS title dump, a 24-hour cache and rate-limited HTTP API calls.
- Media-only directory scan and natural sorting (`ep2` before `ep10`).
- Configurable naming template and Windows-safe filenames.
- Conflict detection, two-phase rename for cycles, preview and undo journal.
- Offline fonts and icons; runtime errors are shown without crashing the app.

## Configuration

Open **Settings** in the application. Configuration is stored with mode `0600` in the user config directory:

```text
%APPDATA%\EasyMediaRenamer\config.json
```

Supported environment override: `TMDB_API_KEY`. `EMR_CACHE_DIR` optionally overrides the cache/undo location for portable or test environments.

TMDB requires a v3 API key from the [TMDB account settings](https://www.themoviedb.org/settings/api). The public AniDB HTTP API identifier for Easy Media Renamer is built in, so AniDB needs no user configuration.

Naming tokens: `{show}`, `{season}`, `{episode}`, `{title}`, `{year}`, `{ext}`.

## Development

Requirements: Go 1.25+, Node.js 20.19+ or 22.12+, npm, WebView2 on Windows.

```powershell
go install github.com/wailsapp/wails/v2/cmd/wails@v2.13.0
Set-Location frontend
npm install
npm run build
Set-Location ..
go test ./...
go vet ./...
wails dev
```

Production build:

```powershell
wails build
```

Optional live provider checks use real network services. Credential-dependent cases are skipped when environment variables are absent:

```powershell
go test -tags live ./internal/provider/...
```

No key file is embedded in the executable; a clean checkout can compile without secrets.
