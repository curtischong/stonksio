import React from 'react';
import styled from 'styled-components';

const TweetContainer = styled.div`
  display: flex;
  flex-direction: column;
  background-color: #16181c;
  margin-top: 8px;
  padding: 8px;
  border-radius: 4px;
  width: 100%;
  position: relative;
  box-sizing: border-box;
`;

const Username = styled.div`
  font-size: 12px;
  color: #474b52;
`;

const Message = styled.div`
  font-size: 16px;
  color: #ffffff;
  margin-top: 4px;
  padding-left: 8px;
  border-left: 1px solid #4E2A84;
`;

const Timestamp = styled.div`
  position: absolute;
  font-size: 12px;
  color: #474b52;
  right: 8px;
`;

interface TweetProps {
  name: string;
  msg: string;
  timestamp: string;
}

const Tweet: React.FC<TweetProps> = ({ name, msg, timestamp }) => {
  return (
    <TweetContainer>
      <Username>
        { name }
      </Username>
      <Message>
        { msg }
      </Message>
      <Timestamp>
        { timestamp }
      </Timestamp>
    </TweetContainer>
  );
};

export default Tweet;
