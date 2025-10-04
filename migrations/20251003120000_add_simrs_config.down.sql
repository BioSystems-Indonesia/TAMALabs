-- remove SIMRS integration config keys
DELETE FROM configs WHERE id = 'SimrsIntegrationEnabled';

DELETE FROM configs WHERE id = 'SimrsDatabaseDSN';