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
			Content: `Commands available to all members of the chat.

<b>Playback</b>
<b>/play</b> [song] — Play music from YouTube, Spotify, SoundCloud, and other supported platforms.
<b>/vplay</b> [song] — Play a video in the group video chat.
<b>/fplay</b> [song] — Play a track immediately, skipping the current queue.
<b>/fvplay</b> [song] — Play a video immediately, skipping the current queue.

<b>General</b>
<b>/start</b> — Start the bot or verify that it is online.
<b>/help</b> — Open the interactive help menu.
<b>/ping</b> — Display the bot's response time and system information.
<b>/privacy</b> — View the bot's privacy policy.
<b>/queue</b> — Display the current playback queue.`,
			Markup: core.BackHelpMenuKeyboard(),
		},
		"help_admin": {
			Title: "Admin Commands",
			Content: `Commands available to chat administrators and authorized users.

<b>Playback Control</b>
<b>/skip</b> — Skip the currently playing track.
<b>/pause</b> — Pause playback.
<b>/resume</b> — Resume playback.
<b>/seek</b> [seconds] — Jump to a specific position in the current track.
<b>/mute</b> — Mute the voice chat audio.
<b>/unmute</b> — Unmute the voice chat audio.

<b>Queue & Access</b>
<b>/remove</b> [index] — Remove a track from the queue by its position.
<b>/loop</b> [0-10] — Repeat the current track the specified number of times.
<b>/auth</b> — Authorize a user to use administrator commands.
<b>/unauth</b> — Remove a user's authorization.
<b>/authlist</b> — Show all authorized users in the current chat.`,
			Markup: core.BackHelpMenuKeyboard(),
		},
		"help_devs": {
			Title: "Developer Commands",
			Content: `Commands intended for bot developers and maintainers.

<b>System</b>
<b>/stats</b> — Display bot, hosting, and database statistics.
<b>/av</b> — List all active voice and video chats.
<b>/clearass</b> — Disconnect and clear all active assistant clients.
<b>/leaveall</b> — Disconnect assistants from every active chat.
<b>/logger</b> — View the current logging configuration.`,
			Markup: core.BackHelpMenuKeyboard(),
		},
		"help_owner": {
			Title: "Chat Owner Commands",
			Content: `Configuration options available to the chat owner.

<b>Chat Settings</b>
<b>/settings</b> — Manage chat settings, including play mode, administrator mode, command auto-delete, and language preferences.`,
			Markup: core.BackHelpMenuKeyboard(),
		},
		"help_playlist": {
			Title: "Playlist Commands",
			Content: `Create, organize, and manage your personal playlists.

<b>Playlist Management</b>
<b>/createplaylist</b> — Create a new playlist.
<b>/deleteplaylist</b> — Delete one of your playlists.
<b>/addtoplaylist</b> — Add a track to a playlist.
<b>/removefromplaylist</b> — Remove a track from a playlist.
<b>/playlistinfo</b> — Display information about a playlist.
<b>/myplaylists</b> — List all of your playlists.`,
			Markup: core.BackHelpMenuKeyboard(),
		},
		"help_autoplay": {
			Title: "Autoplay Commands",
			Content: `Automatically continue playback with recommended tracks.

<b>Autoplay</b>
<b>/autoplay</b> — Enable or disable autoplay. When enabled, recommended tracks are automatically queued when playback ends.`,
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
			"<h3>Welcome, %s!</h3>\n"+
				"<p><b>%s</b> is a fast, reliable, and feature-rich music bot for Telegram voice and video chats.</p>\n\n"+
				"<p><b>Supported platforms:</b> YouTube, Spotify, Apple Music, SoundCloud, Deezer, Twitch, and many more.</p>\n\n"+
				"<p>Select a category below to browse the available commands.</p>",
			user.FirstName,
			c.Me.FirstName,
		)

		richMessage := &td.InputRichMessage{
			Source: &td.RichMessageSourceHtml{
				Text: response,
			},
		}
		_, _ = c.EditMessageText(cb.ChatId, &td.InputMessageRichMessage{Message: richMessage}, cb.MessageId, &td.EditMessageTextOpts{ReplyMarkup: core.HelpMenuKeyboard()})
		return nil
	}

	if strings.Contains(data, "help_back") {
		_ = cb.Answer(c, 0, false, "Returning to main menu...", "")

		response := fmt.Sprintf(
			"<img src=\"%s\"/>\n"+
				"<h3>Welcome, %s!</h3>\n"+
				"<p><b>%s</b> lets you stream high-quality music and video directly in Telegram voice and video chats.</p>\n\n"+
				"<p><b>Supported platforms:</b> YouTube, Spotify, Apple Music, SoundCloud, Deezer, Twitch, and many more.</p>\n\n"+
				"<p>Use the buttons below to add the bot to your group or explore the available commands.</p>",
			config.StartImg,
			user.FirstName,
			c.Me.FirstName,
		)

		richMessage := &td.InputRichMessage{
			Source: &td.RichMessageSourceHtml{
				Text: response,
			},
		}
		_, _ = c.EditMessageText(cb.ChatId, &td.InputMessageRichMessage{Message: richMessage}, cb.MessageId, &td.EditMessageTextOpts{ReplyMarkup: core.AddMeMarkup(c.Me.Usernames.EditableUsername)})
		return nil
	}

	if category, ok := helpCategories[data]; ok {
		_ = cb.Answer(c, 0, false, category.Title, "")
		response := fmt.Sprintf("<h3>%s</h3>\n\n%s\n\n<i>Use the buttons below to go back.</i>", category.Title, category.Content)
		richMessage := &td.InputRichMessage{
			Source: &td.RichMessageSourceHtml{
				Text: response,
			},
		}
		_, _ = c.EditMessageText(cb.ChatId, &td.InputMessageRichMessage{Message: richMessage}, cb.MessageId, &td.EditMessageTextOpts{ReplyMarkup: category.Markup})
		return nil
	}

	_ = cb.Answer(c, 0, true, "Unknown help category.", "")
	return nil
}
