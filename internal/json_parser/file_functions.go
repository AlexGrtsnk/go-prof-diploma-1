package internal

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Orders struct {
	Number   string  `json:"number"`
	Status   string  `json:"status"`
	Accrual  float64 `json:"accrual"`
	Uoloaded string  `json:"uploaded"`
}
type OrdersList []Orders

type UserAnsw struct {
	Current  float64 `json:"Current"`
	Withdraw float64 `json:"withdrawn"`
}

type BalanceAnsw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type DrawAnsw struct {
	Order        string  `json:"order"`
	Sum          float64 `json:"sum"`
	ProccessedAt string  `json:"processed_at"`
}

type WithAnsw struct {
	Order   string  `json:"number"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

type DrawAnswList []DrawAnsw
