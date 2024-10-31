package services

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type NotificationService struct {
	CurrentMessage      string
	lastNotificationSet time.Time
	isRunning           bool
}

func NewNotificationService(w fyne.Window) *NotificationService {
	ns := &NotificationService{isRunning: false}
	ns.Init(w)
	return ns
}

func (ns *NotificationService) Init(w fyne.Window) {
	ns.lastNotificationSet = time.Now()
	ns.isRunning = true
	notificationText := widget.NewLabel("")
	notificationText.TextStyle.Bold = true
	background := canvas.NewRectangle(color.RGBA{0, 0, 0, 100})
	background.SetMinSize(notificationText.MinSize())
	background.Hide()
	content := container.NewStack(background, notificationText)
	boxContainer := container.NewWithoutLayout(content)
	content.Move(fyne.NewPos(5, 5))
	content.Resize(notificationText.MinSize())
	w.Canvas().Overlays().Add(boxContainer)
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		lastMessage := ""
		for range ticker.C {
			if ns.CurrentMessage != lastMessage {
				background.SetMinSize(notificationText.MinSize())
				background.Show()
				notificationText.SetText(ns.CurrentMessage)
				content.Resize(notificationText.MinSize())
				lastMessage = ns.CurrentMessage
			}
			if len(ns.CurrentMessage) > 0 && time.Since(ns.lastNotificationSet) >= 2*time.Second {
				notificationText.SetText("")
				ns.CurrentMessage = ""
				lastMessage = ""
				background.Hide()
			}
		}
	}()
}

func (ns *NotificationService) SetNotification(notification string) {
	ns.CurrentMessage = notification
	ns.lastNotificationSet = time.Now()
}
