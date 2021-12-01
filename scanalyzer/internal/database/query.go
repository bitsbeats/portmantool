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
	err = db.Table("expected_states a").Select("address, port, protocol, a.state expected_state, b.state actual_state, scan_id, comment").Joins("FULL JOIN (?) b USING (address, port, protocol) WHERE a.state IS DISTINCT FROM b.state", currentState(db)).Scan(&diff).Error
	if err != nil {
		return nil, err
	}

	return diff, nil
}

func DiffOne(db *gorm.DB, id time.Time) (diff []Diff, err error) {
	state := db.Model(&ActualState{}).Where(&ActualState{ScanID: id})
	err = db.Table("expected_states a").Select("address, port, protocol, a.state expected_state, b.state actual_state, scan_id, comment").Joins("FULL JOIN (?) b USING (address, port, protocol) WHERE a.state IS DISTINCT FROM b.state", state).Scan(&diff).Error
	if err != nil {
		return nil, err
	}

	return diff, nil
}

func DiffTwo(db *gorm.DB, id1, id2 time.Time) (diff []DiffAB, err error) {
	state1 := db.Model(&ActualState{}).Where(&ActualState{ScanID: id1})
	state2 := db.Model(&ActualState{}).Where(&ActualState{ScanID: id2})
	err = db.Table("(?) a", state1).Select("address, port, protocol, a.state state_a, b.state state_b, a.scan_id scan_a, b.scan_id scan_b").Joins("FULL JOIN (?) b USING (address, port, protocol) WHERE a.state IS DISTINCT FROM b.state", state2).Scan(&diff).Error
	if err != nil {
		return nil, err
	}

	return diff, nil
}

func Expected(db *gorm.DB) (state []ExpectedState, err error) {
	err = db.Find(&state).Error
	if err != nil {
		return nil, err
	}

	return state, nil
}

func Prune(db *gorm.DB, keep time.Time) error {
	current := db.Model(&ActualState{}).Select("MAX(scan_id)").Group("address, port, protocol")
	err := db.Delete(&Scan{}, "id NOT IN (?) AND id < ?", current, keep).Error
	if err != nil {
		return err
	}

	err = db.Delete(&ActualState{}, "ROW(address, port, protocol, scan_id) NOT IN (?) AND scan_id < ?", current.Select("address, port, protocol, MAX(scan_id)"), keep).Error
	if err != nil {
		return err
	}

	return nil
}

func Scans(db *gorm.DB) (scans []Scan, err error) {
	err = db.Find(&scans).Error
	if err != nil {
		return nil, err
	}

	return scans, nil
}

func StateAt(db *gorm.DB, id time.Time) (state []ActualState, err error) {
	err = db.Where(&ActualState{ScanID: id}).Find(&state, id).Error
	if err != nil {
		return nil, err
	}

	return state, nil
}

func currentState(db *gorm.DB) *gorm.DB {
	latestScans := db.Model(&ActualState{}).Select("address addr, port p, protocol proto, MAX(scan_id) max_scan_id").Group("addr, p, proto")
	return db.Model(&ActualState{}).Joins("JOIN (?) latest_scans ON address = addr AND port = p AND protocol = proto AND scan_id = max_scan_id", latestScans)
}
