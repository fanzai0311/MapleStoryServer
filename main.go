package main

import (
	"encoding/hex"
	"log"
	"math/rand"
	"net"
	. "strconv"
	"time"

	"github.com/Francesco149/maplelib"
	"github.com/axgle/mahonia"
	"github.com/lxmgo/config"
)

var (
	MapleVersion = uint16(79)
	ServerName   string
	ServerHost   string
	ServerIP     []byte
	enc          mahonia.Encoder
	dec          mahonia.Decoder
	MapleMaps    map[uint32]*MapleMap
)

func main() {
	enc = mahonia.NewEncoder("gbk")
	dec = mahonia.NewDecoder("gbk")
	//载入配置
	c, err := config.NewConfig("server.ini")
	if err != nil {
		Log("载入服务端配置错误:", err)
		return
	}
	ServerName = c.String("servername")
	ServerHost = c.String("serverhost")
	ServerIP = net.ParseIP(ServerHost)
	ServerIP = ServerIP[len(ServerIP)-4:]
	dbhost := c.String("dbhost")
	dbuser := c.String("dbuser")
	dbpw := c.String("dbpw")
	dbname := c.String("dbname")
	//连接数据库
	connsql(dbhost, dbuser, dbpw, dbname)
	defer db.Close()
	Log(ServerName)
	//载入包头
	loadPacketSH()
	loadPacketRH()

	//启动登录端
	login, err := net.Listen("tcp", ServerHost+":8487")
	CheckError("启动登录端", err)
	defer login.Close()
	Log("启动登录端 成功")
	go LoginSessionOpened(login)

	//启动频道
	channel, err := net.Listen("tcp", ServerHost+":9001")
	CheckError("启动频道端", err)
	defer channel.Close()
	MapleMaps = make(map[uint32]*MapleMap)
	var mm MapleMap
	mm.mapid = 910000000
	mm.characters = make(map[uint32]*MapleCharacter)
	MapleMaps[910000000] = &mm
	Log("启动频道端 成功")
	go ChannelSessionOpened(channel)
	for {
		time.Sleep(time.Second)
	}
}

func LoginSessionOpened(login net.Listener) {
	for {
		c, err := login.Accept()
		if err != nil {
			log.Println("登录端 ", err)
			continue
		}
		Log(c.RemoteAddr().String(), " 连接到登录端")
		ivSend := [4]byte{byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255))}
		ivRecv := [4]byte{byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255))}

		clients := NewMapleClient(ivSend, ivRecv, c)

		c.Write(GetHello(MapleVersion, ivSend, ivRecv))
		go LoginHandle(clients)
	}
}

func LoginHandle(c MapleClient) {
	defer c.Conn.Close()
	c.RC = maplelib.NewCrypt(c.Receive, MapleVersion)
	c.SC = maplelib.NewCrypt(c.Send, MapleVersion)
	for {
		data, err := c.Read()
		if err != nil {
			return

		}
		Log(c.Conn.RemoteAddr().String(), "收到消息: \n"+hex.Dump(data))
		go LoginHandlePro(data, &c)
	}
}

func LoginHandlePro(data []byte, c *MapleClient) {
	var res []byte
	p := maplelib.NewPacket()
	p.Append(data)
	it := p.Begin()
	paketId, err := it.Decode2()
	Checkerr(err)
	switch paketId {
	//LOGIN*
	case RH.LOGIN_PASSWORD: //LOGIN_PASSWORD 登录
		account, err := it.DecodeString()
		Checkerr(err)
		password, err := it.DecodeString()
		Checkerr(err)
		Log("账号:", account, " 密码:", password)
		c.AccountName = account
		loginok := c.login(account, password)
		switch loginok {
		case 0:
			res = getLoginSuccess(c)
			c.Write(res)
			res = getServerList()
			c.Write(res)
			res = getEndOfServerList()
			c.Write(res)
		case -1: // 超过尝试次数
			res = GetLoginFailed(4)
			c.Write(res)
			c.Conn.Close()
		case -2: //需要选择性别
			res = GenderNeeded(c)
			c.Write(res)
		default:
			res = GetLoginFailed(loginok)
			c.Write(res)

		}
	case RH.SERVERLIST_REQUEST: //SERVERLIST_REQUEST 服务器列表请求
		res = getServerList()
		c.Write(res)
		res = getEndOfServerList()
		c.Write(res)
	case RH.SET_GENDER: //SET_GENDER 设置性别
		gender, err := it.Decode1()
		Checkerr(err)
		username, err := it.DecodeString()
		Checkerr(err)
		if gender == 0 {
			c.setGender(0, username)
		} else {
			c.setGender(1, username)
		}
		res = GenderChanged(c)
		c.Write(res)
		res = GetLoginFailed(22)
		c.Write(res)
	case RH.LICENSE_REQUEST: //LicenseRequest 许可协议
		yesno, err := it.Decode1()
		Checkerr(err)
		if yesno == 1 {
			res = LicenseResult()
			c.Write(res)
		} else {
			c.Conn.Close()
		}
	case RH.SERVERSTATUS_REQUEST: //ServerStatus_Request 服务器状态请求
		res = getServerStatus()
		c.Write(res)
	case RH.CHARLIST_REQUEST: //CharList_Request 角色列表请求
		res = getCharList(c)
		c.Write(res)
	case RH.CHECK_CHAR_NAME: //Check_Char_Name 检查角色名
		charname, err := it.DecodeString()
		charname = dec.ConvertString(charname)
		Checkerr(err)
		Log("检查角色名: " + charname)
		res = getCheckName(charname)
		c.Write(res)
	case RH.CREATE_CHAR: //Create_Char 创建角色
		charname, err := it.DecodeString()
		Checkerr(err)
		charname = dec.ConvertString(charname)
		job, err := it.Decode4()
		Checkerr(err)
		face, err := it.Decode4()
		Checkerr(err)
		hair, err := it.Decode4()
		Checkerr(err)
		//			top, err := it.Decode4()
		//			Checkerr(err)
		//			bottom, err := it.Decode4()
		//			Checkerr(err)
		//			shoes, err := it.Decode4()
		//			Checkerr(err)
		//			weapon, err := it.Decode4()
		//			Checkerr(err)
		gender := c.Gender
		var char MapleCharacter
		char.face = face
		char.hair = hair
		char.gender = uint8(gender)
		if job == 2 {
			char.str = 11
			char.dex = 6
			char.int_ = 4
			char.luk = 4
			char.remainingAP = 0
		} else {
			char.str = 4
			char.dex = 4
			char.int_ = 4
			char.luk = 4
			char.remainingAP = 0
		}
		char.name = charname
		char.skincolor = 0
		switch job {
		case 0:
			char.job = 1000
		case 1:
			char.job = 0
		case 2:
			char.job = 2000
		}
		Log(char)
		res = addNewCharEntry(&char)
		c.Write(res)
	case RH.CHAR_SELECT:
		charid, err := it.Decode4()
		Checkerr(err)
		res = getServerIP(charid)
		c.Write(res)
	}

}

func ChannelSessionOpened(channel net.Listener) {
	for {
		c, err := channel.Accept()
		if err != nil {
			log.Println("频道端 ", err)
			continue
		}
		Log(c.RemoteAddr().String(), " 连接到频道端")
		ivSend := [4]byte{byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255))}
		ivRecv := [4]byte{byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255))}

		clients := NewMapleClient(ivSend, ivRecv, c)

		c.Write(GetHello(MapleVersion, ivSend, ivRecv))
		go ChannelHandle(clients)
	}
}

func ChannelHandle(c MapleClient) {
	defer c.Conn.Close()
	c.RC = maplelib.NewCrypt(c.Receive, MapleVersion)
	c.SC = maplelib.NewCrypt(c.Send, MapleVersion)
	for {
		data, err := c.Read()
		if err != nil {
			return

		}
		Log(c.Conn.RemoteAddr().String(), "收到消息: \n"+hex.Dump(data))
		go ChannelHandlePro(data, &c)
	}
}

func ChannelHandlePro(data []byte, c *MapleClient) {

	var res []byte

	p := maplelib.NewPacket()
	p.Append(data)
	it := p.Begin()
	paketId, err := it.Decode2()
	Checkerr(err)
	switch paketId {
	case RH.PLAYER_LOGGEDIN:
		charid, err := it.Decode4()
		Checkerr(err)
		charid = charid
		var char MapleCharacter
		id := rand.Intn(10) + 1
		Checkerr(err)
		char.name = "管理员00" + Itoa(id)
		char.id = uint32(id)
		char.level = 250
		char.job = 900
		char.str = 4
		char.dex = 4
		char.int_ = 4
		char.luk = 4
		char.hp = 50
		char.maxhp = 50
		char.mp = 50
		char.maxmp = 50
		char.remainingAP = 9
		char.mapid = 910000000
		char.spawnPoint = 0
		char.face = 20100
		char.hair = 30000
		char.meso = 0
		char.position.x = 100
		char.position.y = 50
		c.Character = char
		c.Character.client = c
		res = getCharInfo(&c.Character)
		c.Write(res)
		Log("玩家[", c.Character.name, "]登录")
		res = weirdStatusUpdate()
		c.Write(res)

		MapleMaps[c.Character.mapid].characters[c.Character.id] = &c.Character
		c.Character.map_ = &MapleMaps
		MapleMaps[c.Character.mapid].WritetoMapOther(c.Character, spawnPlayerMapobject(&c.Character))
		for _, sc := range MapleMaps[c.Character.mapid].characters {
			if sc.id == c.Character.id {
				continue
			}
			c.Write(spawnPlayerMapobject(sc))
		}
	case RH.MOVE_PLAYER:
		it.Skip(33)
		moves := parseMovement(&it)
		MapleMaps[c.Character.mapid].WritetoMapOther(c.Character, movePlayer(c.Character.id, moves))
	default:
		return
	}

}
