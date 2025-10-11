############################################
# Tilt for product-query-svc + Postgres
############################################

# NOTE: old versions of this guide used `load('ext://restart_process', 'restart_process')`.
# Many Tilt installs don't provide that ext; omit the load to avoid startup errors.

# Settings
namespace = 'gopractice-dev'
svc_name = 'product-query-svc'
# Starlark (Tiltfile) 不支持 Python f-strings，使用字符串连接
docker_ref = svc_name + ':dev'

# Make Tilt create namespace before applying other manifests
k8s_yaml('k8s/namespace.yaml')
# Note: namespace is not a workload resource in Tilt; do not call k8s_resource on it.

# App + DB manifests
k8s_yaml(['k8s/config-app.yaml', 'k8s/postgres.yaml', 'k8s/product-query-svc.yaml'])

# Image build
docker_build(
    ref=docker_ref,
    context='.',
    dockerfile='Dockerfile',
    live_update=[
        # Distroless image: live_update syncs won't hot-reload; expect full rebuilds.
    ]
)

# Wire the k8s Deployment to the local image
allow_k8s_contexts('kind-kind')  # default kind context name

# Resource configuration
k8s_resource(
    workload=svc_name,
    port_forwards=[8080],  # service HTTP
    labels=['backend']
)

k8s_resource(
    workload='postgres',
    port_forwards=[5432],
    labels=['db']
)

# Optional: fast_restart when a small set of files change (if using a sidecar with shell)
# If you have an ext that provides restart_process, you can load it and uncomment the line below
# load('ext://restart_process', 'restart_process')
# restart_process(svc_name, ['touch /tmp/restart'])

# --- Local DB migration step (dev convenience) ---
# Requires: golang-migrate CLI installed locally: https://github.com/golang-migrate/migrate
MIGRATE_PATH = 'apps/product-query-svc/adapters/outbound/postgres/migrations'
# For dockerized migration, use host.docker.internal to reach Tilt port-forward on macOS/Windows.
MIGRATE_URL  = 'postgres://app:app_password@host.docker.internal:5432/productdb?sslmode=disable'

# Run migrations via docker to avoid requiring local migrate CLI.
local_resource(
    name='db-migrate',
    cmd='docker run --rm -v "$(pwd)/' + MIGRATE_PATH + '":/migrations:ro migrate/migrate -path /migrations -database "' + MIGRATE_URL + '" up',
    deps=[MIGRATE_PATH],
    resource_deps=['postgres'],
)

# Ensure app deploys after migrations ran at least once
k8s_resource(
    workload=svc_name,
    resource_deps=['db-migrate']
)
