package services

import (
	"gugcp/database"
	"gugcp/models"
)

func SaveRedeemCode(input models.RedeemInput) {
	query, err := database.DB.Prepare("INSERT INTO wpone_prakerja_redeems(user_id,state,redeem_code,sequence,status) VALUES (?,?,?,?,?)")

	if err != nil {
		panic(err)
	}

	res, err := query.Exec(input.UserID, input.State, input.RedeemCode, input.Sequence, input.Status)

	if err != nil {
		panic(err)
	}

	lastID, err := res.LastInsertId()

	if err != nil || lastID == 0 {
		panic(err)
	}
}

func GetRedeemCode() {

}

func CheckAttendanceStatus() {

}
