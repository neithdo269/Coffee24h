//go:build !windows

package main

func initTray(a *App) {}
func removeTray() {}
func updateTrayStatus(isAwake bool) {}
func setWindowIcons(title string) {}

