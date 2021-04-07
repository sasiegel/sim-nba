package main

import (
	"time"
	"fmt"
	"math/rand"
)

const (
	nTeams = 30
	nGameDays = 81
	nGamesPerDay = 12 // no more than 15
)

type Player struct {
	offScore float64
	defScore float64
}

type Team struct {
	players []Player
	w int
	l int
}

func (t *Team) buildRoster() {
	const (
		nPlayers = 10
		mean = 0
		stdDev = 3
	)
	players := make([]Player, nPlayers)
	for i := 0; i < nPlayers; i++ {
		players[i].offScore = rand.NormFloat64() * stdDev + mean
		players[i].defScore = rand.NormFloat64() * stdDev + mean
	}
	t.players = players
}

func buildNBA() []Team {
	teams := make([]Team, nTeams)
	for i := 0; i < nTeams; i++ {
		teams[i].buildRoster()
		teams[i].w, teams[i].l = 0,0
	}
	return teams
}

func getTeamsIndices() []int {
	teamsIdx := make([]int, nTeams)
	for i := 0; i < nTeams; i++ {
		teamsIdx[i] = i
	}
	return teamsIdx
}

func buildSchedule() [][][]int {
	var schedule [][][]int
	for i := 0; i < nGameDays; i++ {
		var todaysGames [][]int
		availableTeams := getTeamsIndices()
		team1 := -1
		for j := 0; j < nGamesPerDay*2; j++ {
			// get a random index from availableTeams slice
			teamIdx := rand.Intn(len(availableTeams)) 
			if team1 < 0 {
				// assign team1
				team1 = availableTeams[teamIdx]
			} else {
				// assign team2 then append this new matchup to the schedule
				team2 := availableTeams[teamIdx]
				todaysGames = append(todaysGames, []int{team1, team2})
				team1 = -1
			}
			// remove that value from availableTeams slice
			availableTeams = append(availableTeams[:teamIdx], availableTeams[teamIdx+1:]...)
		}
		schedule = append(schedule, todaysGames)
	}
	return schedule
}

func versus(t1 *Team, t2 *Team) {
	const stdDev = 3
	var team1Score, team2Score float64
	for _, player := range t1.players {
		team1Score += player.offScore
		team2Score -= player.defScore
	}
	for _, player := range t2.players {
		team2Score += player.offScore
		team1Score += player.defScore
	}
	team1Draw := rand.NormFloat64() * stdDev + team1Score
	team2Draw := rand.NormFloat64() * stdDev + team2Score
	if team1Draw > team2Draw {
		t1.w++
		t2.l++
	} else {
		t1.l++
		t2.w++
	}
}

func sim(teams []Team, s [][][]int) {
	for _, day := range s {
		for _, game := range day {
			team1 := &teams[game[0]]
			team2 := &teams[game[1]]
			versus(team1, team2)
		}
	}
}

func printTeamRecords(teams []Team) {
	for i, team := range teams {
		fmt.Printf("Team %d: %d-%d\n", i, team.w, team.l)
	}
}

func printWinningTeam(teams []Team) {
	winningTeam := teams[0]
	winningIndex := 0
	for i, team := range teams {
		if team.w > winningTeam.w {
			winningTeam = team
			winningIndex = i
		}
	}
	fmt.Printf("Team %d finished with the best record!\n", winningIndex)
	for i, player := range winningTeam.players {
		fmt.Printf("Player %d: Off %.2f; Def %.2f\n", i, player.offScore, player.defScore)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	teams := buildNBA() 			// slice of teams
	schedule := buildSchedule()		// slice of days > slice of games > slice of matchup
	
	sim(teams, schedule)
	printTeamRecords(teams)
	printWinningTeam(teams)
}