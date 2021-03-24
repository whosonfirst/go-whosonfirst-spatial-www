# go-whosonfirst-reader

Common methods for reading Who's On First documents.

## Examples

_Note that error handling has been removed for the sake of brevity._

### LoadReadCloserFromID

```
import (
	"context"
	"github.com/whosonfirst/go-reader"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	"io"
	"os"
)

func main() {

	ctx := context.Backround()
	wof_id := int64(101736545)

	r_uri := "local:///usr/local/data/whosonfirst-data-admin-ca/data"
	r, _ := reader.NewReader(ctx, r_uri)

	fh, _ := wof_reader.LoadReadCloserFromID(ctx, r, wof_id)
	io.Copy(os.Stdout, fh)
}
```

### LoadBytesFromID

```
import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-reader"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
)

func main() {

	ctx := context.Backround()
	wof_id := int64(101736545)

	r_uri := "local:///usr/local/data/whosonfirst-data-admin-ca/data"
	r, _ := reader.NewReader(ctx, r_uri)

	body, _ := wof_reader.LoadReadCloserFromID(ctx, r, wof_id)
	fmt.Printf("%d bytes\n", len(body))
}
```

### LoadFeatureFromID

```
import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-reader"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"	
)

func main() {

	ctx := context.Backround()
	wof_id := int64(101736545)

	r_uri := "local:///usr/local/data/whosonfirst-data-admin-ca/data"
	r, _ := reader.NewReader(ctx, r_uri)

	f, _ := wof_reader.LoadFeatureFromID(ctx, r, wof_id)
	fmt.Println(f.Name())
}
```

## See also

* https://github.com/whosonfirst/go-reader