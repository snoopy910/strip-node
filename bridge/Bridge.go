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

// BridgeMetaData contains all meta data concerning the Bridge contract.
var BridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"AuthorityChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"chainId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"srcToken\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"peggedToken\",\"type\":\"address\"}],\"name\":\"TokenAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"TokenMinted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chainId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"srcToken\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"peggedToken\",\"type\":\"address\"}],\"name\":\"addToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authority\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_messageHash\",\"type\":\"bytes32\"}],\"name\":\"getEthSignedMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"contractBridgeToken\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_authority\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"contractBridgeToken\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"mintNonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"peggedTokens\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_ethSignedMessageHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_signature\",\"type\":\"bytes\"}],\"name\":\"recoverSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_authority\",\"type\":\"address\"}],\"name\":\"setAuthority\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"splitSignature\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f80fd5b506119f18061001c5f395ff3fe608060405234801561000f575f80fd5b50600436106100e8575f3560e01c8063bf7e214f1161008a578063c4d66de811610064578063c4d66de81461025c578063f2fde38b14610278578063fa3d24f214610294578063fa540801146102b0576100e8565b8063bf7e214f146101de578063bfa4bf9a146101fc578063bfb4ad0c1461022c576100e8565b80638da5cb5b116100c65780638da5cb5b1461012e5780639540bce81461014c57806397aba7f91461017c578063a7bb5803146101ac576100e8565b80632adfefeb146100ec578063715018a6146101085780637a9e5e4b14610112575b5f80fd5b61010660048036038101906101019190610f94565b6102e0565b005b610110610515565b005b61012c60048036038101906101279190611027565b610528565b005b6101366105a9565b6040516101439190611061565b60405180910390f35b61016660048036038101906101619190611027565b6105de565b6040516101739190611089565b60405180910390f35b610196600480360381019061019191906110d5565b6105f3565b6040516101a39190611061565b60405180910390f35b6101c660048036038101906101c1919061112f565b61065d565b6040516101d5939291906111a0565b60405180910390f35b6101e66106c2565b6040516101f39190611061565b60405180910390f35b610216600480360381019061021191906111d5565b6106e5565b6040516102239190611239565b60405180910390f35b610246600480360381019061024191906112f0565b610726565b6040516102539190611061565b60405180910390f35b61027660048036038101906102719190611027565b610793565b005b610292600480360381019061028d9190611027565b610953565b005b6102ae60048036038101906102a991906113c3565b6109d7565b005b6102ca60048036038101906102c59190611454565b610aa1565b6040516102d79190611239565b60405180910390f35b8160015f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20541461035f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610356906114d9565b60405180910390fd5b5f61036c848488886106e5565b90505f61037882610aa1565b90505f61038582856105f3565b90505f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610414576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161040b90611541565b60405180910390fd5b8673ffffffffffffffffffffffffffffffffffffffff166340c10f19878a6040518363ffffffff1660e01b815260040161044f92919061155f565b5f604051808303815f87803b158015610466575f80fd5b505af1158015610478573d5f803e3d5ffd5b5050505060015f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8154809291906104c9906115b3565b91905055507f747bd6dbfd6ceb446b50b008eeade0e74f807993dd969546d7efd6008554b1d087878a8860405161050394939291906115fa565b60405180910390a15050505050505050565b61051d610ad0565b6105265f610b57565b565b610530610ad0565b805f806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507f3430ad8dbed7c32bf49006f0d79d2ab70785ea13ebd4ef7d1b87e487ef08928c8160405161059e9190611061565b60405180910390a150565b5f806105b3610c28565b9050805f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1691505090565b6001602052805f5260405f205f915090505481565b5f805f806106008561065d565b9250925092506001868285856040515f8152602001604052604051610628949392919061163d565b6020604051602081039080840390855afa158015610648573d5f803e3d5ffd5b50505060206040510351935050505092915050565b5f805f60418451146106a4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161069b906116ca565b60405180910390fd5b602084015192506040840151915060608401515f1a90509193909250565b5f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b5f848484846106f2610c4f565b6040516020016107069594939291906117b0565b604051602081830303815290604052805190602001209050949350505050565b600282805160208101820180518482526020830160208501208183528095505050505050818051602081018201805184825260208301602085012081835280955050505050505f915091509054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b5f61079c610c5b565b90505f815f0160089054906101000a900460ff161590505f825f015f9054906101000a900467ffffffffffffffff1690505f808267ffffffffffffffff161480156107e45750825b90505f60018367ffffffffffffffff1614801561081757505f3073ffffffffffffffffffffffffffffffffffffffff163b145b905081158015610825575080155b1561085c576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001855f015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555083156108a9576001855f0160086101000a81548160ff0219169083151502179055505b6108b233610c82565b855f806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550831561094b575f855f0160086101000a81548160ff0219169083151502179055507fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d26001604051610942919061185a565b60405180910390a15b505050505050565b61095b610ad0565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036109cb575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016109c29190611061565b60405180910390fd5b6109d481610b57565b50565b6109df610ad0565b80600286866040516109f29291906118a1565b90815260200160405180910390208484604051610a109291906118a1565b90815260200160405180910390205f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507f7ab7074f7826dda120af25c3924584f012f6cc3108c1185301c5461ac75f80078585858585604051610a929594939291906118e5565b60405180910390a15050505050565b5f81604051602001610ab39190611996565b604051602081830303815290604052805190602001209050919050565b610ad8610c96565b73ffffffffffffffffffffffffffffffffffffffff16610af66105a9565b73ffffffffffffffffffffffffffffffffffffffff1614610b5557610b19610c96565b6040517f118cdaa7000000000000000000000000000000000000000000000000000000008152600401610b4c9190611061565b60405180910390fd5b565b5f610b60610c28565b90505f815f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905082825f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508273ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3505050565b5f7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300905090565b5f804690508091505090565b5f7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00905090565b610c8a610c9d565b610c9381610cdd565b50565b5f33905090565b610ca5610d61565b610cdb576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b610ce5610c9d565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610d55575f6040517f1e4fbdf7000000000000000000000000000000000000000000000000000000008152600401610d4c9190611061565b60405180910390fd5b610d5e81610b57565b50565b5f610d6a610c5b565b5f0160089054906101000a900460ff16905090565b5f604051905090565b5f80fd5b5f80fd5b5f819050919050565b610da281610d90565b8114610dac575f80fd5b50565b5f81359050610dbd81610d99565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610dec82610dc3565b9050919050565b5f610dfd82610de2565b9050919050565b610e0d81610df3565b8114610e17575f80fd5b50565b5f81359050610e2881610e04565b92915050565b610e3781610de2565b8114610e41575f80fd5b50565b5f81359050610e5281610e2e565b92915050565b5f80fd5b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b610ea682610e60565b810181811067ffffffffffffffff82111715610ec557610ec4610e70565b5b80604052505050565b5f610ed7610d7f565b9050610ee38282610e9d565b919050565b5f67ffffffffffffffff821115610f0257610f01610e70565b5b610f0b82610e60565b9050602081019050919050565b828183375f83830152505050565b5f610f38610f3384610ee8565b610ece565b905082815260208101848484011115610f5457610f53610e5c565b5b610f5f848285610f18565b509392505050565b5f82601f830112610f7b57610f7a610e58565b5b8135610f8b848260208601610f26565b91505092915050565b5f805f805f60a08688031215610fad57610fac610d88565b5b5f610fba88828901610daf565b9550506020610fcb88828901610e1a565b9450506040610fdc88828901610e44565b9350506060610fed88828901610daf565b925050608086013567ffffffffffffffff81111561100e5761100d610d8c565b5b61101a88828901610f67565b9150509295509295909350565b5f6020828403121561103c5761103b610d88565b5b5f61104984828501610e44565b91505092915050565b61105b81610de2565b82525050565b5f6020820190506110745f830184611052565b92915050565b61108381610d90565b82525050565b5f60208201905061109c5f83018461107a565b92915050565b5f819050919050565b6110b4816110a2565b81146110be575f80fd5b50565b5f813590506110cf816110ab565b92915050565b5f80604083850312156110eb576110ea610d88565b5b5f6110f8858286016110c1565b925050602083013567ffffffffffffffff81111561111957611118610d8c565b5b61112585828601610f67565b9150509250929050565b5f6020828403121561114457611143610d88565b5b5f82013567ffffffffffffffff81111561116157611160610d8c565b5b61116d84828501610f67565b91505092915050565b61117f816110a2565b82525050565b5f60ff82169050919050565b61119a81611185565b82525050565b5f6060820190506111b35f830186611176565b6111c06020830185611176565b6111cd6040830184611191565b949350505050565b5f805f80608085870312156111ed576111ec610d88565b5b5f6111fa87828801610e44565b945050602061120b87828801610daf565b935050604061121c87828801610daf565b925050606061122d87828801610e1a565b91505092959194509250565b5f60208201905061124c5f830184611176565b92915050565b5f67ffffffffffffffff82111561126c5761126b610e70565b5b61127582610e60565b9050602081019050919050565b5f61129461128f84611252565b610ece565b9050828152602081018484840111156112b0576112af610e5c565b5b6112bb848285610f18565b509392505050565b5f82601f8301126112d7576112d6610e58565b5b81356112e7848260208601611282565b91505092915050565b5f806040838503121561130657611305610d88565b5b5f83013567ffffffffffffffff81111561132357611322610d8c565b5b61132f858286016112c3565b925050602083013567ffffffffffffffff8111156113505761134f610d8c565b5b61135c858286016112c3565b9150509250929050565b5f80fd5b5f80fd5b5f8083601f84011261138357611382610e58565b5b8235905067ffffffffffffffff8111156113a05761139f611366565b5b6020830191508360018202830111156113bc576113bb61136a565b5b9250929050565b5f805f805f606086880312156113dc576113db610d88565b5b5f86013567ffffffffffffffff8111156113f9576113f8610d8c565b5b6114058882890161136e565b9550955050602086013567ffffffffffffffff81111561142857611427610d8c565b5b6114348882890161136e565b9350935050604061144788828901610e44565b9150509295509295909350565b5f6020828403121561146957611468610d88565b5b5f611476848285016110c1565b91505092915050565b5f82825260208201905092915050565b7f4272696467653a206e6f6e6365000000000000000000000000000000000000005f82015250565b5f6114c3600d8361147f565b91506114ce8261148f565b602082019050919050565b5f6020820190508181035f8301526114f0816114b7565b9050919050565b7f4272696467653a20696e76616c6964207369676e6174757265000000000000005f82015250565b5f61152b60198361147f565b9150611536826114f7565b602082019050919050565b5f6020820190508181035f8301526115588161151f565b9050919050565b5f6040820190506115725f830185611052565b61157f602083018461107a565b9392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6115bd82610d90565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036115ef576115ee611586565b5b600182019050919050565b5f60808201905061160d5f830187611052565b61161a6020830186611052565b611627604083018561107a565b611634606083018461107a565b95945050505050565b5f6080820190506116505f830187611176565b61165d6020830186611191565b61166a6040830185611176565b6116776060830184611176565b95945050505050565b7f696e76616c6964207369676e6174757265206c656e67746800000000000000005f82015250565b5f6116b460188361147f565b91506116bf82611680565b602082019050919050565b5f6020820190508181035f8301526116e1816116a8565b9050919050565b5f8160601b9050919050565b5f6116fe826116e8565b9050919050565b5f61170f826116f4565b9050919050565b61172761172282610de2565b611705565b82525050565b5f819050919050565b61174761174282610d90565b61172d565b82525050565b5f819050919050565b5f61177061176b61176684610dc3565b61174d565b610dc3565b9050919050565b5f61178182611756565b9050919050565b5f61179282611777565b9050919050565b6117aa6117a582611788565b611705565b82525050565b5f6117bb8288611716565b6014820191506117cb8287611736565b6020820191506117db8286611736565b6020820191506117eb8285611799565b6014820191506117fb8284611736565b6020820191508190509695505050505050565b5f819050919050565b5f67ffffffffffffffff82169050919050565b5f61184461183f61183a8461180e565b61174d565b611817565b9050919050565b6118548161182a565b82525050565b5f60208201905061186d5f83018461184b565b92915050565b5f81905092915050565b5f6118888385611873565b9350611895838584610f18565b82840190509392505050565b5f6118ad82848661187d565b91508190509392505050565b5f6118c4838561147f565b93506118d1838584610f18565b6118da83610e60565b840190509392505050565b5f6060820190508181035f8301526118fe8187896118b9565b905081810360208301526119138185876118b9565b90506119226040830184611052565b9695505050505050565b7f19457468657265756d205369676e6564204d6573736167653a0a3332000000005f82015250565b5f611960601c83611873565b915061196b8261192c565b601c82019050919050565b5f819050919050565b61199061198b826110a2565b611976565b82525050565b5f6119a082611954565b91506119ac828461197f565b6020820191508190509291505056fea2646970667358221220458e63f5bcfbfc1db8acc2b5ad7f887529be94b853a7b1cb1a36ac2ca891f8d964736f6c63430008190033",
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

// GetMessageHash is a free data retrieval call binding the contract method 0xbfa4bf9a.
//
// Solidity: function getMessageHash(address account, uint256 nonce, uint256 amount, address token) view returns(bytes32)
func (_Bridge *BridgeCaller) GetMessageHash(opts *bind.CallOpts, account common.Address, nonce *big.Int, amount *big.Int, token common.Address) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getMessageHash", account, nonce, amount, token)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetMessageHash is a free data retrieval call binding the contract method 0xbfa4bf9a.
//
// Solidity: function getMessageHash(address account, uint256 nonce, uint256 amount, address token) view returns(bytes32)
func (_Bridge *BridgeSession) GetMessageHash(account common.Address, nonce *big.Int, amount *big.Int, token common.Address) ([32]byte, error) {
	return _Bridge.Contract.GetMessageHash(&_Bridge.CallOpts, account, nonce, amount, token)
}

// GetMessageHash is a free data retrieval call binding the contract method 0xbfa4bf9a.
//
// Solidity: function getMessageHash(address account, uint256 nonce, uint256 amount, address token) view returns(bytes32)
func (_Bridge *BridgeCallerSession) GetMessageHash(account common.Address, nonce *big.Int, amount *big.Int, token common.Address) ([32]byte, error) {
	return _Bridge.Contract.GetMessageHash(&_Bridge.CallOpts, account, nonce, amount, token)
}

// MintNonces is a free data retrieval call binding the contract method 0x9540bce8.
//
// Solidity: function mintNonces(address ) view returns(uint256)
func (_Bridge *BridgeCaller) MintNonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "mintNonces", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MintNonces is a free data retrieval call binding the contract method 0x9540bce8.
//
// Solidity: function mintNonces(address ) view returns(uint256)
func (_Bridge *BridgeSession) MintNonces(arg0 common.Address) (*big.Int, error) {
	return _Bridge.Contract.MintNonces(&_Bridge.CallOpts, arg0)
}

// MintNonces is a free data retrieval call binding the contract method 0x9540bce8.
//
// Solidity: function mintNonces(address ) view returns(uint256)
func (_Bridge *BridgeCallerSession) MintNonces(arg0 common.Address) (*big.Int, error) {
	return _Bridge.Contract.MintNonces(&_Bridge.CallOpts, arg0)
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
