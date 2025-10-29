sudo docker network ls | grep app-network
NET=dbp_app-network
ID=$(sudo docker network inspect "$NET" -f '{{.Id}}')
IFACE=br-${ID:0:12}
sudo wireshark -k -i "$IFACE"