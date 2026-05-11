package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const liveBase = "https://cdn.nba.com/static/json/liveData"
const statsBase = "https://stats.nba.com/stats"
const espnSiteBase = "https://site.api.espn.com/apis/site/v2/sports/basketball/nba"
const espnWebBase = "https://site.web.api.espn.com/apis"
const espnCoreBase = "https://sports.core.api.espn.com/v2/sports/basketball/leagues/nba"
const defaultUpdateSource = "https://github.com/avyayv/nba-cli.git"

type anyMap map[string]any

var teamIDs = map[string]int{
	"atlanta hawks": 1610612737, "boston celtics": 1610612738, "brooklyn nets": 1610612751,
	"charlotte hornets": 1610612766, "chicago bulls": 1610612741, "cleveland cavaliers": 1610612739,
	"dallas mavericks": 1610612742, "denver nuggets": 1610612743, "detroit pistons": 1610612765,
	"golden state warriors": 1610612744, "houston rockets": 1610612745, "indiana pacers": 1610612754,
	"la clippers": 1610612746, "los angeles clippers": 1610612746, "los angeles lakers": 1610612747,
	"memphis grizzlies": 1610612763, "miami heat": 1610612748, "milwaukee bucks": 1610612749,
	"minnesota timberwolves": 1610612750, "new orleans pelicans": 1610612740, "new york knicks": 1610612752,
	"oklahoma city thunder": 1610612760, "orlando magic": 1610612753, "philadelphia 76ers": 1610612755,
	"phoenix suns": 1610612756, "portland trail blazers": 1610612757, "sacramento kings": 1610612758,
	"san antonio spurs": 1610612759, "toronto raptors": 1610612761, "utah jazz": 1610612762,
	"washington wizards": 1610612764,
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() < 1 {
		usage()
		os.Exit(2)
	}
	cmd, args := flag.Arg(0), flag.Args()[1:]
	res, err := run(cmd, args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	printJSON(res)
}

func usage() {
	fmt.Fprintf(os.Stderr, `NBA CLI (Go, no Python required)

Usage: nba <command> [args]

Utilities:
  update [install_path]        Download, rebuild, and install the latest CLI

Live commands:
  live-scores                 Today's scoreboard
  live-game-summary           Compact live score/status/leaders
  live-boxscore <game_id>     Live box score
  live-play-by-play <game_id> Live play-by-play

DARKO:
  darko-leaderboard [limit]   Current DPM leaderboard from darko.app
  darko-player <name|id>      DARKO metrics for a player

ESPN public data:
  espn-search-player <name>       Search ESPN NBA players
  get-espn-player-id <name>       Resolve an ESPN player ID
  get-espn-team-id <team name>    Resolve an ESPN team ID
  espn-player <id|name>           ESPN player card/profile data
  espn-player-bio <id|name>       ESPN awards/biographical extras
  espn-player-stats <id|name>     ESPN player statistics/splits
  espn-player-contracts <id|name> [season_year]
  espn-teams                      ESPN teams index
  espn-team-roster <id|name>      ESPN roster with player biographical data
  espn-injuries [id|name]         ESPN league or team injury reports
  espn-news [limit]               ESPN NBA news
  espn-scoreboard [YYYYMMDD]      ESPN scoreboard
  espn-game-summary <espn_game_id> ESPN box score, leaders, plays, odds, etc.
  espn-standings                  ESPN standings hierarchy

Stats commands:
  get-team-id <team name>
  get-player-id <player name>
  player-profile <player_id>
  player-awards <player_id>
  player-career-stats <player_id>
  player-game-log <player_id>
  league-leaders [PTS|REB|AST|...]
  player-estimated-metrics
  draft-history [year]
  franchise-history
  team-details <team_id>
  league-dash-player-stats <team_id> [Advanced|Basic]
  team-standings
  team-game-log <team_id>
  team-roster <team_id>
  league-dash-team-stats [Advanced|Basic]
  league-schedule <team_id>
  game-win-probability <game_id>
  game-boxscore <game_id>
  league-dash-lineups <team_id> [Advanced|Basic] [2-5]
`)
}

func run(cmd string, args []string) (any, error) {
	switch cmd {
	case "update":
		return updateCLI(args)
	case "live-scores":
		return getJSON(liveBase+"/scoreboard/todaysScoreboard_00.json", nil)
	case "live-game-summary":
		return liveSummary()
	case "live-boxscore":
		return need1(args, func(id string) (any, error) { return getJSON(liveBase+"/boxscore/boxscore_"+id+".json", nil) })
	case "live-play-by-play":
		return need1(args, func(id string) (any, error) { return getJSON(liveBase+"/playbyplay/playbyplay_"+id+".json", nil) })
	case "darko-leaderboard":
		limit := 25
		if len(args) > 0 {
			limit, _ = strconv.Atoi(args[0])
		}
		return darkoLeaderboard(limit)
	case "darko-player":
		return need1(args, darkoPlayer)
	case "espn-search-player":
		if len(args) < 1 {
			return nil, errors.New("player name required")
		}
		return espnSearchPlayers(strings.Join(args, " "), 10)
	case "get-espn-player-id":
		if len(args) < 1 {
			return nil, errors.New("player name required")
		}
		return resolveESPNPlayerID(strings.Join(args, " "))
	case "get-espn-team-id":
		if len(args) < 1 {
			return nil, errors.New("team name required")
		}
		return getESPNTeamID(strings.Join(args, " "))
	case "espn-player":
		return espnPlayer(args)
	case "espn-player-bio":
		return espnPlayerBio(args)
	case "espn-player-stats":
		return espnPlayerStats(args)
	case "espn-player-contracts":
		return espnPlayerContracts(args)
	case "espn-teams":
		return espnTeams()
	case "espn-team-roster":
		if len(args) < 1 {
			return nil, errors.New("team id or name required")
		}
		return espnTeamRoster(strings.Join(args, " "))
	case "espn-injuries":
		return espnInjuries(args)
	case "espn-news":
		return espnNews(args)
	case "espn-scoreboard":
		return espnScoreboard(args)
	case "espn-game-summary":
		return need1(args, espnGameSummary)
	case "espn-standings":
		return espnStandings()
	case "get-team-id":
		if len(args) < 1 {
			return nil, errors.New("team name required")
		}
		return getTeamID(strings.Join(args, " "))
	case "get-player-id":
		if len(args) < 1 {
			return nil, errors.New("player name required")
		}
		return getPlayerID(strings.Join(args, " "))
	case "player-profile":
		return statsNeedID(args, "commonplayerinfo", map[string]string{"PlayerID": "%s"})
	case "player-awards":
		return statsNeedID(args, "playerawards", map[string]string{"PlayerID": "%s", "LeagueID": "__omit__", "SeasonType": "__omit__"})
	case "player-career-stats":
		return statsNeedID(args, "playercareerstats", map[string]string{"PlayerID": "%s", "PerMode": "PerGame"})
	case "player-game-log":
		return statsNeedID(args, "playergamelog", map[string]string{"PlayerID": "%s", "Season": currentSeason(), "SeasonType": "Regular Season"})
	case "league-leaders":
		return leagueLeaders(args)
	case "player-estimated-metrics":
		return stats("playerestimatedmetrics", map[string]string{"Season": currentSeason(), "SeasonType": "Regular Season"})
	case "draft-history":
		return draftHistory(args)
	case "franchise-history":
		return stats("franchisehistory", map[string]string{"SeasonType": "__omit__"})
	case "team-details":
		return statsNeedID(args, "teamdetails", map[string]string{"TeamID": "%s", "SeasonType": "__omit__"})
	case "league-dash-player-stats":
		if len(args) < 1 {
			return nil, errors.New("team_id required")
		}
		return stats("leaguedashplayerstats", dashPlayerParams(args[0], opt(args, 1, "Advanced")))
	case "team-standings":
		return stats("leaguestandingsv3", map[string]string{"LeagueID": "00", "Season": currentSeason(), "SeasonType": "Regular Season"})
	case "team-game-log":
		return statsNeedID(args, "teamgamelog", map[string]string{"TeamID": "%s", "Season": currentSeason(), "SeasonType": "Regular Season"})
	case "team-roster":
		return statsNeedID(args, "commonteamroster", map[string]string{"TeamID": "%s", "Season": currentSeason()})
	case "league-dash-team-stats":
		return stats("leaguedashteamstats", dashTeamParams(opt(args, 0, "Advanced")))
	case "league-schedule":
		return leagueSchedule(args)
	case "game-win-probability":
		return gameWinProbability(args)
	case "game-boxscore":
		return statsNeedID(args, "boxscoreadvancedv3", map[string]string{"GameID": "%s"})
	case "league-dash-lineups":
		if len(args) < 1 {
			return nil, errors.New("team_id required")
		}
		return stats("leaguedashlineups", dashLineupParams(args[0], opt(args, 1, "Advanced"), opt(args, 2, "2")))
	default:
		return nil, fmt.Errorf("unknown command %q", cmd)
	}
}

func need1(args []string, f func(string) (any, error)) (any, error) {
	if len(args) < 1 {
		return nil, errors.New("argument required")
	}
	return f(args[0])
}
func opt(args []string, i int, def string) string {
	if len(args) > i && args[i] != "" {
		return args[i]
	}
	return def
}
func statsNeedID(args []string, endpoint string, params map[string]string) (any, error) {
	if len(args) < 1 {
		return nil, errors.New("id required")
	}
	for k, v := range params {
		if v == "%s" {
			params[k] = args[0]
		}
	}
	return stats(endpoint, params)
}

func updateCLI(args []string) (any, error) {
	if len(args) > 1 {
		return nil, errors.New("usage: nba update [install_path]")
	}
	target, err := updateTarget(args)
	if err != nil {
		return nil, err
	}
	previousChecksum, _ := fileSHA256(target)

	source := os.Getenv("NBA_CLI_UPDATE_SOURCE")
	if source == "" {
		source = os.Getenv("NBA_CLI_UPDATE_URL")
	}
	if source == "" {
		source = defaultUpdateSource
	}

	tmp, err := os.MkdirTemp("", "nba-cli-update-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmp)

	sourceDir := filepath.Join(tmp, "source")
	if err := fetchUpdateSource(source, sourceDir); err != nil {
		return nil, fmt.Errorf("fetch latest source: %w", err)
	}

	cliDir, err := findCLIDir(sourceDir)
	if err != nil {
		return nil, err
	}

	builtBinary := filepath.Join(tmp, "nba")
	build := exec.Command("go", "build", "-o", builtBinary, ".")
	build.Dir = cliDir
	if out, err := build.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("go build failed: %w: %s", err, strings.TrimSpace(string(out)))
	}
	if err := os.Chmod(builtBinary, 0755); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return nil, err
	}
	if err := os.Rename(builtBinary, target); err != nil {
		_ = os.Remove(target)
		if retryErr := os.Rename(builtBinary, target); retryErr != nil {
			return nil, retryErr
		}
	}
	newChecksum, _ := fileSHA256(target)

	currentExecutable, _ := os.Executable()
	return anyMap{
		"updated":                true,
		"changed":                previousChecksum == "" || previousChecksum != newChecksum,
		"installedPath":          target,
		"source":                 source,
		"currentExecutable":      currentExecutable,
		"previousChecksumSha256": previousChecksum,
		"newChecksumSha256":      newChecksum,
	}, nil
}

func updateTarget(args []string) (string, error) {
	if len(args) == 1 {
		return filepath.Abs(args[0])
	}
	if env := os.Getenv("NBA_CLI_INSTALL_PATH"); env != "" {
		return filepath.Abs(env)
	}
	if exe, err := os.Executable(); err == nil && filepath.Base(exe) == "nba" {
		return filepath.Abs(exe)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "bin", "nba"), nil
}

func fetchUpdateSource(source, dst string) error {
	if strings.HasSuffix(strings.ToLower(strings.Split(source, "?")[0]), ".zip") {
		archivePath := filepath.Join(filepath.Dir(dst), "source.zip")
		if err := downloadFile(source, archivePath); err != nil {
			return err
		}
		return unzip(archivePath, dst)
	}

	clone := exec.Command("git", "clone", "--depth", "1", source, dst)
	if out, err := clone.CombinedOutput(); err != nil {
		return fmt.Errorf("git clone failed: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func downloadFile(u, dst string) error {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "nba-cli updater")
	client := &http.Client{Timeout: 2 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 300))
		return fmt.Errorf("%s: %s", resp.Status, string(body))
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func unzip(src, dst string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	cleanDst, err := filepath.Abs(dst)
	if err != nil {
		return err
	}
	for _, f := range r.File {
		target := filepath.Join(cleanDst, f.Name)
		if !strings.HasPrefix(target, cleanDst+string(os.PathSeparator)) {
			return fmt.Errorf("unsafe zip path %q", f.Name)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, f.Mode()); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		in, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			in.Close()
			return err
		}
		_, copyErr := io.Copy(out, in)
		closeInErr := in.Close()
		closeOutErr := out.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeInErr != nil {
			return closeInErr
		}
		if closeOutErr != nil {
			return closeOutErr
		}
	}
	return nil
}

func findCLIDir(root string) (string, error) {
	var cliDir string
	if err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || !d.IsDir() || cliDir != "" {
			return err
		}
		if filepath.Base(path) != "cli" {
			return nil
		}
		if _, err := os.Stat(filepath.Join(path, "go.mod")); err != nil {
			return nil
		}
		if _, err := os.Stat(filepath.Join(path, "main.go")); err != nil {
			return nil
		}
		cliDir = path
		return filepath.SkipDir
	}); err != nil {
		return "", err
	}
	if cliDir == "" {
		return "", errors.New("cli source not found in update source")
	}
	return cliDir, nil
}

func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func getJSON(u string, headers map[string]string) (any, error) {
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("Host", req.URL.Host)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referer", "https://www.nba.com/")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Sec-Ch-Ua", `"Not:A-Brand";v="99", "Google Chrome";v="145", "Chromium";v="145"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s: %s", resp.Status, string(body[:min(len(body), 300)]))
	}
	var out any
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()
	if err := dec.Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func stats(endpoint string, params map[string]string) (any, error) {
	q := url.Values{}
	defaults := map[string]string{"LeagueID": "00", "SeasonType": "Regular Season"}
	for k, v := range defaults {
		q.Set(k, v)
	}
	for k, v := range params {
		if v == "__omit__" {
			q.Del(k)
			continue
		}
		q.Set(k, v)
	}
	return getJSON(statsBase+"/"+endpoint+"?"+q.Encode(), map[string]string{"Origin": "https://www.nba.com"})
}

func currentSeason() string {
	now := time.Now()
	y := now.Year()
	if now.Month() < 8 {
		y--
	}
	return fmt.Sprintf("%d-%02d", y, (y+1)%100)
}

func dashCommon(measureType string) map[string]string {
	return map[string]string{
		"LastNGames": "0", "MeasureType": measureType, "Month": "0", "OpponentTeamID": "0",
		"PaceAdjust": "N", "PerMode": "Per100Possessions", "Period": "0", "PlusMinus": "N",
		"Rank": "N", "Season": currentSeason(), "SeasonType": "Regular Season", "Conference": "",
		"DateFrom": "", "DateTo": "", "Division": "", "GameSegment": "", "LeagueID": "00",
		"Location": "", "Outcome": "", "PORound": "", "SeasonSegment": "", "ShotClockRange": "",
		"TeamID": "", "VsConference": "", "VsDivision": "",
	}
}

func dashPlayerParams(teamID, measureType string) map[string]string {
	p := dashCommon(measureType)
	p["TeamID"] = teamID
	p["College"] = ""
	p["Country"] = ""
	p["DraftPick"] = ""
	p["DraftYear"] = ""
	p["GameScope"] = ""
	p["Height"] = ""
	p["PlayerExperience"] = ""
	p["PlayerPosition"] = ""
	p["StarterBench"] = ""
	p["TwoWay"] = ""
	p["Weight"] = ""
	return p
}

func dashTeamParams(measureType string) map[string]string {
	p := dashCommon(measureType)
	p["GameScope"] = ""
	p["PlayerExperience"] = ""
	p["PlayerPosition"] = ""
	p["StarterBench"] = ""
	p["TwoWay"] = ""
	return p
}

func dashLineupParams(teamID, measureType, groupQuantity string) map[string]string {
	p := dashCommon(measureType)
	p["TeamID"] = teamID
	p["GroupQuantity"] = groupQuantity
	return p
}

func getTeamID(name string) (any, error) {
	id, ok := teamIDs[strings.ToLower(name)]
	if !ok {
		return nil, errors.New("unknown team")
	}
	return id, nil
}
func getPlayerID(name string) (any, error) {
	res, err := stats("commonallplayers", map[string]string{"Season": currentSeason(), "IsOnlyCurrentSeason": "0"})
	if err != nil {
		return nil, err
	}
	rows := resultRows(res)
	needle := strings.ToLower(name)
	for _, r := range rows {
		if len(r) > 2 && strings.Contains(strings.ToLower(fmt.Sprint(r[2])), needle) {
			return r[0], nil
		}
	}
	return nil, errors.New("player not found")
}

func leagueLeaders(args []string) (any, error) {
	return stats("leagueleaders", map[string]string{
		"StatCategory": strings.ToUpper(opt(args, 0, "PTS")),
		"PerMode":      opt(args, 1, "PerGame"),
		"Scope":        "S",
		"Season":       currentSeason(),
		"SeasonType":   "Regular Season",
	})
}

func draftHistory(args []string) (any, error) {
	season := currentSeason()[:4]
	if len(args) > 0 && args[0] != "" {
		season = args[0]
	}
	return stats("drafthistory", map[string]string{
		"Season":      season,
		"College":     "",
		"OverallPick": "",
		"RoundNum":    "",
		"RoundPick":   "",
		"TeamID":      "0",
		"TopX":        "",
		"SeasonType":  "__omit__",
	})
}

func espnSearchPlayers(name string, limit int) (any, error) {
	q := url.Values{}
	q.Set("query", name)
	q.Set("limit", strconv.Itoa(limit))
	q.Set("mode", "prefix")
	q.Set("type", "player")
	return getJSON(espnWebBase+"/common/v3/search?"+q.Encode(), nil)
}

func resolveESPNPlayerID(q string) (string, error) {
	q = strings.TrimSpace(q)
	if q == "" {
		return "", errors.New("player id or name required")
	}
	if isDigits(q) {
		return q, nil
	}
	res, err := espnSearchPlayers(q, 10)
	if err != nil {
		return "", err
	}
	m, _ := res.(map[string]any)
	items, _ := m["items"].([]any)
	for _, item := range items {
		im, _ := item.(map[string]any)
		if fmt.Sprint(im["type"]) != "player" {
			continue
		}
		if leagueAny, ok := im["league"]; ok {
			if league := strings.ToLower(fmt.Sprint(leagueAny)); league != "" && league != "<nil>" && league != "nba" {
				continue
			}
		}
		id := fmt.Sprint(im["id"])
		if id != "" && id != "<nil>" {
			return id, nil
		}
	}
	return "", errors.New("ESPN player not found")
}

func espnPlayer(args []string) (any, error) {
	if len(args) < 1 {
		return nil, errors.New("player id or name required")
	}
	id, err := resolveESPNPlayerID(strings.Join(args, " "))
	if err != nil {
		return nil, err
	}
	return getJSON(espnWebBase+"/common/v3/sports/basketball/nba/athletes/"+url.PathEscape(id), nil)
}

func espnPlayerBio(args []string) (any, error) {
	if len(args) < 1 {
		return nil, errors.New("player id or name required")
	}
	id, err := resolveESPNPlayerID(strings.Join(args, " "))
	if err != nil {
		return nil, err
	}
	return getJSON(espnWebBase+"/common/v3/sports/basketball/nba/athletes/"+url.PathEscape(id)+"/bio", nil)
}

func espnPlayerStats(args []string) (any, error) {
	if len(args) < 1 {
		return nil, errors.New("player id or name required")
	}
	id, err := resolveESPNPlayerID(strings.Join(args, " "))
	if err != nil {
		return nil, err
	}
	return getJSON(espnCoreBase+"/athletes/"+url.PathEscape(id)+"/statistics?lang=en&region=us", nil)
}

func espnPlayerContracts(args []string) (any, error) {
	id, season, err := parseESPNContractArgs(args)
	if err != nil {
		return nil, err
	}
	if season != "" {
		contract, err := getJSON(espnCoreBase+"/athletes/"+url.PathEscape(id)+"/contracts/"+url.PathEscape(season)+"?lang=en&region=us", nil)
		if err != nil {
			return nil, err
		}
		return anyMap{"athleteId": id, "season": season, "contract": contract}, nil
	}

	index, err := getJSON(espnCoreBase+"/athletes/"+url.PathEscape(id)+"/contracts?lang=en&region=us", nil)
	if err != nil {
		return nil, err
	}
	m, _ := index.(map[string]any)
	items, _ := m["items"].([]any)
	contracts := make([]any, 0, len(items))
	for _, item := range items {
		contract, err := derefESPN(item)
		if err != nil {
			contracts = append(contracts, anyMap{"ref": refString(item), "error": err.Error()})
			continue
		}
		contracts = append(contracts, contract)
	}
	return anyMap{"athleteId": id, "count": len(contracts), "contracts": contracts}, nil
}

func parseESPNContractArgs(args []string) (string, string, error) {
	if len(args) < 1 {
		return "", "", errors.New("player id or name required")
	}
	season := ""
	playerArgs := args
	if len(args) > 1 && regexp.MustCompile(`^\d{4}$`).MatchString(args[len(args)-1]) {
		season = args[len(args)-1]
		playerArgs = args[:len(args)-1]
	}
	id, err := resolveESPNPlayerID(strings.Join(playerArgs, " "))
	return id, season, err
}

func espnTeams() (any, error) {
	return getJSON(espnSiteBase+"/teams", nil)
}

func resolveESPNTeamID(q string) (string, error) {
	q = strings.TrimSpace(q)
	if q == "" {
		return "", errors.New("team id or name required")
	}
	if isDigits(q) {
		return q, nil
	}
	v, err := getESPNTeamID(q)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(v), nil
}

func getESPNTeamID(name string) (any, error) {
	team, err := findESPNTeam(name)
	if err != nil {
		return nil, err
	}
	return team["id"], nil
}

func findESPNTeam(name string) (map[string]any, error) {
	res, err := espnTeams()
	if err != nil {
		return nil, err
	}
	needle := normalizeName(name)
	root, _ := res.(map[string]any)
	sports, _ := root["sports"].([]any)
	for _, sport := range sports {
		sm, _ := sport.(map[string]any)
		leagues, _ := sm["leagues"].([]any)
		for _, league := range leagues {
			lm, _ := league.(map[string]any)
			teams, _ := lm["teams"].([]any)
			for _, t := range teams {
				tm, _ := t.(map[string]any)
				team, _ := tm["team"].(map[string]any)
				candidates := []string{"displayName", "shortDisplayName", "abbreviation", "name", "location", "slug"}
				for _, c := range candidates {
					val := normalizeName(fmt.Sprint(team[c]))
					if val == needle || (len(needle) >= 3 && strings.Contains(val, needle)) || (len(val) >= 4 && strings.Contains(needle, val)) {
						return team, nil
					}
				}
			}
		}
	}
	return nil, errors.New("ESPN team not found")
}

func espnTeamRoster(q string) (any, error) {
	id, err := resolveESPNTeamID(q)
	if err != nil {
		return nil, err
	}
	return getJSON(espnSiteBase+"/teams/"+url.PathEscape(id)+"/roster", nil)
}

func espnInjuries(args []string) (any, error) {
	res, err := getJSON(espnSiteBase+"/injuries", nil)
	if err != nil {
		return nil, err
	}
	if len(args) == 0 {
		return res, nil
	}
	id, err := resolveESPNTeamID(strings.Join(args, " "))
	if err != nil {
		return nil, err
	}
	root, _ := res.(map[string]any)
	teams, _ := root["injuries"].([]any)
	for _, t := range teams {
		tm, _ := t.(map[string]any)
		if sameID(tm["id"], id) {
			return tm, nil
		}
	}
	return anyMap{"teamId": id, "injuries": []any{}}, nil
}

func espnNews(args []string) (any, error) {
	limit := 10
	if len(args) > 0 {
		if n, err := strconv.Atoi(args[0]); err == nil && n > 0 {
			limit = n
		}
	}
	q := url.Values{}
	q.Set("limit", strconv.Itoa(limit))
	return getJSON(espnSiteBase+"/news?"+q.Encode(), nil)
}

func espnScoreboard(args []string) (any, error) {
	q := url.Values{}
	if len(args) > 0 && args[0] != "" {
		q.Set("dates", args[0])
	}
	u := espnSiteBase + "/scoreboard"
	if len(q) > 0 {
		u += "?" + q.Encode()
	}
	return getJSON(u, nil)
}

func espnGameSummary(id string) (any, error) {
	q := url.Values{}
	q.Set("event", id)
	return getJSON(espnSiteBase+"/summary?"+q.Encode(), nil)
}

func espnStandings() (any, error) {
	return getJSON(espnWebBase+"/v2/sports/basketball/nba/standings?region=us&lang=en&contentorigin=espn&type=0&level=3&sort=winPercent%3Adesc", nil)
}

func derefESPN(item any) (any, error) {
	ref := refString(item)
	if ref == "" {
		return item, nil
	}
	ref = strings.Replace(ref, "http://sports.core.api.espn.com/", "https://sports.core.api.espn.com/", 1)
	return getJSON(ref, nil)
}

func refString(item any) string {
	m, _ := item.(map[string]any)
	if ref, ok := m["$ref"]; ok && ref != nil {
		return fmt.Sprint(ref)
	}
	return ""
}

func normalizeName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	re := regexp.MustCompile(`[^a-z0-9]+`)
	return strings.Trim(re.ReplaceAllString(s, " "), " ")
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func leagueSchedule(args []string) (any, error) {
	if len(args) < 1 {
		return nil, errors.New("team_id required")
	}
	res, err := stats("scheduleleaguev2", map[string]string{"Season": currentSeason(), "SeasonType": "__omit__"})
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(res)
	var root anyMap
	json.Unmarshal(b, &root)
	tid := args[0]
	var out []anyMap
	ls, _ := root["leagueSchedule"].(map[string]any)
	dates, _ := ls["gameDates"].([]any)
	for _, d := range dates {
		dm, _ := d.(map[string]any)
		games, _ := dm["games"].([]any)
		for _, g := range games {
			gm, _ := g.(map[string]any)
			h, _ := gm["homeTeam"].(map[string]any)
			a, _ := gm["awayTeam"].(map[string]any)
			if sameID(h["teamId"], tid) || sameID(a["teamId"], tid) {
				out = append(out, gm)
			}
		}
	}
	return out, nil
}

func gameWinProbability(args []string) (any, error) {
	if len(args) < 1 {
		return nil, errors.New("game_id required")
	}
	res, err := stats("winprobabilitypbp", map[string]string{"GameID": args[0], "RunType": "each second", "LeagueID": "__omit__", "SeasonType": "__omit__"})
	if err != nil {
		return anyMap{"gameId": args[0], "error": "Win probability feed unavailable for this game", "detail": err.Error()}, nil
	}
	return res, nil
}

func liveSummary() (any, error) {
	res, err := getJSON(liveBase+"/scoreboard/todaysScoreboard_00.json", nil)
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(res)
	var root anyMap
	json.Unmarshal(b, &root)
	sb, _ := root["scoreboard"].(map[string]any)
	games, _ := sb["games"].([]any)
	var out []anyMap
	for _, g := range games {
		gm := g.(map[string]any)
		h := gm["homeTeam"].(map[string]any)
		a := gm["awayTeam"].(map[string]any)
		out = append(out, anyMap{"gameId": gm["gameId"], "status": gm["gameStatusText"], "period": gm["period"], "clock": gm["gameClock"], "away": fmt.Sprintf("%v %v", a["teamTricode"], a["score"]), "home": fmt.Sprintf("%v %v", h["teamTricode"], h["score"]), "leaders": gm["gameLeaders"]})
	}
	return out, nil
}

func fetchDarko() ([]anyMap, error) {
	resp, err := http.Get("https://www.darko.app/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	s := string(body)
	i := strings.Index(s, "players:[")
	if i < 0 {
		return nil, errors.New("DARKO player data not found")
	}
	s = s[i+9:]
	reObj := regexp.MustCompile(`\{[^{}]*player_name:"[^"]+"[^{}]*\}`)
	objs := reObj.FindAllString(s, -1)
	fields := []string{"nba_id", "player_name", "team_name", "position", "date", "season", "dpm", "o_dpm", "d_dpm", "box_dpm", "on_off_dpm", "x_minutes", "x_pace", "x_pts_100", "x_ast_100", "x_fg_pct", "x_fg3_pct", "x_ft_pct"}
	var out []anyMap
	for _, obj := range objs {
		row := anyMap{}
		for _, f := range fields {
			row[f] = darkoField(obj, f)
		}
		out = append(out, row)
	}
	sort.Slice(out, func(i, j int) bool { return num(out[i]["dpm"]) > num(out[j]["dpm"]) })
	return out, nil
}
func darkoLeaderboard(limit int) (any, error) {
	ps, err := fetchDarko()
	if err != nil {
		return nil, err
	}
	if limit > len(ps) {
		limit = len(ps)
	}
	return ps[:limit], nil
}
func darkoPlayer(q string) (any, error) {
	ps, err := fetchDarko()
	if err != nil {
		return nil, err
	}
	n := strings.ToLower(q)
	for _, p := range ps {
		if fmt.Sprint(p["nba_id"]) == q || strings.Contains(strings.ToLower(fmt.Sprint(p["player_name"])), n) {
			return p, nil
		}
	}
	return nil, errors.New("DARKO player not found")
}
func darkoField(obj, name string) any {
	re := regexp.MustCompile(name + `:("[^"]*"|null|-?\.\d+|-?\d+(?:\.\d+)?)`)
	m := re.FindStringSubmatch(obj)
	if len(m) < 2 {
		return nil
	}
	v := m[1]
	if v == "null" {
		return nil
	}
	if strings.HasPrefix(v, "\"") {
		return strings.Trim(v, "\"")
	}
	if strings.HasPrefix(v, ".") {
		v = "0" + v
	}
	if strings.HasPrefix(v, "-.") {
		v = "-0." + strings.TrimPrefix(v, "-.")
	}
	f, _ := strconv.ParseFloat(v, 64)
	return f
}
func num(v any) float64 { f, _ := strconv.ParseFloat(fmt.Sprint(v), 64); return f }

func resultRows(res any) [][]any {
	b, _ := json.Marshal(res)
	var root anyMap
	json.Unmarshal(b, &root)
	sets, _ := root["resultSets"].([]any)
	if len(sets) == 0 {
		sets, _ = root["resultSet"].([]any)
	}
	if len(sets) == 0 {
		return nil
	}
	first, _ := sets[0].(map[string]any)
	rows, _ := first["rowSet"].([]any)
	out := make([][]any, 0, len(rows))
	for _, r := range rows {
		if a, ok := r.([]any); ok {
			out = append(out, a)
		}
	}
	return out
}

func printJSON(v any) { b, _ := json.MarshalIndent(v, "", "  "); fmt.Println(string(b)) }
func sameID(v any, id string) bool {
	want, _ := strconv.ParseInt(id, 10, 64)
	switch x := v.(type) {
	case json.Number:
		got, _ := x.Int64()
		return got == want
	case float64:
		return int64(x) == want
	default:
		return fmt.Sprint(v) == id
	}
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
