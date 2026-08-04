/*
 * TgMusicBot - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/AshokShau/TgMusicBot
 */

package handlers

import (
	"ashokshau/tgmusic/config"
	"fmt"
	"runtime"
	"time"

	"ashokshau/tgmusic/src/core"
	"ashokshau/tgmusic/src/core/db"

	td "github.com/AshokShau/gotdbot"
)

// pingHandler handles the /ping command.
func pingHandler(c *td.Client, m *td.Message) error {

	start := time.Now()

	msg, err := m.ReplyText(c, "Pinging… please wait…", nil)
	if err != nil {
		return err
	}

	latency := time.Since(start).Milliseconds()
	uptime := getFormattedDuration(time.Since(startTime))

	response := fmt.Sprintf(
		"<b>📊 System Performance Metrics</b>\n\n"+
			"<b>Bot Latency:</b> <code>%d ms</code>\n"+
			"<b>Uptime:</b> <code>%s</code>\n"+
			"<b>Go Routines:</b> <code>%d</code>\n",
		latency, uptime, runtime.NumGoroutine(),
	)

	_, err = msg.EditText(c, response, &td.EditTextMessageOpts{ParseMode: "HTML"})
	return err
}

// startHandler handles the /start command.
func startHandler(c *td.Client, m *td.Message) error {
	chatID := m.ChatId

	if m.IsPrivate() {
		go func(chatID int64) {
			_ = db.Instance.AddUser(chatID)
		}(chatID)

		response := fmt.Sprintf(
			"<b>Welcome, %s!</b>\n\n"+
				"<b>%s</b> lets you stream high-quality music and video directly in Telegram voice and video chats.\n\n"+
				"<b>Supported platforms:</b> YouTube, Spotify, Apple Music, SoundCloud, Deezer, Twitch, and many more.\n\n"+
				"Use the buttons below to add the bot to your group or explore the available commands.",
			firstName(c, m),
			c.Me.FirstName,
		)

		_, err := m.ReplyPhoto(c, &td.InputFileRemote{Id: config.StartImg}, &td.SendPhotoOpts{
			Caption:     response,
			ParseMode:   "HTML",
			ReplyMarkup: core.AddMeMarkup(c.Me.Usernames.EditableUsername),
		})

		return err
	}

	go func(chatID int64) {
		_ = db.Instance.AddChat(chatID)
	}(chatID)

	uptime := getFormattedDuration(time.Since(startTime))
	response := fmt.Sprintf(
		"<b>%s is ready!</b>\n\n"+
			"<b>Uptime:</b> <code>%s</code>\n\n"+
			"A feature-rich music bot for your group video chats. Play your favorite tracks seamlessly.",
		c.Me.FirstName,
		uptime,
	)

	_, err := m.ReplyText(c, response, &td.SendTextMessageOpts{
		ReplyMarkup: core.SupportBtn(),
		ParseMode:   "HTML",
	})

	return err
}
