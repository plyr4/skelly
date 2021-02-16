package emojis

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/davidvader/skelly/types"
	"github.com/davidvader/skelly/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

const (
	// iamcalEmojiDataURL is the raw github usercontent url for fetching data for all emojis
	iamcalEmojiDataURL = "https://raw.githubusercontent.com/iamcal/emoji-data/master/emoji.json"
)

var (
	// emojis is the global object for storing all emoji data
	emojis map[string]types.Emoji

	// refreshInterval is the time interval to automatically refresh cached emojis
	// default: 12 hours
	refreshInterval = time.Hour * 12

	//emojiMutex
	emojiMutex = sync.Mutex{}
)

// Load fetches emoji data from github and loads all emojis into memory
func Load() error {

	// lock the emojis cache
	emojiMutex.Lock()

	// unlock the cache when complete
	defer emojiMutex.Unlock()

	// retrieve and merge all emoji from github and slack api
	e, err := getAllEmojis()
	if err != nil {
		err = errors.Wrap(err, "could not get all emojis")
		return err
	}

	logrus.Info("storing emojis in memory")

	// set global var
	emojis = e

	// create a thread to refresh the emojis cache
	go func() {

		// wait for refresh interval
		time.Sleep(refreshInterval)

		logrus.Info("refreshing emojis in memory")

		err := Load()
		if err != nil {
			err = errors.Wrap(err, "could not refresh emojis in memory")
			logrus.Error(err)
		}
	}()
	return nil
}

// getAllEmojis fetches emoji data from github and parses into slice
func getAllEmojis() (map[string]types.Emoji, error) {

	// get all emojis from github
	githubEmojis, err := getGithubEmojis()
	if err != nil {
		err = errors.Wrap(err, "could not get github emojis")
		return nil, err
	}

	// get all custom emojis from slack
	slackEmojis, err := getSlackEmojis()
	if err != nil {
		err = errors.Wrap(err, "could not get custom slack emojis")
		return nil, err
	}

	// merge the two maps
	e := mergeMaps(githubEmojis, slackEmojis)
	return e, nil
}

// getEmojisFromGithub gets emoji data from iamcal github usercontent
func getGithubEmojis() (map[string]types.Emoji, error) {

	logrus.Trace("getting emoji data from github")

	// fetch emoji data from github
	resp, err := http.Get(iamcalEmojiDataURL)
	if err != nil {
		err = errors.Wrap(err, "could not get emoji data")
		return nil, err
	}

	// read the request body
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "could not read request body")
		return nil, err
	}

	// unmarshal emoji data
	eList := []types.Emoji{}

	err = json.Unmarshal(b, &eList)
	if err != nil {
		err = errors.Wrap(err, "could not unmarshal emoji list data")
		return nil, err
	}

	// check for data
	if len(eList) == 0 {
		return nil, errors.New("no github emojis unmarshalled")
	}

	logrus.Tracef("got (%v) emojis from github", len(eList))

	// slice to map
	eMap := map[string]types.Emoji{}
	for _, e := range eList {

		_, ok := eMap[e.ShortName]
		if !ok {
			eMap[e.ShortName] = e
		} else {
			emoji := eMap[e.ShortName]
			emoji.ShortNames = util.Unique(append(emoji.ShortNames, e.ShortNames...))
			eMap[e.ShortName] = emoji
		}
	}

	logrus.Tracef("merged (%v) github emojis into (%v)", len(eList), len(eMap))

	return eMap, nil
}

// getSlackEmojis gets emoji data from iamcal github usercontent
func getSlackEmojis() (map[string]types.Emoji, error) {

	logrus.Trace("getting emoji data from slack")

	// create an api client
	bToken := os.Getenv("SKELLY_BOT_TOKEN")

	api := slack.New(bToken)

	// fetch custom emoji from slack api
	customEmojis, err := api.GetEmoji()
	if err != nil {
		return nil, errors.Wrap(err, "could not get emoji from slack api")
	}

	// check for data
	if len(customEmojis) == 0 {
		return nil, errors.New("no slack emojis unmarshalled")
	}

	logrus.Tracef("got (%v) custom emojis from slack", len(customEmojis))

	// convert slack emoji to internal emoji
	e := map[string]types.Emoji{}

	// convert all alias emojis pointing to the same short name into a slice of shortnames
	for alias, value := range customEmojis {

		// find true emoji shortname
		split := strings.Split(value, ":")

		shortname := alias
		if len(split) > 1 {

			// check for alias pseudo protocol
			if split[0] == "alias" {
				shortname = split[1]
			}
		}

		// add the alias to the emoji map
		_, ok := e[shortname]
		if !ok {
			emoji := types.Emoji{
				ShortName:  shortname,
				ShortNames: []string{alias},
			}
			e[shortname] = emoji
		} else {
			emoji := e[shortname]
			emoji.ShortNames = append(e[shortname].ShortNames, alias)
			e[shortname] = emoji
		}
	}

	logrus.Tracef("merged (%v) slack emojis into (%v)", len(customEmojis), len(e))

	return e, nil
}

// findEmoji retrieves a specific emoji from loaded emoji data
func findEmoji(alias string) (*types.Emoji, error) {

	// lock the emojis cache
	emojiMutex.Lock()

	// unlock the cache when complete
	defer emojiMutex.Unlock()

	// check for valid emoji data
	if len(emojis) == 0 {
		return nil, errors.New("no emoji data loaded")
	}

	// find the emoji by shortname
	var emoji *types.Emoji
	for _, e := range emojis {

		// check all shortnames for the emoji
		for _, s := range e.ShortNames {
			if alias == s {
				emoji = &e
				break
			}
		}

		// found the emoji
		if emoji != nil {
			break
		}
	}

	// check for found emoji
	if emoji == nil {
		return nil, errors.Errorf("could not find emoji(%s)", alias)
	}

	return emoji, nil
}

// GetShortname takes emoji alias and returns the core shortname used for the emoji
func GetShortname(alias string) (string, error) {

	// retrieve the emoji from cached emoji data
	emoji, err := findEmoji(alias)
	if err != nil {
		return "", errors.Wrap(err, "could not find emoji")
	}

	// extract shortname
	shortname := emoji.ShortName
	if shortname == "" {
		return "", errors.New("no shortname")
	}

	return shortname, nil
}

// mergeMaps loops over two emoji maps and merges all duplicate shortnames
func mergeMaps(m1, m2 map[string]types.Emoji) map[string]types.Emoji {

	// merge maps
	m := mergeShortNames(mergeShortNames(map[string]types.Emoji{}, m1), m2)

	logrus.Tracef("merged (%v) github emojis and (%v) slack emojis into (%v)", len(m1), len(m2), len(m))

	return m
}

// mergeShortNames takes two emoji maps, iterates over and merges all duplicate shortnames
func mergeShortNames(m1, m2 map[string]types.Emoji) map[string]types.Emoji {

	// iterate over input map
	for _, e := range m2 {

		// check in new map if emoji exists
		_, ok := m1[e.ShortName]
		if !ok {

			// initialize emoji
			m1[e.ShortName] = e
		} else {

			// merge shortnames
			emoji := m1[e.ShortName]
			logrus.Tracef("merging emoji shortnames (%v) and (%v) for (%s)", e.ShortNames, emoji.ShortNames, emoji.ShortName)
			emoji.ShortNames = util.Unique(append(emoji.ShortNames, e.ShortNames...))
			m1[e.ShortName] = emoji
		}
	}
	return m1
}
