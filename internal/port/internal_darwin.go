//go:build darwin
package port
func NewInspector() PortInspector {
    return DarwinPortInspector{}
}