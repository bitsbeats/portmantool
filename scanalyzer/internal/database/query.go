package database

import (
	"gorm.io/gorm"
)

func CurrentState(db *gorm.DB) (state []ActualState, err error) {
	latestScan := db.Model(&ActualState{}).Select("address, port, protocol, MAX(scan_id) max_scan_id").Group("address, port, protocol")
	err = db.Model(&ActualState{}).Joins("JOIN (?) latest_states ON actual_states.address = latest_states.address AND actual_states.port = latest_states.port AND actual_states.protocol = latest_states.protocol AND scan_id = max_scan_id", latestScan).Scan(&state).Error
	if err != nil {
		return nil, err
	}

	return state, nil
}
