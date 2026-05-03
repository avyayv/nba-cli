from nba_api.stats.endpoints import (
    winprobabilitypbp,
    boxscoreadvancedv3,
    boxscoresummaryv2,
    boxscoreusagev3,
    boxscoredefensivev2,
    scheduleleaguev2
)
from nba_api.live.nba.endpoints import scoreboard
from utils import convert_result_sets_to_tables

def league_schedule(team_id: int) -> dict:
    """
    Get the schedule for all NBA games for a team id. This includes games that haven't happened yet.
    """
    data = scheduleleaguev2.ScheduleLeagueV2().get_dict()

    # now iterate through all games and only pick when homeTeam_teamId or awayTeam_teamId is team_id
    filtered_games = []
    for row in data['leagueSchedule']['gameDates']:
        for game in row['games']:
            home_team_id = game['homeTeam']['teamId']
            away_team_id = game['awayTeam']['teamId']
            if home_team_id == team_id or away_team_id == team_id:
                # Format the game information with team names, scores, seeds, and records
                home_team_tricode = game['homeTeam']['teamTricode']
                home_team_wins = game['homeTeam']['wins']
                home_team_losses = game['homeTeam']['losses']
                home_team_score = game['homeTeam']['score']
                home_team_seed = game['homeTeam']['seed'] if game['homeTeam']['seed'] is not None else "N/A"
            
                away_team_tricode = game['awayTeam']['teamTricode']
                away_team_wins = game['awayTeam']['wins']
                away_team_losses = game['awayTeam']['losses']
                away_team_score = game['awayTeam']['score']
                away_team_seed = game['awayTeam']['seed'] if game['awayTeam']['seed'] is not None else "N/A"
            
                game_info = f"[{game['gameDateEst']}] [{game['gameStatusText']}] [{game['gameId']}] {away_team_score if away_team_score is not None else '-'} {away_team_tricode} (seed #{away_team_seed}, {away_team_wins}-{away_team_losses}) @ {home_team_score if home_team_score is not None else '-'} {home_team_tricode} (seed #{home_team_seed}, {home_team_wins}-{home_team_losses})"
                filtered_games.append(game_info)
    
    return filtered_games


def live_scores() -> dict:
    """
    Get the live scores for all NBA games. This only gets the current day's games, whether they have happened or not.
    """
    return scoreboard.ScoreBoard().get_dict()

def game_win_probability(game_id: str) -> dict:
    """
    Get the win probability for a game id.
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
