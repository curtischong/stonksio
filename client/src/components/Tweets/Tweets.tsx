import React from 'react';
import styled from 'styled-components';

import Tweet from '../Tweet';

import { Tweet as TweetStruct } from '../../types';

const TweetsContainer = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
  overflow-y: auto;
  height: 640px;
`;

interface TweetsProps {
  tweets: TweetStruct[];
}

const Tweets: React.FC<TweetsProps> = ({ tweets }) => {
  const tweetMarkdown = tweets.map(({ name, msg, timestamp: ts }, idx) =>
    <Tweet key={idx} name={name} msg={msg} timestamp={`${ts.getHours()}:${ts.getMinutes().toString().padStart(2, "0")}`} />);

  return (
    <TweetsContainer>
      {tweetMarkdown}
    </TweetsContainer>
  );
}

export default Tweets;
