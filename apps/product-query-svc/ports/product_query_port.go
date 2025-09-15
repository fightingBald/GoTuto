package ports

import ip "github.com/fightingBald/GoTuto/internal/ports"

// 入站端口：应用层使用的查询接口（alias 到 internal 的定义以避免重复）
type ProductQueryPort = ip.ProductService
