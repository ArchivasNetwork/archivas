package consensus

// DifficultyParamsID identifies the difficulty computation parameters
// v1.1.1: Prevents incompatible nodes from joining the network
const (
	// DifficultyParamsID must match genesis file and all peers
	DifficultyParamsID = "v2-normalized-qmax1e12"
	
	// ProtocolVersion for handshake validation
	ProtocolVersion = "v1.1.1-ibd"
)

