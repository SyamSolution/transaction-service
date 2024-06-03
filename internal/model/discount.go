package model

type Discount struct {
	IsContinentSoldOut bool    `json:"isContinentSoldOut"`
	IsContinentDiff    bool    `json:"isContinentDiff"`
	DiscountAmount     float32 `json:"discountAmount"`
}

type ResponseDiscount struct {
	Data float32 `json:"data"`
	Meta Meta    `json:"meta"`
}
