package channels

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"

	"github.com/sipeed/picoclaw/pkg/config"
)

func toChannelHashes(cfg *config.Config) map[string]string {
	result := make(map[string]string)
	ch := cfg.Channels
	// should not be error
	marshal, _ := json.Marshal(ch)
	var channelConfig map[string]map[string]any
	_ = json.Unmarshal(marshal, &channelConfig)

	for key, value := range channelConfig {
		if !value["enabled"].(bool) {
			continue
		}
		hiddenValues(key, value, ch.Get(key))
		valueBytes, _ := json.Marshal(value)
		hash := md5.Sum(valueBytes)
		result[key] = hex.EncodeToString(hash[:])
	}

	return result
}

func hiddenValues(key string, value map[string]any, ch *config.Channel) {
	v, err := ch.GetDecoded()
	if err != nil {
		return
	}
	switch key {
	case "pico":
		value["token"] = v.(*config.PicoSettings).Token.String()
	case "telegram":
		value["token"] = v.(*config.TelegramSettings).Token.String()
	case "discord":
		value["token"] = v.(*config.DiscordSettings).Token.String()
	case "slack":
		value["bot_token"] = v.(*config.SlackSettings).BotToken.String()
		value["app_token"] = v.(*config.SlackSettings).AppToken.String()
	case "matrix":
		value["token"] = v.(*config.MatrixSettings).AccessToken.String()
	case "onebot":
		value["token"] = v.(*config.OneBotSettings).AccessToken.String()
	case "line":
		value["token"] = v.(*config.LINESettings).ChannelAccessToken.String()
		value["secret"] = v.(*config.LINESettings).ChannelSecret.String()
	case "wecom":
		value["secret"] = v.(*config.WeComSettings).Secret.String()
	case "dingtalk":
		value["secret"] = v.(*config.DingTalkSettings).ClientSecret.String()
	case "qq":
		value["secret"] = v.(*config.QQSettings).AppSecret.String()
	case "irc":
		value["password"] = v.(*config.IRCSettings).Password.String()
		value["serv_password"] = v.(*config.IRCSettings).NickServPassword.String()
		value["sasl_password"] = v.(*config.IRCSettings).SASLPassword.String()
	case "feishu":
		value["app_secret"] = v.(*config.FeishuSettings).AppSecret.String()
		value["encrypt_key"] = v.(*config.FeishuSettings).EncryptKey.String()
		value["verification_token"] = v.(*config.FeishuSettings).VerificationToken.String()
	case "teams_webhook":
		// Expose webhook URLs for hash computation (they contain secrets)
		vv := value["webhooks"]
		webhooks := make(map[string]string)
		if vv != nil {
			webhooks = vv.(map[string]string)
		}
		ts := v.(*config.TeamsWebhookSettings)
		for name, target := range ts.Webhooks {
			webhooks[name] = target.WebhookURL.String()
		}
		value["webhooks"] = webhooks
	}
}

func compareChannels(old, news map[string]string) (added, removed []string) {
	for key, newHash := range news {
		if oldHash, ok := old[key]; ok {
			if newHash != oldHash {
				removed = append(removed, key)
				added = append(added, key)
			}
		} else {
			added = append(added, key)
		}
	}
	for key := range old {
		if _, ok := news[key]; !ok {
			removed = append(removed, key)
		}
	}
	return added, removed
}

func toChannelConfig(cfg *config.Config, list []string) (*config.ChannelsConfig, error) {
	result := make(config.ChannelsConfig)
	for _, name := range list {
		bc, ok := cfg.Channels[name]
		if !ok || !bc.Enabled {
			continue
		}
		result[name] = bc
	}
	return &result, nil
}
