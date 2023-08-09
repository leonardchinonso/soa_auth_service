package datasource

import (
    "log"

    "github.com/leonardchinonso/auth_service_cmp7174/config"
)

type DataSource struct {
    *DatabaseContext
    Cfg *map[string]string
}

// InitDataSource initializes the data source
func InitDataSource() (*DataSource, error) {
    log.Println("Initializing Data Sources...")

    // initialize the config
    configMap, err := config.InitConfig()
    if err != nil {
        return nil, err
    }

    // initialize the database
    dbCtx, err := InitDB(configMap)
    if err != nil {
        return nil, err
    }

    return &DataSource{
        DatabaseContext: dbCtx,
        Cfg: configMap,
    }, nil
}

// Close closes all the data sources and releases resources held
func (ds *DataSource) Close() {
    // cancel the context after the database client has closed its connection
    defer ds.CancelFunc()

    defer func() {
        // disconnect from the client
        if err := ds.Client.Disconnect(ds.Ctx); err != nil {
            log.Printf("Error trying to close data sources correctly: %v", err)
            panic(err)
        }
    }()
}
