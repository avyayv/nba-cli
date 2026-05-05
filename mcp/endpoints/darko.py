import re
from typing import Any, Dict, List, Optional

import httpx

DARKO_URL = "https://www.darko.app/"
DEFAULT_FIELDS = [
    "nba_id", "player_name", "team_name", "position", "date", "season",
    "dpm", "o_dpm", "d_dpm", "box_dpm", "on_off_dpm", "x_minutes",
    "x_pace", "x_pts_100", "x_ast_100", "x_fg_pct", "x_fg3_pct", "x_ft_pct",
]


def _extract_player_objects(html: str) -> List[str]:
    marker = "players:["
    start = html.find(marker)
    if start == -1:
        return []
    i = start + len(marker)
    depth = 0
    obj_start: Optional[int] = None
    objects: List[str] = []
    while i < len(html):
        ch = html[i]
        if ch == "{" and depth == 0:
            obj_start = i
            depth = 1
        elif ch == "{" and depth > 0:
            depth += 1
        elif ch == "}" and depth > 0:
            depth -= 1
            if depth == 0 and obj_start is not None:
                objects.append(html[obj_start:i + 1])
                obj_start = None
        elif ch == "]" and depth == 0:
            break
        i += 1
    return objects


def _field(obj: str, name: str) -> Any:
    m = re.search(rf'{re.escape(name)}:("[^"]*"|null|-?\.\d+|-?\d+(?:\.\d+)?)', obj)
    if not m:
        return None
    raw = m.group(1)
    if raw == "null":
        return None
    if raw.startswith('"'):
        return raw[1:-1]
    if raw.startswith(".") or raw.startswith("-."):
        raw = raw.replace(".", "0.", 1) if raw.startswith(".") else raw.replace("-.", "-0.", 1)
    try:
        return float(raw) if "." in raw else int(raw)
    except ValueError:
        return raw


def _load_darko_players() -> List[Dict[str, Any]]:
    html = httpx.get(DARKO_URL, timeout=30).text
    players: List[Dict[str, Any]] = []
    for obj in _extract_player_objects(html):
        row = {field: _field(obj, field) for field in DEFAULT_FIELDS}
        if row.get("player_name"):
            players.append(row)
    return players


def darko_leaderboard(limit: int = 25) -> List[Dict[str, Any]]:
    """Get the current DARKO DPM leaderboard from darko.app."""
    return _load_darko_players()[:limit]


def darko_player(player_name_or_id: str) -> Dict[str, Any]:
    """Get DARKO projection metrics for a player by NBA id or name fragment."""
    needle = str(player_name_or_id).lower()
    for player in _load_darko_players():
        if str(player.get("nba_id")) == needle or needle in str(player.get("player_name", "")).lower():
            return player
    return {"error": f"No DARKO player found for {player_name_or_id}"}
