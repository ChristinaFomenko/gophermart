package model

type Order struct {
	Number uint64
	Status Status
}

type OrderAccrual struct {
	Order   uint64  `json:"order,string"`
	Status  Status  `json:"status"`
	Accrual float32 `json:"accrual,omitempty"`
}
