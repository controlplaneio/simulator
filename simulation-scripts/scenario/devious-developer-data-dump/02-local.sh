set -x

NODE_IP="$(get_node 1)"

run_scp "./_scripts/logo.png" "$NODE_IP:/mnt/psql-data/gitea/public/img/"
run_scp "./_scripts/logo.svg" "$NODE_IP:/mnt/psql-data/gitea/public/img/"
run_scp "./_scripts/favicon.svg" "$NODE_IP:/mnt/psql-data/gitea/public/img/"
run_scp "./_scripts/home.tmpl" "$NODE_IP:/mnt/psql-data/gitea/templates/"