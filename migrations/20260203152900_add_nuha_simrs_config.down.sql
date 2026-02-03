-- Remove Nuha SIMRS configuration entries
DELETE FROM configs WHERE id = 'NuhaIntegrationEnabled';

DELETE FROM configs WHERE id = 'NuhaBaseURL';

DELETE FROM configs WHERE id = 'NuhaSessionID';
