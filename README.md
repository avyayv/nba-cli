# nba-cli

Standalone Go command-line interface for NBA live scores, stats, lineup data, and DARKO projections. All command output is JSON. No Python runtime is required.

## Install

```bash
git clone https://github.com/avyayv/nba-cli
cd nba-cli/cli
mkdir -p ~/.local/bin
go build -o ~/.local/bin/nba .
chmod +x ~/.local/bin/nba
nba --help
```

## Update

```bash
nba update
```

`nba update` fetches the latest source from GitHub, rebuilds the Go CLI, and replaces the installed `nba` binary. Pass an explicit target path if needed: `nba update ~/.local/bin/nba`.

## Local development

```bash
git clone https://github.com/avyayv/nba-cli
cd nba-cli/cli
go build -o nba .
./nba --help
```

## Usage

```bash
# Resolve names to IDs
nba get-player-id "LeBron James"
nba get-team-id "Los Angeles Lakers"

# Player stats
nba player-career-stats 2544
nba player-game-log 2544
nba league-dash-player-stats 1610612747 --measure-type Advanced

# Team stats
nba team-standings
nba team-game-log 1610612747
nba team-roster 1610612747
nba league-dash-team-stats --measure-type Advanced

# Games / live data
nba league-schedule 1610612747
nba live-scores
nba live-game-summary
nba live-boxscore 0022400123
nba live-play-by-play 0022400123
nba game-win-probability 0022400123
nba game-boxscore 0022400123

# DARKO projections
nba darko-leaderboard 10
nba darko-player "Nikola Jokic"

# Keep the installed CLI current
nba update

# Lineups
nba league-dash-lineups 1610612747 --lineup-count 2
```

Run `nba --help` or `nba <command> --help` for the full list.

## Available commands

### Utilities
- `update` — fetch, rebuild, and install the latest CLI
- `get-player-id` — convert a player name to an ID
- `get-team-id` — convert a team name to an ID

### Player
- `player-career-stats` — career statistics for a player
- `player-game-log` — game-by-game logs for a player
- `league-dash-player-stats` — stats for all players on a team

### Team
- `team-standings` — current NBA standings
- `team-game-log` — game log for a team
- `team-roster` — roster for a team
- `league-dash-team-stats` — stats for all NBA teams

### Game / live
- `league-schedule` — schedule for a team's games
- `live-scores` — current live scores
- `live-game-summary` — compact live score/status/leaders for today's games
- `live-boxscore` — live box score for a game
- `live-play-by-play` — live play-by-play for a game
- `game-win-probability` — win probability data for a game
- `game-boxscore` — box score for a game

### Lineup
- `league-dash-lineups` — lineup statistics and combinations

### DARKO
- `darko-leaderboard` — current DARKO DPM leaderboard from darko.app
- `darko-player` — DARKO projection metrics by NBA ID or name fragment
