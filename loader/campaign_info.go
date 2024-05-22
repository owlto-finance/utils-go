package loader

import (
	"fmt"
	"github.com/go-lark/lark"
	"github.com/owlto-finance/utils-go/alert"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"

	"gorm.io/datatypes"
)

type CampaignInfo struct {
	Id        uint64
	Name      string
	StartTime *time.Time
	EndTime   *time.Time
	Status    int8
	ChainId   string
	ChainName string
	Direction int8
	LogoUrl   string
	Tasks     datatypes.JSON
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type CampaignManager struct {
	campaignsInfo    []*CampaignInfo
	campaignsNameMap map[string]*CampaignInfo
	campaignsIdMap   map[uint64]*CampaignInfo

	chatID string
	db     *gorm.DB
	mutex  *sync.RWMutex
}

func NewCampaignManager(db *gorm.DB, chatID string) *CampaignManager {
	return &CampaignManager{
		campaignsInfo:    make([]*CampaignInfo, 0),
		campaignsNameMap: make(map[string]*CampaignInfo),
		campaignsIdMap:   make(map[uint64]*CampaignInfo),
		chatID:           chatID,
		db:               db,
		mutex:            &sync.RWMutex{},
	}
}

func (mgr *CampaignManager) LoadAllCampaignsInfo() {
	var campaignsInfo []*CampaignInfo
	var campaignsNameMap map[string]*CampaignInfo
	var campaignsIdMap map[uint64]*CampaignInfo
	if err := mgr.db.Find(&campaignsInfo).Error; err != nil {
		_, _ = alert.LarkBot.PostText(fmt.Sprintf("db find t_campaign_info err: %v", err), lark.WithChatID(mgr.chatID))
		return
	}
	for _, campaignInfo := range campaignsInfo {
		campaignsNameMap[campaignInfo.Name] = campaignInfo
		campaignsIdMap[campaignInfo.Id] = campaignInfo
	}
	mgr.mutex.Lock()
	mgr.campaignsInfo = campaignsInfo
	mgr.campaignsNameMap = campaignsNameMap
	mgr.campaignsIdMap = campaignsIdMap
	mgr.mutex.Unlock()
	log.Println("load all campaign info: ", len(campaignsInfo))
}

func (mgr *CampaignManager) GetCampaignInfoById(id uint32) *CampaignInfo {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()
	return mgr.campaignsInfo[id]
}

func (mgr *CampaignManager) GetCampaignInfoByName(name string) *CampaignInfo {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()
	return mgr.campaignsNameMap[name]
}

func (mgr *CampaignManager) GetAllCampaigns() []*CampaignInfo {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()
	return mgr.campaignsInfo
}

func (mgr *CampaignManager) GetChatID() string {
	return mgr.chatID
}
