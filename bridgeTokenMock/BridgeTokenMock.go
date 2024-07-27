// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bridgeTokenMock

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

// BridgeTokenMockMetaData contains all meta data concerning the BridgeTokenMock contract.
var BridgeTokenMockMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authority\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimal\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"faucet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"_decimal\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"_authority\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"sweep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f80fd5b50611ada8061001c5f395ff3fe608060405234801561000f575f80fd5b50600436106100fe575f3560e01c806376809ce311610095578063a9059cbb11610064578063a9059cbb14610298578063bf7e214f146102c8578063dd62ed3e146102e6578063de7ea79d14610316576100fe565b806376809ce3146102245780637b56c2b21461024257806395d89b411461025e5780639dc29fac1461027c576100fe565b8063313ce567116100d1578063313ce5671461019e57806340c10f19146101bc5780636ea056a9146101d857806370a08231146101f4576100fe565b806306fdde0314610102578063095ea7b31461012057806318160ddd1461015057806323b872dd1461016e575b5f80fd5b61010a610332565b604051610117919061121b565b60405180910390f35b61013a600480360381019061013591906112d9565b6103d0565b6040516101479190611331565b60405180910390f35b6101586103f2565b6040516101659190611359565b60405180910390f35b61018860048036038101906101839190611372565b610409565b6040516101959190611331565b60405180910390f35b6101a6610437565b6040516101b391906113dd565b60405180910390f35b6101d660048036038101906101d191906112d9565b61044b565b005b6101f260048036038101906101ed91906112d9565b6104e8565b005b61020e600480360381019061020991906113f6565b6105a7565b60405161021b9190611359565b60405180910390f35b61022c6105fa565b60405161023991906113dd565b60405180910390f35b61025c600480360381019061025791906112d9565b61060a565b005b610266610618565b604051610273919061121b565b60405180910390f35b610296600480360381019061029191906112d9565b6106b6565b005b6102b260048036038101906102ad91906112d9565b610753565b6040516102bf9190611331565b60405180910390f35b6102d0610775565b6040516102dd9190611430565b60405180910390f35b61030060048036038101906102fb9190611449565b61079a565b60405161030d9190611359565b60405180910390f35b610330600480360381019061032b91906115dd565b61082a565b005b60605f61033d610a08565b905080600301805461034e906116a6565b80601f016020809104026020016040519081016040528092919081815260200182805461037a906116a6565b80156103c55780601f1061039c576101008083540402835291602001916103c5565b820191905f5260205f20905b8154815290600101906020018083116103a857829003601f168201915b505050505091505090565b5f806103da610a2f565b90506103e7818585610a36565b600191505092915050565b5f806103fc610a08565b9050806002015491505090565b5f80610413610a2f565b9050610420858285610a48565b61042b858585610ada565b60019150509392505050565b5f805f9054906101000a900460ff16905090565b5f60019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146104da576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104d190611720565b60405180910390fd5b6104e48282610bca565b5050565b5f60019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610577576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161056e90611720565b60405180910390fd5b6105a3825f60019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1683610ada565b5050565b5f806105b1610a08565b9050805f015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054915050919050565b5f8054906101000a900460ff1681565b6106148282610bca565b5050565b60605f610623610a08565b9050806004018054610634906116a6565b80601f0160208091040260200160405190810160405280929190818152602001828054610660906116a6565b80156106ab5780601f10610682576101008083540402835291602001916106ab565b820191905f5260205f20905b81548152906001019060200180831161068e57829003601f168201915b505050505091505090565b5f60019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610745576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161073c90611720565b60405180910390fd5b61074f8282610c49565b5050565b5f8061075d610a2f565b905061076a818585610ada565b600191505092915050565b5f60019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b5f806107a4610a08565b9050806001015f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205491505092915050565b5f610833610cc8565b90505f815f0160089054906101000a900460ff161590505f825f015f9054906101000a900467ffffffffffffffff1690505f808267ffffffffffffffff1614801561087b5750825b90505f60018367ffffffffffffffff161480156108ae57505f3073ffffffffffffffffffffffffffffffffffffffff163b145b9050811580156108bc575080155b156108f3576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001855f015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055508315610940576001855f0160086101000a81548160ff0219169083151502179055505b61094a8989610cef565b865f806101000a81548160ff021916908360ff160217905550855f60016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555083156109fd575f855f0160086101000a81548160ff0219169083151502179055507fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d260016040516109f49190611793565b60405180910390a15b505050505050505050565b5f7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace00905090565b5f33905090565b610a438383836001610d05565b505050565b5f610a53848461079a565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114610ad45781811015610ac5578281836040517ffb8f41b2000000000000000000000000000000000000000000000000000000008152600401610abc939291906117ac565b60405180910390fd5b610ad384848484035f610d05565b5b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610b4a575f6040517f96c6fd1e000000000000000000000000000000000000000000000000000000008152600401610b419190611430565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610bba575f6040517fec442f05000000000000000000000000000000000000000000000000000000008152600401610bb19190611430565b60405180910390fd5b610bc5838383610ee2565b505050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610c3a575f6040517fec442f05000000000000000000000000000000000000000000000000000000008152600401610c319190611430565b60405180910390fd5b610c455f8383610ee2565b5050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610cb9575f6040517f96c6fd1e000000000000000000000000000000000000000000000000000000008152600401610cb09190611430565b60405180910390fd5b610cc4825f83610ee2565b5050565b5f7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00905090565b610cf7611111565b610d018282611151565b5050565b5f610d0e610a08565b90505f73ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff1603610d80575f6040517fe602df05000000000000000000000000000000000000000000000000000000008152600401610d779190611430565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1603610df0575f6040517f94280d62000000000000000000000000000000000000000000000000000000008152600401610de79190611430565b60405180910390fd5b82816001015f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508115610edb578373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92585604051610ed29190611359565b60405180910390a35b5050505050565b5f610eeb610a08565b90505f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1603610f3f5781816002015f828254610f33919061180e565b92505081905550611011565b5f815f015f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905082811015610fca578481846040517fe450d38c000000000000000000000000000000000000000000000000000000008152600401610fc1939291906117ac565b60405180910390fd5b828103825f015f8773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2081905550505b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff160361105a5781816002015f82825403925050819055506110a6565b81815f015f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055505b8273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516111039190611359565b60405180910390a350505050565b61111961118d565b61114f576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b611159611111565b5f611162610a08565b90508281600301908161117591906119d5565b508181600401908161118791906119d5565b50505050565b5f611196610cc8565b5f0160089054906101000a900460ff16905090565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f6111ed826111ab565b6111f781856111b5565b93506112078185602086016111c5565b611210816111d3565b840191505092915050565b5f6020820190508181035f83015261123381846111e3565b905092915050565b5f604051905090565b5f80fd5b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6112758261124c565b9050919050565b6112858161126b565b811461128f575f80fd5b50565b5f813590506112a08161127c565b92915050565b5f819050919050565b6112b8816112a6565b81146112c2575f80fd5b50565b5f813590506112d3816112af565b92915050565b5f80604083850312156112ef576112ee611244565b5b5f6112fc85828601611292565b925050602061130d858286016112c5565b9150509250929050565b5f8115159050919050565b61132b81611317565b82525050565b5f6020820190506113445f830184611322565b92915050565b611353816112a6565b82525050565b5f60208201905061136c5f83018461134a565b92915050565b5f805f6060848603121561138957611388611244565b5b5f61139686828701611292565b93505060206113a786828701611292565b92505060406113b8868287016112c5565b9150509250925092565b5f60ff82169050919050565b6113d7816113c2565b82525050565b5f6020820190506113f05f8301846113ce565b92915050565b5f6020828403121561140b5761140a611244565b5b5f61141884828501611292565b91505092915050565b61142a8161126b565b82525050565b5f6020820190506114435f830184611421565b92915050565b5f806040838503121561145f5761145e611244565b5b5f61146c85828601611292565b925050602061147d85828601611292565b9150509250929050565b5f80fd5b5f80fd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6114c5826111d3565b810181811067ffffffffffffffff821117156114e4576114e361148f565b5b80604052505050565b5f6114f661123b565b905061150282826114bc565b919050565b5f67ffffffffffffffff8211156115215761152061148f565b5b61152a826111d3565b9050602081019050919050565b828183375f83830152505050565b5f61155761155284611507565b6114ed565b9050828152602081018484840111156115735761157261148b565b5b61157e848285611537565b509392505050565b5f82601f83011261159a57611599611487565b5b81356115aa848260208601611545565b91505092915050565b6115bc816113c2565b81146115c6575f80fd5b50565b5f813590506115d7816115b3565b92915050565b5f805f80608085870312156115f5576115f4611244565b5b5f85013567ffffffffffffffff81111561161257611611611248565b5b61161e87828801611586565b945050602085013567ffffffffffffffff81111561163f5761163e611248565b5b61164b87828801611586565b935050604061165c878288016115c9565b925050606061166d87828801611292565b91505092959194509250565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806116bd57607f821691505b6020821081036116d0576116cf611679565b5b50919050565b7f427269646765546f6b656e3a20617574686f72697479000000000000000000005f82015250565b5f61170a6016836111b5565b9150611715826116d6565b602082019050919050565b5f6020820190508181035f830152611737816116fe565b9050919050565b5f819050919050565b5f67ffffffffffffffff82169050919050565b5f819050919050565b5f61177d6117786117738461173e565b61175a565b611747565b9050919050565b61178d81611763565b82525050565b5f6020820190506117a65f830184611784565b92915050565b5f6060820190506117bf5f830186611421565b6117cc602083018561134a565b6117d9604083018461134a565b949350505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f611818826112a6565b9150611823836112a6565b925082820190508082111561183b5761183a6117e1565b5b92915050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f6008830261189d7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82611862565b6118a78683611862565b95508019841693508086168417925050509392505050565b5f6118d96118d46118cf846112a6565b61175a565b6112a6565b9050919050565b5f819050919050565b6118f2836118bf565b6119066118fe826118e0565b84845461186e565b825550505050565b5f90565b61191a61190e565b6119258184846118e9565b505050565b5b818110156119485761193d5f82611912565b60018101905061192b565b5050565b601f82111561198d5761195e81611841565b61196784611853565b81016020851015611976578190505b61198a61198285611853565b83018261192a565b50505b505050565b5f82821c905092915050565b5f6119ad5f1984600802611992565b1980831691505092915050565b5f6119c5838361199e565b9150826002028217905092915050565b6119de826111ab565b67ffffffffffffffff8111156119f7576119f661148f565b5b611a0182546116a6565b611a0c82828561194c565b5f60209050601f831160018114611a3d575f8415611a2b578287015190505b611a3585826119ba565b865550611a9c565b601f198416611a4b86611841565b5f5b82811015611a7257848901518255600182019150602085019450602081019050611a4d565b86831015611a8f5784890151611a8b601f89168261199e565b8355505b6001600288020188555050505b50505050505056fea26469706673582212205982daa5572425c84ca966b03ce6cbb3723e5c595b6597fe455b2dc519c8a45464736f6c63430008190033",
}

// BridgeTokenMockABI is the input ABI used to generate the binding from.
// Deprecated: Use BridgeTokenMockMetaData.ABI instead.
var BridgeTokenMockABI = BridgeTokenMockMetaData.ABI

// BridgeTokenMockBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BridgeTokenMockMetaData.Bin instead.
var BridgeTokenMockBin = BridgeTokenMockMetaData.Bin

// DeployBridgeTokenMock deploys a new Ethereum contract, binding an instance of BridgeTokenMock to it.
func DeployBridgeTokenMock(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BridgeTokenMock, error) {
	parsed, err := BridgeTokenMockMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BridgeTokenMockBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BridgeTokenMock{BridgeTokenMockCaller: BridgeTokenMockCaller{contract: contract}, BridgeTokenMockTransactor: BridgeTokenMockTransactor{contract: contract}, BridgeTokenMockFilterer: BridgeTokenMockFilterer{contract: contract}}, nil
}

// BridgeTokenMock is an auto generated Go binding around an Ethereum contract.
type BridgeTokenMock struct {
	BridgeTokenMockCaller     // Read-only binding to the contract
	BridgeTokenMockTransactor // Write-only binding to the contract
	BridgeTokenMockFilterer   // Log filterer for contract events
}

// BridgeTokenMockCaller is an auto generated read-only Go binding around an Ethereum contract.
type BridgeTokenMockCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeTokenMockTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BridgeTokenMockTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeTokenMockFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BridgeTokenMockFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeTokenMockSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BridgeTokenMockSession struct {
	Contract     *BridgeTokenMock  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeTokenMockCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BridgeTokenMockCallerSession struct {
	Contract *BridgeTokenMockCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// BridgeTokenMockTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BridgeTokenMockTransactorSession struct {
	Contract     *BridgeTokenMockTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// BridgeTokenMockRaw is an auto generated low-level Go binding around an Ethereum contract.
type BridgeTokenMockRaw struct {
	Contract *BridgeTokenMock // Generic contract binding to access the raw methods on
}

// BridgeTokenMockCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BridgeTokenMockCallerRaw struct {
	Contract *BridgeTokenMockCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeTokenMockTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BridgeTokenMockTransactorRaw struct {
	Contract *BridgeTokenMockTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridgeTokenMock creates a new instance of BridgeTokenMock, bound to a specific deployed contract.
func NewBridgeTokenMock(address common.Address, backend bind.ContractBackend) (*BridgeTokenMock, error) {
	contract, err := bindBridgeTokenMock(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BridgeTokenMock{BridgeTokenMockCaller: BridgeTokenMockCaller{contract: contract}, BridgeTokenMockTransactor: BridgeTokenMockTransactor{contract: contract}, BridgeTokenMockFilterer: BridgeTokenMockFilterer{contract: contract}}, nil
}

// NewBridgeTokenMockCaller creates a new read-only instance of BridgeTokenMock, bound to a specific deployed contract.
func NewBridgeTokenMockCaller(address common.Address, caller bind.ContractCaller) (*BridgeTokenMockCaller, error) {
	contract, err := bindBridgeTokenMock(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeTokenMockCaller{contract: contract}, nil
}

// NewBridgeTokenMockTransactor creates a new write-only instance of BridgeTokenMock, bound to a specific deployed contract.
func NewBridgeTokenMockTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeTokenMockTransactor, error) {
	contract, err := bindBridgeTokenMock(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeTokenMockTransactor{contract: contract}, nil
}

// NewBridgeTokenMockFilterer creates a new log filterer instance of BridgeTokenMock, bound to a specific deployed contract.
func NewBridgeTokenMockFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeTokenMockFilterer, error) {
	contract, err := bindBridgeTokenMock(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeTokenMockFilterer{contract: contract}, nil
}

// bindBridgeTokenMock binds a generic wrapper to an already deployed contract.
func bindBridgeTokenMock(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BridgeTokenMockMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeTokenMock *BridgeTokenMockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeTokenMock.Contract.BridgeTokenMockCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeTokenMock *BridgeTokenMockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.BridgeTokenMockTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeTokenMock *BridgeTokenMockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.BridgeTokenMockTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeTokenMock *BridgeTokenMockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeTokenMock.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeTokenMock *BridgeTokenMockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeTokenMock *BridgeTokenMockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_BridgeTokenMock *BridgeTokenMockCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BridgeTokenMock.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_BridgeTokenMock *BridgeTokenMockSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _BridgeTokenMock.Contract.Allowance(&_BridgeTokenMock.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_BridgeTokenMock *BridgeTokenMockCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _BridgeTokenMock.Contract.Allowance(&_BridgeTokenMock.CallOpts, owner, spender)
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (_BridgeTokenMock *BridgeTokenMockCaller) Authority(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BridgeTokenMock.contract.Call(opts, &out, "authority")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (_BridgeTokenMock *BridgeTokenMockSession) Authority() (common.Address, error) {
	return _BridgeTokenMock.Contract.Authority(&_BridgeTokenMock.CallOpts)
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (_BridgeTokenMock *BridgeTokenMockCallerSession) Authority() (common.Address, error) {
	return _BridgeTokenMock.Contract.Authority(&_BridgeTokenMock.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_BridgeTokenMock *BridgeTokenMockCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BridgeTokenMock.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_BridgeTokenMock *BridgeTokenMockSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _BridgeTokenMock.Contract.BalanceOf(&_BridgeTokenMock.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_BridgeTokenMock *BridgeTokenMockCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _BridgeTokenMock.Contract.BalanceOf(&_BridgeTokenMock.CallOpts, account)
}

// Decimal is a free data retrieval call binding the contract method 0x76809ce3.
//
// Solidity: function decimal() view returns(uint8)
func (_BridgeTokenMock *BridgeTokenMockCaller) Decimal(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _BridgeTokenMock.contract.Call(opts, &out, "decimal")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimal is a free data retrieval call binding the contract method 0x76809ce3.
//
// Solidity: function decimal() view returns(uint8)
func (_BridgeTokenMock *BridgeTokenMockSession) Decimal() (uint8, error) {
	return _BridgeTokenMock.Contract.Decimal(&_BridgeTokenMock.CallOpts)
}

// Decimal is a free data retrieval call binding the contract method 0x76809ce3.
//
// Solidity: function decimal() view returns(uint8)
func (_BridgeTokenMock *BridgeTokenMockCallerSession) Decimal() (uint8, error) {
	return _BridgeTokenMock.Contract.Decimal(&_BridgeTokenMock.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_BridgeTokenMock *BridgeTokenMockCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _BridgeTokenMock.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_BridgeTokenMock *BridgeTokenMockSession) Decimals() (uint8, error) {
	return _BridgeTokenMock.Contract.Decimals(&_BridgeTokenMock.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_BridgeTokenMock *BridgeTokenMockCallerSession) Decimals() (uint8, error) {
	return _BridgeTokenMock.Contract.Decimals(&_BridgeTokenMock.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_BridgeTokenMock *BridgeTokenMockCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BridgeTokenMock.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_BridgeTokenMock *BridgeTokenMockSession) Name() (string, error) {
	return _BridgeTokenMock.Contract.Name(&_BridgeTokenMock.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_BridgeTokenMock *BridgeTokenMockCallerSession) Name() (string, error) {
	return _BridgeTokenMock.Contract.Name(&_BridgeTokenMock.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_BridgeTokenMock *BridgeTokenMockCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BridgeTokenMock.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_BridgeTokenMock *BridgeTokenMockSession) Symbol() (string, error) {
	return _BridgeTokenMock.Contract.Symbol(&_BridgeTokenMock.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_BridgeTokenMock *BridgeTokenMockCallerSession) Symbol() (string, error) {
	return _BridgeTokenMock.Contract.Symbol(&_BridgeTokenMock.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_BridgeTokenMock *BridgeTokenMockCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BridgeTokenMock.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_BridgeTokenMock *BridgeTokenMockSession) TotalSupply() (*big.Int, error) {
	return _BridgeTokenMock.Contract.TotalSupply(&_BridgeTokenMock.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_BridgeTokenMock *BridgeTokenMockCallerSession) TotalSupply() (*big.Int, error) {
	return _BridgeTokenMock.Contract.TotalSupply(&_BridgeTokenMock.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_BridgeTokenMock *BridgeTokenMockTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_BridgeTokenMock *BridgeTokenMockSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.Approve(&_BridgeTokenMock.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_BridgeTokenMock *BridgeTokenMockTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.Approve(&_BridgeTokenMock.TransactOpts, spender, value)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 amount) returns()
func (_BridgeTokenMock *BridgeTokenMockTransactor) Burn(opts *bind.TransactOpts, from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.contract.Transact(opts, "burn", from, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 amount) returns()
func (_BridgeTokenMock *BridgeTokenMockSession) Burn(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.Burn(&_BridgeTokenMock.TransactOpts, from, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 amount) returns()
func (_BridgeTokenMock *BridgeTokenMockTransactorSession) Burn(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.Burn(&_BridgeTokenMock.TransactOpts, from, amount)
}

// Faucet is a paid mutator transaction binding the contract method 0x7b56c2b2.
//
// Solidity: function faucet(address to, uint256 amount) returns()
func (_BridgeTokenMock *BridgeTokenMockTransactor) Faucet(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.contract.Transact(opts, "faucet", to, amount)
}

// Faucet is a paid mutator transaction binding the contract method 0x7b56c2b2.
//
// Solidity: function faucet(address to, uint256 amount) returns()
func (_BridgeTokenMock *BridgeTokenMockSession) Faucet(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.Faucet(&_BridgeTokenMock.TransactOpts, to, amount)
}

// Faucet is a paid mutator transaction binding the contract method 0x7b56c2b2.
//
// Solidity: function faucet(address to, uint256 amount) returns()
func (_BridgeTokenMock *BridgeTokenMockTransactorSession) Faucet(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.Faucet(&_BridgeTokenMock.TransactOpts, to, amount)
}

// Initialize is a paid mutator transaction binding the contract method 0xde7ea79d.
//
// Solidity: function initialize(string name, string symbol, uint8 _decimal, address _authority) returns()
func (_BridgeTokenMock *BridgeTokenMockTransactor) Initialize(opts *bind.TransactOpts, name string, symbol string, _decimal uint8, _authority common.Address) (*types.Transaction, error) {
	return _BridgeTokenMock.contract.Transact(opts, "initialize", name, symbol, _decimal, _authority)
}

// Initialize is a paid mutator transaction binding the contract method 0xde7ea79d.
//
// Solidity: function initialize(string name, string symbol, uint8 _decimal, address _authority) returns()
func (_BridgeTokenMock *BridgeTokenMockSession) Initialize(name string, symbol string, _decimal uint8, _authority common.Address) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.Initialize(&_BridgeTokenMock.TransactOpts, name, symbol, _decimal, _authority)
}

// Initialize is a paid mutator transaction binding the contract method 0xde7ea79d.
//
// Solidity: function initialize(string name, string symbol, uint8 _decimal, address _authority) returns()
func (_BridgeTokenMock *BridgeTokenMockTransactorSession) Initialize(name string, symbol string, _decimal uint8, _authority common.Address) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.Initialize(&_BridgeTokenMock.TransactOpts, name, symbol, _decimal, _authority)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_BridgeTokenMock *BridgeTokenMockTransactor) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.contract.Transact(opts, "mint", to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_BridgeTokenMock *BridgeTokenMockSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.Mint(&_BridgeTokenMock.TransactOpts, to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_BridgeTokenMock *BridgeTokenMockTransactorSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.Mint(&_BridgeTokenMock.TransactOpts, to, amount)
}

// Sweep is a paid mutator transaction binding the contract method 0x6ea056a9.
//
// Solidity: function sweep(address from, uint256 amount) returns()
func (_BridgeTokenMock *BridgeTokenMockTransactor) Sweep(opts *bind.TransactOpts, from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.contract.Transact(opts, "sweep", from, amount)
}

// Sweep is a paid mutator transaction binding the contract method 0x6ea056a9.
//
// Solidity: function sweep(address from, uint256 amount) returns()
func (_BridgeTokenMock *BridgeTokenMockSession) Sweep(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.Sweep(&_BridgeTokenMock.TransactOpts, from, amount)
}

// Sweep is a paid mutator transaction binding the contract method 0x6ea056a9.
//
// Solidity: function sweep(address from, uint256 amount) returns()
func (_BridgeTokenMock *BridgeTokenMockTransactorSession) Sweep(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.Sweep(&_BridgeTokenMock.TransactOpts, from, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_BridgeTokenMock *BridgeTokenMockTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_BridgeTokenMock *BridgeTokenMockSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.Transfer(&_BridgeTokenMock.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_BridgeTokenMock *BridgeTokenMockTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.Transfer(&_BridgeTokenMock.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_BridgeTokenMock *BridgeTokenMockTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_BridgeTokenMock *BridgeTokenMockSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.TransferFrom(&_BridgeTokenMock.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_BridgeTokenMock *BridgeTokenMockTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _BridgeTokenMock.Contract.TransferFrom(&_BridgeTokenMock.TransactOpts, from, to, value)
}

// BridgeTokenMockApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the BridgeTokenMock contract.
type BridgeTokenMockApprovalIterator struct {
	Event *BridgeTokenMockApproval // Event containing the contract specifics and raw log

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
func (it *BridgeTokenMockApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeTokenMockApproval)
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
		it.Event = new(BridgeTokenMockApproval)
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
func (it *BridgeTokenMockApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeTokenMockApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeTokenMockApproval represents a Approval event raised by the BridgeTokenMock contract.
type BridgeTokenMockApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_BridgeTokenMock *BridgeTokenMockFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*BridgeTokenMockApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _BridgeTokenMock.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &BridgeTokenMockApprovalIterator{contract: _BridgeTokenMock.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_BridgeTokenMock *BridgeTokenMockFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *BridgeTokenMockApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _BridgeTokenMock.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeTokenMockApproval)
				if err := _BridgeTokenMock.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_BridgeTokenMock *BridgeTokenMockFilterer) ParseApproval(log types.Log) (*BridgeTokenMockApproval, error) {
	event := new(BridgeTokenMockApproval)
	if err := _BridgeTokenMock.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeTokenMockInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the BridgeTokenMock contract.
type BridgeTokenMockInitializedIterator struct {
	Event *BridgeTokenMockInitialized // Event containing the contract specifics and raw log

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
func (it *BridgeTokenMockInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeTokenMockInitialized)
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
		it.Event = new(BridgeTokenMockInitialized)
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
func (it *BridgeTokenMockInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeTokenMockInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeTokenMockInitialized represents a Initialized event raised by the BridgeTokenMock contract.
type BridgeTokenMockInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_BridgeTokenMock *BridgeTokenMockFilterer) FilterInitialized(opts *bind.FilterOpts) (*BridgeTokenMockInitializedIterator, error) {

	logs, sub, err := _BridgeTokenMock.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &BridgeTokenMockInitializedIterator{contract: _BridgeTokenMock.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_BridgeTokenMock *BridgeTokenMockFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *BridgeTokenMockInitialized) (event.Subscription, error) {

	logs, sub, err := _BridgeTokenMock.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeTokenMockInitialized)
				if err := _BridgeTokenMock.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_BridgeTokenMock *BridgeTokenMockFilterer) ParseInitialized(log types.Log) (*BridgeTokenMockInitialized, error) {
	event := new(BridgeTokenMockInitialized)
	if err := _BridgeTokenMock.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeTokenMockTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the BridgeTokenMock contract.
type BridgeTokenMockTransferIterator struct {
	Event *BridgeTokenMockTransfer // Event containing the contract specifics and raw log

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
func (it *BridgeTokenMockTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeTokenMockTransfer)
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
		it.Event = new(BridgeTokenMockTransfer)
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
func (it *BridgeTokenMockTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeTokenMockTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeTokenMockTransfer represents a Transfer event raised by the BridgeTokenMock contract.
type BridgeTokenMockTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_BridgeTokenMock *BridgeTokenMockFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BridgeTokenMockTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BridgeTokenMock.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BridgeTokenMockTransferIterator{contract: _BridgeTokenMock.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_BridgeTokenMock *BridgeTokenMockFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *BridgeTokenMockTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BridgeTokenMock.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeTokenMockTransfer)
				if err := _BridgeTokenMock.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_BridgeTokenMock *BridgeTokenMockFilterer) ParseTransfer(log types.Log) (*BridgeTokenMockTransfer, error) {
	event := new(BridgeTokenMockTransfer)
	if err := _BridgeTokenMock.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
