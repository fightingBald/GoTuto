package domain

import dd "github.com/fightingBald/GoTuto/internal/domain"

// 该包是新目录的领域层外观（facade）。为了最小化重复实现，我们将原有 internal/domain 类型作别名导出。

type Product = dd.Product

// 额外的领域相关类型可以在此扩展。
