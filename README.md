[![AWS ECR](https://github.com/korzepadawid/qr-codes-analyzer/actions/workflows/deploy.yml/badge.svg)](https://github.com/korzepadawid/qr-codes-analyzer/actions/workflows/deploy.yml)
[![unit tests](https://github.com/korzepadawid/qr-codes-analyzer/actions/workflows/test.yml/badge.svg)](https://github.com/korzepadawid/qr-codes-analyzer/actions/workflows/test.yml)

# qr-codes-analyzer

![gopher](https://foomandoonian.files.wordpress.com/2012/04/qart1.png)

This is a REST API for analyzing the popularity of your QR codes, developed with Go, PostgreSQL, Redis and AWS.
## Table of content

- [Tech](#tech)
- [How does it work?](#how-does-it-work)
- [Database schema](#database-schema)

## Tech
- [Go](https://golang.org/dl/)
- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/)
- [Gin](https://gin-gonic.com/)
- [sqlc](https://github.com/kyleconroy/sqlc)
- [Testify](https://github.com/stretchr/testify)
- [gomock](https://github.com/golang/mock)
- [migrate](https://github.com/golang-migrate/migrate)
- [S3](https://aws.amazon.com/s3/)
- [ECR](https://aws.amazon.com/ecr/)
- [Docker](https://www.docker.com/)
- [GitHub Actions](https://github.com/features/actions)

## How does it work?
The application works like a proxy for QR codes.
After scanning the QR code, the user is getting redirected to the `/qr-codes/:uuid/redirect` endpoint. The application redirects the end-user to the redirection URL and saves redirection details to the database in parallel.
[![how-does-it-work](https://i.im.ge/2022/09/14/1XuOvy.Diagram-bez-tytulu-drawio.png)](https://im.ge/i/1XuOvy)

- Users can create accounts and sign in with either email or username.

- The application is secured with JWT and RSA.

- Every user can create its own groups of QR codes, to simplify management.

- Users can download CSV files with statistics of the specific QR code. `/qr-codes/:uuid/stats/csv`. File scheme below.

```
uuid,title,url,ipv4,isp,as,city,country,lat,lon,date
0ef9b69d-e1b0-4eb8-9450-684082909c10,a new qr code,https://www.google.pl/,142.250.203.206,Google LLC,AS15169 Google LLC,Warsaw,Poland,52.22970,21.01220,2022-08-23 16:08:33.807904 +0000 UTC
```
- Detailed data about an IPv4 come from the external  API (https://ip-api.com/).

## Database scheme
[![db-scheme](https://i.im.ge/2022/09/14/1X1RQq.Beztytulu.png)](https://im.ge/i/1X1RQq)
