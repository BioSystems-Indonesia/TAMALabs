-- Add backup configuration entries
INSERT INTO
    configs (id, value)
VALUES ('BackupScheduleType', 'daily');

INSERT INTO configs (id, value) VALUES ('BackupInterval', '6');

INSERT INTO configs (id, value) VALUES ('BackupTime', '02:00');