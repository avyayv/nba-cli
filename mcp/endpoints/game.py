from nba_api.stats.endpoints import (
    winprobabilitypbp,
    boxscoreadvancedv3,
    boxscoresummaryv2,
    boxscoreusagev3,
    boxscoredefensivev2,
    scheduleleaguev2
)
from nba_api.live.nba.endpoints import scoreboard, boxscore, playbyplay
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


def live_game_summary() -> list:
    """
    Get a compact live summary for today's games: status, clock, score, leaders, and game ids.
    """
    data = scoreboard.ScoreBoard().get_dict()
    games = data.get('scoreboard', {}).get('games', [])
    result = []
    for game in games:
        home = game.get('homeTeam', {})
        away = game.get('awayTeam', {})
        leaders = game.get('gameLeaders') or {}
        away_leader = leaders.get('awayLeaders') or {}
        home_leader = leaders.get('homeLeaders') or {}
        result.append({
            'gameId': game.get('gameId'),
            'status': game.get('gameStatusText'),
            'period': game.get('period'),
            'clock': game.get('gameClock'),
            'away': f"{away.get('teamTricode')} {away.get('score')}",
            'home': f"{home.get('teamTricode')} {home.get('score')}",
            'awayLeader': away_leader,
            'homeLeader': home_leader,
        })
    return result


def live_boxscore(game_id: str) -> dict:
    """
    Get the live NBA box score for a game id from the live data feed.
    """
    return boxscore.BoxScore(game_id=game_id).get_dict()


def live_play_by_play(game_id: str) -> dict:
    """
    Get live play-by-play actions for a game id from the live data feed.
    """
    return playbyplay.PlayByPlay(game_id=game_id).get_dict()

def game_win_probability(game_id: str) -> dict:
    """
    Get the win probability for a game id.
    """
    try:
        data = winprobabilitypbp.WinProbabilityPBP(game_id=game_id).get_dict()
    except Exception as e:
        return {
            "gameId": game_id,
            "error": "Win probability feed unavailable for this game.",
            "detail": str(e),
        }

    result = []
    for row in data.get('resultSets', [{}])[0].get('rowSet', []):
        if row[1] is None:
            continue
        try:
            line = f"Period {row[7]} ({row[13]}) | Home: {row[4]} pts, {row[2]} win prob | Away: {row[5]} pts, {row[3]} | {row[11]}"
            result.append(line)
        except Exception:
            continue
    return "\n".join(result) if result else {"gameId": game_id, "error": "No win probability rows returned for this game."}

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
