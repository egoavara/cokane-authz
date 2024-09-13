-- +goose Up
ALTER TABLE tuple ADD COLUMN condition_name TEXT, ADD COLUMN condition_context BYTEA;
ALTER TABLE changelog ADD COLUMN condition_name TEXT, ADD COLUMN condition_context BYTEA;