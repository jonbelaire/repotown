module github.com/jonbelaire/repotown/services/parking

go 1.20

require (
	github.com/go-chi/chi/v5 v5.0.10
	github.com/google/uuid v1.3.1
	github.com/jonbelaire/repotown/packages/go-core v0.0.0
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
)

replace github.com/jonbelaire/repotown/packages/go-core => ../../packages/go-core
