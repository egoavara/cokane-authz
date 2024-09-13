-- +goose Up
ALTER TABLE authorization_model ADD COLUMN schema_version TEXT NOT NULL DEFAULT '1.0';