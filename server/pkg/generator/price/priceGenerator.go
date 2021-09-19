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
	defaultAvgTime       = 3 // on avg, generate a new price every x seconds
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
	globalSentiment   float32
	minuteSentiment   float32
	globalMomentum    float32
	minuteMomentum    float32
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
		globalSentiment:   3000,
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
	g.globalSentiment += sentiment.Score * 30

	sentimentDirection := 1.0
	if sentiment.Score < 0 {
		sentimentDirection = -1.0
	}

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
	g.globalMomentum += g.globalSentiment - g.lastPrice
	globalMomentumPrice := g.globalMomentum * 0.05

	noise := float32(rand.NormFloat64() * stdDevForRandPrice)
	newPrice := g.minuteSentiment + globalMomentumPrice + noise
	// the exponential function is always positive. By feeding it -newPrice we will trend up more the closer past 0 the price is
	// offset the exponential function to the right a bit
	newPrice = float32(math.Max(0, float64(newPrice)+3*math.Exp(float64(-newPrice-10.0)))) // bound the newPrice so it can't be lower than 0
	log.Infof("globalSentiment: %f, globalMomentum: %f, price: %f", g.globalSentiment, g.globalMomentum, newPrice)
	g.lastPrice = newPrice
	return &common.Price{
		Asset:      g.asset,
		TradePrice: newPrice,
		Timestamp:  time.Now(),
	}
}
