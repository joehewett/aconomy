package main

import (
	"fmt"
	"strings"

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

	ActionsPerTurn = 3

	WinningGoldAmount = 1000
	MaxTurns          = 100
)

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
	Agents      []Agent
	CurrentTurn int
	Winner      *Agent
}

// NewGame initializes a new game with the specified number of agents
func NewGame(numAgents int) *Game {
	game := &Game{
		Agents:      make([]Agent, numAgents),
		CurrentTurn: 0,
		Winner:      nil,
	}

	for i := 0; i < numAgents; i++ {
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

// Main function to run the game
func main() {
	game := NewGame(3) // Start a game with 3 agents
	RunGame(game)
}

// RunGame manages the main game loop
func RunGame(game *Game) {
	for game.CurrentTurn < MaxTurns && game.Winner == nil {
		for i := range game.Agents {
			if game.Agents[i].Lost {
				continue
			}

			ProcessTurn(&game.Agents[i], game)

			if game.Agents[i].Gold >= WinningGoldAmount {
				game.Winner = &game.Agents[i]
				break
			}

			if isLastAgent(game.Agents, game.Agents[i].ID) {
				game.Winner = &game.Agents[i]
			}
		}
		game.CurrentTurn++
	}

	PrintGameResult(game)
}

func isLastAgent(agents []Agent, agentID int) bool {
	for _, agent := range agents {
		if agent.ID != agentID && !agent.Lost {
			return false
		}
	}

	return true
}

// ProcessTurn handles a single agent's turn
func ProcessTurn(agent *Agent, game *Game) {
	agent.IncrementTurn() // Increment the agent's turn counter
	agent.FeedWorkers()
	agent.ProduceResources()
	agent.DecayWheat()

	err := PerformActions(agent, game)
	if err != nil {
		fmt.Println(err)
	}

	agent.EndTurn(game)

	fmt.Printf("Agent %d's turn ended\n", agent.ID)

	game.displayGameState()
}

// PerformActions allows the agent to take actions during their turn
func PerformActions(agent *Agent, g *Game) error {
	agent.AddTurnLog(getTurnPrompt())

	err := agent.TakeTurn(g, ActionsPerTurn) // This function should be implemented by the AI
	if err != nil {
		return fmt.Errorf("failed to decide action: %w", err)
	}

	return nil
}

func basePrompt() []openai.ChatCompletionMessage {
	sysPrompt := getSystemPrompt()

	fmt.Printf("System prompt: %s\n", sysPrompt)
	return []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: sysPrompt,
		},
	}
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

func (g *Game) displayGameState() {
	fmt.Println("\nCurrent game state:")
	for i, agent := range g.Agents {
		fmt.Printf("\nAgent %d:\n", i)
		fmt.Println("+-----------------------+")
		fmt.Printf("| Gold:     %-4d %s     \n", agent.Gold, getGoldSymbols(agent.Gold))
		fmt.Printf("| Wheat:    %-4d %s     \n", agent.Wheat, getWheatSymbols(agent.Wheat))
		fmt.Printf("| Unoccupied Workers:  %-4d %s     \n", agent.Workers-agent.getOccupiedWorkers(), getWorkerSymbols(agent.Workers))
		fmt.Printf("| Buildings:%-4d %s     \n", len(agent.Buildings), getBuildingSymbols(agent.Buildings))
		fmt.Println("+-----------------------+")
	}
}

func getGoldSymbols(count int) string {
	return strings.Repeat("💰", count/10)
}

func getWheatSymbols(count int) string {
	return strings.Repeat("🌾", count)
}

func getWorkerSymbols(count int) string {
	return strings.Repeat("👷", count)
}

func getBuildingSymbols(b []Building) string {
	s := ""
	for _, building := range b {
		if building.Type == Farm {
			s += "🚜"
			if building.Manned {
				// Tick symbol
				s += "✅"
			} else {
				// Cross symbol
				s += "❌"
			}
		} else if building.Type == Mine {
			s += "⛏"
			if building.Manned {
				// Tick symbol
				s += "✅"
			} else {
				// Cross symbol
				s += "❌"
			}
		}
	}

	return s
}
