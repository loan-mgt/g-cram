-- +migrate Up
ALTER TABLE users DROP COLUMN websocket;

-- +migrate Down
ALTER TABLE users ADD COLUMN websocket TEXT;
