import React from 'react';
import styled from 'styled-components';

import Tweet from '../Tweet';

import { Tweet as TweetStruct } from '../../types';

const TweetsContainer = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
`;

interface TweetsProps {
  tweets: TweetStruct[];
}

const Tweets: React.FC<TweetsProps> = ({ tweets }) => {
  const tweetMarkdown = tweets.map(({ name, msg, timestamp }) => 
    <Tweet name={name} msg={msg} timestamp={timestamp} />);

  return (
    <TweetsContainer>
      {tweetMarkdown}
    </TweetsContainer>
  );
}

export default Tweets;
