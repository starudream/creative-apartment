# Database

- bbolt [GoDoc](https://pkg.go.dev/go.etcd.io/bbolt) | [GitHub](https://github.com/etcd-io/bbolt)

## Buckets

### `config`

预置

### `customer`

```
`${phone}`
```

### `${phone}_house_data_${YYYY}_${type}`

`type` 取值为 `0`, `1`, `2`，分别为 `电`, `水`, `气`

```
`${MMDD}`
```

```json
 {
    "surplus": 0,
    "surplusAmount": 0,
    "unitPrice": 0,
    "lastReadTime": "2022-05-12T00:00:00+08:00"
}
```

### `${phone}_house_stats_${YYYY}_${type}`

```
`${MMDD}`
```

值为一个数组，`0` 为 `量`，`1` 为 `费`

```json
[0, 0]
```
