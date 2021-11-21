package events

import (
	"fmt"
	"math/rand"
	"time"
	"tobio/reacto/constant"
	disc "tobio/reacto/discordHelpers"

	"github.com/bwmarrin/discordgo"
)

var MSG_TO_WATCH string = ""

func OnReady(s *discordgo.Session, _ *discordgo.Ready) {
	disc.SendLog(s, "init")
}

func OnNewMember(s *discordgo.Session, memberJoinEvent *discordgo.GuildMemberAdd) {

	var newUserName string
	if memberJoinEvent.Member.Nick != "" {
		newUserName = memberJoinEvent.Member.Nick
	} else {
		newUserName = memberJoinEvent.User.Username
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	greeting := constant.RandomGreeting(r)
	suggestion := constant.RandomSuggestion(r)
	secondSuggestion := constant.RandomSuggestion(r)
	closing := constant.RandomClosing(r)

	botWelcomeScript := fmt.Sprintf("%s, %s! %s introduce yourself, tell us your coding story.\n %s check out the react-for-roles channel and let us know where you're based!\n %s", greeting, newUserName, suggestion, secondSuggestion, closing)

	welcomeChannel, err := disc.GetChannelByName(s, "off-topic")
	if err != nil {
		fmt.Println("Error finding off-topic channel")
		fmt.Println(err)
	} else {
		s.ChannelMessageSend(welcomeChannel.ID, botWelcomeScript)
	}
}

func OnReaction(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == s.State.User.ID {
		return
	}

	emojiUsed := m.Emoji.MessageFormat()

	if emojiUsed == "🔥" {
		if MSG_TO_WATCH == "" {
			fmt.Println("MSG TO WATCH was set to ", m.MessageID)
			MSG_TO_WATCH = m.MessageID
		}
	}

	fmt.Println(m.Emoji.APIName())

	if m.MessageID == MSG_TO_WATCH {
		fmt.Println("Pretend role given")
	}
}
