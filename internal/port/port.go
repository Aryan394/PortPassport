package port

type PortInfo struct{
	Port int
	Protocol string
	Address string
	PId int
	Process string
	User string
}

type PortInspector interface {
	Find(port int) ([]PortInfo, error)
}