package port

type MockPortInspector struct{}

func (m MockPortInspector) Find(port int) ([]PortInfo, error) {
	return []PortInfo{
		{
			Port: port,
			Protocol: "TCP",
			Address: "127.0.0.1:8080",
			PId: 4218,
			Process: "dotnet",
			User: "Aryan",
		},
	}, nil
}