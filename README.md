# EAGLE


### Introduction

`Eagle` is a tool to get system/server parameters instead of running multiple
commands, it will be ran in each system/server and create a `data.md` as an
output.

Parameters are:

    * Hostname
    * IP Address
    * Count of CPU Cores
    * Percentage of CPU Usage
    * Total Memory
    * Percentage of Memory Usage
    * Disk Size
    * Percentage of Disk Usage
    * Load Average for last 1 , 5 and 15 minutes

If there are lots of servers and you have access to them, also you can use
`eagle` with help of `ansible` to gather an overview of them all.

### To Build

First, make sure that `Go` is already installed on your system.If not check
[here](https://go.dev/doc/install).

To make a binary file of this project run the command below:

```bash
go build --ldflags '-linkmode external -extldflags "-static"' .
```

### To Use Locally

To run `eagle` locally , you only need to run the command below:

```bash
./eagle
```

and to view result, run :

```bash
cat data.md
```

### To Use With Ansible

1. Make sure that `ansible` is already installed on your host. If not , Run the
   commands below:

```bash
sudo apt update
sudo apt install software-properties-common
sudo add-apt-repository --yes --update ppa:ansible/ansible
sudo apt install ansible
```

2. Make sure that you have `ssh access` to the servers **without password**
   that you want to run `eagle`.

```bash
ssh-copy-id <USER>@<SERVER_IP>
```

3. You need to generate a host file for ansible to read hosts from it.So run:

```bash
chmod +x inventory/host-generator.sh
```

then

```bash
./inventory/host-generator.sh
```

follow the structure and generate the host file in inventory directory.

4. Last step to use eagle is to run:

```bash
ansible-playbook --extra-vars "user=<$USER_NAME>" -i inventory/hosts playbooks/eagle.yml
```

**NOTE: `$USER_NAME` must be in hosts and you should have access it via ssh without
password**

5. And Finaly you can check the markdown view result in `result.md`

done.

### Example Result view

| Hostname | IP Address    | Up Time    | CPU Usage Percentage  | Disk Usage Percentage | Memory Usage Percentage | Load Average 1 | Load Average 5 | Load average 15 |
| -------- | ------------- | ---------- | --------------------- | --------------------- | ----------------------- | -------------- | -------------- | --------------- |
| office   | 192.168.1.116 | 119h45m38s | 5.34 % out of 8 cores | 66.62 % out of 109GGB | 30 % out of 15.5GGB     | 0.41           | 0.65           | 0.78            |
