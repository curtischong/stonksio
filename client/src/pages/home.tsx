import React from "react";
import styled from "styled-components";

import Graph from "../components/Graph";
import Heading from "../components/Heading";
import TweetInput from "../components/TweetInput";
import Tweets from "../components/Tweets";
import { initWebsockets } from "./websockets";

const Sidebar = styled.div`
  width: 100%;
`;

const Content = styled.div`
  width: 100%;
`;

const GridContainer = styled.div`
  display: grid;
  grid-column-gap: 40px;
  padding: 24px;
  grid-template-columns: 1fr 2fr;
`;

const Line = styled.div`
  height: 1px;
  width: 100%;
  background-color: #474b52;
  margin: 16px 0;
`;

const Price = styled.div`
  font-size: 16px;
  font-weight: 600;
  color: #474b52;
`;

const HomePage: React.FC = () => {
  initWebsockets();
  return (
    <GridContainer>
      <Sidebar>
        <Heading>Activity</Heading>
        <TweetInput />
        <Line />
        <Tweets />
      </Sidebar>
      <Content>
        <Heading>Ethereum</Heading>
        <Price>US$1234.41</Price>
        <Graph />
      </Content>
    </GridContainer>
  );
};

export default HomePage;
