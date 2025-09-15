package lease

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockDataService struct {
	mock.Mock
}

func (m *mockDataService) List(query *Lease) (*Leases, error) {
	args := m.Called(query)
	return args.Get(0).(*Leases), args.Error(1)
}

func (m *mockDataService) Write(lease *Lease, lastModifiedOn *int64) error {
	args := m.Called(lease, lastModifiedOn)
	return args.Error(0)
}

func (m *mockDataService) Get(leaseID string) (*Lease, error) {
	args := m.Called(leaseID)
	return args.Get(0).(*Lease), args.Error(1)
}

func (m *mockDataService) GetByAccountIDAndPrincipalID(accountID string, principalID string) (*Lease, error) {
	args := m.Called(accountID, principalID)
	return args.Get(0).(*Lease), args.Error(1)
}

func TestListPages_MultiplePages(t *testing.T) {
	mockSvc := &mockDataService{}
	service := &Service{dataSvc: mockSvc}

	firstPage := &Leases{
		{ID: stringPtr("lease1"), AccountID: stringPtr("account1"), PrincipalID: stringPtr("principal1")},
		{ID: stringPtr("lease2"), AccountID: stringPtr("account2"), PrincipalID: stringPtr("principal2")},
	}

	secondPage := &Leases{
		{ID: stringPtr("lease3"), AccountID: stringPtr("account3"), PrincipalID: stringPtr("principal3")},
	}

	query := &Lease{Status: StatusActive.StatusPtr()}

	mockSvc.On("List", mock.MatchedBy(func(q *Lease) bool {
		return q.NextAccountID == nil && q.NextPrincipalID == nil
	})).Return(firstPage, nil).Once().Run(func(args mock.Arguments) {
		q := args.Get(0).(*Lease)
		q.NextAccountID = stringPtr("account2")
		q.NextPrincipalID = stringPtr("principal2")
	})

	mockSvc.On("List", mock.MatchedBy(func(q *Lease) bool {
		return q.NextAccountID != nil && q.NextPrincipalID != nil
	})).Return(secondPage, nil).Once().Run(func(args mock.Arguments) {
		q := args.Get(0).(*Lease)
		q.NextAccountID = nil
		q.NextPrincipalID = nil
	})

	var allLeases []*Lease
	err := service.ListPages(query, func(leases *Leases) bool {
		for _, lease := range *leases {
			l := lease
			allLeases = append(allLeases, &l)
		}
		return true
	})

	assert.NoError(t, err)
	assert.Len(t, allLeases, 3)
	assert.Equal(t, "lease1", *allLeases[0].ID)
	assert.Equal(t, "lease2", *allLeases[1].ID)
	assert.Equal(t, "lease3", *allLeases[2].ID)

	mockSvc.AssertExpectations(t)
}

func TestListPages_SinglePage(t *testing.T) {
	mockSvc := &mockDataService{}
	service := &Service{dataSvc: mockSvc}

	singlePage := &Leases{
		{ID: stringPtr("lease1"), AccountID: stringPtr("account1"), PrincipalID: stringPtr("principal1")},
	}

	query := &Lease{Status: StatusActive.StatusPtr()}

	mockSvc.On("List", mock.MatchedBy(func(q *Lease) bool {
		return q.NextAccountID == nil && q.NextPrincipalID == nil
	})).Return(singlePage, nil).Once().Run(func(args mock.Arguments) {
		q := args.Get(0).(*Lease)
		q.NextAccountID = nil
		q.NextPrincipalID = nil
	})

	var allLeases []*Lease
	err := service.ListPages(query, func(leases *Leases) bool {
		for _, lease := range *leases {
			l := lease
			allLeases = append(allLeases, &l)
		}
		return true
	})

	assert.NoError(t, err)
	assert.Len(t, allLeases, 1)
	assert.Equal(t, "lease1", *allLeases[0].ID)

	mockSvc.AssertExpectations(t)
}

func stringPtr(s string) *string {
	return &s
}
