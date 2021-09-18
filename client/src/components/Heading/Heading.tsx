import React from 'react';
import styled from 'styled-components';

const StyledHeading = styled.h1`
  font-size: 24px;
  color: #ffffff;
  text-align: left;
`;

const Heading: React.FC = ({ children }) => {
  return (
    <StyledHeading>
      { children }
    </StyledHeading>
  );
}

export default Heading;

