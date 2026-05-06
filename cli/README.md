# nba-cli

Standalone Go command-line interface for NBA live scores, stats, and DARKO projections. No Python runtime is required.

## Build

```bash
go build -o nba .
./nba --help
```

## Install

```bash
mkdir -p ~/.local/bin
go build -o ~/.local/bin/nba .
chmod +x ~/.local/bin/nba
nba --help
```

## Update

```bash
nba update
```

`nba update` fetches the latest source from GitHub, rebuilds the Go CLI, and replaces the installed `nba` binary. Pass an explicit target path if needed: `nba update ~/.local/bin/nba`.

## Examples

```bash
./nba live-game-summary
./nba live-boxscore 0022400123
./nba live-play-by-play 0022400123
./nba darko-player "Nikola Jokic"
./nba darko-leaderboard 10
./nba update ./nba
```

See the [top-level README](../README.md) for the full command list.
