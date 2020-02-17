package cmd

type Resolution struct {
	Bin    string `json:"bin"`
	Npm    string `json:"npm"`
	Github string `json:"github"`
}

type Arg struct {
	Description string `json:"description"`
	Type        string `json:"type"`
	IsRequired  bool   `json:"isRequired"`
}

type Flag struct {
	Description string `json:"description"`
	Short       string `json:"short"`
	Long        string `json:"long"`
	Type        string `json:"type"`
	IsRequired  bool   `json:"isRequired"`
	Default     string `json:"default"`
}

type Inputs struct {
	Args  []Arg  `json:"args"`
	Flags []Flag `json:"flags"`
}

type Recipe struct {
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Inputs      Inputs   `json:"inputs"`
}

type Response struct {
	Resolution Resolution `json:"resolution"`
	Keywords   []string   `json:"keywords"`
	Inputs     Inputs     `json:"inputs"`
	Recipes    []Recipe   `json:"recipes"`
}
