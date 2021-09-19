import React, { useState } from 'react';
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

const Button = styled.button`
  width: 100%;
  height: 32px;
  margin-top: 4px;
  border-radius: 4px;
  background-color: white;
  font-weight: 600;
  cursor: pointer;
  font-size: 12px;
`;

interface TweetInputProps {
  onSubmit: Function;
}

const TweetInput: React.FC<TweetInputProps> = ({ onSubmit }) => {
  const [value, setValue] = useState('');

  const handleSubmit = (event: any) => {
    event.preventDefault();
    onSubmit(value);
    setValue('');
  }

  const handleChange = (event: any) => {
    setValue(event.target.value);
  };

  return (
    <div>
      <Input 
        placeholder="What's on your mind?"
        value={value}
        onChange={handleChange}
        maxLength={300}
      />
      <Button onClick={handleSubmit}>
        Send
      </Button>
    </div>
  );
}

export default TweetInput;
