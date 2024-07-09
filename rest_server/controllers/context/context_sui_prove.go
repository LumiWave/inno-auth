package context

type ReqProve struct {
	MaxEpoch                   int64  `json:"maxEpoch"`
	JwtRandomness              string `json:"jwtRandomness"`
	Jwt                        string `json:"jwt"`
	KeyClaimName               string `json:"keyClaimName"`
	ExtendedEphemeralPublicKey string `json:"extendedEphemeralPublicKey"`
	EphemeralPublicKey         string `json:"ephemeralPublicKey"`
	Salt                       string `json:"salt"`
}

type ResProve struct {
	Name        string `json:"name"`
	Message     string `json:"message"`
	ProofPoints struct {
		A []string   `json:"a"`
		B [][]string `json:"b"`
		C []string   `json:"c"`
	} `json:"proofPoints"`
	IssBase64Details struct {
		Value     string `json:"value"`
		IndexMod4 int64  `json:"indexMod4"`
	} `json:"issBase64Details"`
	HeaderBase64 string `json:"headerBase64"`
}

type ReqProveNonce struct {
	EphemeralPublicKey string `json:"ephemeralPublicKey"`
}

type ResProveNonce struct {
	Name        string `json:"name"`
	Message     string `json:"message"`
	ProofPoints struct {
		A []string   `json:"a"`
		B [][]string `json:"b"`
		C []string   `json:"c"`
	} `json:"proofPoints"`
	IssBase64Details struct {
		Value     string `json:"value"`
		IndexMod4 int64  `json:"indexMod4"`
	} `json:"issBase64Details"`
	HeaderBase64 string `json:"headerBase64"`
}
