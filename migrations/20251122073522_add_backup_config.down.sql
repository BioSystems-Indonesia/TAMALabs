-- Remove backup configuration entries
DELETE FROM configs WHERE id = 'BackupScheduleType';

DELETE FROM configs WHERE id = 'BackupInterval';

DELETE FROM configs WHERE id = 'BackupTime';