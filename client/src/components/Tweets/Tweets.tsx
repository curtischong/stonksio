import React from 'react';
import styled from 'styled-components';
import Tweet from '../Tweet';

const TweetsContainer = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
`;

const mockData = [
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

const Tweets: React.FC = () => {
  const tweetMarkdown = mockData.map(({ name, msg, timestamp }) => 
    <Tweet name={name} msg={msg} timestamp={timestamp} />);
  return (
    <TweetsContainer>
      {tweetMarkdown}
    </TweetsContainer>
  );
}

export default Tweets;
