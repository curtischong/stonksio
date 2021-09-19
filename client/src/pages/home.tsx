import axios from 'axios';
import React, { useEffect, useState } from 'react';
import styled from 'styled-components';

import Graph from '../components/Graph';
import Heading from '../components/Heading';
import Modal from '../components/Modal';
import TweetInput from '../components/TweetInput';
import Tweets from '../components/Tweets';

import { Tweet, TweetReponse } from '../types';

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

  const mapTweetResponse = (resp: TweetReponse): Tweet => {
    const ts = new Date(resp.Timestamp)
    return {
      name: resp.Username,
      msg: resp.Body,
      timestamp: `${ts.getHours()}:${ts.getMinutes()}`,
    }
  }

  useEffect(() => {
    const getTweets = () => {
      axios.get("https://stonk.st/api/posts", {
        params: {
          count: 20
        }
      }).then(resp => {
        setTweets(resp.data.map(mapTweetResponse))
      }).catch(err => console.error(err))
    };

    const onTweetReceived = (resp: TweetReponse) => {
      setTweets((prevTweets: Tweet[]): any => {
        return [...prevTweets, mapTweetResponse(resp)];
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
      timestamp: Date.now().toString()
    }
    setTweets((prevTweets: Tweet[]): any => {
      return [newTweet, ...prevTweets];
    });
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
