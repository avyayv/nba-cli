from nba_api.stats.endpoints import playercareerstats, playergamelog, leaguedashplayerstats
from utils import convert_result_sets_to_tables

def player_career_stats(player_id: int) -> dict:
    """
    Get the career stats for a player id
    """
    return playercareerstats.PlayerCareerStats(player_id=player_id).get_dict()

def player_game_log(player_id: int) -> dict:
    """
    Get the game log for a player id
    """
    return playergamelog.PlayerGameLog(player_id=player_id).get_dict()

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