module blockChaim

go 1.13

replace (
	golang.org/x/crypto v0.0.0-20180820150726-614d502a4dac => github.com/golang/crypto v0.0.0-20180820150726-614d502a4dac
	golang.org/x/net v0.0.0-20180821023952-922f4815f713 => github.com/golang/net v0.0.0-20180826012351-8a410e7b638d
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
)

require (
	github.com/boltdb/bolt v1.3.1
	golang.org/x/crypto v0.0.0-20180820150726-614d502a4dac
	golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c // indirect

)
