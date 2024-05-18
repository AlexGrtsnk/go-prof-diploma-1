package internal

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Orders struct {
	Number   string `json:"number"`
	Status   string `json:"status"`
	Accrual  int    `json:"accrual"`
	Uoloaded string `json:"uploaded"`
}
type OrdersList []Orders

type UserAnsw struct {
	Accural  int `json:"accural"`
	Withdraw int `json:"withdrawn"`
}

type BalanceAnsw struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}

type DrawAnsw struct {
	Number       string `json:"number"`
	Sum          int    `json:"sum"`
	ProccessedAt string `json:"processed_at"`
}

type WithAnsw struct {
	Order   string `json:"number"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}

type DrawAnswList []DrawAnsw
