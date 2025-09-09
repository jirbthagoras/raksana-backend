package models

type RequestScannedObject struct {
	Name       string     `json:"name"`
	Categories []Category `json:"categories"`
}

type Category struct {
	Name string `json:"name"`
}

type AIResponseScan struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	ImageKey    string          `json:"image_key"`
	Items       []ResponseItems `json:"items"`
}

type ResponseItems struct {
	Id               int    `json:"id,omitempty"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Value            string `json:"value"`
	HavingGreenprint bool   `json:"having_greenprint"`
}

type AIResponseGreenprint struct {
	Title               string             `json:"'title"`
	Description         string             `json:"description"`
	SustainabilityScore string             `json:"sustainability_score"`
	EstimatedTime       string             `json:"estimated_time"`
	Tools               []ResponseTool     `json:"tools"`
	Materials           []ResponseMaterial `json:"materials"`
	Steps               []ResponseStep     `json:"steps"`
	Text                string
	CreatedAt           string `json:"created_at,omitempty"`
}

type ResponseTool struct {
	ID           int64  `json:"id,omitempty"`
	GreenprintID int64  `json:"greenprint_id,omitempty"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Price        int32  `json:"price"`
}

type ResponseMaterial struct {
	ID           int64  `json:"id,omitempty"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Price        int32  `json:"price"`
	Quantity     int32  `json:"Quantity"`
	GreenprintID int64  `json:"greenprint_id,omitempty"`
}

type ResponseStep struct {
	ID           int64  `json:"id,omitempty"`
	GreenprintID int64  `json:"greenprint_id,omitempty"`
	Description  string `json:"description"`
}
