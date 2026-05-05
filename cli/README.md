# nba-cli

Standalone Go command-line interface for NBA live scores, stats, and DARKO projections. No Python runtime is required.

## Build

```bash
go build -o nba .
./nba --help
```

## Examples

```bash
./nba live-game-summary
./nba live-boxscore 0022400123
./nba live-play-by-play 0022400123
./nba darko-player "Nikola Jokic"
./nba darko-leaderboard 10
```

See the [top-level README](../README.md) for the full command list.
