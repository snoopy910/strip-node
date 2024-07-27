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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"AuthorityChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractISwapRouter\",\"name\":\"swapRouter\",\"type\":\"address\"}],\"name\":\"SwapRouterChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"chainId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"srcToken\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"peggedToken\",\"type\":\"address\"}],\"name\":\"TokenAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"TokenMinted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"TokenSwapped\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chainId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"srcToken\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"peggedToken\",\"type\":\"address\"}],\"name\":\"addToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authority\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_messageHash\",\"type\":\"bytes32\"}],\"name\":\"getEthSignedMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"contractBridgeToken\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getMintMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"structISwapRouter.ExactInputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"getSwapMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_authority\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"contractBridgeToken\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"peggedTokens\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_ethSignedMessageHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_signature\",\"type\":\"bytes\"}],\"name\":\"recoverSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_authority\",\"type\":\"address\"}],\"name\":\"setAuthority\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractISwapRouter\",\"name\":\"_swapRouter\",\"type\":\"address\"}],\"name\":\"setSwapRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"splitSignature\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"structISwapRouter.ExactInputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"swap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"swapRouter\",\"outputs\":[{\"internalType\":\"contractISwapRouter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f80fd5b506126768061001c5f395ff3fe608060405234801561000f575f80fd5b5060043610610114575f3560e01c806397aba7f9116100a0578063c31c9c071161006f578063c31c9c07146102f0578063c4d66de81461030e578063f2fde38b1461032a578063fa3d24f214610346578063fa5408011461036257610114565b806397aba7f914610240578063a7bb580314610270578063bf7e214f146102a2578063bfb4ad0c146102c057610114565b8063715018a6116100e7578063715018a61461019c5780637a9e5e4b146101a65780637b7dfd10146101c25780637ecebe00146101f25780638da5cb5b1461022257610114565b80632adfefeb1461011857806341273657146101345780634f334e891461015057806368cab61d1461016c575b5f80fd5b610132600480360381019061012d91906115d3565b610392565b005b61014e600480360381019061014991906116a1565b6105c7565b005b61016a600480360381019061016591906116ef565b610649565b005b61018660048036038101906101819190611773565b610a7e565b60405161019391906117de565b60405180910390f35b6101a4610b2f565b005b6101c060048036038101906101bb91906117f7565b610b42565b005b6101dc60048036038101906101d79190611822565b610bc3565b6040516101e991906117de565b60405180910390f35b61020c600480360381019061020791906117f7565b610c04565b6040516102199190611895565b60405180910390f35b61022a610c19565b60405161023791906118bd565b60405180910390f35b61025a60048036038101906102559190611900565b610c4e565b60405161026791906118bd565b60405180910390f35b61028a6004803603810190610285919061195a565b610cb8565b604051610299939291906119bc565b60405180910390f35b6102aa610d1d565b6040516102b791906118bd565b60405180910390f35b6102da60048036038101906102d59190611a8f565b610d40565b6040516102e791906118bd565b60405180910390f35b6102f8610dad565b6040516103059190611b60565b60405180910390f35b610328600480360381019061032391906117f7565b610dd2565b005b610344600480360381019061033f91906117f7565b610f92565b005b610360600480360381019061035b9190611bd6565b611016565b005b61037c60048036038101906103779190611c67565b6110e0565b60405161038991906117de565b60405180910390f35b8160025f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205414610411576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161040890611cec565b60405180910390fd5b5f61041e84848888610bc3565b90505f61042a826110e0565b90505f6104378285610c4e565b90505f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146104c6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104bd90611d54565b60405180910390fd5b8673ffffffffffffffffffffffffffffffffffffffff166340c10f19878a6040518363ffffffff1660e01b8152600401610501929190611d72565b5f604051808303815f87803b158015610518575f80fd5b505af115801561052a573d5f803e3d5ffd5b5050505060025f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f81548092919061057b90611dc6565b91905055507f747bd6dbfd6ceb446b50b008eeade0e74f807993dd969546d7efd6008554b1d087878a886040516105b59493929190611e0d565b60405180910390a15050505050505050565b6105cf61110f565b8060015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507f449a1bd1377b6ad637113368e2e67a7ff6920f8700956c81906a2485fed279098160405161063e9190611b60565b60405180910390a150565b8160025f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054146106c8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106bf90611cec565b60405180910390fd5b8273ffffffffffffffffffffffffffffffffffffffff168460600160208101906106f291906117f7565b73ffffffffffffffffffffffffffffffffffffffff1614610748576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161073f90611e9a565b60405180910390fd5b5f610754858585610a7e565b90505f610760826110e0565b90505f61076d8285610c4e565b90505f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146107fc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107f390611d54565b60405180910390fd5b5f875f01602081019061080f91906117f7565b90508073ffffffffffffffffffffffffffffffffffffffff16636ea056a9888a60a001356040518363ffffffff1660e01b8152600401610850929190611d72565b5f604051808303815f87803b158015610867575f80fd5b505af1158015610879573d5f803e3d5ffd5b505050508073ffffffffffffffffffffffffffffffffffffffff1663095ea7b360015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff168a60a001356040518363ffffffff1660e01b81526004016108dd929190611d72565b6020604051808303815f875af11580156108f9573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061091d9190611eed565b505f60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663414bf3898a6040518263ffffffff1660e01b815260040161097991906120ee565b6020604051808303815f875af1158015610995573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906109b9919061211c565b905060025f8973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f815480929190610a0890611dc6565b91905055507fd36cc53ba71bc76a3db3364981f5296dd4ca5eba0e8c89874f2170515bd20d24888a5f016020810190610a4191906117f7565b8b6020016020810190610a5491906117f7565b8c60a0013585604051610a6b959493929190612147565b60405180910390a1505050505050505050565b5f835f016020810190610a9191906117f7565b846020016020810190610aa491906117f7565b856040016020810190610ab79190612198565b866060016020810190610aca91906117f7565b87608001358860a001358960c001358a60e0016020810190610aec91906121c3565b898b610af6611196565b604051602001610b109b9a9998979695949392919061229e565b6040516020818303038152906040528051906020012090509392505050565b610b3761110f565b610b405f6111a2565b565b610b4a61110f565b805f806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507f3430ad8dbed7c32bf49006f0d79d2ab70785ea13ebd4ef7d1b87e487ef08928c81604051610bb891906118bd565b60405180910390a150565b5f84848484610bd0611196565b604051602001610be495949392919061238a565b604051602081830303815290604052805190602001209050949350505050565b6002602052805f5260405f205f915090505481565b5f80610c23611273565b9050805f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1691505090565b5f805f80610c5b85610cb8565b9250925092506001868285856040515f8152602001604052604051610c8394939291906123e8565b6020604051602081039080840390855afa158015610ca3573d5f803e3d5ffd5b50505060206040510351935050505092915050565b5f805f6041845114610cff576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610cf690612475565b60405180910390fd5b602084015192506040840151915060608401515f1a90509193909250565b5f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600382805160208101820180518482526020830160208501208183528095505050505050818051602081018201805184825260208301602085012081835280955050505050505f915091509054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b5f610ddb61129a565b90505f815f0160089054906101000a900460ff161590505f825f015f9054906101000a900467ffffffffffffffff1690505f808267ffffffffffffffff16148015610e235750825b90505f60018367ffffffffffffffff16148015610e5657505f3073ffffffffffffffffffffffffffffffffffffffff163b145b905081158015610e64575080155b15610e9b576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001855f015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055508315610ee8576001855f0160086101000a81548160ff0219169083151502179055505b610ef1336112c1565b855f806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508315610f8a575f855f0160086101000a81548160ff0219169083151502179055507fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d26001604051610f8191906124df565b60405180910390a15b505050505050565b610f9a61110f565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361100a575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161100191906118bd565b60405180910390fd5b611013816111a2565b50565b61101e61110f565b8060038686604051611031929190612526565b9081526020016040518091039020848460405161104f929190612526565b90815260200160405180910390205f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507f7ab7074f7826dda120af25c3924584f012f6cc3108c1185301c5461ac75f800785858585856040516110d195949392919061256a565b60405180910390a15050505050565b5f816040516020016110f2919061261b565b604051602081830303815290604052805190602001209050919050565b6111176112d5565b73ffffffffffffffffffffffffffffffffffffffff16611135610c19565b73ffffffffffffffffffffffffffffffffffffffff1614611194576111586112d5565b6040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161118b91906118bd565b60405180910390fd5b565b5f804690508091505090565b5f6111ab611273565b90505f815f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905082825f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508273ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3505050565b5f7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300905090565b5f7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00905090565b6112c96112dc565b6112d28161131c565b50565b5f33905090565b6112e46113a0565b61131a576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6113246112dc565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603611394575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161138b91906118bd565b60405180910390fd5b61139d816111a2565b50565b5f6113a961129a565b5f0160089054906101000a900460ff16905090565b5f604051905090565b5f80fd5b5f80fd5b5f819050919050565b6113e1816113cf565b81146113eb575f80fd5b50565b5f813590506113fc816113d8565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61142b82611402565b9050919050565b5f61143c82611421565b9050919050565b61144c81611432565b8114611456575f80fd5b50565b5f8135905061146781611443565b92915050565b61147681611421565b8114611480575f80fd5b50565b5f813590506114918161146d565b92915050565b5f80fd5b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6114e58261149f565b810181811067ffffffffffffffff82111715611504576115036114af565b5b80604052505050565b5f6115166113be565b905061152282826114dc565b919050565b5f67ffffffffffffffff821115611541576115406114af565b5b61154a8261149f565b9050602081019050919050565b828183375f83830152505050565b5f61157761157284611527565b61150d565b9050828152602081018484840111156115935761159261149b565b5b61159e848285611557565b509392505050565b5f82601f8301126115ba576115b9611497565b5b81356115ca848260208601611565565b91505092915050565b5f805f805f60a086880312156115ec576115eb6113c7565b5b5f6115f9888289016113ee565b955050602061160a88828901611459565b945050604061161b88828901611483565b935050606061162c888289016113ee565b925050608086013567ffffffffffffffff81111561164d5761164c6113cb565b5b611659888289016115a6565b9150509295509295909350565b5f61167082611421565b9050919050565b61168081611666565b811461168a575f80fd5b50565b5f8135905061169b81611677565b92915050565b5f602082840312156116b6576116b56113c7565b5b5f6116c38482850161168d565b91505092915050565b5f80fd5b5f61010082840312156116e6576116e56116cc565b5b81905092915050565b5f805f806101608587031215611708576117076113c7565b5b5f611715878288016116d0565b94505061010061172787828801611483565b935050610120611739878288016113ee565b92505061014085013567ffffffffffffffff81111561175b5761175a6113cb565b5b611767878288016115a6565b91505092959194509250565b5f805f610140848603121561178b5761178a6113c7565b5b5f611798868287016116d0565b9350506101006117aa86828701611483565b9250506101206117bc868287016113ee565b9150509250925092565b5f819050919050565b6117d8816117c6565b82525050565b5f6020820190506117f15f8301846117cf565b92915050565b5f6020828403121561180c5761180b6113c7565b5b5f61181984828501611483565b91505092915050565b5f805f806080858703121561183a576118396113c7565b5b5f61184787828801611483565b9450506020611858878288016113ee565b9350506040611869878288016113ee565b925050606061187a87828801611459565b91505092959194509250565b61188f816113cf565b82525050565b5f6020820190506118a85f830184611886565b92915050565b6118b781611421565b82525050565b5f6020820190506118d05f8301846118ae565b92915050565b6118df816117c6565b81146118e9575f80fd5b50565b5f813590506118fa816118d6565b92915050565b5f8060408385031215611916576119156113c7565b5b5f611923858286016118ec565b925050602083013567ffffffffffffffff811115611944576119436113cb565b5b611950858286016115a6565b9150509250929050565b5f6020828403121561196f5761196e6113c7565b5b5f82013567ffffffffffffffff81111561198c5761198b6113cb565b5b611998848285016115a6565b91505092915050565b5f60ff82169050919050565b6119b6816119a1565b82525050565b5f6060820190506119cf5f8301866117cf565b6119dc60208301856117cf565b6119e960408301846119ad565b949350505050565b5f67ffffffffffffffff821115611a0b57611a0a6114af565b5b611a148261149f565b9050602081019050919050565b5f611a33611a2e846119f1565b61150d565b905082815260208101848484011115611a4f57611a4e61149b565b5b611a5a848285611557565b509392505050565b5f82601f830112611a7657611a75611497565b5b8135611a86848260208601611a21565b91505092915050565b5f8060408385031215611aa557611aa46113c7565b5b5f83013567ffffffffffffffff811115611ac257611ac16113cb565b5b611ace85828601611a62565b925050602083013567ffffffffffffffff811115611aef57611aee6113cb565b5b611afb85828601611a62565b9150509250929050565b5f819050919050565b5f611b28611b23611b1e84611402565b611b05565b611402565b9050919050565b5f611b3982611b0e565b9050919050565b5f611b4a82611b2f565b9050919050565b611b5a81611b40565b82525050565b5f602082019050611b735f830184611b51565b92915050565b5f80fd5b5f80fd5b5f8083601f840112611b9657611b95611497565b5b8235905067ffffffffffffffff811115611bb357611bb2611b79565b5b602083019150836001820283011115611bcf57611bce611b7d565b5b9250929050565b5f805f805f60608688031215611bef57611bee6113c7565b5b5f86013567ffffffffffffffff811115611c0c57611c0b6113cb565b5b611c1888828901611b81565b9550955050602086013567ffffffffffffffff811115611c3b57611c3a6113cb565b5b611c4788828901611b81565b93509350506040611c5a88828901611483565b9150509295509295909350565b5f60208284031215611c7c57611c7b6113c7565b5b5f611c89848285016118ec565b91505092915050565b5f82825260208201905092915050565b7f4272696467653a206e6f6e6365000000000000000000000000000000000000005f82015250565b5f611cd6600d83611c92565b9150611ce182611ca2565b602082019050919050565b5f6020820190508181035f830152611d0381611cca565b9050919050565b7f4272696467653a20696e76616c6964207369676e6174757265000000000000005f82015250565b5f611d3e601983611c92565b9150611d4982611d0a565b602082019050919050565b5f6020820190508181035f830152611d6b81611d32565b9050919050565b5f604082019050611d855f8301856118ae565b611d926020830184611886565b9392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f611dd0826113cf565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611e0257611e01611d99565b5b600182019050919050565b5f608082019050611e205f8301876118ae565b611e2d60208301866118ae565b611e3a6040830185611886565b611e476060830184611886565b95945050505050565b7f4272696467653a2066726f6d00000000000000000000000000000000000000005f82015250565b5f611e84600c83611c92565b9150611e8f82611e50565b602082019050919050565b5f6020820190508181035f830152611eb181611e78565b9050919050565b5f8115159050919050565b611ecc81611eb8565b8114611ed6575f80fd5b50565b5f81519050611ee781611ec3565b92915050565b5f60208284031215611f0257611f016113c7565b5b5f611f0f84828501611ed9565b91505092915050565b5f611f266020840184611483565b905092915050565b611f3781611421565b82525050565b5f62ffffff82169050919050565b611f5481611f3d565b8114611f5e575f80fd5b50565b5f81359050611f6f81611f4b565b92915050565b5f611f836020840184611f61565b905092915050565b611f9481611f3d565b82525050565b5f611fa860208401846113ee565b905092915050565b611fb9816113cf565b82525050565b611fc881611402565b8114611fd2575f80fd5b50565b5f81359050611fe381611fbf565b92915050565b5f611ff76020840184611fd5565b905092915050565b61200881611402565b82525050565b610100820161201f5f830183611f18565b61202b5f850182611f2e565b506120396020830183611f18565b6120466020850182611f2e565b506120546040830183611f75565b6120616040850182611f8b565b5061206f6060830183611f18565b61207c6060850182611f2e565b5061208a6080830183611f9a565b6120976080850182611fb0565b506120a560a0830183611f9a565b6120b260a0850182611fb0565b506120c060c0830183611f9a565b6120cd60c0850182611fb0565b506120db60e0830183611fe9565b6120e860e0850182611fff565b50505050565b5f610100820190506121025f83018461200e565b92915050565b5f81519050612116816113d8565b92915050565b5f60208284031215612131576121306113c7565b5b5f61213e84828501612108565b91505092915050565b5f60a08201905061215a5f8301886118ae565b61216760208301876118ae565b61217460408301866118ae565b6121816060830185611886565b61218e6080830184611886565b9695505050505050565b5f602082840312156121ad576121ac6113c7565b5b5f6121ba84828501611f61565b91505092915050565b5f602082840312156121d8576121d76113c7565b5b5f6121e584828501611fd5565b91505092915050565b5f8160601b9050919050565b5f612204826121ee565b9050919050565b5f612215826121fa565b9050919050565b61222d61222882611421565b61220b565b82525050565b5f8160e81b9050919050565b5f61224982612233565b9050919050565b61226161225c82611f3d565b61223f565b82525050565b5f819050919050565b61228161227c826113cf565b612267565b82525050565b61229861229382611402565b6121fa565b82525050565b5f6122a9828e61221c565b6014820191506122b9828d61221c565b6014820191506122c9828c612250565b6003820191506122d9828b61221c565b6014820191506122e9828a612270565b6020820191506122f98289612270565b6020820191506123098288612270565b6020820191506123198287612287565b6014820191506123298286612270565b602082019150612339828561221c565b6014820191506123498284612270565b6020820191508190509c9b505050505050505050505050565b5f61236c82611b2f565b9050919050565b61238461237f82612362565b61220b565b82525050565b5f612395828861221c565b6014820191506123a58287612270565b6020820191506123b58286612270565b6020820191506123c58285612373565b6014820191506123d58284612270565b6020820191508190509695505050505050565b5f6080820190506123fb5f8301876117cf565b61240860208301866119ad565b61241560408301856117cf565b61242260608301846117cf565b95945050505050565b7f696e76616c6964207369676e6174757265206c656e67746800000000000000005f82015250565b5f61245f601883611c92565b915061246a8261242b565b602082019050919050565b5f6020820190508181035f83015261248c81612453565b9050919050565b5f819050919050565b5f67ffffffffffffffff82169050919050565b5f6124c96124c46124bf84612493565b611b05565b61249c565b9050919050565b6124d9816124af565b82525050565b5f6020820190506124f25f8301846124d0565b92915050565b5f81905092915050565b5f61250d83856124f8565b935061251a838584611557565b82840190509392505050565b5f612532828486612502565b91508190509392505050565b5f6125498385611c92565b9350612556838584611557565b61255f8361149f565b840190509392505050565b5f6060820190508181035f83015261258381878961253e565b9050818103602083015261259881858761253e565b90506125a760408301846118ae565b9695505050505050565b7f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f82015250565b5f6125e5601c836124f8565b91506125f0826125b1565b601c82019050919050565b5f819050919050565b612615612610826117c6565b6125fb565b82525050565b5f612625826125d9565b91506126318284612604565b6020820191508190509291505056fea2646970667358221220f470cd5037538cc6e4de72d3abce3afb16e307317d64666a6a00a2c006d329a764736f6c63430008190033",
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
