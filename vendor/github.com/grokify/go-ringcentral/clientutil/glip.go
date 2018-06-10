package clientutil

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	rc "github.com/grokify/go-ringcentral/client"
)

var rxMultiSpace = regexp.MustCompile(`\\s+`)

type GlipApiUtil struct {
	ApiClient *rc.APIClient
}

func AtMention(personId string) string {
	return fmt.Sprintf("![:Person](%v)", personId)
}

func (apiUtil *GlipApiUtil) GlipGroupMemberCount(groupId string) (int64, error) {
	if apiUtil.ApiClient == nil {
		return int64(-1), fmt.Errorf("GlipApiUtil is missing RingCentral ApiClient")
	}
	groupId = strings.ToLower(strings.TrimSpace(groupId))
	grp, resp, err := apiUtil.ApiClient.GlipApi.LoadGroup(context.Background(), groupId)
	if err != nil {
		return int64(-1), err
	} else if resp.StatusCode >= 300 {
		return int64(-1), fmt.Errorf("Glip API Response Code [%v]", resp.StatusCode)
	}
	return int64(len(grp.Members)), nil
}

type GlipInfoAtMentionOrGroupOfTwoInfo struct {
	PersonId       string
	AtMentions     []rc.GlipMentionsInfo
	PersonName     string
	FuzzyAtMention bool
	TextRaw        string
	GroupId        string
}

func (apiUtil *GlipApiUtil) AtMentionedOrGroupOfTwoFuzzy(info GlipInfoAtMentionOrGroupOfTwoInfo) (bool, error) {
	if IsAtMentioned(info.PersonId, info.AtMentions) ||
		IsAtMentionedFuzzy(info.PersonName, info.TextRaw) {
		return true, nil
	}
	count, err := apiUtil.GlipGroupMemberCount(info.GroupId)
	if err != nil || count != int64(2) {
		return false, err
	}
	return true, nil
}

func IsAtMentionedFuzzy(personName, textRaw string) bool {
	personName = strings.ToLower(strings.TrimSpace(personName))
	rx, err := regexp.Compile(`(\A|\W)@` + personName + `\b`)
	if err != nil {
		return false
	}
	str := rx.FindString(strings.ToLower(textRaw))
	if len(str) > 0 {
		return true
	}
	return false
}

func IsAtMentionedGlipdown(personId, textRaw string) bool {
	personIdMarkdownLc := strings.ToLower(AtMention(personId))
	if strings.Index(strings.ToLower(textRaw), personIdMarkdownLc) == -1 {
		return false
	}
	return true
}

func PrefixAtMentionUnlessMentioned(personId, text string) string {
	personId = strings.TrimSpace(personId)
	if len(personId) > 0 && !IsAtMentionedGlipdown(personId, text) {
		return AtMention(personId) + " " + text
	}
	return text
}

// DirectMessage means a group of 2 or a team of 2
func (apiUtil *GlipApiUtil) AtMentionedOrGroupOfTwo(userId, groupId string, mentions []rc.GlipMentionsInfo) (bool, error) {
	if IsAtMentioned(userId, mentions) {
		return true, nil
	}

	count, err := apiUtil.GlipGroupMemberCount(groupId)
	if err != nil {
		return false, err
	}
	if count == int64(2) {
		return true, nil
	}
	return false, nil
}

func IsAtMentioned(userId string, mentions []rc.GlipMentionsInfo) bool {
	for _, mention := range mentions {
		if userId == mention.Id {
			return true
		}
	}
	return false
}

func GlipCreatePostIsEmpty(post rc.GlipCreatePost) bool {
	if len(strings.TrimSpace(post.Text)) == 0 && len(post.Attachments) == 0 {
		return false
	}
	return true
}

func StripAtMentionAll(id, personName, text string) string {
	noAtMention := StripAtMention(id, text)
	return StripAtMentionFuzzy(personName, noAtMention)
}

func StripAtMention(id, text string) string {
	rx := regexp.MustCompile(fmt.Sprintf("!\\[:Person\\]\\(%v\\)", id))
	noAtMention := rx.ReplaceAllString(text, " ")
	return strings.TrimSpace(rxMultiSpace.ReplaceAllString(noAtMention, " "))
}

func StripAtMentionFuzzy(personName, text string) string {
	rx, err := regexp.Compile(`(?i)(\A|\W)@` + personName + `\b`)
	if err != nil {
		return text
	}
	noAtMention := rx.ReplaceAllString(text, " ")
	return strings.TrimSpace(rxMultiSpace.ReplaceAllString(noAtMention, " "))
}
