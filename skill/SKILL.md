---
name: nba-cli
description: Guidance for AI agents using or documenting this repository's NBA data CLI. Focus on what every command does, how to call it, and the JSON shapes returned from NBA live data, NBA Stats, DARKO, and ESPN public APIs.
---

# NBA CLI Usage Skill

This repository provides an `nba` command for fetching NBA live data, official NBA Stats data, lineup data, DARKO projections, and ESPN public data such as contracts, salaries, injuries, bios, rosters, news, scoreboards, standings, and game summaries. When using or documenting it, focus on the data workflow: choose a command, pass IDs/names as needed, and consume JSON.

## Core behavior

- All command output is JSON printed to stdout.
- Do not add or document table/pretty-table/non-JSON output modes.
- IDs vary by data source:
  - NBA team IDs look like `1610612747`.
  - NBA player IDs look like `2544`.
  - NBA game IDs look like `0022400123`.
  - ESPN team/player/game IDs are shorter ESPN IDs, for example an ESPN event ID like `401705298`.
- Name lookup commands are helpers for turning human-readable names into IDs before calling stats, roster, contract, injury, or profile commands.
- Most NBA Stats commands return the raw NBA Stats API envelope, commonly with `resource`, `parameters`, and `resultSets`/`resultSet` containing `headers` and `rowSet` arrays.
- NBA live commands return JSON from `cdn.nba.com`; some commands are raw feed passthroughs and `live-game-summary` is a compact normalized summary.
- ESPN commands return JSON from ESPN public site/web/core APIs. Some ESPN core endpoints return `$ref` links; commands such as `espn-player-contracts` dereference contract items when listing all seasons.
- DARKO commands return parsed projection objects from `darko.app`; this source is unofficial and can break if the frontend payload changes.

## How to run commands

Use:

```bash
nba <command> [args]
```

When working from an uninstalled checkout, use the same commands via:

```bash
cd cli
go run . <command> [args]
```

Prefer examples that demonstrate what data comes back rather than install/update steps. If installation must be mentioned, use the prebuilt GitHub release binary or the existing README instructions.

## Common workflow

1. Resolve a team or player name to the ID for the data source you plan to use.
2. Use NBA IDs with official NBA Stats, schedule, roster, lineup, or game-log commands.
3. Use live commands for today's NBA live feed and NBA game IDs.
4. Use ESPN commands for contracts/salaries, injuries, bios, ESPN rosters, ESPN scoreboards, news, standings, and ESPN game summaries.
5. Use DARKO commands for projection-style player metrics.

Examples:

```bash
nba get-team-id "Los Angeles Lakers"
nba team-roster 1610612747
nba league-schedule 1610612747
nba league-dash-lineups 1610612747 Advanced 2
nba get-espn-player-id "LeBron James"
nba espn-player-contracts "LeBron James" 2026
nba espn-injuries lakers
```

## Commands and return shapes

### Utilities

```bash
nba update [install_path]
```

Fetches the latest source from GitHub or `NBA_CLI_UPDATE_SOURCE` / `NBA_CLI_UPDATE_URL`, rebuilds the Go CLI, and installs it to `install_path`, `NBA_CLI_INSTALL_PATH`, the current `nba` executable path, or `~/.local/bin/nba`. Returns JSON with fields such as `updated`, `changed`, `installedPath`, `source`, `currentExecutable`, and SHA-256 checksums.

### Name / ID lookup

```bash
nba get-team-id <team name>
nba get-player-id <player name>
nba get-espn-team-id <team name>
nba get-espn-player-id <player name>
```

- `get-team-id` returns an NBA team ID JSON number from the built-in NBA team map.
- `get-player-id` searches official NBA Stats `commonallplayers` for the current season context and returns an NBA player ID.
- `get-espn-team-id` searches ESPN teams by display name, short name, abbreviation, location, or slug and returns an ESPN team ID.
- `get-espn-player-id` searches ESPN player results and returns an ESPN player ID string.

### Live games

```bash
nba live-scores
```

Returns the raw daily NBA live scoreboard JSON, including `scoreboard.games[]` with game IDs, teams, scores, status, period, clock, and leaders.

```bash
nba live-game-summary
```

Returns a compact array like:

```json
[
  {
    "gameId": "0022400123",
    "status": "Final",
    "period": 4,
    "clock": "",
    "away": "LAL 110",
    "home": "DEN 117",
    "leaders": { }
  }
]
```

```bash
nba live-boxscore <game_id>
nba live-play-by-play <game_id>
```

Return raw NBA live feed JSON for one game. `live-boxscore` includes game status, teams, players, scores, and player/team box score fields. `live-play-by-play` includes ordered game actions/events.

### DARKO projections

```bash
nba darko-leaderboard [limit]
nba darko-player <name|nba_id>
```

`darko-leaderboard` returns an array sorted by `dpm` descending; default limit is `25`. `darko-player` returns a single projection object matching an NBA ID or name fragment. Fields include:

```json
{
  "nba_id": 203999,
  "player_name": "Nikola Jokic",
  "team_name": "DEN",
  "position": "C",
  "date": "...",
  "season": 2025,
  "dpm": 7.1,
  "o_dpm": 6.3,
  "d_dpm": 0.8,
  "box_dpm": 6.9,
  "on_off_dpm": 7.4,
  "x_minutes": 34.0,
  "x_pace": 99.0,
  "x_pts_100": 32.0,
  "x_ast_100": 11.0,
  "x_fg_pct": 0.6,
  "x_fg3_pct": 0.4,
  "x_ft_pct": 0.8
}
```

Values vary by DARKO payload and may be `null` when unavailable.

### Official NBA player stats and profiles

```bash
nba player-profile <player_id>
nba player-awards <player_id>
nba player-career-stats <player_id>
nba player-game-log <player_id>
nba league-dash-player-stats <team_id> [Advanced|Basic]
nba league-leaders [PTS|REB|AST|...] [PerGame|Totals|Per48|...]
nba player-estimated-metrics
```

All return official NBA Stats JSON envelopes with `resultSets`/`resultSet`, `headers`, and `rowSet` data unless the NBA endpoint shape differs.

- `player-profile`: `commonplayerinfo` profile/bio fields for one NBA player ID.
- `player-awards`: `playerawards` award history for one NBA player ID.
- `player-career-stats`: `playercareerstats` career and season-level per-game rows for one player.
- `player-game-log`: `playergamelog` current regular-season game-by-game rows for one player.
- `league-dash-player-stats`: `leaguedashplayerstats` team-filtered dashboard stats; default measure type is `Advanced`, optional `Basic` is accepted.
- `league-leaders`: `leagueleaders` leaderboard rows; default stat category is `PTS`, default per mode is `PerGame`.
- `player-estimated-metrics`: `playerestimatedmetrics` for the current regular season.

Use `headers` to interpret each `rowSet` row.

### Official NBA team, league, and history stats

```bash
nba team-standings
nba team-game-log <team_id>
nba team-roster <team_id>
nba team-details <team_id>
nba league-dash-team-stats [Advanced|Basic]
nba draft-history [year]
nba franchise-history
```

All return official NBA Stats JSON envelopes unless noted.

- `team-standings`: `leaguestandingsv3` current regular-season standings.
- `team-game-log`: `teamgamelog` current regular-season game logs for one NBA team ID.
- `team-roster`: `commonteamroster` current-season roster rows for one NBA team ID.
- `team-details`: `teamdetails` arena, ownership, social, and franchise background fields for one NBA team ID.
- `league-dash-team-stats`: `leaguedashteamstats` metrics for all NBA teams; default measure type is `Advanced`.
- `draft-history`: `drafthistory` for the supplied four-digit year or the current season's starting year by default.
- `franchise-history`: `franchisehistory` historical NBA franchise records.

### Schedule, game stats, and lineups

```bash
nba league-schedule <team_id>
nba game-win-probability <game_id>
nba game-boxscore <game_id>
nba league-dash-lineups <team_id> [Advanced|Basic] [2-5]
```

- `league-schedule`: calls NBA Stats `scheduleleaguev2`, filters to games involving the requested NBA team ID, and returns an array of game objects with game ID/date/status and home/away team objects.
- `game-win-probability`: calls NBA Stats `winprobabilitypbp` with `RunType=each second`; if unavailable, returns JSON with `gameId`, `error`, and `detail` instead of crashing.
- `game-boxscore`: calls NBA Stats `boxscoreadvancedv3` and returns advanced box score JSON for one NBA game ID.
- `league-dash-lineups`: calls NBA Stats `leaguedashlineups` for one NBA team ID; default measure type is `Advanced`, default lineup/group size is `2`, and the group size can be `2` through `5`.

### ESPN player, contract, and salary data

```bash
nba espn-search-player <name>
nba espn-player <id|name>
nba espn-player-bio <id|name>
nba espn-player-stats <id|name>
nba espn-player-contracts <id|name> [season_year]
```

- `espn-search-player`: returns ESPN common search JSON for NBA players, defaulting to a limit of 10 inside the command implementation.
- `espn-player`: resolves a name or accepts an ESPN player ID and returns ESPN athlete card/profile data.
- `espn-player-bio`: returns ESPN biographical extras and awards data for the player.
- `espn-player-stats`: returns ESPN core athlete statistics/splits JSON.
- `espn-player-contracts`: returns ESPN contract/salary data. With a four-digit `season_year`, returns `{ "athleteId", "season", "contract" }` for that season. Without a season, returns `{ "athleteId", "count", "contracts" }` and dereferences each ESPN contract `$ref`; failed dereferences are represented as objects with `ref` and `error`.

Example:

```bash
nba espn-player-contracts "LeBron James" 2026
```

### ESPN teams, rosters, injuries, news, games, and standings

```bash
nba espn-teams
nba espn-team-roster <id|name>
nba espn-injuries [id|name]
nba espn-news [limit]
nba espn-scoreboard [YYYYMMDD]
nba espn-game-summary <espn_game_id>
nba espn-standings
```

- `espn-teams`: returns the ESPN NBA teams index from the public site API.
- `espn-team-roster`: resolves an ESPN team ID/name and returns the ESPN roster with player biographical data.
- `espn-injuries`: with no argument, returns league-wide ESPN injury reports; with a team ID/name, filters to that team's injury object or returns `{ "teamId", "injuries": [] }` when no injuries are listed.
- `espn-news`: returns ESPN NBA news feed JSON; default limit is `10` and an optional positive integer overrides it.
- `espn-scoreboard`: returns ESPN scoreboard/events JSON; optional date format is `YYYYMMDD`.
- `espn-game-summary`: returns ESPN summary JSON for an ESPN event ID, including game summary sections such as box score, leaders, plays, odds, header, competitions, and related data when ESPN provides them.
- `espn-standings`: returns ESPN standings hierarchy sorted by win percentage.

## Data source notes

- Live data comes from `https://cdn.nba.com/static/json/liveData`.
- NBA Stats data comes from `https://stats.nba.com/stats` and needs browser-like headers. Missing default parameters can produce `500`, `EOF`, or timeout errors.
- ESPN site data comes from `https://site.api.espn.com/apis/site/v2/sports/basketball/nba`.
- ESPN web/common data comes from `https://site.web.api.espn.com/apis`.
- ESPN core data comes from `https://sports.core.api.espn.com/v2/sports/basketball/leagues/nba`.
- DARKO data is parsed from the embedded Svelte payload at `https://www.darko.app/`.

## Maintenance notes for agents

- If command behavior changes, update this skill plus `README.md` and `cli/README.md`.
- Every CLI command in `cli/main.go` should be documented in this skill to at least describe arguments and returned JSON at a high level.
- Keep output JSON-only.
- Run `gofmt` and `go test ./...` from `cli/` after code changes. Documentation-only changes do not require Go formatting.
