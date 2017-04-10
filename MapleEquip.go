package main

type MapleEquip struct {
	upgradeSlots                            byte
	level, flag                             byte
	locked                                  byte
	str, dex, _int, luk                     uint16
	hp, mp, watk, matk, wdef, mdef          uint16
	acc, avoid, hands, speed, jump, vicious uint16
	partnerName                             string
	partnerId, partnerUniqueId              uint32
}
