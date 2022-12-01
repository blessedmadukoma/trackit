package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// createRandomBudget creates a random budget in the database
func createRandomBudget(t *testing.T) Budget {
	user := createRandomUser(t)

	arg := CreateBudgetParams{
		UserID:        user.ID,
		InitialAmount: 100000,
		CurrentAmount: 0,
		BudgetName:    "New Budget",
		Description:   "This is the new budget",
		StartDate:     time.Now().UTC(),
		EndDate:       time.Now().UTC().Add(100 * time.Hour),
	}

	budget, err := testQueries.CreateBudget(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, budget)

	require.Equal(t, arg.BudgetName, budget.BudgetName)
	require.Equal(t, arg.InitialAmount, budget.InitialAmount)
	require.Equal(t, arg.Description, budget.Description)
	require.Equal(t, arg.UserID, budget.UserID)

	require.NotZero(t, budget.ID)
	require.NotZero(t, budget.CreatedAt)

	return budget
}

// TestCreateBudget tests the createRandomBudget method which creates a new budget
func TestCreateBudget(t *testing.T) {
	createRandomBudget(t)
}

// TestGetBudgetByID gets the budget information through the ID
func TestGetBudgetByID(t *testing.T) {
	budget1 := createRandomBudget(t)

	budget2, err := testQueries.GetBudgetByID(context.Background(), budget1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, budget2)

	require.Equal(t, budget1.ID, budget2.ID)
	require.Equal(t, budget1.BudgetName, budget2.BudgetName)
	require.Equal(t, budget1.InitialAmount, budget2.InitialAmount)
	require.Equal(t, budget1.Description, budget2.Description)
	require.Equal(t, budget1.UserID, budget2.UserID)

	require.WithinDuration(t, budget1.CreatedAt, budget2.CreatedAt, time.Second)
}

// TestGetBudgetByUserID gets the budget information through the User ID
func TestGetBudgetByUserID(t *testing.T) {
	budget1 := createRandomBudget(t)

	budget2, err := testQueries.GetBudgetByUserID(context.Background(), budget1.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, budget2)

	require.Equal(t, budget1.ID, budget2.ID)
	require.Equal(t, budget1.BudgetName, budget2.BudgetName)
	require.Equal(t, budget1.InitialAmount, budget2.InitialAmount)
	require.Equal(t, budget1.Description, budget2.Description)
	require.Equal(t, budget1.StartDate, budget2.StartDate)
	require.Equal(t, budget1.EndDate, budget2.EndDate)
	require.Equal(t, budget1.UserID, budget2.UserID)

	require.WithinDuration(t, budget1.CreatedAt, budget2.CreatedAt, time.Second)
}

// TestListBudgets lists all the budgets
func TestListBudgets(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomBudget(t)
	}

	arg := ListBudgetsParams{
		Limit:  5,
		Offset: 5,
	}

	budgets, err := testQueries.ListBudgets(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, budgets, 5)

	for _, budget := range budgets {
		require.NotEmpty(t, budget)
	}
}

// TestUpdateBudget tests the UpdateBudget method which updates the details of the budget
func TestUpdateBudget(t *testing.T) {
	budget1 := createRandomBudget(t)

	arg := UpdateBudgetParams{
		ID:            budget1.ID,
		UserID: budget1.UserID,
		// InitialAmount: 10,
		CurrentAmount: 20000,
		BudgetName:    "New Budget",
		Description:   "This is the new budget",
		StartDate:     time.Now().UTC(),
		EndDate:       time.Now().UTC().Add(100 * time.Hour),
	}

	budget2, err := testQueries.UpdateBudget(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, budget2)

	require.Equal(t, budget1.ID, budget2.ID)
	require.Equal(t, budget1.BudgetName, budget2.BudgetName)
	// require.Equal(t, budget1.InitialAmount, budget2.InitialAmount)
	require.Equal(t, budget1.Description, budget2.Description)
	// require.Equal(t, budget1.StartDate, budget2.StartDate)
	// require.Equal(t, budget1.EndDate, budget2.EndDate)
	require.Equal(t, budget1.UserID, budget2.UserID)

	require.WithinDuration(t, budget1.CreatedAt, budget2.CreatedAt, time.Second)
}
