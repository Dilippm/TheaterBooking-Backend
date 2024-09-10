// stripe.go
package config

import (
	
	"github.com/stripe/stripe-go"
)

func InitStripe() {
	stripe.Key = STRIPE_KEY
}
