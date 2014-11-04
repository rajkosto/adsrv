package ads

import (
	"database/sql"
)

func GetZoneByName(db *sql.DB, zoneName string) (zone ZoneElem, err error) {
	var rows *sql.Rows
	var exists bool

	rows, err = db.Query("SELECT zoneId, name FROM zones WHERE name = ?", zoneName)
	defer rows.Close()
	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(&zone.Id, &zone.Name)
		if err != nil {
			return
		}
		exists = true
	}
	err = rows.Err()
	if err != nil {
		return
	}
	rows.Close()

	if exists == false {
		var res sql.Result
		res, err = db.Exec("INSERT INTO zones (name) VALUES (?)", zoneName)
		if err != nil {
			return
		}

		var bigId int64
		bigId, err = res.LastInsertId()
		if err != nil {
			return
		}

		zone.Id = uint32(bigId)
		zone.Name = zoneName
	} else {
		rows, err = db.Query("SELECT targetId FROM targets WHERE zoneId = ?", zone.Id)
		defer rows.Close()
		if err != nil {
			return
		}

		zone.Targets = make([]uint32, 0, 40)
		for rows.Next() {
			var targId uint32
			err = rows.Scan(&targId)
			if err != nil {
				return
			}
			zone.Targets = append(zone.Targets, targId)
		}
		err = rows.Err()
		if err != nil {
			return
		}
		rows.Close()
	}

	return
}

func AddZoneVisit(db *sql.DB, sessionId uint32, zone *ZoneElem, timestamp uint64) error {
	_, err := db.Exec("INSERT INTO zone_visits (sessionId,zoneId,timestampMs) VALUES (?,?,?)", sessionId, zone.Id, timestamp)
	return err
}
