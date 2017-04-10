package main

import (
	"bytes"
	"encoding/binary"
	//"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	. "strconv"
	//"time"

	"github.com/Francesco149/maplelib"
)

var CHAR_INFO_MAGIC []byte = []byte{0xFF, 0xC9, 0x9A, 0x3B}
var FILETIMESTAMP []byte = []byte{0x10, 0x3D, 0x4E, 0x6C, 0xBD, 0xA7, 0xD2, 0x01}
var FILETIMESTAMP2 []byte = []byte{0x40, 0xB2, 0x4E, 0x6C, 0xBD, 0xA7, 0xD2, 0x01}

//握手包
func GetHello(mapleVersion uint16, sendIv [4]byte, recIv [4]byte) []byte {
	mw := new(bytes.Buffer)
	var data = []interface{}{
		uint16(0x0d),
		uint16(mapleVersion),
		[]byte{0, 0},
		recIv,
		sendIv,
		byte(4),
	}
	for _, v := range data {
		err := binary.Write(mw, binary.LittleEndian, v)

		CheckError("GetHello", err)
	}
	return mw.Bytes()
}

// LOGIN 代码段开始

/*
	登录失败
	reason
	0 - 成功
	2 - 永久封号
	3 - 已被封号
	4 - 密码错误
	5 - 账号不存在
	7 - 已经登录
	22- 许可协议
*/
func GetLoginFailed(reason int32) []byte {

	p := NewPacket()
	p.Encode2(SH.LOGIN_STATUS) //LOGIN_STATUS 登录状态包头
	p.Encode4s(reason)
	p.Encode2(0)
	if reason == 2 {
		p.Append(FILETIMESTAMP)
	}
	return []byte(p)
}

// 需要选择性别
func GenderNeeded(c *MapleClient) []byte {
	p := NewPacket()
	p.Encode2(SH.CHOOSE_GENDER) //CHOOSE_GENDER  选择性别包头
	p.EncodeString(c.AccountName)
	return []byte(p)
}

// 选择性别返回
func GenderChanged(c *MapleClient) []byte {
	p := NewPacket()
	p.Encode2(SH.GENDER_SET) // GENDER_SET 选择性别返回
	p.Encode1(0)
	p.EncodeString(c.AccountName)
	p.EncodeString(FormatInt(int64(c.AccId), 10))
	return []byte(p)
}

//用户协议返回
func LicenseResult() []byte {
	p := NewPacket()
	p.Encode2(SH.LICENSE_RESULT) //LICENSE_RESULT 用户协议返回
	p.Encode1(1)
	return []byte(p)
}

//登录成功
func getLoginSuccess(c *MapleClient) []byte {
	p := NewPacket()
	p.Encode2(SH.LOGIN_STATUS) //LOGIN_STATUS 登录状态包头
	p.Encode4(1632938240)
	p.Encode1(0)
	p.Encode1s(c.Gender)
	p.Encode2(0)
	p.EncodeGBKMapleString(c.AccountName)
	//b, err := hex.DecodeString("00 00 00 03 01 00 00 00 E2 ED A3 7A FA C9 01")
	//Checkerr(err)
	//p.EncodeBuffer(b)
	p.Append([]byte{0x00, 0x00, 0x00, 0x03, 0x01, 0x00, 0x00, 0x00, 0xE2, 0xED, 0xA3, 0x7A, 0xFA, 0xC9, 0x01})

	p.Encode4(0)
	p.Encode8(0)
	p.EncodeGBKMapleString(Itoa(c.AccId))
	p.EncodeGBKMapleString(c.AccountName)
	p.Encode1(1)
	return []byte(p)
}

//服务器列表
func getServerList() []byte {
	p := NewPacket()
	p.Encode2(SH.SERVERLIST) // SERVERLIST 服务器列表
	p.Encode1(0)
	p.EncodeString(ServerName)
	p.Encode1s(int8(rand.Intn(4)))

	p.EncodeGBKMapleString("欢迎来到" + ServerName)
	p.Encode1(0x64)
	p.Encode1(0)
	p.Encode1(0x64)
	p.Encode1(0)

	p.Encode1(1)
	p.Encode4(500)

	p.EncodeGBKMapleString(ServerName + "-1")
	p.Encode4(1200)
	p.Encode1(0)
	p.Encode2(0)

	p.Encode2(0)
	return []byte(p)
}

//服务器列表结尾
func getEndOfServerList() []byte {
	p := NewPacket()
	p.Encode2(SH.SERVERLIST)
	p.Encode1(0xFF)
	return []byte(p)
}

//服务器状态
func getServerStatus() []byte {
	p := NewPacket()
	p.Encode2(SH.SERVERSTATUS)
	p.Encode1s(int8(rand.Intn(2)))
	return []byte(p)
}

//角色列表
func getCharList(c *MapleClient) []byte {
	p := NewPacket()
	p.Encode2(SH.CHARLIST)
	p.Encode1(0)
	p.Encode4(0)
	p.Encode1(1)

	var char MapleCharacter
	char.name = "管理员001"
	char.id = 1
	char.level = 1
	char.str = 4
	char.dex = 4
	char.int_ = 4
	char.luk = 4
	char.hp = 50
	char.maxhp = 50
	char.mp = 50
	char.maxmp = 50
	char.remainingAP = 9
	char.mapid = 102000000
	char.spawnPoint = 11
	char.face = 20100
	char.hair = 30000
	c.Character = char
	p.Append(addCharEntry(&char))
	p.Encode2(3)
	p.Encode4(6)
	return []byte(p)
}

//检查角色名
func getCheckName(charname string) []byte {
	p := NewPacket()
	p.Encode2(SH.CHAR_NAME_RESPONSE)
	p.EncodeGBKMapleString(charname)
	p.Encode1(0)
	return []byte(p)
}

//新建角色
func addNewCharEntry(char *MapleCharacter) []byte {
	p := NewPacket()
	p.Encode2(SH.ADD_NEW_CHAR_ENTRY)
	p.Encode1(1)
	p.Append(addCharEntry(char))
	return []byte(p)
}

//角色
func addCharEntry(char *MapleCharacter) []byte {
	p := NewPacket()
	p.Append(addCharStatus(char))
	p.Append(addCharLook(char))
	p.Encode1(0)
	return []byte(p[4:])
}

//角色状态
func addCharStatus(char *MapleCharacter) []byte {
	p := NewPacket()
	p.Encode4(char.id)
	p.EncodeNameString(char.name)
	p.Encode1(char.gender)
	p.Encode1(char.skincolor)
	p.Encode4(char.face)
	p.Encode4(char.hair)
	p.Encode8(0)
	p.Encode8(0)
	p.Encode8(0)
	p.Encode1(char.level)
	p.Encode2(char.job)
	p.Encode2(char.str)
	p.Encode2(char.dex)
	p.Encode2(char.int_)
	p.Encode2(char.luk)
	p.Encode2(char.hp)
	p.Encode2(char.maxhp)
	p.Encode2(char.mp)
	p.Encode2(char.maxmp)
	p.Encode2(char.remainingAP)
	p.Encode2(char.remainingSP)
	p.Encode4(char.exp)
	p.Encode2(char.fame)
	p.Encode4(0)
	p.Append(FILETIMESTAMP)
	p.Encode4(char.mapid)
	p.Encode1(char.spawnPoint)
	return []byte(p[4:])
}

//角色外观
func addCharLook(char *MapleCharacter) []byte {
	p := NewPacket()
	p.Encode1(char.gender)
	p.Encode1(char.skincolor)
	p.Encode4(char.face)
	p.Encode1(1)
	p.Encode4(char.hair)
	p.Encode1(0xFF)
	p.Encode1(0xFF)
	p.Encode4(0)
	p.Encode4(0)
	p.Encode8(0)
	return []byte(p[4:])
}

//获取频道IP
func getServerIP(charid uint32) []byte {
	p := NewPacket()
	p.Encode2(SH.SERVER_IP)
	p.Encode2(0)
	p.Append(ServerIP)
	p.Encode2(9001)
	p.Encode4(charid)
	p.Encode1(1)
	p.Encode4(0)
	return []byte(p)
}

//获取角色详情
func getCharInfo(char *MapleCharacter) []byte {
	p := NewPacket()
	p.Encode2(SH.WARP_TO_MAP)
	p.Encode4(0) //频道ID-1
	p.Encode1(0)
	p.Encode1(1)
	p.Encode1(1)
	p.Encode2(0)
	p.Encode4(0)
	p.Append(FILETIMESTAMP)
	p.Encode8(^uint64(0))
	p.Encode1(0)
	p.Append(addCharStatus(char))
	p.Encode1(20) // 好友上限
	p.Append(addInventoryInfo(char))
	p.Append(addSkillRecord(char))
	p.Append(addQuestRecord(char))
	p.Append(addRingInfo(char))
	p.Append(addTeleportRockRecord(char))
	p.Encode2(0)
	p.Encode8(0)
	p.Encode1(0)
	p.Encode4(0)
	p.Append(FILETIMESTAMP2)
	return []byte(p)
}

func addInventoryInfo(char *MapleCharacter) []byte {
	p := NewPacket()
	p.Encode1(1)
	p.EncodeGBKMapleString(char.name)
	p.Encode4(char.meso)
	p.Encode4(char.id)
	p.Encode8(0)
	p.Encode1(24) //EQUIP slots 装备栏上限
	p.Encode1(24) //USE slots 消耗栏上限
	p.Encode1(24) //SETUP 设置栏上限
	p.Encode1(24) //ETC 其他栏上限
	p.Encode1(50) //CASH 现金栏上限
	p.Append(FILETIMESTAMP)
	//equipped
	p.Encode1(0) //equipped cash start
	p.Encode1(0) //equip
	p.Encode1(0) //use
	p.Encode1(0) //setup
	p.Encode1(0) //etc
	p.Encode1(0) //cash
	return []byte(p[4:])
}
func addSkillRecord(char *MapleCharacter) []byte {
	p := NewPacket()
	p.Encode1(0)
	p.Encode2(0) //skills count
	p.Encode2(0) //skillsCD count
	return []byte(p[4:])
}
func addQuestRecord(char *MapleCharacter) []byte {
	p := NewPacket()
	p.Encode2(0) //quest started count
	p.Encode2(0) //quest completed count
	return []byte(p[4:])
}
func addRingInfo(char *MapleCharacter) []byte {
	p := NewPacket()
	p.Encode2(0)
	p.Encode2(0)
	p.Encode2(0)
	p.Encode2(0)
	return []byte(p[4:])
}
func addTeleportRockRecord(char *MapleCharacter) []byte {
	p := NewPacket()

	for i := 0; i < 15; i++ {
		p.Append(CHAR_INFO_MAGIC)
	}
	return []byte(p[4:])
}

func weirdStatusUpdate() []byte {
	p := NewPacket()
	p.Encode2(SH.UPDATE_STATUS)
	p.Encode1(0)
	p.Encode1(0x38)
	p.Encode2(0)
	p.Encode8(0)
	p.Encode8(0)
	p.Encode8(0)
	p.Encode1(0)
	p.Encode1(1)
	return []byte(p)
}

//召唤玩家
func spawnPlayerMapobject(char *MapleCharacter) []byte {
	p := NewPacket()
	p.Encode2(SH.SPAWN_PLAYER)
	p.Encode4(char.id)
	p.Encode1(0)
	p.EncodeGBKMapleString(char.name)
	p.EncodeGBKMapleString("")
	p.Encode4(0)
	p.Encode2(0)
	p.Encode4(0)
	p.Encode1(0)
	p.Encode1(0xE0)
	p.Encode1(0x1F)
	p.Encode1(0)
	p.Encode1(0)

	p.Encode4(0)
	p.Encode4(0)
	p.Encode4(0)
	p.Encode4(0)
	p.Encode2(0)
	char_spawn_magic := []byte{0xE8, 0xC2, 0x74, 0x7E}
	p.Append(char_spawn_magic) //1
	p.Encode8(0)
	p.Encode2(0)
	p.Encode1(0)
	p.Append(char_spawn_magic) //2
	p.Encode8(0)
	p.Encode2(0)
	p.Encode1(0)
	p.Append(char_spawn_magic) //3
	p.Encode2(0)
	p.Encode1(0)
	p.Append(char_spawn_magic) //4
	p.Encode8(0)
	p.Encode1(0)
	p.Encode8(0)
	p.Append(char_spawn_magic) //5
	p.Encode1(0)
	p.Encode1(1)
	p.Append([]byte{0x41, 0x9A, 0x70, 0x07})
	p.Encode8(0)
	p.Encode1(0)
	p.Append(char_spawn_magic) //6
	p.Encode8(0)
	p.Encode4(0)
	p.Encode1(0)
	p.Append(char_spawn_magic) //7
	p.Encode8(0)
	p.Encode2(0)
	p.Encode1(0)
	p.Append(char_spawn_magic) //8
	p.Encode1(0)
	p.Encode2(char.job)
	p.Append(addCharLook(char))
	p.Encode8(0)
	p.Encode4(0)
	p.Encode4(^uint32(0))
	p.Encode4(0)
	p.Encode2(char.position.x) //x
	p.Encode2(char.position.y) //y
	p.Encode1(0)               //姿势
	p.Encode1(0)
	p.Encode2(0)
	p.Encode4(1)
	p.Encode8(0)
	p.Encode1(0)
	p.Encode2(0)
	p.Encode4(0)
	return []byte(p)
}

//匹配移动封包
func parseMovement(it *maplelib.PacketIterator) []byte {
	p := NewPacket()
	num, err := it.Decode1()
	Checkerr(err)
	p.Encode1(num)
	for i := 0; i < int(num); i++ {
		command, err := it.Decode1()
		Checkerr(err)
		switch command {
		case 0, 5, 17:
			x, err := it.Decode2()
			Checkerr(err)
			y, err := it.Decode2()
			Checkerr(err)
			ppsx, err := it.Decode2()
			Checkerr(err)
			ppsy, err := it.Decode2()
			Checkerr(err)
			unk, err := it.Decode2()
			Checkerr(err)
			newstate, err := it.Decode1()
			Checkerr(err)
			duration, err := it.Decode2()
			Checkerr(err)

			p.Encode1(command)
			p.Encode2(x)
			p.Encode2(y)
			p.Encode2(ppsx)
			p.Encode2(ppsy)
			p.Encode2(unk)
			p.Encode1(newstate)
			p.Encode2(duration)
		case 1, 2, 6, 12, 13, 16:
			x, err := it.Decode2()
			Checkerr(err)
			y, err := it.Decode2()
			Checkerr(err)
			newstate, err := it.Decode1()
			Checkerr(err)
			duration, err := it.Decode2()
			Checkerr(err)
			p.Encode1(command)
			p.Encode2(x)
			p.Encode2(y)
			p.Encode1(newstate)
			p.Encode2(duration)
		case 3, 4, 7, 8, 9, 14:
			x, err := it.Decode2()
			Checkerr(err)
			y, err := it.Decode2()
			Checkerr(err)
			ppsx, err := it.Decode2()
			Checkerr(err)
			ppsy, err := it.Decode2()
			Checkerr(err)
			newstate, err := it.Decode1()
			Checkerr(err)
			p.Encode1(command)
			p.Encode2(x)
			p.Encode2(y)
			p.Encode2(ppsx)
			p.Encode2(ppsy)
			p.Encode1(newstate)
		case 10:
			wui, err := it.Decode1()
			Checkerr(err)
			p.Encode1(10)
			p.Encode1(wui)
		case 11:
			x, err := it.Decode2()
			Checkerr(err)
			y, err := it.Decode2()
			Checkerr(err)
			unk, err := it.Decode2()
			Checkerr(err)
			newstate, err := it.Decode1()
			Checkerr(err)
			duration, err := it.Decode2()
			Checkerr(err)

			p.Encode1(command)
			p.Encode2(x)
			p.Encode2(y)
			p.Encode2(unk)
			p.Encode1(newstate)
			p.Encode2(duration)
		case 15:
			x, err := it.Decode2()
			Checkerr(err)
			y, err := it.Decode2()
			Checkerr(err)
			ppsx, err := it.Decode2()
			Checkerr(err)
			ppsy, err := it.Decode2()
			Checkerr(err)
			unk, err := it.Decode2()
			Checkerr(err)
			fh, err := it.Decode2()
			Checkerr(err)
			newstate, err := it.Decode1()
			Checkerr(err)
			duration, err := it.Decode2()
			Checkerr(err)

			p.Encode1(command)
			p.Encode2(x)
			p.Encode2(y)
			p.Encode2(ppsx)
			p.Encode2(ppsy)
			p.Encode2(unk)
			p.Encode2(fh)
			p.Encode1(newstate)
			p.Encode2(duration)
		case 20, 21, 22:
			unk, err := it.Decode2()
			Checkerr(err)
			newstate, err := it.Decode1()
			Checkerr(err)
			p.Encode1(command)
			p.Encode1(newstate)
			p.Encode2(unk)
		}
	}
	Log(p.String())
	return []byte(p[4:])
}

//玩家移动
func movePlayer(cid uint32, moves []byte) []byte {
	p := NewPacket()
	p.Encode2(SH.MOVE_PLAYER)
	p.Encode4(cid)
	p.Encode4(0)

	p.Append(moves)

	return []byte(p)
}

func NewPacket() maplelib.Packet {
	p := maplelib.NewPacket()
	p.Encode4(0)
	return p
}

func Log(v ...interface{}) {
	log.Println(v...)
}

func CheckError(tag string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, tag+": %s", err.Error())
		os.Exit(1)
	}
}

func Checkerr(err error) {
	if err != nil {
		panic(err)
	}
}
