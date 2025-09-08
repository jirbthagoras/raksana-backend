package models

type RequestScannedObject struct {
	Name       string     `json:"name"`
	Categories []Category `json:"categories"`
}

type Category struct {
	Name string `json:"name"`
}
