import React from 'react';
import { Container, Row } from 'react-bootstrap';
import {
  faUsers,
  faChalkboardTeacher,
  faQuestionCircle,
} from '@fortawesome/free-solid-svg-icons';

import ControlButton from './ControlButton';
import Column from './Column';
import columnName from './columnName';

import './home.sass';
import logo from '../../assets/homeheader.svg';

const Home = () => (
  <div id="Home">
    <div className="Logo">
      <img alt="logo" src={logo} />
    </div>
    <div className="ControlBar">
      <Container>
        <Row className="Controls">
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
        </Row>
      </Container>
    </div>

    <Row className="ColumnsRow">
      <Column name={columnName.INBOX} />
      <Column name={columnName.READY} />
      <Column name={columnName.GO} />
      <Column name={columnName.DONE} />
    </Row>
  </div>
);

export default Home;
