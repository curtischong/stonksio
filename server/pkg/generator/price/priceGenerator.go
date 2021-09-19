package price

import (
	"math"
	"math/rand"
	"stonksio/pkg/common"
	"stonksio/pkg/database"
	"stonksio/pkg/sentiment"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	defaultAvgTime       = 2 // on avg, generate a new price every x seconds
	defaultStartingPrice = 3000.0
	stdDevForRandPrice   = 3.0
)

type PriceGenerator struct {
	cockroachDbClient *database.CockroachDbClient
	gcpClient         *sentiment.GcpClient
	out               chan<- *common.Price
	logger            *log.Logger
	asset             string
	lastPrice         float32
	sentimentPrice    float32
	minuteSentiment   float32
	globalMomentum    float32
	averageTime       int
}

func NewPriceGenerator(
	cockroachDbClient *database.CockroachDbClient,
	gcpClient *sentiment.GcpClient,
	postChannel chan<- *common.Price,
	asset string,
) *PriceGenerator {
	lastPrice, err := cockroachDbClient.GetLatestPrice(asset)
	if err != nil {
		log.Errorf("couldn't find last price. err=%s", err)
		lastPrice = defaultStartingPrice
	}

	return &PriceGenerator{
		cockroachDbClient: cockroachDbClient,
		gcpClient:         gcpClient,
		lastPrice:         lastPrice,
		out:               postChannel,
		asset:             asset,
		logger:            log.New(),
		averageTime:       defaultAvgTime,
		globalMomentum:    0,
		sentimentPrice:    3000,
	}
}

func (g *PriceGenerator) Start() {
	go g.generatePrices()
}

func (g *PriceGenerator) generatePrices() {
	defer func() {
		if err := recover(); err != nil {
			g.logger.Warnf("Price job died, restarting. err=%s", err)
			go g.generatePrices()
		}
	}()

	for {
		delay := time.Duration(rand.Intn(g.averageTime)) * time.Second
		time.Sleep(delay)

		g.out <- g.getNewPrice()
	}
}

func (g *PriceGenerator) GetNewPriceFromPostSentiment(
	postBody string,
) (*common.Price, error) {
	sentiment, err := g.gcpClient.CalculateSentiment(postBody)
	if err != nil {
		return nil, err
	}

	// in the future the coefficient we multiply the score by should be dependent on who says it
	g.sentimentPrice += sentiment.Score * 30

	sentimentDirection := 1.0
	if sentiment.Score < 0 {
		sentimentDirection = -1.0
	}

	// TODO: FIX BUG WITH SIGMOID
	// since magnitude is from [0, inf) I'm using a sigmoid so it doesn't explode
	// Thus, the minute sentiment will now be from [-20, 20]
	g.minuteSentiment += float32(sentimentDirection*float64(sentiment.Magnitude)) * 20
	return g.getNewPrice(), nil
}

func (g *PriceGenerator) sigmoid(x float32) float32 {
	return float32(1 / (1 + math.Exp(-float64(x))))
}

func (g *PriceGenerator) getNewPrice() *common.Price {
	// this is to bring is back to the default starting price
	// we are using bang-bang control to over-compensate
	targetPrice := g.sentimentPrice
	// if we don't math.pow this, the momentum changes too quickly
	diff := targetPrice - g.lastPrice
	momentumDirection := 1.0
	if diff < 0 {
		momentumDirection = -1.0
	}
	//g.globalMomentum += float32(momentumDirection*math.Pow(math.Abs(float64(diff)), 0.4)) * 0.2
	//g.globalMomentum += float32(momentumDirection * math.Log(math.Abs(float64(diff))))

	// is negative when abs(diff) is small
	absDiff := math.Abs(float64(diff))
	// if the diff is < 50, the slowMomentumTerm is negative
	slowMomentumTerm := float64(g.sigmoid(float32(absDiff+50))*10 - 5)
	// if the diff is < 10, the slowerMomentumTerm is negative
	slowerMomentumTerm := float64(g.sigmoid(float32(absDiff+10)) - 0.5)

	newMomentum := float64(g.globalMomentum) + momentumDirection*(math.Log(absDiff+1.1))
	if absDiff < 200 {
		if math.Abs(float64(g.globalMomentum)) > 20 {
			// we want to slow down the momentum
			newMomentum = float64(g.globalMomentum)*0.5 + slowMomentumTerm
		} else {
			newMomentum = float64(g.globalMomentum)*0.5 + slowerMomentumTerm
		}
	}
	newBoundedMomentum := math.Max(-100.0, math.Min(100.0, newMomentum))

	g.globalMomentum = float32(newBoundedMomentum)

	globalMomentumPrice := g.lastPrice + g.globalMomentum //float32(math.Pow(float64(g.globalMomentum), 1.0))

	noise := float32(rand.NormFloat64() * stdDevForRandPrice)
	newPrice := g.minuteSentiment + globalMomentumPrice + noise
	// the exponential function is always positive. By feeding it -newPrice we will trend up more the closer past 0 the price is
	// offset the exponential function to the right a bit
	newPrice = float32(math.Max(0, float64(newPrice)+3*math.Exp(float64(-newPrice-10.0)))) // bound the newPrice so it can't be lower than 0
	log.Infof("sentimentPrice: %f, globalMomentum: %f, diff: %f, price: %f", g.sentimentPrice, g.globalMomentum, diff, newPrice)
	g.lastPrice = newPrice
	return &common.Price{
		Asset:      g.asset,
		TradePrice: newPrice,
		Timestamp:  time.Now(),
	}
}
