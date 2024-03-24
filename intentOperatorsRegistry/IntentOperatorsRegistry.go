// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package intentoperatorsregistry

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

// IntentOperatorsRegistryMetaData contains all meta data concerning the IntentOperatorsRegistry contract.
var IntentOperatorsRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maximumSigners\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"publickey\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"added\",\"type\":\"bool\"}],\"name\":\"SignerUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAXIMUM_SIGNERS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_publickey\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"}],\"name\":\"addSigner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_publickey\",\"type\":\"bytes32\"}],\"name\":\"removeSigner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"signers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"whitelisted\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a0604052348015600e575f80fd5b50604051610e58380380610e588339818101604052810190602e9190606d565b8060808181525050506093565b5f80fd5b5f819050919050565b604f81603f565b81146058575f80fd5b50565b5f815190506067816048565b92915050565b5f60208284031215607f57607e603b565b5b5f608a84828501605b565b91505092915050565b608051610dad6100ab5f395f61057a0152610dad5ff3fe608060405234801561000f575f80fd5b5060043610610086575f3560e01c80638da5cb5b116100595780638da5cb5b146100ea5780639194a60914610108578063ab0bba8614610124578063f2fde38b1461014257610086565b8063141774ef1461008a578063715018a6146100ba5780638129fc1c146100c45780638cc6f44c146100ce575b5f80fd5b6100a4600480360381019061009f91906108fe565b61015e565b6040516100b19190610943565b60405180910390f35b6100c2610183565b005b6100cc610196565b005b6100e860048036038101906100e391906108fe565b610316565b005b6100f2610427565b6040516100ff919061099b565b60405180910390f35b610122600480360381019061011d9190610a15565b61045c565b005b61012c610578565b6040516101399190610a8a565b60405180910390f35b61015c60048036038101906101579190610acd565b61059c565b005b5f602052805f5260405f205f91509050805f015f9054906101000a900460ff16905081565b61018b610620565b6101945f6106a7565b565b5f61019f610778565b90505f815f0160089054906101000a900460ff161590505f825f015f9054906101000a900467ffffffffffffffff1690505f808267ffffffffffffffff161480156101e75750825b90505f60018367ffffffffffffffff1614801561021a57505f3073ffffffffffffffffffffffffffffffffffffffff163b145b905081158015610228575080155b1561025f576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001855f015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555083156102ac576001855f0160086101000a81548160ff0219169083151502179055505b6102b53361079f565b831561030f575f855f0160086101000a81548160ff0219169083151502179055507fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d260016040516103069190610b4d565b60405180910390a15b5050505050565b61031e610620565b5f801b8103610362576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161035990610bc0565b60405180910390fd5b600115155f808381526020019081526020015f205f015f9054906101000a900460ff161515146103c7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103be90610c28565b60405180910390fd5b5f808281526020019081526020015f205f8082015f6101000a81549060ff02191690555050807f6eaf83ec4eec8fa4159f63480f8bc9e3f2e39f3fed2e5856d8d103268680e6f05f60405161041c9190610c69565b60405180910390a250565b5f806104316107b3565b9050805f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1691505090565b610464610620565b5f801b83036104a8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161049f90610bc0565b60405180910390fd5b5f15155f808581526020019081526020015f205f015f9054906101000a900460ff1615151461050c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161050390610cdf565b60405180910390fd5b60015f808581526020019081526020015f205f015f6101000a81548160ff021916908315150217905550827f6eaf83ec4eec8fa4159f63480f8bc9e3f2e39f3fed2e5856d8d103268680e6f08383600160405161056b93929190610d47565b60405180910390a2505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b6105a4610620565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610614575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161060b919061099b565b60405180910390fd5b61061d816106a7565b50565b6106286107da565b73ffffffffffffffffffffffffffffffffffffffff16610646610427565b73ffffffffffffffffffffffffffffffffffffffff16146106a5576106696107da565b6040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161069c919061099b565b60405180910390fd5b565b5f6106b06107b3565b90505f815f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905082825f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508273ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3505050565b5f7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00905090565b6107a76107e1565b6107b081610821565b50565b5f7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300905090565b5f33905090565b6107e96108a5565b61081f576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6108296107e1565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610899575f6040517f1e4fbdf7000000000000000000000000000000000000000000000000000000008152600401610890919061099b565b60405180910390fd5b6108a2816106a7565b50565b5f6108ae610778565b5f0160089054906101000a900460ff16905090565b5f80fd5b5f80fd5b5f819050919050565b6108dd816108cb565b81146108e7575f80fd5b50565b5f813590506108f8816108d4565b92915050565b5f60208284031215610913576109126108c3565b5b5f610920848285016108ea565b91505092915050565b5f8115159050919050565b61093d81610929565b82525050565b5f6020820190506109565f830184610934565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6109858261095c565b9050919050565b6109958161097b565b82525050565b5f6020820190506109ae5f83018461098c565b92915050565b5f80fd5b5f80fd5b5f80fd5b5f8083601f8401126109d5576109d46109b4565b5b8235905067ffffffffffffffff8111156109f2576109f16109b8565b5b602083019150836001820283011115610a0e57610a0d6109bc565b5b9250929050565b5f805f60408486031215610a2c57610a2b6108c3565b5b5f610a39868287016108ea565b935050602084013567ffffffffffffffff811115610a5a57610a596108c7565b5b610a66868287016109c0565b92509250509250925092565b5f819050919050565b610a8481610a72565b82525050565b5f602082019050610a9d5f830184610a7b565b92915050565b610aac8161097b565b8114610ab6575f80fd5b50565b5f81359050610ac781610aa3565b92915050565b5f60208284031215610ae257610ae16108c3565b5b5f610aef84828501610ab9565b91505092915050565b5f819050919050565b5f67ffffffffffffffff82169050919050565b5f819050919050565b5f610b37610b32610b2d84610af8565b610b14565b610b01565b9050919050565b610b4781610b1d565b82525050565b5f602082019050610b605f830184610b3e565b92915050565b5f82825260208201905092915050565b7f7075626c6963206b657920697320656d707479000000000000000000000000005f82015250565b5f610baa601383610b66565b9150610bb582610b76565b602082019050919050565b5f6020820190508181035f830152610bd781610b9e565b9050919050565b7f7369676e6572206973206e6f74206164646564207965740000000000000000005f82015250565b5f610c12601783610b66565b9150610c1d82610bde565b602082019050919050565b5f6020820190508181035f830152610c3f81610c06565b9050919050565b50565b5f610c545f83610b66565b9150610c5f82610c46565b5f82019050919050565b5f6040820190508181035f830152610c8081610c49565b9050610c8f6020830184610934565b92915050565b7f7369676e657220697320616c72656164792061646465640000000000000000005f82015250565b5f610cc9601783610b66565b9150610cd482610c95565b602082019050919050565b5f6020820190508181035f830152610cf681610cbd565b9050919050565b828183375f83830152505050565b5f601f19601f8301169050919050565b5f610d268385610b66565b9350610d33838584610cfd565b610d3c83610d0b565b840190509392505050565b5f6040820190508181035f830152610d60818587610d1b565b9050610d6f6020830184610934565b94935050505056fea2646970667358221220bd055e4311094b521a5af5eee0080575055dd9845ef79bde9157a25a22fded0864736f6c63430008190033",
}

// IntentOperatorsRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use IntentOperatorsRegistryMetaData.ABI instead.
var IntentOperatorsRegistryABI = IntentOperatorsRegistryMetaData.ABI

// IntentOperatorsRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use IntentOperatorsRegistryMetaData.Bin instead.
var IntentOperatorsRegistryBin = IntentOperatorsRegistryMetaData.Bin

// DeployIntentOperatorsRegistry deploys a new Ethereum contract, binding an instance of IntentOperatorsRegistry to it.
func DeployIntentOperatorsRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, maximumSigners *big.Int) (common.Address, *types.Transaction, *IntentOperatorsRegistry, error) {
	parsed, err := IntentOperatorsRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(IntentOperatorsRegistryBin), backend, maximumSigners)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &IntentOperatorsRegistry{IntentOperatorsRegistryCaller: IntentOperatorsRegistryCaller{contract: contract}, IntentOperatorsRegistryTransactor: IntentOperatorsRegistryTransactor{contract: contract}, IntentOperatorsRegistryFilterer: IntentOperatorsRegistryFilterer{contract: contract}}, nil
}

// IntentOperatorsRegistry is an auto generated Go binding around an Ethereum contract.
type IntentOperatorsRegistry struct {
	IntentOperatorsRegistryCaller     // Read-only binding to the contract
	IntentOperatorsRegistryTransactor // Write-only binding to the contract
	IntentOperatorsRegistryFilterer   // Log filterer for contract events
}

// IntentOperatorsRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type IntentOperatorsRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IntentOperatorsRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IntentOperatorsRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IntentOperatorsRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IntentOperatorsRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IntentOperatorsRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IntentOperatorsRegistrySession struct {
	Contract     *IntentOperatorsRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// IntentOperatorsRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IntentOperatorsRegistryCallerSession struct {
	Contract *IntentOperatorsRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// IntentOperatorsRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IntentOperatorsRegistryTransactorSession struct {
	Contract     *IntentOperatorsRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// IntentOperatorsRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type IntentOperatorsRegistryRaw struct {
	Contract *IntentOperatorsRegistry // Generic contract binding to access the raw methods on
}

// IntentOperatorsRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IntentOperatorsRegistryCallerRaw struct {
	Contract *IntentOperatorsRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// IntentOperatorsRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IntentOperatorsRegistryTransactorRaw struct {
	Contract *IntentOperatorsRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIntentOperatorsRegistry creates a new instance of IntentOperatorsRegistry, bound to a specific deployed contract.
func NewIntentOperatorsRegistry(address common.Address, backend bind.ContractBackend) (*IntentOperatorsRegistry, error) {
	contract, err := bindIntentOperatorsRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IntentOperatorsRegistry{IntentOperatorsRegistryCaller: IntentOperatorsRegistryCaller{contract: contract}, IntentOperatorsRegistryTransactor: IntentOperatorsRegistryTransactor{contract: contract}, IntentOperatorsRegistryFilterer: IntentOperatorsRegistryFilterer{contract: contract}}, nil
}

// NewIntentOperatorsRegistryCaller creates a new read-only instance of IntentOperatorsRegistry, bound to a specific deployed contract.
func NewIntentOperatorsRegistryCaller(address common.Address, caller bind.ContractCaller) (*IntentOperatorsRegistryCaller, error) {
	contract, err := bindIntentOperatorsRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IntentOperatorsRegistryCaller{contract: contract}, nil
}

// NewIntentOperatorsRegistryTransactor creates a new write-only instance of IntentOperatorsRegistry, bound to a specific deployed contract.
func NewIntentOperatorsRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*IntentOperatorsRegistryTransactor, error) {
	contract, err := bindIntentOperatorsRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IntentOperatorsRegistryTransactor{contract: contract}, nil
}

// NewIntentOperatorsRegistryFilterer creates a new log filterer instance of IntentOperatorsRegistry, bound to a specific deployed contract.
func NewIntentOperatorsRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*IntentOperatorsRegistryFilterer, error) {
	contract, err := bindIntentOperatorsRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IntentOperatorsRegistryFilterer{contract: contract}, nil
}

// bindIntentOperatorsRegistry binds a generic wrapper to an already deployed contract.
func bindIntentOperatorsRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IntentOperatorsRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IntentOperatorsRegistry *IntentOperatorsRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IntentOperatorsRegistry.Contract.IntentOperatorsRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IntentOperatorsRegistry *IntentOperatorsRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.IntentOperatorsRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IntentOperatorsRegistry *IntentOperatorsRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.IntentOperatorsRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IntentOperatorsRegistry *IntentOperatorsRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IntentOperatorsRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IntentOperatorsRegistry *IntentOperatorsRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IntentOperatorsRegistry *IntentOperatorsRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.contract.Transact(opts, method, params...)
}

// MAXIMUMSIGNERS is a free data retrieval call binding the contract method 0xab0bba86.
//
// Solidity: function MAXIMUM_SIGNERS() view returns(uint256)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryCaller) MAXIMUMSIGNERS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IntentOperatorsRegistry.contract.Call(opts, &out, "MAXIMUM_SIGNERS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXIMUMSIGNERS is a free data retrieval call binding the contract method 0xab0bba86.
//
// Solidity: function MAXIMUM_SIGNERS() view returns(uint256)
func (_IntentOperatorsRegistry *IntentOperatorsRegistrySession) MAXIMUMSIGNERS() (*big.Int, error) {
	return _IntentOperatorsRegistry.Contract.MAXIMUMSIGNERS(&_IntentOperatorsRegistry.CallOpts)
}

// MAXIMUMSIGNERS is a free data retrieval call binding the contract method 0xab0bba86.
//
// Solidity: function MAXIMUM_SIGNERS() view returns(uint256)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryCallerSession) MAXIMUMSIGNERS() (*big.Int, error) {
	return _IntentOperatorsRegistry.Contract.MAXIMUMSIGNERS(&_IntentOperatorsRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IntentOperatorsRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IntentOperatorsRegistry *IntentOperatorsRegistrySession) Owner() (common.Address, error) {
	return _IntentOperatorsRegistry.Contract.Owner(&_IntentOperatorsRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryCallerSession) Owner() (common.Address, error) {
	return _IntentOperatorsRegistry.Contract.Owner(&_IntentOperatorsRegistry.CallOpts)
}

// Signers is a free data retrieval call binding the contract method 0x141774ef.
//
// Solidity: function signers(bytes32 ) view returns(bool whitelisted)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryCaller) Signers(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _IntentOperatorsRegistry.contract.Call(opts, &out, "signers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Signers is a free data retrieval call binding the contract method 0x141774ef.
//
// Solidity: function signers(bytes32 ) view returns(bool whitelisted)
func (_IntentOperatorsRegistry *IntentOperatorsRegistrySession) Signers(arg0 [32]byte) (bool, error) {
	return _IntentOperatorsRegistry.Contract.Signers(&_IntentOperatorsRegistry.CallOpts, arg0)
}

// Signers is a free data retrieval call binding the contract method 0x141774ef.
//
// Solidity: function signers(bytes32 ) view returns(bool whitelisted)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryCallerSession) Signers(arg0 [32]byte) (bool, error) {
	return _IntentOperatorsRegistry.Contract.Signers(&_IntentOperatorsRegistry.CallOpts, arg0)
}

// AddSigner is a paid mutator transaction binding the contract method 0x9194a609.
//
// Solidity: function addSigner(bytes32 _publickey, string url) returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistryTransactor) AddSigner(opts *bind.TransactOpts, _publickey [32]byte, url string) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.contract.Transact(opts, "addSigner", _publickey, url)
}

// AddSigner is a paid mutator transaction binding the contract method 0x9194a609.
//
// Solidity: function addSigner(bytes32 _publickey, string url) returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistrySession) AddSigner(_publickey [32]byte, url string) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.AddSigner(&_IntentOperatorsRegistry.TransactOpts, _publickey, url)
}

// AddSigner is a paid mutator transaction binding the contract method 0x9194a609.
//
// Solidity: function addSigner(bytes32 _publickey, string url) returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistryTransactorSession) AddSigner(_publickey [32]byte, url string) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.AddSigner(&_IntentOperatorsRegistry.TransactOpts, _publickey, url)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistryTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistrySession) Initialize() (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.Initialize(&_IntentOperatorsRegistry.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistryTransactorSession) Initialize() (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.Initialize(&_IntentOperatorsRegistry.TransactOpts)
}

// RemoveSigner is a paid mutator transaction binding the contract method 0x8cc6f44c.
//
// Solidity: function removeSigner(bytes32 _publickey) returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistryTransactor) RemoveSigner(opts *bind.TransactOpts, _publickey [32]byte) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.contract.Transact(opts, "removeSigner", _publickey)
}

// RemoveSigner is a paid mutator transaction binding the contract method 0x8cc6f44c.
//
// Solidity: function removeSigner(bytes32 _publickey) returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistrySession) RemoveSigner(_publickey [32]byte) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.RemoveSigner(&_IntentOperatorsRegistry.TransactOpts, _publickey)
}

// RemoveSigner is a paid mutator transaction binding the contract method 0x8cc6f44c.
//
// Solidity: function removeSigner(bytes32 _publickey) returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistryTransactorSession) RemoveSigner(_publickey [32]byte) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.RemoveSigner(&_IntentOperatorsRegistry.TransactOpts, _publickey)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.RenounceOwnership(&_IntentOperatorsRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.RenounceOwnership(&_IntentOperatorsRegistry.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.TransferOwnership(&_IntentOperatorsRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.TransferOwnership(&_IntentOperatorsRegistry.TransactOpts, newOwner)
}

// IntentOperatorsRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the IntentOperatorsRegistry contract.
type IntentOperatorsRegistryInitializedIterator struct {
	Event *IntentOperatorsRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *IntentOperatorsRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IntentOperatorsRegistryInitialized)
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
		it.Event = new(IntentOperatorsRegistryInitialized)
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
func (it *IntentOperatorsRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IntentOperatorsRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IntentOperatorsRegistryInitialized represents a Initialized event raised by the IntentOperatorsRegistry contract.
type IntentOperatorsRegistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*IntentOperatorsRegistryInitializedIterator, error) {

	logs, sub, err := _IntentOperatorsRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &IntentOperatorsRegistryInitializedIterator{contract: _IntentOperatorsRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *IntentOperatorsRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _IntentOperatorsRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IntentOperatorsRegistryInitialized)
				if err := _IntentOperatorsRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) ParseInitialized(log types.Log) (*IntentOperatorsRegistryInitialized, error) {
	event := new(IntentOperatorsRegistryInitialized)
	if err := _IntentOperatorsRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IntentOperatorsRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the IntentOperatorsRegistry contract.
type IntentOperatorsRegistryOwnershipTransferredIterator struct {
	Event *IntentOperatorsRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *IntentOperatorsRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IntentOperatorsRegistryOwnershipTransferred)
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
		it.Event = new(IntentOperatorsRegistryOwnershipTransferred)
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
func (it *IntentOperatorsRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IntentOperatorsRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IntentOperatorsRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the IntentOperatorsRegistry contract.
type IntentOperatorsRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*IntentOperatorsRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IntentOperatorsRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &IntentOperatorsRegistryOwnershipTransferredIterator{contract: _IntentOperatorsRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IntentOperatorsRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IntentOperatorsRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IntentOperatorsRegistryOwnershipTransferred)
				if err := _IntentOperatorsRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*IntentOperatorsRegistryOwnershipTransferred, error) {
	event := new(IntentOperatorsRegistryOwnershipTransferred)
	if err := _IntentOperatorsRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IntentOperatorsRegistrySignerUpdatedIterator is returned from FilterSignerUpdated and is used to iterate over the raw logs and unpacked data for SignerUpdated events raised by the IntentOperatorsRegistry contract.
type IntentOperatorsRegistrySignerUpdatedIterator struct {
	Event *IntentOperatorsRegistrySignerUpdated // Event containing the contract specifics and raw log

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
func (it *IntentOperatorsRegistrySignerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IntentOperatorsRegistrySignerUpdated)
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
		it.Event = new(IntentOperatorsRegistrySignerUpdated)
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
func (it *IntentOperatorsRegistrySignerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IntentOperatorsRegistrySignerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IntentOperatorsRegistrySignerUpdated represents a SignerUpdated event raised by the IntentOperatorsRegistry contract.
type IntentOperatorsRegistrySignerUpdated struct {
	Publickey [32]byte
	Url       string
	Added     bool
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSignerUpdated is a free log retrieval operation binding the contract event 0x6eaf83ec4eec8fa4159f63480f8bc9e3f2e39f3fed2e5856d8d103268680e6f0.
//
// Solidity: event SignerUpdated(bytes32 indexed publickey, string url, bool added)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) FilterSignerUpdated(opts *bind.FilterOpts, publickey [][32]byte) (*IntentOperatorsRegistrySignerUpdatedIterator, error) {

	var publickeyRule []interface{}
	for _, publickeyItem := range publickey {
		publickeyRule = append(publickeyRule, publickeyItem)
	}

	logs, sub, err := _IntentOperatorsRegistry.contract.FilterLogs(opts, "SignerUpdated", publickeyRule)
	if err != nil {
		return nil, err
	}
	return &IntentOperatorsRegistrySignerUpdatedIterator{contract: _IntentOperatorsRegistry.contract, event: "SignerUpdated", logs: logs, sub: sub}, nil
}

// WatchSignerUpdated is a free log subscription operation binding the contract event 0x6eaf83ec4eec8fa4159f63480f8bc9e3f2e39f3fed2e5856d8d103268680e6f0.
//
// Solidity: event SignerUpdated(bytes32 indexed publickey, string url, bool added)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) WatchSignerUpdated(opts *bind.WatchOpts, sink chan<- *IntentOperatorsRegistrySignerUpdated, publickey [][32]byte) (event.Subscription, error) {

	var publickeyRule []interface{}
	for _, publickeyItem := range publickey {
		publickeyRule = append(publickeyRule, publickeyItem)
	}

	logs, sub, err := _IntentOperatorsRegistry.contract.WatchLogs(opts, "SignerUpdated", publickeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IntentOperatorsRegistrySignerUpdated)
				if err := _IntentOperatorsRegistry.contract.UnpackLog(event, "SignerUpdated", log); err != nil {
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

// ParseSignerUpdated is a log parse operation binding the contract event 0x6eaf83ec4eec8fa4159f63480f8bc9e3f2e39f3fed2e5856d8d103268680e6f0.
//
// Solidity: event SignerUpdated(bytes32 indexed publickey, string url, bool added)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) ParseSignerUpdated(log types.Log) (*IntentOperatorsRegistrySignerUpdated, error) {
	event := new(IntentOperatorsRegistrySignerUpdated)
	if err := _IntentOperatorsRegistry.contract.UnpackLog(event, "SignerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
