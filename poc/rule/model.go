package rule

const (
	AffectContent          = "text"
	AffectDirectory        = "directory"
	AffectURL              = "url"
	AffectAppendParameter  = "appendparam"
	AffectReplaceParameter = "replaceparam"
	AffectServer           = "server"
	AffectScript           = "script"
)

type Rule struct {
	Method          string            `json:"method"`
	Path            string            `json:"path"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	Search          string            `json:"search"`
	FollowRedirects bool              `json:"follow_redirects"`
	Expression      string            `json:"expression"`
}

type Detail struct {
	Author      string   `json:"author"`
	Links       []string `json:"links"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
}

// Rules 和 Groups 只能存在一个
type Poc struct {
	Params	[]string	 	 `json:"params"`
	Name   string            `json:"name"`
	Set    map[string]string `json:"set"`
	Rules  []Rule           `json:"rules"`
	Groups  map[string][]Rule `json:"groups"`
	Detail Detail            `json:"detail"`
}

type Plugin struct {
	VulId   string `gorm:"column:vul_id"` // 漏洞编号
	Affects string `gorm:"column:affects"`   // 影响类型  dir/server/param/url/content
	JsonPoc *Poc   `gorm:"column:json_poc"`  // json规则
	Enable  bool   `gorm:"column:enable"`    // 是否启用
}