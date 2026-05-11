# nba-cli

Standalone Go command-line interface for NBA live scores, stats, lineup data, DARKO projections, and ESPN public data including contracts, injuries, news, player bios, rosters, scoreboards, and game summaries. All command output is JSON. No Python runtime is required.

## Install

### Prebuilt release binary

Every push to `main` publishes Linux, macOS, and Windows binaries to [GitHub Releases](https://github.com/avyayv/nba-cli/releases/latest).

macOS Apple Silicon:

```bash
mkdir -p ~/.local/bin
curl -fsSL https://github.com/avyayv/nba-cli/releases/latest/download/nba_darwin_arm64.tar.gz | tar -xz
mv nba ~/.local/bin/nba
chmod +x ~/.local/bin/nba
nba --help
```

Use `nba_darwin_amd64.tar.gz` for Intel Macs, `nba_linux_amd64.tar.gz` / `nba_linux_arm64.tar.gz` for Linux, or the Windows `.zip` assets.

### From source

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
nba league-dash-player-stats 1610612747 Advanced

# Team stats
nba team-standings
nba team-game-log 1610612747
nba team-roster 1610612747
nba league-dash-team-stats Advanced

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

# ESPN public data
nba espn-search-player "LeBron James"
nba get-espn-player-id "LeBron James"
nba espn-player-contracts "LeBron James" 2026
nba espn-player-bio "Nikola Jokic"
nba espn-player-stats "Nikola Jokic"
nba espn-injuries lakers
nba espn-news 5
nba espn-scoreboard 20260511
nba espn-game-summary 401705298
nba espn-standings

# Keep the installed CLI current
nba update

# Lineups
nba league-dash-lineups 1610612747 Advanced 2
```

Run `nba --help` for the full list.

## Available commands

### Utilities
- `update` — fetch, rebuild, and install the latest CLI
- `get-player-id` — convert a player name to an ID
- `get-team-id` — convert a team name to an ID

### Player
- `player-profile` — official NBA common player profile and bio fields
- `player-awards` — official NBA award history for a player
- `player-career-stats` — career statistics for a player
- `player-game-log` — game-by-game logs for a player
- `league-dash-player-stats` — stats for all players on a team
- `league-leaders` — official NBA leaderboards by stat category
- `player-estimated-metrics` — official NBA estimated metrics

### Team / league
- `team-standings` — current NBA standings
- `team-game-log` — game log for a team
- `team-roster` — roster for a team
- `team-details` — arena, ownership, and franchise background
- `league-dash-team-stats` — stats for all NBA teams
- `draft-history` — official NBA draft history by year
- `franchise-history` — historical NBA franchise records

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

### ESPN public data
- `espn-search-player` / `get-espn-player-id` — resolve ESPN NBA player IDs
- `get-espn-team-id` — resolve ESPN team IDs
- `espn-player` — ESPN player card/profile data
- `espn-player-bio` — ESPN awards and biographical extras
- `espn-player-stats` — ESPN player statistics and splits
- `espn-player-contracts` — ESPN contract/salary data by player, optionally by season year
- `espn-teams` / `espn-team-roster` — ESPN team index and rosters
- `espn-injuries` — league-wide or team-specific injury reports
- `espn-news` — ESPN NBA news feed
- `espn-scoreboard` — ESPN scoreboard/events feed
- `espn-game-summary` — ESPN game summary, box score, leaders, plays, odds, and related data
- `espn-standings` — ESPN standings hierarchy
