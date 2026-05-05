# NBA Tools

NBA stats — player, team, game, and lineup data — exposed two ways:

- **[`mcp/`](./mcp)** — an MCP server for use with Claude Desktop, Claude Code, and any MCP-compatible client.
- **[`cli/`](./cli)** — a standalone Go command-line tool. It does not require Python.

The MCP server uses `nba-api`; the Go CLI calls NBA Stats/live JSON feeds directly.

## MCP Server

### Claude Desktop

Add this to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "nba": {
      "command": "uvx",
      "args": ["--from", "git+https://github.com/avyayv/nba-mcp-server#subdirectory=mcp", "nba-mcp-server"]
    }
  }
}
```

### Claude Code

```bash
claude mcp add -s user nba -- uvx --from "git+https://github.com/avyayv/nba-mcp-server#subdirectory=mcp" nba-mcp-server
```

### Local development

```bash
git clone https://github.com/avyayv/nba-mcp-server
cd nba-mcp-server/mcp
uv sync
uv run nba-mcp-server
```

## CLI

### Install

```bash
go install github.com/avyayv/nba-mcp-server/cli@latest
nba --help
```

Or for local development:

```bash
git clone https://github.com/avyayv/nba-mcp-server
cd nba-mcp-server/cli
go build -o nba .
./nba --help
```

### Usage

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

# Games
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

# Lineups
nba league-dash-lineups 1610612747 --lineup-count 2

# All CLI output is JSON
nba player-career-stats 2544
```

Run `nba --help` or `nba <command> --help` for the full list.

## Available endpoints

Both the MCP server and CLI expose the same set of operations:

### Utilities
- `get_player_id` — convert a player name to an ID
- `get_team_id` — convert a team name to an ID

### Player
- `player_career_stats` — career statistics for a player
- `player_game_log` — game-by-game logs for a player
- `league_dash_player_stats` — advanced stats for all players on a team

### Team
- `team_standings` — current NBA standings
- `team_game_log` — game log for a team
- `team_roster` — roster for a team
- `league_dash_team_stats` — advanced stats for all NBA teams

### Game / live
- `league_schedule` — schedule for a team's games
- `live_scores` — current live scores
- `live_game_summary` — compact live score/status/leaders for today's games
- `live_boxscore` — live box score for a game
- `live_play_by_play` — live play-by-play for a game
- `game_win_probability` — win probability data for a game
- `game_boxscore` — box score for a game

### Lineup
- `league_dash_lineups` — lineup statistics and combinations

### DARKO
- `darko_leaderboard` — current DARKO DPM leaderboard from darko.app
- `darko_player` — DARKO projection metrics by NBA id or name fragment
