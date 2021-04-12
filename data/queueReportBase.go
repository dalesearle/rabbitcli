package data

type QueueReportBase struct {
	connections []*Connection
	Queue *Queue
}

func NewQueueReportBase(connection *Connection) *QueueReportBase {
	return &QueueReportBase{
		connections: []*Connection{connection},
	}
}

func (qc *QueueReportBase) AddConnection(connection *Connection) {
	if qc.connections == nil {
		qc.connections = make([]*Connection, 0)
	}
	qc.connections = append(qc.connections, connection)
}

func (qc *QueueReportBase) Connections() []*Connection {
	return qc.connections
}

func (qc *QueueReportBase) IsMissingConnection() bool {
	return qc.connections == nil
}

func (qc *QueueReportBase) IsMissingQueue() bool {
	return qc.Queue == nil
}
