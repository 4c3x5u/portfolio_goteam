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

import inboxHeader from '../../assets/inboxHeader.svg';
import readyHeader from '../../assets/readyHeader.svg';
import goHeader from '../../assets/goHeader.svg';
import doneHeader from '../../assets/doneHeader.svg';

import './home.sass';
import logo from '../../assets/homeHeader.svg';

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
      <Column name="inbox" header={inboxHeader} />
      <Column name="ready" header={readyHeader} />
      <Column name="go" header={goHeader} />
      <Column name="done" header={doneHeader} />
    </Row>

    <div className="Footer">
      <Container>
        <Row className="Row">
          <Col className="Col" xs={6}>
            <a href="https:FIX THIS">
              <FontAwesomeIcon icon={faGithub} />
              PROJECT
            </a>
          </Col>
          <Col className="Col" xs={6}>
            <a href="https:FIX THIS">
              <FontAwesomeIcon icon={faLinkedin} />
              AUTHOR
            </a>
          </Col>
        </Row>
      </Container>
    </div>
  </div>
);

export default Home;
