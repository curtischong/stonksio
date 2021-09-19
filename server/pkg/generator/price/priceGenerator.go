package price

import (
	"fmt"
	"math"
	"math/rand"
	"stonksio/pkg/config"
	"stonksio/pkg/database"
	"stonksio/pkg/sentiment"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	defaultLastPrice   = 3000.0
	stdDevForRandPrice = 3.0
)

type PriceGenerator struct {
	cockroachDbClient *database.CockroachDbClient
	gcpClient         sentiment.GcpClient
	lastPrice         float32
	globalSentiment   float32
	minuteSentiment   float32
	averageTime       int
	logger            *log.Logger
}

func NewPriceGenerator(
	config config.FeedConfig,
	cockroachDbClient *database.CockroachDbClient,
	gcpClient sentiment.GcpClient,
	asset string,
) *PriceGenerator {
	lastPrice, err := cockroachDbClient.GetLatestOhlc(asset)
	if err != nil {
		log.Errorf("couldn't find last price. err=%s", err)
		lastPrice = defaultLastPrice
	}

	return &PriceGenerator{
		cockroachDbClient: cockroachDbClient,
		gcpClient:         gcpClient,
		lastPrice:         lastPrice,
		logger:            log.New(),
		averageTime:       config.AverageTime,
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

		err := g.writeNewRandomPrice()
		if err != nil {
			log.Errorf("couldn't generate random price. err=%s", err)
		}
	}
}

func (g *PriceGenerator) getNewPrice() float32 {
	return g.lastPrice + g.globalSentiment*0.1 + g.minuteSentiment*0.1
}

func (g *PriceGenerator) writeNewRandomPrice() error {
	newPrice := g.getNewPrice() + float32(rand.NormFloat64()*stdDevForRandPrice)
	err := g.cockroachDbClient.InsertPrice("ETH", newPrice)
	if err != nil {
		return fmt.Errorf("cannot insert price err=%s", err)
	}
	return nil
}

func (g *PriceGenerator) writeNewPriceFromPostSentiment(
	postBody string,
) error {
	sentiment, err := g.gcpClient.CalculateSentiment(postBody)
	if err != nil {
		return err
	}
	g.globalSentiment += float32(sentiment.Score * 0.1)
	// since magnitude is from [0, inf) I'm using a sqrt so it doesn't explode
	// TODO: consider switching to log if need be
	g.minuteSentiment += float32(math.Sqrt(float64(sentiment.Magnitude)) * 0.1)

	err = g.cockroachDbClient.InsertPrice("ETH", g.getNewPrice())
	if err != nil {
		return fmt.Errorf("cannot insert price err=%s", err)
	}
	return nil
}
