# Connect to Bastion Host via VS Code Remote SSH

## Prerequisites

- VS Code installed
- **Remote - SSH** extension installed (by Microsoft)
- Private key file (`.pem`) available on your machine

---

## Step 1: Install Remote - SSH Extension

1. Open VS Code
2. Go to Extensions (`Ctrl+Shift+X`)
3. Search for `Remote - SSH`
4. Install the one by **Microsoft**

---

## Step 2: Open SSH Config File

1. Press `Ctrl+Shift+P`
2. Type `Remote-SSH: Open SSH Configuration File`
3. Select the config file (usually `C:\Users\<you>\.ssh\config`)

---

## Step 3: Add Bastion Host Entry

Paste the following into the config file:

```
Host ignek-ubuntu
    HostName 52.31.233.63
    User ubuntu
    IdentityFile "C:\Users\Dell\OneDrive - TechAnek Technologies (techanek.com)\Meet\techanek\keys\ignek\<your-key-file>.pem"
```

> Replace `<your-key-file>.pem` with your actual private key filename.

> **Important:** Always wrap the `IdentityFile` path in double quotes if it contains spaces.

Save the file with `Ctrl+S`.

---

## Step 4: Connect

1. Open the **Remote Explorer** panel from the left sidebar
2. Click the refresh icon to reload hosts
3. You will see `ignek-ubuntu` listed under **SSH**
4. Click the **arrow icon** next to it to connect
5. VS Code will open a new window connected to the bastion host

---

## Troubleshooting

### Error: `keyword identityfile extra arguments at end of line`

Your `IdentityFile` path has spaces and is not quoted. Wrap it in double quotes:

```
IdentityFile "C:\path with spaces\your-key.pem"
```

---

### Error: `Permission denied (publickey)`

The key file permissions may be too open. If using Git Bash or WSL, run:

```bash
chmod 400 /path/to/your-key.pem
```

On Windows, right-click the `.pem` file > Properties > Security > remove all users except yourself.

---

### Error: `Bad configuration options` or `terminating`

- Check the config file for typos
- Make sure there are no extra spaces or characters after the key path
- Make sure the `IdentityFile` points to the actual key **file**, not a folder

---

## Final SSH Config Reference

```
Host ignek-ubuntu
    HostName 52.31.233.63
    User ubuntu
    IdentityFile "C:\Users\Dell\OneDrive - TechAnek Technologies (techanek.com)\Meet\techanek\keys\ignek\your-key.pem"
```
