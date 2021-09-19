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
	defaultAvgTime     = 2 // on avg, generate a new price every x seconds
	defaultLastPrice   = 3000.0
	stdDevForRandPrice = 3.0
)

type PriceGenerator struct {
	cockroachDbClient *database.CockroachDbClient
	gcpClient         sentiment.GcpClient
	out               chan<- *common.Price
	logger            *log.Logger
	asset             string
	lastPrice         float32
	globalSentiment   float32
	minuteSentiment   float32
	averageTime       int
}

func NewPriceGenerator(
	cockroachDbClient *database.CockroachDbClient,
	gcpClient sentiment.GcpClient,
	postChannel chan<- *common.Price,
	asset string,
) *PriceGenerator {
	lastPrice, err := cockroachDbClient.GetLatestPrice(asset)
	if err != nil {
		log.Errorf("couldn't find last price. err=%s", err)
		lastPrice = defaultLastPrice
	}

	return &PriceGenerator{
		cockroachDbClient: cockroachDbClient,
		gcpClient:         gcpClient,
		lastPrice:         lastPrice,
		out:               postChannel,
		asset:             asset,
		logger:            log.New(),
		averageTime:       defaultAvgTime,
	}
}

func (g *PriceGenerator) Start() {
	defer func() {
		if err := recover(); err != nil {
			g.logger.Warnf("Price Start job died, restarting. err=%s", err)
			go g.Start()
		}
	}()

	for {
		delay := time.Duration(rand.Intn(g.averageTime)) * time.Second
		time.Sleep(delay)

		tradePrice := g.getNewRandomTradePrice()
		g.out <- &common.Price{
			Asset:      g.asset,
			TradePrice: tradePrice,
			Timestamp:  time.Now(),
		}
	}
}

func (g *PriceGenerator) getNewRandomTradePrice() float32 {
	return g.getNewPrice() + float32(rand.NormFloat64()*stdDevForRandPrice)
}

func (g *PriceGenerator) GetNewPriceFromPostSentiment(
	postBody string,
) (float32, error) {
	sentiment, err := g.gcpClient.CalculateSentiment(postBody)
	if err != nil {
		return 0, err
	}
	g.globalSentiment += float32(sentiment.Score * 0.1)
	// since magnitude is from [0, inf) I'm using a sqrt so it doesn't explode
	// TODO: consider switching to log if need be
	g.minuteSentiment += float32(math.Sqrt(float64(sentiment.Magnitude)) * 0.1)
	return g.getNewPrice(), nil
}

func (g *PriceGenerator) getNewPrice() float32 {
	return g.lastPrice + g.globalSentiment*0.1 + g.minuteSentiment*0.1
}
