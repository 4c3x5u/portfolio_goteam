import React from 'react';
import { Button, Container } from 'react-bootstrap';
import './home.sass';
import logo from '../../assets/homeheader.svg';

const Home = () => (
  <div id="Home">
    <div className="Header">
      <img alt="logo" src={logo} />
    </div>
    <div className="ControlBar">
      <Container>
        <ul>
          <li>
            <Button
              className="Button"
              onClick={() => console.log('team button clicked')}
              aria-label="team controls"
            >
              TEAM
            </Button>
          </li>
          <li>
            <Button
              className="Button"
              onClick={() => console.log('boards button clicked')}
              aria-label="boards controls"
            >
              BOARDS
            </Button>
          </li>
          <li>
            <Button
              className="Button"
              onClick={() => console.log('help button clicked')}
              aria-label="help button"
            >
              HELP
            </Button>
          </li>
        </ul>
      </Container>
    </div>
  </div>
);

export default Home;
