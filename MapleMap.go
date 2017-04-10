package main

type MapleMap struct {
	mapid      uint32
	characters map[uint32]*MapleCharacter
}

func (this *MapleMap) WritetoMapOther(player MapleCharacter, data []byte) {
	chars := this.characters
	for char := range chars {
		if chars[char].id == player.id {
			continue
		}
		chars[char].client.Write(data)
	}
}

func (this *MapleMap) WritetoMap(data []byte) {
	chars := this.characters
	for char := range chars {
		chars[char].client.Write(data)
	}
}
