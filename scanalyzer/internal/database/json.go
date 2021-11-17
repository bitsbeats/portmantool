package database

import (
	"encoding/json"
)

type dict map[string]interface{}

func (d Diff) MarshalJSON() ([]byte, error) {
	return json.Marshal(dict{
		"address":        d.Address,
		"port":           d.Port,
		"protocol":       d.Protocol,
		"expected_state": d.ExpectedState,
		"actual_state":   d.ActualState,
		"scan_id":        d.ScanID,
		"comment":        d.Comment,
	})
}

func (d DiffAB) MarshalJSON() ([]byte, error) {
	return json.Marshal(dict{
		"address":  d.Address,
		"port":     d.Port,
		"protocol": d.Protocol,
		"state_a":  d.StateA,
		"state_b":  d.StateB,
		"scan_a":   d.ScanA,
		"scan_b":   d.ScanB,
	})
}

func (s ActualState) MarshalJSON() ([]byte, error) {
	return json.Marshal(dict{
		"address":  s.Address,
		"port":     s.Port,
		"protocol": s.Protocol,
		"state":    s.State,
		"scan_id":  s.ScanID.Unix(),
	})
}

func (s ExpectedState) MarshalJSON() ([]byte, error) {
	return json.Marshal(dict{
		"address":  s.Address,
		"port":     s.Port,
		"protocol": s.Protocol,
		"state":    s.State,
		"comment":  s.Comment,
	})
}

func (s Scan) MarshalJSON() ([]byte, error) {
	return json.Marshal(dict{"id": s.ID.Unix()})
}
