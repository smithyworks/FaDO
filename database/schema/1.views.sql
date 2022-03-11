CREATE VIEW buckets_faas_deployments AS
  SELECT x.bucket_id, x.bucket_name, array_agg(x.faas_id) AS faas_ids, array_agg(x.faas_url) AS faas_urls
  FROM (
    SELECT b.bucket_id, b.name AS bucket_name, fd.faas_id, fd.url AS faas_url
      FROM buckets b
      INNER JOIN replica_bucket_locations rbl ON rbl.bucket_id = b.bucket_id
      LEFT JOIN storage_deployments sd ON sd.storage_id = rbl.storage_id
      LEFT JOIN clusters c ON c.cluster_id = sd.cluster_id
      LEFT JOIN faas_deployments fd ON fd.cluster_id = c.cluster_id
    UNION
    SELECT b.bucket_id, b.name AS bucket_name, fd.faas_id, fd.url AS faas_url
      FROM buckets b
      LEFT JOIN storage_deployments sd ON sd.storage_id = b.storage_id
      LEFT JOIN clusters c ON c.cluster_id = sd.cluster_id
      LEFT JOIN faas_deployments fd ON fd.cluster_id = c.cluster_id
  ) AS x GROUP BY x.bucket_id, x.bucket_name;

CREATE VIEW existing_bucket_locations AS
  SELECT bucket_id, array_agg(storage_id) AS storage_ids
  FROM replica_bucket_locations
  GROUP BY bucket_id;

CREATE VIEW bucket_replications AS
  SELECT rbl.bucket_id, b.name AS bucket_name,
    b.storage_id AS src_storage_id, src_sd.alias AS src_storage_alias, src_sd.minio_deployment_id AS src_minio_deployment_id,
    rbl.storage_id AS dst_storage_id, dst_sd.alias AS dst_storage_alias, dst_sd.minio_deployment_id AS dst_minio_deployment_id
  FROM replica_bucket_locations rbl
  LEFT JOIN buckets b ON b.bucket_id = rbl.bucket_id
  LEFT JOIN storage_deployments src_sd ON src_sd.storage_id = b.storage_id
  LEFT JOIN storage_deployments dst_sd ON dst_sd.storage_id = rbl.storage_id;
