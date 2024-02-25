package sequencer

type Operation struct {
	SerializedTxn string
	DataToSign    string
	ChainId       string
	KeyCurve      string
}

type Intent struct {
	Operations    []Operation
	Signature     string
	Identity      string
	IdentityCurve string
}

func StartSequencer() {

}
