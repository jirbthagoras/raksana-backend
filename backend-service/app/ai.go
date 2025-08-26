package app

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

const (
	TrashScanner       int8 = 0
	Ecoach             int8 = 1
	Recap              int8 = 2
	ChallengeGenerator int8 = 3
)

type AIClient struct {
	Genai *genai.Client
}

func InitAiClient(cnf *viper.Viper) *AIClient {
	ctx := context.Background()

	geminiApiKey := cnf.GetString("GEMINI_API_KEY")

	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiApiKey))
	if err != nil {
		panic(err)
	}

	return &AIClient{
		Genai: client,
	}
}

func initModel(client *genai.Client, cnf *viper.Viper, modelType int8) (*genai.GenerativeModel, error) {
	model := cnf.GetString("MODEL")
	generativeModel := client.GenerativeModel(model)

	var systemInstruction = ""
	switch modelType {
	case TrashScanner:
		systemInstruction = cnf.GetString("TRASH_SCANNER_SYSTEM_INSTRUCTION")
		trashScannerConfig(generativeModel)
	case Recap:
		systemInstruction = cnf.GetString("RECAP_SYSTEM_INSTRUCTION")
		recapConfig(generativeModel)
	case Ecoach:
		systemInstruction = cnf.GetString("HABIT_GENERATOR_SYSTEM_INSTRUCTION")
		ecoachConfig(generativeModel)
	case ChallengeGenerator:
		systemInstruction = cnf.GetString("CHALLENGE_GENERATOR_SYSTEM_INSTRUCTION")
		challengeGeneratorConfig(generativeModel)
	}

	generativeModel.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(systemInstruction),
		},
	}

	return generativeModel, nil
}

func trashScannerConfig(generativeModel *genai.GenerativeModel) {
	generativeModel.SetTemperature(1.6)
	generativeModel.SetTopK(40)
	generativeModel.SetTopP(0.95)
	generativeModel.SetMaxOutputTokens(8192)
	generativeModel.ResponseMIMEType = "application/json"
	generativeModel.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title": {
				Type: genai.TypeString,
			},
			"description": {
				Type: genai.TypeString,
			},
			"recycling_ideas": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"name": {
							Type: genai.TypeString,
						},
						"type": {
							Type: genai.TypeString,
							Enum: []string{"recycle", "reuse", "upcycle"},
						},
						"description": {
							Type: genai.TypeString,
						},
						"value": {
							Type: genai.TypeString,
							Enum: []string{"high", "mid", "low"},
						},
					},
					Required: []string{"name", "type", "description", "value"},
				},
			},
		},
		Required: []string{"title", "description", "recycling_ideas"},
	}
}
func ecoachConfig(generativeModel *genai.GenerativeModel) {
	generativeModel.SetTemperature(1.6)
	generativeModel.SetTopK(40)
	generativeModel.SetTopP(0.95)
	generativeModel.SetMaxOutputTokens(8192)
	generativeModel.ResponseMIMEType = "application/json"
	generativeModel.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"habits_generated": {
				Type: genai.TypeInteger,
			},
			"task_per_day": {
				Type: genai.TypeInteger,
			},
			"habits": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"name": {
							Type: genai.TypeString,
						},
						"description": {
							Type: genai.TypeString,
						},
						"difficulty": {
							Type: genai.TypeString,
							Enum: []string{"hard", "normal", "easy"},
						},
					},
					Required: []string{"name", "description", "difficulty"},
				},
			},
		},
		Required: []string{"habits_generated", "task_per_day"},
	}
}
func recapConfig(generativeModel *genai.GenerativeModel) {
	generativeModel.SetTemperature(1.6)
	generativeModel.SetTopK(40)
	generativeModel.SetTopP(0.95)
	generativeModel.SetMaxOutputTokens(8192)
	generativeModel.ResponseMIMEType = "application/json"
	generativeModel.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"growth_rating": {
				Type: genai.TypeString,
				Enum: []string{"bad", "decent", "good", "excellent"},
			},
			"description": {
				Type: genai.TypeString,
			},
			"tips": {
				Type: genai.TypeString,
			},
		},
	}
}
func challengeGeneratorConfig(generativeModel *genai.GenerativeModel) {
	generativeModel.SetTemperature(1.6)
	generativeModel.SetTopK(40)
	generativeModel.SetTopP(0.95)
	generativeModel.SetMaxOutputTokens(8192)
	generativeModel.ResponseMIMEType = "application/json"
	generativeModel.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"name": {
				Type: genai.TypeString,
			},
			"description": {
				Type: genai.TypeString,
			},
		},
	}
}
