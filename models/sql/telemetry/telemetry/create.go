package telemetry

import (
	"fmt"
	"time"

	"github.com/zk101/nixie/lib/storage"
)

// Create stores a telemetry model into the database
func (m *Model) Create(sqlStore storage.SQL) error {
	insertQuery := "INSERT INTO `%s` (telemetry_user_string, telemetry_client_string, telemetry_client_version, telemetry_date, telemetry_data, telemetry_status) VALUES ('%s', '%s', '%s', '%d', '%s', '1')"
	tableName := time.Unix(m.Date, 0).UTC().Format("200601021500") + "_telemetry"

	if _, err := sqlStore.Exec(fmt.Sprintf(insertQuery, tableName, m.User, m.Client, m.Version, m.Date, m.Data)); err != nil {
		return err
	}

	return nil
}

// EOF
