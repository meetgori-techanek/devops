# Increase AWS EBS Volume and Extend LVM Filesystem (Ubuntu)

This document describes the steps to increase an AWS EBS volume and expand the Linux filesystem when the disk uses **LVM**.  
Example case: Increasing storage from **220 GB → 250 GB** for an ELK data volume mounted at `/data`.

---

# 1. Increase EBS Volume in AWS

1. Go to **AWS Console**
2. Navigate to **EC2 → Elastic Block Store → Volumes**
3. Select the volume attached to the instance
4. Click **Modify Volume**
5. Change size from **220 GB → 250 GB**
6. Apply changes

Wait until the volume modification state becomes **completed/optimizing**.

---

# 2. Verify New Disk Size on Server

SSH into the instance and check the disk size.

```bash
lsblk
````

Example output:

```
nvme1n1         250G disk
└─datavg-datalv 220G lvm  /data
```

Check filesystem usage:

```bash
df -h
```

Example:

```
/dev/mapper/datavg-datalv 216G 179G 28G 87% /data
```

---

# 3. Resize the LVM Physical Volume

Extend the physical volume to recognize the new disk size.

```bash
sudo pvresize /dev/nvme1n1
```

Verify:

```bash
sudo pvs
```

Expected:

```
PV           VG     PSize    PFree
/dev/nvme1n1 datavg <250.00g 30.00g
```

---

# 4. Extend the Logical Volume

Allocate all free space in the volume group to the logical volume.

```bash
sudo lvextend -l +100%FREE -r /dev/datavg/datalv
```

Explanation:

* `-l +100%FREE` → Use all remaining space
* `-r` → Automatically resize filesystem

---

# 5. Verify Filesystem Expansion

```bash
df -h
```

Example result:

```
/dev/mapper/datavg-datalv 246G 179G 56G 77% /data
```

---

# 6. Verify LVM Layout (Optional)

```bash
sudo pvs
sudo vgs
sudo lvs
```

Expected:

| Component      | Size           |
| -------------- | -------------- |
| Disk           | 250 GB         |
| Volume Group   | 250 GB         |
| Logical Volume | ~250 GB        |
| Filesystem     | ~246 GB usable |

---

# 7. Notes for ELK Servers

* Elasticsearch typically stores data under `/data`
* Ensure disk usage stays **below 85%** to avoid shard allocation issues
* Increasing disk space does **not require Elasticsearch restart**


---

# Summary

| Step | Action                         |
| ---- | ------------------------------ |
| 1    | Increase EBS volume in AWS     |
| 2    | Verify disk size with `lsblk`  |
| 3    | Resize LVM physical volume     |
| 4    | Extend logical volume          |
| 5    | Verify filesystem with `df -h` |

Final result: `/data` expanded from **216 GB → 246 GB usable storage**.

```
```
