package colx

import (
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/giangnamnabka/btcutil/base58"
)

// ConvertCOLXAddress decode ltc address and convert to BTC address
// nolint:gocyclo // keep it
func (b *Bridge) ConvertCOLXAddress(addr, net string) (address btcutil.Address, err error) {
	bchainConfig := &chaincfg.MainNetParams
	cchainConfig := b.GetChainParams()

	// Serialized public keys are either 65 bytes (130 hex chars) if
	// uncompressed/hybrid or 33 bytes (66 hex chars) if compressed.
	if len(addr) == 130 || len(addr) == 66 {
		serializedPubKey, errf := hex.DecodeString(addr)
		if errf != nil {
			return nil, errf
		}
		return btcutil.NewAddressPubKey(serializedPubKey, bchainConfig)
	}

	// Switch on decoded length to determine the type.
	decoded, netID, err := base58.CheckDecode(addr)
	if err != nil {
		if err == base58.ErrChecksum {
			return nil, btcutil.ErrChecksumMismatch
		}
		return nil, errors.New("decoded address is of unknown format")
	}
	switch len(decoded) {
	case 20: // P2PKH or P2SH
		isP2PKH := netID == cchainConfig.PubKeyHashAddrID
		isP2SH := netID == cchainConfig.ScriptHashAddrID
		switch hash160 := decoded; {
		case isP2PKH && isP2SH:
			return nil, btcutil.ErrAddressCollision
		case isP2PKH:
			return btcutil.NewAddressPubKeyHash(hash160, bchainConfig)
		case isP2SH:
			return btcutil.NewAddressScriptHashFromHash(hash160, bchainConfig)
		default:
			return nil, btcutil.ErrUnknownAddressType
		}

	default:
		return nil, errors.New("decoded address is of unknown size")
	}
}
