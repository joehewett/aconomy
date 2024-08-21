package main

import (
	"encoding/json"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// Agent represents a player in the game
type Agent struct {
	ID        int
	Gold      int
	Wheat     int
	Workers   int
	Buildings []Building
	Prompt    []openai.ChatCompletionMessage
	Turn      int
	Lost      bool
}

func (a *Agent) IncrementTurn() {
	a.AddTurnLog(fmt.Sprintf("Incrementing turn to %d", a.Turn+1))
	a.Turn++

	gameState := a.getGameState()
	a.AddTurnLog(fmt.Sprintf("Current game state: %v", gameState))
	a.AddTurnLog("Performing mandatory start-of-turn actions...")
}

func (a *Agent) TakeTurn(g *Game, actionCount int) (t *AgentTurn, e error) {
	turn := AgentTurn{
		Turn:    a.Turn,
		AgentID: a.ID,
		StartState: State{
			Gold:      a.Gold,
			Wheat:     a.Wheat,
			Workers:   a.Workers,
			Buildings: a.Buildings,
		},
	}

	// Set the error on the returned turn if one occurs
	defer func() {
		if e != nil {
			turn.Error = e
		}
	}()

	a.AddTurnLog(fmt.Sprintf("Current state: Gold: %d, Wheat: %d, Workers: %d, Buildings: %+v", a.Gold, a.Wheat, a.Workers, a.Buildings))
	a.AddTurnLog("First, please outline your strategy for this turn. Afterwards, you will be prompted to take actions one by one.")

	strategy, err := getReasoningFromLM(a.Prompt)
	if err != nil {
		return &turn, fmt.Errorf("failed to call LM: %w", err)
	}

	a.AddAgentMessage(strategy)
	turn.Strategy = strategy

	actionsLeft := actionCount

	// Loop until the agent has no actions left
	for {
		if actionsLeft <= 0 {
			break
		}

		a.AddTurnLog(fmt.Sprintf("Please choose your next action. You have %d actions left for this turn", actionsLeft))

		toolCall, err := getToolCall(a.Prompt)
		if err != nil {
			return &turn, fmt.Errorf("failed to get tool call: %w", err)
		}

		turn.Action = toolCall.Function.Name

		a.AddTurnLog(fmt.Sprintf("You chose to take action: %v", toolCall.Function.Name))

		err = a.TakeAction(g, *toolCall)
		if err != nil {
			return &turn, fmt.Errorf("failed to take action: %w", err)
		}

		actionsLeft--
	}

	a.AddTurnLog("Your turn has ended. Please explain your reasoning for your actions, how you think your turn went, any unforseen issues that arose or future issues you see arising, and any other thoughts you have.")
	postRationalisation, err := getReasoningFromLM(a.Prompt)
	if err != nil {
		return &turn, fmt.Errorf("failed to call LM: %w", err)
	}

	a.AddAgentMessage(postRationalisation)
	turn.PostRationalisation = postRationalisation

	a.AddTurnLog("Your turn has now ended. Waiting for other agents to finish their turns...")

	turn.FullPrompt = a.Prompt

	turn.EndState = State{
		Gold:      a.Gold,
		Wheat:     a.Wheat,
		Workers:   a.Workers,
		Buildings: a.Buildings,
	}

	return &turn, nil
}

func (a *Agent) EndTurn(g *Game) {
	g.broadcastMessage(
		fmt.Sprintf(
			"Agent %d has ended their turn with %d gold, %d wheat, %d workers, and %d buildings",
			a.ID, a.Gold, a.Wheat, a.Workers, len(a.Buildings),
		),
		a.ID,
	)

	// Check if the agent is in a losing state, and if so, mark them as lost and tell the other agents
	if a.Gold == 0 && a.Wheat == 0 && a.Workers == 0 {
		a.Lost = true
		g.broadcastMessage(fmt.Sprintf("Agent %d has been eliminated from the game", a.ID), a.ID)
	}

}

func (a *Agent) getGameState() string {
	buildingsString := ""
	occupiedWorkers := 0
	for i, building := range a.Buildings {
		buildingsString = fmt.Sprintf("Building %d: %s / Manned: %t\n", i, building.Type, building.Manned)
		if building.Manned {
			occupiedWorkers++
		}
	}

	return fmt.Sprintf("Gold: %d\nWheat: %d\nTotal Workers: %d\nUnoccupied Workers: %d\nBuildings: %s", a.Gold, a.Wheat, a.Workers, a.Workers-occupiedWorkers, buildingsString)

}

func (a *Agent) TakeAction(g *Game, toolCall openai.ToolCall) error {
	argsString := toolCall.Function.Arguments
	// Convert the arguments to a map for easier access

	argMap := map[string]interface{}{}

	err := json.Unmarshal([]byte(argsString), &argMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal arguments: %w", err)
	}

	// fmt.Printf("ARGS: %v\n", argMap)

	switch toolCall.Function.Name {
	case "give_resources":
		targetAgent := argMap["target_agent"].(float64)
		resourceType := argMap["resource"].(map[string]interface{})["type"].(string)
		resourceAmount := argMap["resource"].(map[string]interface{})["amount"].(float64)
		resource := Resource{
			Type:   resourceType,
			Amount: int(resourceAmount),
		}
		if resource.Amount <= 0 {
			return fmt.Errorf("resource amount must be greater than 0")
		}

		a.GiveResource(g, int(targetAgent), resource)
	case "send_message":
		targetAgent := argMap["target_agent"].(float64)
		message := argMap["message"].(string)
		a.SendMessage(g, int(targetAgent), message)
	case "buy_building":
		buildingType := argMap["building_type"].(string)
		a.BuyBuilding(g, buildingType)
	case "buy_worker":
		count := argMap["count"].(float64)
		if count <= 0 {
			return fmt.Errorf("worker count must be greater than 0")
		}

		a.BuyWorkers(g, int(count))
	case "end_turn":
		// No action needed for this tool
		a.AddTurnLog("Ending turn early")
	case "unman_building":
		buildingType := argMap["building_type"].(string)
		a.UnmanBuilding(g, buildingType)
	case "man_building":
		buildingType := argMap["building_type"].(string)
		a.ManBuilding(g, buildingType)
	default:
		a.AddTurnLog(fmt.Sprintf("Unknown tool name: %s", toolCall.Function.Name))
	}

	return nil
}

func (a *Agent) ManBuilding(g *Game, buildingType string) {
	occupiedWorkers := a.getOccupiedWorkers()
	if occupiedWorkers >= a.Workers {
		a.AddTurnLog(fmt.Sprintf("Unable to man %s building, all workers are already occupied", buildingType))

		return
	}

	for i, building := range a.Buildings {
		if !(building.Type == buildingType) || building.Manned {
			continue
		}

		a.Buildings[i].Manned = true
		a.AddTurnLog(fmt.Sprintf("Manned a %s", buildingType))

		return
	}

	a.AddTurnLog(fmt.Sprintf("Unable to man %s building, no unoccupied buildings of that type found", buildingType))
}

func (a *Agent) UnmanBuilding(g *Game, buildingType string) {
	occupiedWorkers := a.getOccupiedWorkers()

	if occupiedWorkers <= 0 {
		a.AddTurnLog("Unable to unman building, no workers are currently occupied")
	}

	for i, building := range a.Buildings {
		if !(building.Type == buildingType) || !building.Manned {
			continue
		}

		a.Buildings[i].Manned = false
		a.AddTurnLog(fmt.Sprintf("Unmanned a %s, freed up 1 worker", buildingType))

		return
	}

	a.AddTurnLog(fmt.Sprintf("Unable to unman %s building, no occupied buildings of that type found", buildingType))
}

func (a *Agent) getOccupiedWorkers() int {
	occupiedWorkers := 0
	for _, building := range a.Buildings {
		if building.Manned {
			occupiedWorkers++
		}
	}

	return occupiedWorkers
}

func (a *Agent) SendMessage(g *Game, targetAgent int, message string) {
	// Find the target agent
	message = fmt.Sprintf("You have received a message from Agent %d! The message says: %s", a.ID, message)

	err := g.sendMessage(targetAgent, message)
	if err != nil {
		a.AddTurnLog(fmt.Sprintf("Failed to send message to Agent %d: %s", targetAgent, err))
		return
	}
}

// GiveResource transfers resources from one agent to another
func (a *Agent) GiveResource(g *Game, targetAgent int, resource Resource) {
	a.AddTurnLog(fmt.Sprintf("Attempting to give %d %s to Agent %d", resource.Amount, resource.Type, targetAgent))

	recipient := g.Agents[targetAgent]
	if recipient.Lost {
		a.AddTurnLog(fmt.Sprintf("Failed to give %d %s to Agent %d, recipient is eliminated", resource.Amount, resource.Type, targetAgent))
		return
	}

	if resource.Type == Gold {
		if a.Gold >= resource.Amount {
			a.Gold -= resource.Amount
			recipient.Gold += resource.Amount
			a.SendMessage(g, targetAgent, fmt.Sprintf("You have received %d gold from Agent %d", resource.Amount, a.ID))
			a.AddTurnLog(fmt.Sprintf("Transferred %d %s to Agent %d", resource.Amount, resource.Type, targetAgent))
			return
		}
	} else if resource.Type == Wheat {
		if a.Wheat >= resource.Amount {
			a.Wheat -= resource.Amount
			recipient.Wheat += resource.Amount
			a.SendMessage(g, targetAgent, fmt.Sprintf("You have received %d wheat from Agent %d", resource.Amount, a.ID))
			a.AddTurnLog(fmt.Sprintf("Transferred %d %s to Agent %d", resource.Amount, resource.Type, targetAgent))
			return
		}
	}

	a.AddTurnLog(fmt.Sprintf("Failed to give %d %s to Agent %d, not enough resources", resource.Amount, resource.Type, targetAgent))

}

// BuyWorkers adds workers to the agent if they can afford it
func (a *Agent) BuyWorkers(g *Game, count int) {
	a.AddTurnLog(fmt.Sprintf("Attempting to buy %d workers", count))

	cost := count * WorkerCost

	if cost > a.Gold {
		a.AddTurnLog(fmt.Sprintf("Failed to buy %d workers, not enough gold", count))
		return
	}

	a.Gold -= cost
	a.Workers += count
	a.AddTurnLog(fmt.Sprintf("Bought %d workers for %d gold", count, cost))
}

// BuyBuilding adds a building to the agent if they can afford it
func (a *Agent) BuyBuilding(g *Game, buildingType string) {
	a.AddTurnLog(fmt.Sprintf("Attempting to buy a %s", buildingType))
	var cost int
	switch buildingType {
	case Farm:
		cost = FarmCost
	case Mine:
		cost = MineCost
	default:
		return
	}

	if cost > a.Gold {
		a.AddTurnLog(fmt.Sprintf("Failed to buy a %s, not enough gold", buildingType))

		return
	}

	a.Gold -= cost
	a.Buildings = append(a.Buildings, Building{Type: buildingType, Manned: false})
	a.AddTurnLog(fmt.Sprintf("Bought a %s for %d gold", buildingType, cost))
	a.AddTurnLog(fmt.Sprintf("Buildings after purchase: %v", a.Buildings))
}

// FeedWorkers deducts wheat for each worker (double if the worker is working in a building) and kills unfed workers
func (a *Agent) FeedWorkers() {
	occupiedWorkers := a.getOccupiedWorkers()
	a.AddTurnLog(fmt.Sprintf("Attempting to feed %d workers with %d wheat", a.Workers, a.Wheat))
	workersFed, workersUnfed := 0, 0
	wheatNeeded := a.Workers*WheatPerWorker + (occupiedWorkers * WheatPerWorker)
	if a.Wheat >= wheatNeeded {
		a.Wheat -= wheatNeeded
	} else {
		workersFed = a.Wheat / WheatPerWorker
		workersUnfed = a.Workers - workersFed
		a.Workers = workersFed
		a.Wheat = 0

		// Ensure that the count of manned buildings is not more than the number of workers
		if occupiedWorkers > a.Workers {
			for i, building := range a.Buildings {
				if building.Manned {
					a.Buildings[i].Manned = false
					a.AddTurnLog(fmt.Sprintf("An occupied worker died due to starvation, a %s building is now unmanned", building.Type))
					occupiedWorkers--
					if occupiedWorkers == a.Workers {
						break
					}
				}
			}
		}
	}

	a.AddTurnLog(fmt.Sprintf("Fed %d workers, %d wheat remaining, %d workers died", a.Workers, a.Wheat, workersUnfed))
}

// ProduceResources generates resources from manned buildings
func (a *Agent) ProduceResources() {
	producedWheat, producedGold := 0, 0
	for i := range a.Buildings {
		if a.Buildings[i].Manned {
			switch a.Buildings[i].Type {
			case Farm:
				producedWheat += FarmProduction
				a.Wheat += FarmProduction
			case Mine:
				producedGold += MineProduction
				a.Gold += MineProduction
			}
		}
	}

	a.AddTurnLog(fmt.Sprintf("Produced %d wheat and %d gold from buildings", producedWheat, producedGold))
}

// DecayWheat reduces the agent's wheat by the decay rate
func (a *Agent) DecayWheat() {
	decayAmount := int(float64(a.Wheat) * WheatDecayRate)
	a.Wheat -= decayAmount

	a.AddTurnLog(fmt.Sprintf("Decayed %d wheat", decayAmount))
}

func (a *Agent) AddTurnLog(log string) {
	msg := fmt.Sprintf("Turn %d: %s\n", a.Turn, log)
	fmt.Printf("Agent %d: %s\n", a.ID, log)
	// a.TurnLog = append(a.TurnLog, log)
	a.Prompt = append(a.Prompt, openai.ChatCompletionMessage{Role: "system", Content: msg})
}

func (a *Agent) AddAgentMessage(msg string) {
	// fmt.Println(msg)
	a.Prompt = append(a.Prompt, openai.ChatCompletionMessage{Role: "assistant", Content: msg})
}
