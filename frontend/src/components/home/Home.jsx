import React from 'react';
import { Container } from 'react-bootstrap';
import {
  faUsers,
  faChalkboardTeacher,
  faQuestionCircle,
} from '@fortawesome/free-solid-svg-icons';
import ControlButton from './ControlButton';
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
          <ControlButton
            name="team"
            action={() => console.log('team button clicked')}
            icon={faUsers}
          />
          <ControlButton
            name="boards"
            action={() => console.log('boards button clicked')}
            icon={faChalkboardTeacher}
          />
          <ControlButton
            name="help"
            action={() => console.log('boards button clicked')}
            icon={faQuestionCircle}
          />
        </ul>
      </Container>
    </div>
  </div>
);

export default Home;
