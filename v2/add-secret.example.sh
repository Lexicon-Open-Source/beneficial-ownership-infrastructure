echo "your_db_user" | docker secret create db_user -
echo "your_db_password" | docker secret create db_password -
echo "your_nats_user" | docker secret create nats_user -
echo "your_nats_password" | docker secret create nats_password -
echo "admin" | docker secret create grafana_admin_user -
echo "admin_password" | docker secret create grafana_admin_password -