package main

import (
	"github.com/lxmgo/config"
)

type SendHeader struct {
	//Login 登录
	LOGIN_STATUS       uint16
	LICENSE_RESULT     uint16
	CHOOSE_GENDER      uint16
	GENDER_SET         uint16
	SERVERSTATUS       uint16
	SERVERLIST         uint16
	CHARLIST           uint16
	SERVER_IP          uint16
	CHAR_NAME_RESPONSE uint16
	ADD_NEW_CHAR_ENTRY uint16
	CHANNEL_SELECTED   uint16
	//Channel 频道
	WARP_TO_MAP   uint16
	UPDATE_STATUS uint16
	SPAWN_PLAYER  uint16
	MOVE_PLAYER   uint16
}

type ReceiveHeader struct {
	//Login 登录
	LOGIN_PASSWORD       uint16
	SERVERLIST_REQUEST   uint16
	LICENSE_REQUEST      uint16
	SET_GENDER           uint16
	SERVERSTATUS_REQUEST uint16
	CHARLIST_REQUEST     uint16
	CHAR_SELECT          uint16
	CHECK_CHAR_NAME      uint16
	CREATE_CHAR          uint16
	//Channel 频道
	PLAYER_LOGGEDIN uint16
	CHANGE_MAP      uint16
	MOVE_PLAYER     uint16
}

var SH SendHeader
var RH ReceiveHeader

func loadPacketSH() {
	s, err := config.NewConfig("send.ini")
	if err != nil {
		Log("载入发送包头配置错误:", err)
		return
	}
	SH.LOGIN_STATUS = s.Uint16("LOGIN_STATUS")
	SH.LICENSE_RESULT = s.Uint16("LICENSE_RESULT")
	SH.CHOOSE_GENDER = s.Uint16("CHOOSE_GENDER")
	SH.GENDER_SET = s.Uint16("GENDER_SET")
	SH.SERVERSTATUS = s.Uint16("SERVERSTATUS")
	SH.SERVERLIST = s.Uint16("SERVERLIST")
	SH.CHARLIST = s.Uint16("CHARLIST")
	SH.SERVER_IP = s.Uint16("SERVER_IP")
	SH.CHAR_NAME_RESPONSE = s.Uint16("CHAR_NAME_RESPONSE")
	SH.ADD_NEW_CHAR_ENTRY = s.Uint16("ADD_NEW_CHAR_ENTRY")
	SH.CHANNEL_SELECTED = s.Uint16("CHANNEL_SELECTED")
	SH.WARP_TO_MAP = s.Uint16("WARP_TO_MAP")
	SH.UPDATE_STATUS = s.Uint16("UPDATE_STATUS")
	SH.SPAWN_PLAYER = s.Uint16("SPAWN_PLAYER")
	SH.MOVE_PLAYER = s.Uint16("MOVE_PLAYER")
}

func loadPacketRH() {
	s, err := config.NewConfig("receive.ini")
	if err != nil {
		Log("载入返回包头配置错误:", err)
		return
	}
	RH.LOGIN_PASSWORD = s.Uint16("LOGIN_PASSWORD")
	RH.SERVERLIST_REQUEST = s.Uint16("SERVERLIST_REQUEST")
	RH.LICENSE_REQUEST = s.Uint16("LICENSE_REQUEST")
	RH.SET_GENDER = s.Uint16("SET_GENDER")
	RH.SERVERSTATUS_REQUEST = s.Uint16("SERVERSTATUS_REQUEST")
	RH.CHARLIST_REQUEST = s.Uint16("CHARLIST_REQUEST")
	RH.CHAR_SELECT = s.Uint16("CHAR_SELECT")
	RH.CHECK_CHAR_NAME = s.Uint16("CHECK_CHAR_NAME")
	RH.CREATE_CHAR = s.Uint16("CREATE_CHAR")
	RH.PLAYER_LOGGEDIN = s.Uint16("PLAYER_LOGGEDIN")
	RH.CHANGE_MAP = s.Uint16("CHANGE_MAP")
	RH.MOVE_PLAYER = s.Uint16("MOVE_PLAYER")
}
