//go:build integration
// +build integration

package persistence

import (
	"github.com/c0llinn/ebook-store/internal/common"
	config2 "github.com/c0llinn/ebook-store/internal/config"
	"github.com/c0llinn/ebook-store/internal/shop"
	"github.com/c0llinn/ebook-store/test"
	"github.com/c0llinn/ebook-store/test/factory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	repo OrderRepository
}

func (s *OrderRepositoryTestSuite) SetupTest() {
	test.SetEnvironmentVariables()
	config2.InitLogger()
	config2.LoadMigrations("file:../../../migrations")

	conn := config2.NewConnection()
	s.repo = OrderRepository{conn}
}

func TestOrderRepositoryRun(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (s *OrderRepositoryTestSuite) TearDownTest() {
	s.repo.db.Delete(&shop.Order{}, "1 = 1")
}

func (s *OrderRepositoryTestSuite) TestFindByQuery_WithEmptyQuery() {
	order1 := factory.NewOrder()
	order2 := factory.NewOrder()
	order3 := factory.NewOrder()

	err := s.repo.Create(&order1)
	require.Nil(s.T(), err)

	err = s.repo.Create(&order2)
	require.Nil(s.T(), err)

	err = s.repo.Create(&order3)
	require.Nil(s.T(), err)

	actual, err := s.repo.FindByQuery(shop.OrderQuery{})
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), 0, actual.Limit)
	assert.Equal(s.T(), 0, actual.Offset)
	assert.Equal(s.T(), int64(3), actual.TotalOrders)
	assert.Equal(s.T(), order3.ID, actual.Orders[0].ID)
	assert.Equal(s.T(), order2.ID, actual.Orders[1].ID)
	assert.Equal(s.T(), order1.ID, actual.Orders[2].ID)
}

func (s *OrderRepositoryTestSuite) TestFindByQuery_WithLimit() {
	order1 := factory.NewOrder()
	order2 := factory.NewOrder()
	order3 := factory.NewOrder()

	err := s.repo.Create(&order1)
	require.Nil(s.T(), err)

	err = s.repo.Create(&order2)
	require.Nil(s.T(), err)

	err = s.repo.Create(&order3)
	require.Nil(s.T(), err)

	actual, err := s.repo.FindByQuery(shop.OrderQuery{Limit: 1})
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), 1, actual.Limit)
	assert.Equal(s.T(), 0, actual.Offset)
	assert.Equal(s.T(), int64(3), actual.TotalOrders)
	assert.Equal(s.T(), order3.ID, actual.Orders[0].ID)
	assert.Len(s.T(), actual.Orders, 1)
}

func (s *OrderRepositoryTestSuite) TestFindByQuery_WithOffset() {
	order1 := factory.NewOrder()
	order2 := factory.NewOrder()
	order3 := factory.NewOrder()

	err := s.repo.Create(&order1)
	require.Nil(s.T(), err)

	err = s.repo.Create(&order2)
	require.Nil(s.T(), err)

	err = s.repo.Create(&order3)
	require.Nil(s.T(), err)

	actual, err := s.repo.FindByQuery(shop.OrderQuery{Offset: 1})
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), 0, actual.Limit)
	assert.Equal(s.T(), 1, actual.Offset)
	assert.Equal(s.T(), int64(3), actual.TotalOrders)
	assert.Equal(s.T(), order2.ID, actual.Orders[0].ID)
	assert.Equal(s.T(), order1.ID, actual.Orders[1].ID)
	assert.Len(s.T(), actual.Orders, 2)
}

func (s *OrderRepositoryTestSuite) TestFindByQuery_WithStatus() {
	order1 := factory.NewOrder()
	order1.Status = shop.Paid
	order2 := factory.NewOrder()
	order3 := factory.NewOrder()

	err := s.repo.Create(&order1)
	require.Nil(s.T(), err)

	err = s.repo.Create(&order2)
	require.Nil(s.T(), err)

	err = s.repo.Create(&order3)
	require.Nil(s.T(), err)

	actual, err := s.repo.FindByQuery(shop.OrderQuery{Status: shop.Paid})
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), 0, actual.Limit)
	assert.Equal(s.T(), 0, actual.Offset)
	assert.Equal(s.T(), int64(1), actual.TotalOrders)
	assert.Equal(s.T(), order1.ID, actual.Orders[0].ID)
	assert.Len(s.T(), actual.Orders, 1)
}

func (s *OrderRepositoryTestSuite) TestFindByQuery_WithBookID() {
	order1 := factory.NewOrder()
	order2 := factory.NewOrder()
	order3 := factory.NewOrder()

	err := s.repo.Create(&order1)
	require.Nil(s.T(), err)

	err = s.repo.Create(&order2)
	require.Nil(s.T(), err)

	err = s.repo.Create(&order3)
	require.Nil(s.T(), err)

	actual, err := s.repo.FindByQuery(shop.OrderQuery{BookID: order1.BookID})
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), 0, actual.Limit)
	assert.Equal(s.T(), 0, actual.Offset)
	assert.Equal(s.T(), int64(1), actual.TotalOrders)
	assert.Equal(s.T(), order1.ID, actual.Orders[0].ID)
	assert.Len(s.T(), actual.Orders, 1)
}

func (s *OrderRepositoryTestSuite) TestFindByQuery_WithUserID() {
	order1 := factory.NewOrder()
	order2 := factory.NewOrder()
	order3 := factory.NewOrder()

	err := s.repo.Create(&order1)
	require.Nil(s.T(), err)

	err = s.repo.Create(&order2)
	require.Nil(s.T(), err)

	err = s.repo.Create(&order3)
	require.Nil(s.T(), err)

	actual, err := s.repo.FindByQuery(shop.OrderQuery{UserID: order1.UserID})
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), 0, actual.Limit)
	assert.Equal(s.T(), 0, actual.Offset)
	assert.Equal(s.T(), int64(1), actual.TotalOrders)
	assert.Equal(s.T(), order1.ID, actual.Orders[0].ID)
	assert.Len(s.T(), actual.Orders, 1)
}

func (s *OrderRepositoryTestSuite) TestFindByID_Successfully() {
	order := factory.NewOrder()

	err := s.repo.Create(&order)
	require.Nil(s.T(), err)

	actual, err := s.repo.FindByID(order.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), order.ID, actual.ID)
}

func (s *OrderRepositoryTestSuite) TestFindByID_WithError() {
	_, err := s.repo.FindByID(uuid.NewString())

	assert.IsType(s.T(), &common.ErrEntityNotFound{}, err)
}

func (s *OrderRepositoryTestSuite) TestCreate_Successfully() {
	order := factory.NewOrder()

	err := s.repo.Create(&order)
	assert.Nil(s.T(), err)

	persisted, err := s.repo.FindByID(order.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), order.ID, persisted.ID)
}

func (s *OrderRepositoryTestSuite) TestCreate_WithError() {
	order := factory.NewOrder()

	err := s.repo.Create(&order)
	assert.Nil(s.T(), err)

	err = s.repo.Create(&order)
	assert.NotNil(s.T(), err)
}

func (s *OrderRepositoryTestSuite) TestUpdate_Successfully() {
	order := factory.NewOrder()

	err := s.repo.Create(&order)
	assert.Nil(s.T(), err)

	order.Status = shop.Paid
	err = s.repo.Update(&order)
	assert.Nil(s.T(), err)

	persisted, err := s.repo.FindByID(order.ID)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), shop.Paid, persisted.Status)
}