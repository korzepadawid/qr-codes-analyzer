-- name: CreateRedirect :one
INSERT INTO redirects (qr_code_uuid, ipv4, isp, autonomous_sys, lat, lon, city, country)
VALUES (sqlc.arg(qr_code_uuid), sqlc.arg(ipv4), sqlc.arg(isp), sqlc.arg(autonomous_sys), sqlc.arg(lat), sqlc.arg(lon),
        sqlc.arg(city), sqlc.arg(country))
RETURNING *;
