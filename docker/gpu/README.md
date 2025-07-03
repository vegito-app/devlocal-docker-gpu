# NVIDIA GPU Docker Setup for Debian Bookworm

This directory contains helper scripts to set up NVIDIA GPU support for Docker containers on **Debian Bookworm** systems. These scripts automate the installation of the NVIDIA drivers and the NVIDIA Container Toolkit so you can leverage GPU acceleration inside Docker containers.

## Scripts

⚠️ System Requirements
* Debian Bookworm
* NVIDIA GPU hardware
* Docker installed on the system

⚠️ Limitations
* These scripts are designed for bare metal hosts, not for virtualized environments without direct GPU passthrough.

### 1️⃣ `nvidia-gpu-debian-driver-install.sh`

✅ **Purpose:**  
Installs the NVIDIA GPU driver and required firmware for Debian Bookworm.

✅ **What it does:**
- Checks if an NVIDIA GPU is present on the system (`lspci`).
- Verifies that the system is running Debian Bookworm.
- Configures `apt` to include `contrib`, `non-free`, and `non-free-firmware` repositories.
- Installs the `nvidia-driver` and necessary firmware packages.
- Advises a reboot to activate the NVIDIA driver.

✅ **Usage:**
```bash
sudo ./nvidia-gpu-debian-driver-install.sh
```

➡️ A reboot is required after this script to load the NVIDIA kernel modules.

⸻

2️⃣ nvidia-gpu-container-setup.sh

✅ Purpose:
Installs and configures the NVIDIA Container Toolkit for Docker.

✅ What it does:
	•	Validates NVIDIA driver installation by running nvidia-smi.
	•	Confirms NVIDIA devices are available at /dev/nvidia*.
	•	Adds the NVIDIA Docker APT repository using a signed keyring.
	•	Installs the nvidia-container-toolkit.
	•	Configures the Docker daemon to use NVIDIA runtime.
	•	Restarts Docker (if it fails, asks you to restart manually).
	•	Runs a test container (nvidia/cuda:12.2.0-base-ubuntu20.04) to validate GPU access from Docker.

✅ Usage:

sudo ./nvidia-gpu-container-setup.sh

➡️ This script does not require a reboot but does restart the Docker service.

## Example workflow

### Install NVIDIA drivers

```bash
sudo ./nvidia-gpu-debian-driver-install.sh
```
### Reboot the machine

```bash
sudo reboot
```

### Set up NVIDIA Container Toolkit

```bash
sudo ./nvidia-gpu-container-setup.sh
```

⸻

License MIT

These scripts are provided as-is without any warranty. Use at your own risk.
