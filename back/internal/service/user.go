package service

import (
	"context"
	"database/sql"
	"fmt"
	"loan-mgt/g-cram/internal/db"
	"loan-mgt/g-cram/internal/db/sqlc"
)

func UpdateOrCreateUser(ctx context.Context, db *db.Store, id, refreshToken, sha string) error {
	fmt.Println("Updating or creating user... token hash: ", sha, " id: ", id, " token: ", refreshToken)
	_, err := db.GetUser(ctx, id)
	if err != nil {
		if err != sql.ErrNoRows {
			return fmt.Errorf("error getting user: %w", err)
		}

		// Create new user
		arg := sqlc.CreateUserParams{
			ID:    id,
			Token: sql.NullString{String: refreshToken, Valid: true},
			TokenHash: sql.NullString{
				String: sha,
				Valid:  true,
			},
		}
		fmt.Println("create", arg)
		if err = db.CreateUser(ctx, arg); err != nil {
			return fmt.Errorf("error creating user: %w", err)
		}
		return nil
	}

	// Update existing user
	arg := sqlc.UpdateUserTokenParams{
		ID:    id,
		Token: sql.NullString{String: refreshToken, Valid: true},
		TokenHash: sql.NullString{
			String: sha,
			Valid:  true,
		},
	}
	if err = db.UpdateUserToken(ctx, arg); err != nil {
		return fmt.Errorf("error updating user token: %w", err)
	}

	return nil
}

func AddSubscriptionToUser(ctx context.Context, db *db.Store, sub string, subscription string) error {
	arg := sqlc.UpdateUserSubscriptionParams{
		ID:           sub,
		Subscription: sql.NullString{String: subscription, Valid: true},
	}
	if err := db.UpdateUserSubscription(ctx, arg); err != nil {
		return fmt.Errorf("error updating user subscription: %w", err)
	}
	return nil
}

func RemoveSubscriptionFromUser(ctx context.Context, db *db.Store, sub string) error {
	arg := sqlc.UpdateUserSubscriptionParams{
		ID:           sub,
		Subscription: sql.NullString{String: "", Valid: false},
	}
	if err := db.UpdateUserSubscription(ctx, arg); err != nil {
		return fmt.Errorf("error updating user subscription: %w", err)
	}
	return nil
}
