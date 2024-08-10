// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bridge

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ISwapRouterExactInputSingleParams is an auto generated low-level Go binding around an user-defined struct.
type ISwapRouterExactInputSingleParams struct {
	TokenIn           common.Address
	TokenOut          common.Address
	Fee               *big.Int
	Recipient         common.Address
	Deadline          *big.Int
	AmountIn          *big.Int
	AmountOutMinimum  *big.Int
	SqrtPriceLimitX96 *big.Int
}

// BridgeMetaData contains all meta data concerning the Bridge contract.
var BridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"AuthorityChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractISwapRouter\",\"name\":\"swapRouter\",\"type\":\"address\"}],\"name\":\"SwapRouterChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"chainId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"srcToken\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"peggedToken\",\"type\":\"address\"}],\"name\":\"TokenAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TokenBurned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"TokenMinted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"TokenSwapped\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chainId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"srcToken\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"peggedToken\",\"type\":\"address\"}],\"name\":\"addToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authority\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"contractBridgeToken\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"contractBridgeToken\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getBurnMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_messageHash\",\"type\":\"bytes32\"}],\"name\":\"getEthSignedMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"contractBridgeToken\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getMintMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"structISwapRouter.ExactInputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"getSwapMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_authority\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"contractBridgeToken\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"peggedTokens\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_ethSignedMessageHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_signature\",\"type\":\"bytes\"}],\"name\":\"recoverSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_authority\",\"type\":\"address\"}],\"name\":\"setAuthority\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractISwapRouter\",\"name\":\"_swapRouter\",\"type\":\"address\"}],\"name\":\"setSwapRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"splitSignature\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"structISwapRouter.ExactInputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"swap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"swapRouter\",\"outputs\":[{\"internalType\":\"contractISwapRouter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f80fd5b50612a148061001c5f395ff3fe608060405234801561000f575f80fd5b506004361061012a575f3560e01c806397aba7f9116100ab578063c4d66de81161006f578063c4d66de814610340578063dfc344811461035c578063f2fde38b1461038c578063fa3d24f2146103a8578063fa540801146103c45761012a565b806397aba7f914610272578063a7bb5803146102a2578063bf7e214f146102d4578063bfb4ad0c146102f2578063c31c9c07146103225761012a565b8063715018a6116100f2578063715018a6146101ce5780637a9e5e4b146101d85780637b7dfd10146101f45780637ecebe00146102245780638da5cb5b146102545761012a565b80632adfefeb1461012e57806335ae976e1461014a57806341273657146101665780634f334e891461018257806368cab61d1461019e575b5f80fd5b610148600480360381019061014391906118a9565b6103f4565b005b610164600480360381019061015f919061193c565b610629565b005b610180600480360381019061017b9190611a0a565b61085c565b005b61019c60048036038101906101979190611a58565b6108de565b005b6101b860048036038101906101b39190611adc565b610d13565b6040516101c59190611b47565b60405180910390f35b6101d6610dc4565b005b6101f260048036038101906101ed9190611b60565b610dd7565b005b61020e60048036038101906102099190611b8b565b610e58565b60405161021b9190611b47565b60405180910390f35b61023e60048036038101906102399190611b60565b610e99565b60405161024b9190611bfe565b60405180910390f35b61025c610eae565b6040516102699190611c26565b60405180910390f35b61028c60048036038101906102879190611c69565b610ee3565b6040516102999190611c26565b60405180910390f35b6102bc60048036038101906102b79190611cc3565b610f4d565b6040516102cb93929190611d25565b60405180910390f35b6102dc610fb2565b6040516102e99190611c26565b60405180910390f35b61030c60048036038101906103079190611df8565b610fd5565b6040516103199190611c26565b60405180910390f35b61032a611042565b6040516103379190611ec9565b60405180910390f35b61035a60048036038101906103559190611b60565b611067565b005b61037660048036038101906103719190611b8b565b611227565b6040516103839190611b47565b60405180910390f35b6103a660048036038101906103a19190611b60565b611268565b005b6103c260048036038101906103bd9190611f3f565b6112ec565b005b6103de60048036038101906103d99190611fd0565b6113b6565b6040516103eb9190611b47565b60405180910390f35b8160025f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205414610473576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161046a90612055565b60405180910390fd5b5f61048084848888610e58565b90505f61048c826113b6565b90505f6104998285610ee3565b90505f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610528576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161051f906120bd565b60405180910390fd5b8673ffffffffffffffffffffffffffffffffffffffff166340c10f19878a6040518363ffffffff1660e01b81526004016105639291906120db565b5f604051808303815f87803b15801561057a575f80fd5b505af115801561058c573d5f803e3d5ffd5b5050505060025f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8154809291906105dd9061212f565b91905055507f747bd6dbfd6ceb446b50b008eeade0e74f807993dd969546d7efd6008554b1d087878a886040516106179493929190612176565b60405180910390a15050505050505050565b8160025f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054146106a8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161069f90612055565b60405180910390fd5b5f6106b586848787611227565b90505f6106c1826113b6565b90505f6106ce8285610ee3565b90505f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161461075d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610754906120bd565b60405180910390fd5b8573ffffffffffffffffffffffffffffffffffffffff16639dc29fac89896040518363ffffffff1660e01b81526004016107989291906120db565b5f604051808303815f87803b1580156107af575f80fd5b505af11580156107c1573d5f803e3d5ffd5b5050505060025f8973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8154809291906108129061212f565b91905055507fbfa41556980d157c24e8632dbb78958f8759a86b4acdea421f93dc7259fb55db88878960405161084a939291906121b9565b60405180910390a15050505050505050565b6108646113e5565b8060015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507f449a1bd1377b6ad637113368e2e67a7ff6920f8700956c81906a2485fed27909816040516108d39190611ec9565b60405180910390a150565b8160025f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541461095d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161095490612055565b60405180910390fd5b8273ffffffffffffffffffffffffffffffffffffffff168460600160208101906109879190611b60565b73ffffffffffffffffffffffffffffffffffffffff16146109dd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109d490612238565b60405180910390fd5b5f6109e9858585610d13565b90505f6109f5826113b6565b90505f610a028285610ee3565b90505f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610a91576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a88906120bd565b60405180910390fd5b5f875f016020810190610aa49190611b60565b90508073ffffffffffffffffffffffffffffffffffffffff16636ea056a9888a60a001356040518363ffffffff1660e01b8152600401610ae59291906120db565b5f604051808303815f87803b158015610afc575f80fd5b505af1158015610b0e573d5f803e3d5ffd5b505050508073ffffffffffffffffffffffffffffffffffffffff1663095ea7b360015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff168a60a001356040518363ffffffff1660e01b8152600401610b729291906120db565b6020604051808303815f875af1158015610b8e573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610bb2919061228b565b505f60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663414bf3898a6040518263ffffffff1660e01b8152600401610c0e919061248c565b6020604051808303815f875af1158015610c2a573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610c4e91906124ba565b905060025f8973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f815480929190610c9d9061212f565b91905055507fd36cc53ba71bc76a3db3364981f5296dd4ca5eba0e8c89874f2170515bd20d24888a5f016020810190610cd69190611b60565b8b6020016020810190610ce99190611b60565b8c60a0013585604051610d009594939291906124e5565b60405180910390a1505050505050505050565b5f835f016020810190610d269190611b60565b846020016020810190610d399190611b60565b856040016020810190610d4c9190612536565b866060016020810190610d5f9190611b60565b87608001358860a001358960c001358a60e0016020810190610d819190612561565b898b610d8b61146c565b604051602001610da59b9a9998979695949392919061263c565b6040516020818303038152906040528051906020012090509392505050565b610dcc6113e5565b610dd55f611478565b565b610ddf6113e5565b805f806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507f3430ad8dbed7c32bf49006f0d79d2ab70785ea13ebd4ef7d1b87e487ef08928c81604051610e4d9190611c26565b60405180910390a150565b5f84848484610e6561146c565b604051602001610e79959493929190612728565b604051602081830303815290604052805190602001209050949350505050565b6002602052805f5260405f205f915090505481565b5f80610eb8611549565b9050805f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1691505090565b5f805f80610ef085610f4d565b9250925092506001868285856040515f8152602001604052604051610f189493929190612786565b6020604051602081039080840390855afa158015610f38573d5f803e3d5ffd5b50505060206040510351935050505092915050565b5f805f6041845114610f94576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f8b90612813565b60405180910390fd5b602084015192506040840151915060608401515f1a90509193909250565b5f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600382805160208101820180518482526020830160208501208183528095505050505050818051602081018201805184825260208301602085012081835280955050505050505f915091509054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b5f611070611570565b90505f815f0160089054906101000a900460ff161590505f825f015f9054906101000a900467ffffffffffffffff1690505f808267ffffffffffffffff161480156110b85750825b90505f60018367ffffffffffffffff161480156110eb57505f3073ffffffffffffffffffffffffffffffffffffffff163b145b9050811580156110f9575080155b15611130576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001855f015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550831561117d576001855f0160086101000a81548160ff0219169083151502179055505b61118633611597565b855f806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550831561121f575f855f0160086101000a81548160ff0219169083151502179055507fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d26001604051611216919061287d565b60405180910390a15b505050505050565b5f8484848461123461146c565b604051602001611248959493929190612728565b604051602081830303815290604052805190602001209050949350505050565b6112706113e5565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036112e0575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016112d79190611c26565b60405180910390fd5b6112e981611478565b50565b6112f46113e5565b80600386866040516113079291906128c4565b908152602001604051809103902084846040516113259291906128c4565b90815260200160405180910390205f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507f7ab7074f7826dda120af25c3924584f012f6cc3108c1185301c5461ac75f800785858585856040516113a7959493929190612908565b60405180910390a15050505050565b5f816040516020016113c891906129b9565b604051602081830303815290604052805190602001209050919050565b6113ed6115ab565b73ffffffffffffffffffffffffffffffffffffffff1661140b610eae565b73ffffffffffffffffffffffffffffffffffffffff161461146a5761142e6115ab565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016114619190611c26565b60405180910390fd5b565b5f804690508091505090565b5f611481611549565b90505f815f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905082825f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508273ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3505050565b5f7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300905090565b5f7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00905090565b61159f6115b2565b6115a8816115f2565b50565b5f33905090565b6115ba611676565b6115f0576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6115fa6115b2565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361166a575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016116619190611c26565b60405180910390fd5b61167381611478565b50565b5f61167f611570565b5f0160089054906101000a900460ff16905090565b5f604051905090565b5f80fd5b5f80fd5b5f819050919050565b6116b7816116a5565b81146116c1575f80fd5b50565b5f813590506116d2816116ae565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f611701826116d8565b9050919050565b5f611712826116f7565b9050919050565b61172281611708565b811461172c575f80fd5b50565b5f8135905061173d81611719565b92915050565b61174c816116f7565b8114611756575f80fd5b50565b5f8135905061176781611743565b92915050565b5f80fd5b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6117bb82611775565b810181811067ffffffffffffffff821117156117da576117d9611785565b5b80604052505050565b5f6117ec611694565b90506117f882826117b2565b919050565b5f67ffffffffffffffff82111561181757611816611785565b5b61182082611775565b9050602081019050919050565b828183375f83830152505050565b5f61184d611848846117fd565b6117e3565b90508281526020810184848401111561186957611868611771565b5b61187484828561182d565b509392505050565b5f82601f8301126118905761188f61176d565b5b81356118a084826020860161183b565b91505092915050565b5f805f805f60a086880312156118c2576118c161169d565b5b5f6118cf888289016116c4565b95505060206118e08882890161172f565b94505060406118f188828901611759565b9350506060611902888289016116c4565b925050608086013567ffffffffffffffff811115611923576119226116a1565b5b61192f8882890161187c565b9150509295509295909350565b5f805f805f60a086880312156119555761195461169d565b5b5f61196288828901611759565b9550506020611973888289016116c4565b94505060406119848882890161172f565b9350506060611995888289016116c4565b925050608086013567ffffffffffffffff8111156119b6576119b56116a1565b5b6119c28882890161187c565b9150509295509295909350565b5f6119d9826116f7565b9050919050565b6119e9816119cf565b81146119f3575f80fd5b50565b5f81359050611a04816119e0565b92915050565b5f60208284031215611a1f57611a1e61169d565b5b5f611a2c848285016119f6565b91505092915050565b5f80fd5b5f6101008284031215611a4f57611a4e611a35565b5b81905092915050565b5f805f806101608587031215611a7157611a7061169d565b5b5f611a7e87828801611a39565b945050610100611a9087828801611759565b935050610120611aa2878288016116c4565b92505061014085013567ffffffffffffffff811115611ac457611ac36116a1565b5b611ad08782880161187c565b91505092959194509250565b5f805f6101408486031215611af457611af361169d565b5b5f611b0186828701611a39565b935050610100611b1386828701611759565b925050610120611b25868287016116c4565b9150509250925092565b5f819050919050565b611b4181611b2f565b82525050565b5f602082019050611b5a5f830184611b38565b92915050565b5f60208284031215611b7557611b7461169d565b5b5f611b8284828501611759565b91505092915050565b5f805f8060808587031215611ba357611ba261169d565b5b5f611bb087828801611759565b9450506020611bc1878288016116c4565b9350506040611bd2878288016116c4565b9250506060611be38782880161172f565b91505092959194509250565b611bf8816116a5565b82525050565b5f602082019050611c115f830184611bef565b92915050565b611c20816116f7565b82525050565b5f602082019050611c395f830184611c17565b92915050565b611c4881611b2f565b8114611c52575f80fd5b50565b5f81359050611c6381611c3f565b92915050565b5f8060408385031215611c7f57611c7e61169d565b5b5f611c8c85828601611c55565b925050602083013567ffffffffffffffff811115611cad57611cac6116a1565b5b611cb98582860161187c565b9150509250929050565b5f60208284031215611cd857611cd761169d565b5b5f82013567ffffffffffffffff811115611cf557611cf46116a1565b5b611d018482850161187c565b91505092915050565b5f60ff82169050919050565b611d1f81611d0a565b82525050565b5f606082019050611d385f830186611b38565b611d456020830185611b38565b611d526040830184611d16565b949350505050565b5f67ffffffffffffffff821115611d7457611d73611785565b5b611d7d82611775565b9050602081019050919050565b5f611d9c611d9784611d5a565b6117e3565b905082815260208101848484011115611db857611db7611771565b5b611dc384828561182d565b509392505050565b5f82601f830112611ddf57611dde61176d565b5b8135611def848260208601611d8a565b91505092915050565b5f8060408385031215611e0e57611e0d61169d565b5b5f83013567ffffffffffffffff811115611e2b57611e2a6116a1565b5b611e3785828601611dcb565b925050602083013567ffffffffffffffff811115611e5857611e576116a1565b5b611e6485828601611dcb565b9150509250929050565b5f819050919050565b5f611e91611e8c611e87846116d8565b611e6e565b6116d8565b9050919050565b5f611ea282611e77565b9050919050565b5f611eb382611e98565b9050919050565b611ec381611ea9565b82525050565b5f602082019050611edc5f830184611eba565b92915050565b5f80fd5b5f80fd5b5f8083601f840112611eff57611efe61176d565b5b8235905067ffffffffffffffff811115611f1c57611f1b611ee2565b5b602083019150836001820283011115611f3857611f37611ee6565b5b9250929050565b5f805f805f60608688031215611f5857611f5761169d565b5b5f86013567ffffffffffffffff811115611f7557611f746116a1565b5b611f8188828901611eea565b9550955050602086013567ffffffffffffffff811115611fa457611fa36116a1565b5b611fb088828901611eea565b93509350506040611fc388828901611759565b9150509295509295909350565b5f60208284031215611fe557611fe461169d565b5b5f611ff284828501611c55565b91505092915050565b5f82825260208201905092915050565b7f4272696467653a206e6f6e6365000000000000000000000000000000000000005f82015250565b5f61203f600d83611ffb565b915061204a8261200b565b602082019050919050565b5f6020820190508181035f83015261206c81612033565b9050919050565b7f4272696467653a20696e76616c6964207369676e6174757265000000000000005f82015250565b5f6120a7601983611ffb565b91506120b282612073565b602082019050919050565b5f6020820190508181035f8301526120d48161209b565b9050919050565b5f6040820190506120ee5f830185611c17565b6120fb6020830184611bef565b9392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f612139826116a5565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361216b5761216a612102565b5b600182019050919050565b5f6080820190506121895f830187611c17565b6121966020830186611c17565b6121a36040830185611bef565b6121b06060830184611bef565b95945050505050565b5f6060820190506121cc5f830186611c17565b6121d96020830185611c17565b6121e66040830184611bef565b949350505050565b7f4272696467653a2066726f6d00000000000000000000000000000000000000005f82015250565b5f612222600c83611ffb565b915061222d826121ee565b602082019050919050565b5f6020820190508181035f83015261224f81612216565b9050919050565b5f8115159050919050565b61226a81612256565b8114612274575f80fd5b50565b5f8151905061228581612261565b92915050565b5f602082840312156122a05761229f61169d565b5b5f6122ad84828501612277565b91505092915050565b5f6122c46020840184611759565b905092915050565b6122d5816116f7565b82525050565b5f62ffffff82169050919050565b6122f2816122db565b81146122fc575f80fd5b50565b5f8135905061230d816122e9565b92915050565b5f61232160208401846122ff565b905092915050565b612332816122db565b82525050565b5f61234660208401846116c4565b905092915050565b612357816116a5565b82525050565b612366816116d8565b8114612370575f80fd5b50565b5f813590506123818161235d565b92915050565b5f6123956020840184612373565b905092915050565b6123a6816116d8565b82525050565b61010082016123bd5f8301836122b6565b6123c95f8501826122cc565b506123d760208301836122b6565b6123e460208501826122cc565b506123f26040830183612313565b6123ff6040850182612329565b5061240d60608301836122b6565b61241a60608501826122cc565b506124286080830183612338565b612435608085018261234e565b5061244360a0830183612338565b61245060a085018261234e565b5061245e60c0830183612338565b61246b60c085018261234e565b5061247960e0830183612387565b61248660e085018261239d565b50505050565b5f610100820190506124a05f8301846123ac565b92915050565b5f815190506124b4816116ae565b92915050565b5f602082840312156124cf576124ce61169d565b5b5f6124dc848285016124a6565b91505092915050565b5f60a0820190506124f85f830188611c17565b6125056020830187611c17565b6125126040830186611c17565b61251f6060830185611bef565b61252c6080830184611bef565b9695505050505050565b5f6020828403121561254b5761254a61169d565b5b5f612558848285016122ff565b91505092915050565b5f602082840312156125765761257561169d565b5b5f61258384828501612373565b91505092915050565b5f8160601b9050919050565b5f6125a28261258c565b9050919050565b5f6125b382612598565b9050919050565b6125cb6125c6826116f7565b6125a9565b82525050565b5f8160e81b9050919050565b5f6125e7826125d1565b9050919050565b6125ff6125fa826122db565b6125dd565b82525050565b5f819050919050565b61261f61261a826116a5565b612605565b82525050565b612636612631826116d8565b612598565b82525050565b5f612647828e6125ba565b601482019150612657828d6125ba565b601482019150612667828c6125ee565b600382019150612677828b6125ba565b601482019150612687828a61260e565b602082019150612697828961260e565b6020820191506126a7828861260e565b6020820191506126b78287612625565b6014820191506126c7828661260e565b6020820191506126d782856125ba565b6014820191506126e7828461260e565b6020820191508190509c9b505050505050505050505050565b5f61270a82611e98565b9050919050565b61272261271d82612700565b6125a9565b82525050565b5f61273382886125ba565b601482019150612743828761260e565b602082019150612753828661260e565b6020820191506127638285612711565b601482019150612773828461260e565b6020820191508190509695505050505050565b5f6080820190506127995f830187611b38565b6127a66020830186611d16565b6127b36040830185611b38565b6127c06060830184611b38565b95945050505050565b7f696e76616c6964207369676e6174757265206c656e67746800000000000000005f82015250565b5f6127fd601883611ffb565b9150612808826127c9565b602082019050919050565b5f6020820190508181035f83015261282a816127f1565b9050919050565b5f819050919050565b5f67ffffffffffffffff82169050919050565b5f61286761286261285d84612831565b611e6e565b61283a565b9050919050565b6128778161284d565b82525050565b5f6020820190506128905f83018461286e565b92915050565b5f81905092915050565b5f6128ab8385612896565b93506128b883858461182d565b82840190509392505050565b5f6128d08284866128a0565b91508190509392505050565b5f6128e78385611ffb565b93506128f483858461182d565b6128fd83611775565b840190509392505050565b5f6060820190508181035f8301526129218187896128dc565b905081810360208301526129368185876128dc565b90506129456040830184611c17565b9695505050505050565b7f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f82015250565b5f612983601c83612896565b915061298e8261294f565b601c82019050919050565b5f819050919050565b6129b36129ae82611b2f565b612999565b82525050565b5f6129c382612977565b91506129cf82846129a2565b6020820191508190509291505056fea26469706673582212208c0480f7692cb5fde3c4dcbb3f13aa819767ce9916f8457fd7d2fb80aa6b204864736f6c63430008190033",
}

// BridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use BridgeMetaData.ABI instead.
var BridgeABI = BridgeMetaData.ABI

// BridgeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BridgeMetaData.Bin instead.
var BridgeBin = BridgeMetaData.Bin

// DeployBridge deploys a new Ethereum contract, binding an instance of Bridge to it.
func DeployBridge(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Bridge, error) {
	parsed, err := BridgeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BridgeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Bridge{BridgeCaller: BridgeCaller{contract: contract}, BridgeTransactor: BridgeTransactor{contract: contract}, BridgeFilterer: BridgeFilterer{contract: contract}}, nil
}

// Bridge is an auto generated Go binding around an Ethereum contract.
type Bridge struct {
	BridgeCaller     // Read-only binding to the contract
	BridgeTransactor // Write-only binding to the contract
	BridgeFilterer   // Log filterer for contract events
}

// BridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type BridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BridgeSession struct {
	Contract     *Bridge           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BridgeCallerSession struct {
	Contract *BridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// BridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BridgeTransactorSession struct {
	Contract     *BridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type BridgeRaw struct {
	Contract *Bridge // Generic contract binding to access the raw methods on
}

// BridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BridgeCallerRaw struct {
	Contract *BridgeCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BridgeTransactorRaw struct {
	Contract *BridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridge creates a new instance of Bridge, bound to a specific deployed contract.
func NewBridge(address common.Address, backend bind.ContractBackend) (*Bridge, error) {
	contract, err := bindBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bridge{BridgeCaller: BridgeCaller{contract: contract}, BridgeTransactor: BridgeTransactor{contract: contract}, BridgeFilterer: BridgeFilterer{contract: contract}}, nil
}

// NewBridgeCaller creates a new read-only instance of Bridge, bound to a specific deployed contract.
func NewBridgeCaller(address common.Address, caller bind.ContractCaller) (*BridgeCaller, error) {
	contract, err := bindBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeCaller{contract: contract}, nil
}

// NewBridgeTransactor creates a new write-only instance of Bridge, bound to a specific deployed contract.
func NewBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeTransactor, error) {
	contract, err := bindBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeTransactor{contract: contract}, nil
}

// NewBridgeFilterer creates a new log filterer instance of Bridge, bound to a specific deployed contract.
func NewBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeFilterer, error) {
	contract, err := bindBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeFilterer{contract: contract}, nil
}

// bindBridge binds a generic wrapper to an already deployed contract.
func bindBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BridgeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bridge *BridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bridge.Contract.BridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bridge *BridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bridge *BridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bridge *BridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bridge *BridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bridge *BridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bridge.Contract.contract.Transact(opts, method, params...)
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (_Bridge *BridgeCaller) Authority(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "authority")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (_Bridge *BridgeSession) Authority() (common.Address, error) {
	return _Bridge.Contract.Authority(&_Bridge.CallOpts)
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (_Bridge *BridgeCallerSession) Authority() (common.Address, error) {
	return _Bridge.Contract.Authority(&_Bridge.CallOpts)
}

// GetBurnMessageHash is a free data retrieval call binding the contract method 0xdfc34481.
//
// Solidity: function getBurnMessageHash(address account, uint256 nonce, uint256 amount, address token) view returns(bytes32)
func (_Bridge *BridgeCaller) GetBurnMessageHash(opts *bind.CallOpts, account common.Address, nonce *big.Int, amount *big.Int, token common.Address) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getBurnMessageHash", account, nonce, amount, token)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetBurnMessageHash is a free data retrieval call binding the contract method 0xdfc34481.
//
// Solidity: function getBurnMessageHash(address account, uint256 nonce, uint256 amount, address token) view returns(bytes32)
func (_Bridge *BridgeSession) GetBurnMessageHash(account common.Address, nonce *big.Int, amount *big.Int, token common.Address) ([32]byte, error) {
	return _Bridge.Contract.GetBurnMessageHash(&_Bridge.CallOpts, account, nonce, amount, token)
}

// GetBurnMessageHash is a free data retrieval call binding the contract method 0xdfc34481.
//
// Solidity: function getBurnMessageHash(address account, uint256 nonce, uint256 amount, address token) view returns(bytes32)
func (_Bridge *BridgeCallerSession) GetBurnMessageHash(account common.Address, nonce *big.Int, amount *big.Int, token common.Address) ([32]byte, error) {
	return _Bridge.Contract.GetBurnMessageHash(&_Bridge.CallOpts, account, nonce, amount, token)
}

// GetEthSignedMessageHash is a free data retrieval call binding the contract method 0xfa540801.
//
// Solidity: function getEthSignedMessageHash(bytes32 _messageHash) pure returns(bytes32)
func (_Bridge *BridgeCaller) GetEthSignedMessageHash(opts *bind.CallOpts, _messageHash [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getEthSignedMessageHash", _messageHash)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetEthSignedMessageHash is a free data retrieval call binding the contract method 0xfa540801.
//
// Solidity: function getEthSignedMessageHash(bytes32 _messageHash) pure returns(bytes32)
func (_Bridge *BridgeSession) GetEthSignedMessageHash(_messageHash [32]byte) ([32]byte, error) {
	return _Bridge.Contract.GetEthSignedMessageHash(&_Bridge.CallOpts, _messageHash)
}

// GetEthSignedMessageHash is a free data retrieval call binding the contract method 0xfa540801.
//
// Solidity: function getEthSignedMessageHash(bytes32 _messageHash) pure returns(bytes32)
func (_Bridge *BridgeCallerSession) GetEthSignedMessageHash(_messageHash [32]byte) ([32]byte, error) {
	return _Bridge.Contract.GetEthSignedMessageHash(&_Bridge.CallOpts, _messageHash)
}

// GetMintMessageHash is a free data retrieval call binding the contract method 0x7b7dfd10.
//
// Solidity: function getMintMessageHash(address account, uint256 nonce, uint256 amount, address token) view returns(bytes32)
func (_Bridge *BridgeCaller) GetMintMessageHash(opts *bind.CallOpts, account common.Address, nonce *big.Int, amount *big.Int, token common.Address) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getMintMessageHash", account, nonce, amount, token)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetMintMessageHash is a free data retrieval call binding the contract method 0x7b7dfd10.
//
// Solidity: function getMintMessageHash(address account, uint256 nonce, uint256 amount, address token) view returns(bytes32)
func (_Bridge *BridgeSession) GetMintMessageHash(account common.Address, nonce *big.Int, amount *big.Int, token common.Address) ([32]byte, error) {
	return _Bridge.Contract.GetMintMessageHash(&_Bridge.CallOpts, account, nonce, amount, token)
}

// GetMintMessageHash is a free data retrieval call binding the contract method 0x7b7dfd10.
//
// Solidity: function getMintMessageHash(address account, uint256 nonce, uint256 amount, address token) view returns(bytes32)
func (_Bridge *BridgeCallerSession) GetMintMessageHash(account common.Address, nonce *big.Int, amount *big.Int, token common.Address) ([32]byte, error) {
	return _Bridge.Contract.GetMintMessageHash(&_Bridge.CallOpts, account, nonce, amount, token)
}

// GetSwapMessageHash is a free data retrieval call binding the contract method 0x68cab61d.
//
// Solidity: function getSwapMessageHash((address,address,uint24,address,uint256,uint256,uint256,uint160) params, address from, uint256 nonce) view returns(bytes32)
func (_Bridge *BridgeCaller) GetSwapMessageHash(opts *bind.CallOpts, params ISwapRouterExactInputSingleParams, from common.Address, nonce *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getSwapMessageHash", params, from, nonce)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetSwapMessageHash is a free data retrieval call binding the contract method 0x68cab61d.
//
// Solidity: function getSwapMessageHash((address,address,uint24,address,uint256,uint256,uint256,uint160) params, address from, uint256 nonce) view returns(bytes32)
func (_Bridge *BridgeSession) GetSwapMessageHash(params ISwapRouterExactInputSingleParams, from common.Address, nonce *big.Int) ([32]byte, error) {
	return _Bridge.Contract.GetSwapMessageHash(&_Bridge.CallOpts, params, from, nonce)
}

// GetSwapMessageHash is a free data retrieval call binding the contract method 0x68cab61d.
//
// Solidity: function getSwapMessageHash((address,address,uint24,address,uint256,uint256,uint256,uint160) params, address from, uint256 nonce) view returns(bytes32)
func (_Bridge *BridgeCallerSession) GetSwapMessageHash(params ISwapRouterExactInputSingleParams, from common.Address, nonce *big.Int) ([32]byte, error) {
	return _Bridge.Contract.GetSwapMessageHash(&_Bridge.CallOpts, params, from, nonce)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_Bridge *BridgeCaller) Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "nonces", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_Bridge *BridgeSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _Bridge.Contract.Nonces(&_Bridge.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_Bridge *BridgeCallerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _Bridge.Contract.Nonces(&_Bridge.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Bridge *BridgeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Bridge *BridgeSession) Owner() (common.Address, error) {
	return _Bridge.Contract.Owner(&_Bridge.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Bridge *BridgeCallerSession) Owner() (common.Address, error) {
	return _Bridge.Contract.Owner(&_Bridge.CallOpts)
}

// PeggedTokens is a free data retrieval call binding the contract method 0xbfb4ad0c.
//
// Solidity: function peggedTokens(string , string ) view returns(address)
func (_Bridge *BridgeCaller) PeggedTokens(opts *bind.CallOpts, arg0 string, arg1 string) (common.Address, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "peggedTokens", arg0, arg1)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PeggedTokens is a free data retrieval call binding the contract method 0xbfb4ad0c.
//
// Solidity: function peggedTokens(string , string ) view returns(address)
func (_Bridge *BridgeSession) PeggedTokens(arg0 string, arg1 string) (common.Address, error) {
	return _Bridge.Contract.PeggedTokens(&_Bridge.CallOpts, arg0, arg1)
}

// PeggedTokens is a free data retrieval call binding the contract method 0xbfb4ad0c.
//
// Solidity: function peggedTokens(string , string ) view returns(address)
func (_Bridge *BridgeCallerSession) PeggedTokens(arg0 string, arg1 string) (common.Address, error) {
	return _Bridge.Contract.PeggedTokens(&_Bridge.CallOpts, arg0, arg1)
}

// RecoverSigner is a free data retrieval call binding the contract method 0x97aba7f9.
//
// Solidity: function recoverSigner(bytes32 _ethSignedMessageHash, bytes _signature) pure returns(address)
func (_Bridge *BridgeCaller) RecoverSigner(opts *bind.CallOpts, _ethSignedMessageHash [32]byte, _signature []byte) (common.Address, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "recoverSigner", _ethSignedMessageHash, _signature)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RecoverSigner is a free data retrieval call binding the contract method 0x97aba7f9.
//
// Solidity: function recoverSigner(bytes32 _ethSignedMessageHash, bytes _signature) pure returns(address)
func (_Bridge *BridgeSession) RecoverSigner(_ethSignedMessageHash [32]byte, _signature []byte) (common.Address, error) {
	return _Bridge.Contract.RecoverSigner(&_Bridge.CallOpts, _ethSignedMessageHash, _signature)
}

// RecoverSigner is a free data retrieval call binding the contract method 0x97aba7f9.
//
// Solidity: function recoverSigner(bytes32 _ethSignedMessageHash, bytes _signature) pure returns(address)
func (_Bridge *BridgeCallerSession) RecoverSigner(_ethSignedMessageHash [32]byte, _signature []byte) (common.Address, error) {
	return _Bridge.Contract.RecoverSigner(&_Bridge.CallOpts, _ethSignedMessageHash, _signature)
}

// SplitSignature is a free data retrieval call binding the contract method 0xa7bb5803.
//
// Solidity: function splitSignature(bytes sig) pure returns(bytes32 r, bytes32 s, uint8 v)
func (_Bridge *BridgeCaller) SplitSignature(opts *bind.CallOpts, sig []byte) (struct {
	R [32]byte
	S [32]byte
	V uint8
}, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "splitSignature", sig)

	outstruct := new(struct {
		R [32]byte
		S [32]byte
		V uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.R = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.S = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.V = *abi.ConvertType(out[2], new(uint8)).(*uint8)

	return *outstruct, err

}

// SplitSignature is a free data retrieval call binding the contract method 0xa7bb5803.
//
// Solidity: function splitSignature(bytes sig) pure returns(bytes32 r, bytes32 s, uint8 v)
func (_Bridge *BridgeSession) SplitSignature(sig []byte) (struct {
	R [32]byte
	S [32]byte
	V uint8
}, error) {
	return _Bridge.Contract.SplitSignature(&_Bridge.CallOpts, sig)
}

// SplitSignature is a free data retrieval call binding the contract method 0xa7bb5803.
//
// Solidity: function splitSignature(bytes sig) pure returns(bytes32 r, bytes32 s, uint8 v)
func (_Bridge *BridgeCallerSession) SplitSignature(sig []byte) (struct {
	R [32]byte
	S [32]byte
	V uint8
}, error) {
	return _Bridge.Contract.SplitSignature(&_Bridge.CallOpts, sig)
}

// SwapRouter is a free data retrieval call binding the contract method 0xc31c9c07.
//
// Solidity: function swapRouter() view returns(address)
func (_Bridge *BridgeCaller) SwapRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "swapRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SwapRouter is a free data retrieval call binding the contract method 0xc31c9c07.
//
// Solidity: function swapRouter() view returns(address)
func (_Bridge *BridgeSession) SwapRouter() (common.Address, error) {
	return _Bridge.Contract.SwapRouter(&_Bridge.CallOpts)
}

// SwapRouter is a free data retrieval call binding the contract method 0xc31c9c07.
//
// Solidity: function swapRouter() view returns(address)
func (_Bridge *BridgeCallerSession) SwapRouter() (common.Address, error) {
	return _Bridge.Contract.SwapRouter(&_Bridge.CallOpts)
}

// AddToken is a paid mutator transaction binding the contract method 0xfa3d24f2.
//
// Solidity: function addToken(string chainId, string srcToken, address peggedToken) returns()
func (_Bridge *BridgeTransactor) AddToken(opts *bind.TransactOpts, chainId string, srcToken string, peggedToken common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "addToken", chainId, srcToken, peggedToken)
}

// AddToken is a paid mutator transaction binding the contract method 0xfa3d24f2.
//
// Solidity: function addToken(string chainId, string srcToken, address peggedToken) returns()
func (_Bridge *BridgeSession) AddToken(chainId string, srcToken string, peggedToken common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.AddToken(&_Bridge.TransactOpts, chainId, srcToken, peggedToken)
}

// AddToken is a paid mutator transaction binding the contract method 0xfa3d24f2.
//
// Solidity: function addToken(string chainId, string srcToken, address peggedToken) returns()
func (_Bridge *BridgeTransactorSession) AddToken(chainId string, srcToken string, peggedToken common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.AddToken(&_Bridge.TransactOpts, chainId, srcToken, peggedToken)
}

// Burn is a paid mutator transaction binding the contract method 0x35ae976e.
//
// Solidity: function burn(address from, uint256 amount, address token, uint256 nonce, bytes signature) returns()
func (_Bridge *BridgeTransactor) Burn(opts *bind.TransactOpts, from common.Address, amount *big.Int, token common.Address, nonce *big.Int, signature []byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "burn", from, amount, token, nonce, signature)
}

// Burn is a paid mutator transaction binding the contract method 0x35ae976e.
//
// Solidity: function burn(address from, uint256 amount, address token, uint256 nonce, bytes signature) returns()
func (_Bridge *BridgeSession) Burn(from common.Address, amount *big.Int, token common.Address, nonce *big.Int, signature []byte) (*types.Transaction, error) {
	return _Bridge.Contract.Burn(&_Bridge.TransactOpts, from, amount, token, nonce, signature)
}

// Burn is a paid mutator transaction binding the contract method 0x35ae976e.
//
// Solidity: function burn(address from, uint256 amount, address token, uint256 nonce, bytes signature) returns()
func (_Bridge *BridgeTransactorSession) Burn(from common.Address, amount *big.Int, token common.Address, nonce *big.Int, signature []byte) (*types.Transaction, error) {
	return _Bridge.Contract.Burn(&_Bridge.TransactOpts, from, amount, token, nonce, signature)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _authority) returns()
func (_Bridge *BridgeTransactor) Initialize(opts *bind.TransactOpts, _authority common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "initialize", _authority)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _authority) returns()
func (_Bridge *BridgeSession) Initialize(_authority common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.Initialize(&_Bridge.TransactOpts, _authority)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _authority) returns()
func (_Bridge *BridgeTransactorSession) Initialize(_authority common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.Initialize(&_Bridge.TransactOpts, _authority)
}

// Mint is a paid mutator transaction binding the contract method 0x2adfefeb.
//
// Solidity: function mint(uint256 amount, address token, address to, uint256 nonce, bytes signature) returns()
func (_Bridge *BridgeTransactor) Mint(opts *bind.TransactOpts, amount *big.Int, token common.Address, to common.Address, nonce *big.Int, signature []byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "mint", amount, token, to, nonce, signature)
}

// Mint is a paid mutator transaction binding the contract method 0x2adfefeb.
//
// Solidity: function mint(uint256 amount, address token, address to, uint256 nonce, bytes signature) returns()
func (_Bridge *BridgeSession) Mint(amount *big.Int, token common.Address, to common.Address, nonce *big.Int, signature []byte) (*types.Transaction, error) {
	return _Bridge.Contract.Mint(&_Bridge.TransactOpts, amount, token, to, nonce, signature)
}

// Mint is a paid mutator transaction binding the contract method 0x2adfefeb.
//
// Solidity: function mint(uint256 amount, address token, address to, uint256 nonce, bytes signature) returns()
func (_Bridge *BridgeTransactorSession) Mint(amount *big.Int, token common.Address, to common.Address, nonce *big.Int, signature []byte) (*types.Transaction, error) {
	return _Bridge.Contract.Mint(&_Bridge.TransactOpts, amount, token, to, nonce, signature)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Bridge *BridgeTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Bridge *BridgeSession) RenounceOwnership() (*types.Transaction, error) {
	return _Bridge.Contract.RenounceOwnership(&_Bridge.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Bridge *BridgeTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Bridge.Contract.RenounceOwnership(&_Bridge.TransactOpts)
}

// SetAuthority is a paid mutator transaction binding the contract method 0x7a9e5e4b.
//
// Solidity: function setAuthority(address _authority) returns()
func (_Bridge *BridgeTransactor) SetAuthority(opts *bind.TransactOpts, _authority common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "setAuthority", _authority)
}

// SetAuthority is a paid mutator transaction binding the contract method 0x7a9e5e4b.
//
// Solidity: function setAuthority(address _authority) returns()
func (_Bridge *BridgeSession) SetAuthority(_authority common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.SetAuthority(&_Bridge.TransactOpts, _authority)
}

// SetAuthority is a paid mutator transaction binding the contract method 0x7a9e5e4b.
//
// Solidity: function setAuthority(address _authority) returns()
func (_Bridge *BridgeTransactorSession) SetAuthority(_authority common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.SetAuthority(&_Bridge.TransactOpts, _authority)
}

// SetSwapRouter is a paid mutator transaction binding the contract method 0x41273657.
//
// Solidity: function setSwapRouter(address _swapRouter) returns()
func (_Bridge *BridgeTransactor) SetSwapRouter(opts *bind.TransactOpts, _swapRouter common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "setSwapRouter", _swapRouter)
}

// SetSwapRouter is a paid mutator transaction binding the contract method 0x41273657.
//
// Solidity: function setSwapRouter(address _swapRouter) returns()
func (_Bridge *BridgeSession) SetSwapRouter(_swapRouter common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.SetSwapRouter(&_Bridge.TransactOpts, _swapRouter)
}

// SetSwapRouter is a paid mutator transaction binding the contract method 0x41273657.
//
// Solidity: function setSwapRouter(address _swapRouter) returns()
func (_Bridge *BridgeTransactorSession) SetSwapRouter(_swapRouter common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.SetSwapRouter(&_Bridge.TransactOpts, _swapRouter)
}

// Swap is a paid mutator transaction binding the contract method 0x4f334e89.
//
// Solidity: function swap((address,address,uint24,address,uint256,uint256,uint256,uint160) params, address from, uint256 nonce, bytes signature) returns()
func (_Bridge *BridgeTransactor) Swap(opts *bind.TransactOpts, params ISwapRouterExactInputSingleParams, from common.Address, nonce *big.Int, signature []byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "swap", params, from, nonce, signature)
}

// Swap is a paid mutator transaction binding the contract method 0x4f334e89.
//
// Solidity: function swap((address,address,uint24,address,uint256,uint256,uint256,uint160) params, address from, uint256 nonce, bytes signature) returns()
func (_Bridge *BridgeSession) Swap(params ISwapRouterExactInputSingleParams, from common.Address, nonce *big.Int, signature []byte) (*types.Transaction, error) {
	return _Bridge.Contract.Swap(&_Bridge.TransactOpts, params, from, nonce, signature)
}

// Swap is a paid mutator transaction binding the contract method 0x4f334e89.
//
// Solidity: function swap((address,address,uint24,address,uint256,uint256,uint256,uint160) params, address from, uint256 nonce, bytes signature) returns()
func (_Bridge *BridgeTransactorSession) Swap(params ISwapRouterExactInputSingleParams, from common.Address, nonce *big.Int, signature []byte) (*types.Transaction, error) {
	return _Bridge.Contract.Swap(&_Bridge.TransactOpts, params, from, nonce, signature)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Bridge *BridgeTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Bridge *BridgeSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.TransferOwnership(&_Bridge.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Bridge *BridgeTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.TransferOwnership(&_Bridge.TransactOpts, newOwner)
}

// BridgeAuthorityChangedIterator is returned from FilterAuthorityChanged and is used to iterate over the raw logs and unpacked data for AuthorityChanged events raised by the Bridge contract.
type BridgeAuthorityChangedIterator struct {
	Event *BridgeAuthorityChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeAuthorityChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeAuthorityChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeAuthorityChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeAuthorityChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeAuthorityChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeAuthorityChanged represents a AuthorityChanged event raised by the Bridge contract.
type BridgeAuthorityChanged struct {
	Authority common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAuthorityChanged is a free log retrieval operation binding the contract event 0x3430ad8dbed7c32bf49006f0d79d2ab70785ea13ebd4ef7d1b87e487ef08928c.
//
// Solidity: event AuthorityChanged(address authority)
func (_Bridge *BridgeFilterer) FilterAuthorityChanged(opts *bind.FilterOpts) (*BridgeAuthorityChangedIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "AuthorityChanged")
	if err != nil {
		return nil, err
	}
	return &BridgeAuthorityChangedIterator{contract: _Bridge.contract, event: "AuthorityChanged", logs: logs, sub: sub}, nil
}

// WatchAuthorityChanged is a free log subscription operation binding the contract event 0x3430ad8dbed7c32bf49006f0d79d2ab70785ea13ebd4ef7d1b87e487ef08928c.
//
// Solidity: event AuthorityChanged(address authority)
func (_Bridge *BridgeFilterer) WatchAuthorityChanged(opts *bind.WatchOpts, sink chan<- *BridgeAuthorityChanged) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "AuthorityChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeAuthorityChanged)
				if err := _Bridge.contract.UnpackLog(event, "AuthorityChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAuthorityChanged is a log parse operation binding the contract event 0x3430ad8dbed7c32bf49006f0d79d2ab70785ea13ebd4ef7d1b87e487ef08928c.
//
// Solidity: event AuthorityChanged(address authority)
func (_Bridge *BridgeFilterer) ParseAuthorityChanged(log types.Log) (*BridgeAuthorityChanged, error) {
	event := new(BridgeAuthorityChanged)
	if err := _Bridge.contract.UnpackLog(event, "AuthorityChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Bridge contract.
type BridgeInitializedIterator struct {
	Event *BridgeInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeInitialized represents a Initialized event raised by the Bridge contract.
type BridgeInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Bridge *BridgeFilterer) FilterInitialized(opts *bind.FilterOpts) (*BridgeInitializedIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &BridgeInitializedIterator{contract: _Bridge.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Bridge *BridgeFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *BridgeInitialized) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeInitialized)
				if err := _Bridge.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Bridge *BridgeFilterer) ParseInitialized(log types.Log) (*BridgeInitialized, error) {
	event := new(BridgeInitialized)
	if err := _Bridge.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Bridge contract.
type BridgeOwnershipTransferredIterator struct {
	Event *BridgeOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeOwnershipTransferred represents a OwnershipTransferred event raised by the Bridge contract.
type BridgeOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Bridge *BridgeFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BridgeOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BridgeOwnershipTransferredIterator{contract: _Bridge.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Bridge *BridgeFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BridgeOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeOwnershipTransferred)
				if err := _Bridge.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Bridge *BridgeFilterer) ParseOwnershipTransferred(log types.Log) (*BridgeOwnershipTransferred, error) {
	event := new(BridgeOwnershipTransferred)
	if err := _Bridge.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeSwapRouterChangedIterator is returned from FilterSwapRouterChanged and is used to iterate over the raw logs and unpacked data for SwapRouterChanged events raised by the Bridge contract.
type BridgeSwapRouterChangedIterator struct {
	Event *BridgeSwapRouterChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeSwapRouterChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeSwapRouterChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeSwapRouterChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeSwapRouterChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeSwapRouterChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeSwapRouterChanged represents a SwapRouterChanged event raised by the Bridge contract.
type BridgeSwapRouterChanged struct {
	SwapRouter common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSwapRouterChanged is a free log retrieval operation binding the contract event 0x449a1bd1377b6ad637113368e2e67a7ff6920f8700956c81906a2485fed27909.
//
// Solidity: event SwapRouterChanged(address swapRouter)
func (_Bridge *BridgeFilterer) FilterSwapRouterChanged(opts *bind.FilterOpts) (*BridgeSwapRouterChangedIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "SwapRouterChanged")
	if err != nil {
		return nil, err
	}
	return &BridgeSwapRouterChangedIterator{contract: _Bridge.contract, event: "SwapRouterChanged", logs: logs, sub: sub}, nil
}

// WatchSwapRouterChanged is a free log subscription operation binding the contract event 0x449a1bd1377b6ad637113368e2e67a7ff6920f8700956c81906a2485fed27909.
//
// Solidity: event SwapRouterChanged(address swapRouter)
func (_Bridge *BridgeFilterer) WatchSwapRouterChanged(opts *bind.WatchOpts, sink chan<- *BridgeSwapRouterChanged) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "SwapRouterChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeSwapRouterChanged)
				if err := _Bridge.contract.UnpackLog(event, "SwapRouterChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSwapRouterChanged is a log parse operation binding the contract event 0x449a1bd1377b6ad637113368e2e67a7ff6920f8700956c81906a2485fed27909.
//
// Solidity: event SwapRouterChanged(address swapRouter)
func (_Bridge *BridgeFilterer) ParseSwapRouterChanged(log types.Log) (*BridgeSwapRouterChanged, error) {
	event := new(BridgeSwapRouterChanged)
	if err := _Bridge.contract.UnpackLog(event, "SwapRouterChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeTokenAddedIterator is returned from FilterTokenAdded and is used to iterate over the raw logs and unpacked data for TokenAdded events raised by the Bridge contract.
type BridgeTokenAddedIterator struct {
	Event *BridgeTokenAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeTokenAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeTokenAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeTokenAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeTokenAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeTokenAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeTokenAdded represents a TokenAdded event raised by the Bridge contract.
type BridgeTokenAdded struct {
	ChainId     string
	SrcToken    string
	PeggedToken common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTokenAdded is a free log retrieval operation binding the contract event 0x7ab7074f7826dda120af25c3924584f012f6cc3108c1185301c5461ac75f8007.
//
// Solidity: event TokenAdded(string chainId, string srcToken, address peggedToken)
func (_Bridge *BridgeFilterer) FilterTokenAdded(opts *bind.FilterOpts) (*BridgeTokenAddedIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "TokenAdded")
	if err != nil {
		return nil, err
	}
	return &BridgeTokenAddedIterator{contract: _Bridge.contract, event: "TokenAdded", logs: logs, sub: sub}, nil
}

// WatchTokenAdded is a free log subscription operation binding the contract event 0x7ab7074f7826dda120af25c3924584f012f6cc3108c1185301c5461ac75f8007.
//
// Solidity: event TokenAdded(string chainId, string srcToken, address peggedToken)
func (_Bridge *BridgeFilterer) WatchTokenAdded(opts *bind.WatchOpts, sink chan<- *BridgeTokenAdded) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "TokenAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeTokenAdded)
				if err := _Bridge.contract.UnpackLog(event, "TokenAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTokenAdded is a log parse operation binding the contract event 0x7ab7074f7826dda120af25c3924584f012f6cc3108c1185301c5461ac75f8007.
//
// Solidity: event TokenAdded(string chainId, string srcToken, address peggedToken)
func (_Bridge *BridgeFilterer) ParseTokenAdded(log types.Log) (*BridgeTokenAdded, error) {
	event := new(BridgeTokenAdded)
	if err := _Bridge.contract.UnpackLog(event, "TokenAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeTokenBurnedIterator is returned from FilterTokenBurned and is used to iterate over the raw logs and unpacked data for TokenBurned events raised by the Bridge contract.
type BridgeTokenBurnedIterator struct {
	Event *BridgeTokenBurned // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeTokenBurnedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeTokenBurned)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeTokenBurned)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeTokenBurnedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeTokenBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeTokenBurned represents a TokenBurned event raised by the Bridge contract.
type BridgeTokenBurned struct {
	User   common.Address
	Token  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTokenBurned is a free log retrieval operation binding the contract event 0xbfa41556980d157c24e8632dbb78958f8759a86b4acdea421f93dc7259fb55db.
//
// Solidity: event TokenBurned(address user, address token, uint256 amount)
func (_Bridge *BridgeFilterer) FilterTokenBurned(opts *bind.FilterOpts) (*BridgeTokenBurnedIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "TokenBurned")
	if err != nil {
		return nil, err
	}
	return &BridgeTokenBurnedIterator{contract: _Bridge.contract, event: "TokenBurned", logs: logs, sub: sub}, nil
}

// WatchTokenBurned is a free log subscription operation binding the contract event 0xbfa41556980d157c24e8632dbb78958f8759a86b4acdea421f93dc7259fb55db.
//
// Solidity: event TokenBurned(address user, address token, uint256 amount)
func (_Bridge *BridgeFilterer) WatchTokenBurned(opts *bind.WatchOpts, sink chan<- *BridgeTokenBurned) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "TokenBurned")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeTokenBurned)
				if err := _Bridge.contract.UnpackLog(event, "TokenBurned", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTokenBurned is a log parse operation binding the contract event 0xbfa41556980d157c24e8632dbb78958f8759a86b4acdea421f93dc7259fb55db.
//
// Solidity: event TokenBurned(address user, address token, uint256 amount)
func (_Bridge *BridgeFilterer) ParseTokenBurned(log types.Log) (*BridgeTokenBurned, error) {
	event := new(BridgeTokenBurned)
	if err := _Bridge.contract.UnpackLog(event, "TokenBurned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeTokenMintedIterator is returned from FilterTokenMinted and is used to iterate over the raw logs and unpacked data for TokenMinted events raised by the Bridge contract.
type BridgeTokenMintedIterator struct {
	Event *BridgeTokenMinted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeTokenMintedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeTokenMinted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeTokenMinted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeTokenMintedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeTokenMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeTokenMinted represents a TokenMinted event raised by the Bridge contract.
type BridgeTokenMinted struct {
	Token  common.Address
	To     common.Address
	Amount *big.Int
	Nonce  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTokenMinted is a free log retrieval operation binding the contract event 0x747bd6dbfd6ceb446b50b008eeade0e74f807993dd969546d7efd6008554b1d0.
//
// Solidity: event TokenMinted(address token, address to, uint256 amount, uint256 nonce)
func (_Bridge *BridgeFilterer) FilterTokenMinted(opts *bind.FilterOpts) (*BridgeTokenMintedIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "TokenMinted")
	if err != nil {
		return nil, err
	}
	return &BridgeTokenMintedIterator{contract: _Bridge.contract, event: "TokenMinted", logs: logs, sub: sub}, nil
}

// WatchTokenMinted is a free log subscription operation binding the contract event 0x747bd6dbfd6ceb446b50b008eeade0e74f807993dd969546d7efd6008554b1d0.
//
// Solidity: event TokenMinted(address token, address to, uint256 amount, uint256 nonce)
func (_Bridge *BridgeFilterer) WatchTokenMinted(opts *bind.WatchOpts, sink chan<- *BridgeTokenMinted) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "TokenMinted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeTokenMinted)
				if err := _Bridge.contract.UnpackLog(event, "TokenMinted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTokenMinted is a log parse operation binding the contract event 0x747bd6dbfd6ceb446b50b008eeade0e74f807993dd969546d7efd6008554b1d0.
//
// Solidity: event TokenMinted(address token, address to, uint256 amount, uint256 nonce)
func (_Bridge *BridgeFilterer) ParseTokenMinted(log types.Log) (*BridgeTokenMinted, error) {
	event := new(BridgeTokenMinted)
	if err := _Bridge.contract.UnpackLog(event, "TokenMinted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeTokenSwappedIterator is returned from FilterTokenSwapped and is used to iterate over the raw logs and unpacked data for TokenSwapped events raised by the Bridge contract.
type BridgeTokenSwappedIterator struct {
	Event *BridgeTokenSwapped // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeTokenSwappedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeTokenSwapped)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeTokenSwapped)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeTokenSwappedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeTokenSwappedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeTokenSwapped represents a TokenSwapped event raised by the Bridge contract.
type BridgeTokenSwapped struct {
	User      common.Address
	TokenIn   common.Address
	TokenOut  common.Address
	AmountIn  *big.Int
	AmountOut *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTokenSwapped is a free log retrieval operation binding the contract event 0xd36cc53ba71bc76a3db3364981f5296dd4ca5eba0e8c89874f2170515bd20d24.
//
// Solidity: event TokenSwapped(address user, address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOut)
func (_Bridge *BridgeFilterer) FilterTokenSwapped(opts *bind.FilterOpts) (*BridgeTokenSwappedIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "TokenSwapped")
	if err != nil {
		return nil, err
	}
	return &BridgeTokenSwappedIterator{contract: _Bridge.contract, event: "TokenSwapped", logs: logs, sub: sub}, nil
}

// WatchTokenSwapped is a free log subscription operation binding the contract event 0xd36cc53ba71bc76a3db3364981f5296dd4ca5eba0e8c89874f2170515bd20d24.
//
// Solidity: event TokenSwapped(address user, address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOut)
func (_Bridge *BridgeFilterer) WatchTokenSwapped(opts *bind.WatchOpts, sink chan<- *BridgeTokenSwapped) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "TokenSwapped")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeTokenSwapped)
				if err := _Bridge.contract.UnpackLog(event, "TokenSwapped", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTokenSwapped is a log parse operation binding the contract event 0xd36cc53ba71bc76a3db3364981f5296dd4ca5eba0e8c89874f2170515bd20d24.
//
// Solidity: event TokenSwapped(address user, address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOut)
func (_Bridge *BridgeFilterer) ParseTokenSwapped(log types.Log) (*BridgeTokenSwapped, error) {
	event := new(BridgeTokenSwapped)
	if err := _Bridge.contract.UnpackLog(event, "TokenSwapped", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
