package services

import (
	"context"
	"database/sql"
	"fmt"
	"gugcp/database"
	"gugcp/models"
)

func SaveRedeemCode(ctx context.Context, input models.RedeemInput) (models.Redeem, error) {
	tx, err := database.DB.BeginTx(ctx, nil)

	if err != nil {
		return models.Redeem{}, fmt.Errorf("error when creating transaction: %v", err)
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO wpone_prakerja_redeems(user_id,state,redeem_code,sequence,status) VALUES (?,?,?,?,?)",
		input.UserID, input.State, input.RedeemCode, input.Sequence, 0,
	)

	if err != nil {
		return models.Redeem{}, fmt.Errorf("error when saving redeem code: %v", err)
	}

	var redeem models.Redeem

	result := tx.QueryRowContext(ctx, "SELECT * FROM wpone_prakerja_redeems WHERE state = ?", input.State)

	if err := result.Scan(&redeem.ID, &redeem.UserID, &redeem.State, &redeem.RedeemCode, &redeem.Sequence, &redeem.Status); err != nil {
		if err == sql.ErrNoRows {
			return models.Redeem{}, fmt.Errorf("redeem is not exists: %v", err)
		}
		return models.Redeem{}, fmt.Errorf("error when getting redeem: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return models.Redeem{}, fmt.Errorf("error when starting transaction: %v", err)
	}

	return redeem, nil
}

func GetRedeemCode() {

}

func CheckAttendanceStatus() {

}
