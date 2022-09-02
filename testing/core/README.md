## some key pair
```
Validate1
Seckey: 20504139dcba1f398acde20333673a7dfeb0a780416b9a2738bcf4a5f76e0c26
Pubkey: a892612cf51956116275b2a17a65802c849c570b5c8327bef7985f703b96c2a812f55c53b0192c4a5081307676cecbe9
Address: eae5739b41db1e42e7a0301da93c4697663f3469
```
```
Validate2
Seckey: 1e7459f6d9e645c3fe562b21ec071dcc216b74915bc5f6aa044a547138349694
Pubkey: a5550e6beea6224c113d5a013e0c18b09c403be4b6bc7a78f2a975e3afcf484e9ac4811a190666a4bd4d3a798e5dcc3e
Address: 2d0950b45dbc84f40c6fa00edbe22ce8836963ca
```
```
Validate3
Seckey: 6b0701cb53e84d9e695cff69fb3e0d87c7be3f4d75134becf7f7d690a04d6b4b
Pubkey: 83ba90eb84a0704560de9583b6cff7de94ec017b0138c1dac0e325fe5024b801cdcd9be6d43631afdb24d23ccdf859ca
Address: d9f1bcc35777192012db45cddebd682e5312da0c
```
```
node1
Seckey: 52c5089108db8b63252baf64876393a8d36b39d9f6ad41f848dd0ff9f5d9db2f
Pubkey: a5d2fa3ceb485ae2426e481b442801ad08864f476b9609ee4c382afa15e91904e2a8c830c04d9f73e8fc87530daf3f5a
Address: 1dca49f3183d83eb6d23475f97855daa887b9828
```
```
node2
Seckey: 73bc0c94423d8bbdac8cd2e04852090114d26601bbe85f0f7677081ad256897f
Pubkey: a90070ff2a19e3140ae51ab9ee45c47d1a1789b417f5d4812d3c8049571a9cc8e24c6ae701845723e52c4472ce14900e
Address: 38990c0ec3cc8f32fb893dc59a676c25eca01404
```
```
node3
Seckey: 3fc5dc5c1e885b1558b97d2b93b6e8b16cec41e2dd86b039ebf7d08d558f1975
Pubkey: b4b7f7f0dfd87263755411571fa6d3f5434c07196861b3aecc9850bfd624f70cccf94bc79fab3741ab1789807cde84b2
Address: d142c6c25e9d2cb439cbb40cd7dc468f3858be65
```
```
miner1
Seckey: 73729e721d718cce82109a56ebd445d22224d090bcf4d2a833b114d31a726937
Pubkey: 8a4cd49435a377d6c1e7692dbcb7fc4e8d4d5eeade474be32ae5dcd6c0fb4987d143ae81d84c9b1d4ea67fa0458c9f32
Address: 25ce79a3148482130059718b9275401917bb249b
```
```
miner2
Seckey: 06b7fc84d79aa9d341c15ca823fdd5b1191b75128fe7766b84ad4d3d100cd7e7
Pubkey: b13d542ba748134ae3ab8270d0234566c6963b009a4a0c9119b7bbcae654b050d203ea4f90f4f0e33e42f3b587df7765
Address: 848b518c2c0c5963d41488ca56f7612c7d85b244
```
```
miner3
Seckey: 62e69260d7731e0178d56a2f996e32ad9d32cfc1462f27c76d10337129b21575
Pubkey: a3f988e87cbffab8112ba8fe316ff606a4184bf69ccf213f8a456016fcf91b42fed1dfc4d996c1cfd8ca3d18cb4700de
Address: 133ebed2b953db9ae64abf2989fce709bebcfccd
```
```
explorer
Seckey: 62d5d86cc43ade13f5f9dba8552de5ccfc4d6b72e5b73cc619f8512618c1b693
Pubkey: 8d645c68921b6307459d4c7b18a86306d945df529a1e1bfc5739e95ac92cabe9664e7eb0f78ac285333a2cba82887dfd
Address: edd76cbc5faea7bdd7bf9e7e2565b44deeb2ac6f
```
```


build proto
protoc --go_out=. ./proto/*.proto

## TEST
```
export CGO_CFLAGS="-O -D__BLST_PORTABLE__"
go test -v -coverprofile cover.out -coverpkg ./... ./...  
go tool cover -html=cover.out  

go test -cpu 1,2,4,8 -benchmem -v ./... -bench=.
```

## BUILD
```
export CGO_CFLAGS="-O -D__BLST_PORTABLE__"  
protoc --go_out=. ./proto/*.proto
go build .

```