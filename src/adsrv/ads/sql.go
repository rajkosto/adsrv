package ads

import (
	"database/sql"
)

func GetZoneByName(db *sql.DB, zoneName string) (zone ZoneElem, err error) {
	var rows *sql.Rows
	var exists bool

	rows, err = db.Query("SELECT zoneId, Name FROM zones WHERE Name = ?", zoneName)
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
		res, err = db.Exec("INSERT INTO zones (Name) VALUES (?)", zoneName)
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

func GetTargetsByZoneId(db *sql.DB, zoneId uint32) (targets []TargetElem, err error) {
	var rows *sql.Rows

	rows, err = db.Query("SELECT targetId, Name FROM targets WHERE zoneId = ?", zoneId)
	defer rows.Close()
	if err != nil {
		return
	}

	targets = make([]TargetElem, 0, 40)
	for rows.Next() {
		var elem TargetElem
		err = rows.Scan(&elem.Id, &elem.Name)
		if err != nil {
			return
		}
		targets = append(targets, elem)
	}
	err = rows.Err()
	if err != nil {
		return
	}
	rows.Close()

	for targetIdx, target := range targets {
		rows, err = db.Query("SELECT mediaId FROM target_medias WHERE targetId = ?", target.Id)
		defer rows.Close()
		if err != nil {
			return
		}

		target.Medias = make([]uint32, 0, 40)
		for rows.Next() {
			var mediaId uint32
			err = rows.Scan(&mediaId)
			if err != nil {
				return
			}
			target.Medias = append(target.Medias, mediaId)
		}
		targets[targetIdx] = target
		err = rows.Err()
		if err != nil {
			return
		}
		rows.Close()
	}

	return
}

func GetMediaById(db *sql.DB, mediaId uint32) (media MediaElem, err error) {
	var rows *sql.Rows

	rows, err = db.Query("SELECT mediaId, Name, mimeType, fileSize, fileMd5, filePath, medias.crexId, minSize, rotDuration, minAngleDeg " +
						"FROM medias INNER JOIN crexs ON medias.crexId = crexs.crexId WHERE mediaId = ? LIMIT 1", mediaId)
	defer rows.Close()
	if err != nil {
		return
	}

	for rows.Next() {
		var crex Crex
		err = rows.Scan(&media.id, &media.name, &media.unkShort, &media.fileSize, &media.fileMd5, &media.filePath, &crex.id, &crex.minSize, &crex.rotDuration, &crex.minAngleDeg)
		media.crex = crex
		if err != nil {
			return
		}
	}
	err = rows.Err()
	if err != nil {
		return
	}
	rows.Close()
	return
}