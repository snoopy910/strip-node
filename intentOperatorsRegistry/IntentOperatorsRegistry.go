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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maximumSigners\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"publickey\",\"type\":\"bytes32\"}],\"name\":\"SignedAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"publickey\",\"type\":\"bytes32\"}],\"name\":\"SignedRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"publickey\",\"type\":\"bytes32\"}],\"name\":\"SignerBlacklisted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"publickey\",\"type\":\"bytes32\"}],\"name\":\"SignerWhitelisted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_publickey\",\"type\":\"bytes32\"}],\"name\":\"addSigner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"publickey\",\"type\":\"bytes32\"}],\"name\":\"blacklistSigner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_publickey\",\"type\":\"bytes32\"}],\"name\":\"removeSigner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"signers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"exists\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"whitelisted\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"publickey\",\"type\":\"bytes32\"}],\"name\":\"whitelistSigner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a0604052348015600e575f80fd5b50604051610e94380380610e948339818101604052810190602e9190606d565b8060808181525050506093565b5f80fd5b5f819050919050565b604f81603f565b81146058575f80fd5b50565b5f815190506067816048565b92915050565b5f60208284031215607f57607e603b565b5b5f608a84828501605b565b91505092915050565b608051610dec6100a85f395f5050610dec5ff3fe608060405234801561000f575f80fd5b5060043610610091575f3560e01c80638da5cb5b116100645780638da5cb5b146100f657806394efc8ee14610114578063c46b824814610130578063eb49850c1461014c578063f2fde38b1461016857610091565b8063141774ef14610095578063715018a6146100c65780638129fc1c146100d05780638cc6f44c146100da575b5f80fd5b6100af60048036038101906100aa9190610a57565b610184565b6040516100bd929190610a9c565b60405180910390f35b6100ce6101bb565b005b6100d86101ce565b005b6100f460048036038101906100ef9190610a57565b61034e565b005b6100fe610471565b60405161010b9190610b02565b60405180910390f35b61012e60048036038101906101299190610a57565b6104a6565b005b61014a60048036038101906101459190610a57565b610513565b005b61016660048036038101906101619190610a57565b61068d565b005b610182600480360381019061017d9190610b45565b6106f9565b005b5f602052805f5260405f205f91509050805f015f9054906101000a900460ff1690805f0160019054906101000a900460ff16905082565b6101c361077d565b6101cc5f610804565b565b5f6101d76108d5565b90505f815f0160089054906101000a900460ff161590505f825f015f9054906101000a900467ffffffffffffffff1690505f808267ffffffffffffffff1614801561021f5750825b90505f60018367ffffffffffffffff1614801561025257505f3073ffffffffffffffffffffffffffffffffffffffff163b145b905081158015610260575080155b15610297576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001855f015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555083156102e4576001855f0160086101000a81548160ff0219169083151502179055505b6102ed336108fc565b8315610347575f855f0160086101000a81548160ff0219169083151502179055507fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2600160405161033e9190610bc5565b60405180910390a15b5050505050565b61035661077d565b5f801b810361039a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161039190610c38565b60405180910390fd5b600115155f808381526020019081526020015f205f015f9054906101000a900460ff161515146103ff576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103f690610ca0565b60405180910390fd5b5f808281526020019081526020015f205f8082015f6101000a81549060ff02191690555f820160016101000a81549060ff021916905550507f202ba177372096a533cb0be65537787905a2c9a9b25538d8d9f578706b412cb3816040516104669190610ccd565b60405180910390a150565b5f8061047b610910565b9050805f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1691505090565b6104ae61077d565b60015f808381526020019081526020015f205f0160016101000a81548160ff0219169083151502179055507f684d8290b28f7ee1d9799d0632bb71110e2f2c8feddb5493fb872b8b57faa927816040516105089190610ccd565b60405180910390a150565b61051b61077d565b5f801b810361055f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161055690610c38565b60405180910390fd5b5f15155f808381526020019081526020015f205f015f9054906101000a900460ff161515146105c3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105ba90610d30565b60405180910390fd5b600115155f808381526020019081526020015f205f0160019054906101000a900460ff16151514610629576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161062090610d98565b60405180910390fd5b60015f808381526020019081526020015f205f015f6101000a81548160ff0219169083151502179055507fb748e944a84031f3ffd5ffcc9af7992d594d84493ff33509de4614dc6ecd1dc9816040516106829190610ccd565b60405180910390a150565b61069561077d565b5f805f8381526020019081526020015f205f0160016101000a81548160ff0219169083151502179055507ffd0fd0ce237fc8c6c5ea5042cba831db42434ca670d70ead573412793ad2b48c816040516106ee9190610ccd565b60405180910390a150565b61070161077d565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610771575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016107689190610b02565b60405180910390fd5b61077a81610804565b50565b610785610937565b73ffffffffffffffffffffffffffffffffffffffff166107a3610471565b73ffffffffffffffffffffffffffffffffffffffff1614610802576107c6610937565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016107f99190610b02565b60405180910390fd5b565b5f61080d610910565b90505f815f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905082825f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508273ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3505050565b5f7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00905090565b61090461093e565b61090d8161097e565b50565b5f7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300905090565b5f33905090565b610946610a02565b61097c576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b61098661093e565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036109f6575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016109ed9190610b02565b60405180910390fd5b6109ff81610804565b50565b5f610a0b6108d5565b5f0160089054906101000a900460ff16905090565b5f80fd5b5f819050919050565b610a3681610a24565b8114610a40575f80fd5b50565b5f81359050610a5181610a2d565b92915050565b5f60208284031215610a6c57610a6b610a20565b5b5f610a7984828501610a43565b91505092915050565b5f8115159050919050565b610a9681610a82565b82525050565b5f604082019050610aaf5f830185610a8d565b610abc6020830184610a8d565b9392505050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610aec82610ac3565b9050919050565b610afc81610ae2565b82525050565b5f602082019050610b155f830184610af3565b92915050565b610b2481610ae2565b8114610b2e575f80fd5b50565b5f81359050610b3f81610b1b565b92915050565b5f60208284031215610b5a57610b59610a20565b5b5f610b6784828501610b31565b91505092915050565b5f819050919050565b5f67ffffffffffffffff82169050919050565b5f819050919050565b5f610baf610baa610ba584610b70565b610b8c565b610b79565b9050919050565b610bbf81610b95565b82525050565b5f602082019050610bd85f830184610bb6565b92915050565b5f82825260208201905092915050565b7f7075626c6963206b657920697320656d707479000000000000000000000000005f82015250565b5f610c22601383610bde565b9150610c2d82610bee565b602082019050919050565b5f6020820190508181035f830152610c4f81610c16565b9050919050565b7f7075626c6963206b6579206973206e6f742072656769737465726564000000005f82015250565b5f610c8a601c83610bde565b9150610c9582610c56565b602082019050919050565b5f6020820190508181035f830152610cb781610c7e565b9050919050565b610cc781610a24565b82525050565b5f602082019050610ce05f830184610cbe565b92915050565b7f7075626c6963206b657920697320616c726561647920726567697374657265645f82015250565b5f610d1a602083610bde565b9150610d2582610ce6565b602082019050919050565b5f6020820190508181035f830152610d4781610d0e565b9050919050565b7f7369676e6572206973206e6f742077686974656c6973746564000000000000005f82015250565b5f610d82601983610bde565b9150610d8d82610d4e565b602082019050919050565b5f6020820190508181035f830152610daf81610d76565b905091905056fea26469706673582212208052032b40f8474931aa900a4903c6dc9ac1a452cbe9585c68329a0a4db3e9e064736f6c63430008190033",
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
// Solidity: function signers(bytes32 ) view returns(bool exists, bool whitelisted)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryCaller) Signers(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Exists      bool
	Whitelisted bool
}, error) {
	var out []interface{}
	err := _IntentOperatorsRegistry.contract.Call(opts, &out, "signers", arg0)

	outstruct := new(struct {
		Exists      bool
		Whitelisted bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Exists = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Whitelisted = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// Signers is a free data retrieval call binding the contract method 0x141774ef.
//
// Solidity: function signers(bytes32 ) view returns(bool exists, bool whitelisted)
func (_IntentOperatorsRegistry *IntentOperatorsRegistrySession) Signers(arg0 [32]byte) (struct {
	Exists      bool
	Whitelisted bool
}, error) {
	return _IntentOperatorsRegistry.Contract.Signers(&_IntentOperatorsRegistry.CallOpts, arg0)
}

// Signers is a free data retrieval call binding the contract method 0x141774ef.
//
// Solidity: function signers(bytes32 ) view returns(bool exists, bool whitelisted)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryCallerSession) Signers(arg0 [32]byte) (struct {
	Exists      bool
	Whitelisted bool
}, error) {
	return _IntentOperatorsRegistry.Contract.Signers(&_IntentOperatorsRegistry.CallOpts, arg0)
}

// AddSigner is a paid mutator transaction binding the contract method 0xc46b8248.
//
// Solidity: function addSigner(bytes32 _publickey) returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistryTransactor) AddSigner(opts *bind.TransactOpts, _publickey [32]byte) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.contract.Transact(opts, "addSigner", _publickey)
}

// AddSigner is a paid mutator transaction binding the contract method 0xc46b8248.
//
// Solidity: function addSigner(bytes32 _publickey) returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistrySession) AddSigner(_publickey [32]byte) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.AddSigner(&_IntentOperatorsRegistry.TransactOpts, _publickey)
}

// AddSigner is a paid mutator transaction binding the contract method 0xc46b8248.
//
// Solidity: function addSigner(bytes32 _publickey) returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistryTransactorSession) AddSigner(_publickey [32]byte) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.AddSigner(&_IntentOperatorsRegistry.TransactOpts, _publickey)
}

// BlacklistSigner is a paid mutator transaction binding the contract method 0xeb49850c.
//
// Solidity: function blacklistSigner(bytes32 publickey) returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistryTransactor) BlacklistSigner(opts *bind.TransactOpts, publickey [32]byte) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.contract.Transact(opts, "blacklistSigner", publickey)
}

// BlacklistSigner is a paid mutator transaction binding the contract method 0xeb49850c.
//
// Solidity: function blacklistSigner(bytes32 publickey) returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistrySession) BlacklistSigner(publickey [32]byte) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.BlacklistSigner(&_IntentOperatorsRegistry.TransactOpts, publickey)
}

// BlacklistSigner is a paid mutator transaction binding the contract method 0xeb49850c.
//
// Solidity: function blacklistSigner(bytes32 publickey) returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistryTransactorSession) BlacklistSigner(publickey [32]byte) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.BlacklistSigner(&_IntentOperatorsRegistry.TransactOpts, publickey)
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

// WhitelistSigner is a paid mutator transaction binding the contract method 0x94efc8ee.
//
// Solidity: function whitelistSigner(bytes32 publickey) returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistryTransactor) WhitelistSigner(opts *bind.TransactOpts, publickey [32]byte) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.contract.Transact(opts, "whitelistSigner", publickey)
}

// WhitelistSigner is a paid mutator transaction binding the contract method 0x94efc8ee.
//
// Solidity: function whitelistSigner(bytes32 publickey) returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistrySession) WhitelistSigner(publickey [32]byte) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.WhitelistSigner(&_IntentOperatorsRegistry.TransactOpts, publickey)
}

// WhitelistSigner is a paid mutator transaction binding the contract method 0x94efc8ee.
//
// Solidity: function whitelistSigner(bytes32 publickey) returns()
func (_IntentOperatorsRegistry *IntentOperatorsRegistryTransactorSession) WhitelistSigner(publickey [32]byte) (*types.Transaction, error) {
	return _IntentOperatorsRegistry.Contract.WhitelistSigner(&_IntentOperatorsRegistry.TransactOpts, publickey)
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
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSignedAdded is a free log retrieval operation binding the contract event 0xb748e944a84031f3ffd5ffcc9af7992d594d84493ff33509de4614dc6ecd1dc9.
//
// Solidity: event SignedAdded(bytes32 publickey)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) FilterSignedAdded(opts *bind.FilterOpts) (*IntentOperatorsRegistrySignedAddedIterator, error) {

	logs, sub, err := _IntentOperatorsRegistry.contract.FilterLogs(opts, "SignedAdded")
	if err != nil {
		return nil, err
	}
	return &IntentOperatorsRegistrySignedAddedIterator{contract: _IntentOperatorsRegistry.contract, event: "SignedAdded", logs: logs, sub: sub}, nil
}

// WatchSignedAdded is a free log subscription operation binding the contract event 0xb748e944a84031f3ffd5ffcc9af7992d594d84493ff33509de4614dc6ecd1dc9.
//
// Solidity: event SignedAdded(bytes32 publickey)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) WatchSignedAdded(opts *bind.WatchOpts, sink chan<- *IntentOperatorsRegistrySignedAdded) (event.Subscription, error) {

	logs, sub, err := _IntentOperatorsRegistry.contract.WatchLogs(opts, "SignedAdded")
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

// ParseSignedAdded is a log parse operation binding the contract event 0xb748e944a84031f3ffd5ffcc9af7992d594d84493ff33509de4614dc6ecd1dc9.
//
// Solidity: event SignedAdded(bytes32 publickey)
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
// Solidity: event SignedRemoved(bytes32 publickey)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) FilterSignedRemoved(opts *bind.FilterOpts) (*IntentOperatorsRegistrySignedRemovedIterator, error) {

	logs, sub, err := _IntentOperatorsRegistry.contract.FilterLogs(opts, "SignedRemoved")
	if err != nil {
		return nil, err
	}
	return &IntentOperatorsRegistrySignedRemovedIterator{contract: _IntentOperatorsRegistry.contract, event: "SignedRemoved", logs: logs, sub: sub}, nil
}

// WatchSignedRemoved is a free log subscription operation binding the contract event 0x202ba177372096a533cb0be65537787905a2c9a9b25538d8d9f578706b412cb3.
//
// Solidity: event SignedRemoved(bytes32 publickey)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) WatchSignedRemoved(opts *bind.WatchOpts, sink chan<- *IntentOperatorsRegistrySignedRemoved) (event.Subscription, error) {

	logs, sub, err := _IntentOperatorsRegistry.contract.WatchLogs(opts, "SignedRemoved")
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
// Solidity: event SignedRemoved(bytes32 publickey)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) ParseSignedRemoved(log types.Log) (*IntentOperatorsRegistrySignedRemoved, error) {
	event := new(IntentOperatorsRegistrySignedRemoved)
	if err := _IntentOperatorsRegistry.contract.UnpackLog(event, "SignedRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IntentOperatorsRegistrySignerBlacklistedIterator is returned from FilterSignerBlacklisted and is used to iterate over the raw logs and unpacked data for SignerBlacklisted events raised by the IntentOperatorsRegistry contract.
type IntentOperatorsRegistrySignerBlacklistedIterator struct {
	Event *IntentOperatorsRegistrySignerBlacklisted // Event containing the contract specifics and raw log

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
func (it *IntentOperatorsRegistrySignerBlacklistedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IntentOperatorsRegistrySignerBlacklisted)
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
		it.Event = new(IntentOperatorsRegistrySignerBlacklisted)
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
func (it *IntentOperatorsRegistrySignerBlacklistedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IntentOperatorsRegistrySignerBlacklistedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IntentOperatorsRegistrySignerBlacklisted represents a SignerBlacklisted event raised by the IntentOperatorsRegistry contract.
type IntentOperatorsRegistrySignerBlacklisted struct {
	Publickey [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSignerBlacklisted is a free log retrieval operation binding the contract event 0xfd0fd0ce237fc8c6c5ea5042cba831db42434ca670d70ead573412793ad2b48c.
//
// Solidity: event SignerBlacklisted(bytes32 publickey)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) FilterSignerBlacklisted(opts *bind.FilterOpts) (*IntentOperatorsRegistrySignerBlacklistedIterator, error) {

	logs, sub, err := _IntentOperatorsRegistry.contract.FilterLogs(opts, "SignerBlacklisted")
	if err != nil {
		return nil, err
	}
	return &IntentOperatorsRegistrySignerBlacklistedIterator{contract: _IntentOperatorsRegistry.contract, event: "SignerBlacklisted", logs: logs, sub: sub}, nil
}

// WatchSignerBlacklisted is a free log subscription operation binding the contract event 0xfd0fd0ce237fc8c6c5ea5042cba831db42434ca670d70ead573412793ad2b48c.
//
// Solidity: event SignerBlacklisted(bytes32 publickey)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) WatchSignerBlacklisted(opts *bind.WatchOpts, sink chan<- *IntentOperatorsRegistrySignerBlacklisted) (event.Subscription, error) {

	logs, sub, err := _IntentOperatorsRegistry.contract.WatchLogs(opts, "SignerBlacklisted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IntentOperatorsRegistrySignerBlacklisted)
				if err := _IntentOperatorsRegistry.contract.UnpackLog(event, "SignerBlacklisted", log); err != nil {
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

// ParseSignerBlacklisted is a log parse operation binding the contract event 0xfd0fd0ce237fc8c6c5ea5042cba831db42434ca670d70ead573412793ad2b48c.
//
// Solidity: event SignerBlacklisted(bytes32 publickey)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) ParseSignerBlacklisted(log types.Log) (*IntentOperatorsRegistrySignerBlacklisted, error) {
	event := new(IntentOperatorsRegistrySignerBlacklisted)
	if err := _IntentOperatorsRegistry.contract.UnpackLog(event, "SignerBlacklisted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IntentOperatorsRegistrySignerWhitelistedIterator is returned from FilterSignerWhitelisted and is used to iterate over the raw logs and unpacked data for SignerWhitelisted events raised by the IntentOperatorsRegistry contract.
type IntentOperatorsRegistrySignerWhitelistedIterator struct {
	Event *IntentOperatorsRegistrySignerWhitelisted // Event containing the contract specifics and raw log

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
func (it *IntentOperatorsRegistrySignerWhitelistedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IntentOperatorsRegistrySignerWhitelisted)
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
		it.Event = new(IntentOperatorsRegistrySignerWhitelisted)
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
func (it *IntentOperatorsRegistrySignerWhitelistedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IntentOperatorsRegistrySignerWhitelistedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IntentOperatorsRegistrySignerWhitelisted represents a SignerWhitelisted event raised by the IntentOperatorsRegistry contract.
type IntentOperatorsRegistrySignerWhitelisted struct {
	Publickey [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSignerWhitelisted is a free log retrieval operation binding the contract event 0x684d8290b28f7ee1d9799d0632bb71110e2f2c8feddb5493fb872b8b57faa927.
//
// Solidity: event SignerWhitelisted(bytes32 publickey)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) FilterSignerWhitelisted(opts *bind.FilterOpts) (*IntentOperatorsRegistrySignerWhitelistedIterator, error) {

	logs, sub, err := _IntentOperatorsRegistry.contract.FilterLogs(opts, "SignerWhitelisted")
	if err != nil {
		return nil, err
	}
	return &IntentOperatorsRegistrySignerWhitelistedIterator{contract: _IntentOperatorsRegistry.contract, event: "SignerWhitelisted", logs: logs, sub: sub}, nil
}

// WatchSignerWhitelisted is a free log subscription operation binding the contract event 0x684d8290b28f7ee1d9799d0632bb71110e2f2c8feddb5493fb872b8b57faa927.
//
// Solidity: event SignerWhitelisted(bytes32 publickey)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) WatchSignerWhitelisted(opts *bind.WatchOpts, sink chan<- *IntentOperatorsRegistrySignerWhitelisted) (event.Subscription, error) {

	logs, sub, err := _IntentOperatorsRegistry.contract.WatchLogs(opts, "SignerWhitelisted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IntentOperatorsRegistrySignerWhitelisted)
				if err := _IntentOperatorsRegistry.contract.UnpackLog(event, "SignerWhitelisted", log); err != nil {
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

// ParseSignerWhitelisted is a log parse operation binding the contract event 0x684d8290b28f7ee1d9799d0632bb71110e2f2c8feddb5493fb872b8b57faa927.
//
// Solidity: event SignerWhitelisted(bytes32 publickey)
func (_IntentOperatorsRegistry *IntentOperatorsRegistryFilterer) ParseSignerWhitelisted(log types.Log) (*IntentOperatorsRegistrySignerWhitelisted, error) {
	event := new(IntentOperatorsRegistrySignerWhitelisted)
	if err := _IntentOperatorsRegistry.contract.UnpackLog(event, "SignerWhitelisted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
