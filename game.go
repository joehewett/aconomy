package main

import (
	"fmt"

	"github.com/gorilla/websocket"
	openai "github.com/sashabaranov/go-openai"
)

// Resource types
const (
	Gold  = "Gold"
	Wheat = "Wheat"
)

// Building types
const (
	Farm = "Farm"
	Mine = "Mine"
)

// Game constants
const (
	NumAgents = 3

	StartingGold     = 50
	StartingWorkers  = 1
	StartingBuilding = 0
	StartingWheat    = 10

	WorkerCost = 10
	FarmCost   = 20
	MineCost   = 30

	WheatPerWorker = 1
	WheatDecayRate = 0.1

	FarmProduction = 3
	MineProduction = 5

	ActionsPerTurn = 1

	WinningGoldAmount = 1000
	MaxTurns          = 100
)

// [0]                 <- turn 0
//    [0]              <- agent 0
//       [AgentID]
//       [StartState]
//       [Strategy]
//       [Action]
//       [EndState]
//       [FullPrompt]
//    [1] ...
//    [2] ...
//    [3] ...

// Resource represents a quantity of a specific resource
type Resource struct {
	Type   string
	Amount int
}

// Building represents a production building
type Building struct {
	Type   string
	Manned bool
}

// Game represents the overall game state
type Game struct {
	Agents       []Agent
	GameLog      GameLog
	CurrentTurn  int
	Winner       *Agent
	Websocket    *websocket.Conn
	Done         chan struct{}
	OpenAIapiKey string
}

type GameLog []AgentTurn

type AgentTurn struct {
	Turn                int
	AgentID             int
	StartState          State
	Strategy            string
	Action              string
	EndState            State
	PostRationalisation string
	FullPrompt          []openai.ChatCompletionMessage
	Error               error
}

type State struct {
	Gold      int
	Wheat     int
	Workers   int
	Buildings []Building
}

// NewGame initializes a new game with the specified number of agents
func NewGame(conn *websocket.Conn, openAIapiKey string) *Game {
	game := &Game{
		Agents:       make([]Agent, NumAgents),
		GameLog:      GameLog{},
		CurrentTurn:  0,
		Winner:       nil,
		Websocket:    conn,
		Done:         make(chan struct{}),
		OpenAIapiKey: openAIapiKey,
	}

	for i := 0; i < NumAgents; i++ {
		game.Agents[i] = Agent{
			ID:        i,
			Gold:      StartingGold,
			Wheat:     StartingWheat,
			Workers:   StartingWorkers,
			Buildings: []Building{},
			Prompt:    basePrompt(),
			Lost:      false,
		}
	}

	return game
}

// RunGame manages the main game loop
func RunGame(game *Game) {
	for game.CurrentTurn < MaxTurns && game.Winner == nil {

		for i := range game.Agents {
			var agentTurn AgentTurn
			select {
			case <-game.Done:
				fmt.Printf("Game end detected, breaking out of game loop\n")
				break
			default:
			}

			if game.Agents[i].Lost {
				continue
			}

			agentTurn = ProcessTurn(&game.Agents[i], game)

			if game.Agents[i].Gold >= WinningGoldAmount {
				game.Winner = &game.Agents[i]
				break
			}

			if isLastAgent(game.Agents, game.Agents[i].ID) {
				game.Winner = &game.Agents[i]
			}

			err := game.PushGameState(agentTurn)
			if err != nil {
				fmt.Printf("Failed to push game state: %v\n", err)
				break
			}
		}

		game.CurrentTurn++

	}

	fmt.Printf("Game loop exiting after %d turns\n", game.CurrentTurn)
}

// ProcessTurn handles a single agent's turn
func ProcessTurn(agent *Agent, game *Game) AgentTurn {
	agent.IncrementTurn() // Increment the agent's turn counter
	agent.FeedWorkers()
	agent.ProduceResources()
	agent.DecayWheat()

	agent.AddTurnLog(getTurnPrompt())

	agentTurn, err := agent.TakeTurn(game, ActionsPerTurn)
	if err != nil {
		fmt.Printf("Agent %d failed to take turn: %v\n", agent.ID, err)
		game.End()
	}

	agent.EndTurn(game)

	fmt.Printf("Agent %d's turn ended\n", agent.ID)

	return *agentTurn
}

// PrintGameResult displays the final game state
func PrintGameResult(game *Game) {
	fmt.Printf("Game ended after %d turns\n", game.CurrentTurn)
	if game.Winner != nil {
		fmt.Printf("Winner: Agent %d\n", game.Winner.ID)
	} else {
		fmt.Println("No winner (max turns reached)")
	}

	for _, agent := range game.Agents {
		fmt.Printf("Agent %d: Gold=%d, Wheat=%d, Workers=%d, Buildings=%d\n",
			agent.ID, agent.Gold, agent.Wheat, agent.Workers, len(agent.Buildings))
	}
}

func (g *Game) End() {
	fmt.Printf("Ending game after %d turns\n", g.CurrentTurn)
	close(g.Done)
}

func (g *Game) PushGameState(agentTurn AgentTurn) error {
	g.GameLog = append(g.GameLog, agentTurn)
	if err := g.Websocket.WriteJSON(agentTurn); err != nil {
		return fmt.Errorf("failed to write game state to websocket: %w", err)
	}

	return nil
}

func (g *Game) broadcastMessage(message string, fromAgentID int) {
	for i := range g.Agents {
		if i != fromAgentID {
			g.Agents[i].AddTurnLog(message)
		}
	}
}

func (g *Game) sendMessage(targetAgentID int, message string) error {
	recipient := g.Agents[targetAgentID]
	if recipient.Lost {
		return fmt.Errorf("cannot send message to lost agent")
	}

	g.Agents[targetAgentID].AddTurnLog(message)

	return nil
}

func isLastAgent(agents []Agent, agentID int) bool {
	for _, agent := range agents {
		if agent.ID != agentID && !agent.Lost {
			return false
		}
	}

	return true
}
