package connection
// ServerConn is..
type ServerConn struct {
	Name string
	AddressString string
}
// JSONServerConn is..
type JSONServerConn struct {
	Name string `json:"name"`
	AddressString string `json:"addressString"`
}
