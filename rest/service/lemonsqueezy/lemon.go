package lemonsqueezy

import "os"

type LemonSqueezy struct {
	apikey string
}

var Lemon LemonSqueezy = LemonSqueezy{
	apikey: os.Getenv("LEMONSQUEEZY_API_KEY"),
}
