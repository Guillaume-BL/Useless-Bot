<<<<<<< HEAD
package main

import (
	"encoding/json"

	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	//"time"

	"./dgvoice"

	"github.com/bwmarrin/discordgo"
)

// Struct

type Sounds struct {
	Sounds []Sound `json:"sounds"`
}

type Sound struct {
	Text  string `json:"text"`
	Path  string `json:"path"`
	Regex string `json:"regex"`
}

// Variables used for command line parameters
var (
	Token  string
	dg     *discordgo.Session
	sounds Sounds
)

var isPlaying bool = false

var capybara_resource = []string{
	"https://cdn.britannica.com/s:800x450,c:crop/94/194294-138-B2CF7780/overview-capybara.jpg",
	"https://www.rainforest-alliance.org/sites/default/files/styles/750w_585h/public/2016-09/capybara.jpg",
	"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcR10Sx9MBImh81eaR2tajNQDI6T1olm-ymw0pqoxab01eV5Ny6Q6Q&s",
	"https://66.media.tumblr.com/tumblr_mcdjthzQN21r03kk7o1_500.jpg",
	"https://66.media.tumblr.com/00fc839d7634e8ef0f8f05475c107584/tumblr_mmny371ixu1r03kk7o1_500.jpg",
	"https://sociorocketnewsen.files.wordpress.com/2016/03/capy-4.jpg",
	"https://i.pinimg.com/originals/7e/60/52/7e6052b52add1adbde9fb2888312e898.jpg",
	"https://upload.wikimedia.org/wikipedia/commons/thumb/b/bc/Bristol.zoo.capybara.arp.jpg/1200px-Bristol.zoo.capybara.arp.jpg",
	"https://image.jimcdn.com/app/cms/image/transf/dimension=1920x10000:format=jpg/path/s5dde8bff85c81b2f/image/iafb738a20225a5d5/version/1476089104/capybara-cabiai-famille-bebe-parent.jpg",
	"https://upload.wikimedia.org/wikipedia/commons/thumb/e/e1/Cattle_tyrant_%28Machetornis_rixosa%29_on_Capybara.jpg/800px-Cattle_tyrant_%28Machetornis_rixosa%29_on_Capybara.jpg",
	"https://i.pinimg.com/originals/11/3f/26/113f26c0109eb801a665d6b17925b0be.jpg",
	"https://i.pinimg.com/originals/3f/49/1d/3f491d524ef380ff7d6e040de803d73e.jpg",
}

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	var err error
	// Create a new Discord session using the provided bot token.
	dg, err = discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	loadSounds()

	dg.AddHandler(messageEvent)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func loadSounds() {
	jsonFile, err := os.Open("sounds.json")
	defer jsonFile.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err := json.Unmarshal(byteValue, &sounds); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Sounds Loaded")
		for _, s := range sounds.Sounds {
			fmt.Printf("File : %s Match on : %s\n", s.Path, s.Regex)
		}
	}
}

func voiceEvent(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
	fmt.Println("Echo")
	fmt.Println(vs.UserID)

}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageEvent(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	g, _ := s.Guild(m.GuildID)
	if g != nil {

		content_lowered := strings.ToLower(m.Content)

		if strings.Contains(content_lowered, "!help") {
			response := &discordgo.MessageEmbed{
				Title: "The Useless Bot - SoundBox",
			}
			for _, s := range sounds.Sounds {
				response.Description += s.Text
				response.Description += "\n"
			}
			s.ChannelMessageSendEmbed(m.ChannelID, response)
		}

		if strings.Contains(content_lowered, "krjone") {
			s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
		}

		v := g.VoiceStates
		var voiceStateOfUser *discordgo.VoiceState
		for _, ve := range v {
			if ve.UserID == m.Author.ID {
				voiceStateOfUser = ve
			}
		}

		if voiceStateOfUser != nil {

			for _, sound := range sounds.Sounds {
				if re, err := regexp.Compile(sound.Regex); err != nil {
					fmt.Println(err)
				} else {

					if re.Match([]byte(content_lowered)) {
						response := &discordgo.MessageEmbed{
							Title:       "The Useless Bot - SoundBox",
							Description: sound.Text,
						}
						s.ChannelMessageSendEmbed(m.ChannelID, response)
						playSound(s, m.ChannelID, voiceStateOfUser.GuildID, voiceStateOfUser.ChannelID, sound.Path)
						break
					}
				}
			}
		}
		go sendEmojiIfContains(s, m, "\U0001F409", "dragon")
		go sendEmojiIfContains(s, m, getEmojiFromFromGuild(g, "capybara"), "capybara")
		go sendCapybaraPicture(s, m)
	} else {
		go sendPM(s, m)
		return
	}

}

func getEmojiFromFromGuild(g *discordgo.Guild, emoji_name string) string {
	emojiId := ""
	if g != nil {
		for _, i := range g.Emojis {
			if strings.Contains(strings.ToLower(i.Name), strings.ToLower(emoji_name)) {
				emojiId = i.APIName()
			}
		}
	}
	return emojiId

}

func sendEmojiIfContains(s *discordgo.Session, m *discordgo.MessageCreate, emojiID, match string) {
	is_match := true
	content_lowered := strings.ToLower(m.Content)
	for _, c := range match {
		if strings.IndexRune(content_lowered, c) == -1 {
			is_match = false
		}
	}
	if is_match {
		if err := s.MessageReactionAdd(m.ChannelID, m.Message.ID, emojiID); err != nil {
			fmt.Println(err)
		}
	}
}

func sendCapybaraPicture(s *discordgo.Session, m *discordgo.MessageCreate) {
	content_lowered := strings.ToLower(m.Content)
	if strings.Contains(content_lowered, "capybara") {
		i := rand.Intn(len(capybara_resource) - 1)
		image := &discordgo.MessageEmbedImage{
			URL: capybara_resource[i],
		}
		response := &discordgo.MessageEmbed{
			Image: image,
		}
		s.ChannelMessageSendEmbed(m.ChannelID, response)
	}
}

func sendPM(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelTyping(m.ChannelID)
	time.Sleep(1 * time.Second)
	s.ChannelMessageSend(m.ChannelID, "Pourquoi tu m'envoie un message "+m.Author.Username)
	s.ChannelTyping(m.ChannelID)
	time.Sleep(3 * time.Second)
	s.ChannelMessageSend(m.ChannelID, "ça sert à rien putain, t'es con ou quoi ?")
}

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, MessageChannelID, guildID, vocalChannelID, filename string) (err error) {
	if isPlaying == false {

		isPlaying = true

		// Join the provided voice channel.
		vc, err := s.ChannelVoiceJoin(guildID, vocalChannelID, false, true)
		if err != nil {
			return err
		}

		// Sleep for a specified amount of time before playing the sound
		//time.Sleep(50 * time.Millisecond)

		// Start speaking.
		vc.Speaking(true)
		c := make(chan bool)
		dgvoice.PlayAudioFile(vc, filename, c)
		c <- true

		// Stop speaking
		vc.Speaking(false)

		// Disconnect from the provided voice channel.
		//vc.Disconnect()

		isPlaying = false

		//time.Sleep(250 * time.Second)

		return nil
	} else {
		response := &discordgo.MessageEmbed{
			Description: "Useless Bot already doing useless things !! ",
		}
		s.ChannelMessageSendEmbed(MessageChannelID, response)

		return nil
	}

}
=======
package main

import (
	"encoding/json"

	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	//"time"

	"./dgvoice"

	"github.com/bwmarrin/discordgo"
)

// Struct

type Sounds struct {
	Sounds []Sound `json:"sounds"`
}

type Sound struct {
	Text  string `json:"text"`
	Path  string `json:"path"`
	Regex string `json:"regex"`
}

// Variables used for command line parameters
var (
	Token  string
	dg     *discordgo.Session
	sounds Sounds
)

var isPlaying bool = false

var capybara_resource = []string{
	"https://cdn.britannica.com/s:800x450,c:crop/94/194294-138-B2CF7780/overview-capybara.jpg",
	"https://www.rainforest-alliance.org/sites/default/files/styles/750w_585h/public/2016-09/capybara.jpg",
	"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcR10Sx9MBImh81eaR2tajNQDI6T1olm-ymw0pqoxab01eV5Ny6Q6Q&s",
	"https://66.media.tumblr.com/tumblr_mcdjthzQN21r03kk7o1_500.jpg",
	"https://66.media.tumblr.com/00fc839d7634e8ef0f8f05475c107584/tumblr_mmny371ixu1r03kk7o1_500.jpg",
	"https://sociorocketnewsen.files.wordpress.com/2016/03/capy-4.jpg",
	"https://i.pinimg.com/originals/7e/60/52/7e6052b52add1adbde9fb2888312e898.jpg",
	"https://upload.wikimedia.org/wikipedia/commons/thumb/b/bc/Bristol.zoo.capybara.arp.jpg/1200px-Bristol.zoo.capybara.arp.jpg",
	"https://image.jimcdn.com/app/cms/image/transf/dimension=1920x10000:format=jpg/path/s5dde8bff85c81b2f/image/iafb738a20225a5d5/version/1476089104/capybara-cabiai-famille-bebe-parent.jpg",
	"https://upload.wikimedia.org/wikipedia/commons/thumb/e/e1/Cattle_tyrant_%28Machetornis_rixosa%29_on_Capybara.jpg/800px-Cattle_tyrant_%28Machetornis_rixosa%29_on_Capybara.jpg",
	"https://i.pinimg.com/originals/11/3f/26/113f26c0109eb801a665d6b17925b0be.jpg",
	"https://i.pinimg.com/originals/3f/49/1d/3f491d524ef380ff7d6e040de803d73e.jpg",
}

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	var err error
	// Create a new Discord session using the provided bot token.
	dg, err = discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	loadSounds()

	dg.AddHandler(messageEvent)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func loadSounds() {
	jsonFile, err := os.Open("sounds.json")
	defer jsonFile.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err := json.Unmarshal(byteValue, &sounds); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Sounds Loaded")
		for _, s := range sounds.Sounds {
			fmt.Printf("File : %s Match on : %s\n", s.Path, s.Regex)
		}
	}
}

func voiceEvent(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
	fmt.Println("Echo")
	fmt.Println(vs.UserID)

}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageEvent(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	g, _ := s.Guild(m.GuildID)
	if g != nil {

		content_lowered := strings.ToLower(m.Content)

		if strings.Contains(content_lowered, "!help") {
			response := &discordgo.MessageEmbed{
				Title: "The Useless Bot - SoundBox",
			}
			for _, s := range sounds.Sounds {
				response.Description += s.Text
				response.Description += "\n"
			}
			s.ChannelMessageSendEmbed(m.ChannelID, response)
		}

		if strings.Contains(content_lowered, "krjone") {
			s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
		}

		v := g.VoiceStates
		var voiceStateOfUser *discordgo.VoiceState
		for _, ve := range v {
			if ve.UserID == m.Author.ID {
				voiceStateOfUser = ve
			}
		}

		if voiceStateOfUser != nil {

			for _, sound := range sounds.Sounds {
				if re, err := regexp.Compile(sound.Regex); err != nil {
					fmt.Println(err)
				} else {

					if re.Match([]byte(content_lowered)) {
						response := &discordgo.MessageEmbed{
							Title:       "The Useless Bot - SoundBox",
							Description: sound.Text,
						}
						s.ChannelMessageSendEmbed(m.ChannelID, response)
						playSound(s, m.ChannelID, voiceStateOfUser.GuildID, voiceStateOfUser.ChannelID, sound.Path)
						break
					}
				}
			}
		}
		go sendEmojiIfContains(s, m, "\U0001F409", "dragon")
		go sendEmojiIfContains(s, m, getEmojiFromFromGuild(g, "capybara"), "capybara")
		go sendCapybaraPicture(s, m)
	} else {
		go sendPM(s, m)
		return
	}

}

func getEmojiFromFromGuild(g *discordgo.Guild, emoji_name string) string {
	emojiId := ""
	if g != nil {
		for _, i := range g.Emojis {
			if strings.Contains(strings.ToLower(i.Name), strings.ToLower(emoji_name)) {
				emojiId = i.APIName()
			}
		}
	}
	return emojiId

}

func sendEmojiIfContains(s *discordgo.Session, m *discordgo.MessageCreate, emojiID, match string) {
	is_match := true
	content_lowered := strings.ToLower(m.Content)
	for _, c := range match {
		if strings.IndexRune(content_lowered, c) == -1 {
			is_match = false
		}
	}
	if is_match {
		if err := s.MessageReactionAdd(m.ChannelID, m.Message.ID, emojiID); err != nil {
			fmt.Println(err)
		}
	}
}

func sendCapybaraPicture(s *discordgo.Session, m *discordgo.MessageCreate) {
	content_lowered := strings.ToLower(m.Content)
	if strings.Contains(content_lowered, "capybara") {
		i := rand.Intn(len(capybara_resource) - 1)
		image := &discordgo.MessageEmbedImage{
			URL: capybara_resource[i],
		}
		response := &discordgo.MessageEmbed{
			Image: image,
		}
		s.ChannelMessageSendEmbed(m.ChannelID, response)
	}
}

func sendPM(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelTyping(m.ChannelID)
	time.Sleep(1 * time.Second)
	s.ChannelMessageSend(m.ChannelID, "Pourquoi tu m'envoie un message "+m.Author.Username)
	s.ChannelTyping(m.ChannelID)
	time.Sleep(3 * time.Second)
	s.ChannelMessageSend(m.ChannelID, "ça sert à rien putain, t'es con ou quoi ?")
}

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, MessageChannelID, guildID, vocalChannelID, filename string) (err error) {
	if isPlaying == false {

		isPlaying = true

		// Join the provided voice channel.
		vc, err := s.ChannelVoiceJoin(guildID, vocalChannelID, false, true)
		if err != nil {
			return err
		}

		// Sleep for a specified amount of time before playing the sound
		//time.Sleep(50 * time.Millisecond)

		// Start speaking.
		vc.Speaking(true)
		c := make(chan bool)
		dgvoice.PlayAudioFile(vc, filename, c)
		c <- true

		// Stop speaking
		vc.Speaking(false)

		// Disconnect from the provided voice channel.
		//vc.Disconnect()

		isPlaying = false

		//time.Sleep(250 * time.Second)

		return nil
	} else {
		response := &discordgo.MessageEmbed{
			Description: "Useless Bot already doing useless things !! ",
		}
		s.ChannelMessageSendEmbed(MessageChannelID, response)

		return nil
	}

}
>>>>>>> a70cba95e57aa41a3777597c483d9ba69fb65c51
