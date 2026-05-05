from mcp.server.fastmcp import FastMCP
from utils import get_player_id, get_team_id
from endpoints import (
    player_career_stats,
    player_game_log,
    league_dash_player_stats,
    team_standings,
    team_game_log,
    team_roster,
    league_dash_team_stats,
    league_schedule,
    live_scores,
    game_win_probability,
    game_boxscore,
    live_boxscore,
    live_play_by_play,
    live_game_summary,
    league_dash_lineups,
    darko_leaderboard,
    darko_player
)

# Initialize FastMCP server
mcp = FastMCP("nba")

# Register utility functions
mcp.tool(name="get_player_id")(get_player_id)
mcp.tool(name="get_team_id")(get_team_id)

# Register player endpoints
mcp.tool(name="player_career_stats")(player_career_stats)
mcp.tool(name="player_game_log")(player_game_log)
mcp.tool(name="league_dash_player_stats")(league_dash_player_stats)

# Register team endpoints
mcp.tool(name="team_standings")(team_standings)
mcp.tool(name="team_game_log")(team_game_log)
mcp.tool(name="team_roster")(team_roster)
mcp.tool(name="league_dash_team_stats")(league_dash_team_stats)

# Register game endpoints
mcp.tool(name="league_schedule")(league_schedule)
mcp.tool(name="live_scores")(live_scores)
mcp.tool(name="game_win_probability")(game_win_probability)
mcp.tool(name="game_boxscore")(game_boxscore)
mcp.tool(name="live_boxscore")(live_boxscore)
mcp.tool(name="live_play_by_play")(live_play_by_play)
mcp.tool(name="live_game_summary")(live_game_summary)

# Register lineup endpoints
mcp.tool(name="league_dash_lineups")(league_dash_lineups)

# Register DARKO endpoints
mcp.tool(name="darko_leaderboard")(darko_leaderboard)
mcp.tool(name="darko_player")(darko_player)

def main():
    mcp.run(transport='stdio')


if __name__ == "__main__":
    main()
