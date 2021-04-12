package data

type ConnectionReportBase struct {
	Name  string
	Connections []*Connection
}

func NewConnectionReportBase(connection *Connection) *ConnectionReportBase {
	cons := make([]*Connection,0)
	cons = append(cons, connection)
	return &ConnectionReportBase{
		Name:  connection.ProvidedName,
		Connections: cons,
	}
}

func (cc *ConnectionReportBase) AddConnection(connection *Connection) {
	cc.Connections = append(cc.Connections, connection)
}

func (cc *ConnectionReportBase) ConnectionCount() int {
	return len(cc.Connections)
}

func (cc *ConnectionReportBase) String() string {
	return cc.Name
}
