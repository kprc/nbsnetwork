package flowkey



type FlowKey struct {
	stationId string
	serialNo uint64
}

func NewFlowKey(s string,sn uint64) *FlowKey {
	return &FlowKey{stationId:s,serialNo:sn}
}

func (fk *FlowKey)GetStationId() string  {
	return fk.stationId
}

func (fk *FlowKey)GetSerialNo() uint64  {
	return fk.serialNo
}

func (fk *FlowKey)SetStationId(s string)  {
	fk.stationId = s
}

func (fk *FlowKey)SetSerialNo(sn uint64)  {
	fk.serialNo = sn
}

