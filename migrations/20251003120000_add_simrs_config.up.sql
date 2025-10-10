-- add SIMRS integration config keys with default values
INSERT OR IGNORE INTO
    configs (id, value)
VALUES (
        'SimrsIntegrationEnabled',
        'false'
    );

INSERT OR IGNORE INTO
    configs (id, value)
VALUES (
        'SimrsDatabaseDSN',
        'root:secret@tcp(localhost:3306)/simrs_db?charset=utf8mb4&parseTime=True&loc=Local'
    );