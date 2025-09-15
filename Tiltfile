############################################
# Tilt for product-query-svc + Postgres
############################################

# NOTE: old versions of this guide used `load('ext://restart_process', 'restart_process')`.
# Many Tilt installs don't provide that ext; omit the load to avoid startup errors.

# Settings
namespace = 'marketplace-dev'
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
        sync('apps/product-query-svc', '/app'),  # if you run from alpine w/ shell; for distroless, restart only
        sync('backend', '/app'),                 # fallback sync path if you change wiring/main
        # If using a wrapper script to restart, call it here. For distroless, a full image rebuild is fine.
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
