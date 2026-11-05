//go:build !windows

package main

var currentCupEdition = "gdu"

func initTray(a *App) {}
func removeTray() {}
func updateTrayStatus(isAwake bool) {}
func setWindowIcons(title string) {}
func switchCupIcon(edition string) {}
func showTrayNotification(title, message string) {}

