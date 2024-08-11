package main

import (
	"log"
	"strings"
	"text/template"
)

var templateData = struct {
	WorkerCost        int
	FarmCost          int
	MineCost          int
	StartingGold      int
	StartingWheat     int
	StartingWorkers   int
	StartingBuilding  int
	ActionsPerTurn    int
	FarmProduction    int
	MineProduction    int
	WheatDecayRate    float64
	WinningGoldAmount int
	MaxTurn           int
	WheatPerWorker    int
}{
	WorkerCost:        WorkerCost,
	FarmCost:          FarmCost,
	MineCost:          MineCost,
	StartingGold:      StartingGold,
	StartingWheat:     StartingWheat,
	StartingWorkers:   StartingWorkers,
	StartingBuilding:  StartingBuilding,
	ActionsPerTurn:    ActionsPerTurn,
	FarmProduction:    FarmProduction,
	MineProduction:    MineProduction,
	WheatDecayRate:    WheatDecayRate,
	WinningGoldAmount: WinningGoldAmount,
	MaxTurn:           MaxTurns,
	WheatPerWorker:    WheatPerWorker,
}

func getSystemPrompt() string {

	var systemPromptTemplate = `
"You are an AI agent participating in a resource management and negotiation game called Aconomy. Your goal is to accumulate 1000 gold before any other agent. Here are the key details of the game:

1. Resources: There are two main resources - Gold and Wheat.
2. Buildings: You can build Farms (produce wheat) and Mines (produce gold).
3. Workers: You need workers to operate buildings. Each worker consumes 1 wheat per turn, or two if they are working in a building.
4. Starting conditions: You begin with {{ .StartingGold }} Gold, {{ .StartingWheat }} Wheat, {{ .StartingWorkers }} Workers, and {{ .StartingBuilding }} Buildings.
5. Actions: Each turn, you can perform {{ .ActionsPerTurn }} actions from the following:
   - Give resources (gold or wheat) to another agent
   - Buy workers ({{ .WorkerCost }} gold each)
   - Buy buildings (Farm: {{ .FarmCost }} gold, Mine: {{ .MineCost }} gold)
   - Send a message to another agent
   - End your turn early
   - Man a building with a worker so that it produces resources (workers manning a building consume 2*{{ .WheatPerWorker }} wheat per turn)
   - Unman a building so that it stops producing resources
6. Production:
   - A manned Farm produces {{ .FarmProduction }} wheat per turn
   - A manned Mine produces {{ .MineProduction }} gold per turn
7. Wheat decays at {{ .WheatDecayRate }}*totalWheat per turn
8. If you can't feed your workers, they will starve
9. The game ends when an agent reaches {{ .WinningGoldAmount }} gold or after {{ .MaxTurn }} turns

Your task is to make strategic decisions to grow your economy, manage your resources, and negotiate with other agents. Remember:

- Trades are based on trust; there's no mechanism to enforce agreements
- You can communicate freely with other agents to negotiate deals
- Balance short-term gains with long-term strategy
- Monitor your wheat production to ensure you can feed your workers. Unfed workers will die instantly.
- Consider the actions of other agents and adapt your strategy accordingly

In each turn, you will receive the current game state and must first strategise about your plan, then you will be given a chance to choose your actions.

Good luck, and may the best strategist win!
`

	var templ = template.Must(template.New("systemPrompt").Parse(systemPromptTemplate))

	tempWriter := new(strings.Builder)

	err := templ.Execute(tempWriter, templateData)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	return tempWriter.String()
}

func getTurnPrompt() string {
	var turnPromptTemplate = `
   It is now your turn to take actions. Remember, you can perform any {{ .ActionsPerTurn }} actions from the following:
- Give resources (gold or wheat) to another agent
- Buy workers ({{ .WorkerCost }} gold each)
- Buy buildings (Farm: {{ .FarmCost }} gold, Mine: {{ .MineCost }} gold)
- Send a message to another agent
- End your turn early
- Man a building with a worker so that it produces resources
- Unman a building so that it stops producing resources
You can perform any combination of these actions, up to 3 total actions per turn. Please choose one action at a time. To choose an action, please return one of the provided Tool Calls.
Please explain your reasoning for each action you take.
`

	var templ = template.Must(template.New("turnPrompt").Parse(turnPromptTemplate))

	tempWriter := new(strings.Builder)

	err := templ.Execute(tempWriter, templateData)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	return tempWriter.String()
}
