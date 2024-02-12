package server

type Response struct {
	InjectivePeggo struct {
		Status            bool   `json:"status"`
		PendingBatchSize  int    `json:"pending-batch-size"`
		PendingValsetSize int    `json:"pending-valset-size"`
		Nonce             string `json:"nonce(event / observed)"`
		NonceFalseCount   int    `json:"nonce_false_count"`
	} `json:"injective-peggo"`
}
