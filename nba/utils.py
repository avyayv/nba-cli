from typing import Any, Dict
from nba_api.stats.static import players, teams

def get_player_id(player_name: str) -> str:
    """
    Get the player id for a player name
    """
    player_dict = players.find_players_by_full_name(player_name)
    return player_dict[0]['id']

def get_team_id(team_name: str) -> str:
    """
    Get the team id for a team name. i.e "Los Angeles Lakers"
    """
    team_dict = teams.find_teams_by_full_name(team_name)
    return team_dict[0]['id']

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