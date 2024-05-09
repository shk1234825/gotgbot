package message

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

func All(_ *gotgbot.Message) bool {
	return true
}

func FromUserID(id int64) filters.Message {
	return func(m *gotgbot.Message) bool {
		return m.From.Id == id
	}
}

func FromUsername(name string) filters.Message {
	return func(m *gotgbot.Message) bool {
		return m.From.Username == name
	}
}

func ChatUsername(name string) filters.Message {
	return func(m *gotgbot.Message) bool {
		return m.Chat.Username != "" && m.Chat.Username == name
	}
}

func ChatID(id int64) filters.Message {
	return func(m *gotgbot.Message) bool {
		return m.Chat.Id == id
	}
}

func ForwardFromUserID(id int64) filters.Message {
	return func(m *gotgbot.Message) bool {
		if m.ForwardOrigin == nil {
			return false
		}
		u := m.ForwardOrigin.MergeMessageOrigin().SenderUser
		return u != nil && u.Id == id
	}
}

func ForwardFromChatID(id int64) filters.Message {
	return func(m *gotgbot.Message) bool {
		if m.ForwardOrigin == nil {
			return false
		}
		c := m.ForwardOrigin.MergeMessageOrigin().Chat
		return c != nil && c.Id == id
	}
}

func Regex(p string) (filters.Message, error) {
	r, err := regexp.Compile(p)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %w", err)
	}
	return func(m *gotgbot.Message) bool {
		return r.MatchString(m.GetText())
	}, nil
}

func Reply(msg *gotgbot.Message) bool {
	return msg.ReplyToMessage != nil
}

func ChatType(t string) filters.Message {
	return func(m *gotgbot.Message) bool {
		return m.Chat.Type == t
	}
}

func Private(msg *gotgbot.Message) bool {
	return msg.Chat.Type == "private"
}

func Group(msg *gotgbot.Message) bool {
	return msg.Chat.Type == "group"
}

func Supergroup(msg *gotgbot.Message) bool {
	return msg.Chat.Type == "supergroup"
}

func Channel(msg *gotgbot.Message) bool {
	return msg.Chat.Type == "channel"
}

func Business(msg *gotgbot.Message) bool {
	return msg.BusinessConnectionId != ""
}

func Forwarded(msg *gotgbot.Message) bool {
	return msg.ForwardOrigin != nil
}

func Text(msg *gotgbot.Message) bool {
	return msg.Text != ""
}

func HasPrefix(prefix string) filters.Message {
	return func(msg *gotgbot.Message) bool {
		return strings.HasPrefix(msg.GetText(), prefix)

	}
}

func HasSuffix(suffix string) filters.Message {
	return func(msg *gotgbot.Message) bool {
		return strings.HasSuffix(msg.GetText(), suffix)
	}
}

func Contains(contains string) filters.Message {
	return func(msg *gotgbot.Message) bool {
		return strings.Contains(msg.GetText(), contains)
	}
}

func Equal(eq string) filters.Message {
	return func(msg *gotgbot.Message) bool {
		return msg.GetText() == eq
	}
}

func Caption(msg *gotgbot.Message) bool {
	return msg.Caption != ""
}

// Command returns true if the message starts with a bot_command entitiy.
func Command(msg *gotgbot.Message) bool {
	ents := msg.GetEntities()
	return len(ents) > 0 && ents[0].Type == "bot_command" && ents[0].Offset == 0
}

// CommandName returns true if the message starts with a bot_command entity matching the name provided.
// The bot object is required to ensure that the /command@username format is respected.
func CommandName(b *gotgbot.Bot, name string) func(msg *gotgbot.Message) bool {
	return CommandNameTriggers(b, name, []rune("/"))
}

// CommandNameTriggers returns true is the message starts with one of the triggers provided, and expects position 0 to
// either be a bot_command entity or no entity at all.
// This means that this filter will not trigger if it has a a code/bold/italics entity at position 0, allowing bot users
// to "escape" commands.
//
// Note: if you are looking to standardise the triggers across your bot, it could be a good idea to create a helper
// function around this to define all your triggers in one place.
func CommandNameTriggers(b *gotgbot.Bot, command string, triggers []rune) func(msg *gotgbot.Message) bool {
	return func(msg *gotgbot.Message) bool {
		text := msg.GetText()

		var msgCmd string
		for _, t := range triggers {
			if r, _ := utf8.DecodeRuneInString(text); r != t {
				continue
			}

			split := strings.Split(strings.ToLower(strings.Fields(text)[0]), "@")
			if len(split) > 1 && split[1] != strings.ToLower(b.User.Username) {
				return false
			}
			msgCmd = split[0][1:]
			break
		}
		if msgCmd == "" {
			return false
		}

		ents := msg.GetEntities()
		if len(ents) != 0 && ents[0].Offset == 0 && ents[0].Type != "bot_command" {
			return false
		}

		return msgCmd == command
	}
}

func Animation(msg *gotgbot.Message) bool {
	return msg.Animation != nil
}

func Audio(msg *gotgbot.Message) bool {
	return msg.Audio != nil
}

func Document(msg *gotgbot.Message) bool {
	return msg.Document != nil
}

func Photo(msg *gotgbot.Message) bool {
	return len(msg.Photo) > 0
}

func Sticker(msg *gotgbot.Message) bool {
	return msg.Sticker != nil
}

func Video(msg *gotgbot.Message) bool {
	return msg.Video != nil
}

func VideoNote(msg *gotgbot.Message) bool {
	return msg.VideoNote != nil
}

func Voice(msg *gotgbot.Message) bool {
	return msg.Voice != nil
}

func Contact(msg *gotgbot.Message) bool {
	return msg.Contact != nil
}

func Dice(msg *gotgbot.Message) bool {
	return msg.Dice != nil
}

func DiceValue(msg *gotgbot.Message, val int64) bool {
	return msg.Dice != nil && msg.Dice.Value == val
}

func Game(msg *gotgbot.Message) bool {
	return msg.Game != nil
}

func Poll(msg *gotgbot.Message) bool {
	return msg.Poll != nil
}

func Venue(msg *gotgbot.Message) bool {
	return msg.Venue != nil
}

func Location(msg *gotgbot.Message) bool {
	return msg.Location != nil
}

func NewChatMembers(msg *gotgbot.Message) bool {
	return msg.NewChatMembers != nil
}

func LeftChatMember(msg *gotgbot.Message) bool {
	return msg.LeftChatMember != nil
}

func PinnedMessage(msg *gotgbot.Message) bool {
	return msg.PinnedMessage != nil
}

func ViaBot(msg *gotgbot.Message) bool {
	return msg.ViaBot != nil
}

func Entities(m *gotgbot.Message) bool {
	return len(m.Entities) > 0
}

func Entity(entType string) filters.Message {
	return func(m *gotgbot.Message) bool {
		for _, ent := range m.Entities {
			if ent.Type == entType {
				return true
			}
		}
		return false
	}
}

func CaptionEntities(m *gotgbot.Message) bool {
	return len(m.CaptionEntities) > 0
}

func CaptionEntity(entType string) filters.Message {
	return func(m *gotgbot.Message) bool {
		for _, ent := range m.CaptionEntities {
			if ent.Type == entType {
				return true
			}
		}
		return false
	}
}

func Migrate(msg *gotgbot.Message) bool {
	return msg.MigrateFromChatId != 0 || msg.MigrateToChatId != 0
}

func MigrateFrom(msg *gotgbot.Message) bool {
	return msg.MigrateFromChatId != 0
}

func MigrateTo(msg *gotgbot.Message) bool {
	return msg.MigrateToChatId != 0
}

func ReplyMarkup(msg *gotgbot.Message) bool {
	return msg.ReplyMarkup != nil
}

func MediaGroup(msg *gotgbot.Message) bool {
	return msg.MediaGroupId != ""
}

func IsAutomaticForward(msg *gotgbot.Message) bool {
	return msg.IsAutomaticForward
}

func UsersShared(msg *gotgbot.Message) bool {
	return msg.UsersShared != nil
}

func ChatShared(msg *gotgbot.Message) bool {
	return msg.ChatShared != nil
}

func Story(msg *gotgbot.Message) bool {
	return msg.Story != nil
}

func TopicEdited(msg *gotgbot.Message) bool {
	return msg.ForumTopicEdited != nil
}

func TopicCreated(msg *gotgbot.Message) bool {
	return msg.ForumTopicCreated != nil
}

func TopicClosed(msg *gotgbot.Message) bool {
	return msg.ForumTopicClosed != nil
}

func TopicReopened(msg *gotgbot.Message) bool {
	return msg.ForumTopicReopened != nil
}

func TopicAction(msg *gotgbot.Message) bool {
	return TopicEdited(msg) || TopicCreated(msg) ||
		TopicClosed(msg) || TopicReopened(msg)
}
