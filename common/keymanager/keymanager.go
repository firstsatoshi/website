package keymanager

import (
	"crypto/sha256"
	"fmt"
	"sync"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// Refacotr FROM btcutil https://github.com/modood/btckeygen/blob/master/main.go

// Purpose BIP43 - Purpose Field for Deterministic Wallets
// https://github.com/bitcoin/bips/blob/master/bip-0043.mediawiki
//
// Purpose is a constant set to 44' (or 0x8000002C) following the BIP43 recommendation.
// It indicates that the subtree of this node is used according to this specification.
//
// What does 44' mean in BIP44?
// https://bitcoin.stackexchange.com/questions/74368/what-does-44-mean-in-bip44
//
// 44' means that hardened keys should be used. The distinguisher for whether
// a key a given index is hardened is that the index is greater than 2^31,
// which is 2147483648. In hex, that is 0x80000000. That is what the apostrophe (') means.
// The 44 comes from adding it to 2^31 to get the final hardened key index.
// In hex, 44 is 2C, so 0x80000000 + 0x2C = 0x8000002C.
type Purpose = uint32

const (
	PurposeBIP44 Purpose = 0x8000002C // 44' BIP44
	PurposeBIP49 Purpose = 0x80000031 // 49' BIP49
	PurposeBIP84 Purpose = 0x80000054 // 84' BIP84
)

// CoinType SLIP-0044 : Registered coin types for BIP-0044
// https://github.com/satoshilabs/slips/blob/master/slip-0044.md
type CoinType = uint32

const (
	CoinTypeBTC CoinType = 0x80000000
	CoinTypeLTC CoinType = 0x80000002
	CoinTypeETH CoinType = 0x8000003c
	CoinTypeEOS CoinType = 0x800000c2
)

const (
	Apostrophe uint32 = 0x80000000 // 0'
)

type Key struct {
	path     string
	bip32Key *bip32.Key
}

func (k *Key) getWifKeyAndAddress(compress bool, chainCfg chaincfg.Params) (wif, legacyAddress, p2trAddress string, err error) {
	prvKey, _ := btcec.PrivKeyFromBytes(k.bip32Key.Key)

	wif, legacyAddress, p2trAddress, _, err = generateFromBytes(prvKey, compress, chainCfg)
	return
}

// https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki
// bip44 define the following 5 levels in BIP32 path:
// m / purpose' / coin_type' / account' / change / address_index

func (k *Key) GetPath() string {
	return k.path
}

type KeyManager struct {
	mnemonic   string
	passphrase string
	keys       map[string]*bip32.Key

	chainCfg chaincfg.Params
	mux      sync.Mutex
}

// NewKeyManager return new key manager
// bitSize has to be a multiple 32 and be within the inclusive range of {128, 256}
// 128: 12 phrases
// 256: 24 phrases
// func newKeyManager(bitSize int, mnemonic string) (*KeyManager, error) {

// 	if mnemonic == "" {
// 		entropy, err := bip39.NewEntropy(bitSize)
// 		if err != nil {
// 			return nil, err
// 		}
// 		mnemonic, err = bip39.NewMnemonic(entropy)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	km := &KeyManager{
// 		mnemonic: mnemonic,
// 		keys:     make(map[string]*bip32.Key, 0),
// 	}
// 	return km, nil
// }

func NewKeyManagerFromSeed(seed string, chainCfg chaincfg.Params) (*KeyManager, error) {

	if seed == "" {
		seed = "youngqqcn@163.com20230529"
	}

	msg := []byte(seed)
	salt := []byte{0x20, 0x23, 0x05, 0x9, 0xC, 0x25}
	msg = append(msg, salt...)

	h := sha256.Sum256([]byte(msg))

	// 128 bits  for 12 words mnemonic
	entropy := h[:16]
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}

	km := &KeyManager{
		mnemonic:   mnemonic,
		passphrase: "godbless@you20230529",
		chainCfg:   chainCfg,
		keys:       make(map[string]*bip32.Key, 0),
	}
	return km, nil
}

func (km *KeyManager) getMnemonic() string {
	return km.mnemonic
}

func (km *KeyManager) getPassphrase() string {
	return km.passphrase
}

func (km *KeyManager) getSeed() []byte {
	return bip39.NewSeed(km.getMnemonic(), km.getPassphrase())
}

func (km *KeyManager) getKey(path string) (*bip32.Key, bool) {
	km.mux.Lock()
	defer km.mux.Unlock()

	key, ok := km.keys[path]
	return key, ok
}

func (km *KeyManager) setKey(path string, key *bip32.Key) {
	km.mux.Lock()
	defer km.mux.Unlock()

	km.keys[path] = key
}

func (km *KeyManager) getMasterKey() (*bip32.Key, error) {
	path := "m"

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	key, err := bip32.NewMasterKey(km.getSeed())
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

func (km *KeyManager) getPurposeKey(purpose uint32) (*bip32.Key, error) {
	path := fmt.Sprintf(`m/%d'`, purpose-Apostrophe)

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	parent, err := km.getMasterKey()
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(purpose)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

func (km *KeyManager) getCoinTypeKey(purpose, coinType uint32) (*bip32.Key, error) {
	path := fmt.Sprintf(`m/%d'/%d'`, purpose-Apostrophe, coinType-Apostrophe)

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	parent, err := km.getPurposeKey(purpose)
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(coinType)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

func (km *KeyManager) getAccountKey(purpose, coinType, account uint32) (*bip32.Key, error) {
	path := fmt.Sprintf(`m/%d'/%d'/%d'`, purpose-Apostrophe, coinType-Apostrophe, account)

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	parent, err := km.getCoinTypeKey(purpose, coinType)
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(account + Apostrophe)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

// getChangeKey ...
// https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki#change
// change constant 0 is used for external chain
// change constant 1 is used for internal chain (also known as change addresses)
func (km *KeyManager) getChangeKey(purpose, coinType, account, change uint32) (*bip32.Key, error) {
	path := fmt.Sprintf(`m/%d'/%d'/%d'/%d`, purpose-Apostrophe, coinType-Apostrophe, account, change)

	key, ok := km.getKey(path)
	if ok {
		return key, nil
	}

	parent, err := km.getAccountKey(purpose, coinType, account)
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(change)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return key, nil
}

func (km *KeyManager) getPrivateKey(purpose, coinType, account, change, index uint32) (*Key, error) {
	path := fmt.Sprintf(`m/%d'/%d'/%d'/%d/%d`, purpose-Apostrophe, coinType-Apostrophe, account, change, index)

	key, ok := km.getKey(path)
	if ok {
		return &Key{path: path, bip32Key: key}, nil
	}

	parent, err := km.getChangeKey(purpose, coinType, account, change)
	if err != nil {
		return nil, err
	}

	key, err = parent.NewChildKey(index)
	if err != nil {
		return nil, err
	}

	km.setKey(path, key)

	return &Key{path: path, bip32Key: key}, nil
}

func (km *KeyManager) GetWifKeyAndAddresss(accoutIndex, addrIndex uint32) (wif, p2trAddress string, err error) {
	k, err := km.getPrivateKey(PurposeBIP44, CoinTypeBTC, accoutIndex, 0, addrIndex)
	if err != nil {
		return
	}
	compressed := true
	wif, _, p2trAddress, err = k.getWifKeyAndAddress(compressed, km.chainCfg)
	if err != nil {
		return
	}
	return
}

// func Generate(compress bool) (wif, address, segwitBech32, segwitNested string, err error) {
// 	prvKey, err := btcec.NewPrivateKey()
// 	if err != nil {
// 		return "", "", "", "", err
// 	}
// 	return generateFromBytes(prvKey, compress)
// }

func generateFromBytes(prvKey *btcec.PrivateKey, compress bool, chainCfg chaincfg.Params) (wif, address, taprootBech32, segwitNested string, err error) {
	// generate the wif(wallet import format) string
	btcwif, err := btcutil.NewWIF(prvKey, &chainCfg, compress)
	if err != nil {
		return "", "", "", "", err
	}
	wif = btcwif.String()

	// generate a normal p2pkh address
	serializedPubKey := btcwif.SerializePubKey()
	addressPubKey, err := btcutil.NewAddressPubKey(serializedPubKey, &chainCfg)
	if err != nil {
		return "", "", "", "", err
	}
	address = addressPubKey.EncodeAddress()

	addressWitnessPubKeyHash, err := btcutil.NewAddressTaproot(
		schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(btcwif.PrivKey.PubKey())),
		&chainCfg)

	if err != nil {
		return "", "", "", "", err
	}
	taprootBech32 = addressWitnessPubKeyHash.EncodeAddress()

	// generate an address which is
	// backwards compatible to Bitcoin nodes running 0.6.0 onwards, but
	// allows us to take advantage of segwit's scripting improvments,
	// and malleability fixes.
	serializedScript, err := txscript.PayToAddrScript(addressWitnessPubKeyHash)
	if err != nil {
		return "", "", "", "", err
	}
	addressScriptHash, err := btcutil.NewAddressScriptHash(serializedScript, &chainCfg)
	if err != nil {
		return "", "", "", "", err
	}
	segwitNested = addressScriptHash.EncodeAddress()

	return wif, address, taprootBech32, segwitNested, nil
}

// func xmain() {
// 	compress := true // generate a compressed public key
// 	bip39 := flag.Bool("bip39", false, "mnemonic code for generating deterministic keys")
// 	pass := flag.String("pass", "", "protect bip39 mnemonic with a passphrase")
// 	number := flag.Int("n", 10, "set number of keys to generate")
// 	mnemonic := flag.String("mnemonic", "", "optional list of words to re-generate a root key")

// 	flag.Parse()

// 	if !*bip39 {
// 		fmt.Printf("\n%-34s %-52s %-42s %s\n", "Bitcoin Address", "WIF(Wallet Import Format)", "SegWit(bech32)", "SegWit(nested)")
// 		fmt.Println(strings.Repeat("-", 165))

// 		for i := 0; i < *number; i++ {
// 			wif, address, segwitBech32, segwitNested, err := Generate(compress)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			fmt.Printf("%-34s %s %s %s\n", address, wif, segwitBech32, segwitNested)
// 		}
// 		fmt.Println()
// 		return
// 	}

// 	km, err := NewKeyManager(128, *pass, *mnemonic)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	masterKey, err := km.GetMasterKey()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	passphrase := km.GetPassphrase()
// 	if passphrase == "" {
// 		passphrase = "<none>"
// 	}
// 	fmt.Printf("\n%-18s %s\n", "BIP39 Mnemonic:", km.GetMnemonic())
// 	fmt.Printf("%-18s %s\n", "BIP39 Passphrase:", passphrase)
// 	fmt.Printf("%-18s %x\n", "BIP39 Seed:", km.GetSeed())
// 	fmt.Printf("%-18s %s\n", "BIP32 Root Key:", masterKey.B58Serialize())

// 	fmt.Printf("\n%-18s %-34s %-52s\n", "Path(BIP44)", "Bitcoin Address", "WIF(Wallet Import Format)")
// 	fmt.Println(strings.Repeat("-", 106))
// 	for i := 0; i < *number; i++ {
// 		key, err := km.GetKey(PurposeBIP44, CoinTypeBTC, 0, 0, uint32(i))
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		wif, address, _, _, err := key.Encode(compress)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		fmt.Printf("%-18s %-34s %s\n", key.GetPath(), address, wif)
// 	}

// 	fmt.Printf("\n%-18s %-34s %s\n", "Path(BIP49)", "SegWit(nested)", "WIF(Wallet Import Format)")
// 	fmt.Println(strings.Repeat("-", 106))
// 	for i := 0; i < *number; i++ {
// 		key, err := km.GetKey(PurposeBIP49, CoinTypeBTC, 0, 0, uint32(i))
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		wif, _, _, segwitNested, err := key.Encode(compress)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		fmt.Printf("%-18s %s %s\n", key.GetPath(), segwitNested, wif)
// 	}

// 	fmt.Printf("\n%-18s %-42s %s\n", "Path(BIP84)", "SegWit(bech32)", "WIF(Wallet Import Format)")
// 	fmt.Println(strings.Repeat("-", 114))
// 	for i := 0; i < *number; i++ {
// 		key, err := km.GetKey(PurposeBIP84, CoinTypeBTC, 0, 0, uint32(i))
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		wif, _, segwitBech32, _, err := key.Encode(compress)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		fmt.Printf("%-18s %s %s\n", key.GetPath(), segwitBech32, wif)
// 	}
// 	fmt.Println()
// }
