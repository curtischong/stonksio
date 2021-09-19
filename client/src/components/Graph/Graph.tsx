import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { VictoryChart, VictoryLine, VictoryAxis, VictoryTooltip, VictoryArea } from 'victory';
import { Price, ServerPrice } from '../../types';

import pusher from '../../utils/pusher';

const GraphContainer = styled.div`
  position: relative;
  width: 100%;
`;

const PriceText = styled.div`
  font-size: 16px;
  font-weight: 600;
  color: #474b52;
`;


let max = 0;
let min = 10000000000;

const mapFromPriceResponse = (resp: ServerPrice): Price => {
  const newPrice: number = parseFloat(resp.TradePrice);

  max = Math.max(max, newPrice);
  min = Math.min(min, newPrice);

  return {
    x: new Date(resp.Timestamp),
    y: newPrice
  };
}

const Graph: React.FC = () => {
  const [prices, setPrices] = useState([]);
  const [currentPrice, setCurrentPrice] = useState(0.00);

  useEffect(() => {
    const onPriceReceived = (resp: ServerPrice) => {
      const newPrice: Price = mapFromPriceResponse(resp);
      setPrices((prevPrices: Price[]): any => {
        return [...prevPrices, newPrice];
      });
      setCurrentPrice(newPrice.y);
    };

    const getPrices = () => {
      fetch("https://stonk.st/api/prices/eth?window=5m").then(resp => {
        return resp.json()
      }).then(json => {
        setPrices(json.map(mapFromPriceResponse));
      }).catch(err => console.error(err));
    };

    const setupPusher = () => {
      const channel = pusher().subscribe("prices");
      channel.bind('new-price', onPriceReceived);
    };

    getPrices();
    setupPusher();

    return (): void => {
      pusher().unbind('new-price', onPriceReceived);
    };
  }, []);

  return (
    <>
      <svg style={{ height: 0 }}>
        <defs>
          <linearGradient id="gradient" x1="0" x2="0" y1="0" y2="1">
            <stop offset="0%" stopColor="#4E2A84"/>
            <stop offset="100%" stopColor="#1a1b20"/>
          </linearGradient>
        </defs>
      </svg>
      <PriceText>US${currentPrice}</PriceText>
      <GraphContainer>
        <VictoryChart
          scale={{ x: "time" }}
          maxDomain={{ 
            y: max + max/1000
          }}
          minDomain={{ 
            y: min - min/1000
          }}
        >
          <VictoryAxis 
            style={{
              axis: {
                stroke: 'white',
              },
              tickLabels: {
                fontSize: 8,
                padding: 5,
                fill: 'white'
              },
            }}
          />
          <VictoryAxis 
            dependentAxis
            style={{
              axis: {
                stroke: 'white',
              },
              tickLabels: {
                fontSize: 6,
                padding: 5,
                fill: 'white'
              },
            }}
          />
          <VictoryArea
            style={{
              data: { 
                stroke: "#4E2A84",
                strokeWidth: 0.5,
                fill: "url(#gradient)",
                fillOpacity: 0.5
              }
            }}
            data={prices}
            animate={{
              duration: 1000,
              onLoad: { 
                duration: 1000 
              }
            }}
          />
        </VictoryChart>
      </GraphContainer>
    </>
  );
}

export default Graph;
