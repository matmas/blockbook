// +build unittest

package capricoin

import (
	"blockbook/bchain"
	"blockbook/bchain/coins/btc"
	"encoding/hex"
	"math/big"
	"os"
	"reflect"
	"testing"

	"github.com/martinboehm/btcutil/chaincfg"
)

func TestMain(m *testing.M) {
	c := m.Run()
	chaincfg.ResetParams()
	os.Exit(c)
}

func Test_GetAddrDescFromAddress(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "P2PKH",
			args:    args{address: "Ce9aQjmHtMivsvUmy9dZW7YPnc3Xos1ew6"},
			want:    "76a914eddc1ed2170822e3aae1d16f0547248300b0cd0188ac",
			wantErr: false,
		},
		{
			name:    "P2SH",
			args:    args{address: "FRSxKiRpyDkLdus5awTsRyekb5QqMVfiKA"},
			want:    "a914d731bc3d8479ba40ba33fe5d31793d6f091b2f9a87",
			wantErr: false,
		},
	}
	parser := NewCapricoinParser(GetChainParams("main"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.GetAddrDescFromAddress(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddrDescFromAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("GetAddrDescFromAddress() = %v, want %v", h, tt.want)
			}
		})
	}
}

func Test_GetAddrDescFromVout(t *testing.T) {
	type args struct {
		vout bchain.Vout
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "P2PKH",
			args:    args{vout: bchain.Vout{ScriptPubKey: bchain.ScriptPubKey{Hex: "76a914eddc1ed2170822e3aae1d16f0547248300b0cd0188ac"}}},
			want:    "76a914eddc1ed2170822e3aae1d16f0547248300b0cd0188ac",
			wantErr: false,
		},
		{
			name:    "P2PK compressed CTGbsJhVNSJAGTSrVxCbs546CLJPRtkb82",
			args:    args{vout: bchain.Vout{ScriptPubKey: bchain.ScriptPubKey{Hex: "2102be69865a3551c2defe6a55860250847819d0261e83c5abeb27df23c7a791032dac"}}},
			want:    "76a9147686dd9a5a7e3ac121273d4420ef4685db31d6be88ac",
			wantErr: false,
		},
		{
			name:    "P2PK uncompressed CZ2oddmtUKtwWpisXukHY37QnzGkk3Pb1w",
			args:    args{vout: bchain.Vout{ScriptPubKey: bchain.ScriptPubKey{Hex: "4104be69865a3551c2defe6a55860250847819d0261e83c5abeb27df23c7a791032d5604bea92bc7ad1cd83b0e096eab3118d084aaca39522d909c5c99e77447ef84ac"}}},
			want:    "76a914b5bb9e017818127e6083daf4461088b604e9606588ac",
			wantErr: false,
		},
		{
			name:    "P2SH",
			args:    args{vout: bchain.Vout{ScriptPubKey: bchain.ScriptPubKey{Hex: "a914d731bc3d8479ba40ba33fe5d31793d6f091b2f9a87"}}},
			want:    "a914d731bc3d8479ba40ba33fe5d31793d6f091b2f9a87",
			wantErr: false,
		},
	}
	parser := NewCapricoinParser(GetChainParams("main"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.GetAddrDescFromVout(&tt.args.vout)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddrDescFromVout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("GetAddrDescFromVout() = %v, want %v", h, tt.want)
			}
		})
	}
}

func Test_GetAddressesFromAddrDesc(t *testing.T) {
	type args struct {
		script string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		want2   bool
		wantErr bool
	}{
		{
			name:    "P2PKH",
			args:    args{script: "76a914eddc1ed2170822e3aae1d16f0547248300b0cd0188ac"},
			want:    []string{"Ce9aQjmHtMivsvUmy9dZW7YPnc3Xos1ew6"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "P2PK compressed",
			args:    args{script: "2102be69865a3551c2defe6a55860250847819d0261e83c5abeb27df23c7a791032dac"},
			want:    []string{"CTGbsJhVNSJAGTSrVxCbs546CLJPRtkb82"},
			want2:   false,
			wantErr: false,
		},
		{
			name:    "P2PK uncompressed",
			args:    args{script: "4104be69865a3551c2defe6a55860250847819d0261e83c5abeb27df23c7a791032d5604bea92bc7ad1cd83b0e096eab3118d084aaca39522d909c5c99e77447ef84ac"},
			want:    []string{"CZ2oddmtUKtwWpisXukHY37QnzGkk3Pb1w"},
			want2:   false,
			wantErr: false,
		},
		{
			name:    "P2SH",
			args:    args{script: "a914d731bc3d8479ba40ba33fe5d31793d6f091b2f9a87"},
			want:    []string{"FRSxKiRpyDkLdus5awTsRyekb5QqMVfiKA"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "OP_RETURN ascii",
			args:    args{script: "6a0461686f6a"},
			want:    []string{"OP_RETURN (ahoj)"},
			want2:   false,
			wantErr: false,
		},
		{
			name:    "OP_RETURN OP_PUSHDATA1 ascii",
			args:    args{script: "6a4c0b446c6f7568792074657874"},
			want:    []string{"OP_RETURN (Dlouhy text)"},
			want2:   false,
			wantErr: false,
		},
		{
			name:    "OP_RETURN hex",
			args:    args{script: "6a072020f1686f6a20"},
			want:    []string{"OP_RETURN 2020f1686f6a20"},
			want2:   false,
			wantErr: false,
		},
	}

	parser := NewCapricoinParser(GetChainParams("main"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.script)
			got, got2, err := parser.GetAddressesFromAddrDesc(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddressesFromAddrDesc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

var (
	testTx1 bchain.Tx
	testTxPacked1 = "0a20b7724cdac78bd4d8ec7e19c9c7b53adfd095afc4ae478fc79a7936a0a6a4a29b12e60101000000eb5a025b0137699807c5b6e84c9fd558226c723f4b33aa53421456535b4431d0434b9d6dea010000006b483045022100ca490f7225874a7177a33c920b5a6ca2fa0e23e5aa6a59d8345335d92d3b478702202ce6504c93c643f15071e2051dc6b39ff2262cff4caef86a3403dac8197219fb012103b23813d1e0e783e825fe3a754935e0d060e6e4c6ace28a3e28fe91ac92585d95ffffffff0278499fe3078b00001976a914c1cd13a152e14586714e3081f7b02563fc3fbefb88ac00ca9a3b000000001976a9144ede88857a16efb6a088c7add2687c40d250875788ac0000000018bfb689d805200028c0c4073299010a001220ea6d9d4b43d031445b5356144253aa334b3f726c2258d59f4ce8b6c5079869371801226b483045022100ca490f7225874a7177a33c920b5a6ca2fa0e23e5aa6a59d8345335d92d3b478702202ce6504c93c643f15071e2051dc6b39ff2262cff4caef86a3403dac8197219fb012103b23813d1e0e783e825fe3a754935e0d060e6e4c6ace28a3e28fe91ac92585d9528ffffffff0f3a490a068b07e39f497810001a1976a914c1cd13a152e14586714e3081f7b02563fc3fbefb88ac222243613863663350715946317576705a56796d574c747a78766f596f4d7056564b76343a470a043b9aca0010011a1976a9144ede88857a16efb6a088c7add2687c40d250875788ac222243506575727454324343666178755856516f41327354534b6f6f4b6f316f79766353400148ebb589d805"
)

func init() {
	testTx1 = bchain.Tx{
		Hex:       "01000000eb5a025b0137699807c5b6e84c9fd558226c723f4b33aa53421456535b4431d0434b9d6dea010000006b483045022100ca490f7225874a7177a33c920b5a6ca2fa0e23e5aa6a59d8345335d92d3b478702202ce6504c93c643f15071e2051dc6b39ff2262cff4caef86a3403dac8197219fb012103b23813d1e0e783e825fe3a754935e0d060e6e4c6ace28a3e28fe91ac92585d95ffffffff0278499fe3078b00001976a914c1cd13a152e14586714e3081f7b02563fc3fbefb88ac00ca9a3b000000001976a9144ede88857a16efb6a088c7add2687c40d250875788ac00000000",
		Txid:      "b7724cdac78bd4d8ec7e19c9c7b53adfd095afc4ae478fc79a7936a0a6a4a29b",
		Version:   1,
		Time:      1526881003,
		LockTime:  0,
		Blocktime: 1526881087,
		Vin: []bchain.Vin{
			{
				Txid:     "ea6d9d4b43d031445b5356144253aa334b3f726c2258d59f4ce8b6c507986937",
				Vout:     1,
				ScriptSig: bchain.ScriptSig{
					Hex: "483045022100ca490f7225874a7177a33c920b5a6ca2fa0e23e5aa6a59d8345335d92d3b478702202ce6504c93c643f15071e2051dc6b39ff2262cff4caef86a3403dac8197219fb012103b23813d1e0e783e825fe3a754935e0d060e6e4c6ace28a3e28fe91ac92585d95",
				},
				Sequence: 4294967295,
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat: *big.NewInt(152865999899000),
				N:        0,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "76a914c1cd13a152e14586714e3081f7b02563fc3fbefb88ac",
					Addresses: []string{
						"Ca8cf3PqYF1uvpZVymWLtzxvoYoMpVVKv4",
					},
				},
			},
			{
				ValueSat: *big.NewInt(1000000000),
				N:        1,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "76a9144ede88857a16efb6a088c7add2687c40d250875788ac",
					Addresses: []string{
						"CPeurtT2CCfaxuXVQoA2sTSKooKo1oyvcS",
					},
				},
			},
		},
	}
}

func Test_PackTx(t *testing.T) {
	type args struct {
		tx        bchain.Tx
		height    uint32
		blockTime int64
		parser    *CapricoinParser
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "cpc-1",
			args: args{
				tx:        testTx1,
				height:    123456,
				blockTime: 1526881087,
				parser:    NewCapricoinParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    testTxPacked1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.parser.PackTx(&tt.args.tx, tt.args.height, tt.args.blockTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("packTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("packTx() = %v, want %v", h, tt.want)
			}
		})
	}
}

func Test_UnpackTx(t *testing.T) {
	type args struct {
		packedTx string
		parser   *CapricoinParser
	}
	tests := []struct {
		name    string
		args    args
		want    *bchain.Tx
		want1   uint32
		wantErr bool
	}{
		{
			name: "cpc-1",
			args: args{
				packedTx: testTxPacked1,
				parser:   NewCapricoinParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    &testTx1,
			want1:   123456,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.packedTx)
			got, got1, err := tt.args.parser.UnpackTx(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpackTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unpackTx() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("unpackTx() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
