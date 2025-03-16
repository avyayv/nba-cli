from typing import Any, Dict, List
from mcp.server.fastmcp import FastMCP
from nba_api.stats.endpoints import (
    playercareerstats, 
    playergamelog, leaguestandings, 
    commonteamroster, 
    teamgamelog, 
    winprobabilitypbp,
    boxscoreadvancedv3,
    boxscoresummaryv2, 
    boxscoreusagev3,
    boxscoredefensivev2,
    leaguedashlineups,
    leaguedashteamstats,
    leaguedashplayerstats,
)

from nba_api.live.nba.endpoints import scoreboard  

from nba_api.stats.static import players, teams

# Initialize FastMCP server
mcp = FastMCP("nba")

@mcp.tool(name="get_player_id")
def get_player_id(player_name: str) -> str:
    """
    Get the player id for a player name
    """
    player_dict = players.find_players_by_full_name(player_name)
    return player_dict[0]['id']

@mcp.tool(name="get_team_id")
def get_team_id(team_name: str) -> str:
    """
    Get the team id for a team name. i.e "Los Angeles Lakers"
    """
    team_dict = teams.find_teams_by_full_name(team_name)
    return team_dict[0]['id']

@mcp.tool(name="player_career_stats")
def player_career_stats(player_id: int) -> dict:
    """
    Get the career stats for a player id
    """
    return playercareerstats.PlayerCareerStats(player_id=player_id).get_dict()

@mcp.tool(name="live_scores")
def live_scores() -> dict:
    """
    Get the live scores for all NBA games
    """
    return scoreboard.ScoreBoard().get_dict()

@mcp.tool(name="player_game_log")
def player_game_log(player_id: int) -> dict:
    """
    Get the game log for a player id
    """
    return playergamelog.PlayerGameLog(player_id=player_id).get_dict()

@mcp.tool(name="team_standings")
def team_standings() -> dict:
    """
    Get the team standings for all NBA teams
    """
    return leaguestandings.LeagueStandings().get_dict()

@mcp.tool(name="team_game_log")
def team_game_log(team_id: int) -> dict:
    """
    Get the game log for a team id
    """
    return teamgamelog.TeamGameLog(team_id=team_id).get_dict()

@mcp.tool(name="team_roster")
def team_roster(team_id: int) -> dict:
    """
    Get the roster for a team id
    """
    return commonteamroster.CommonTeamRoster(team_id=team_id).get_dict()

@mcp.tool(name="game_win_probability")
def game_win_probability(game_id: str) -> dict:
    """
    Get the win probability for a game id
    """
    data = winprobabilitypbp.WinProbabilityPBP(game_id=game_id).get_dict()
    result = []
    for row in data['resultSets'][0]['rowSet']:
        if row[1] is None:
            continue
        print(row)
        try: 
            line = f"Period {row[7]} ({row[13]}) | Home: {row[4]} pts, {row[2]} win prob | Away: {row[5]} pts, {row[3]} | {row[11]}"
            result.append(line)
        except Exception as e:
            continue
    return "\n".join(result)

def convert_result_sets_to_tables(data: Dict[str, Any]) -> str:
    """
    Helper function to convert NBA API result sets into formatted table strings.
    
    Args:
        data: Dictionary containing NBA API response with 'resultSets' key
        
    Returns:
        Formatted string containing all tables from the result sets
    """
    result = []
    if 'resultSets' not in data:
        return ""
        
    for result_set in data['resultSets']:
        headers = result_set['headers']
        rows = result_set['rowSet']
        if not rows:
            continue
            
        # Filter out rank columns
        rank_indices = [i for i, h in enumerate(headers) if 'RANK' in str(h)]
        filtered_headers = [h for i, h in enumerate(headers) if i not in rank_indices]
        filtered_rows = [[val for i, val in enumerate(row) if i not in rank_indices] for row in rows]
            
        # Create header row
        result.append("\n" + result_set['name'])
        result.append("-" * len(result_set['name']))
        result.append(" | ".join(str(h) for h in filtered_headers))
        result.append("-" * len(" | ".join(str(h) for h in filtered_headers)))
        
        # Add data rows
        for row in filtered_rows:
            result.append(" | ".join(str(val) if val is not None else "" for val in row))
    return "\n".join(result)

@mcp.tool(name="game_boxscore")
def game_boxscore(game_id: str) -> dict:
    """
    Get the boxscore for a game id
    """
    advanced = boxscoreadvancedv3.BoxScoreAdvancedV3(game_id=game_id).get_dict()
    summary = boxscoresummaryv2.BoxScoreSummaryV2(game_id=game_id).get_dict()
    usage = boxscoreusagev3.BoxScoreUsageV3(game_id=game_id).get_dict()
    defensive = boxscoredefensivev2.BoxScoreDefensiveV2(game_id=game_id).get_dict()
    
    result = []
    for data in [advanced, summary, usage, defensive]:
        table = convert_result_sets_to_tables(data)
        if table:
            result.append(table)
                
    return "\n".join(result)

@mcp.tool(name="league_dash_player_stats") 
def league_dash_player_stats(team_id: int, measure_type: str = 'Advanced') -> dict:
    """
    Get advanced stats for all players on a certain team. 

    measure_type can be "Advanced" or "Basic"
    """

    data = leaguedashplayerstats.LeagueDashPlayerStats(
        per_mode_detailed='Per100Possessions',
        measure_type_detailed_defense=measure_type,
        team_id_nullable=team_id
    ).get_dict()
    
    return convert_result_sets_to_tables(data)

@mcp.tool(name="league_dash_team_stats")
def league_dash_team_stats(measure_type: str = 'Advanced') -> dict:
    """
    Get the league dash team stats for all NBA teams. This should be called for every query on any specific 
    query to figure out what is good and what is bad.  

    If you want to get the stats for a specific team, you can use the team_id parameter.
    """
    data = leaguedashteamstats.LeagueDashTeamStats(
        per_mode_detailed='Per100Possessions',
        measure_type_detailed_defense=measure_type,
    ).get_dict()
    
    return convert_result_sets_to_tables(data)

@mcp.tool(name="league_dash_lineups")
def league_dash_lineups(team_id: int, measure_type: str = 'Advanced', lineup_count: int = 2) -> dict:

    """
    Get the league dash lineups for all NBA teams. This should be called for every query on any specific 
    query to figure out what is good and what is bad. 

    If you are trying to find things like how two players play together, you should use this tool and say lineup_count=2 and see if you can find.

    measure_type_detailed_defense can be "Advanced" or "Basic"

    team_id is the id of the team you want to get the lineups for. i.e "Los Angeles Lakers"
    """
    data = leaguedashlineups.LeagueDashLineups(
        group_quantity=lineup_count,
        per_mode_detailed='Per100Possessions',
        measure_type_detailed_defense=measure_type,
        team_id_nullable=team_id
    ).get_dict()
    
    return convert_result_sets_to_tables(data)

if __name__ == "__main__":
    mcp.run(transport='stdio')

