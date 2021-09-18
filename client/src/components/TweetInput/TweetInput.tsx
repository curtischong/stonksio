import React from 'react';
import styled from 'styled-components';

const Input = styled.textarea`
  width: 100%;
  background-color: #25272d;
  border: none;
  overflow-y: auto;
  resize: none;
  height: 64px;
  max-height: 64px;
  border-radius: 4px;
  box-sizing: border-box;
  padding: 8px;
  color: #ffffff;
  ::placeholder,
  ::-webkit-input-placeholder {
    color: #474b52;
  }
  :-ms-input-placeholder {
    color: #474b52;
  }
  border-left: 4px solid #4E2A84;
`;

const TweetInput: React.FC = () => {
  return <Input placeholder="What's on your mind?"/>
}

export default TweetInput;
