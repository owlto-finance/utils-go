package loader

import (
	"database/sql"
	"github.com/owlto-finance/utils-go/log"
)

type MakerAddressGroupPO struct {
	Id        int64
	GroupName string
	Env       string
}

type MakerAddressPO struct {
	Id      int64
	GroupId int64
	Backend Backend
	Address string
}

type MakerAddress struct {
	GroupId   int64
	GroupName string
	Env       string
	Addresses []*MakerAddressPO
}

type MakerAddressManager struct {
	groupIdAddress map[int64]*MakerAddress
	envGroup       map[string][]*MakerAddress

	db *sql.DB
}

func NewMakerAddressManager(db *sql.DB) *MakerAddressManager {
	return &MakerAddressManager{
		groupIdAddress: make(map[int64]*MakerAddress),
		envGroup:       make(map[string][]*MakerAddress),
		db:             db,
	}
}

func (mgr *MakerAddressManager) LoadAllMakerAddresses() {
	// Query the database for all maker address groups
	groupRows, err := mgr.db.Query("SELECT id, group_name, env FROM maker_address_groups")
	if err != nil || groupRows == nil {
		log.Errorf("select maker_address_groups error: %v", err)
		return
	}
	defer groupRows.Close()

	groups := make(map[int64]*MakerAddress)
	for groupRows.Next() {
		var group MakerAddressGroupPO
		if err = groupRows.Scan(&group.Id, &group.GroupName, &group.Env); err != nil {
			log.Errorf("scan maker_address_groups row error: %v", err)
			continue
		}

		makerAddress := &MakerAddress{
			GroupId:   group.Id,
			GroupName: group.GroupName,
			Env:       group.Env,
			Addresses: []*MakerAddressPO{},
		}
		groups[group.Id] = makerAddress
	}

	// Check for errors from iterating over rows
	if err = groupRows.Err(); err != nil {
		log.Errorf("get next maker_address_groups row error: %v", err)
		return
	}

	// Query the database for all maker addresses
	addressRows, err := mgr.db.Query("SELECT id, group_id, backend, address FROM maker_addresses")
	if err != nil || addressRows == nil {
		log.Errorf("select maker_addresses error: %v", err)
		return
	}
	defer addressRows.Close()

	for addressRows.Next() {
		var address MakerAddressPO
		if err = addressRows.Scan(&address.Id, &address.GroupId, &address.Backend, &address.Address); err != nil {
			log.Errorf("scan maker_addresses row error: %v", err)
			continue
		}

		if group, ok := groups[address.GroupId]; ok {
			group.Addresses = append(group.Addresses, &address)
		}
	}

	if err = addressRows.Err(); err != nil {
		log.Errorf("get next maker_addresses row error: %v", err)
		return
	}

	mgr.groupIdAddress = groups
	envGroup := make(map[string][]*MakerAddress)
	for _, group := range groups {
		envGroup[group.Env] = append(envGroup[group.Env], group)
	}
	mgr.envGroup = envGroup
	log.Infof("load all maker addresses: %d", len(groups))
}

func (mgr *MakerAddressManager) GetMakerAddressesByEnv(env string) []*MakerAddress {
	return mgr.envGroup[env]
}

func (mgr *MakerAddressManager) GetMakerAddressByGroupId(groupId int64) *MakerAddress {
	return mgr.groupIdAddress[groupId]
}
