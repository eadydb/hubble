package accident

// Accident 事故回放
type Accident struct {
	Name       string             `json:"name"`        // 事故名称
	Level      AccidentLevel      `json:"level"`       // 事故级别
	Point      AccidentCheckpoint `json:"point"`       // 事故类型
	CheckPoint []CheckPoint       `json:"check_point"` // 事故点
}

// AccidentLevel accident level
type AccidentLevel int

// accident level define
const (
	P0 AccidentLevel = 0 // p0 accident
	P1 AccidentLevel = 1 // p1 accident
	P2 AccidentLevel = 2 // p2 accident
	P3 AccidentLevel = 3 // p3 accident
)

// AccidentCheckpoint  accident point
type AccidentCheckpoint string

const (
	DataBase AccidentCheckpoint = "database"
	Redis    AccidentCheckpoint = "redis"
	kafka    AccidentCheckpoint = "kakfa"
	NetWork  AccidentCheckpoint = "network"
	RabbitMq AccidentCheckpoint = "rabbitmq"
	Cpu      AccidentCheckpoint = "cpu"
	Memory   AccidentCheckpoint = "memory"
)
