package vuid

var (
	uuid UUid
)

type UUid interface {
	Generate() (id int64)
	SetWorkId(id int64)
}

func newGenerator() UUid{
	// 默认生成 UUid 的方法使用雪花算法
	return NewSnowFlake()
}

// 获取 int64 Uuid
func GenUUid() (id int64) {
	// 随机生成一个 workID 生成一个 uuid
	return genWithId(1)
}

func GenUUidWithId(workId int64){
	genWithId(workId)
}

// 指定一个自己的 work id 生成uuid, work 最大值为 1023
func genWithId(workId int64) (id int64) {
	if uuid == nil {
		uuid = newGenerator()
	}
	uuid.SetWorkId(workId)
	id = uuid.Generate()
	return
}
