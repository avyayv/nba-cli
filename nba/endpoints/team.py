from nba_api.stats.endpoints import (
    leaguestandings,
    commonteamroster,
    teamgamelog,
    leaguedashteamstats
)
from utils import convert_result_sets_to_tables

def team_standings() -> dict:
    """
    Get the team standings for all NBA teams. This only gets the current day's standings, whether they have happened or not.
    """
    return leaguestandings.LeagueStandings().get_dict()

def team_game_log(team_id: int) -> dict:
    """
    Get the game log for a team id
    """
    return teamgamelog.TeamGameLog(team_id=team_id).get_dict()

def team_roster(team_id: int) -> dict:
    """
    Get the roster for a team id.
    """
    return commonteamroster.CommonTeamRoster(team_id=team_id).get_dict()

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