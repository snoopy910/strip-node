// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package signerhub

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

// SignerhubMetaData contains all meta data concerning the Signerhub contract.
var SignerhubMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_relayerRegistry\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"publickey\",\"type\":\"bytes32\"}],\"name\":\"SignedAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"publickey\",\"type\":\"bytes32\"}],\"name\":\"SignedRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"publickey\",\"type\":\"bytes32\"}],\"name\":\"SignerBlacklisted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"publickey\",\"type\":\"bytes32\"}],\"name\":\"SignerWhitelisted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_publickey\",\"type\":\"bytes32\"}],\"name\":\"addSigner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"publickey\",\"type\":\"bytes32\"}],\"name\":\"blacklistSigner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"publickeys\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_publickey\",\"type\":\"bytes32\"}],\"name\":\"removeSigner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"signers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"exists\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"whitelisted\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"startKey\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"publickey\",\"type\":\"bytes32\"}],\"name\":\"whitelistSigner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610b86380380610b8683398101604081905261002f91610040565b6001600160a01b0316608052610070565b60006020828403121561005257600080fd5b81516001600160a01b038116811461006957600080fd5b9392505050565b608051610af46100926000396000818161033801526105990152610af46000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c806394efc8ee1161008c578063f2327b8511610066578063f2327b85146101b9578063f2fde38b146101d9578063f910cff8146101ec578063fc7e9c6f146101f557600080fd5b806394efc8ee14610180578063c46b824814610193578063eb49850c146101a657600080fd5b8063141774ef146100d45780635738bdd014610129578063715018a6146101405780638129fc1c1461014a5780638cc6f44c146101525780638da5cb5b14610165575b600080fd5b6101076100e2366004610993565b6068602052600090815260409020805460019091015460ff8082169161010090041683565b6040805193845291151560208401521515908201526060015b60405180910390f35b61013260695481565b604051908152602001610120565b6101486101fe565b005b610148610212565b610148610160366004610993565b61032d565b6033546040516001600160a01b039091168152602001610120565b61014861018e366004610993565b610536565b6101486101a1366004610993565b61058e565b6101486101b4366004610993565b610784565b6101326101c7366004610993565b60676020526000908152604090205481565b6101486101e73660046109ac565b6107d8565b61013260665481565b61013260655481565b61020661084e565b61021060006108a8565b565b600054610100900460ff16158080156102325750600054600160ff909116105b8061024c5750303b15801561024c575060005460ff166001145b6102b45760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084015b60405180910390fd5b6000805460ff1916600117905580156102d7576000805461ff0019166101001790555b426069556102e36108fa565b801561032a576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498906020015b60405180910390a15b50565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016148061036e57506033546001600160a01b031633145b6103b15760405162461bcd60e51b81526020600482015260146024820152734e6f742052656c6179657220526567697374727960601b60448201526064016102ab565b600081815260686020526040902060019081015460ff161515146104175760405162461bcd60e51b815260206004820152601c60248201527f7075626c6963206b6579206973206e6f7420726567697374657265640000000060448201526064016102ab565b600081815260686020526040812060018101805460ff191690558054919055805b600160655461044791906109f2565b81116104d7576067600061045c836001610a09565b81526020019081526020016000205460676000838152602001908152602001600020819055506000606760008360016104959190610a09565b815260208082019290925260409081016000908120939093558383526067825280832054835260689091529020819055806104cf81610a21565b915050610438565b50606580549060006104e883610a3a565b91905055506104f8606554610929565b606655426069556040518281527f202ba177372096a533cb0be65537787905a2c9a9b25538d8d9f578706b412cb39060200160405180910390a15050565b61053e61084e565b60008181526068602052604090819020600101805461ff001916610100179055517f684d8290b28f7ee1d9799d0632bb71110e2f2c8feddb5493fb872b8b57faa927906103219083815260200190565b336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614806105cf57506033546001600160a01b031633145b6106125760405162461bcd60e51b81526020600482015260146024820152734e6f742052656c6179657220526567697374727960601b60448201526064016102ab565b60008181526068602052604090206001015460ff16156106745760405162461bcd60e51b815260206004820181905260248201527f7075626c6963206b657920697320616c7265616479207265676973746572656460448201526064016102ab565b6000818152606860205260409020600190810154610100900460ff161515146106df5760405162461bcd60e51b815260206004820152601960248201527f7369676e6572206973206e6f742077686974656c69737465640000000000000060448201526064016102ab565b60008181526068602081815260408084206001818101805460ff191690911790556065805486526067845291852086905581548686529390925290829055909161072883610a21565b9190505550610738606554610929565b606655426069556065547f0e4e5b3ba228cadf0c65d31f94d29259d1b37003b9dccdd6b60b23e89804f8f190610770906001906109f2565b604080519182526020820184905201610321565b61078c61084e565b60008181526068602052604090819020600101805461ff0019169055517ffd0fd0ce237fc8c6c5ea5042cba831db42434ca670d70ead573412793ad2b48c906103219083815260200190565b6107e061084e565b6001600160a01b0381166108455760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016102ab565b61032a816108a8565b6033546001600160a01b031633146102105760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016102ab565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166109215760405162461bcd60e51b81526004016102ab90610a51565b610210610963565b6000816001148061093a5750816002145b1561094757506001919050565b610952600283610a9c565b61095d906001610a09565b92915050565b600054610100900460ff1661098a5760405162461bcd60e51b81526004016102ab90610a51565b610210336108a8565b6000602082840312156109a557600080fd5b5035919050565b6000602082840312156109be57600080fd5b81356001600160a01b03811681146109d557600080fd5b9392505050565b634e487b7160e01b600052601160045260246000fd5b600082821015610a0457610a046109dc565b500390565b60008219821115610a1c57610a1c6109dc565b500190565b600060018201610a3357610a336109dc565b5060010190565b600081610a4957610a496109dc565b506000190190565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b600082610ab957634e487b7160e01b600052601260045260246000fd5b50049056fea2646970667358221220d493394c43069cfa8268da4f7f36dc7cdeeb1d61b3397d989de330a806d06eb364736f6c634300080d0033",
}

// SignerhubABI is the input ABI used to generate the binding from.
// Deprecated: Use SignerhubMetaData.ABI instead.
var SignerhubABI = SignerhubMetaData.ABI

// SignerhubBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SignerhubMetaData.Bin instead.
var SignerhubBin = SignerhubMetaData.Bin

// DeploySignerhub deploys a new Ethereum contract, binding an instance of Signerhub to it.
func DeploySignerhub(auth *bind.TransactOpts, backend bind.ContractBackend, _relayerRegistry common.Address) (common.Address, *types.Transaction, *Signerhub, error) {
	parsed, err := SignerhubMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SignerhubBin), backend, _relayerRegistry)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Signerhub{SignerhubCaller: SignerhubCaller{contract: contract}, SignerhubTransactor: SignerhubTransactor{contract: contract}, SignerhubFilterer: SignerhubFilterer{contract: contract}}, nil
}

// Signerhub is an auto generated Go binding around an Ethereum contract.
type Signerhub struct {
	SignerhubCaller     // Read-only binding to the contract
	SignerhubTransactor // Write-only binding to the contract
	SignerhubFilterer   // Log filterer for contract events
}

// SignerhubCaller is an auto generated read-only Go binding around an Ethereum contract.
type SignerhubCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignerhubTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SignerhubTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignerhubFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SignerhubFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignerhubSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SignerhubSession struct {
	Contract     *Signerhub        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SignerhubCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SignerhubCallerSession struct {
	Contract *SignerhubCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// SignerhubTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SignerhubTransactorSession struct {
	Contract     *SignerhubTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// SignerhubRaw is an auto generated low-level Go binding around an Ethereum contract.
type SignerhubRaw struct {
	Contract *Signerhub // Generic contract binding to access the raw methods on
}

// SignerhubCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SignerhubCallerRaw struct {
	Contract *SignerhubCaller // Generic read-only contract binding to access the raw methods on
}

// SignerhubTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SignerhubTransactorRaw struct {
	Contract *SignerhubTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSignerhub creates a new instance of Signerhub, bound to a specific deployed contract.
func NewSignerhub(address common.Address, backend bind.ContractBackend) (*Signerhub, error) {
	contract, err := bindSignerhub(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Signerhub{SignerhubCaller: SignerhubCaller{contract: contract}, SignerhubTransactor: SignerhubTransactor{contract: contract}, SignerhubFilterer: SignerhubFilterer{contract: contract}}, nil
}

// NewSignerhubCaller creates a new read-only instance of Signerhub, bound to a specific deployed contract.
func NewSignerhubCaller(address common.Address, caller bind.ContractCaller) (*SignerhubCaller, error) {
	contract, err := bindSignerhub(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SignerhubCaller{contract: contract}, nil
}

// NewSignerhubTransactor creates a new write-only instance of Signerhub, bound to a specific deployed contract.
func NewSignerhubTransactor(address common.Address, transactor bind.ContractTransactor) (*SignerhubTransactor, error) {
	contract, err := bindSignerhub(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SignerhubTransactor{contract: contract}, nil
}

// NewSignerhubFilterer creates a new log filterer instance of Signerhub, bound to a specific deployed contract.
func NewSignerhubFilterer(address common.Address, filterer bind.ContractFilterer) (*SignerhubFilterer, error) {
	contract, err := bindSignerhub(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SignerhubFilterer{contract: contract}, nil
}

// bindSignerhub binds a generic wrapper to an already deployed contract.
func bindSignerhub(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SignerhubMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Signerhub *SignerhubRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Signerhub.Contract.SignerhubCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Signerhub *SignerhubRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Signerhub.Contract.SignerhubTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Signerhub *SignerhubRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Signerhub.Contract.SignerhubTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Signerhub *SignerhubCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Signerhub.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Signerhub *SignerhubTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Signerhub.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Signerhub *SignerhubTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Signerhub.Contract.contract.Transact(opts, method, params...)
}

// CurrentThreshold is a free data retrieval call binding the contract method 0xf910cff8.
//
// Solidity: function currentThreshold() view returns(uint256)
func (_Signerhub *SignerhubCaller) CurrentThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Signerhub.contract.Call(opts, &out, "currentThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentThreshold is a free data retrieval call binding the contract method 0xf910cff8.
//
// Solidity: function currentThreshold() view returns(uint256)
func (_Signerhub *SignerhubSession) CurrentThreshold() (*big.Int, error) {
	return _Signerhub.Contract.CurrentThreshold(&_Signerhub.CallOpts)
}

// CurrentThreshold is a free data retrieval call binding the contract method 0xf910cff8.
//
// Solidity: function currentThreshold() view returns(uint256)
func (_Signerhub *SignerhubCallerSession) CurrentThreshold() (*big.Int, error) {
	return _Signerhub.Contract.CurrentThreshold(&_Signerhub.CallOpts)
}

// NextIndex is a free data retrieval call binding the contract method 0xfc7e9c6f.
//
// Solidity: function nextIndex() view returns(uint256)
func (_Signerhub *SignerhubCaller) NextIndex(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Signerhub.contract.Call(opts, &out, "nextIndex")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextIndex is a free data retrieval call binding the contract method 0xfc7e9c6f.
//
// Solidity: function nextIndex() view returns(uint256)
func (_Signerhub *SignerhubSession) NextIndex() (*big.Int, error) {
	return _Signerhub.Contract.NextIndex(&_Signerhub.CallOpts)
}

// NextIndex is a free data retrieval call binding the contract method 0xfc7e9c6f.
//
// Solidity: function nextIndex() view returns(uint256)
func (_Signerhub *SignerhubCallerSession) NextIndex() (*big.Int, error) {
	return _Signerhub.Contract.NextIndex(&_Signerhub.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Signerhub *SignerhubCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Signerhub.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Signerhub *SignerhubSession) Owner() (common.Address, error) {
	return _Signerhub.Contract.Owner(&_Signerhub.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Signerhub *SignerhubCallerSession) Owner() (common.Address, error) {
	return _Signerhub.Contract.Owner(&_Signerhub.CallOpts)
}

// Publickeys is a free data retrieval call binding the contract method 0xf2327b85.
//
// Solidity: function publickeys(uint256 ) view returns(bytes32)
func (_Signerhub *SignerhubCaller) Publickeys(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Signerhub.contract.Call(opts, &out, "publickeys", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Publickeys is a free data retrieval call binding the contract method 0xf2327b85.
//
// Solidity: function publickeys(uint256 ) view returns(bytes32)
func (_Signerhub *SignerhubSession) Publickeys(arg0 *big.Int) ([32]byte, error) {
	return _Signerhub.Contract.Publickeys(&_Signerhub.CallOpts, arg0)
}

// Publickeys is a free data retrieval call binding the contract method 0xf2327b85.
//
// Solidity: function publickeys(uint256 ) view returns(bytes32)
func (_Signerhub *SignerhubCallerSession) Publickeys(arg0 *big.Int) ([32]byte, error) {
	return _Signerhub.Contract.Publickeys(&_Signerhub.CallOpts, arg0)
}

// Signers is a free data retrieval call binding the contract method 0x141774ef.
//
// Solidity: function signers(bytes32 ) view returns(uint256 index, bool exists, bool whitelisted)
func (_Signerhub *SignerhubCaller) Signers(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Index       *big.Int
	Exists      bool
	Whitelisted bool
}, error) {
	var out []interface{}
	err := _Signerhub.contract.Call(opts, &out, "signers", arg0)

	outstruct := new(struct {
		Index       *big.Int
		Exists      bool
		Whitelisted bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Index = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Exists = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.Whitelisted = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err

}

// Signers is a free data retrieval call binding the contract method 0x141774ef.
//
// Solidity: function signers(bytes32 ) view returns(uint256 index, bool exists, bool whitelisted)
func (_Signerhub *SignerhubSession) Signers(arg0 [32]byte) (struct {
	Index       *big.Int
	Exists      bool
	Whitelisted bool
}, error) {
	return _Signerhub.Contract.Signers(&_Signerhub.CallOpts, arg0)
}

// Signers is a free data retrieval call binding the contract method 0x141774ef.
//
// Solidity: function signers(bytes32 ) view returns(uint256 index, bool exists, bool whitelisted)
func (_Signerhub *SignerhubCallerSession) Signers(arg0 [32]byte) (struct {
	Index       *big.Int
	Exists      bool
	Whitelisted bool
}, error) {
	return _Signerhub.Contract.Signers(&_Signerhub.CallOpts, arg0)
}

// StartKey is a free data retrieval call binding the contract method 0x5738bdd0.
//
// Solidity: function startKey() view returns(uint256)
func (_Signerhub *SignerhubCaller) StartKey(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Signerhub.contract.Call(opts, &out, "startKey")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StartKey is a free data retrieval call binding the contract method 0x5738bdd0.
//
// Solidity: function startKey() view returns(uint256)
func (_Signerhub *SignerhubSession) StartKey() (*big.Int, error) {
	return _Signerhub.Contract.StartKey(&_Signerhub.CallOpts)
}

// StartKey is a free data retrieval call binding the contract method 0x5738bdd0.
//
// Solidity: function startKey() view returns(uint256)
func (_Signerhub *SignerhubCallerSession) StartKey() (*big.Int, error) {
	return _Signerhub.Contract.StartKey(&_Signerhub.CallOpts)
}

// AddSigner is a paid mutator transaction binding the contract method 0xc46b8248.
//
// Solidity: function addSigner(bytes32 _publickey) returns()
func (_Signerhub *SignerhubTransactor) AddSigner(opts *bind.TransactOpts, _publickey [32]byte) (*types.Transaction, error) {
	return _Signerhub.contract.Transact(opts, "addSigner", _publickey)
}

// AddSigner is a paid mutator transaction binding the contract method 0xc46b8248.
//
// Solidity: function addSigner(bytes32 _publickey) returns()
func (_Signerhub *SignerhubSession) AddSigner(_publickey [32]byte) (*types.Transaction, error) {
	return _Signerhub.Contract.AddSigner(&_Signerhub.TransactOpts, _publickey)
}

// AddSigner is a paid mutator transaction binding the contract method 0xc46b8248.
//
// Solidity: function addSigner(bytes32 _publickey) returns()
func (_Signerhub *SignerhubTransactorSession) AddSigner(_publickey [32]byte) (*types.Transaction, error) {
	return _Signerhub.Contract.AddSigner(&_Signerhub.TransactOpts, _publickey)
}

// BlacklistSigner is a paid mutator transaction binding the contract method 0xeb49850c.
//
// Solidity: function blacklistSigner(bytes32 publickey) returns()
func (_Signerhub *SignerhubTransactor) BlacklistSigner(opts *bind.TransactOpts, publickey [32]byte) (*types.Transaction, error) {
	return _Signerhub.contract.Transact(opts, "blacklistSigner", publickey)
}

// BlacklistSigner is a paid mutator transaction binding the contract method 0xeb49850c.
//
// Solidity: function blacklistSigner(bytes32 publickey) returns()
func (_Signerhub *SignerhubSession) BlacklistSigner(publickey [32]byte) (*types.Transaction, error) {
	return _Signerhub.Contract.BlacklistSigner(&_Signerhub.TransactOpts, publickey)
}

// BlacklistSigner is a paid mutator transaction binding the contract method 0xeb49850c.
//
// Solidity: function blacklistSigner(bytes32 publickey) returns()
func (_Signerhub *SignerhubTransactorSession) BlacklistSigner(publickey [32]byte) (*types.Transaction, error) {
	return _Signerhub.Contract.BlacklistSigner(&_Signerhub.TransactOpts, publickey)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Signerhub *SignerhubTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Signerhub.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Signerhub *SignerhubSession) Initialize() (*types.Transaction, error) {
	return _Signerhub.Contract.Initialize(&_Signerhub.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Signerhub *SignerhubTransactorSession) Initialize() (*types.Transaction, error) {
	return _Signerhub.Contract.Initialize(&_Signerhub.TransactOpts)
}

// RemoveSigner is a paid mutator transaction binding the contract method 0x8cc6f44c.
//
// Solidity: function removeSigner(bytes32 _publickey) returns()
func (_Signerhub *SignerhubTransactor) RemoveSigner(opts *bind.TransactOpts, _publickey [32]byte) (*types.Transaction, error) {
	return _Signerhub.contract.Transact(opts, "removeSigner", _publickey)
}

// RemoveSigner is a paid mutator transaction binding the contract method 0x8cc6f44c.
//
// Solidity: function removeSigner(bytes32 _publickey) returns()
func (_Signerhub *SignerhubSession) RemoveSigner(_publickey [32]byte) (*types.Transaction, error) {
	return _Signerhub.Contract.RemoveSigner(&_Signerhub.TransactOpts, _publickey)
}

// RemoveSigner is a paid mutator transaction binding the contract method 0x8cc6f44c.
//
// Solidity: function removeSigner(bytes32 _publickey) returns()
func (_Signerhub *SignerhubTransactorSession) RemoveSigner(_publickey [32]byte) (*types.Transaction, error) {
	return _Signerhub.Contract.RemoveSigner(&_Signerhub.TransactOpts, _publickey)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Signerhub *SignerhubTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Signerhub.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Signerhub *SignerhubSession) RenounceOwnership() (*types.Transaction, error) {
	return _Signerhub.Contract.RenounceOwnership(&_Signerhub.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Signerhub *SignerhubTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Signerhub.Contract.RenounceOwnership(&_Signerhub.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Signerhub *SignerhubTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Signerhub.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Signerhub *SignerhubSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Signerhub.Contract.TransferOwnership(&_Signerhub.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Signerhub *SignerhubTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Signerhub.Contract.TransferOwnership(&_Signerhub.TransactOpts, newOwner)
}

// WhitelistSigner is a paid mutator transaction binding the contract method 0x94efc8ee.
//
// Solidity: function whitelistSigner(bytes32 publickey) returns()
func (_Signerhub *SignerhubTransactor) WhitelistSigner(opts *bind.TransactOpts, publickey [32]byte) (*types.Transaction, error) {
	return _Signerhub.contract.Transact(opts, "whitelistSigner", publickey)
}

// WhitelistSigner is a paid mutator transaction binding the contract method 0x94efc8ee.
//
// Solidity: function whitelistSigner(bytes32 publickey) returns()
func (_Signerhub *SignerhubSession) WhitelistSigner(publickey [32]byte) (*types.Transaction, error) {
	return _Signerhub.Contract.WhitelistSigner(&_Signerhub.TransactOpts, publickey)
}

// WhitelistSigner is a paid mutator transaction binding the contract method 0x94efc8ee.
//
// Solidity: function whitelistSigner(bytes32 publickey) returns()
func (_Signerhub *SignerhubTransactorSession) WhitelistSigner(publickey [32]byte) (*types.Transaction, error) {
	return _Signerhub.Contract.WhitelistSigner(&_Signerhub.TransactOpts, publickey)
}

// SignerhubInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Signerhub contract.
type SignerhubInitializedIterator struct {
	Event *SignerhubInitialized // Event containing the contract specifics and raw log

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
func (it *SignerhubInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SignerhubInitialized)
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
		it.Event = new(SignerhubInitialized)
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
func (it *SignerhubInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SignerhubInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SignerhubInitialized represents a Initialized event raised by the Signerhub contract.
type SignerhubInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Signerhub *SignerhubFilterer) FilterInitialized(opts *bind.FilterOpts) (*SignerhubInitializedIterator, error) {

	logs, sub, err := _Signerhub.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &SignerhubInitializedIterator{contract: _Signerhub.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Signerhub *SignerhubFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *SignerhubInitialized) (event.Subscription, error) {

	logs, sub, err := _Signerhub.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SignerhubInitialized)
				if err := _Signerhub.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Signerhub *SignerhubFilterer) ParseInitialized(log types.Log) (*SignerhubInitialized, error) {
	event := new(SignerhubInitialized)
	if err := _Signerhub.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SignerhubOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Signerhub contract.
type SignerhubOwnershipTransferredIterator struct {
	Event *SignerhubOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SignerhubOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SignerhubOwnershipTransferred)
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
		it.Event = new(SignerhubOwnershipTransferred)
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
func (it *SignerhubOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SignerhubOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SignerhubOwnershipTransferred represents a OwnershipTransferred event raised by the Signerhub contract.
type SignerhubOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Signerhub *SignerhubFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SignerhubOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Signerhub.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SignerhubOwnershipTransferredIterator{contract: _Signerhub.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Signerhub *SignerhubFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SignerhubOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Signerhub.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SignerhubOwnershipTransferred)
				if err := _Signerhub.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Signerhub *SignerhubFilterer) ParseOwnershipTransferred(log types.Log) (*SignerhubOwnershipTransferred, error) {
	event := new(SignerhubOwnershipTransferred)
	if err := _Signerhub.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SignerhubSignedAddedIterator is returned from FilterSignedAdded and is used to iterate over the raw logs and unpacked data for SignedAdded events raised by the Signerhub contract.
type SignerhubSignedAddedIterator struct {
	Event *SignerhubSignedAdded // Event containing the contract specifics and raw log

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
func (it *SignerhubSignedAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SignerhubSignedAdded)
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
		it.Event = new(SignerhubSignedAdded)
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
func (it *SignerhubSignedAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SignerhubSignedAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SignerhubSignedAdded represents a SignedAdded event raised by the Signerhub contract.
type SignerhubSignedAdded struct {
	Index     *big.Int
	Publickey [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSignedAdded is a free log retrieval operation binding the contract event 0x0e4e5b3ba228cadf0c65d31f94d29259d1b37003b9dccdd6b60b23e89804f8f1.
//
// Solidity: event SignedAdded(uint256 index, bytes32 publickey)
func (_Signerhub *SignerhubFilterer) FilterSignedAdded(opts *bind.FilterOpts) (*SignerhubSignedAddedIterator, error) {

	logs, sub, err := _Signerhub.contract.FilterLogs(opts, "SignedAdded")
	if err != nil {
		return nil, err
	}
	return &SignerhubSignedAddedIterator{contract: _Signerhub.contract, event: "SignedAdded", logs: logs, sub: sub}, nil
}

// WatchSignedAdded is a free log subscription operation binding the contract event 0x0e4e5b3ba228cadf0c65d31f94d29259d1b37003b9dccdd6b60b23e89804f8f1.
//
// Solidity: event SignedAdded(uint256 index, bytes32 publickey)
func (_Signerhub *SignerhubFilterer) WatchSignedAdded(opts *bind.WatchOpts, sink chan<- *SignerhubSignedAdded) (event.Subscription, error) {

	logs, sub, err := _Signerhub.contract.WatchLogs(opts, "SignedAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SignerhubSignedAdded)
				if err := _Signerhub.contract.UnpackLog(event, "SignedAdded", log); err != nil {
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

// ParseSignedAdded is a log parse operation binding the contract event 0x0e4e5b3ba228cadf0c65d31f94d29259d1b37003b9dccdd6b60b23e89804f8f1.
//
// Solidity: event SignedAdded(uint256 index, bytes32 publickey)
func (_Signerhub *SignerhubFilterer) ParseSignedAdded(log types.Log) (*SignerhubSignedAdded, error) {
	event := new(SignerhubSignedAdded)
	if err := _Signerhub.contract.UnpackLog(event, "SignedAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SignerhubSignedRemovedIterator is returned from FilterSignedRemoved and is used to iterate over the raw logs and unpacked data for SignedRemoved events raised by the Signerhub contract.
type SignerhubSignedRemovedIterator struct {
	Event *SignerhubSignedRemoved // Event containing the contract specifics and raw log

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
func (it *SignerhubSignedRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SignerhubSignedRemoved)
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
		it.Event = new(SignerhubSignedRemoved)
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
func (it *SignerhubSignedRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SignerhubSignedRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SignerhubSignedRemoved represents a SignedRemoved event raised by the Signerhub contract.
type SignerhubSignedRemoved struct {
	Publickey [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSignedRemoved is a free log retrieval operation binding the contract event 0x202ba177372096a533cb0be65537787905a2c9a9b25538d8d9f578706b412cb3.
//
// Solidity: event SignedRemoved(bytes32 publickey)
func (_Signerhub *SignerhubFilterer) FilterSignedRemoved(opts *bind.FilterOpts) (*SignerhubSignedRemovedIterator, error) {

	logs, sub, err := _Signerhub.contract.FilterLogs(opts, "SignedRemoved")
	if err != nil {
		return nil, err
	}
	return &SignerhubSignedRemovedIterator{contract: _Signerhub.contract, event: "SignedRemoved", logs: logs, sub: sub}, nil
}

// WatchSignedRemoved is a free log subscription operation binding the contract event 0x202ba177372096a533cb0be65537787905a2c9a9b25538d8d9f578706b412cb3.
//
// Solidity: event SignedRemoved(bytes32 publickey)
func (_Signerhub *SignerhubFilterer) WatchSignedRemoved(opts *bind.WatchOpts, sink chan<- *SignerhubSignedRemoved) (event.Subscription, error) {

	logs, sub, err := _Signerhub.contract.WatchLogs(opts, "SignedRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SignerhubSignedRemoved)
				if err := _Signerhub.contract.UnpackLog(event, "SignedRemoved", log); err != nil {
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
func (_Signerhub *SignerhubFilterer) ParseSignedRemoved(log types.Log) (*SignerhubSignedRemoved, error) {
	event := new(SignerhubSignedRemoved)
	if err := _Signerhub.contract.UnpackLog(event, "SignedRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SignerhubSignerBlacklistedIterator is returned from FilterSignerBlacklisted and is used to iterate over the raw logs and unpacked data for SignerBlacklisted events raised by the Signerhub contract.
type SignerhubSignerBlacklistedIterator struct {
	Event *SignerhubSignerBlacklisted // Event containing the contract specifics and raw log

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
func (it *SignerhubSignerBlacklistedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SignerhubSignerBlacklisted)
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
		it.Event = new(SignerhubSignerBlacklisted)
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
func (it *SignerhubSignerBlacklistedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SignerhubSignerBlacklistedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SignerhubSignerBlacklisted represents a SignerBlacklisted event raised by the Signerhub contract.
type SignerhubSignerBlacklisted struct {
	Publickey [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSignerBlacklisted is a free log retrieval operation binding the contract event 0xfd0fd0ce237fc8c6c5ea5042cba831db42434ca670d70ead573412793ad2b48c.
//
// Solidity: event SignerBlacklisted(bytes32 publickey)
func (_Signerhub *SignerhubFilterer) FilterSignerBlacklisted(opts *bind.FilterOpts) (*SignerhubSignerBlacklistedIterator, error) {

	logs, sub, err := _Signerhub.contract.FilterLogs(opts, "SignerBlacklisted")
	if err != nil {
		return nil, err
	}
	return &SignerhubSignerBlacklistedIterator{contract: _Signerhub.contract, event: "SignerBlacklisted", logs: logs, sub: sub}, nil
}

// WatchSignerBlacklisted is a free log subscription operation binding the contract event 0xfd0fd0ce237fc8c6c5ea5042cba831db42434ca670d70ead573412793ad2b48c.
//
// Solidity: event SignerBlacklisted(bytes32 publickey)
func (_Signerhub *SignerhubFilterer) WatchSignerBlacklisted(opts *bind.WatchOpts, sink chan<- *SignerhubSignerBlacklisted) (event.Subscription, error) {

	logs, sub, err := _Signerhub.contract.WatchLogs(opts, "SignerBlacklisted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SignerhubSignerBlacklisted)
				if err := _Signerhub.contract.UnpackLog(event, "SignerBlacklisted", log); err != nil {
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
func (_Signerhub *SignerhubFilterer) ParseSignerBlacklisted(log types.Log) (*SignerhubSignerBlacklisted, error) {
	event := new(SignerhubSignerBlacklisted)
	if err := _Signerhub.contract.UnpackLog(event, "SignerBlacklisted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SignerhubSignerWhitelistedIterator is returned from FilterSignerWhitelisted and is used to iterate over the raw logs and unpacked data for SignerWhitelisted events raised by the Signerhub contract.
type SignerhubSignerWhitelistedIterator struct {
	Event *SignerhubSignerWhitelisted // Event containing the contract specifics and raw log

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
func (it *SignerhubSignerWhitelistedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SignerhubSignerWhitelisted)
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
		it.Event = new(SignerhubSignerWhitelisted)
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
func (it *SignerhubSignerWhitelistedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SignerhubSignerWhitelistedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SignerhubSignerWhitelisted represents a SignerWhitelisted event raised by the Signerhub contract.
type SignerhubSignerWhitelisted struct {
	Publickey [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSignerWhitelisted is a free log retrieval operation binding the contract event 0x684d8290b28f7ee1d9799d0632bb71110e2f2c8feddb5493fb872b8b57faa927.
//
// Solidity: event SignerWhitelisted(bytes32 publickey)
func (_Signerhub *SignerhubFilterer) FilterSignerWhitelisted(opts *bind.FilterOpts) (*SignerhubSignerWhitelistedIterator, error) {

	logs, sub, err := _Signerhub.contract.FilterLogs(opts, "SignerWhitelisted")
	if err != nil {
		return nil, err
	}
	return &SignerhubSignerWhitelistedIterator{contract: _Signerhub.contract, event: "SignerWhitelisted", logs: logs, sub: sub}, nil
}

// WatchSignerWhitelisted is a free log subscription operation binding the contract event 0x684d8290b28f7ee1d9799d0632bb71110e2f2c8feddb5493fb872b8b57faa927.
//
// Solidity: event SignerWhitelisted(bytes32 publickey)
func (_Signerhub *SignerhubFilterer) WatchSignerWhitelisted(opts *bind.WatchOpts, sink chan<- *SignerhubSignerWhitelisted) (event.Subscription, error) {

	logs, sub, err := _Signerhub.contract.WatchLogs(opts, "SignerWhitelisted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SignerhubSignerWhitelisted)
				if err := _Signerhub.contract.UnpackLog(event, "SignerWhitelisted", log); err != nil {
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
func (_Signerhub *SignerhubFilterer) ParseSignerWhitelisted(log types.Log) (*SignerhubSignerWhitelisted, error) {
	event := new(SignerhubSignerWhitelisted)
	if err := _Signerhub.contract.UnpackLog(event, "SignerWhitelisted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
