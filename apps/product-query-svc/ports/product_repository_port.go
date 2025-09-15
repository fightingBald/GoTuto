package ports

import ip "github.com/fightingBald/GoTuto/internal/ports"

// Outbound port (alias to existing internal port to avoid duplication)
type ProductRepository = ip.ProductRepo
