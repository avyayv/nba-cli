# nba-cli

Standalone Go command-line interface for NBA live scores, stats, and DARKO projections. No Python runtime is required.

## Build

```bash
go build -o nba .
./nba --help
```

## Install

### Prebuilt release binary

Every push to `main` publishes Linux, macOS, and Windows binaries to [GitHub Releases](https://github.com/avyayv/nba-cli/releases/latest).

macOS Apple Silicon:

```bash
mkdir -p ~/.local/bin
curl -fsSL https://github.com/avyayv/nba-cli/releases/latest/download/nba_darwin_arm64.tar.gz | tar -xz
mv nba ~/.local/bin/nba
chmod +x ~/.local/bin/nba
nba --help
```

Use `nba_darwin_amd64.tar.gz` for Intel Macs, `nba_linux_amd64.tar.gz` / `nba_linux_arm64.tar.gz` for Linux, or the Windows `.zip` assets.

### From source

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
