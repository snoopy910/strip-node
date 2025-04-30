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

// ISwapRouterExactInputParams is an auto generated low-level Go binding around an user-defined struct.
type ISwapRouterExactInputParams struct {
	Path             []byte
	Recipient        common.Address
	Deadline         *big.Int
	AmountIn         *big.Int
	AmountOutMinimum *big.Int
}

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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"AuthorityChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractISwapRouter\",\"name\":\"swapRouter\",\"type\":\"address\"}],\"name\":\"SwapRouterChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"chainId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"srcToken\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"peggedToken\",\"type\":\"address\"}],\"name\":\"TokenAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TokenBurned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"TokenMinted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"TokenSwapped\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chainId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"srcToken\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"peggedToken\",\"type\":\"address\"}],\"name\":\"addToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authority\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"contractBridgeToken\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"contractBridgeToken\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getBurnMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_messageHash\",\"type\":\"bytes32\"}],\"name\":\"getEthSignedMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"contractBridgeToken\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getMintMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"structISwapRouter.ExactInputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"getSwapMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"}],\"internalType\":\"structISwapRouter.ExactInputParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"getSwapMultiplePoolsMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"}],\"name\":\"getTokenOutFromPath\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_authority\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"contractBridgeToken\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"peggedTokens\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_ethSignedMessageHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_signature\",\"type\":\"bytes\"}],\"name\":\"recoverSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_authority\",\"type\":\"address\"}],\"name\":\"setAuthority\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractISwapRouter\",\"name\":\"_swapRouter\",\"type\":\"address\"}],\"name\":\"setSwapRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"splitSignature\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"structISwapRouter.ExactInputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"swap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"}],\"internalType\":\"structISwapRouter.ExactInputParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"swapMultiplePools\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"swapRouter\",\"outputs\":[{\"internalType\":\"contractISwapRouter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b50611e158061001f6000396000f3fe608060405234801561001057600080fd5b506004361061014d5760003560e01c80638da5cb5b116100c3578063c4d66de81161007c578063c4d66de81461033d578063d627762c14610350578063dfc344811461021f578063f2fde38b14610363578063fa3d24f214610376578063fa5408011461038957600080fd5b80638da5cb5b1461025257806397aba7f914610282578063a7bb580314610295578063bf7e214f146102c6578063bfb4ad0c146102d9578063c31c9c071461032a57600080fd5b806368cab61d1161011557806368cab61d146101b357806370166f3d146101d9578063715018a6146102045780637a9e5e4b1461020c5780637b7dfd101461021f5780637ecebe001461023257600080fd5b80632adfefeb146101525780632b7aa59e1461016757806335ae976e1461017a578063412736571461018d5780634f334e89146101a0575b600080fd5b610165610160366004611523565b61039c565b005b6101656101753660046115b0565b61051d565b610165610188366004611639565b610897565b61016561019b366004611673565b6109fa565b6101656101ae3660046116aa565b610a57565b6101c66101c1366004611718565b610db3565b6040519081526020015b60405180910390f35b6101ec6101e736600461175b565b610e4e565b6040516001600160a01b0390911681526020016101d0565b610165610ec0565b61016561021a366004611673565b610ed4565b6101c661022d366004611797565b610f2a565b6101c6610240366004611673565b60026020526000908152604090205481565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03166101ec565b6101ec6102903660046117e1565b610f8b565b6102a86102a336600461175b565b61100b565b60408051938452602084019290925260ff16908201526060016101d0565b6000546101ec906001600160a01b031681565b6101ec6102e7366004611827565b815160208184018101805160038252928201948201949094209190935281518083018401805192815290840192909301919091209152546001600160a01b031681565b6001546101ec906001600160a01b031681565b61016561034b366004611673565b61107f565b6101c661035e366004611880565b6111a8565b610165610371366004611673565b6111f0565b610165610384366004611931565b61122e565b6101c66103973660046119b4565b6112db565b6001600160a01b03831660009081526002602052604090205482146103dc5760405162461bcd60e51b81526004016103d3906119cd565b60405180910390fd5b60006103ea84848888610f2a565b905060006103f7826112db565b905060006104058285610f8b565b6000549091506001600160a01b038083169116146104355760405162461bcd60e51b81526004016103d3906119f4565b6040516340c10f1960e01b81526001600160a01b038781166004830152602482018a90528816906340c10f1990604401600060405180830381600087803b15801561047f57600080fd5b505af1158015610493573d6000803e3d6000fd5b505050506001600160a01b03861660009081526002602052604081208054916104bb83611a41565b9091555050604080516001600160a01b03808a16825288166020820152908101899052606081018690527f747bd6dbfd6ceb446b50b008eeade0e74f807993dd969546d7efd6008554b1d0906080015b60405180910390a15050505050505050565b6001600160a01b03831660009081526002602052604090205482146105545760405162461bcd60e51b81526004016103d3906119cd565b6001600160a01b03831661056e6040870160208801611673565b6001600160a01b0316146105b35760405162461bcd60e51b815260206004820152600c60248201526b4272696467653a2066726f6d60a01b60448201526064016103d3565b60006105c1868686866111a8565b905060006105ce826112db565b905060006105dc8285610f8b565b6000549091506001600160a01b0380831691161461060c5760405162461bcd60e51b81526004016103d3906119f4565b604051636ea056a960e01b81526001600160a01b03878116600483015260608a01356024830152889190821690636ea056a990604401600060405180830381600087803b15801561065c57600080fd5b505af1158015610670573d6000803e3d6000fd5b505060015460405163095ea7b360e01b81526001600160a01b03918216600482015260608d01356024820152908416925063095ea7b391506044016020604051808303816000875af11580156106ca573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106ee9190611a5a565b5060015460405163c04b8d5960e01b81526000916001600160a01b03169063c04b8d5990610720908d90600401611aa5565b6020604051808303816000875af115801561073f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107639190611b46565b905089608001358110156107b95760405162461bcd60e51b815260206004820152601a60248201527f4272696467653a20616d6f756e74206f757420746f6f206c6f7700000000000060448201526064016103d3565b60006108026107c88c80611b5f565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610e4e92505050565b6001600160a01b038a16600090815260026020526040812080549293509061082983611a41565b9091555050604080516001600160a01b038b811682528c8116602083015283168183015260608d810135908201526080810184905290517fd36cc53ba71bc76a3db3364981f5296dd4ca5eba0e8c89874f2170515bd20d249181900360a00190a15050505050505050505050565b6001600160a01b03851660009081526002602052604090205482146108ce5760405162461bcd60e51b81526004016103d3906119cd565b60006108dc86848787610f2a565b905060006108e9826112db565b905060006108f78285610f8b565b6000549091506001600160a01b038083169116146109275760405162461bcd60e51b81526004016103d3906119f4565b604051632770a7eb60e21b81526001600160a01b03898116600483015260248201899052871690639dc29fac90604401600060405180830381600087803b15801561097157600080fd5b505af1158015610985573d6000803e3d6000fd5b505050506001600160a01b03881660009081526002602052604081208054916109ad83611a41565b9091555050604080516001600160a01b03808b168252881660208201529081018890527fbfa41556980d157c24e8632dbb78958f8759a86b4acdea421f93dc7259fb55db9060600161050b565b610a0261132e565b600180546001600160a01b0319166001600160a01b0383169081179091556040519081527f449a1bd1377b6ad637113368e2e67a7ff6920f8700956c81906a2485fed27909906020015b60405180910390a150565b6001600160a01b0383166000908152600260205260409020548214610a8e5760405162461bcd60e51b81526004016103d3906119cd565b6001600160a01b038316610aa86080860160608701611673565b6001600160a01b031614610aed5760405162461bcd60e51b815260206004820152600c60248201526b4272696467653a2066726f6d60a01b60448201526064016103d3565b6000610afa858585610db3565b90506000610b07826112db565b90506000610b158285610f8b565b6000549091506001600160a01b03808316911614610b455760405162461bcd60e51b81526004016103d3906119f4565b6000610b546020890189611673565b604051636ea056a960e01b81526001600160a01b03898116600483015260a08b0135602483015291925090821690636ea056a990604401600060405180830381600087803b158015610ba557600080fd5b505af1158015610bb9573d6000803e3d6000fd5b505060015460405163095ea7b360e01b81526001600160a01b03918216600482015260a08c01356024820152908416925063095ea7b391506044016020604051808303816000875af1158015610c13573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c379190611a5a565b5060015460405163414bf38960e01b81526000916001600160a01b03169063414bf38990610c69908c90600401611bb8565b6020604051808303816000875af1158015610c88573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610cac9190611b46565b90508860c00135811015610d025760405162461bcd60e51b815260206004820152601a60248201527f4272696467653a20616d6f756e74206f757420746f6f206c6f7700000000000060448201526064016103d3565b6001600160a01b0388166000908152600260205260408120805491610d2683611a41565b909155507fd36cc53ba71bc76a3db3364981f5296dd4ca5eba0e8c89874f2170515bd20d24905088610d5b60208c018c611673565b610d6b60408d0160208e01611673565b604080516001600160a01b039485168152928416602084015292169181019190915260a0808c01356060830152608082018490520160405180910390a1505050505050505050565b6000610dc26020850185611673565b610dd26040860160208701611673565b610de26060870160408801611c5c565b610df26080880160608901611673565b608088013560a089013560c08a0135610e126101008c0160e08d01611673565b898b46604051602001610e2f9b9a99989796959493929190611c77565b6040516020818303038152906040528051906020012090509392505050565b6000601482511015610e9b5760405162461bcd60e51b8152602060048201526016602482015275109c9a5919d94e881c185d1a081d1bdbc81cda1bdc9d60521b60448201526064016103d3565b600060148351610eab9190611d13565b9290920160200151600160601b900492915050565b610ec861132e565b610ed26000611389565b565b610edc61132e565b600080546001600160a01b0319166001600160a01b0383169081179091556040519081527f3430ad8dbed7c32bf49006f0d79d2ab70785ea13ebd4ef7d1b87e487ef08928c90602001610a4c565b6040516bffffffffffffffffffffffff19606086811b82166020840152603483018690526054830185905283901b16607482015246608882015260009060a8015b604051602081830303815290604052805190602001209050949350505050565b600080600080610f9a8561100b565b6040805160008152602081018083528b905260ff8316918101919091526060810184905260808101839052929550909350915060019060a0016020604051602081039080840390855afa158015610ff5573d6000803e3d6000fd5b5050506020604051035193505050505b92915050565b600080600083516041146110615760405162461bcd60e51b815260206004820152601860248201527f696e76616c6964207369676e6174757265206c656e677468000000000000000060448201526064016103d3565b50505060208101516040820151606090920151909260009190911a90565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a008054600160401b810460ff1615906001600160401b03166000811580156110c45750825b90506000826001600160401b031660011480156110e05750303b155b9050811580156110ee575080155b1561110c5760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561113657845460ff60401b1916600160401b1785555b61113f336113fa565b600080546001600160a01b0319166001600160a01b03881617905583156111a057845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b505050505050565b6000836111b58680611b5f565b6111c56040890160208a01611673565b604089013560608a013560808b0135888a46604051602001610f6b9a99989796959493929190611d26565b6111f861132e565b6001600160a01b03811661122257604051631e4fbdf760e01b8152600060048201526024016103d3565b61122b81611389565b50565b61123661132e565b8060038686604051611249929190611d8d565b90815260200160405180910390208484604051611267929190611d8d565b90815260405190819003602001812080546001600160a01b03939093166001600160a01b0319909316929092179091557f7ab7074f7826dda120af25c3924584f012f6cc3108c1185301c5461ac75f8007906112cc9087908790879087908790611d9d565b60405180910390a15050505050565b6040517f19457468657265756d205369676e6564204d6573736167653a0a3332000000006020820152603c8101829052600090605c01604051602081830303815290604052805190602001209050919050565b336113607f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b031614610ed25760405163118cdaa760e01b81523360048201526024016103d3565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3505050565b61140261140b565b61122b81611454565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0054600160401b900460ff16610ed257604051631afcd79f60e31b815260040160405180910390fd5b6111f861140b565b6001600160a01b038116811461122b57600080fd5b803561147c8161145c565b919050565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126114a857600080fd5b81356001600160401b03808211156114c2576114c2611481565b604051601f8301601f19908116603f011681019082821181831017156114ea576114ea611481565b8160405283815286602085880101111561150357600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080600080600060a0868803121561153b57600080fd5b85359450602086013561154d8161145c565b9350604086013561155d8161145c565b92506060860135915060808601356001600160401b0381111561157f57600080fd5b61158b88828901611497565b9150509295509295909350565b600060a082840312156115aa57600080fd5b50919050565b600080600080600060a086880312156115c857600080fd5b85356001600160401b03808211156115df57600080fd5b6115eb89838a01611598565b9650602088013591506115fd8261145c565b90945060408701359061160f8261145c565b909350606087013592506080870135908082111561162c57600080fd5b5061158b88828901611497565b600080600080600060a0868803121561165157600080fd5b853561165c8161145c565b945060208601359350604086013561155d8161145c565b60006020828403121561168557600080fd5b81356116908161145c565b9392505050565b600061010082840312156115aa57600080fd5b60008060008061016085870312156116c157600080fd5b6116cb8686611697565b93506101008501356116dc8161145c565b925061012085013591506101408501356001600160401b0381111561170057600080fd5b61170c87828801611497565b91505092959194509250565b6000806000610140848603121561172e57600080fd5b6117388585611697565b92506101008401356117498161145c565b92959294505050610120919091013590565b60006020828403121561176d57600080fd5b81356001600160401b0381111561178357600080fd5b61178f84828501611497565b949350505050565b600080600080608085870312156117ad57600080fd5b84356117b88161145c565b9350602085013592506040850135915060608501356117d68161145c565b939692955090935050565b600080604083850312156117f457600080fd5b8235915060208301356001600160401b0381111561181157600080fd5b61181d85828601611497565b9150509250929050565b6000806040838503121561183a57600080fd5b82356001600160401b038082111561185157600080fd5b61185d86838701611497565b9350602085013591508082111561187357600080fd5b5061181d85828601611497565b6000806000806080858703121561189657600080fd5b84356001600160401b038111156118ac57600080fd5b6118b887828801611598565b94505060208501356118c98161145c565b925060408501356118d98161145c565b9396929550929360600135925050565b60008083601f8401126118fb57600080fd5b5081356001600160401b0381111561191257600080fd5b60208301915083602082850101111561192a57600080fd5b9250929050565b60008060008060006060868803121561194957600080fd5b85356001600160401b038082111561196057600080fd5b61196c89838a016118e9565b9097509550602088013591508082111561198557600080fd5b50611992888289016118e9565b90945092505060408601356119a68161145c565b809150509295509295909350565b6000602082840312156119c657600080fd5b5035919050565b6020808252600d908201526c4272696467653a206e6f6e636560981b604082015260600190565b60208082526019908201527f4272696467653a20696e76616c6964207369676e617475726500000000000000604082015260600190565b634e487b7160e01b600052601160045260246000fd5b600060018201611a5357611a53611a2b565b5060010190565b600060208284031215611a6c57600080fd5b8151801515811461169057600080fd5b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6020815260008235601e19843603018112611abf57600080fd5b83016020810190356001600160401b03811115611adb57600080fd5b803603821315611aea57600080fd5b60a06020850152611aff60c085018284611a7c565b915050611b0e60208501611471565b6001600160a01b0381166040850152506040840135606084015260608401356080840152608084013560a08401528091505092915050565b600060208284031215611b5857600080fd5b5051919050565b6000808335601e19843603018112611b7657600080fd5b8301803591506001600160401b03821115611b9057600080fd5b60200191503681900382131561192a57600080fd5b803562ffffff8116811461147c57600080fd5b61010081018235611bc88161145c565b6001600160a01b039081168352602084013590611be48261145c565b16602083015262ffffff611bfa60408501611ba5565b166040830152611c0c60608401611471565b6001600160a01b0381166060840152506080830135608083015260a083013560a083015260c083013560c0830152611c4660e08401611471565b6001600160a01b03811660e08401525092915050565b600060208284031215611c6e57600080fd5b61169082611ba5565b60006bffffffffffffffffffffffff19808e60601b168352808d60601b16601484015262ffffff60e81b8c60e81b166028840152808b60601b16602b84015289603f84015288605f84015287607f840152808760601b16609f840152508460b3830152611cf860d383018560601b6bffffffffffffffffffffffff19169052565b5060e7810191909152610107019a9950505050505050505050565b8181038181111561100557611005611a2b565b60006bffffffffffffffffffffffff19808d60601b1683528a8c60148501376060998a1b81169a909201601481019a909a5250602889019690965260488801949094526068870192909252608886015290921b1660a883015260bc82015260dc0192915050565b8183823760009101908152919050565b606081526000611db1606083018789611a7c565b8281036020840152611dc4818688611a7c565b91505060018060a01b0383166040830152969550505050505056fea2646970667358221220a3727284756541f5da4407b9a5cdb62c122746d9a744b36bf90f3388720cdcc664736f6c63430008190033",
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

// GetSwapMultiplePoolsMessageHash is a free data retrieval call binding the contract method 0xd627762c.
//
// Solidity: function getSwapMultiplePoolsMessageHash((bytes,address,uint256,uint256,uint256) params, address tokenIn, address from, uint256 nonce) view returns(bytes32)
func (_Bridge *BridgeCaller) GetSwapMultiplePoolsMessageHash(opts *bind.CallOpts, params ISwapRouterExactInputParams, tokenIn common.Address, from common.Address, nonce *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getSwapMultiplePoolsMessageHash", params, tokenIn, from, nonce)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetSwapMultiplePoolsMessageHash is a free data retrieval call binding the contract method 0xd627762c.
//
// Solidity: function getSwapMultiplePoolsMessageHash((bytes,address,uint256,uint256,uint256) params, address tokenIn, address from, uint256 nonce) view returns(bytes32)
func (_Bridge *BridgeSession) GetSwapMultiplePoolsMessageHash(params ISwapRouterExactInputParams, tokenIn common.Address, from common.Address, nonce *big.Int) ([32]byte, error) {
	return _Bridge.Contract.GetSwapMultiplePoolsMessageHash(&_Bridge.CallOpts, params, tokenIn, from, nonce)
}

// GetSwapMultiplePoolsMessageHash is a free data retrieval call binding the contract method 0xd627762c.
//
// Solidity: function getSwapMultiplePoolsMessageHash((bytes,address,uint256,uint256,uint256) params, address tokenIn, address from, uint256 nonce) view returns(bytes32)
func (_Bridge *BridgeCallerSession) GetSwapMultiplePoolsMessageHash(params ISwapRouterExactInputParams, tokenIn common.Address, from common.Address, nonce *big.Int) ([32]byte, error) {
	return _Bridge.Contract.GetSwapMultiplePoolsMessageHash(&_Bridge.CallOpts, params, tokenIn, from, nonce)
}

// GetTokenOutFromPath is a free data retrieval call binding the contract method 0x70166f3d.
//
// Solidity: function getTokenOutFromPath(bytes path) pure returns(address tokenOut)
func (_Bridge *BridgeCaller) GetTokenOutFromPath(opts *bind.CallOpts, path []byte) (common.Address, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getTokenOutFromPath", path)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetTokenOutFromPath is a free data retrieval call binding the contract method 0x70166f3d.
//
// Solidity: function getTokenOutFromPath(bytes path) pure returns(address tokenOut)
func (_Bridge *BridgeSession) GetTokenOutFromPath(path []byte) (common.Address, error) {
	return _Bridge.Contract.GetTokenOutFromPath(&_Bridge.CallOpts, path)
}

// GetTokenOutFromPath is a free data retrieval call binding the contract method 0x70166f3d.
//
// Solidity: function getTokenOutFromPath(bytes path) pure returns(address tokenOut)
func (_Bridge *BridgeCallerSession) GetTokenOutFromPath(path []byte) (common.Address, error) {
	return _Bridge.Contract.GetTokenOutFromPath(&_Bridge.CallOpts, path)
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

// SwapMultiplePools is a paid mutator transaction binding the contract method 0x2b7aa59e.
//
// Solidity: function swapMultiplePools((bytes,address,uint256,uint256,uint256) params, address tokenIn, address from, uint256 nonce, bytes signature) returns()
func (_Bridge *BridgeTransactor) SwapMultiplePools(opts *bind.TransactOpts, params ISwapRouterExactInputParams, tokenIn common.Address, from common.Address, nonce *big.Int, signature []byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "swapMultiplePools", params, tokenIn, from, nonce, signature)
}

// SwapMultiplePools is a paid mutator transaction binding the contract method 0x2b7aa59e.
//
// Solidity: function swapMultiplePools((bytes,address,uint256,uint256,uint256) params, address tokenIn, address from, uint256 nonce, bytes signature) returns()
func (_Bridge *BridgeSession) SwapMultiplePools(params ISwapRouterExactInputParams, tokenIn common.Address, from common.Address, nonce *big.Int, signature []byte) (*types.Transaction, error) {
	return _Bridge.Contract.SwapMultiplePools(&_Bridge.TransactOpts, params, tokenIn, from, nonce, signature)
}

// SwapMultiplePools is a paid mutator transaction binding the contract method 0x2b7aa59e.
//
// Solidity: function swapMultiplePools((bytes,address,uint256,uint256,uint256) params, address tokenIn, address from, uint256 nonce, bytes signature) returns()
func (_Bridge *BridgeTransactorSession) SwapMultiplePools(params ISwapRouterExactInputParams, tokenIn common.Address, from common.Address, nonce *big.Int, signature []byte) (*types.Transaction, error) {
	return _Bridge.Contract.SwapMultiplePools(&_Bridge.TransactOpts, params, tokenIn, from, nonce, signature)
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
