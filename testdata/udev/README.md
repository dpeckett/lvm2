# LVM2 Udev Rules

For the containerized integration tests we need to run a minimal Udev daemon,

On paper LVM2 can work without Udev but many operations have an implicit requirement on it.
