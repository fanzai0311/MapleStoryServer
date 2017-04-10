package main

type MapleCharacter struct {
	id                  uint32
	accid               int
	name                string
	level               uint8
	str, dex, luk, int_ uint16
	remainingAP         uint16
	remainingSP         uint16
	face                uint32
	hair                uint32
	gender              uint8
	skincolor           uint8
	gmlevel             uint8
	job                 uint16
	hp                  uint16
	maxhp               uint16
	mp                  uint16
	maxmp               uint16
	exp                 uint32
	fame                uint16
	mapid               uint32
	spawnPoint          uint8
	meso                uint32
	map_                *map[uint32]*MapleMap
	client              *MapleClient
	position            MaplePosition
}

type MapleInventory struct {
	inventory     map[string]string
	slotLimit     int8
	inventoryType int8
}
type MaplePosition struct {
	x uint16
	y uint16
}

const (
	UNDEFINED int8 = 0
	EQUIP     int8 = 1
	USE       int8 = 2
	SETUP     int8 = 3
	ETC       int8 = 4
	CASH      int8 = 5
	EQUIPPED  int8 = -1
)

func (c *MapleCharacter) getMap() *MapleMap {
	return MapleMaps[c.mapid]
}
