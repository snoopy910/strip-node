// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package solversregistry

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

// SolversRegistryMetaData contains all meta data concerning the SolversRegistry contract.
var SolversRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"domain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"whitelisted\",\"type\":\"bool\"}],\"name\":\"SolverUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"solvers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"whitelisted\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"domain\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"whitelisted\",\"type\":\"bool\"}],\"name\":\"updateSolver\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f80fd5b50610c1d8061001c5f395ff3fe608060405234801561000f575f80fd5b5060043610610060575f3560e01c80631275d62e146100645780631968fa5314610080578063715018a6146100b05780638129fc1c146100ba5780638da5cb5b146100c4578063f2fde38b146100e2575b5f80fd5b61007e600480360381019061007991906107ae565b6100fe565b005b61009a60048036038101906100959190610943565b6101db565b6040516100a79190610999565b60405180910390f35b6100b8610218565b005b6100c261022b565b005b6100cc6103ab565b6040516100d991906109f1565b60405180910390f35b6100fc60048036038101906100f79190610a34565b6103e0565b005b610106610464565b5f838390500361014b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161014290610ab9565b60405180910390fd5b60405180602001604052808215158152505f848460405161016d929190610b05565b90815260200160405180910390205f820151815f015f6101000a81548160ff0219169083151502179055509050507f4c40c26c445e2acb57aa0f02b1bc6d90de42da3d0e9d181967d3a8d142a666a68383836040516101ce93929190610b49565b60405180910390a1505050565b5f818051602081018201805184825260208301602085012081835280955050505050505f91509050805f015f9054906101000a900460ff16905081565b610220610464565b6102295f6104eb565b565b5f6102346105bc565b90505f815f0160089054906101000a900460ff161590505f825f015f9054906101000a900467ffffffffffffffff1690505f808267ffffffffffffffff1614801561027c5750825b90505f60018367ffffffffffffffff161480156102af57505f3073ffffffffffffffffffffffffffffffffffffffff163b145b9050811580156102bd575080155b156102f4576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001855f015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055508315610341576001855f0160086101000a81548160ff0219169083151502179055505b61034a336105e3565b83156103a4575f855f0160086101000a81548160ff0219169083151502179055507fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2600160405161039b9190610bce565b60405180910390a15b5050505050565b5f806103b56105f7565b9050805f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1691505090565b6103e8610464565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610458575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161044f91906109f1565b60405180910390fd5b610461816104eb565b50565b61046c61061e565b73ffffffffffffffffffffffffffffffffffffffff1661048a6103ab565b73ffffffffffffffffffffffffffffffffffffffff16146104e9576104ad61061e565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016104e091906109f1565b60405180910390fd5b565b5f6104f46105f7565b90505f815f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905082825f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508273ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3505050565b5f7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00905090565b6105eb610625565b6105f481610665565b50565b5f7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300905090565b5f33905090565b61062d6106e9565b610663576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b61066d610625565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036106dd575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016106d491906109f1565b60405180910390fd5b6106e6816104eb565b50565b5f6106f26105bc565b5f0160089054906101000a900460ff16905090565b5f604051905090565b5f80fd5b5f80fd5b5f80fd5b5f80fd5b5f80fd5b5f8083601f84011261073957610738610718565b5b8235905067ffffffffffffffff8111156107565761075561071c565b5b60208301915083600182028301111561077257610771610720565b5b9250929050565b5f8115159050919050565b61078d81610779565b8114610797575f80fd5b50565b5f813590506107a881610784565b92915050565b5f805f604084860312156107c5576107c4610710565b5b5f84013567ffffffffffffffff8111156107e2576107e1610714565b5b6107ee86828701610724565b935093505060206108018682870161079a565b9150509250925092565b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6108558261080f565b810181811067ffffffffffffffff821117156108745761087361081f565b5b80604052505050565b5f610886610707565b9050610892828261084c565b919050565b5f67ffffffffffffffff8211156108b1576108b061081f565b5b6108ba8261080f565b9050602081019050919050565b828183375f83830152505050565b5f6108e76108e284610897565b61087d565b9050828152602081018484840111156109035761090261080b565b5b61090e8482856108c7565b509392505050565b5f82601f83011261092a57610929610718565b5b813561093a8482602086016108d5565b91505092915050565b5f6020828403121561095857610957610710565b5b5f82013567ffffffffffffffff81111561097557610974610714565b5b61098184828501610916565b91505092915050565b61099381610779565b82525050565b5f6020820190506109ac5f83018461098a565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6109db826109b2565b9050919050565b6109eb816109d1565b82525050565b5f602082019050610a045f8301846109e2565b92915050565b610a13816109d1565b8114610a1d575f80fd5b50565b5f81359050610a2e81610a0a565b92915050565b5f60208284031215610a4957610a48610710565b5b5f610a5684828501610a20565b91505092915050565b5f82825260208201905092915050565b7f646f6d61696e20697320656d70747900000000000000000000000000000000005f82015250565b5f610aa3600f83610a5f565b9150610aae82610a6f565b602082019050919050565b5f6020820190508181035f830152610ad081610a97565b9050919050565b5f81905092915050565b5f610aec8385610ad7565b9350610af98385846108c7565b82840190509392505050565b5f610b11828486610ae1565b91508190509392505050565b5f610b288385610a5f565b9350610b358385846108c7565b610b3e8361080f565b840190509392505050565b5f6040820190508181035f830152610b62818587610b1d565b9050610b71602083018461098a565b949350505050565b5f819050919050565b5f67ffffffffffffffff82169050919050565b5f819050919050565b5f610bb8610bb3610bae84610b79565b610b95565b610b82565b9050919050565b610bc881610b9e565b82525050565b5f602082019050610be15f830184610bbf565b9291505056fea26469706673582212203a5fea5ac9444c785e74c601b8ced5292f53415eb5b1b2f1de61fa9ac3d8a96264736f6c63430008190033",
}

// SolversRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use SolversRegistryMetaData.ABI instead.
var SolversRegistryABI = SolversRegistryMetaData.ABI

// SolversRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SolversRegistryMetaData.Bin instead.
var SolversRegistryBin = SolversRegistryMetaData.Bin

// DeploySolversRegistry deploys a new Ethereum contract, binding an instance of SolversRegistry to it.
func DeploySolversRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SolversRegistry, error) {
	parsed, err := SolversRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SolversRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SolversRegistry{SolversRegistryCaller: SolversRegistryCaller{contract: contract}, SolversRegistryTransactor: SolversRegistryTransactor{contract: contract}, SolversRegistryFilterer: SolversRegistryFilterer{contract: contract}}, nil
}

// SolversRegistry is an auto generated Go binding around an Ethereum contract.
type SolversRegistry struct {
	SolversRegistryCaller     // Read-only binding to the contract
	SolversRegistryTransactor // Write-only binding to the contract
	SolversRegistryFilterer   // Log filterer for contract events
}

// SolversRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type SolversRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SolversRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SolversRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SolversRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SolversRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SolversRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SolversRegistrySession struct {
	Contract     *SolversRegistry  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SolversRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SolversRegistryCallerSession struct {
	Contract *SolversRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// SolversRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SolversRegistryTransactorSession struct {
	Contract     *SolversRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// SolversRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type SolversRegistryRaw struct {
	Contract *SolversRegistry // Generic contract binding to access the raw methods on
}

// SolversRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SolversRegistryCallerRaw struct {
	Contract *SolversRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// SolversRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SolversRegistryTransactorRaw struct {
	Contract *SolversRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSolversRegistry creates a new instance of SolversRegistry, bound to a specific deployed contract.
func NewSolversRegistry(address common.Address, backend bind.ContractBackend) (*SolversRegistry, error) {
	contract, err := bindSolversRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SolversRegistry{SolversRegistryCaller: SolversRegistryCaller{contract: contract}, SolversRegistryTransactor: SolversRegistryTransactor{contract: contract}, SolversRegistryFilterer: SolversRegistryFilterer{contract: contract}}, nil
}

// NewSolversRegistryCaller creates a new read-only instance of SolversRegistry, bound to a specific deployed contract.
func NewSolversRegistryCaller(address common.Address, caller bind.ContractCaller) (*SolversRegistryCaller, error) {
	contract, err := bindSolversRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SolversRegistryCaller{contract: contract}, nil
}

// NewSolversRegistryTransactor creates a new write-only instance of SolversRegistry, bound to a specific deployed contract.
func NewSolversRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*SolversRegistryTransactor, error) {
	contract, err := bindSolversRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SolversRegistryTransactor{contract: contract}, nil
}

// NewSolversRegistryFilterer creates a new log filterer instance of SolversRegistry, bound to a specific deployed contract.
func NewSolversRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*SolversRegistryFilterer, error) {
	contract, err := bindSolversRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SolversRegistryFilterer{contract: contract}, nil
}

// bindSolversRegistry binds a generic wrapper to an already deployed contract.
func bindSolversRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SolversRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SolversRegistry *SolversRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SolversRegistry.Contract.SolversRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SolversRegistry *SolversRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SolversRegistry.Contract.SolversRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SolversRegistry *SolversRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SolversRegistry.Contract.SolversRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SolversRegistry *SolversRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SolversRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SolversRegistry *SolversRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SolversRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SolversRegistry *SolversRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SolversRegistry.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SolversRegistry *SolversRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SolversRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SolversRegistry *SolversRegistrySession) Owner() (common.Address, error) {
	return _SolversRegistry.Contract.Owner(&_SolversRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SolversRegistry *SolversRegistryCallerSession) Owner() (common.Address, error) {
	return _SolversRegistry.Contract.Owner(&_SolversRegistry.CallOpts)
}

// Solvers is a free data retrieval call binding the contract method 0x1968fa53.
//
// Solidity: function solvers(string ) view returns(bool whitelisted)
func (_SolversRegistry *SolversRegistryCaller) Solvers(opts *bind.CallOpts, arg0 string) (bool, error) {
	var out []interface{}
	err := _SolversRegistry.contract.Call(opts, &out, "solvers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Solvers is a free data retrieval call binding the contract method 0x1968fa53.
//
// Solidity: function solvers(string ) view returns(bool whitelisted)
func (_SolversRegistry *SolversRegistrySession) Solvers(arg0 string) (bool, error) {
	return _SolversRegistry.Contract.Solvers(&_SolversRegistry.CallOpts, arg0)
}

// Solvers is a free data retrieval call binding the contract method 0x1968fa53.
//
// Solidity: function solvers(string ) view returns(bool whitelisted)
func (_SolversRegistry *SolversRegistryCallerSession) Solvers(arg0 string) (bool, error) {
	return _SolversRegistry.Contract.Solvers(&_SolversRegistry.CallOpts, arg0)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_SolversRegistry *SolversRegistryTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SolversRegistry.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_SolversRegistry *SolversRegistrySession) Initialize() (*types.Transaction, error) {
	return _SolversRegistry.Contract.Initialize(&_SolversRegistry.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_SolversRegistry *SolversRegistryTransactorSession) Initialize() (*types.Transaction, error) {
	return _SolversRegistry.Contract.Initialize(&_SolversRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SolversRegistry *SolversRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SolversRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SolversRegistry *SolversRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _SolversRegistry.Contract.RenounceOwnership(&_SolversRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SolversRegistry *SolversRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SolversRegistry.Contract.RenounceOwnership(&_SolversRegistry.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SolversRegistry *SolversRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SolversRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SolversRegistry *SolversRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SolversRegistry.Contract.TransferOwnership(&_SolversRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SolversRegistry *SolversRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SolversRegistry.Contract.TransferOwnership(&_SolversRegistry.TransactOpts, newOwner)
}

// UpdateSolver is a paid mutator transaction binding the contract method 0x1275d62e.
//
// Solidity: function updateSolver(string domain, bool whitelisted) returns()
func (_SolversRegistry *SolversRegistryTransactor) UpdateSolver(opts *bind.TransactOpts, domain string, whitelisted bool) (*types.Transaction, error) {
	return _SolversRegistry.contract.Transact(opts, "updateSolver", domain, whitelisted)
}

// UpdateSolver is a paid mutator transaction binding the contract method 0x1275d62e.
//
// Solidity: function updateSolver(string domain, bool whitelisted) returns()
func (_SolversRegistry *SolversRegistrySession) UpdateSolver(domain string, whitelisted bool) (*types.Transaction, error) {
	return _SolversRegistry.Contract.UpdateSolver(&_SolversRegistry.TransactOpts, domain, whitelisted)
}

// UpdateSolver is a paid mutator transaction binding the contract method 0x1275d62e.
//
// Solidity: function updateSolver(string domain, bool whitelisted) returns()
func (_SolversRegistry *SolversRegistryTransactorSession) UpdateSolver(domain string, whitelisted bool) (*types.Transaction, error) {
	return _SolversRegistry.Contract.UpdateSolver(&_SolversRegistry.TransactOpts, domain, whitelisted)
}

// SolversRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the SolversRegistry contract.
type SolversRegistryInitializedIterator struct {
	Event *SolversRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *SolversRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SolversRegistryInitialized)
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
		it.Event = new(SolversRegistryInitialized)
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
func (it *SolversRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SolversRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SolversRegistryInitialized represents a Initialized event raised by the SolversRegistry contract.
type SolversRegistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_SolversRegistry *SolversRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*SolversRegistryInitializedIterator, error) {

	logs, sub, err := _SolversRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &SolversRegistryInitializedIterator{contract: _SolversRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_SolversRegistry *SolversRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *SolversRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _SolversRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SolversRegistryInitialized)
				if err := _SolversRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_SolversRegistry *SolversRegistryFilterer) ParseInitialized(log types.Log) (*SolversRegistryInitialized, error) {
	event := new(SolversRegistryInitialized)
	if err := _SolversRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SolversRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SolversRegistry contract.
type SolversRegistryOwnershipTransferredIterator struct {
	Event *SolversRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SolversRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SolversRegistryOwnershipTransferred)
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
		it.Event = new(SolversRegistryOwnershipTransferred)
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
func (it *SolversRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SolversRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SolversRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the SolversRegistry contract.
type SolversRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SolversRegistry *SolversRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SolversRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SolversRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SolversRegistryOwnershipTransferredIterator{contract: _SolversRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SolversRegistry *SolversRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SolversRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SolversRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SolversRegistryOwnershipTransferred)
				if err := _SolversRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_SolversRegistry *SolversRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*SolversRegistryOwnershipTransferred, error) {
	event := new(SolversRegistryOwnershipTransferred)
	if err := _SolversRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SolversRegistrySolverUpdatedIterator is returned from FilterSolverUpdated and is used to iterate over the raw logs and unpacked data for SolverUpdated events raised by the SolversRegistry contract.
type SolversRegistrySolverUpdatedIterator struct {
	Event *SolversRegistrySolverUpdated // Event containing the contract specifics and raw log

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
func (it *SolversRegistrySolverUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SolversRegistrySolverUpdated)
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
		it.Event = new(SolversRegistrySolverUpdated)
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
func (it *SolversRegistrySolverUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SolversRegistrySolverUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SolversRegistrySolverUpdated represents a SolverUpdated event raised by the SolversRegistry contract.
type SolversRegistrySolverUpdated struct {
	Domain      string
	Whitelisted bool
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterSolverUpdated is a free log retrieval operation binding the contract event 0x4c40c26c445e2acb57aa0f02b1bc6d90de42da3d0e9d181967d3a8d142a666a6.
//
// Solidity: event SolverUpdated(string domain, bool whitelisted)
func (_SolversRegistry *SolversRegistryFilterer) FilterSolverUpdated(opts *bind.FilterOpts) (*SolversRegistrySolverUpdatedIterator, error) {

	logs, sub, err := _SolversRegistry.contract.FilterLogs(opts, "SolverUpdated")
	if err != nil {
		return nil, err
	}
	return &SolversRegistrySolverUpdatedIterator{contract: _SolversRegistry.contract, event: "SolverUpdated", logs: logs, sub: sub}, nil
}

// WatchSolverUpdated is a free log subscription operation binding the contract event 0x4c40c26c445e2acb57aa0f02b1bc6d90de42da3d0e9d181967d3a8d142a666a6.
//
// Solidity: event SolverUpdated(string domain, bool whitelisted)
func (_SolversRegistry *SolversRegistryFilterer) WatchSolverUpdated(opts *bind.WatchOpts, sink chan<- *SolversRegistrySolverUpdated) (event.Subscription, error) {

	logs, sub, err := _SolversRegistry.contract.WatchLogs(opts, "SolverUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SolversRegistrySolverUpdated)
				if err := _SolversRegistry.contract.UnpackLog(event, "SolverUpdated", log); err != nil {
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

// ParseSolverUpdated is a log parse operation binding the contract event 0x4c40c26c445e2acb57aa0f02b1bc6d90de42da3d0e9d181967d3a8d142a666a6.
//
// Solidity: event SolverUpdated(string domain, bool whitelisted)
func (_SolversRegistry *SolversRegistryFilterer) ParseSolverUpdated(log types.Log) (*SolversRegistrySolverUpdated, error) {
	event := new(SolversRegistrySolverUpdated)
	if err := _SolversRegistry.contract.UnpackLog(event, "SolverUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
