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
	"strings"

	"ashokshau/tgmusic/src/core"

	td "github.com/AshokShau/gotdbot"
)

func getHelpCategories() map[string]struct {
	Title   string
	Content string
	Markup  *td.ReplyMarkupInlineKeyboard
} {
	return map[string]struct {
		Title   string
		Content string
		Markup  *td.ReplyMarkupInlineKeyboard
	}{
		"help_user": {
			Title: "User Commands",
			Content: "Commands available to all members of the chat.\n\n" +
				"<b>Playback</b>\n" +
				"/play [song]: Play music from YouTube, Spotify, SoundCloud, and other supported platforms.\n" +
				"/vplay [song]: Play a video in the group video chat.\n" +
				"/fplay [song]: Play a track immediately, skipping the current queue.\n" +
				"/fvplay [song]: Play a video immediately, skipping the current queue.\n\n" +
				"<b>General</b>\n" +
				"/start: Start the bot or verify that it is online.\n" +
				"/help: Open the interactive help menu.\n" +
				"/ping: Display the bot's response time and system information.\n" +
				"/privacy: View the bot's privacy policy.\n" +
				"/queue: Display the current playback queue.",
			Markup: core.BackHelpMenuKeyboard(),
		},
		"help_admin": {
			Title: "Admin Commands",
			Content: "Commands available to chat administrators and authorized users.\n\n" +
				"<b>Playback Control</b>\n" +
				"/skip: Skip the currently playing track.\n" +
				"/pause: Pause playback.\n" +
				"/resume: Resume playback.\n" +
				"/seek [seconds]: Jump to a specific position in the current track.\n" +
				"/mute: Mute the voice chat audio.\n" +
				"/unmute: Unmute the voice chat audio.\n\n" +
				"<b>Queue & Access</b>\n" +
				"/remove [index]: Remove a track from the queue by its position.\n" +
				"/loop [0-10]: Repeat the current track the specified number of times.\n" +
				"/auth: Authorize a user to use administrator commands.\n" +
				"/unauth: Remove a user's authorization.\n" +
				"/authlist: Show all authorized users in the current chat.",
			Markup: core.BackHelpMenuKeyboard(),
		},
		"help_devs": {
			Title: "Developer Commands",
			Content: "Commands intended for bot developers and maintainers.\n\n" +
				"<b>System</b>\n" +
				"/stats: Display bot, hosting, and database statistics.\n" +
				"/av: List all active voice and video chats.\n" +
				"/clearass: Disconnect and clear all active assistant clients.\n" +
				"/leaveall: Disconnect assistants from every active chat.\n" +
				"/logger: View the current logging configuration.",
			Markup: core.BackHelpMenuKeyboard(),
		},
		"help_owner": {
			Title: "Chat Owner Commands",
			Content: "Configuration options available to the chat owner.\n\n" +
				"<b>Chat Settings</b>\n" +
				"/settings: Manage chat settings, including play mode, administrator mode, command auto-delete, and language preferences.",
			Markup: core.BackHelpMenuKeyboard(),
		},
		"help_playlist": {
			Title: "Playlist Commands",
			Content: "Create, organize, and manage your personal playlists.\n\n" +
				"<b>Playlist Management</b>\n" +
				"/createplaylist: Create a new playlist.\n" +
				"/deleteplaylist: Delete one of your playlists.\n" +
				"/addtoplaylist: Add a track to a playlist.\n" +
				"/removefromplaylist: Remove a track from a playlist.\n" +
				"/playlistinfo: Display information about a playlist.\n" +
				"/myplaylists: List all of your playlists.",
			Markup: core.BackHelpMenuKeyboard(),
		},
		"help_autoplay": {
			Title: "Autoplay Commands",
			Content: "Automatically continue playback with recommended tracks.\n\n" +
				"<b>Autoplay</b>\n" +
				"/autoplay: Enable or disable autoplay. When enabled, recommended tracks are automatically queued when playback ends.",
			Markup: core.BackHelpMenuKeyboard(),
		},
	}
}

func helpCallbackHandler(c *td.Client, cb *td.UpdateNewCallbackQuery) error {
	data := cb.DataString()

	user, err := c.GetUser(cb.SenderUserId)
	if err != nil {
		user = &td.User{FirstName: "User", Id: cb.SenderUserId}
	}

	helpCategories := getHelpCategories()

	if strings.Contains(data, "help_all") {
		_ = cb.Answer(c, 0, false, "Opening help menu...", "")
		response := fmt.Sprintf(
			"<b>Welcome, %s!</b>\n\n"+
				"<b>%s</b> is a fast, reliable, and feature-rich music bot for Telegram voice and video chats.\n\n"+
				"<b>Supported platforms:</b> YouTube, Spotify, Apple Music, SoundCloud, Deezer, Twitch, and many more.\n\n"+
				"Select a category below to browse the available commands.",
			user.FirstName,
			c.Me.FirstName,
		)

		_, _ = cb.EditMessageText(c, response, &td.EditTextMessageOpts{ReplyMarkup: core.HelpMenuKeyboard(), ParseMode: "HTML"})
		return nil
	}

	if strings.Contains(data, "help_back") {
		_ = cb.Answer(c, 0, false, "Returning to main menu...", "")

		response := fmt.Sprintf(
			"<b>Welcome, %s!</b>\n\n"+
				"<b>%s</b> lets you stream high-quality music and video directly in Telegram voice and video chats.\n\n"+
				"<b>Supported platforms:</b> YouTube, Spotify, Apple Music, SoundCloud, Deezer, Twitch, and many more.\n\n"+
				"Use the buttons below to add the bot to your group or explore the available commands.",
			user.FirstName,
			c.Me.FirstName,
		)

		caption, err := c.GetFormattedText(response, nil, "HTML")
		if err != nil {
			return err
		}

		content := &td.InputMessagePhoto{
			Photo: &td.InputPhoto{
				Photo: &td.InputFileRemote{Id: config.StartImg},
			},
			Caption: caption,
		}

		msg, err := cb.GetMessage(c)
		if err != nil {
			return err
		}

		_, _ = msg.EditMedia(c, content, &td.EditMessageMediaOpts{ReplyMarkup: core.AddMeMarkup(c.Me.Usernames.EditableUsername)})
		return nil
	}

	if category, ok := helpCategories[data]; ok {
		_ = cb.Answer(c, 0, false, category.Title, "")
		response := fmt.Sprintf("<b>%s</b>\n\n%s", category.Title, category.Content)
		_, _ = cb.EditMessageText(c, response, &td.EditTextMessageOpts{ReplyMarkup: category.Markup, ParseMode: "HTML"})
		return nil
	}

	_ = cb.Answer(c, 0, true, "Unknown help category.", "")
	return nil
}
