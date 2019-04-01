# cashAccountAPId
cashAccountAPId is a web-based API to access proof-of-inclusion for a CashAccount.
cashAccountDBd is required for cashAccountAPId to run.

More about CashAccount: https://gitlab.com/cash-accounts/specification/blob/master/SPECIFICATION.md
More about cashAccountDBd: https://github.com/abck/cashAccountDBd

### 1. Installation

Install cashAccountDBd ( see https://github.com/abck/cashAccountDBd/blob/master/README.md )

Install cashAccountAPId via go get
```
go get github.com/abck/cashAccountAPId
```

Run cashAccountAPId once to create the default configuration.
```
~/go/bin/cashAccountAPId
```

Edit the configuration, located at `~/.cashAccount/cashAccountAPId.conf`
Sample configuration is:
```
[cashAccountAPId]
rpchost=127.0.0.1:8334
rpcendpoint=ws
rpcuser=USERNAME
rpcpass=PASSWORD
webserverbindaddr=0.0.0.0:8585

```

Start the API-Server:
```
~/go/bin/cashAccountAPId
```

cashAccountAPId is now running.
You should now see an output like this:
```
2019/01/28 20:06:30 Opening database...
2019/01/28 20:06:30 Connecting to bchd...
2019/01/28 20:06:30 Starting webserver...
```

### 2. Updating
To update to a newer version of cashAccountAPId run:
```
go get -u github.com/abck/cashAccountAPId
```

### 2. Updating
To update to a newer version of cashAccountAPId run:
```
go get -u github.com/abck/cashAccountAPId
```

### 3. Web-Endpoints
#### `/lookup/accountNumber/accountName/`
Example: `/lookup/100/jonathan/` will return:
``` 
{
  "identifier": "jonathan#100",
  "block": 563720,
  "results": [
    {
      "transaction": "01000000017cc04d29109cb43a0bfade3993b5840b6f68f22e09d4806f26fe9e7d772fc72f010000006a47304402207b0da3150bf9a44a8fae7333f4d5b03ba1297dd21c8641880efea9eaa9e1e89d022028d2a8d840771d4a87c84f0b217d02f17c3f57832fe0c37ac27b5afc27004fbe41210355f64f0ed04944eb477b33dcb46bb45453b8988bba1862698abe7343c6f0e2c6ffffffff020000000000000000256a0401010101084a6f6e617468616e1501ebdeb6430f3d16a9c6758d6c0d7a400c8e6bbee4c00c1600000000001976a914efd03e75f2aedb19261b39a6c8361c7bccd9f4f088ac00000000",
      "inclusion_proof": "0000c0204895ef83b69cd64a7901fee54858461bfc0c09cb74a6b0000000000000000000cdd0c434d6ee0aee331016d1f9d2bef0bbcd26e003b784a39f392f2f2a50bfda6f532e5c6ea304184467d54bd300000009c33356693bc1c8928b96621d832c510e6e239ae52f0ee27ab28a44643313e0c65eb8a4982058bb5e86cf84653c797d51741a90ca6c1a2518c1ba871402cb76e31faf65208ce05d15d4720cc209c02cea0270107c59af586add632f128fde00e16951e442a7aa1fafe6464ae63cd1797a9c6c0b24128d0ec61ec856b7ef04672036715eae4e2df35b3d30295df5cbbd71198d9ebb94918fe00eaf047edf1f0d5912b92bd2afd8911d303d691a7e5b873dab49e4f7b09b74d83a9025c18ac11f5926aa34a21c258cf49dd10909bec8c3d603648c2d86f90938ddc39470ee4ca98b1fb673951453599354f1d56d8f710c89d2bde5835a66b565df382d04c3e6eac3f98b52bb82f830a80d718e25f6adfcb93ec2b89cf2d56d610ce9238373b2cfab03bb1a00"
    }
  ]
}
```
Possible errors are: 
1. 404 - Input is invalid, there can't be any valid data for the input supplied
2. 500 - Something went wrong with cashAccountAPId (check STDOUT of the cashAccountAPId)

