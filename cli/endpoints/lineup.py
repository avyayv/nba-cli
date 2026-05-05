from nba_api.stats.endpoints import leaguedashlineups
from utils import convert_result_sets_to_tables

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