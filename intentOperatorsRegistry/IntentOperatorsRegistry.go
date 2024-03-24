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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maximumSigners\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"publickey\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"}],\"name\":\"SignedAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"publickey\",\"type\":\"bytes32\"}],\"name\":\"SignedRemoved\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAXIMUM_SIGNERS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_publickey\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"}],\"name\":\"addSigner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_publickey\",\"type\":\"bytes32\"}],\"name\":\"removeSigner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"signers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"whitelisted\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a0604052348015600e575f80fd5b50604051610de0380380610de08339818101604052810190602e9190606d565b8060808181525050506093565b5f80fd5b5f819050919050565b604f81603f565b81146058575f80fd5b50565b5f815190506067816048565b92915050565b5f60208284031215607f57607e603b565b5b5f608a84828501605b565b91505092915050565b608051610d356100ab5f395f6105770152610d355ff3fe608060405234801561000f575f80fd5b5060043610610086575f3560e01c80638da5cb5b116100595780638da5cb5b146100ea5780639194a60914610108578063ab0bba8614610124578063f2fde38b1461014257610086565b8063141774ef1461008a578063715018a6146100ba5780638129fc1c146100c45780638cc6f44c146100ce575b5f80fd5b6100a4600480360381019061009f91906108fb565b61015e565b6040516100b19190610940565b60405180910390f35b6100c2610183565b005b6100cc610196565b005b6100e860048036038101906100e391906108fb565b610316565b005b6100f261041c565b6040516100ff9190610998565b60405180910390f35b610122600480360381019061011d9190610a12565b610451565b005b61012c610575565b6040516101399190610a87565b60405180910390f35b61015c60048036038101906101579190610aca565b610599565b005b5f602052805f5260405f205f91509050805f015f9054906101000a900460ff16905081565b61018b61061d565b6101945f6106a4565b565b5f61019f610775565b90505f815f0160089054906101000a900460ff161590505f825f015f9054906101000a900467ffffffffffffffff1690505f808267ffffffffffffffff161480156101e75750825b90505f60018367ffffffffffffffff1614801561021a57505f3073ffffffffffffffffffffffffffffffffffffffff163b145b905081158015610228575080155b1561025f576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001855f015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555083156102ac576001855f0160086101000a81548160ff0219169083151502179055505b6102b53361079c565b831561030f575f855f0160086101000a81548160ff0219169083151502179055507fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d260016040516103069190610b4a565b60405180910390a15b5050505050565b61031e61061d565b5f801b8103610362576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161035990610bbd565b60405180910390fd5b600115155f808381526020019081526020015f205f015f9054906101000a900460ff161515146103c7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103be90610c25565b60405180910390fd5b5f808281526020019081526020015f205f8082015f6101000a81549060ff02191690555050807f202ba177372096a533cb0be65537787905a2c9a9b25538d8d9f578706b412cb360405160405180910390a250565b5f806104266107b0565b9050805f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1691505090565b61045961061d565b5f801b830361049d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161049490610bbd565b60405180910390fd5b5f15155f808581526020019081526020015f205f015f9054906101000a900460ff16151514610501576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104f890610c8d565b60405180910390fd5b60015f808581526020019081526020015f205f015f6101000a81548160ff021916908315150217905550818160405161053b929190610ce7565b6040518091039020837fd993cade8966730f67b2b2c9d7689b17c3b727d61bf30f122bb3dcc5325fbec560405160405180910390a3505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b6105a161061d565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610611575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016106089190610998565b60405180910390fd5b61061a816106a4565b50565b6106256107d7565b73ffffffffffffffffffffffffffffffffffffffff1661064361041c565b73ffffffffffffffffffffffffffffffffffffffff16146106a2576106666107d7565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016106999190610998565b60405180910390fd5b565b5f6106ad6107b0565b90505f815f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905082825f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508273ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3505050565b5f7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00905090565b6107a46107de565b6107ad8161081e565b50565b5f7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300905090565b5f33905090565b6107e66108a2565b61081c576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6108266107de565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610896575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161088d9190610998565b60405180910390fd5b61089f816106a4565b50565b5f6108ab610775565b5f0160089054906101000a900460ff16905090565b5f80fd5b5f80fd5b5f819050919050565b6108da816108c8565b81146108e4575f80fd5b50565b5f813590506108f5816108d1565b92915050565b5f602082840312156109105761090f6108c0565b5b5f61091d848285016108e7565b91505092915050565b5f8115159050919050565b61093a81610926565b82525050565b5f6020820190506109535f830184610931565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61098282610959565b9050919050565b61099281610978565b82525050565b5f6020820190506109ab5f830184610989565b92915050565b5f80fd5b5f80fd5b5f80fd5b5f8083601f8401126109d2576109d16109b1565b5b8235905067ffffffffffffffff8111156109ef576109ee6109b5565b5b602083019150836001820283011115610a0b57610a0a6109b9565b5b9250929050565b5f805f60408486031215610a2957610a286108c0565b5b5f610a36868287016108e7565b935050602084013567ffffffffffffffff811115610a5757610a566108c4565b5b610a63868287016109bd565b92509250509250925092565b5f819050919050565b610a8181610a6f565b82525050565b5f602082019050610a9a5f830184610a78565b92915050565b610aa981610978565b8114610ab3575f80fd5b50565b5f81359050610ac481610aa0565b92915050565b5f60208284031215610adf57610ade6108c0565b5b5f610aec84828501610ab6565b91505092915050565b5f819050919050565b5f67ffffffffffffffff82169050919050565b5f819050919050565b5f610b34610b2f610b2a84610af5565b610b11565b610afe565b9050919050565b610b4481610b1a565b82525050565b5f602082019050610b5d5f830184610b3b565b92915050565b5f82825260208201905092915050565b7f7075626c6963206b657920697320656d707479000000000000000000000000005f82015250565b5f610ba7601383610b63565b9150610bb282610b73565b602082019050919050565b5f6020820190508181035f830152610bd481610b9b565b9050919050565b7f7369676e6572206973206e6f74206164646564207965740000000000000000005f82015250565b5f610c0f601783610b63565b9150610c1a82610bdb565b602082019050919050565b5f6020820190508181035f830152610c3c81610c03565b9050919050565b7f7369676e657220697320616c72656164792061646465640000000000000000005f82015250565b5f610c77601783610b63565b9150610c8282610c43565b602082019050919050565b5f6020820190508181035f830152610ca481610c6b565b9050919050565b5f81905092915050565b828183375f83830152505050565b5f610cce8385610cab565b9350610cdb838584610cb5565b82840190509392505050565b5f610cf3828486610cc3565b9150819050939250505056fea2646970667358221220a9fb4a459eb0d5a2f7bac5ea1961620a91781cad813ac55bac43ebeb8954cc9c64736f6c63430008190033",
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

// IntentOperatorsRegistrySignedAddedIterator is returned from FilterSignedAdded and is used to iterate over the raw logs and unpacked data for SignedAdded events raised by the IntentOperatorsRegistry contract.
type IntentOperatorsRegistrySignedAddedIterator struct {
	Event *IntentOperatorsRegistrySignedAdded // Event containing the contract specifics and raw log

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
func (it *IntentOperatorsRegistrySignedAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IntentOperatorsRegistrySignedAdded)
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
		it.Event = new(IntentOperatorsRegistrySignedAdded)
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
func (it *IntentOperatorsRegistrySignedAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IntentOperatorsRegistrySignedAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IntentOperatorsRegistrySignedAdded represents a SignedAdded event raised by the IntentOperatorsRegistry contract.
type IntentOperatorsRegistrySignedAdded struct {
	Publickey [32]byte
	Url       common.Hash
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSignedAdded is a free log retrieval operation binding the contract event 0xd993cade8966730f67b2b2c9d7689b17c3b727d61bf30f122bb3dcc5325fbec5.
//
// Solidity: event SignedAdded(bytes32 indexed publickey, string indexed url)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) FilterSignedAdded(opts *bind.FilterOpts, publickey [][32]byte, url []string) (*IntentOperatorsRegistrySignedAddedIterator, error) {

	var publickeyRule []interface{}
	for _, publickeyItem := range publickey {
		publickeyRule = append(publickeyRule, publickeyItem)
	}
	var urlRule []interface{}
	for _, urlItem := range url {
		urlRule = append(urlRule, urlItem)
	}

	logs, sub, err := _IntentOperatorsRegistry.contract.FilterLogs(opts, "SignedAdded", publickeyRule, urlRule)
	if err != nil {
		return nil, err
	}
	return &IntentOperatorsRegistrySignedAddedIterator{contract: _IntentOperatorsRegistry.contract, event: "SignedAdded", logs: logs, sub: sub}, nil
}

// WatchSignedAdded is a free log subscription operation binding the contract event 0xd993cade8966730f67b2b2c9d7689b17c3b727d61bf30f122bb3dcc5325fbec5.
//
// Solidity: event SignedAdded(bytes32 indexed publickey, string indexed url)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) WatchSignedAdded(opts *bind.WatchOpts, sink chan<- *IntentOperatorsRegistrySignedAdded, publickey [][32]byte, url []string) (event.Subscription, error) {

	var publickeyRule []interface{}
	for _, publickeyItem := range publickey {
		publickeyRule = append(publickeyRule, publickeyItem)
	}
	var urlRule []interface{}
	for _, urlItem := range url {
		urlRule = append(urlRule, urlItem)
	}

	logs, sub, err := _IntentOperatorsRegistry.contract.WatchLogs(opts, "SignedAdded", publickeyRule, urlRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IntentOperatorsRegistrySignedAdded)
				if err := _IntentOperatorsRegistry.contract.UnpackLog(event, "SignedAdded", log); err != nil {
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

// ParseSignedAdded is a log parse operation binding the contract event 0xd993cade8966730f67b2b2c9d7689b17c3b727d61bf30f122bb3dcc5325fbec5.
//
// Solidity: event SignedAdded(bytes32 indexed publickey, string indexed url)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) ParseSignedAdded(log types.Log) (*IntentOperatorsRegistrySignedAdded, error) {
	event := new(IntentOperatorsRegistrySignedAdded)
	if err := _IntentOperatorsRegistry.contract.UnpackLog(event, "SignedAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IntentOperatorsRegistrySignedRemovedIterator is returned from FilterSignedRemoved and is used to iterate over the raw logs and unpacked data for SignedRemoved events raised by the IntentOperatorsRegistry contract.
type IntentOperatorsRegistrySignedRemovedIterator struct {
	Event *IntentOperatorsRegistrySignedRemoved // Event containing the contract specifics and raw log

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
func (it *IntentOperatorsRegistrySignedRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IntentOperatorsRegistrySignedRemoved)
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
		it.Event = new(IntentOperatorsRegistrySignedRemoved)
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
func (it *IntentOperatorsRegistrySignedRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IntentOperatorsRegistrySignedRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IntentOperatorsRegistrySignedRemoved represents a SignedRemoved event raised by the IntentOperatorsRegistry contract.
type IntentOperatorsRegistrySignedRemoved struct {
	Publickey [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSignedRemoved is a free log retrieval operation binding the contract event 0x202ba177372096a533cb0be65537787905a2c9a9b25538d8d9f578706b412cb3.
//
// Solidity: event SignedRemoved(bytes32 indexed publickey)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) FilterSignedRemoved(opts *bind.FilterOpts, publickey [][32]byte) (*IntentOperatorsRegistrySignedRemovedIterator, error) {

	var publickeyRule []interface{}
	for _, publickeyItem := range publickey {
		publickeyRule = append(publickeyRule, publickeyItem)
	}

	logs, sub, err := _IntentOperatorsRegistry.contract.FilterLogs(opts, "SignedRemoved", publickeyRule)
	if err != nil {
		return nil, err
	}
	return &IntentOperatorsRegistrySignedRemovedIterator{contract: _IntentOperatorsRegistry.contract, event: "SignedRemoved", logs: logs, sub: sub}, nil
}

// WatchSignedRemoved is a free log subscription operation binding the contract event 0x202ba177372096a533cb0be65537787905a2c9a9b25538d8d9f578706b412cb3.
//
// Solidity: event SignedRemoved(bytes32 indexed publickey)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) WatchSignedRemoved(opts *bind.WatchOpts, sink chan<- *IntentOperatorsRegistrySignedRemoved, publickey [][32]byte) (event.Subscription, error) {

	var publickeyRule []interface{}
	for _, publickeyItem := range publickey {
		publickeyRule = append(publickeyRule, publickeyItem)
	}

	logs, sub, err := _IntentOperatorsRegistry.contract.WatchLogs(opts, "SignedRemoved", publickeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IntentOperatorsRegistrySignedRemoved)
				if err := _IntentOperatorsRegistry.contract.UnpackLog(event, "SignedRemoved", log); err != nil {
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

// ParseSignedRemoved is a log parse operation binding the contract event 0x202ba177372096a533cb0be65537787905a2c9a9b25538d8d9f578706b412cb3.
//
// Solidity: event SignedRemoved(bytes32 indexed publickey)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) ParseSignedRemoved(log types.Log) (*IntentOperatorsRegistrySignedRemoved, error) {
	event := new(IntentOperatorsRegistrySignedRemoved)
	if err := _IntentOperatorsRegistry.contract.UnpackLog(event, "SignedRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
