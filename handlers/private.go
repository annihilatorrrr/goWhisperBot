/**
 * ezWhisperBot - A Telegram bot for sending whisper messages
 * Copyright (C) 2021  Roj Serbest
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package handlers

import (
	"fmt"
	"os"
	"time"

	"goWhisperBot/whispers"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

const (
	defaultText = "This bot allows you to send whisper messages, " +
		"works only in inline mode\n\n" +
		"<a href=\"https://github.com/rojserbest/ezWhisperBot\">Source Code</a>" +
		" | <a href=\"https://t.me/rojserbest\">Developer</a>" +
		" | <a href=\"https://t.me/ezupdev\">Support Chat</a>"
	learnText = "This bot works only in inline mode, a example use would be like " +
		"this:\n\n" +
		"- Write a whisper to @username\n" +
		"<code>@ezWhisperBot @username some text here</code>\n\n" +
		"- Whisper to the first one who open it (can also be used in PM)\n" +
		"<code>@ezWhisperBot some text here</code>\n\n" +
		"If the username is <code>@all</code>, anyone who open the whisper first can read it."
)

var (
	defaultReplyMarkup = gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				gotgbot.InlineKeyboardButton{
					Text:              "Select a Chat to Try",
					SwitchInlineQuery: " ",
				},
				gotgbot.InlineKeyboardButton{
					Text:                         "Try in This Chat",
					SwitchInlineQueryCurrentChat: " ",
				},
			}, {
				gotgbot.InlineKeyboardButton{
					Text:         "My Whispers",
					CallbackData: "list_whispers",
				},
			},
		},
	}
	learnReplyMarkup = gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				gotgbot.InlineKeyboardButton{
					Text:         "Next",
					CallbackData: "learn_next",
				},
			},
		},
	}
)

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	text, replyMarkup := defaultText, defaultReplyMarkup

	if strings.HasSuffix(ctx.EffectiveMessage.Text, "learn") {
		text, replyMarkup = learnText, learnReplyMarkup
	}

	ctx.EffectiveMessage.Reply(
		b,
		text,
		&gotgbot.SendMessageOpts{
			ParseMode:             "HTML",
			DisableWebPagePreview: true,
			ReplyMarkup:           replyMarkup,
		},
	)

	return nil
}

func back(b *gotgbot.Bot, ctx *ext.Context) error {
	ctx.EffectiveMessage.EditText(
		b,
		defaultText,
		&gotgbot.EditMessageTextOpts{
			DisableWebPagePreview: true,
			ParseMode:             "HTML",
			ReplyMarkup:           defaultReplyMarkup,
		},
	)

	if ctx.CallbackQuery.Data == "learn_next" {
		ctx.CallbackQuery.Answer(
			b,
			&gotgbot.AnswerCallbackQueryOpts{
				Text: "🤖 Now you may try it in inline mode",
			},
		)
	}

	return nil
}

func saveWhispers(b *gotgbot.Bot, ctx *ext.Context) error {
	if strings.EqualFold(ctx.EffectiveMessage.From.Username, os.Getenv("SUDO")) {
		return nil
	}

	whispers.SaveWhispers()
	return nil
}

func myWhispers(b *gotgbot.Bot, ctx *ext.Context) error {
	_whispers, text := []whispers.Whisper{}, ""

	for _, whisper := range whispers.Whispers.Whispers {
		if whisper.Sender == ctx.EffectiveUser.Id {
			_whispers = append(_whispers, whisper)
		}
	}

	if len(_whispers) == 0 {
		text = "You don't have any whispers"
	} else {
		text = fmt.Sprintf("You have %d whispers", len(_whispers))
	}

	ctx.EffectiveMessage.EditText(
		b,
		text,
		&gotgbot.EditMessageTextOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						gotgbot.InlineKeyboardButton{
							Text:         "🗑 Delete My Whispers",
							CallbackData: "delete_my_whispers",
						},
					},
					{
						gotgbot.InlineKeyboardButton{
							Text:         "◀️ Back to Main Page",
							CallbackData: "start",
						},
					},
				},
			},
		},
	)
	return nil
}

func deleteMyWhispers(b *gotgbot.Bot, ctx *ext.Context) error {
	_whispers := []whispers.Whisper{}
	ids := []string{}

	for id, whisper := range whispers.Whispers.Whispers {
		if whisper.Sender == ctx.EffectiveUser.Id {
			_whispers = append(_whispers, whisper)
			ids = append(ids, id)
		}
	}

	if len(_whispers) == 0 {
		ctx.CallbackQuery.Answer(
			b,
			&gotgbot.AnswerCallbackQueryOpts{
				Text: "You don't have any whispers",
			},
		)
	} else {
		for _, id := range ids {
			delete(whispers.Whispers.Whispers, id)
		}

		ctx.CallbackQuery.Answer(
			b,
			&gotgbot.AnswerCallbackQueryOpts{
				Text: fmt.Sprintf("Removed %d whispers", len(_whispers)),
			},
		)
		now := time.Now().UTC().String()
		ctx.EffectiveMessage.EditText(
			b,
			"Your whispers has been removed at "+now,
			nil,
		)
	}

	return nil
}
