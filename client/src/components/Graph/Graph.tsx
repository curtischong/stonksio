import React from 'react';
import styled from 'styled-components';
import { VictoryChart, VictoryLine, VictoryAxis, VictoryTooltip } from 'victory';

const GraphContainer = styled.div`
  position: relative;
  width: 100%;
`;

const Graph: React.FC = () => {
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
          data={[
            { x: 1, y: 2 },
            { x: 2, y: 3 },
            { x: 3, y: 5 },
            { x: 4, y: 4 },
            { x: 5, y: 60 },
            { x: 6, y: 20 },
            { x: 7, y: 30 },
            { x: 8, y: 50 },
            { x: 9, y: 40 },
            { x: 10, y: 100 }
          ]}
          animate={{
            duration: 2000,
          }}
        />
      </VictoryChart>
    </GraphContainer>
  );
}

export default Graph;
