// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package owlto20

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

// Owlto20MetaData contains all meta data concerning the Owlto20 contract.
var Owlto20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AmountError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LengthError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotOwnerError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressError\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"target\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"isOwltoTransfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"target\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"ercToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// Owlto20ABI is the input ABI used to generate the binding from.
// Deprecated: Use Owlto20MetaData.ABI instead.
var Owlto20ABI = Owlto20MetaData.ABI

// Owlto20 is an auto generated Go binding around an Ethereum contract.
type Owlto20 struct {
	Owlto20Caller     // Read-only binding to the contract
	Owlto20Transactor // Write-only binding to the contract
	Owlto20Filterer   // Log filterer for contract events
}

// Owlto20Caller is an auto generated read-only Go binding around an Ethereum contract.
type Owlto20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Owlto20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Owlto20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Owlto20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Owlto20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Owlto20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Owlto20Session struct {
	Contract     *Owlto20          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Owlto20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Owlto20CallerSession struct {
	Contract *Owlto20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// Owlto20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Owlto20TransactorSession struct {
	Contract     *Owlto20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// Owlto20Raw is an auto generated low-level Go binding around an Ethereum contract.
type Owlto20Raw struct {
	Contract *Owlto20 // Generic contract binding to access the raw methods on
}

// Owlto20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Owlto20CallerRaw struct {
	Contract *Owlto20Caller // Generic read-only contract binding to access the raw methods on
}

// Owlto20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Owlto20TransactorRaw struct {
	Contract *Owlto20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewOwlto20 creates a new instance of Owlto20, bound to a specific deployed contract.
func NewOwlto20(address common.Address, backend bind.ContractBackend) (*Owlto20, error) {
	contract, err := bindOwlto20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Owlto20{Owlto20Caller: Owlto20Caller{contract: contract}, Owlto20Transactor: Owlto20Transactor{contract: contract}, Owlto20Filterer: Owlto20Filterer{contract: contract}}, nil
}

// NewOwlto20Caller creates a new read-only instance of Owlto20, bound to a specific deployed contract.
func NewOwlto20Caller(address common.Address, caller bind.ContractCaller) (*Owlto20Caller, error) {
	contract, err := bindOwlto20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Owlto20Caller{contract: contract}, nil
}

// NewOwlto20Transactor creates a new write-only instance of Owlto20, bound to a specific deployed contract.
func NewOwlto20Transactor(address common.Address, transactor bind.ContractTransactor) (*Owlto20Transactor, error) {
	contract, err := bindOwlto20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Owlto20Transactor{contract: contract}, nil
}

// NewOwlto20Filterer creates a new log filterer instance of Owlto20, bound to a specific deployed contract.
func NewOwlto20Filterer(address common.Address, filterer bind.ContractFilterer) (*Owlto20Filterer, error) {
	contract, err := bindOwlto20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Owlto20Filterer{contract: contract}, nil
}

// bindOwlto20 binds a generic wrapper to an already deployed contract.
func bindOwlto20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := Owlto20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Owlto20 *Owlto20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Owlto20.Contract.Owlto20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Owlto20 *Owlto20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owlto20.Contract.Owlto20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Owlto20 *Owlto20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owlto20.Contract.Owlto20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Owlto20 *Owlto20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Owlto20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Owlto20 *Owlto20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owlto20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Owlto20 *Owlto20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Owlto20.Contract.contract.Transact(opts, method, params...)
}

// IsOwltoTransfer is a free data retrieval call binding the contract method 0x5a62795e.
//
// Solidity: function isOwltoTransfer() pure returns(bool)
func (_Owlto20 *Owlto20Caller) IsOwltoTransfer(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Owlto20.contract.Call(opts, &out, "isOwltoTransfer")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwltoTransfer is a free data retrieval call binding the contract method 0x5a62795e.
//
// Solidity: function isOwltoTransfer() pure returns(bool)
func (_Owlto20 *Owlto20Session) IsOwltoTransfer() (bool, error) {
	return _Owlto20.Contract.IsOwltoTransfer(&_Owlto20.CallOpts)
}

// IsOwltoTransfer is a free data retrieval call binding the contract method 0x5a62795e.
//
// Solidity: function isOwltoTransfer() pure returns(bool)
func (_Owlto20 *Owlto20CallerSession) IsOwltoTransfer() (bool, error) {
	return _Owlto20.Contract.IsOwltoTransfer(&_Owlto20.CallOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0x2952438e.
//
// Solidity: function transfer(string target, address ercToken, address maker, uint256 amount) payable returns()
func (_Owlto20 *Owlto20Transactor) Transfer(opts *bind.TransactOpts, target string, ercToken common.Address, maker common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Owlto20.contract.Transact(opts, "transfer", target, ercToken, maker, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0x2952438e.
//
// Solidity: function transfer(string target, address ercToken, address maker, uint256 amount) payable returns()
func (_Owlto20 *Owlto20Session) Transfer(target string, ercToken common.Address, maker common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Owlto20.Contract.Transfer(&_Owlto20.TransactOpts, target, ercToken, maker, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0x2952438e.
//
// Solidity: function transfer(string target, address ercToken, address maker, uint256 amount) payable returns()
func (_Owlto20 *Owlto20TransactorSession) Transfer(target string, ercToken common.Address, maker common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Owlto20.Contract.Transfer(&_Owlto20.TransactOpts, target, ercToken, maker, amount)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Owlto20 *Owlto20Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Owlto20.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Owlto20 *Owlto20Session) Receive() (*types.Transaction, error) {
	return _Owlto20.Contract.Receive(&_Owlto20.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Owlto20 *Owlto20TransactorSession) Receive() (*types.Transaction, error) {
	return _Owlto20.Contract.Receive(&_Owlto20.TransactOpts)
}

// Owlto20DepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the Owlto20 contract.
type Owlto20DepositIterator struct {
	Event *Owlto20Deposit // Event containing the contract specifics and raw log

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
func (it *Owlto20DepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Owlto20Deposit)
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
		it.Event = new(Owlto20Deposit)
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
func (it *Owlto20DepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Owlto20DepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Owlto20Deposit represents a Deposit event raised by the Owlto20 contract.
type Owlto20Deposit struct {
	User      common.Address
	Token     common.Address
	Maker     common.Address
	Target    string
	Amount    *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x673a534e56ef22312f97f00524e3ab12066b624575e63f01a9b579ce40cffac9.
//
// Solidity: event Deposit(address indexed user, address indexed token, address indexed maker, string target, uint256 amount, uint256 timestamp)
func (_Owlto20 *Owlto20Filterer) FilterDeposit(opts *bind.FilterOpts, user []common.Address, token []common.Address, maker []common.Address) (*Owlto20DepositIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}

	logs, sub, err := _Owlto20.contract.FilterLogs(opts, "Deposit", userRule, tokenRule, makerRule)
	if err != nil {
		return nil, err
	}
	return &Owlto20DepositIterator{contract: _Owlto20.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x673a534e56ef22312f97f00524e3ab12066b624575e63f01a9b579ce40cffac9.
//
// Solidity: event Deposit(address indexed user, address indexed token, address indexed maker, string target, uint256 amount, uint256 timestamp)
func (_Owlto20 *Owlto20Filterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *Owlto20Deposit, user []common.Address, token []common.Address, maker []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}

	logs, sub, err := _Owlto20.contract.WatchLogs(opts, "Deposit", userRule, tokenRule, makerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Owlto20Deposit)
				if err := _Owlto20.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0x673a534e56ef22312f97f00524e3ab12066b624575e63f01a9b579ce40cffac9.
//
// Solidity: event Deposit(address indexed user, address indexed token, address indexed maker, string target, uint256 amount, uint256 timestamp)
func (_Owlto20 *Owlto20Filterer) ParseDeposit(log types.Log) (*Owlto20Deposit, error) {
	event := new(Owlto20Deposit)
	if err := _Owlto20.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
