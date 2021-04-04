import React from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faUsers,
  faChalkboardTeacher,
  faQuestionCircle,
} from '@fortawesome/free-solid-svg-icons';
import { faGithub, faLinkedin } from '@fortawesome/free-brands-svg-icons';

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

    <div className="Footer">
      <Container>
        <Row className="Row">
          <Col className="Col" xs={4}>
            <a href="https:FIX THIS">
              <FontAwesomeIcon icon={faGithub} />
              <span>PROJECT</span>
            </a>
          </Col>
          <Col className="Col" xs={4}>
            <a href="https:FIX THIS">
              <FontAwesomeIcon icon={faLinkedin} />
              <span>AUTHOR</span>
            </a>
          </Col>
        </Row>
      </Container>
    </div>
  </div>
);

export default Home;
