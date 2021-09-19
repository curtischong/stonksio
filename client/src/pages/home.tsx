import React, { useEffect, useState } from 'react';
import styled from 'styled-components';

import Graph from '../components/Graph';
import Heading from '../components/Heading';
import Modal from '../components/Modal';
import TweetInput from '../components/TweetInput';
import Tweets from '../components/Tweets';

import { Tweet, ServerTweet } from '../types';

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

const HomePage: React.FC = () => {
  const [tweets, setTweets] = useState([]);
  const [username, setUsername] = useState('');

  const mapFromTweetResponse = (resp: ServerTweet): Tweet => {
    return {
      name: resp.Username,
      msg: resp.Body,
      timestamp: new Date(resp.Timestamp),
    };
  }

  const mapToTweetRequest = (tweet: Tweet): ServerTweet => {
    return {
      Username: tweet.name,
      Body: tweet.msg,
      Timestamp: tweet.timestamp.toISOString(),
    };
  }

  useEffect(() => {
    const getTweets = () => {
      fetch("https://stonk.st/api/posts?count=20").then(resp => {
        return resp.json()
      }).then(json => {
        setTweets(json.map(mapFromTweetResponse));
      }).catch(err => console.error(err));
    };

    const onTweetReceived = (resp: ServerTweet) => {
      setTweets((prevTweets: Tweet[]): any => {
        return [mapFromTweetResponse(resp), ...prevTweets, mapFromTweetResponse(resp)];
      });
    };

    const setupPusher = () => {
      const channel = pusher().subscribe("post");
      channel.bind('new-post', onTweetReceived);
    };

    getTweets();
    setupPusher();

    return (): void => {
      pusher().unbind('new-post', onTweetReceived);
    };
  }, []);

  const submitTweet = (message: string) => {
    const newTweet: Tweet = {
      name: username,
      msg: message,
      timestamp: new Date()
    }
    fetch(
      "https://stonk.st/api/post",
      {
        method: "POST",
        headers: {
          'Content-Type': 'application/json'
        },
        mode: "no-cors",
        body: JSON.stringify(mapToTweetRequest(newTweet))
      }).then(resp => {}).catch(err => console.error(err));
  };

  const onClose = (newUsername: string) => {
    setUsername(newUsername);
  };

  return (
    <>
      {username === '' && <Modal onClose={onClose} />}
      <GridContainer>
        <Sidebar>
          <Heading>
            Activity
          </Heading>
          <TweetInput onSubmit={submitTweet} />
          <Line />
          <Tweets tweets={tweets} />
        </Sidebar>
        <Content>
          <Heading>Ethereum</Heading>
          <Price>US$1234.41</Price>
          <Graph />
        </Content>
      </GridContainer>
    </>
  );
};

export default HomePage;
