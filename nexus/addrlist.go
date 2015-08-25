package nexus

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/coredef"
	"log"
)

type RegionData struct {
	*coredef.Region
	Session cellnet.CellID
}

// TODO 线程安全

var (

	// regionid的映射
	id2regionMap map[int32]*RegionData = make(map[int32]*RegionData)

	// Region间的连接点(nodeid)的映射
	ses2regionMap map[cellnet.CellID]*RegionData = make(map[cellnet.CellID]*RegionData)

	// 事件通知
	Event *cellnet.EventDispatcher = cellnet.NewEventDispatcher()
)

func AddRegion(ses cellnet.CellID, profile *coredef.Region) {

	rd := &RegionData{
		Region:  profile,
		Session: ses,
	}

	id2regionMap[profile.GetID()] = rd
	ses2regionMap[ses] = rd

	Event.Invoke("OnAddRegion", rd)

	log.Printf("[nexus] add region: %d@%s node: %v", profile.GetID(), profile.GetAddress(), ses.String())
}

func GetRegion(regionID int32) *RegionData {
	if v, ok := id2regionMap[regionID]; ok {
		return v
	}

	return nil
}

func getRegionBySession(nid cellnet.CellID) *RegionData {
	if v, ok := ses2regionMap[nid]; ok {
		return v
	}

	return nil
}

func RemoveRegion(nid cellnet.CellID) (*RegionData, bool) {

	rd := getRegionBySession(nid)
	if rd == nil {
		return nil, false
	}

	Event.Invoke("OnRemoveRegion", rd)

	delete(id2regionMap, rd.GetID())
	delete(ses2regionMap, nid)

	log.Println("[nexus] remove region: %d@%s", rd.GetID(), rd.GetAddress())

	return rd, true

}

func IterateRegion(callback func(*RegionData)) {

	for _, v := range id2regionMap {

		callback(v)
	}

}
