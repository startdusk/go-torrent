package types

const SHALEN int = 20

// Port to listen on
const Port uint16 = 6881

type InfoHash = [SHALEN]byte

type PeerID = [SHALEN]byte

type PieceHashes = []Piece
type Piece = [SHALEN]byte
