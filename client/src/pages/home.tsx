import React, { useEffect, useState } from 'react';
import styled from 'styled-components';

import Graph from '../components/Graph';
import Heading from '../components/Heading';
import TweetInput from '../components/TweetInput';
import Tweets from '../components/Tweets';

import { Tweet } from '../types';

import pusher from '../utils/pusher';

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

const mockData: any = [
  {
    name: "a",
    msg: "An essential piece of #Ethereum’s Serenity upgrade, the Beacon Chain’s deposit contract, is live. This begins a transition to #Eth2.",
    timestamp: "12:00"
  },
  {
    name: "b",
    msg: "More than 1000 hackers from around the world are staked and beginning to hack today at ETHOnline!",
    timestamp: "12:00"
  },
  {
    name: "c",
    msg: "We're days away from #ETHOnline—the biggest Ethereum event of the year! It's a hackathon with multiple single-day conferences on NFTs, DAOs, and The Merge.",
    timestamp: "12:00"
  },
  {
    name: "d",
    msg: "We're pleased to announce we've chosen SpruceID to lead the effort to standardize Sign-in with Ethereum!",
    timestamp: "12:00"
  }
];

const HomePage: React.FC = () => {
  const [tweets, setTweets] = useState([]);

  useEffect(() => {
    const getTweets = () => {
      setTweets(mockData);
    };

    const onTweetReceived = (tweet: Tweet) => {
      setTweets((prevTweets: Tweet[]): any => {
        return [...prevTweets, tweet];
      });
    };

    const setupPusher = () => {
      const channel = pusher().subscribe("tweets");
      channel.bind('newTweet', onTweetReceived);
    };

    getTweets();
    setupPusher();

    return (): void => {
      pusher().unbind('newTweet', onTweetReceived);
    };
  }, []);

  const submitTweet = (message: string) => {
    const newTweet: Tweet = {
      name: 'Daniel',
      msg: message,
      timestamp: Date.now().toString()
    }
    setTweets((prevTweets: Tweet[]): any => {
      return [newTweet, ...prevTweets];
    });
  };

  return (
    <GridContainer>
      <Sidebar>
        <Heading>
          Activity
        </Heading>
        <TweetInput onSubmit={submitTweet}/>
        <Line/>
        <Tweets tweets={tweets}/>
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
