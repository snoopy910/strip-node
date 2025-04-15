package libs

type Operation struct {
	ID               int64  `json:"id"`
	SerializedTxn    string `json:"serializedTxn"`
	DataToSign       string `json:"dataToSign"`
	ChainId          string `json:"chainId"`
	GenesisHash      string `json:"genesisHash"`
	KeyCurve         string `json:"keyCurve"`
	Status           string `json:"status"`
	Result           string `json:"result"`
	Type             string `json:"type"`
	Solver           string `json:"solver"`
	SolverMetadata   string `json:"solverMetadata"`
	SolverDataToSign string `json:"solverDataToSign"`
	SolverOutput     string `json:"solverOutput"`
}

type Intent struct {
	ID            int64       `json:"id"`
	Operations    []Operation `json:"operations"`
	Signature     string      `json:"signature"`
	Identity      string      `json:"identity"`
	IdentityCurve string      `json:"identityCurve"`
	Status        string      `json:"status"`
	Expiry        uint64      `json:"expiry"`
	CreatedAt     uint64      `json:"createdAt"`
}
