package main

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
	jsonschema "github.com/sashabaranov/go-openai/jsonschema"
)

var tools = getToolDefinitions()

func validateAPIKey(apiKey string) error {
	client := openai.NewClient(apiKey)
	_, err := client.ListModels(context.Background())
	return err
}

func getReasoningFromLM(prompt []openai.ChatCompletionMessage, key string) (string, error) {
	client := openai.NewClient(key)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:      openai.GPT3Dot5Turbo,
			Messages:   prompt,
			Tools:      tools,
			ToolChoice: "none",
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}

func getToolCall(messages []openai.ChatCompletionMessage, key string) (*openai.ToolCall, error) {
	client := openai.NewClient(os.Getenv(key))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:      openai.GPT3Dot5Turbo,
			Messages:   messages,
			Tools:      tools,
			ToolChoice: "required",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned")
	}

	choice := resp.Choices[0]

	if len(choice.Message.ToolCalls) == 0 {
		return nil, fmt.Errorf("no tool calls returned")
	}

	toolCall := choice.Message.ToolCalls[0]

	return &toolCall, nil
}

func giveResourcesTool() openai.Tool {
	params := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"target_agent": {
				Type:        jsonschema.Integer,
				Description: "The ID of the agent to give resources to",
			},
			"resource": {
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"type": {
						Type:        jsonschema.String,
						Description: "The type of resource to give (gold or wheat)",
					},
					"amount": {
						Type:        jsonschema.Integer,
						Description: "The amount of the resource to give",
					},
				},
			},
		},
		Required: []string{"target_agent", "resource"},
	}

	f := openai.FunctionDefinition{
		Name:        "give_resources",
		Description: "Give resources to another agent",
		Parameters:  params,
	}

	return openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f,
	}
}

func sendMessageTool() openai.Tool {
	params := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"target_agent": {
				Type:        jsonschema.Integer,
				Description: "The ID of the agent to send a message to",
			},
			"message": {
				Type:        jsonschema.String,
				Description: "The message to send to the target agent",
			},
		},
		Required: []string{"target_agent", "message"},
	}

	f := openai.FunctionDefinition{
		Name:        "send_message",
		Description: "Send a message to another agent",
		Parameters:  params,
	}

	return openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f,
	}
}

func buyBuildingTool() openai.Tool {
	params := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"building_type": {
				Type:        jsonschema.String,
				Description: "The type of building to buy (Farm or Mine)",
			},
		},
		Required: []string{"building_type"},
	}

	f := openai.FunctionDefinition{
		Name:        "buy_building",
		Description: "Buy a building",
		Parameters:  params,
	}

	return openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f,
	}
}

func buyWorkerTool() openai.Tool {
	params := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"count": {
				Type:        jsonschema.Integer,
				Description: "The number of workers to buy",
			},
		},
		Required: []string{"count"},
	}

	f := openai.FunctionDefinition{
		Name:        "buy_worker",
		Description: "Buy workers",
		Parameters:  params,
	}

	return openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f,
	}
}

func endTurnTool() openai.Tool {
	f := openai.FunctionDefinition{
		Name:        "end_turn",
		Description: "End the current turn even if you still have actions left. You do not need to use this tool, as the turn will end automatically after all actions are taken.",
	}

	return openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f,
	}
}

func manBuildingTool() openai.Tool {
	params := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"building_type": {
				Type:        jsonschema.String,
				Description: "The type of building that you want to allocate workers to",
			},
		},
		Required: []string{"building_type"},
	}

	f := openai.FunctionDefinition{
		Name:        "man_building",
		Description: "Allocate workers to a building",
		Parameters:  params,
	}

	return openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f,
	}
}

func unmanBuildingTool() openai.Tool {
	params := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"building_type": {
				Type:        jsonschema.String,
				Description: "The type of building that you want to deallocate workers from",
			},
		},
		Required: []string{"building_type"},
	}

	f := openai.FunctionDefinition{
		Name:        "unman_building",
		Description: "Deallocate workers from a building, freeing them up to work in other buildings",
		Parameters:  params,
	}

	return openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f,
	}
}

func getToolDefinitions() []openai.Tool {
	return []openai.Tool{
		giveResourcesTool(),
		sendMessageTool(),
		buyBuildingTool(),
		buyWorkerTool(),
		manBuildingTool(),
		unmanBuildingTool(),
		endTurnTool(),
	}
}
