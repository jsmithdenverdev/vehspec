package main

type Vehicle struct {
	Vin               string `json:"vin"`
	Year              int64  `json:"year"`
	Make              string `json:"make"`
	Model             string `json:"model"`
	Trim              string `json:"trim"`
	Price             int64  `json:"price"`
	Type              string `json:"type"`
	Mileage           int64  `json:"mileage"`
	Doors             int64  `json:"doors"`
	Body              string `json:"body"`
	Mpg               Mpg    `json:"mpg"`
	IsNew             bool   `json:"is_new"`
	CertifiedPreowned bool   `json:"certified_preowned"`
	Colors            Colors `json:"colors"`
	Engine            Engine `json:"engine"`
	Transmission      string `json:"transmission"`
}

type Colors struct {
	Exterior       string `json:"exterior"`
	ExteriorDetail string `json:"exterior_detail"`
	Interior       string `json:"interior"`
	InteriorDetail string `json:"interior_detail"`
}

type Engine struct {
	Fuel         string  `json:"fuel"`
	Type         string  `json:"type"`
	Cylinders    int64   `json:"cylinders"`
	Displacement float64 `json:"displacement"`
}

type Mpg struct {
	City    int64 `json:"city"`
	Highway int64 `json:"highway"`
}
