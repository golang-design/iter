# iter

Iterator utils for Go.

```go
import "golang.design/x/iter"
```

## Batched Slice Iterator

```go
it := iter.NewBatchFromSlice[T](s, 42)
for batch, ok := it.Next(); ok; batch, ok = it.Next() {
    // do something with batch
}
```

## Batched `gorm` Iterator

```go
it := NewGormIter[T](tx, batchSize)
for batch, ok := it.Next(); ok; batch, ok = it.Next() {
	// Process the batch.
	...
	// Stop the iteration if necessary.
	if ... {
		it.Stop()
		break
	}
}
if err := it.Err(); err != nil {
	// Handle the error.
}
```

## License

MIT | &copy; 2023 The golang.design Initiative Authors, written by [Changkun Ou](https://changkun.de).