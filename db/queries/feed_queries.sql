-- name: InsertFeed :exec
INSERT INTO feed
    (id, system_pk, auto_update_enabled, auto_update_period, config)
VALUES
    (sqlc.arg(id), sqlc.arg(system_pk), sqlc.arg(auto_update_enabled), 
     sqlc.arg(auto_update_period), sqlc.arg(config));

-- name: UpdateFeed :exec
UPDATE feed
SET auto_update_enabled = sqlc.arg(auto_update_enabled), 
    auto_update_period = sqlc.arg(auto_update_period), 
    config = sqlc.arg(config)
WHERE pk = sqlc.arg(feed_pk);

-- name: DeleteFeed :exec
DELETE FROM feed WHERE pk = sqlc.arg(pk);


-- name: ListAutoUpdateFeedsForSystem :many
SELECT feed.id, feed.auto_update_period
FROM feed
    INNER JOIN system ON system.pk = feed.system_pk
WHERE feed.auto_update_enabled
    AND system.id = sqlc.arg(system_id);
