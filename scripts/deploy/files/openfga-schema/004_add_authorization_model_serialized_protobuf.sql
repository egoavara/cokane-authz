-- +goose Up
ALTER TABLE authorization_model ADD COLUMN serialized_protobuf BYTEA;