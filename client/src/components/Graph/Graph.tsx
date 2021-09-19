import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { VictoryChart, VictoryLine, VictoryAxis, VictoryTooltip } from 'victory';
import { Price, ServerPrice } from '../../types';

import pusher from '../../utils/pusher';

const GraphContainer = styled.div`
  position: relative;
  width: 100%;
`;

const mapFromPriceResponse = (resp: ServerPrice): Price => {
  return {
    x: resp.Timestamp,
    y: parseFloat(resp.TradePrice)
  };
}

const Graph: React.FC = () => {
  const [prices, setPrices] = useState([]);

  useEffect(() => {
    const onPriceReceived = (resp: ServerPrice) => {
      setPrices((prevPrices: Price[]): any => {
        return [...prevPrices, mapFromPriceResponse(resp)];
      });
    };

    const getTweets = () => {
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

    getTweets();
    setupPusher();

    return (): void => {
      pusher().unbind('new-price', onPriceReceived);
    };
  }, []);

  return (
    <GraphContainer>
      <VictoryChart>
        <VictoryAxis 
          style={{
            axis: {
              stroke: 'white'
            },
            tickLabels: {
              fill: 'white'
            },
            grid: { stroke: "#818e99", strokeWidth: 0.5 }
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
            grid: { stroke: "#818e99", strokeWidth: 0.5 }
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
  );
}

export default Graph;
