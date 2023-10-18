// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	context "context"
	big "math/big"

	gobcy "github.com/blockcypher/gobcy/v2"

	mock "github.com/stretchr/testify/mock"

	model "github.com/aalexanderkevin/crypto-wallet/model"
)

// Bitcoin is an autogenerated mock type for the Bitcoin type
type Bitcoin struct {
	mock.Mock
}

// CheckAddress provides a mock function with given fields: address
func (_m *Bitcoin) CheckAddress(address *string) bool {
	ret := _m.Called(address)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*string) bool); ok {
		r0 = rf(address)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// CreateWebhookConfirmedTx provides a mock function with given fields: ctx, address
func (_m *Bitcoin) CreateWebhookConfirmedTx(ctx context.Context, address *string) (*gobcy.Hook, error) {
	ret := _m.Called(ctx, address)

	var r0 *gobcy.Hook
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *string) (*gobcy.Hook, error)); ok {
		return rf(ctx, address)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string) *gobcy.Hook); ok {
		r0 = rf(ctx, address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gobcy.Hook)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string) error); ok {
		r1 = rf(ctx, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteWebhook provides a mock function with given fields: ctx, id
func (_m *Bitcoin) DeleteWebhook(ctx context.Context, id *string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetBalance provides a mock function with given fields: ctx, address
func (_m *Bitcoin) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	ret := _m.Called(ctx, address)

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*big.Int, error)); ok {
		return rf(ctx, address)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *big.Int); ok {
		r0 = rf(ctx, address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTx provides a mock function with given fields: ctx, txhash
func (_m *Bitcoin) GetTx(ctx context.Context, txhash string) (*gobcy.TX, error) {
	ret := _m.Called(ctx, txhash)

	var r0 *gobcy.TX
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*gobcy.TX, error)); ok {
		return rf(ctx, txhash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *gobcy.TX); ok {
		r0 = rf(ctx, txhash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gobcy.TX)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, txhash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWallet provides a mock function with given fields: ctx, seedPhrase
func (_m *Bitcoin) GetWallet(ctx context.Context, seedPhrase *string) (*model.BtcHdWallet, error) {
	ret := _m.Called(ctx, seedPhrase)

	var r0 *model.BtcHdWallet
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *string) (*model.BtcHdWallet, error)); ok {
		return rf(ctx, seedPhrase)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string) *model.BtcHdWallet); ok {
		r0 = rf(ctx, seedPhrase)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.BtcHdWallet)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string) error); ok {
		r1 = rf(ctx, seedPhrase)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendTx provides a mock function with given fields: ctx, wallet, txOpts
func (_m *Bitcoin) SendTx(ctx context.Context, wallet *model.BtcHdWallet, txOpts *model.TxOpts) (*model.Transaction, error) {
	ret := _m.Called(ctx, wallet, txOpts)

	var r0 *model.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.BtcHdWallet, *model.TxOpts) (*model.Transaction, error)); ok {
		return rf(ctx, wallet, txOpts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.BtcHdWallet, *model.TxOpts) *model.Transaction); ok {
		r0 = rf(ctx, wallet, txOpts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.BtcHdWallet, *model.TxOpts) error); ok {
		r1 = rf(ctx, wallet, txOpts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewBitcoin creates a new instance of Bitcoin. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBitcoin(t interface {
	mock.TestingT
	Cleanup(func())
}) *Bitcoin {
	mock := &Bitcoin{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
