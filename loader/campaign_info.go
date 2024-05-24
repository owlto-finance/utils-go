package loader

import (
	"fmt"
	"github.com/go-lark/lark"
	"github.com/owlto-finance/utils-go/alert"
	"gorm.io/gorm"
	"log"
	"time"

	"gorm.io/datatypes"
)

type CampaignInfo struct {
	*campaignInfoPO
	Rule *campaignRulesPO
}

type campaignInfoPO struct {
	Id        uint64
	Name      string
	StartTime *time.Time
	EndTime   *time.Time
	Status    int8
	ChainId   string
	ChainName string
	Direction int8
	LogoUrl   string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (*campaignInfoPO) TableName() string {
	return "t_campaign_info"
}

type campaignRulesPO struct {
	Id   uint64
	Rule datatypes.JSON
}

func (*campaignRulesPO) TableName() string {
	return "t_campaign_rules"
}

type CampaignManager struct {
	campaignsInfo    []*CampaignInfo
	campaignsNameMap map[string]*CampaignInfo
	campaignsIdMap   map[uint64]*CampaignInfo

	chatID  string
	larkBot *alert.Bot
	db      *gorm.DB
}

func NewCampaignManager(db *gorm.DB, larkBot *alert.Bot, chatID string) *CampaignManager {
	return &CampaignManager{
		campaignsInfo:    make([]*CampaignInfo, 0),
		campaignsNameMap: make(map[string]*CampaignInfo),
		campaignsIdMap:   make(map[uint64]*CampaignInfo),
		chatID:           chatID,
		larkBot:          larkBot,
		db:               db,
	}
}

func (mgr *CampaignManager) LoadAllCampaignsInfo() {
	var campaignsInfoPOs []*campaignInfoPO
	var campaignRulesPOs []*campaignRulesPO

	var campaignsInfo []*CampaignInfo
	var campaignsNameMap = make(map[string]*CampaignInfo)
	var campaignsIdMap = make(map[uint64]*CampaignInfo)
	if err := mgr.db.Model(&campaignInfoPO{}).Find(&campaignsInfoPOs).Error; err != nil {
		_, _ = mgr.larkBot.PostText(fmt.Sprintf("db find t_campaign_info err: %v", err), lark.WithChatID(mgr.chatID))
		return
	}
	if err := mgr.db.Model(&campaignRulesPO{}).Find(&campaignRulesPOs).Error; err != nil {
		_, _ = mgr.larkBot.PostText(fmt.Sprintf("db find t_campaign_rules err: %v", err), lark.WithChatID(mgr.chatID))
		return
	}

	var ruleIdMap = make(map[uint64]*campaignRulesPO)
	for _, rule := range campaignRulesPOs {
		ruleIdMap[rule.Id] = rule
	}

	for _, po := range campaignsInfoPOs {
		tmp := &CampaignInfo{
			campaignInfoPO: po,
			Rule:           ruleIdMap[po.Id],
		}
		campaignsInfo = append(campaignsInfo, tmp)
		campaignsNameMap[po.Name] = tmp
		campaignsIdMap[po.Id] = tmp
	}
	mgr.campaignsInfo = campaignsInfo
	mgr.campaignsNameMap = campaignsNameMap
	mgr.campaignsIdMap = campaignsIdMap
	log.Println("load all campaign info: ", len(campaignsInfo))
}

func (mgr *CampaignManager) GetCampaignInfoById(id uint64) *CampaignInfo {
	return mgr.campaignsIdMap[id]
}

func (mgr *CampaignManager) GetCampaignInfoByName(name string) *CampaignInfo {
	return mgr.campaignsNameMap[name]
}

func (mgr *CampaignManager) GetAllCampaigns() []*CampaignInfo {
	return mgr.campaignsInfo
}

func (mgr *CampaignManager) GetChatID() string {
	return mgr.chatID
}
