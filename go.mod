module github.com/F0hor/gator

go 1.26.2

replace github.com/F0hor/config v0.0.0 => ./internal/config

require github.com/F0hor/config v0.0.0

replace github.com/F0hor/database v0.0.0 => ./internal/database

require github.com/F0hor/database v0.0.0


require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/lib/pq v1.12.3 // indirect
)
