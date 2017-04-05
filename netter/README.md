# Design idea

List of (docker network, physical interface) pairs -> ip link set (physical interface) master (bridge name for(docker network))

# Dockergen part

- We need all of the hashy names for the bridges for the networks we are going to attach routers to (i.e. br-1234567890deadbeef)
- Dockergen can enumerate the networks pretty easily

List of Docker networks -> bridge names 
