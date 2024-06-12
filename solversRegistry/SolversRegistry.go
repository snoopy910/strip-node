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
	ABI: "[{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"domain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"whitelisted\",\"type\":\"bool\"}],\"name\":\"SolverUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"domain\",\"type\":\"string\"}],\"name\":\"getChains\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"solvers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"exists\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"whitelisted\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"lastUpdated\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSolvers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"domain\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"whitelisted\",\"type\":\"bool\"},{\"internalType\":\"uint256[]\",\"name\":\"chains\",\"type\":\"uint256[]\"}],\"name\":\"updateSolver\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f80fd5b506111338061001c5f395ff3fe608060405234801561000f575f80fd5b5060043610610086575f3560e01c80638da5cb5b116100595780638da5cb5b146101005780638e7663411461011e578063b91bacc61461013c578063f2fde38b1461015857610086565b8063083603c81461008a5780631968fa53146100ba578063715018a6146100ec5780638129fc1c146100f6575b5f80fd5b6100a4600480360381019061009f9190610992565b610174565b6040516100b19190610a9d565b60405180910390f35b6100d460048036038101906100cf9190610bf5565b6101ed565b6040516100e393929190610c65565b60405180910390f35b6100f4610242565b005b6100fe610255565b005b6101086103d5565b6040516101159190610cd9565b60405180910390f35b61012661040a565b6040516101339190610cf2565b60405180910390f35b61015660048036038101906101519190610e1f565b610410565b005b610172600480360381019061016d9190610ed6565b610593565b005b60605f8383604051610187929190610f2f565b90815260200160405180910390206001018054806020026020016040519081016040528092919081815260200182805480156101e057602002820191905f5260205f20905b8154815260200190600101908083116101cc575b5050505050905092915050565b5f818051602081018201805184825260208301602085012081835280955050505050505f91509050805f015f9054906101000a900460ff1690805f0160019054906101000a900460ff16908060020154905083565b61024a610617565b6102535f61069e565b565b5f61025e61076f565b90505f815f0160089054906101000a900460ff161590505f825f015f9054906101000a900467ffffffffffffffff1690505f808267ffffffffffffffff161480156102a65750825b90505f60018367ffffffffffffffff161480156102d957505f3073ffffffffffffffffffffffffffffffffffffffff163b145b9050811580156102e7575080155b1561031e576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001855f015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550831561036b576001855f0160086101000a81548160ff0219169083151502179055505b61037433610796565b83156103ce575f855f0160086101000a81548160ff0219169083151502179055507fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d260016040516103c59190610f9c565b60405180910390a15b5050505050565b5f806103df6107aa565b9050805f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1691505090565b60015481565b610418610617565b5f848490500361045d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104549061100f565b60405180910390fd5b5f848460405161046e929190610f2f565b90815260200160405180910390205f015f9054906101000a900460ff166104a75760015f8154809291906104a19061105a565b91905055505b60405180608001604052806001151581526020018315158152602001828152602001428152505f85856040516104de929190610f2f565b90815260200160405180910390205f820151815f015f6101000a81548160ff0219169083151502179055506020820151815f0160016101000a81548160ff02191690831515021790555060408201518160010190805190602001906105449291906108ba565b50606082015181600201559050507f4c40c26c445e2acb57aa0f02b1bc6d90de42da3d0e9d181967d3a8d142a666a6848484604051610585939291906110cd565b60405180910390a150505050565b61059b610617565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361060b575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016106029190610cd9565b60405180910390fd5b6106148161069e565b50565b61061f6107d1565b73ffffffffffffffffffffffffffffffffffffffff1661063d6103d5565b73ffffffffffffffffffffffffffffffffffffffff161461069c576106606107d1565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016106939190610cd9565b60405180910390fd5b565b5f6106a76107aa565b90505f815f015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905082825f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508273ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3505050565b5f7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00905090565b61079e6107d8565b6107a781610818565b50565b5f7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300905090565b5f33905090565b6107e061089c565b610816576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6108206107d8565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610890575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016108879190610cd9565b60405180910390fd5b6108998161069e565b50565b5f6108a561076f565b5f0160089054906101000a900460ff16905090565b828054828255905f5260205f209081019282156108f4579160200282015b828111156108f35782518255916020019190600101906108d8565b5b5090506109019190610905565b5090565b5b8082111561091c575f815f905550600101610906565b5090565b5f604051905090565b5f80fd5b5f80fd5b5f80fd5b5f80fd5b5f80fd5b5f8083601f84011261095257610951610931565b5b8235905067ffffffffffffffff81111561096f5761096e610935565b5b60208301915083600182028301111561098b5761098a610939565b5b9250929050565b5f80602083850312156109a8576109a7610929565b5b5f83013567ffffffffffffffff8111156109c5576109c461092d565b5b6109d18582860161093d565b92509250509250929050565b5f81519050919050565b5f82825260208201905092915050565b5f819050602082019050919050565b5f819050919050565b610a1881610a06565b82525050565b5f610a298383610a0f565b60208301905092915050565b5f602082019050919050565b5f610a4b826109dd565b610a5581856109e7565b9350610a60836109f7565b805f5b83811015610a90578151610a778882610a1e565b9750610a8283610a35565b925050600181019050610a63565b5085935050505092915050565b5f6020820190508181035f830152610ab58184610a41565b905092915050565b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b610b0782610ac1565b810181811067ffffffffffffffff82111715610b2657610b25610ad1565b5b80604052505050565b5f610b38610920565b9050610b448282610afe565b919050565b5f67ffffffffffffffff821115610b6357610b62610ad1565b5b610b6c82610ac1565b9050602081019050919050565b828183375f83830152505050565b5f610b99610b9484610b49565b610b2f565b905082815260208101848484011115610bb557610bb4610abd565b5b610bc0848285610b79565b509392505050565b5f82601f830112610bdc57610bdb610931565b5b8135610bec848260208601610b87565b91505092915050565b5f60208284031215610c0a57610c09610929565b5b5f82013567ffffffffffffffff811115610c2757610c2661092d565b5b610c3384828501610bc8565b91505092915050565b5f8115159050919050565b610c5081610c3c565b82525050565b610c5f81610a06565b82525050565b5f606082019050610c785f830186610c47565b610c856020830185610c47565b610c926040830184610c56565b949350505050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610cc382610c9a565b9050919050565b610cd381610cb9565b82525050565b5f602082019050610cec5f830184610cca565b92915050565b5f602082019050610d055f830184610c56565b92915050565b610d1481610c3c565b8114610d1e575f80fd5b50565b5f81359050610d2f81610d0b565b92915050565b5f67ffffffffffffffff821115610d4f57610d4e610ad1565b5b602082029050602081019050919050565b610d6981610a06565b8114610d73575f80fd5b50565b5f81359050610d8481610d60565b92915050565b5f610d9c610d9784610d35565b610b2f565b90508083825260208201905060208402830185811115610dbf57610dbe610939565b5b835b81811015610de85780610dd48882610d76565b845260208401935050602081019050610dc1565b5050509392505050565b5f82601f830112610e0657610e05610931565b5b8135610e16848260208601610d8a565b91505092915050565b5f805f8060608587031215610e3757610e36610929565b5b5f85013567ffffffffffffffff811115610e5457610e5361092d565b5b610e608782880161093d565b94509450506020610e7387828801610d21565b925050604085013567ffffffffffffffff811115610e9457610e9361092d565b5b610ea087828801610df2565b91505092959194509250565b610eb581610cb9565b8114610ebf575f80fd5b50565b5f81359050610ed081610eac565b92915050565b5f60208284031215610eeb57610eea610929565b5b5f610ef884828501610ec2565b91505092915050565b5f81905092915050565b5f610f168385610f01565b9350610f23838584610b79565b82840190509392505050565b5f610f3b828486610f0b565b91508190509392505050565b5f819050919050565b5f67ffffffffffffffff82169050919050565b5f819050919050565b5f610f86610f81610f7c84610f47565b610f63565b610f50565b9050919050565b610f9681610f6c565b82525050565b5f602082019050610faf5f830184610f8d565b92915050565b5f82825260208201905092915050565b7f646f6d61696e20697320656d70747900000000000000000000000000000000005f82015250565b5f610ff9600f83610fb5565b915061100482610fc5565b602082019050919050565b5f6020820190508181035f83015261102681610fed565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61106482610a06565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036110965761109561102d565b5b600182019050919050565b5f6110ac8385610fb5565b93506110b9838584610b79565b6110c283610ac1565b840190509392505050565b5f6040820190508181035f8301526110e68185876110a1565b90506110f56020830184610c47565b94935050505056fea2646970667358221220b07385d47b3d7d4fa0d90a29eba00a748445baa603993325ea13b341586ef8bd64736f6c63430008190033",
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

// GetChains is a free data retrieval call binding the contract method 0x083603c8.
//
// Solidity: function getChains(string domain) view returns(uint256[])
func (_SolversRegistry *SolversRegistryCaller) GetChains(opts *bind.CallOpts, domain string) ([]*big.Int, error) {
	var out []interface{}
	err := _SolversRegistry.contract.Call(opts, &out, "getChains", domain)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetChains is a free data retrieval call binding the contract method 0x083603c8.
//
// Solidity: function getChains(string domain) view returns(uint256[])
func (_SolversRegistry *SolversRegistrySession) GetChains(domain string) ([]*big.Int, error) {
	return _SolversRegistry.Contract.GetChains(&_SolversRegistry.CallOpts, domain)
}

// GetChains is a free data retrieval call binding the contract method 0x083603c8.
//
// Solidity: function getChains(string domain) view returns(uint256[])
func (_SolversRegistry *SolversRegistryCallerSession) GetChains(domain string) ([]*big.Int, error) {
	return _SolversRegistry.Contract.GetChains(&_SolversRegistry.CallOpts, domain)
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
// Solidity: function solvers(string ) view returns(bool exists, bool whitelisted, uint256 lastUpdated)
func (_SolversRegistry *SolversRegistryCaller) Solvers(opts *bind.CallOpts, arg0 string) (struct {
	Exists      bool
	Whitelisted bool
	LastUpdated *big.Int
}, error) {
	var out []interface{}
	err := _SolversRegistry.contract.Call(opts, &out, "solvers", arg0)

	outstruct := new(struct {
		Exists      bool
		Whitelisted bool
		LastUpdated *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Exists = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Whitelisted = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.LastUpdated = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Solvers is a free data retrieval call binding the contract method 0x1968fa53.
//
// Solidity: function solvers(string ) view returns(bool exists, bool whitelisted, uint256 lastUpdated)
func (_SolversRegistry *SolversRegistrySession) Solvers(arg0 string) (struct {
	Exists      bool
	Whitelisted bool
	LastUpdated *big.Int
}, error) {
	return _SolversRegistry.Contract.Solvers(&_SolversRegistry.CallOpts, arg0)
}

// Solvers is a free data retrieval call binding the contract method 0x1968fa53.
//
// Solidity: function solvers(string ) view returns(bool exists, bool whitelisted, uint256 lastUpdated)
func (_SolversRegistry *SolversRegistryCallerSession) Solvers(arg0 string) (struct {
	Exists      bool
	Whitelisted bool
	LastUpdated *big.Int
}, error) {
	return _SolversRegistry.Contract.Solvers(&_SolversRegistry.CallOpts, arg0)
}

// TotalSolvers is a free data retrieval call binding the contract method 0x8e766341.
//
// Solidity: function totalSolvers() view returns(uint256)
func (_SolversRegistry *SolversRegistryCaller) TotalSolvers(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SolversRegistry.contract.Call(opts, &out, "totalSolvers")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSolvers is a free data retrieval call binding the contract method 0x8e766341.
//
// Solidity: function totalSolvers() view returns(uint256)
func (_SolversRegistry *SolversRegistrySession) TotalSolvers() (*big.Int, error) {
	return _SolversRegistry.Contract.TotalSolvers(&_SolversRegistry.CallOpts)
}

// TotalSolvers is a free data retrieval call binding the contract method 0x8e766341.
//
// Solidity: function totalSolvers() view returns(uint256)
func (_SolversRegistry *SolversRegistryCallerSession) TotalSolvers() (*big.Int, error) {
	return _SolversRegistry.Contract.TotalSolvers(&_SolversRegistry.CallOpts)
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

// UpdateSolver is a paid mutator transaction binding the contract method 0xb91bacc6.
//
// Solidity: function updateSolver(string domain, bool whitelisted, uint256[] chains) returns()
func (_SolversRegistry *SolversRegistryTransactor) UpdateSolver(opts *bind.TransactOpts, domain string, whitelisted bool, chains []*big.Int) (*types.Transaction, error) {
	return _SolversRegistry.contract.Transact(opts, "updateSolver", domain, whitelisted, chains)
}

// UpdateSolver is a paid mutator transaction binding the contract method 0xb91bacc6.
//
// Solidity: function updateSolver(string domain, bool whitelisted, uint256[] chains) returns()
func (_SolversRegistry *SolversRegistrySession) UpdateSolver(domain string, whitelisted bool, chains []*big.Int) (*types.Transaction, error) {
	return _SolversRegistry.Contract.UpdateSolver(&_SolversRegistry.TransactOpts, domain, whitelisted, chains)
}

// UpdateSolver is a paid mutator transaction binding the contract method 0xb91bacc6.
//
// Solidity: function updateSolver(string domain, bool whitelisted, uint256[] chains) returns()
func (_SolversRegistry *SolversRegistryTransactorSession) UpdateSolver(domain string, whitelisted bool, chains []*big.Int) (*types.Transaction, error) {
	return _SolversRegistry.Contract.UpdateSolver(&_SolversRegistry.TransactOpts, domain, whitelisted, chains)
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
