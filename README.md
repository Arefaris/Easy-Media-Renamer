# Easy Media Renamer

Desktop media renamer built with Go, Wails v2 and React. It parses release filenames, matches them against TVmaze, TMDB or AniDB, shows a safe preview, then applies the rename. Every operation is revertible.

A console binary (`emr`) exposes the same engine for scripting.

## Features

**Matching**

- Automatic file-to-episode matching — no manual reordering required. Filenames are parsed for `S01E01`, `1x01`, `Season 1 Episode 1`, air dates, anime absolute numbering and multi-episode ranges (`S01E01-E03`).
- Five strategies, tried in order of reliability: season+episode, absolute number, air date, episode-only, and positional fallback. Each result carries a strategy label and a confidence score; low-confidence matches are highlighted in the UI.
- Files that match nothing stay visible as `unmatched` and are never renamed.
- Manual drag & drop still works — auto-match only sets the starting point.

**Sources**

- TVmaze — no credentials required.
- TMDB — TV shows and movies; movies resolve to a single entry.
- AniDB — HTTPS title dump, 24-hour cache, rate-limited API calls. The public client identifier is built in, so no user configuration is needed.

**Files**

- Recursive scanning with a depth limit; hidden and system directories are skipped.
- Media-only filtering and natural sorting (`ep2` before `ep10`).
- Actions: `rename`, `move`, `copy`, `hardlink`, `symlink`, `test` (dry run).
- Cross-volume moves fall back to copy+delete with a size check.
- Two-phase rename handles cycles such as `a.mkv → b.mkv` / `b.mkv → a.mkv`.
- Conflict detection before anything is written, plus a re-check at apply time.
- History of the last 50 operations; any entry can be reverted, not just the most recent.

**Naming**

- Token-based templates with conditional blocks, subfolders and Plex/Kodi/Jellyfin presets.
- Windows-safe filenames: invalid characters, control characters, trailing dots/spaces, reserved device names and the 255-byte limit are all handled per path segment.
- Pure-Go MediaInfo (no ffmpeg required) for resolution, codecs, channels and HDR.

**Other**

- CRC32, MD5, SHA1 and SHA256 checksums; verification of 8-digit CRC32 embedded in anime release names.
- Offline fonts and icons; runtime errors are surfaced without crashing the app.

## Naming templates

### Tokens

| Group | Tokens |
|---|---|
| Show | `{n}` `{show}` · `{y}` `{year}` |
| Episode | `{t}` `{title}` · `{airdate}` · `{abs}` |
| Numbering | `{s}` `{e}` (unpadded) · `{season}` `{episode}` `{s00}` `{e00}` (padded) · `{s00e00}` → `S01E05` · `{sxe}` → `1x05` |
| File | `{ext}` · `{fn}` (original name without extension) |
| Media | `{vf}` `{resolution}` · `{vc}` · `{ac}` · `{channels}` · `{hd}` · `{duration}` |
| Release | `{source}` · `{group}` · `{crc32}` |

Release and media tokens are read from the filename and the file itself, so they survive renaming.

### Conditional blocks

Text in `[ ... ]` disappears entirely when any token inside it is empty:

```text
{n} - {s00e00} - {t}[ - {vf}][ - {group}]

Frieren - S01E05 - Journey's End - 1080p - SubsPlease.mkv
Frieren - S01E05 - Journey's End.mkv          (no media info, no group)
```

For literal square brackets — the norm in anime naming — escape them:

```text
\[{group}\] {n} - {e00}     →     [SubsPlease] Frieren - 05.mkv
```

### Subfolders

A `/` in the template creates directories:

```text
{n} ({y})/Season {s}/{n} - {s00e00} - {t}

Frieren (2023)/Season 1/Frieren - S01E05 - Journey's End.mkv
```

Slashes inside token *values* are sanitized, so an episode titled `Us/Them` becomes `Us-Them` rather than a stray folder.

### Presets

| Preset | Template |
|---|---|
| Plex | `{n} ({y})/Season {s}/{n} - {s00e00} - {t}` |
| Kodi | `{n} ({y})/Season {s}/{n} S{s00}E{e00} {t}` |
| Jellyfin | `{n} ({y})/Season {s}/{n} S{s00}E{e00} {t}` |
| Simple | `{n} - {s00e00} - {t}` |

The extension is appended automatically unless the template contains `{ext}`, in which case you control placement yourself.

## MediaInfo

Parsed in pure Go — no ffmpeg or MediaInfo.dll dependency:

| Container | Support |
|---|---|
| `.mkv`, `.webm` | resolution, video/audio codec, channels, duration, HDR (PQ / HLG) |
| `.mp4`, `.m4v`, `.mov` | resolution, video/audio codec, duration |
| others | media tokens resolve to empty; conditional blocks collapse |

Files are probed lazily — only when the active template actually uses a media token.

## Configuration

Open **Settings** in the application. Configuration is stored with mode `0600` in the user config directory:

```text
%APPDATA%\EasyMediaRenamer\config.json
```

| Key | Meaning |
|---|---|
| `tmdb_api_key` | TMDB v3 API key |
| `name_template` | active naming template |
| `media_extensions` | scanned extensions |
| `include_specials` | include season 0 |
| `recursive` / `max_depth` | recursive scanning and its depth limit |
| `action` | `rename`, `move`, `copy`, `hardlink`, `symlink`, `test` |
| `auto_match` | run matching automatically |
| `preset` | selected naming preset |
| `language` | metadata language |

Environment overrides: `TMDB_API_KEY` for the key, `EMR_CACHE_DIR` to relocate the cache and history for portable or test setups.

TMDB requires a v3 API key from the [TMDB account settings](https://www.themoviedb.org/settings/api). AniDB needs no user configuration.

## Command line

`emr` shares the engine with the GUI and is built as a console application.

```powershell
emr -list -r ./Media
emr -check -algo sha256 ./Media
emr -rename -r --db TMDB --action move --dry-run ./Downloads
```

| Flag | Meaning |
|---|---|
| `-list` | list matching media files |
| `-check` | compute checksums (`-algo crc32\|md5\|sha1\|sha256`) |
| `-rename` | match and rename |
| `-r` | scan recursively |
| `-db` | `TVmaze`, `TMDB` or `AniDB` |
| `-format` | naming template (defaults to the configured one) |
| `-action` | file operation |
| `-filter` | comma-separated extensions |
| `-lang` | metadata language |
| `-dry-run` | print operations without writing |

`-rename` infers the show name from the filenames themselves and picks the closest search result, so no query argument is needed. Files that cannot be matched are skipped rather than renamed.

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

Production builds:

```powershell
wails build            # GUI  -> build/bin/Easy Media Renamer.exe
go build -o build/bin/emr.exe ./cmd/emr
```

Optional live provider checks hit real network services. Credential-dependent cases are skipped when the corresponding environment variables are absent:

```powershell
go test -tags live ./internal/provider/...
```

No key file is embedded in the executable; a clean checkout compiles without secrets.

## Project layout

```text
app.go                 Wails bindings (thin)
cmd/emr/               console binary
internal/
  config/              runtime configuration
  match/               filename parsing, matching strategies, fuzzy title compare
  media/               scanning, naming, rename actions, history, presets
  mediainfo/           pure-Go MKV and MP4 probing
  provider/            TVmaze, TMDB, AniDB behind one interface
  verify/              checksums and embedded-CRC verification
frontend/src/          React UI
```
