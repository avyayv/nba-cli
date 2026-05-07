---
name: nba-cli
description: Guidance for AI agents using or documenting this repository's NBA data CLI. Focus on what commands do, how to call them, and the JSON shapes they return.
---

# NBA CLI Usage Skill

This repository provides an `nba` command for fetching NBA live data, NBA Stats data, lineup data, and DARKO projections. When using or documenting it, focus on the data workflow: choose a command, pass IDs/names as needed, and consume JSON.

## Core behavior

- All command output is JSON printed to stdout.
- Do not add or document table/pretty-table/non-JSON output modes.
- IDs are NBA IDs:
  - team IDs look like `1610612747`.
  - player IDs look like `2544`.
  - game IDs look like `0022400123`.
- Name lookup commands are helpers for turning human-readable names into IDs before calling stats commands.
- Most NBA Stats commands return the raw NBA Stats API envelope, commonly with `resource`, `parameters`, and `resultSets`/`resultSet` containing `headers` and `rowSet` arrays.
- Live commands return NBA live-data JSON from `cdn.nba.com`; some commands are raw feed passthroughs and `live-game-summary` is a compact normalized summary.
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

Prefer examples that demonstrate what data comes back rather than install/update steps. If installation must be mentioned, use the Homebrew tap path for end users rather than source-build instructions.

## Common workflow

1. Resolve a team or player name to an ID.
2. Use the ID with stats, schedule, roster, lineup, or game-log commands.
3. Use live commands for today's games and game IDs.
4. Use DARKO commands for projection-style player metrics.

Example:

```bash
nba get-team-id "Los Angeles Lakers"
nba team-roster 1610612747
nba league-schedule 1610612747
nba league-dash-lineups 1610612747 Advanced 2
```

## Commands and return shapes

### Name / ID lookup

```bash
nba get-team-id <team name>
nba get-player-id <player name>
```

Returns a JSON number. Examples: `1610612747` for Lakers, `2544` for LeBron James.

### Live games

```bash
nba live-scores
```

Returns the raw daily scoreboard JSON, including `scoreboard.games[]` with game IDs, teams, scores, status, period, clock, and leaders.

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

Return raw live feed JSON for one game. Box score data includes game status, teams, players, scores, and player/team box score fields. Play-by-play data includes ordered game actions/events.

### DARKO projections

```bash
nba darko-leaderboard [limit]
nba darko-player <name|nba_id>
```

`darko-leaderboard` returns an array sorted by `dpm` descending. `darko-player` returns a single projection object. Fields include:

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

### Player stats

```bash
nba player-career-stats <player_id>
nba player-game-log <player_id>
nba league-dash-player-stats <team_id> [Advanced|Basic]
```

Returns NBA Stats JSON. Use `headers` to interpret each row in `rowSet`.

- `player-career-stats`: career and season-level player stats.
- `player-game-log`: current-season game-by-game rows for a player.
- `league-dash-player-stats`: team-filtered player dashboard stats; default measure type is `Advanced`.

### Team stats

```bash
nba team-standings
nba team-game-log <team_id>
nba team-roster <team_id>
nba league-dash-team-stats [Advanced|Basic]
```

Returns NBA Stats JSON with standings, team game logs, roster rows, or team dashboard metrics. `league-dash-team-stats` defaults to `Advanced`.

### Schedule / game stats

```bash
nba league-schedule <team_id>
nba game-win-probability <game_id>
nba game-boxscore <game_id>
```

- `league-schedule` returns an array of schedule game objects for the requested team, filtered from the league schedule feed. Each object includes game ID/date/status and home/away team objects.
- `game-win-probability` returns NBA Stats win-probability play-by-play data when available. If unavailable, it returns JSON with `gameId`, `error`, and `detail` instead of crashing.
- `game-boxscore` returns NBA Stats advanced box score JSON for the game.

### Lineups

```bash
nba league-dash-lineups <team_id> [Advanced|Basic] [2-5]
```

Returns NBA Stats lineup dashboard JSON for the team. The second argument is measure type, default `Advanced`; the third is lineup size/group quantity, default `2`.

## Data source notes

- Live data comes from `https://cdn.nba.com/static/json/liveData`.
- NBA Stats data comes from `https://stats.nba.com/stats` and needs browser-like headers. Missing default parameters can produce `500`, `EOF`, or timeout errors.
- DARKO data is parsed from the embedded Svelte payload at `https://www.darko.app/`.

## Maintenance notes for agents

- If command behavior changes, update this skill plus `README.md` and `cli/README.md`.
- Keep output JSON-only.
- Run `gofmt` and `go test ./...` from `cli/` after code changes.
