package main

import (
	"database/sql"
	"encoding/hex"
	"log"
	"net"

	"github.com/Francesco149/maplelib"
	//"github.com/go-sql-driver/mysql"
)

type MapleClient struct {
	Channels    int
	Send        [4]byte
	Receive     [4]byte
	Conn        net.Conn
	AccountName string
	LoginTry    int8
	LoginStatus int8
	AccId       int
	RC          maplelib.Crypt //Receive Crypt
	SC          maplelib.Crypt //Send Crypt
	Gender      int8
	Character   MapleCharacter
}

func NewMapleClient(send [4]byte, receive [4]byte, conn net.Conn) MapleClient {
	var c MapleClient
	c.Channels = -1
	c.Send = send
	c.Receive = receive
	c.Conn = conn
	return c
}

func (this *MapleClient) login(name string, pwd string) int32 {
	row := db.QueryRow("SELECT id,password,loginstatus,gender FROM accounts WHERE username = ? ", name)
	var (
		id          int
		password    string
		loginstatus int8
		gender      int8
	)
	err := row.Scan(&id, &password, &loginstatus, &gender)
	if err != nil {
		if err != sql.ErrNoRows {
			Log(err)
		}
		return 5 //账号不存在
	}
	if password != pwd {
		this.LoginTry += 1
		if this.LoginTry >= 5 {
			return -1 //超过尝试次数
		}
		return 4 // 密码错误
	}

	if gender == -1 {
		return -2 // 需要选择性别
	}
	if loginstatus != 0 {
		return 7 // 已经登录
	}
	this.AccId = id
	this.LoginStatus = 1
	this.Gender = gender
	return 0 //登录成功

}
func (this *MapleClient) setGender(gender int, username string) {
	Log("UPDATE accounts SET gender = ? WHERE username = ?", gender, this.AccountName == username)
	_, err := db.Exec("UPDATE accounts SET gender = ? WHERE username = ?", gender, username)
	Checkerr(err)
}
func (this *MapleClient) Write(data []byte) {
	log.Printf("发送消息:\n" + hex.Dump(data[4:]))
	this.SC.Encrypt(data)
	this.SC.Shuffle()
	this.Conn.Write(data)

}

func (this *MapleClient) Read() ([]byte, error) {
	buffer := make([]byte, 2048)
	n, err := this.Conn.Read(buffer)
	if err != nil {
		log.Println(this.Conn.RemoteAddr().String(), " 连接断开: ", err)
		return nil, err
	}

	data := buffer[4:n]
	this.RC.Decrypt(data)
	this.RC.Shuffle()
	return data, nil
}
