package main

type LifeMovement interface {
	serialize() []byte
}

type AbsoluteLifeMovement struct {
	command  byte
	x        uint16
	y        uint16
	newstate byte
	duration uint16
	unk      uint16
	ppsx     uint16 //PixelsPerSecond  像素/秒
	ppsy     uint16
}
type RelativeLifeMovement struct {
	command  byte
	x        uint16
	y        uint16
	newstate byte
	duration uint16
}

type TeleportMovement struct {
	command  byte
	x        uint16
	y        uint16
	ppsx     uint16
	ppsy     uint16
	newstate byte
}

type ChairMovement struct {
	command  byte
	x        uint16
	y        uint16
	unk      uint16
	newstate byte
	duration uint16
}

type JumpDownMovement struct {
	command  byte
	x        uint16
	y        uint16
	ppsx     uint16
	ppsy     uint16
	unk      uint16
	fh       uint16
	newstate byte
	duration uint16
}

func (this AbsoluteLifeMovement) serialize() []byte {
	p := NewPacket()
	p.Encode1(this.command)
	p.Encode2(this.x)
	p.Encode2(this.y)
	p.Encode2(this.ppsx)
	p.Encode2(this.ppsy)
	p.Encode2(this.unk)
	p.Encode1(this.newstate)
	p.Encode2(this.duration)
	return []byte(p[4:])
}

func (this RelativeLifeMovement) serialize() []byte {
	p := NewPacket()
	p.Encode1(this.command)
	p.Encode2(this.x)
	p.Encode2(this.y)
	p.Encode1(this.newstate)
	p.Encode2(this.duration)
	return []byte(p[4:])
}
func (this TeleportMovement) serialize() []byte {
	p := NewPacket()
	p.Encode1(this.command)
	p.Encode2(this.x)
	p.Encode2(this.y)
	p.Encode1(this.newstate)
	return []byte(p[4:])
}
func (this ChairMovement) serialize() []byte {
	p := NewPacket()
	p.Encode1(this.command)
	p.Encode2(this.x)
	p.Encode2(this.y)
	p.Encode1(this.newstate)
	p.Encode2(this.duration)
	return []byte(p[4:])
}
func (this JumpDownMovement) serialize() []byte {
	p := NewPacket()
	p.Encode1(this.command)
	p.Encode2(this.x)
	p.Encode2(this.y)
	p.Encode1(this.newstate)
	p.Encode2(this.duration)
	return []byte(p[4:])
}
