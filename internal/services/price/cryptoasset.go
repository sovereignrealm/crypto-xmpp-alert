package price

type CrytpoAsset string

const (
	CrytpoAssetBitcoin           CrytpoAsset = "Bitcoin"
	CrytpoAssetEthereum          CrytpoAsset = "Ethereum"
	CrytpoAssetCardano           CrytpoAsset = "Cardano"
	CrytpoAssetPolkadot          CrytpoAsset = "Polkadot"
	CrytpoAssetUnsupportedCrypto CrytpoAsset = "UnsupportedCrypto"
)
