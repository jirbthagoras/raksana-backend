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
	Name             string `json:"name"`
	Description      string `json:"description"`
	Value            string `json:"value"`
	HavingGreenprint bool   `json:"having_greenprint"`
}
