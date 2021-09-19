import React, { useState } from 'react';
import styled from 'styled-components';
import Heading from '../Heading';

const ModalContainer = styled.div`
  display: flex;
  position: absolute;
  height: 100vh;
  width: 100%;
  background-color: rgba(0,0,0,0.8);
  z-index: 100;
  align-items: center;
  justify-content: center;
`;

const ModalContent = styled.div`
  display: flex;
  max-width: 400px;
  min-width: 300px;
  background-color: #1a1b20;
  border-radius: 4px;
  height: auto;
  padding: 32px;
  flex-direction: column;
`;

const SubText = styled.div`
  color: #474b52;
`;

const Input = styled.input`
  width: 100%;
  border-radius: 4px;
  border: none;
  margin-top: 8px;
  height: 32px;
  padding: 0 8px;
  box-sizing: border-box;
  background-color: #25272d;
  color: #ffffff;
`;

const Button = styled.button`
  width: 100%;
  height: 32px;
  margin-top: 8px;
  border-radius: 4px;
  background-color: white;
  font-weight: 600;
  cursor: pointer;
  font-size: 12px;
`;

interface ModalProps {
  onClose: Function;
}

const Modal: React.FC<ModalProps> = ({ onClose }) => {
  const [username, setUsername] = useState('');

  return (
    <ModalContainer>
      <ModalContent>
        <Heading>
          Welcome to stonk.st!
        </Heading>
        <SubText>
          Please enter a username to continue:
        </SubText>
        <Input placeholder="username" value={username} onChange={(e) => setUsername(e.target.value)}/>
        <Button onClick={() => onClose(username)}> Enter </Button>
      </ModalContent>
    </ModalContainer>
  );
};

export default Modal;
