package transaction

type Transaction struct {
	TxId   string
	Input  []Input
	Output []Output
}

type Input struct {
	PrevTxId    string
	OutputIndex int
	ScriptSig   string
}

type Output struct {
	Value        float64
	ScriptPubKey string
}
