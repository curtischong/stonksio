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
	defaultAvgTime       = 1 // on avg, generate a new price every x seconds
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
	g.globalSentiment += float32(sentiment.Score * 0.1)
	// since magnitude is from [0, inf) I'm using a sqrt so it doesn't explode
	// TODO: consider switching to log if need be
	g.minuteSentiment += float32(math.Sqrt(float64(sentiment.Magnitude)) * 0.1)
	return g.getNewPrice(), nil
}

func (g *PriceGenerator) getNewPrice() *common.Price {
	sentimentPrice := g.globalSentiment*5 + g.minuteSentiment*5
	// this is to bring is back to the default starting price
	// we are using bang-bang control to over-compensate
	g.globalMomentum += defaultStartingPrice - g.lastPrice
	momentumPrice := g.globalMomentum*0.05 + float32(rand.NormFloat64()*6)

	noise := float32(rand.NormFloat64() * stdDevForRandPrice)
	newPrice := sentimentPrice + momentumPrice + noise
	// the exponential function is always positive. By feeding it -newPrice we will trend up more the closer past 0 the price is
	newPrice = float32(math.Max(0, float64(newPrice)+3*math.Exp(float64(-newPrice)))) // bound the newPrice so it can't be lower than 0
	log.Infof("globalSentiment: %f, globalMomentum: %f, price: %f", g.globalSentiment, g.globalMomentum, newPrice)
	g.lastPrice = newPrice
	return &common.Price{
		Asset:      g.asset,
		TradePrice: newPrice,
		Timestamp:  time.Now(),
	}
}
