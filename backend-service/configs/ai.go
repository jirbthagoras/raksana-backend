package configs

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

const (
	TrashScanner int8 = 0
	Ecoach       int8 = 1
	RecapMonthly int8 = 2
	RecapWeekly  int8 = 3
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

func InitModel(client *genai.Client, cnf *viper.Viper, modelType int8) (*genai.GenerativeModel, error) {
	model := cnf.GetString("MODEL")
	generativeModel := client.GenerativeModel(model)

	var systemInstruction = ""
	switch modelType {
	case TrashScanner:
		systemInstruction = cnf.GetString("TRASH_SCANNER_SYSTEM_INSTRUCTION")
		trashScannerConfig(generativeModel)
	case RecapMonthly:
		systemInstruction = cnf.GetString("MONTHLY_RECAP_SYSTEM_INSTRUCTION")
		recapConfig(generativeModel)
	case RecapWeekly:
		systemInstruction = cnf.GetString("WEEKLY_RECAP_SYSTEM_INSTRUCTION")
		recapConfig(generativeModel)
	case Ecoach:
		systemInstruction = cnf.GetString("ECOACH_SYSTEM_INSTRUCTION")
		ecoachConfig(generativeModel)
	default:
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
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
				Type: genai.TypeString},
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
			"name": {
				Type: genai.TypeString,
			},
			"expected_task": {
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
		Required: []string{"expected_task", "task_per_day"},
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
				Enum: []string{"1", "2", "3", "4", "5"},
			},
			"summary": {
				Type: genai.TypeString,
			},
			"tips": {
				Type: genai.TypeString,
			},
		},
	}
}

// func challengeGeneratorConfig(generativeModel *genai.GenerativeModel) {
// 	generativeModel.SetTemperature(1.6)
// 	generativeModel.SetTopK(40)
// 	generativeModel.SetTopP(0.95)
// 	generativeModel.SetMaxOutputTokens(8192)
// 	generativeModel.ResponseMIMEType = "application/json"
// 	generativeModel.ResponseSchema = &genai.Schema{
// 		Type: genai.TypeObject,
// 		Properties: map[string]*genai.Schema{
// 			"name": {
// 				Type: genai.TypeString,
// 			},
// 			"description": {
// 				Type: genai.TypeString,
// 			},
// 		},
// 	}
// }
