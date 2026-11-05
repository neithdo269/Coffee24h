# Build Directory

The build directory is used to house all the build files and assets for your application. 

The structure is:

* bin - Output directory
* darwin - macOS specific files
* linux - Linux specific files and package builders
* windows - Windows specific files

## Mac

The `darwin` directory holds files specific to Mac builds:
- `Info.plist` - Main macOS application plist.
- `Info.dev.plist` - Plist for dev mode.
- `entitlements.plist` - Security and hardened runtime entitlements.
- `build-dmg.sh` - Automated Universal DMG installer packager.

## Linux

The `linux` directory holds packaging files for Linux distributions:
- `coffee24h.desktop` - FreeDesktop standard application entry for app launchers (GNOME, KDE).
- `coffee24h.png` - Standard 256x256 application icon.
- `build-deb.sh` - Automated Debian/Ubuntu (`.deb`) package builder with dependency definitions.
- `build-appimage.sh` - Standalone portable Linux AppImage packager.

## Windows

The `windows` directory contains the manifest and rc files used when building with `wails build`.
These may be customised for your application. To return these files to the default state, simply delete them and
build with `wails build`.

- `icon.ico` - The icon used for the application. This is used when building using `wails build`. If you wish to
  use a different icon, simply replace this file with your own. If it is missing, a new `icon.ico` file
  will be created using the `appicon.png` file in the build directory.
- `installer/*` - The files used to create the Windows installer. These are used when building using `wails build`.
- `info.json` - Application details used for Windows builds. The data here will be used by the Windows installer,
  as well as the application itself (right click the exe -> properties -> details)
- `wails.exe.manifest` - The main application manifest file.