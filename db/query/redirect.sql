-- name: CreateRedirect :one
INSERT INTO redirects (qr_code_uuid, ipv4, isp, autonomous_sys, lat, lon, city, country)
VALUES (sqlc.arg(qr_code_uuid), sqlc.arg(ipv4), sqlc.arg(isp), sqlc.arg(autonomous_sys), sqlc.arg(lat), sqlc.arg(lon),
        sqlc.arg(city), sqlc.arg(country))
RETURNING *;

-- name: GetQRCodeRedirectEntries :many
SELECT qc.uuid            AS "uuid",
       qc.redirection_url AS "url",
       qc.title           AS "title",
       r.ipv4             AS "ipv4",
       r.isp              AS "isp",
       r.autonomous_sys   AS "autonomous_sys",
       r.city             AS "city",
       r.lat              AS "lat",
       r.lon              AS "lon",
       r.country          AS "country",
       r.created_at       AS "date"
FROM qr_codes qc
         INNER JOIN redirects r on qc.uuid = r.qr_code_uuid
WHERE qc.uuid = sqlc.arg(uuid)
  AND qc.owner = sqlc.arg(owner)
ORDER BY r.created_at DESC;