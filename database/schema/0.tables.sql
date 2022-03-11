CREATE TABLE policies (
  policy_id              serial       PRIMARY KEY,
  name                   text         NOT NULL UNIQUE,
  default_value          jsonb        NOT NULL
);

CREATE TABLE global_policies (
  policy_id              serial       NOT NULL UNIQUE REFERENCES policies
                                      ON DELETE CASCADE,
  value                  jsonb        NOT NULL
);

CREATE TABLE clusters (
  cluster_id             serial       PRIMARY KEY,
  name                   text         NOT NULL UNIQUE
);

CREATE TABLE clusters_policies (
  cluster_id             int          NOT NULL REFERENCES clusters
                                      ON DELETE CASCADE,
  policy_id              int          NOT NULL REFERENCES policies
                                      ON DELETE CASCADE,
  value                  jsonb        NOT NULL,

  UNIQUE (cluster_id, policy_id)
);

CREATE TABLE faas_deployments (
  faas_id                serial       PRIMARY KEY,
  cluster_id             int          NOT NULL REFERENCES clusters
                                      ON DELETE CASCADE,
  url                    text         NOT NULL UNIQUE
);

CREATE TABLE storage_deployments (
  storage_id             serial       PRIMARY KEY,
  cluster_id             int          NOT NULL REFERENCES clusters
                                      ON DELETE CASCADE,
  minio_deployment_id    text         NOT NULL UNIQUE,
  alias                  text         NOT NULL UNIQUE,
  endpoint               text         NOT NULL UNIQUE,
  access_key             text         NOT NULL,
  secret_key             text         NOT NULL,
  use_ssl                boolean      NOT NULL,
  sqs_arn                text         NOT NULL,
  management_url         text         NOT NULL DEFAULT ''
);

CREATE TABLE buckets (
  bucket_id              serial       PRIMARY KEY,
  storage_id             int          NOT NULL REFERENCES storage_deployments
                                      ON DELETE CASCADE,
  name                   text         NOT NULL UNIQUE
);

CREATE TABLE buckets_policies (
  bucket_id              int          NOT NULL REFERENCES buckets
                                      ON DELETE CASCADE,
  policy_id              int          NOT NULL REFERENCES policies
                                      ON DELETE CASCADE,
  value                  jsonb        NOT NULL,

  UNIQUE (bucket_id, policy_id)
);

CREATE TABLE replica_bucket_locations (
  bucket_id              int          NOT NULL REFERENCES buckets
                                      ON DELETE CASCADE,
  storage_id             int          NOT NULL REFERENCES storage_deployments
                                      ON DELETE CASCADE,

  UNIQUE (bucket_id, storage_id)
);

CREATE TABLE objects (
  object_id              serial       PRIMARY KEY,
  bucket_id              int          NOT NULL REFERENCES buckets
                                      ON DELETE CASCADE,
  name                   text         NOT NULL,

  UNIQUE (bucket_id, name)
);
