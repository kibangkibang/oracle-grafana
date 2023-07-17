package plugin

import (
	"database/sql"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	go_ora "github.com/sijms/go-ora/v2"
)

type OracleDatasourceConnection struct {
	connection *sql.DB
}

func (c *OracleDatasourceConnection) Connect(settings *OracleDatasourceSettings) error {
	var err error
	if !c.IsConnected() {
		urlOptions := map[string]string{}
		if len(settings.O_sid) > 0 {
			urlOptions["SID"] = settings.O_sid
		}
		connectionString := go_ora.BuildUrl(settings.O_hostname, settings.O_port, settings.O_service, settings.O_user, settings.O_password, urlOptions)
		log.DefaultLogger.Debug(connectionString)

		connection, conErr := sql.Open("oracle", connectionString)
		if conErr != nil {
			err = conErr
			log.DefaultLogger.Error("Error connecting to Oracle: ", err)
		} else {
			c.connection = connection
			err = c.Ping()
		}
	}
	return err
}

func (c *OracleDatasourceConnection) Disconnect() error {
	var err error
	if c.IsConnected() {
		err = c.connection.Close()
		if err != nil {
			log.DefaultLogger.Error("Error closing Oracle connection: ", err)
		}
		c.connection = nil
	}
	return err
}

func (c *OracleDatasourceConnection) IsConnected() bool {
	if c.connection != nil {
		return true
	} else {
		return false
	}
}

func (c *OracleDatasourceConnection) Ping() error {
	var err error
	if c.IsConnected() {
		err = c.connection.Ping()
		if err != nil {
			log.DefaultLogger.Error("Error pingging Oracle connection: ", err)
		}
	}
	return err
}

func (c *OracleDatasourceConnection) Reconnect(settings *OracleDatasourceSettings) error {
	var err error
	if c.IsConnected() {
		err = c.Disconnect()
	}
	err = c.Connect(settings)
	return err
}
