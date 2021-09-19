import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { VictoryChart, VictoryLine, VictoryAxis, VictoryTooltip } from 'victory';
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

const mapFromPriceResponse = (resp: ServerPrice): Price => {
  return {
    x: new Date(resp.Timestamp),
    y: parseFloat(resp.TradePrice)
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
      <PriceText>US${currentPrice}</PriceText>
      <GraphContainer>
        <VictoryChart
          scale={{ x: "time" }}
        >
          <VictoryAxis 
            style={{
              axis: {
                stroke: 'white'
              },
              tickLabels: {
                fill: 'white'
              },
            }}
          />
          <VictoryAxis 
            dependentAxis
            style={{
              axis: {
                stroke: 'white'
              },
              tickLabels: {
                fill: 'white'
              },
            }}
          />
          <VictoryLine
            labelComponent={
              <VictoryTooltip
                  constrainToVisibleArea
                  cornerRadius={0}
                  flyoutStyle={{
                    fill: "transparent",
                    strokeWidth: 0
                  }}
                  pointerLength={0}
                  style={{
                    fontSize: 16,
                    fill: "#ffffff",
                  }}
                />
            }
            style={{
              data: { stroke: "#c43a31" }
            }}
            data={prices}
            animate={{
              duration: 2000,
            }}
          />
        </VictoryChart>
      </GraphContainer>
    </>
  );
}

export default Graph;
