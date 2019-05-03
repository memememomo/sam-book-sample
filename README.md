# サーバレスで構築するWebAPI

「サーバレスで構築するWebAPI」のサンプルプログラムです。  
https://techbookfest.org/event/tbf06/circle/40470002

## 正誤表

### structのバッククォート


```
P30

(誤)
// Request の構造体定義
type PostHelloRequest struct {
        Name string ‘json:"name"‘
}
// Response の構造体定義
type HelloMessageResponse struct {
        Message string ‘json:"message"‘
}

P30

(正)
// Request の構造体定義
type PostHelloRequest struct {
        Name string `json:"name"`
}
// Response の構造体定義
type HelloMessageResponse struct {
        Message string `json:"message"`
}
```

```
P43 

(誤)
// Entity としてのデータ構造
type Product struct {
        BaseModel
        Name        string    ‘dynamo:"Name"‘
        Price       int       ‘dynamo:"Price"‘
        ReleaseDate time.Time ‘dynamo:"ReleaseDate"‘
}

(正)
// Entity としてのデータ構造
type Product struct {
        BaseModel
        Name        string    `dynamo:"Name"`
        Price       int       `dynamo:"Price"`
        ReleaseDate time.Time `dynamo:"ReleaseDate"`
}
```

```
P44

(誤)
type Product struct {
        BaseModel
        Name        string ‘dynamo:"Name"‘
        Price       int    ‘dynamo:"Price"‘
        ReleaseDate time.Time ‘dynamo:"ReleaseDate"‘
}

(正)
type Product struct {
        BaseModel
        Name        string `dynamo:"Name"`
        Price       int    `dynamo:"Price"`
        ReleaseDate time.Time `dynamo:"ReleaseDate"`
}
```


### 必要なライブラリのインポート

```
P44 リスト 3.15: Model を実装 (models/product.go)

(誤) 
import (
        "sam-book-sample/db"
        "time"
)

(正)
import (
	"sam-book-sample/db"
	"github.com/memememomo/nomof"  //追加

	"time"

	"github.com/pkg/errors"        //追加
	"github.com/guregu/dynamo"     //追加
)
```

```
P47  リスト 3.17: Create メソッドのテスト (models/product_test.go)

(追記)

import (
	"sam-book-sample/db"

	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/guregu/dynamo"
)
```
