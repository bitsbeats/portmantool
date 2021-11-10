package database

import (
	"time"

	"gorm.io/gorm"
)

type (
	Diff struct {
		Target
		ExpectedState	string
		ActualState	string
		ScanID		time.Time
		Comment		string
	}

	DiffAB struct {
		Target
		StateA	string
		StateB	string
		ScanA	time.Time
		ScanB	time.Time
		Comment	string
	}
)

func CurrentState(db *gorm.DB) (state []ActualState, err error) {
	err = currentState(db).Scan(&state).Error
	if err != nil {
		return nil, err
	}

	return state, nil
}

func DiffExpected(db *gorm.DB) (diff []Diff, err error) {
	err = db.Table("expected_states a").Select("COALESCE(a.address, b.address) address, COALESCE(a.port, b.port) port, COALESCE(a.protocol, b.protocol) protocol, a.state expected_state, b.state actual_state, scan_id, comment").Joins("FULL JOIN (?) b ON a.address = b.address AND a.port = b.port AND a.protocol = b.protocol AND a.state <> b.state", currentState(db)).Scan(&diff).Error
	if err != nil {
		return nil, err
	}

	return diff, nil
}

func Scans(db *gorm.DB) (scans []Scan, err error) {
	err = db.Model(&Scan{}).Find(&scans).Error
	if err != nil {
		return nil, err
	}

	return scans, nil
}

func StateAt(db *gorm.DB, id time.Time) ([]ActualState, error) {
	var scan Scan
	err := db.Joins("ActualState").Find(&scan, id).Error
	if err != nil {
		return nil, err
	}

	return scan.Ports, nil
}

func currentState(db *gorm.DB) *gorm.DB {
	latestScans := db.Model(&ActualState{}).Select("address addr, port p, protocol proto, MAX(scan_id) max_scan_id").Group("addr, p, proto")
	return db.Model(&ActualState{}).Joins("JOIN (?) latest_scans ON address = addr AND port = p AND protocol = proto AND scan_id = max_scan_id", latestScans)
}
