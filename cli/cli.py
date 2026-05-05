import json
from typing import Any

import click

from utils import (
    convert_result_sets_to_tables,
    get_player_id as _get_player_id,
    get_team_id as _get_team_id,
)
from endpoints import (
    player_career_stats as _player_career_stats,
    player_game_log as _player_game_log,
    league_dash_player_stats as _league_dash_player_stats,
    team_standings as _team_standings,
    team_game_log as _team_game_log,
    team_roster as _team_roster,
    league_dash_team_stats as _league_dash_team_stats,
    league_schedule as _league_schedule,
    live_scores as _live_scores,
    game_win_probability as _game_win_probability,
    game_boxscore as _game_boxscore,
    league_dash_lineups as _league_dash_lineups,
)


def _format(result: Any, raw: bool) -> str:
    if raw:
        return json.dumps(result, indent=2, default=str)
    if isinstance(result, str):
        return result
    if isinstance(result, list):
        return "\n".join(str(item) for item in result)
    if isinstance(result, dict) and "resultSets" in result:
        return convert_result_sets_to_tables(result)
    return json.dumps(result, indent=2, default=str)


def _emit(ctx: click.Context, result: Any) -> None:
    click.echo(_format(result, ctx.obj["raw"]))


MEASURE_TYPES = click.Choice(["Advanced", "Basic"], case_sensitive=False)


@click.group(help="NBA stats CLI — query player, team, game, and lineup data.")
@click.option("--raw", is_flag=True, help="Print raw JSON instead of formatted tables.")
@click.pass_context
def cli(ctx: click.Context, raw: bool) -> None:
    ctx.ensure_object(dict)
    ctx.obj["raw"] = raw


@cli.command("get-player-id", help="Resolve a player name to an NBA player ID.")
@click.argument("player_name")
@click.pass_context
def get_player_id(ctx: click.Context, player_name: str) -> None:
    _emit(ctx, _get_player_id(player_name))


@cli.command("get-team-id", help='Resolve a team name to an NBA team ID (e.g. "Los Angeles Lakers").')
@click.argument("team_name")
@click.pass_context
def get_team_id(ctx: click.Context, team_name: str) -> None:
    _emit(ctx, _get_team_id(team_name))


@cli.command("player-career-stats", help="Career statistics for a player.")
@click.argument("player_id", type=int)
@click.pass_context
def player_career_stats(ctx: click.Context, player_id: int) -> None:
    _emit(ctx, _player_career_stats(player_id))


@cli.command("player-game-log", help="Game-by-game logs for a player.")
@click.argument("player_id", type=int)
@click.pass_context
def player_game_log(ctx: click.Context, player_id: int) -> None:
    _emit(ctx, _player_game_log(player_id))


@cli.command("league-dash-player-stats", help="Per-100-possessions stats for all players on a team.")
@click.argument("team_id", type=int)
@click.option("--measure-type", type=MEASURE_TYPES, default="Advanced", show_default=True)
@click.pass_context
def league_dash_player_stats(ctx: click.Context, team_id: int, measure_type: str) -> None:
    _emit(ctx, _league_dash_player_stats(team_id, measure_type))


@cli.command("team-standings", help="Current NBA standings.")
@click.pass_context
def team_standings(ctx: click.Context) -> None:
    _emit(ctx, _team_standings())


@cli.command("team-game-log", help="Game log for a team.")
@click.argument("team_id", type=int)
@click.pass_context
def team_game_log(ctx: click.Context, team_id: int) -> None:
    _emit(ctx, _team_game_log(team_id))


@cli.command("team-roster", help="Roster for a team.")
@click.argument("team_id", type=int)
@click.pass_context
def team_roster(ctx: click.Context, team_id: int) -> None:
    _emit(ctx, _team_roster(team_id))


@cli.command("league-dash-team-stats", help="Per-100-possessions stats for all NBA teams.")
@click.option("--measure-type", type=MEASURE_TYPES, default="Advanced", show_default=True)
@click.pass_context
def league_dash_team_stats(ctx: click.Context, measure_type: str) -> None:
    _emit(ctx, _league_dash_team_stats(measure_type))


@cli.command("league-schedule", help="Schedule for a team's games (past and upcoming).")
@click.argument("team_id", type=int)
@click.pass_context
def league_schedule(ctx: click.Context, team_id: int) -> None:
    _emit(ctx, _league_schedule(team_id))


@cli.command("live-scores", help="Today's live scoreboard.")
@click.pass_context
def live_scores(ctx: click.Context) -> None:
    _emit(ctx, _live_scores())


@cli.command("game-win-probability", help="Win probability play-by-play for a game.")
@click.argument("game_id")
@click.pass_context
def game_win_probability(ctx: click.Context, game_id: str) -> None:
    _emit(ctx, _game_win_probability(game_id))


@cli.command("game-boxscore", help="Advanced box score for a game.")
@click.argument("game_id")
@click.pass_context
def game_boxscore(ctx: click.Context, game_id: str) -> None:
    _emit(ctx, _game_boxscore(game_id))


@cli.command("league-dash-lineups", help="Lineup statistics and combinations for a team.")
@click.argument("team_id", type=int)
@click.option("--measure-type", type=MEASURE_TYPES, default="Advanced", show_default=True)
@click.option("--lineup-count", type=int, default=2, show_default=True, help="Number of players in lineup (2-5).")
@click.pass_context
def league_dash_lineups(ctx: click.Context, team_id: int, measure_type: str, lineup_count: int) -> None:
    _emit(ctx, _league_dash_lineups(team_id, measure_type, lineup_count))


def main() -> None:
    cli(obj={})


if __name__ == "__main__":
    main()
