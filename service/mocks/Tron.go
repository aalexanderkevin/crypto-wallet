// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	context "context"

	api "github.com/fbsobreira/gotron-sdk/pkg/proto/api"

	core "github.com/fbsobreira/gotron-sdk/pkg/proto/core"

	mock "github.com/stretchr/testify/mock"

	model "github.com/aalexanderkevin/crypto-wallet/model"

	service "github.com/aalexanderkevin/crypto-wallet/service"
)

// Tron is an autogenerated mock type for the Tron type
type Tron struct {
	mock.Mock
}

// CheckAddress provides a mock function with given fields: address
func (_m *Tron) CheckAddress(address string) error {
	ret := _m.Called(address)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(address)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Close provides a mock function with given fields:
func (_m *Tron) Close() {
	_m.Called()
}

// GetBalance provides a mock function with given fields: ctx, address
func (_m *Tron) GetBalance(ctx context.Context, address *string) (*int64, error) {
	ret := _m.Called(ctx, address)

	var r0 *int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *string) (*int64, error)); ok {
		return rf(ctx, address)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string) *int64); ok {
		r0 = rf(ctx, address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*int64)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string) error); ok {
		r1 = rf(ctx, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetConfirmedTxAddress provides a mock function with given fields: ctx, address
func (_m *Tron) GetConfirmedTxAddress(ctx context.Context, address *string) (*service.GetTransactionResponse, error) {
	ret := _m.Called(ctx, address)

	var r0 *service.GetTransactionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *string) (*service.GetTransactionResponse, error)); ok {
		return rf(ctx, address)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string) *service.GetTransactionResponse); ok {
		r0 = rf(ctx, address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*service.GetTransactionResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string) error); ok {
		r1 = rf(ctx, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTx provides a mock function with given fields: ctx, txhash
func (_m *Tron) GetTx(ctx context.Context, txhash string) (*core.TransactionInfo, error) {
	ret := _m.Called(ctx, txhash)

	var r0 *core.TransactionInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*core.TransactionInfo, error)); ok {
		return rf(ctx, txhash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *core.TransactionInfo); ok {
		r0 = rf(ctx, txhash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.TransactionInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, txhash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTxByAccountAddress provides a mock function with given fields: ctx, address, filter
func (_m *Tron) GetTxByAccountAddress(ctx context.Context, address *string, filter *service.GetTxByAccountAddressFilter) (*service.GetTransactionResponse, error) {
	ret := _m.Called(ctx, address, filter)

	var r0 *service.GetTransactionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *string, *service.GetTxByAccountAddressFilter) (*service.GetTransactionResponse, error)); ok {
		return rf(ctx, address, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string, *service.GetTxByAccountAddressFilter) *service.GetTransactionResponse); ok {
		r0 = rf(ctx, address, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*service.GetTransactionResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string, *service.GetTxByAccountAddressFilter) error); ok {
		r1 = rf(ctx, address, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUnconfirmedTxAddress provides a mock function with given fields: ctx, address
func (_m *Tron) GetUnconfirmedTxAddress(ctx context.Context, address *string) (*service.GetTransactionResponse, error) {
	ret := _m.Called(ctx, address)

	var r0 *service.GetTransactionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *string) (*service.GetTransactionResponse, error)); ok {
		return rf(ctx, address)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string) *service.GetTransactionResponse); ok {
		r0 = rf(ctx, address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*service.GetTransactionResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string) error); ok {
		r1 = rf(ctx, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWallet provides a mock function with given fields: ctx, seedPhrase
func (_m *Tron) GetWallet(ctx context.Context, seedPhrase *string) *model.TrxHdWallet {
	ret := _m.Called(ctx, seedPhrase)

	var r0 *model.TrxHdWallet
	if rf, ok := ret.Get(0).(func(context.Context, *string) *model.TrxHdWallet); ok {
		r0 = rf(ctx, seedPhrase)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TrxHdWallet)
		}
	}

	return r0
}

// SendTx provides a mock function with given fields: ctx, txOpts, wallet
func (_m *Tron) SendTx(ctx context.Context, txOpts *model.TxOpts, wallet *model.TrxHdWallet) (*api.TransactionExtention, error) {
	ret := _m.Called(ctx, txOpts, wallet)

	var r0 *api.TransactionExtention
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.TxOpts, *model.TrxHdWallet) (*api.TransactionExtention, error)); ok {
		return rf(ctx, txOpts, wallet)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.TxOpts, *model.TrxHdWallet) *api.TransactionExtention); ok {
		r0 = rf(ctx, txOpts, wallet)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.TransactionExtention)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.TxOpts, *model.TrxHdWallet) error); ok {
		r1 = rf(ctx, txOpts, wallet)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTron creates a new instance of Tron. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTron(t interface {
	mock.TestingT
	Cleanup(func())
}) *Tron {
	mock := &Tron{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
