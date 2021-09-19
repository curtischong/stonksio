package ohlc

import (
	"container/list"
	"stonksio/pkg/common"
	"stonksio/pkg/database"
	"stonksio/pkg/websocket"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const maxListSize = 60

type OHLCManager struct {
	lock              sync.RWMutex
	logger            *log.Logger
	ohlcs             *list.List
	cockroachDbClient *database.CockroachDbClient
	pusherClient      *websocket.PusherClient
}

func NewOHLCManager(
	cockroachDbClient *database.CockroachDbClient,
	pusherClient *websocket.PusherClient,
) *OHLCManager {
	m := &OHLCManager{
		logger:            log.New(),
		ohlcs:             list.New(),
		cockroachDbClient: cockroachDbClient,
		pusherClient:      pusherClient,
	}
	go m.init()
	return m
}

func (m *OHLCManager) init() {
	m.lock.Lock()
	defer m.lock.Unlock()
	ohlcs, err := m.cockroachDbClient.GetOHLCs(maxListSize)
	if err != nil {
		log.Errorf("error loading OHLCs, err=%s", err)
		return
	}

	for i := range ohlcs {
		m.ohlcs.PushFront(&ohlcs[i])
	}
}

func (m *OHLCManager) HandlePrice(price *common.Price) {
	m.lock.Lock()
	defer m.lock.Unlock()
	ohlc := m.ohlcs.Front().Value.(*common.OHLC)

	if ohlc.StartTime.Add(time.Minute).Before(price.Timestamp) {
		if m.ohlcs.Len() == maxListSize {
			m.ohlcs.Remove(m.ohlcs.Back())
		}

		ohlc = &common.OHLC{
			Open:      price.TradePrice,
			High:      price.TradePrice,
			Low:       price.TradePrice,
			Close:     price.TradePrice,
			StartTime: price.Timestamp.Truncate(time.Minute),
		}
		m.ohlcs.PushFront(ohlc)
	} else {
		if ohlc.Low > price.TradePrice {
			ohlc.Low = price.TradePrice
		}
		if ohlc.High < price.TradePrice {
			ohlc.High = price.TradePrice
		}
		ohlc.Close = price.TradePrice
	}

	m.pusherClient.PushOHLC(ohlc)
}

func (m *OHLCManager) GetOHLCs() []common.OHLC {
	m.lock.RLock()
	defer m.lock.RUnlock()
	ohlcs := make([]common.OHLC, 0, m.ohlcs.Len())

	// oldest first
	for e := m.ohlcs.Back(); e != m.ohlcs.Front(); e = e.Prev() {
		ohlcs = append(ohlcs, *e.Value.(*common.OHLC))
	}

	return ohlcs
}
