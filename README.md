# NBA MCP Server

MCP server that wraps the NBA API for player, team, game, and lineup statistics.

## Installation

### Claude Desktop

Add this to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "nba": {
      "command": "uvx",
      "args": ["--from", "git+https://github.com/avyayv/nba-mcp-server", "nba-mcp-server"]
    }
  }
}
```

### Claude Code

```bash
claude mcp add -s user nba -- uvx --from git+https://github.com/avyayv/nba-mcp-server nba-mcp-server
```

### Local Development

```bash
git clone https://github.com/avyayv/nba-mcp-server
cd nba-mcp-server
uv sync
uv run nba-mcp-server
```

## Available Tools

### Utility Functions
- `get_player_id` - Convert player name to ID
- `get_team_id` - Convert team name to ID

### Player Endpoints
- `player_career_stats` - Career statistics for a player
- `player_game_log` - Game-by-game logs for a player
- `league_dash_player_stats` - Advanced stats for all players on a team

### Team Endpoints
- `team_standings` - Current NBA standings
- `team_game_log` - Game log for a team
- `team_roster` - Roster for a team
- `league_dash_team_stats` - Advanced stats for all NBA teams

### Game Endpoints
- `league_schedule` - Schedule for a team's games
- `live_scores` - Current live scores
- `game_win_probability` - Win probability data for a game
- `game_boxscore` - Box score for a game

### Lineup Endpoints
- `league_dash_lineups` - Lineup statistics and combinations
